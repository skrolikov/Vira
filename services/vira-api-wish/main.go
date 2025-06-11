package main

import (
	"context"
	"database/sql"
	"net/http"
	"vira-api-wish/internal/handlers"
	"vira-api-wish/internal/repo"
	"vira-api-wish/internal/service"
	"vira-api-wish/internal/viraid"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	config "github.com/skrolikov/vira-config"
	log "github.com/skrolikov/vira-logger"
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

	baseLogger.Info("🚀 Запуск Vira-Wish")

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

	// БД Postgres
	db, err := sql.Open("postgres", cfg.WishPostgresDSN)
	if err != nil {
		baseLogger.Fatal("db connect: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			baseLogger.Error("failed to close database connection: %v", err)
		}
	}()

	// Проверка соединения с БД
	if err := db.Ping(); err != nil {
		baseLogger.Fatal("db ping failed: %v", err)
	}

	// Клиент Vira-ID
	idClient := viraid.NewClient(cfg.ViraIDEndpoint)

	// Репозиторий и сервис
	upr := repo.NewUserProfileRepo(db)
	authSvc := service.NewAuthService(idClient, upr)

	r := chi.NewRouter()

	r.Post("/register", handlers.RegisterHandler(authSvc))
	r.Post("/login", handlers.LoginHandler(authSvc))

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

	baseLogger.Info("✅ Vira-Wish запущен на порту %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		baseLogger.Fatal("❌ Ошибка запуска сервера: %v", err)
	}
}
