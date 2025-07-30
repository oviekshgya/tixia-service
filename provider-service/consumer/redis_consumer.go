package consumer

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"
	"tixia-service/provider-service/mockapi"

	"github.com/redis/go-redis/v9"
)

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

var (
	ctx         = context.Background()
	streamKey   = "flight.search.requested"
	resultKey   = "flight.search.results"
	consumerGrp = "provider-group"
	consumerID  = "provider-1"

	rdb = redis.NewClient(&redis.Options{
		Addr: getEnv("REDIS_ADDR", "localhost:6379"),
	})
)

// StartFlightSearchConsumer runs a Redis Stream consumer loop
func StartFlightSearchConsumer() {
	_ = rdb.XGroupCreateMkStream(ctx, streamKey, consumerGrp, "0")

	for {
		streams, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    consumerGrp,
			Consumer: consumerID,
			Streams:  []string{streamKey, ">"},
			Block:    5 * time.Second,
			Count:    10,
		}).Result()

		if err != nil && err != redis.Nil {
			log.Println("Error reading stream:", err)
			continue
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				go handleMessage(msg)
			}
		}
	}
}

// =====================
// Handle individual message
// =====================
func handleMessage(msg redis.XMessage) {
	defer func() {
		// Acknowledge message after handling
		_ = rdb.XAck(ctx, streamKey, consumerGrp, msg.ID)
	}()

	// Parse message values
	data := make(map[string]interface{})
	for k, v := range msg.Values {
		data[k] = v
	}

	searchID, _ := data["search_id"].(string)
	from, _ := data["from"].(string)
	to, _ := data["to"].(string)
	date, _ := data["date"].(string)
	passengers, _ := data["passengers"].(string) // redis stores all values as strings

	log.Printf("Processing search: %s → %s on %s (%spax)\n", from, to, date, passengers)

	// Simulate calling external API (mock)
	results := mockapi.MockSearchFlights(from, to, date)

	// Prepare result payload
	resultData := map[string]interface{}{
		"search_id": searchID,
		"status":    "completed",
		"results":   results,
	}

	// Convert slice to JSON string to store as Redis Stream value
	resultJSON, _ := json.Marshal(resultData)

	_, err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: resultKey,
		Values: map[string]interface{}{
			"search_id": searchID,
			"status":    "completed",
			"results":   string(resultJSON),
		},
	}).Result()

	if err != nil {
		log.Println("Error publishing result:", err)
	}
}
