package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	Port       string
	JWTSecret  string
	Env        string // "development", "production", "test"
}

// LoadConfig загружает конфигурацию из переменных окружения.
func LoadConfig() (*Config, error) {
	_ = godotenv.Load() // .env не обязателен, но полезен в dev

	cfg := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		Port:       os.Getenv("PORT"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		Env:        os.Getenv("ENV"),
	}

	// Минимальная валидация
	if cfg.Port == "" || cfg.JWTSecret == "" {
		return nil, errors.New("missing required environment variables (PORT, JWT_SECRET)")
	}

	return cfg, nil
}
