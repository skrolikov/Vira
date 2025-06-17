package types

// AuthRequest содержит данные для входа
// swagger:model AuthRequest
type AuthRequest struct {
	Username string `json:"username" example:"john_doe"`  // Имя пользователя
	Password string `json:"password" example:"secret123"` // Пароль
}

// LoginRequest содержит данные для логина
// swagger:model LoginRequest
type LoginRequest struct {
	Username string `json:"username" example:"john_doe"`  // Имя пользователя
	Password string `json:"password" example:"secret123"` // Пароль
}

// AuthResponse содержит токены и информацию о пользователе
// swagger:model AuthResponse
type AuthResponse struct {
	Tokens TokenPair `json:"tokens"` // Пара токенов доступа
	User   UserInfo  `json:"user"`   // Данные пользователя
}

// TokenPair содержит новую пару токенов доступа
// swagger:model TokenPair
type TokenPair struct {
	Access  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsIn..."`  // Access токен
	Refresh string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsIn..."` // Refresh токен
}

// UserInfo содержит основные данные пользователя
// swagger:model UserInfo
type UserInfo struct {
	ID       string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"` // Уникальный ID пользователя
	Username string `json:"username" example:"john_doe"`                       // Имя пользователя
	Role     string `json:"role" example:"user"`                               // Роль пользователя
}

// RefreshRequest содержит refresh токен для обновления access токена
// swagger:model RefreshRequest
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsIn..."` // Refresh токен
}

// RegisterRequest содержит данные для регистрации нового пользователя
// swagger:model RegisterRequest
type RegisterRequest struct {
	Username string `json:"username" example:"john_doe"`                // Имя пользователя
	Password string `json:"password" example:"secret123"`               // Пароль
	Email    string `json:"email,omitempty" example:"john@example.com"` // Email (опционально)
}

// ErrorResponse стандартная структура ошибки
// swagger:model ErrorResponse
type ErrorResponse struct {
	Error string `json:"error" example:"Неверный токен"` // Текст ошибки
}

// LogoutRequest содержит refresh токен сессии для выхода
// swagger:model LogoutRequest
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsIn..."` // Refresh токен
}
