package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
	logger *zap.Logger
}

func NewRedisCache(client *redis.Client, ttl time.Duration, logger *zap.Logger) *RedisCache {
	return &RedisCache{
		client: client,
		ttl:    ttl,
		logger: logger,
	}
}

// GetOrSet проверяет кэш и, если промах — вызывает fetch() и сохраняет результат
func (c *RedisCache) GetOrSet(ctx context.Context, key string, target any, fetch func() (any, error)) error {
	// Попробовать получить из кэша
	val, err := c.client.Get(ctx, key).Result()
	if err == nil {
		return json.Unmarshal([]byte(val), target)
	}

	if err != redis.Nil {
		c.logger.Warn("Redis Get error", zap.Error(err))
		// Возвращать nil, чтобы не падал весь API
	}

	// Получаем данные заново
	data, err := fetch()
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// 🕒 Устанавливаем TTL (например, 10 минут)
	err = c.client.Set(ctx, key, bytes, 10*time.Minute).Err()
	if err != nil {
		c.logger.Warn("Redis Set error", zap.Error(err))
	}

	// Присваиваем результат
	encoded, _ := json.Marshal(data)
	return json.Unmarshal(encoded, target)
}

func (c *RedisCache) ClearPrefix(ctx context.Context, prefix string) error {
	iter := c.client.Scan(ctx, 0, prefix+"*", 0).Iterator()
	deleted := 0

	for iter.Next(ctx) {
		key := iter.Val()
		if err := c.client.Del(ctx, key).Err(); err == nil {
			c.logger.Info("🧹 Удалён кэш по префиксу", zap.String("key", key))
			deleted++
		} else {
			c.logger.Warn("⚠️ Не удалось удалить ключ", zap.String("key", key), zap.Error(err))
		}
	}

	if err := iter.Err(); err != nil {
		c.logger.Error("🚨 Ошибка при сканировании ключей Redis", zap.String("prefix", prefix), zap.Error(err))
		return err
	}

	c.logger.Info("🧼 Очистка по префиксу завершена", zap.String("prefix", prefix), zap.Int("удалено", deleted))
	return nil
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		c.logger.Warn("⚠️ Не удалось удалить ключ", zap.String("key", key), zap.Error(err))
	} else {
		c.logger.Info("🗑 Удалён ключ", zap.String("key", key))
	}
	return err
}
