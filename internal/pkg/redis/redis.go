package redis

import "github.com/redis/go-redis/v9"

type Redis struct {
	client *redis.Client
}

func NewClient(dsn string) (*Redis, error) {
	var client *redis.Client

	opts, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, err
	}

	client = redis.NewClient(opts)

	return &Redis{client}, nil
}

func (r *Redis) Client() redis.UniversalClient {
	return r.client
}

func (r *Redis) MakeRedisClient() interface{} {
	return r.client
}
