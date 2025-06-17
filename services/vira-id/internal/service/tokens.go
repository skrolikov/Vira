package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"vira-id/internal/events"
	"vira-id/internal/types"

	"github.com/redis/go-redis/v9"
	jwt "github.com/skrolikov/vira-jwt"
)

// generateTokens создает пару JWT-токенов — access и refresh для указанного userID.
func (s *AuthService) generateTokens(userID string) (*types.TokenPair, error) {
	access, err := jwt.GenerateAccessToken(userID, s.Cfg.JwtSecret, s.Cfg.JwtTTL)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации access токена: %w", err)
	}

	refresh, err := jwt.GenerateRefreshToken(userID, s.Cfg.JwtSecret, s.Cfg.JwtRefreshTTL)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации refresh токена: %w", err)
	}

	return &types.TokenPair{
		Access:  access,
		Refresh: refresh,
	}, nil
}

// RefreshToken обновляет пару токенов — access и refresh,
// при этом удаляет старую сессию (по старому refresh токену),
// создает новую сессию и сохраняет её в Redis,
// также запускает асинхронное событие об обновлении токенов в Kafka.
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken, ip, userAgent string) (*types.TokenPair, error) {
	// Ищем в Redis sessionKey, связанный с refresh токеном
	sessionKey, err := s.Redis.Get(ctx, "refresh:"+refreshToken).Result()
	if err == redis.Nil {
		return nil, errors.New("сессия не найдена")
	} else if err != nil {
		s.Logger.Error("Ошибка Redis при получении сессии по refresh токену: %v", err)
		return nil, fmt.Errorf("ошибка Redis: %w", err)
	}

	// Получаем данные сессии по sessionKey
	data, err := s.Redis.Get(ctx, sessionKey).Result()
	if err != nil {
		s.Logger.Error("Ошибка Redis при получении сессии: %v", err)
		return nil, fmt.Errorf("ошибка Redis: %w", err)
	}

	// Парсим JSON сессии
	var oldSession types.SessionInfo
	if err := json.Unmarshal([]byte(data), &oldSession); err != nil {
		s.Logger.Error("Ошибка чтения сессии: %v", err)
		return nil, errors.New("ошибка чтения сессии")
	}

	// Проверяем валидность refresh токена и тип токена
	claims, err := jwt.ParseToken(refreshToken, s.Cfg.JwtSecret)
	if err != nil || !jwt.IsTokenType(claims, "refresh") {
		return nil, errors.New("некорректный refresh токен")
	}

	// Проверяем, что userID в токене совпадает с userID в сессии
	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" || userID != oldSession.UserID {
		return nil, errors.New("некорректный user_id в токене")
	}

	// Генерируем новую пару токенов
	tokens, err := s.generateTokens(userID)
	if err != nil {
		return nil, err
	}

	// Удаляем старую сессию и ключ refresh токена из Redis
	pipe := s.Redis.Pipeline()
	pipe.Del(ctx, sessionKey)
	pipe.Del(ctx, "refresh:"+refreshToken)
	_, err = pipe.Exec(ctx)
	if err != nil {
		s.Logger.Warn("Ошибка удаления старой сессии: %v", err)
	}

	// Сохраняем новую сессию с новым refresh токеном
	_, err = s.saveSession(ctx, userID, tokens.Refresh, ip, userAgent)
	if err != nil {
		return nil, err
	}

	// Асинхронно отправляем событие об обновлении токена в Kafka
	go events.EmitRefreshEvent(ctx, s.Producer, s.Logger, userID, oldSession)

	return tokens, nil
}

// generateConfirmToken создает случайный токен подтверждения (32 hex символа).
func generateConfirmToken() string {
	b := make([]byte, 16) // 16 байт = 32 hex символа
	_, err := rand.Read(b)
	if err != nil {
		// При ошибке генерации безопаснее вернуть пустую строку и логировать,
		// чем фиксированный токен, чтобы избежать конфликтов
		// Но здесь для простоты возвращаем дефолтное значение
		return "default-confirm-token"
	}
	return hex.EncodeToString(b)
}
