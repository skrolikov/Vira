package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"vira-id/internal/types"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// saveSession сохраняет информацию о сессии пользователя в Redis.
// В Redis создается два ключа:
// 1) session:{userID}:{sessionID} — сериализованная структура сессии с информацией (IP, устройство, время и т.д.),
// 2) refresh:{refreshToken} — содержит ссылку на sessionKey для быстрого поиска сессии по refresh токену.
// TTL для обоих ключей берется из конфигурации (время жизни refresh токена).
func (s *AuthService) saveSession(ctx context.Context, userID, refreshToken, ip, userAgent string) (string, error) {
	// Генерируем уникальный ID сессии
	sessionID := uuid.NewString()
	// Формируем ключ для сессии
	sessionKey := "session:" + userID + ":" + sessionID
	// Формируем ключ для индексации refresh токена
	refreshKey := "refresh:" + refreshToken

	// Создаем структуру сессии для сохранения
	session := types.SessionInfo{
		ID:        sessionID,
		UserID:    userID,
		Token:     refreshToken,
		IP:        ip,
		Device:    userAgent,
		LoginTime: time.Now(),
	}

	// Сериализуем структуру сессии в JSON
	sessionData, err := json.Marshal(session)
	if err != nil {
		s.Logger.Error("Ошибка сериализации сессии: %v", err)
		return "", fmt.Errorf("ошибка сервера")
	}

	// Создаем pipeline для пакетной отправки нескольких команд в Redis
	pipe := s.Redis.Pipeline()
	pipe.Set(ctx, sessionKey, sessionData, s.Cfg.JwtRefreshTTL)
	pipe.Set(ctx, refreshKey, sessionKey, s.Cfg.JwtRefreshTTL)

	// Выполняем команды
	_, err = pipe.Exec(ctx)
	if err != nil {
		s.Logger.Error("Ошибка сохранения сессии и индекса refresh в Redis: %v", err)
		return "", fmt.Errorf("ошибка сохранения сессии: %w", err)
	}

	// Возвращаем ID сессии
	return sessionID, nil
}

// GetUserSessions возвращает все активные сессии пользователя, которые хранятся в Redis.
// Использует команду SCAN с паттерном session:{userID}:* для поиска ключей сессий.
// По каждому найденному ключу получает данные, десериализует и добавляет в результат.
func (s *AuthService) GetUserSessions(ctx context.Context, userID string) ([]types.SessionInfo, error) {
	var sessions []types.SessionInfo
	var cursor uint64

	for {
		// Ищем ключи сессий по паттерну с помощью SCAN (пакетами по 100)
		keys, nextCursor, err := s.Redis.Scan(ctx, cursor, "session:"+userID+":*", 100).Result()
		if err != nil {
			return nil, fmt.Errorf("ошибка поиска сессий: %w", err)
		}
		cursor = nextCursor

		for _, key := range keys {
			data, err := s.Redis.Get(ctx, key).Result()
			if err != nil {
				if err == redis.Nil {
					// Ключ мог быть удален после SCAN, пропускаем
					continue
				}
				return nil, fmt.Errorf("ошибка чтения сессии %s: %w", key, err)
			}

			var session types.SessionInfo
			// Преобразуем JSON в структуру
			if err := json.Unmarshal([]byte(data), &session); err != nil {
				s.Logger.Error("Ошибка разбора сессии из Redis: %v", err)
				continue
			}
			sessions = append(sessions, session)
		}

		// Если курсор 0 — обход завершён
		if cursor == 0 {
			break
		}
	}

	return sessions, nil
}
