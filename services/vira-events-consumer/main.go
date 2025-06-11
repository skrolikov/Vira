package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/skrolikov/vira-config"
	db "github.com/skrolikov/vira-db"
	kafka "github.com/skrolikov/vira-kafka"
	log "github.com/skrolikov/vira-logger"
)

type UserLoggedInEvent struct {
	Type     string `json:"type"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IP       string `json:"ip"`
	UA       string `json:"ua"`
	Time     string `json:"time"`
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

	// Инициализация базы данных
	if _, err := db.Init(ctx, cfg); err != nil {
		logger.Fatal("❌ Ошибка инициализации базы: %v", err)
	}
	userLoginRepo := db.NewUserLoginRepository(db.Get())
	logger.Info("✅ База данных инициализирована")

	// Инициализация Kafka consumer
	consumer := kafka.NewConsumer(kafka.ConsumerConfig{
		Brokers: []string{cfg.KafkaAddr},
		Topic:   "vira-events",
		GroupID: "vira-events-consumer",
	}, logger)
	defer func() {
		if err := consumer.Close(); err != nil {
			logger.Error("Ошибка закрытия Kafka consumer: %v", err)
		} else {
			logger.Info("Kafka consumer закрыт")
		}
	}()

	logger.Info("🛰️ Kafka consumer готов к работе")

	// Обработка сигналов завершения
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigchan
		logger.Info("Получен сигнал %v, завершаем работу...", sig)
		cancel()
	}()

	// Основной цикл обработки сообщений
	for {
		msg, err := consumer.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				logger.Info("Контекст отменён, завершаем обработку")
				break
			}
			logger.Error("Ошибка чтения сообщения Kafka: %v", err)
			continue
		}

		go func(msg []byte) {
			var event UserLoggedInEvent
			if err := json.Unmarshal(msg, &event); err != nil {
				logger.Error("Ошибка парсинга события: %v", err)
				return
			}

			if event.Type != "UserLoggedIn" {
				// Игнорируем другие события
				return
			}

			parsedTime, err := time.Parse(time.RFC3339, event.Time)
			if err != nil {
				logger.Error("Ошибка парсинга времени события: %v, ставим текущее время", err)
				parsedTime = time.Now()
			}

			if err := userLoginRepo.Save(ctx, event.UserID, event.Username, event.IP, event.UA, parsedTime); err != nil {
				logger.Error("Ошибка сохранения события в БД: %v", err)
				return
			}

			logger.Info("Событие UserLoggedIn сохранено: user_id=%s username=%s ip=%s ua=%s",
				event.UserID, event.Username, event.IP, event.UA)
		}(msg.Value)
	}

	logger.Info("Завершение работы vira-events")
}
