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

// EventType — тип события в системе сообщений
//
// swagger:model EventType
type EventType string

const (
	// UserRegisteredEvent — событие регистрации пользователя
	UserRegisteredEvent EventType = "user.registered"

	// UserLoggedInEvent — событие входа пользователя
	UserLoggedInEvent EventType = "user.logged_in"

	// UserUpdatedEvent — событие обновления данных пользователя
	UserUpdatedEvent EventType = "user.updated"
)

// UserEventPayload — структура полезной нагрузки пользовательского события
//
// swagger:model UserEventPayload
type UserEventPayload struct {
	// Уникальный ID события (генерируется автоматически)
	// example: evt_1655000000000000000_ab12cd34
	EventID string `json:"event_id"`

	// Тип события
	// example: user.registered
	EventType EventType `json:"event_type"`

	// ID пользователя
	// example: 123e4567-e89b-12d3-a456-426614174000
	UserID string `json:"user_id"`

	// Имя пользователя
	// example: sergey
	Username string `json:"username"`

	// IP-адрес пользователя
	// example: 192.168.1.101
	IP string `json:"ip"`

	// Информация об устройстве (User-Agent)
	// example: Mozilla/5.0 (Windows NT 10.0; Win64; x64)
	Device string `json:"device"`

	// Время события в UTC
	// example: 2025-06-12T14:40:00Z
	Timestamp time.Time `json:"timestamp"`

	// Дополнительные метаданные события
	Metadata Metadata `json:"metadata"`
}

// Metadata — тип для хранения дополнительных данных события
type Metadata map[string]interface{}

// EventEmitter — интерфейс отправки событий
type EventEmitter interface {
	// EmitUserEvent отправляет событие с типом и полезной нагрузкой
	EmitUserEvent(ctx context.Context, eventType EventType, payload UserEventPayload) error
}

// KafkaEventEmitter — реализация EventEmitter для Kafka
type KafkaEventEmitter struct {
	producer *kafka.Producer
	logger   *log.Logger
}

// NewKafkaEventEmitter создаёт новый эмиттер для Kafka
func NewKafkaEventEmitter(producer *kafka.Producer, logger *log.Logger) *KafkaEventEmitter {
	return &KafkaEventEmitter{producer: producer, logger: logger}
}

// EmitUserEvent отправляет событие в Kafka, сериализуя полезную нагрузку в JSON
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

// generateEventID генерирует уникальный идентификатор события
func generateEventID() string {
	return fmt.Sprintf("evt_%d_%s", time.Now().UnixNano(), randomString(8))
}

// randomString генерирует случайную строку длиной length из hex-символов
func randomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

// EmitUserRegisteredEvent отправляет событие регистрации пользователя в Kafka
func EmitUserRegisteredEvent(
	ctx context.Context,
	producer *kafka.Producer,
	logger *log.Logger,
	userID, username, ip, device string,
	ts time.Time,
) error {
	// Создаём отдельный контекст с таймаутом, чтобы не зависеть от контекста запроса
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

// EmitUserLoggedInEvent отправляет событие входа пользователя в Kafka
func EmitUserLoggedInEvent(
	ctx context.Context,
	producer *kafka.Producer,
	logger *log.Logger,
	userID, username, ip, device string,
) {
	// Создаём отдельный контекст с таймаутом, чтобы не зависеть от контекста запроса
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
