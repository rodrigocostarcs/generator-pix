package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisAdapter implementa CacheAdapter usando Redis
type RedisAdapter struct {
	client *redis.Client
}

// NewRedisAdapter cria uma nova instância do adaptador Redis
func NewRedisAdapter(host string, port string, password string, db int) *RedisAdapter {
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
	})

	return &RedisAdapter{
		client: client,
	}
}

// Get implementa a interface CacheAdapter
func (r *RedisAdapter) Get(ctx context.Context, key string) ([]byte, error) {
	return r.client.Get(ctx, key).Bytes()
}

// Set implementa a interface CacheAdapter
func (r *RedisAdapter) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Delete implementa a interface CacheAdapter
func (r *RedisAdapter) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Close fecha a conexão com o Redis
func (r *RedisAdapter) Close() error {
	return r.client.Close()
}
