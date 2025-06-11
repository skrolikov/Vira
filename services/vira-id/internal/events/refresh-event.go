package events

import (
	"context"
	"encoding/json"
	"time"

	kafka "github.com/skrolikov/vira-kafka"
	log "github.com/skrolikov/vira-logger"
)

type RefreshEvent struct {
	Type   string `json:"type"`
	UserID string `json:"user_id"`
	IP     string `json:"ip"`
	UA     string `json:"ua"`
	Time   string `json:"time"`
}

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
