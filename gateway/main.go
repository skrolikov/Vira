package main

import (
	"net/http"

	"vira-gateway/internal/router"

	config "github.com/skrolikov/vira-config"
	log "github.com/skrolikov/vira-logger"
)

func main() {
	cfg := config.Load()

	// ✅ Инициализация логгера
	logger := log.New(log.Config{
		Level:      log.DEBUG,
		JsonOutput: false,
		ShowCaller: true,
		Color:      true,
		OutputFile: "", // можно "gateway.log"
		MaxSizeMB:  10,
		MaxBackups: 3,
		MaxAgeDays: 28,
		Compress:   true,
	})

	logger.Info("🚀 Запуск Vira-Gateway")

	r := router.Setup(cfg, logger)

	logger.Info("✅ Vira-Gateway запущен на порту %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		logger.Fatal("❌ Ошибка запуска gateway: %v", err)
	}
}
