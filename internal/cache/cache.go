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
func (c *RedisCache) GetOrSet(ctx context.Context, key string, dest any, fetch func() (any, error)) error {
	val, err := c.client.Get(ctx, key).Result()
	if err == nil {
		c.logger.Info("📦 Кэш HIT", zap.String("key", key))
		return json.Unmarshal([]byte(val), dest)
	}

	if err != redis.Nil {
		c.logger.Error("🚨 Ошибка при чтении из Redis", zap.String("key", key), zap.Error(err))
		return err
	}

	c.logger.Info("💨 Кэш MISS — получаем заново", zap.String("key", key))

	result, err := fetch()
	if err != nil {
		c.logger.Error("❌ Ошибка при получении данных", zap.String("key", key), zap.Error(err))
		return err
	}

	data, err := json.Marshal(result)
	if err != nil {
		c.logger.Error("❌ Ошибка сериализации результата", zap.Error(err))
		return err
	}

	if err := c.client.Set(ctx, key, data, c.ttl).Err(); err != nil {
		c.logger.Warn("⚠️ Не удалось сохранить в Redis", zap.String("key", key), zap.Error(err))
	}

	// декодируем обратно
	raw, _ := json.Marshal(result)
	return json.Unmarshal(raw, dest)
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
