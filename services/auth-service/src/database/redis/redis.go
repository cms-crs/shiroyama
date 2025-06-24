package redis

import (
	"authservice/src/config"
	"github.com/redis/go-redis/v9"
)

func MustConnect(cfg *config.Config) *redis.Client {
	rdb, err := Connect(cfg)
	if err != nil {
		panic(err)
	}

	return rdb
}

func Connect(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // or "redis:6379" if running from another service in Docker
	})

	return rdb, nil
}
