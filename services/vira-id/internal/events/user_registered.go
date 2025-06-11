package events

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"vira-id/internal/metrics"

	kafka "github.com/skrolikov/vira-kafka"
	log "github.com/skrolikov/vira-logger"
)

// EventType тип события для системы сообщений
type EventType string

const (
	UserRegisteredEvent EventType = "user.registered"
	UserLoggedInEvent   EventType = "user.logged_in"
	UserUpdatedEvent    EventType = "user.updated"
)

// UserEventPayload структура события пользователя
type UserEventPayload struct {
	EventID   string    `json:"event_id"`   // Уникальный ID события
	EventType EventType `json:"event_type"` // Тип события
	UserID    string    `json:"user_id"`    // ID пользователя
	Username  string    `json:"username"`   // Имя пользователя
	IP        string    `json:"ip"`         // IP-адрес
	Device    string    `json:"device"`     // Информация об устройстве (User-Agent)
	Timestamp time.Time `json:"timestamp"`  // Время события
	Metadata  Metadata  `json:"metadata"`   // Дополнительные метаданные
}

// Metadata доп. данные события
type Metadata map[string]interface{}

// EventEmitter интерфейс отправки событий
type EventEmitter interface {
	EmitUserEvent(ctx context.Context, eventType EventType, payload UserEventPayload) error
}

// KafkaEventEmitter реализация для Kafka
type KafkaEventEmitter struct {
	producer *kafka.Producer
	logger   *log.Logger
}

func NewKafkaEventEmitter(producer *kafka.Producer, logger *log.Logger) *KafkaEventEmitter {
	return &KafkaEventEmitter{producer: producer, logger: logger}
}

func (e *KafkaEventEmitter) EmitUserEvent(ctx context.Context, eventType EventType, payload UserEventPayload) error {
	if payload.UserID == "" {
		return errors.New("userID не может быть пустым")
	}

	payload.EventID = generateEventID()
	payload.EventType = eventType
	payload.Timestamp = time.Now().UTC()

	data, err := json.Marshal(payload)
	if err != nil {
		e.logger.Error("Ошибка сериализации события %s: %v", eventType, err)
		return fmt.Errorf("ошибка сериализации события: %w", err)
	}

	if err := e.producer.Send(ctx, string(eventType), data); err != nil {
		metrics.KafkaErrors.WithLabelValues(string(eventType)).Inc()
		e.logger.Error("Ошибка отправки события %s: %v", eventType, err)
		return fmt.Errorf("ошибка отправки события: %w", err)
	}

	metrics.KafkaEvents.WithLabelValues(string(eventType)).Inc()
	e.logger.Debug("Событие %s успешно отправлено для пользователя %s", eventType, payload.UserID)
	return nil
}

// Вспомогательные функции генерации ID

func generateEventID() string {
	return fmt.Sprintf("evt_%d_%s", time.Now().UnixNano(), randomString(8))
}

func randomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

// Функции для отправки конкретных событий

func EmitUserRegisteredEvent(
	ctx context.Context,
	producer *kafka.Producer,
	logger *log.Logger,
	userID, username, ip, device string,
	ts time.Time,
) error {
	// Создаём отдельный контекст для Kafka, чтобы не зависеть от ctx запроса
	ctxKafka, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	emitter := NewKafkaEventEmitter(producer, logger)
	payload := UserEventPayload{
		UserID:    userID,
		Username:  username,
		IP:        ip,
		Device:    device,
		Timestamp: ts,
		Metadata:  Metadata{"source": "auth_service"},
	}
	return emitter.EmitUserEvent(ctxKafka, UserRegisteredEvent, payload)
}

func EmitUserLoggedInEvent(
	ctx context.Context,
	producer *kafka.Producer,
	logger *log.Logger,
	userID, username, ip, device string,
) {
	// Отдельный контекст с таймаутом для Kafka, чтобы не зависеть от контекста запроса
	ctxKafka, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	emitter := NewKafkaEventEmitter(producer, logger)
	payload := UserEventPayload{
		UserID:    userID,
		Username:  username,
		IP:        ip,
		Device:    device,
		Timestamp: time.Now().UTC(),
		Metadata:  Metadata{"source": "auth_service"},
	}
	if err := emitter.EmitUserEvent(ctxKafka, UserLoggedInEvent, payload); err != nil {
		logger.Error("Ошибка при отправке события входа: %v", err)
	}
}
