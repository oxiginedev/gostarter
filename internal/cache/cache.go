package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Has(ctx context.Context, key string) (bool, error)
	Forget(ctx context.Context, key string) error
}
