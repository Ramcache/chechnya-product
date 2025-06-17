// main.go
package main

import (
	"chechnya-product/config"
	"chechnya-product/internal/app"
	"chechnya-product/internal/cache"
	"chechnya-product/internal/db"
	"chechnya-product/internal/logger"
	"context"
	"errors"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"time"
)

// @title Chechnya Product API
// @version 5.0
// @description Backend for products shop
// @host localhost:8080
// @BasePath /api
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

	// ✅ Проверка аргументов
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		logger.Sugar().Info("Running goose migrations...")
		if err := goose.Up(dbConn.DB, "migrations"); err != nil {
			logger.Fatal("Failed to apply migrations", zap.Error(err))
		}
		logger.Sugar().Info("Migrations completed.")
		return
	}
	redisClient := redis.NewClient(cfg.GetRedisOptions())

	logger.Sugar().Infow("🔌 Подключение к Redis", "addr", cfg.RedisAddr)

	logger.Sugar().Info("✅ Успешное подключение к Redis")

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		logger.Fatal("Не удалось подключиться к Redis", zap.Error(err))
	}

	// Обёртка RedisCache
	ttl := time.Duration(cfg.RedisTTLMinutes) * time.Minute
	redisCache := cache.NewRedisCache(redisClient, ttl, logger)

	// 🚀 Запуск сервера
	srv := app.NewServer(cfg, logger, dbConn, redisCache)
	logger.Sugar().Infow("Server is running", "port", cfg.Port)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}
