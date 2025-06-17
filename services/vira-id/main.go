package main

import (
	"context"
	"net/http"
	"time"

	"vira-id/internal/handlers"
	"vira-id/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	config "github.com/skrolikov/vira-config"
	db "github.com/skrolikov/vira-db"
	kafka "github.com/skrolikov/vira-kafka"
	log "github.com/skrolikov/vira-logger"
	middleware "github.com/skrolikov/vira-middleware"
	redisdb "github.com/skrolikov/vira-redisdb"

	kafkago "github.com/segmentio/kafka-go"
)

// @title Vira ID API
// @version 1.0
// @description Сервис авторизации Vira.
// @termsOfService http://example.com/terms/

// @contact.name Поддержка Vira
// @contact.url http://vira.example.com/support
// @contact.email support@vira.example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host vira-id:8000
// @BasePath /
func main() {
	cfg := config.Load()
	ctx := context.Background()

	baseLogger := log.New(log.Config{
		Level:      log.DEBUG,
		JsonOutput: false,
		ShowCaller: true,
		Color:      true,
		OutputFile: "",
		MaxSizeMB:  10,
		MaxBackups: 3,
		MaxAgeDays: 28,
		Compress:   true,
	})

	baseLogger.Info("🚀 Запуск Vira-ID")

	redisConn, err := redisdb.New(ctx, redisdb.Config{
		Addr:     cfg.RedisAddr,
		Password: "",
		DB:       cfg.RedisDB,
	}, baseLogger.WithFields(map[string]any{"component": "redis"}))
	if err != nil {
		baseLogger.Fatal("❌ Ошибка подключения к Redis: %v", err)
	}
	defer redisConn.Close()
	rdb := redisConn.Client()

	if _, err := db.Init(ctx, cfg); err != nil {
		baseLogger.Fatal("❌ Ошибка инициализации базы: %v", err)
	}

	dbConn, err := db.Get()
	if err != nil {
		baseLogger.Fatal("❌ Ошибка получения соединения с БД: %v", err)
	}

	userRepo := db.NewUserRepository(dbConn)

	kafkaLogger := baseLogger.WithFields(map[string]any{"component": "kafka"})

	// ✅ Конфигурация Kafka Producer
	producer := kafka.NewProducer(kafka.ProducerConfig{
		Brokers:      []string{cfg.KafkaAddr},
		Topic:        "vira-events",
		BatchTimeout: 100 * time.Millisecond,
		Async:        false,

		RequiredAcks: kafkago.RequireAll,
		Compression:  kafkago.Snappy,
		MaxAttempts:  5,

		Logger: kafkaLogger,
		Tracer: nil, // или otel.Tracer("vira-id"), если уже подключён
	})
	defer producer.Close()

	authService := service.NewAuthService(cfg, userRepo, rdb, producer, baseLogger)

	r := chi.NewRouter()

	r.Use(middleware.RequestID())
	r.Use(middleware.ContextLogger(baseLogger))

	r.Post("/login", handlers.LoginHandler(authService))
	r.Post("/register", handlers.RegisterHandler(authService))
	r.Post("/refresh", handlers.RefreshHandler(authService))

	// Новый маршрут подтверждения
	r.Get("/confirm", handlers.ConfirmUserHandler(authService))

	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	r.Get("/redis-test", func(w http.ResponseWriter, r *http.Request) {
		logger := baseLogger.WithContext(r.Context())
		if err := rdb.Set(r.Context(), "test_key", "123", 0).Err(); err != nil {
			logger.WithFields(map[string]any{"err": err}).Error("Ошибка при записи в Redis")
			http.Error(w, "ошибка Redis", http.StatusInternalServerError)
			return
		}
		val, _ := rdb.Get(r.Context(), "test_key").Result()
		w.Write([]byte("Redis работает, значение: " + val))
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(cfg, baseLogger))
		r.Get("/me", handlers.MeHandler(userRepo))
		r.Post("/logout", handlers.LogoutHandler(cfg, rdb))
		r.Get("/sessions", handlers.SessionsHandler(cfg, rdb))
		r.Delete("/sessions/{id}", handlers.DeleteSessionHandler(cfg, rdb))
	})

	baseLogger.Info("✅ Vira-ID запущен на порту %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		baseLogger.Fatal("❌ Ошибка запуска сервера: %v", err)
	}
}
