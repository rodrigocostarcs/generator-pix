package cache

import (
	"context"
	"encoding/json"
	"time"
)

// CacheAdapter define a interface para o sistema de cache
type CacheAdapter interface {
	// Get recupera um valor do cache
	Get(ctx context.Context, key string) ([]byte, error)

	// Set armazena um valor no cache com um tempo de expiração
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Delete remove um valor do cache
	Delete(ctx context.Context, key string) error
}

// GetObject recupera e desserializa um objeto do cache
func GetObject(adapter CacheAdapter, ctx context.Context, key string, obj interface{}) error {
	data, err := adapter.Get(ctx, key)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, obj)
}

// SetObject serializa e armazena um objeto no cache
func SetObject(adapter CacheAdapter, ctx context.Context, key string, obj interface{}, expiration time.Duration) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	return adapter.Set(ctx, key, data, expiration)
}
