package events

import (
	"context"
	"encoding/json"
	"time"

	kafka "github.com/skrolikov/vira-kafka"
	log "github.com/skrolikov/vira-logger"
)

// RefreshEvent представляет событие обновления токена пользователя.
//
// swagger:model RefreshEvent
type RefreshEvent struct {
	// Тип события
	// example: UserTokenRefreshed
	Type string `json:"type"`

	// Уникальный идентификатор пользователя
	// example: 123e4567-e89b-12d3-a456-426614174000
	UserID string `json:"user_id"`

	// IP-адрес пользователя при обновлении токена
	// example: 192.168.1.100
	IP string `json:"ip"`

	// User-Agent или устройство пользователя
	// example: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)
	UA string `json:"ua"`

	// Время события в формате RFC3339
	// example: 2025-06-12T14:40:00Z
	Time string `json:"time"`
}

// EmitRefreshEvent формирует и отправляет событие обновления токена пользователя в Kafka.
//
// Параметры:
//   - ctx: контекст запроса для отмены или таймаута
//   - producer: Kafka-продюсер для отправки сообщения
//   - logger: логгер для записи предупреждений и ошибок
//   - userID: уникальный идентификатор пользователя
//   - session: сессия пользователя, реализующая интерфейс LoginSession для получения IP и User-Agent
//
// В случае ошибки сериализации или отправки события логирует предупреждение.
func EmitRefreshEvent(ctx context.Context, producer *kafka.Producer, logger *log.Logger, userID string, session LoginSession) {
	event := RefreshEvent{
		Type:   "UserTokenRefreshed",
		UserID: userID,
		IP:     session.GetIP(),
		UA:     session.GetDevice(),
		Time:   time.Now().Format(time.RFC3339),
	}

	data, err := json.Marshal(event)
	if err != nil {
		logger.Warn("Ошибка маршалинга события refresh: %v", err)
		return
	}

	if err := producer.Send(ctx, userID, data); err != nil {
		logger.Warn("Ошибка отправки Kafka события refresh: %v", err)
	}
}
