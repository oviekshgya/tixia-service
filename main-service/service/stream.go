package service

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb                   = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	Ctx                   = context.Background()
	SearchRequestedStream = "flight.search.requested"
)

func PublishSearchRequest(values map[string]interface{}) error {
	_, err := Rdb.XAdd(Ctx, &redis.XAddArgs{
		Stream: SearchRequestedStream,
		Values: values,
	}).Result()

	if err != nil {
		log.Println("Error publishing to Redis stream:", err)
		return err
	}
	return nil
}
