package main

import (
	"context"
	"net/http"

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
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	// создаём базовый логгер (без вызова WithContext)
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

	// Redis
	redisConn, err := redisdb.New(ctx, redisdb.Config{
		Addr:     cfg.RedisAddr,
		Password: "",
		DB:       cfg.RedisDB,
	}, baseLogger.WithFields(map[string]any{
		"component": "redis",
	}))
	if err != nil {
		baseLogger.Fatal("❌ Ошибка подключения к Redis: %v", err)
	}
	defer redisConn.Close()
	rdb := redisConn.Client()

	// DB
	if _, err := db.Init(ctx, cfg); err != nil {
		baseLogger.Fatal("❌ Ошибка инициализации базы: %v", err)
	}
	userRepo := db.NewUserRepository(db.Get())

	// Kafka
	kafkaLogger := baseLogger.WithFields(map[string]any{
		"component": "kafka",
	})
	producer := kafka.NewProducer(kafka.ProducerConfig{
		Brokers:      []string{cfg.KafkaAddr},
		Topic:        "vira-events",
		BatchTimeout: 100,
		Async:        false,
	}, kafkaLogger)

	defer func() {
		if err := producer.Close(); err != nil {
			kafkaLogger.Error("Ошибка закрытия Kafka продюсера: %v", err)
		}
	}()

	kafkaLogger.Info("🛰️ Kafka producer готов к работе")

	authService := service.NewAuthService(cfg, userRepo, rdb, producer, baseLogger)
	// Маршруты
	r := chi.NewRouter()

	r.Use(middleware.RequestID())               // если реализована
	r.Use(middleware.ContextLogger(baseLogger)) // если реализована

	r.Post("/login", handlers.LoginHandler(authService))
	r.Post("/register", handlers.RegisterHandler(authService))
	r.Post("/refresh", handlers.RefreshHandler(authService))

	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	r.Get("/redis-test", func(w http.ResponseWriter, r *http.Request) {
		logger := baseLogger.WithContext(r.Context())
		if err := rdb.Set(r.Context(), "test_key", "123", 0).Err(); err != nil {
			logger.WithFields(map[string]any{
				"err": err,
			}).Error("Ошибка при записи в Redis")
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
