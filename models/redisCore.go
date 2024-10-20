package models

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var (
	RedisDb *redis.Client
)

func init() {
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := RedisDb.Ping(ctx).Result()
	if err != nil {
		println(err)
	}
}
