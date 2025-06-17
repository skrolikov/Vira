// cmd/consumer/main.go
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	// если нужен какой-то сервис
	config "github.com/skrolikov/vira-config"
	db "github.com/skrolikov/vira-db"
	kafka "github.com/skrolikov/vira-kafka"
	log "github.com/skrolikov/vira-logger"
)

type UserLoggedInEvent struct {
	Type      string `json:"type"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	IP        string `json:"ip"`
	UA        string `json:"ua"`
	Time      string `json:"time"`
	SessionID string `json:"session_id"`
	Success   bool   `json:"success"`
	Device    string `json:"device"`
}

func main() {
	cfg := config.Load()

	logger := log.New(log.Config{
		Level:      log.DEBUG,
		JsonOutput: false,
		ShowCaller: true,
		Color:      true,
	})

	logger.Info("🚀 Запуск consumer vira-events")
	logger.Info("KafkaAddr: %s", cfg.KafkaAddr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// graceful shutdown по сигналам
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigchan
		logger.Info("Получен сигнал %v, завершаем работу...", sig)
		cancel()
	}()

	// --- Инициализация БД ---
	if _, err := db.Init(ctx, cfg); err != nil {
		logger.Fatal("❌ Ошибка инициализации базы: %v", err)
	}
	dbConn, err := db.Get()
	if err != nil {
		logger.Fatal("❌ Ошибка получения соединения с БД: %v", err)
	}
	userLoginRepo := db.NewUserLoginRepository(dbConn)
	logger.Info("✅ База данных инициализирована")

	// --- Инициализация Kafka consumer ---
	consumer := kafka.NewConsumer(kafka.ConsumerConfig{
		Brokers:           []string{cfg.KafkaAddr},
		Topic:             "vira-events",
		GroupID:           "vira-events-consumer",
		MinBytes:          1 * 1024,         // 1KB
		MaxBytes:          10 * 1024 * 1024, // 10MB
		MaxWait:           500 * time.Millisecond,
		CommitInterval:    1 * time.Second,
		HeartbeatInterval: 3 * time.Second,
		SessionTimeout:    30 * time.Second,
		StartOffset:       -1, // LastOffset
		MaxRetryAttempts:  3,
		RetryBackoff:      200 * time.Millisecond,

		// при желании можно прокинуть метрики и трейсер:
		Metrics: nil,
		Tracer:  nil,

		Logger: logger.WithFields(map[string]any{"component": "kafka-consumer"}),
	})
	defer func() {
		if err := consumer.Close(); err != nil {
			logger.Error("Ошибка закрытия Kafka consumer: %v", err)
		}
	}()

	logger.Info("🛰️ Kafka consumer готов к работе")

	// --- Основной цикл обработки через каналы ---
	msgCh, errCh := consumer.Consume(ctx)
outer:
	for {
		select {
		case <-ctx.Done():
			logger.Info("Контекст отменён, выходим из цикла")
			break outer

		case err := <-errCh:
			logger.Error("Ошибка в consumer: %v", err)

		case msg, ok := <-msgCh:
			if !ok {
				logger.Info("Канал сообщений закрыт, выходим")
				break outer
			}

			// декодируем
			var event UserLoggedInEvent
			if err := consumer.DecodeMessage(msg, &event); err != nil {
				logger.Error("Ошибка декодирования сообщения: %v", err)
				consumer.CommitMessage(ctx, msg) // пропускаем «сломанное» сообщение
				continue
			}

			if event.Type != "UserLoggedIn" {
				// просто метим как обработанное
				consumer.CommitMessage(ctx, msg)
				continue
			}

			// парсим время
			parsedTime, err := time.Parse(time.RFC3339, event.Time)
			if err != nil {
				logger.Warn("Некорректное время в событии %q, ставим now: %v", event.Time, err)
				parsedTime = time.Now()
			}

			// сохраняем в БД
			if err := userLoginRepo.Save(ctx,
				event.UserID,
				event.Username,
				event.IP,
				event.UA,
				event.SessionID,
				parsedTime,
				event.Success,
				event.Device,
			); err != nil {
				logger.Error("Ошибка сохранения события в БД: %v", err)
				continue
			}

			// подтверждаем оффсет
			if err := consumer.CommitMessage(ctx, msg); err != nil {
				logger.Error("Не удалось закоммитить сообщение: %v", err)
			}

			logger.Info("✅ UserLoggedIn сохранён: user_id=%s username=%s", event.UserID, event.Username)
		}
	}

	logger.Info("🔚 Завершение работы vira-events consumer")
}
