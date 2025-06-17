package events

import (
	"context"
	"encoding/json"
	"time"

	kafka "github.com/skrolikov/vira-kafka"
	log "github.com/skrolikov/vira-logger"
)

// LoginEvent представляет событие входа пользователя в систему.
// Это событие публикуется в Kafka для последующей обработки.
//
// swagger:model LoginEvent
type LoginEvent struct {
	// Тип события
	// example: UserLoggedIn
	Type string `json:"type"`

	// Уникальный идентификатор пользователя
	// example: 123e4567-e89b-12d3-a456-426614174000
	UserID string `json:"user_id"`

	// Имя пользователя
	// example: sergei
	Username string `json:"username"`

	// IP-адрес пользователя при входе
	// example: 192.168.1.100
	IP string `json:"ip"`

	// User-Agent или устройство пользователя
	// example: Mozilla/5.0 (Windows NT 10.0; Win64; x64)
	UA string `json:"ua"`

	// Время события в формате RFC3339
	// example: 2025-06-12T14:35:00Z
	Time string `json:"time"`
}

// LoginSession описывает интерфейс для получения информации о сессии пользователя при входе.
type LoginSession interface {
	// GetIP возвращает IP-адрес пользователя.
	GetIP() string

	// GetDevice возвращает строку с информацией о устройстве пользователя (User-Agent).
	GetDevice() string
}

// EmitLoginEvent формирует и отправляет событие входа пользователя в Kafka.
//
// Параметры:
//   - ctx: контекст запроса для отмены/таймаута
//   - producer: Kafka-продюсер для отправки сообщения
//   - logger: логгер для записи предупреждений и ошибок
//   - userID: уникальный идентификатор пользователя
//   - username: имя пользователя
//   - s: сессия входа, реализующая интерфейс LoginSession
//
// В случае ошибки сериализации или отправки события записывает предупреждение в лог.
func EmitLoginEvent(ctx context.Context, producer *kafka.Producer, logger *log.Logger, userID, username string, s LoginSession) {
	event := LoginEvent{
		Type:     "UserLoggedIn",
		UserID:   userID,
		Username: username,
		IP:       s.GetIP(),
		UA:       s.GetDevice(),
		Time:     time.Now().Format(time.RFC3339),
	}

	data, err := json.Marshal(event)
	if err != nil {
		logger.Warn("Ошибка маршалинга события входа: %v", err)
		return
	}

	if err := producer.Send(ctx, userID, data); err != nil {
		logger.Warn("Ошибка отправки Kafka события входа: %v", err)
	}
}
