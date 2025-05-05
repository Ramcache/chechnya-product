package config

import (
	"errors"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"os"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	Port          string
	JWTSecret     string
	Env           string
	RedisAddr     string
	RedisPassword string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		Port:          os.Getenv("PORT"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		Env:           os.Getenv("ENV"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}

	if cfg.Port == "" || cfg.JWTSecret == "" || cfg.RedisAddr == "" {
		return nil, errors.New("missing required environment variables (PORT, JWT_SECRET, REDIS_ADDR)")
	}

	return cfg, nil
}

func (c *Config) GetRedisOptions() *redis.Options {
	return &redis.Options{
		Addr:     c.RedisAddr,
		Password: c.RedisPassword,
		DB:       0,
	}
}
