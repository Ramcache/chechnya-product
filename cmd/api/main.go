package main

import (
	"chechnya-product/config"
	"chechnya-product/internal/app"
	"chechnya-product/internal/db"
	"chechnya-product/internal/logger"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	// 🪵 Логгер
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// ⚙️ Конфиг
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// 📦 Подключение к базе
	dbConn, err := db.NewPostgresDB(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer dbConn.Close()

	// 🐥 Применяем миграции через goose
	if err := goose.Up(dbConn.DB, "migrations"); err != nil {
		logger.Fatal("Failed to apply migrations", zap.Error(err))
	}

	// 🚀 Запуск сервера
	srv := app.NewServer(cfg, logger, dbConn)
	logger.Sugar().Infow("Server is running", "port", cfg.Port)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}
