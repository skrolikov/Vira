package types

// AuthRequest содержит данные для входа
type AuthRequest struct {
	Username string `json:"username" example:"john_doe"`  // Имя пользователя
	Password string `json:"password" example:"secret123"` // Пароль
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse содержит токены и информацию о пользователе
type AuthResponse struct {
	Tokens TokenPair `json:"tokens"`
	User   UserInfo  `json:"user"`
}

// TokenPair содержит новую пару токенов доступа
type TokenPair struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

// UserInfo содержит основные данные пользователя
type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// RefreshRequest содержит refresh токен для обновления access токена
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOi..."` // Refresh токен
}

// RegisterRequest содержит данные для регистрации нового пользователя
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email,omitempty"`
}

// ErrorResponse стандартная структура ошибки
type ErrorResponse struct {
	Error string `json:"error" example:"Неверный токен"`
}

// LogoutRequest содержит refresh токен сессии для выхода
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsIn..."` // Refresh токен
}
