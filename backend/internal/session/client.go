package session

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})

	if err := Rdb.Ping(Ctx).Err(); err != nil {
		log.Fatalf("Не удалось подключится к Redis: %v", err)
	}
}
