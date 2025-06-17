package events

import (
	"context"
	"encoding/json"
	"time"

	kafka "github.com/skrolikov/vira-kafka"
	log "github.com/skrolikov/vira-logger"
)

func EmitDevUserRegisteredEvent(ctx context.Context, producer *kafka.Producer, logger *log.Logger, userID, city string, registeredAt time.Time) error {
	event := map[string]interface{}{
		"type":          "dev.user.registered",
		"user_id":       userID,
		"city":          city,
		"registered_at": registeredAt.Format(time.RFC3339),
	}

	data, err := json.Marshal(event)
	if err != nil {
		logger.Warn("Ошибка маршалинга события регистрации: %v", err)
		return err
	}

	return producer.Send(ctx, "vira-dev.users", data)
}

func EmitDevUserLoggedInEvent(ctx context.Context, producer *kafka.Producer, logger *log.Logger, userID, ip, userAgent string, loginAt time.Time) error {
	event := map[string]interface{}{
		"type":       "dev.user.logged_in",
		"user_id":    userID,
		"ip":         ip,
		"user_agent": userAgent,
		"login_at":   loginAt.Format(time.RFC3339),
	}

	data, err := json.Marshal(event)
	if err != nil {
		logger.Warn("Ошибка маршалинга события входа: %v", err)
		return err
	}

	return producer.Send(ctx, "vira-dev.users", data)
}
