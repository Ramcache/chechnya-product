package main

import (
	"chechnya-product/config"
	"chechnya-product/internal/app"
	"chechnya-product/internal/logger"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	// 🪵 Инициализация логгера
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// ⚙️ Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// 🚀 Запуск HTTP-сервера
	srv := app.NewServer(cfg, logger)
	logger.Sugar().Infow("Server is running", "port", cfg.Port)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}
