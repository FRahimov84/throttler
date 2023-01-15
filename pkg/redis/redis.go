package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func New(host, port, pass string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: pass,
		DB:       0,
	})
	cmd := client.Ping(context.Background())

	_, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
