package cache

import "context"

type RedisCache struct {
}

func NewRedisCache() *RedisCache {
	return nil
}

func (rc *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	return nil
}
