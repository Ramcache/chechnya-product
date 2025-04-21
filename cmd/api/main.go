package main

import (
	"chechnya-product/config"
	"chechnya-product/internal/app"
	"chechnya-product/internal/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os/exec"
)

// 🔌 Запуск WhatsApp бота как отдельного процесса
func startWhatsAppBot() error {
	cmd := exec.Command("node", "whatsapp-bot/index.js")
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()

	if err := cmd.Start(); err != nil {
		return err
	}

	log.Println("📲 WhatsApp бот запущен как отдельный процесс")
	return nil
}

func main() {
	// 🪵 Инициализация логгера
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// 🧪 Запуск бота WhatsApp
	if err := startWhatsAppBot(); err != nil {
		logger.Fatal("Failed to start WhatsApp bot", zap.Error(err))
	}

	// ⚙️ Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// 🧠 Подключение к Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       0,
	})

	// 🚀 Запуск HTTP-сервера
	srv := app.NewServer(cfg, logger, redisClient)
	logger.Sugar().Infow("Server is running", "port", cfg.Port)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}
