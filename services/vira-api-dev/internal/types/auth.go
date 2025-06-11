package types

import (
	"context"
	"errors"
	"time"
)

// Токены, получаемые от vira-id
type TokenPair struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

// Ответ от vira-id (обязательно соответствует тому, что возвращает сервис vira-id)
type AuthResponse struct {
	Tokens TokenPair `json:"tokens"`
	User   struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Role     string `json:"role"`
	} `json:"user"`
}

// Запросы, проксируемые в vira-id
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email,omitempty"`
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Локальный профиль для vira-dev
type UserProfile struct {
	UserID   string    `json:"user_id"`
	City     string    `json:"city"`
	JoinedAt time.Time `json:"joined_at"`
}

// Данные, передаваемые из клиента доп. к RegisterRequest
type ProfileData struct {
	City string `json:"city"`
}

// Ошибка «профиль не найден»
var ErrProfileNotFound = errors.New("profile not found")

// Ответ для vira-dev
type DevAuthResponse struct {
	Tokens  TokenPair   `json:"tokens"`
	Profile UserProfile `json:"profile"`
}

// Интерфейс репозитория профилей
type UserProfileRepository interface {
	Exists(ctx context.Context, userID string) (bool, error)
	Create(ctx context.Context, profile UserProfile) error
	GetByUserID(ctx context.Context, userID string) (UserProfile, error)
}
