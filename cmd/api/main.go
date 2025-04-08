package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"myshop/config"
	"myshop/internal/db"
	"myshop/internal/middleware"
)

func main() {
	cfg := config.LoadConfig()

	// логгер
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// БД
	database, err := db.NewPostgresDB(cfg)
	if err != nil {
		logger.Fatal("Не удалось подключиться к БД", zap.Error(err))
	}
	defer database.Close()

	// Роутер
	r := mux.NewRouter()
	r.Use(middleware.LoggerMiddleware(logger))

	// TODO: здесь будут роуты

	log.Printf("Сервер запущен на порту %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
