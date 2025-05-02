package redis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	rdb  *redis.Client
	rctx = context.Background()
)

func InitRedis() *redis.Client {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := rdb.Ping(rctx).Result()
	if err != nil {
		log.Fatal("failed to connect to redis: ", err)
	}
	return rdb
}
