package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"vira-id/internal/events"
	"vira-id/internal/types"
	"vira-id/internal/validators"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	config "github.com/skrolikov/vira-config"
	db "github.com/skrolikov/vira-db"
	hash "github.com/skrolikov/vira-hash"
	jwt "github.com/skrolikov/vira-jwt"
	kafka "github.com/skrolikov/vira-kafka"
	log "github.com/skrolikov/vira-logger"
)

// AuthService сервис аутентификации и регистрации
type AuthService struct {
	Cfg      *config.Config
	Repo     db.UserRepository
	Redis    *redis.Client
	Producer *kafka.Producer
	Logger   *log.Logger
}

func NewAuthService(
	cfg *config.Config,
	repo db.UserRepository,
	rdb *redis.Client,
	producer *kafka.Producer,
	logger *log.Logger,
) *AuthService {
	return &AuthService{
		Cfg:      cfg,
		Repo:     repo,
		Redis:    rdb,
		Producer: producer,
		Logger:   logger,
	}
}

// Register регистрирует пользователя
func (s *AuthService) Register(
	ctx context.Context,
	req types.RegisterRequest,
	ip, userAgent string,
) (*types.AuthResponse, error) {
	if err := s.validateRegistration(req); err != nil {
		return nil, err
	}

	hashedPass, err := hash.HashPassword(req.Password)
	if err != nil {
		s.Logger.Error("Ошибка хеширования пароля: %v", err)
		return nil, fmt.Errorf("ошибка сервера: %w", err)
	}

	userID, err := s.createUser(ctx, req, hashedPass)
	if err != nil {
		return nil, err
	}

	tokens, err := s.generateTokens(userID)
	if err != nil {
		return nil, err
	}

	// Создаем новую сессию
	_, err = s.saveSession(ctx, userID, tokens.Refresh, ip, userAgent)
	if err != nil {
		return nil, err
	}

	// Отправка события регистрации в Kafka асинхронно
	go s.emitRegistrationEvent(ctx, userID, req.Username, ip, userAgent)

	return &types.AuthResponse{
		Tokens: *tokens,
		User: types.UserInfo{
			ID:       userID,
			Username: req.Username,
			Role:     "user",
		},
	}, nil
}

func (s *AuthService) validateRegistration(req types.RegisterRequest) error {
	if err := validators.ValidateCredentials(req.Username, req.Password); err != nil {
		return fmt.Errorf("невалидные учетные данные: %w", err)
	}

	if req.Email != "" {
		if err := validators.ValidateEmail(req.Email); err != nil {
			return fmt.Errorf("невалидный email: %w", err)
		}
	}

	return nil
}

func (s *AuthService) createUser(ctx context.Context, req types.RegisterRequest, hashedPass string) (string, error) {
	confirmToken := uuid.NewString()
	userID, err := s.Repo.CreateUserExtended(
		req.Username,
		hashedPass,
		req.Email,
		"user",
		false,
		confirmToken,
	)
	if err != nil {
		if errors.Is(err, db.ErrUserExists) {
			return "", errors.New("пользователь с таким именем или email уже существует")
		}
		s.Logger.Error("Ошибка базы данных: %v", err)
		return "", fmt.Errorf("ошибка при создании пользователя: %w", err)
	}
	return userID, nil
}

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

func (s *AuthService) saveSession(ctx context.Context, userID, refreshToken, ip, userAgent string) (string, error) {
	sessionID := uuid.NewString()
	sessionKey := "session:" + userID + ":" + sessionID
	refreshKey := "refresh:" + refreshToken

	session := types.SessionInfo{
		ID:        sessionID,
		UserID:    userID,
		Token:     refreshToken,
		IP:        ip,
		Device:    userAgent,
		LoginTime: time.Now(),
	}

	sessionData, err := json.Marshal(session)
	if err != nil {
		s.Logger.Error("Ошибка сериализации сессии: %v", err)
		return "", fmt.Errorf("ошибка сервера")
	}

	pipe := s.Redis.Pipeline()
	pipe.Set(ctx, sessionKey, sessionData, s.Cfg.JwtRefreshTTL)
	pipe.Set(ctx, refreshKey, sessionKey, s.Cfg.JwtRefreshTTL)
	_, err = pipe.Exec(ctx)
	if err != nil {
		s.Logger.Error("Ошибка сохранения сессии и индекса refresh в Redis: %v", err)
		return "", fmt.Errorf("ошибка сохранения сессии: %w", err)
	}

	return sessionID, nil
}

func (s *AuthService) GetUserSessions(ctx context.Context, userID string) ([]types.SessionInfo, error) {
	var sessions []types.SessionInfo

	var cursor uint64
	for {
		// Ищем ключи сессий по паттерну
		keys, nextCursor, err := s.Redis.Scan(ctx, cursor, "session:"+userID+":*", 100).Result()
		if err != nil {
			return nil, fmt.Errorf("ошибка поиска сессий: %w", err)
		}
		cursor = nextCursor

		for _, key := range keys {
			data, err := s.Redis.Get(ctx, key).Result()
			if err != nil {
				if err == redis.Nil {
					continue
				}
				return nil, fmt.Errorf("ошибка чтения сессии %s: %w", key, err)
			}

			var session types.SessionInfo
			if err := json.Unmarshal([]byte(data), &session); err != nil {
				s.Logger.Error("Ошибка разбора сессии из Redis: %v", err)
				continue
			}
			sessions = append(sessions, session)
		}

		if cursor == 0 {
			break
		}
	}

	return sessions, nil
}

func (s *AuthService) emitRegistrationEvent(ctx context.Context, userID, username, ip, device string) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := events.EmitUserRegisteredEvent(ctxTimeout, s.Producer, s.Logger, userID, username, ip, device, time.Now())
	if err != nil {
		s.Logger.Error("Ошибка отправки события регистрации: %v", err)
	}
}

// Login авторизует пользователя
func (s *AuthService) Login(
	ctx context.Context,
	req types.LoginRequest,
	ip, userAgent string,
) (*types.AuthResponse, error) {
	user, err := s.Repo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	if !hash.CheckPasswordHash(user.PasswordHash, req.Password) {
		return nil, errors.New("неверный пароль")
	}

	tokens, err := s.generateTokens(user.ID)
	if err != nil {
		return nil, err
	}

	// Создаем новую сессию
	_, err = s.saveSession(ctx, user.ID, tokens.Refresh, ip, userAgent)
	if err != nil {
		return nil, err
	}

	// Отправка события входа
	go events.EmitUserLoggedInEvent(ctx, s.Producer, s.Logger, user.ID, user.Username, ip, userAgent)

	return &types.AuthResponse{
		Tokens: *tokens,
		User: types.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Role:     user.Role,
		},
	}, nil
}

// RefreshToken обновляет access и refresh токены, создает новую сессию и удаляет старую.
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken, ip, userAgent string) (*types.TokenPair, error) {
	sessionKey, err := s.Redis.Get(ctx, "refresh:"+refreshToken).Result()
	if err == redis.Nil {
		return nil, errors.New("сессия не найдена")
	} else if err != nil {
		s.Logger.Error("Ошибка Redis при получении сессии по refresh токену: %v", err)
		return nil, fmt.Errorf("ошибка Redis: %w", err)
	}

	data, err := s.Redis.Get(ctx, sessionKey).Result()
	if err != nil {
		s.Logger.Error("Ошибка Redis при получении сессии: %v", err)
		return nil, fmt.Errorf("ошибка Redis: %w", err)
	}

	var oldSession types.SessionInfo
	if err := json.Unmarshal([]byte(data), &oldSession); err != nil {
		s.Logger.Error("Ошибка чтения сессии: %v", err)
		return nil, errors.New("ошибка чтения сессии")
	}

	// Проверяем валидность токена (через jwt.ParseToken) и userID
	claims, err := jwt.ParseToken(refreshToken, s.Cfg.JwtSecret)
	if err != nil || !jwt.IsTokenType(claims, "refresh") {
		return nil, errors.New("некорректный refresh токен")
	}
	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" || userID != oldSession.UserID {
		return nil, errors.New("некорректный user_id в токене")
	}

	// Генерируем новые токены
	tokens, err := s.generateTokens(userID)
	if err != nil {
		return nil, err
	}

	// Удаляем старую сессию и refreshKey
	pipe := s.Redis.Pipeline()
	pipe.Del(ctx, sessionKey)
	pipe.Del(ctx, "refresh:"+refreshToken)
	_, err = pipe.Exec(ctx)
	if err != nil {
		s.Logger.Warn("Ошибка удаления старой сессии: %v", err)
	}

	// Создаем новую сессию
	_, err = s.saveSession(ctx, userID, tokens.Refresh, ip, userAgent)
	if err != nil {
		return nil, err
	}

	// Отправка Kafka-события обновления токена
	go events.EmitRefreshEvent(ctx, s.Producer, s.Logger, userID, oldSession)

	return tokens, nil
}
