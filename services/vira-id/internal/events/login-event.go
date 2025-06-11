package events

import (
	"context"
	"encoding/json"
	"time"

	kafka "github.com/skrolikov/vira-kafka"
	log "github.com/skrolikov/vira-logger"
)

type LoginEvent struct {
	Type     string `json:"type"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IP       string `json:"ip"`
	UA       string `json:"ua"`
	Time     string `json:"time"`
}

type LoginSession interface {
	GetIP() string
	GetDevice() string
}

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
