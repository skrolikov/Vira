package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"vira-api-dev/internal/handlers"
	"vira-api-dev/internal/repo"
	"vira-api-dev/internal/service"
	"vira-api-dev/internal/viraid"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	kafkago "github.com/segmentio/kafka-go"
	config "github.com/skrolikov/vira-config"
	kafka "github.com/skrolikov/vira-kafka"
	log "github.com/skrolikov/vira-logger"
	middleware "github.com/skrolikov/vira-middleware"
	redisdb "github.com/skrolikov/vira-redisdb"
)

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

	baseLogger.Info("🚀 Запуск Vira-DEV")

	// Подключение к Redis
	redisLogger := baseLogger.WithFields(map[string]any{"component": "redis"})
	redisConn, err := redisdb.New(ctx, redisdb.Config{
		Addr:     cfg.RedisAddr,
		Password: "",
		DB:       cfg.RedisDB,
	}, redisLogger)
	if err != nil {
		redisLogger.Fatal("❌ Ошибка подключения к Redis: %v", err)
	}
	defer redisConn.Close()
	rdb := redisConn.Client()

	// Подключение к PostgreSQL
	db, err := sql.Open("postgres", cfg.DevPostgresDSN)
	if err != nil {
		baseLogger.Fatal("❌ Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		baseLogger.Fatal("❌ Ping к БД не прошёл: %v", err)
	}

	// Настройка Kafka Producer
	kafkaLogger := baseLogger.WithFields(map[string]any{"component": "kafka"})
	producer := kafka.NewProducer(kafka.ProducerConfig{
		Brokers:      []string{cfg.KafkaAddr},
		Topic:        "vira-dev.users",
		BatchTimeout: 100 * time.Millisecond,
		Async:        false,

		RequiredAcks: kafkago.RequireAll,
		Compression:  kafkago.Snappy,
		MaxAttempts:  5,

		Logger: kafkaLogger,
		Tracer: nil, // future: otel.Tracer("vira-api-dev")
	})
	defer func() {
		if err := producer.Close(); err != nil {
			kafkaLogger.Error("Ошибка закрытия Kafka продюсера: %v", err)
		}
	}()

	// Инициализация зависимостей
	idClient := viraid.NewClient(cfg.ViraIDEndpoint)
	userRepo := repo.NewUserProfileRepo(db)
	authService := service.NewAuthService(idClient, userRepo, producer, baseLogger)

	// HTTP роутер
	r := chi.NewRouter()
	r.Use(middleware.RequestID())
	r.Use(middleware.ContextLogger(baseLogger))

	r.Post("/register", handlers.RegisterHandler(authService))
	r.Post("/login", handlers.LoginHandler(authService))

	// Пример Redis-маршрута
	r.Get("/redis-test", func(w http.ResponseWriter, r *http.Request) {
		logger := baseLogger.WithContext(r.Context())
		if err := rdb.Set(r.Context(), "test_key", "123", 0).Err(); err != nil {
			logger.WithFields(map[string]any{"err": err}).Error("Ошибка Redis")
			http.Error(w, "ошибка Redis", http.StatusInternalServerError)
			return
		}
		val, _ := rdb.Get(r.Context(), "test_key").Result()
		w.Write([]byte("Redis работает, значение: " + val))
	})

	// Метрики Prometheus
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	// Запуск HTTP-сервера
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		baseLogger.Info("✅ Vira-DEV запущен на порту %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			baseLogger.Fatal("❌ Ошибка запуска сервера: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	baseLogger.Info("🛑 Остановка сервера...")

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxTimeout); err != nil {
		baseLogger.Error("Ошибка при остановке сервера: %v", err)
	} else {
		baseLogger.Info("✅ Сервер успешно остановлен")
	}
}
