package types

import "time"

// SessionsResponse содержит список сессий с курсором для пагинации
// swagger:model SessionsResponse
type SessionsResponse struct {
	Cursor   uint64        `json:"cursor" example:"0"` // Курсор для следующей страницы (0 — начало)
	Sessions []SessionInfo `json:"sessions"`           // Список сессий пользователя
}

// SessionInfo содержит данные о конкретной сессии пользователя
// swagger:model SessionInfo
type SessionInfo struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`          // Уникальный ID сессии
	UserID    string    `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`     // ID пользователя
	Token     string    `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`    // Токен сессии (может быть обрезан)
	IP        string    `json:"ip" example:"192.168.1.10"`                                  // IP адрес сессии
	Device    string    `json:"device" example:"Mozilla/5.0 (Windows NT 10.0; Win64; x64)"` // Информация о устройстве / User-Agent
	LoginTime time.Time `json:"login_time" example:"2025-06-12T14:22:35Z"`                  // Время входа в сессию (ISO 8601)
}

func (s SessionInfo) GetIP() string {
	return s.IP
}

func (s SessionInfo) GetDevice() string {
	return s.Device
}
