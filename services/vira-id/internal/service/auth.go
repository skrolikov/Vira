package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"vira-id/internal/events"
	"vira-id/internal/types"

	"github.com/redis/go-redis/v9"
	config "github.com/skrolikov/vira-config"
	db "github.com/skrolikov/vira-db"
	hash "github.com/skrolikov/vira-hash"
	kafka "github.com/skrolikov/vira-kafka"
	log "github.com/skrolikov/vira-logger"
)

// AuthService — сервис аутентификации и регистрации пользователей.
// Инкапсулирует логику работы с пользователями, хранение сессий,
// генерацию токенов и отправку событий.
type AuthService struct {
	Cfg      *config.Config    // Конфигурация приложения
	Repo     db.UserRepository // Репозиторий пользователей (интерфейс к БД)
	Redis    *redis.Client     // Клиент Redis для хранения сессий
	Producer *kafka.Producer   // Kafka-продюсер для отправки событий
	Logger   *log.Logger       // Логгер для записи логов
}

// NewAuthService — конструктор для AuthService, инициализирует поля.
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

// Register — регистрация нового пользователя.
// Выполняет проверку данных, хеширует пароль, создаёт пользователя в БД,
// генерирует токены, сохраняет сессию в Redis и отправляет событие регистрации.
func (s *AuthService) Register(
	ctx context.Context,
	req types.RegisterRequest,
	ip, userAgent string,
) (*types.AuthResponse, error) {
	// Валидация данных регистрации (имплементировать validateRegistration)
	if err := s.validateRegistration(req); err != nil {
		return nil, err
	}

	// Проверяем уникальность username
	exists, err := s.Repo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки существующего пользователя: %w", err)
	}
	if exists {
		return nil, db.ErrDuplicateUsername
	}

	// Если указан email — проверяем уникальность email
	if req.Email != "" {
		existsEmail, err := s.Repo.ExistsByEmail(req.Email)
		if err != nil {
			return nil, fmt.Errorf("ошибка проверки email: %w", err)
		}
		if existsEmail {
			return nil, db.ErrDuplicateEmail
		}
	}

	// Хешируем пароль
	hashedPass, err := hash.HashPassword(req.Password)
	if err != nil {
		s.Logger.Error("Ошибка хеширования пароля: %v", err)
		return nil, fmt.Errorf("ошибка сервера: %w", err)
	}

	// Генерация токена подтверждения (нужно реализовать generateConfirmToken)
	confirmToken := generateConfirmToken()

	// Создаём пользователя с ролью user, не подтверждённого, с токеном подтверждения
	userID, err := s.Repo.CreateUserExtended(req.Username, hashedPass, req.Email, "user", false, confirmToken)
	if err != nil {
		return nil, err
	}

	// TODO: Отправить email с confirmToken (интеграция почты)

	// Генерируем JWT или другие токены доступа и обновления
	tokens, err := s.generateTokens(userID)
	if err != nil {
		return nil, err
	}

	// Сохраняем сессию с refresh-токеном, IP и User-Agent в Redis
	_, err = s.saveSession(ctx, userID, tokens.Refresh, ip, userAgent)
	if err != nil {
		return nil, err
	}

	// Асинхронно отправляем событие регистрации в Kafka
	go s.emitRegistrationEvent(ctx, userID, req.Username, ip, userAgent)

	// Возвращаем токены и базовую информацию о пользователе
	return &types.AuthResponse{
		Tokens: *tokens,
		User: types.UserInfo{
			ID:       userID,
			Username: req.Username,
			Role:     "user",
		},
	}, nil
}

// ConfirmUser — подтверждение пользователя по email и токену подтверждения.
// Вызывает метод репозитория для подтверждения и обрабатывает ошибки.
func (s *AuthService) ConfirmUser(ctx context.Context, email, token string) error {
	err := s.Repo.ConfirmUser(email, token)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return errors.New("пользователь не найден или уже подтвержден")
		}
		return fmt.Errorf("ошибка подтверждения пользователя: %w", err)
	}
	return nil
}

// Login — аутентификация пользователя по username и паролю.
// Проверяет пользователя, пароль, подтверждение, генерирует токены,
// сохраняет сессию и отправляет событие входа.
func (s *AuthService) Login(
	ctx context.Context,
	req types.LoginRequest,
	ip, userAgent string,
) (*types.AuthResponse, error) {
	user, err := s.Repo.GetUserByUsername(req.Username)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return nil, errors.New("пользователь не найден")
		}
		return nil, err
	}

	if !hash.CheckPasswordHash(user.PasswordHash, req.Password) {
		return nil, errors.New("неверный пароль")
	}

	if !user.Confirmed {
		return nil, errors.New("пользователь не подтверждён. Проверьте почту")
	}

	tokens, err := s.generateTokens(user.ID)
	if err != nil {
		return nil, err
	}

	_, err = s.saveSession(ctx, user.ID, tokens.Refresh, ip, userAgent)
	if err != nil {
		return nil, err
	}

	// Обновляем время последнего входа
	user.LastLoginAt = sql.NullTime{Time: time.Now(), Valid: true}
	if err := s.Repo.UpdateUser(user); err != nil {
		s.Logger.Error("Ошибка обновления времени последнего входа: %v", err)
	}

	// Отправляем событие входа асинхронно
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

// ChangePassword — смена пароля пользователя.
// Проверяет старый пароль, хеширует новый, обновляет в БД.
func (s *AuthService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if !hash.CheckPasswordHash(user.PasswordHash, oldPassword) {
		return errors.New("старый пароль неверный")
	}

	newHash, err := hash.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.Repo.UpdatePassword(userID, newHash); err != nil {
		return err
	}

	return nil
}
