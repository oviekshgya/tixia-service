package handler

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	"tixia-service/main-service/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	//rdb          = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	rdb = redis.NewClient(&redis.Options{
		Addr: getEnv("REDIS_ADDR", "localhost:6379"),
	})
	ctx          = context.Background()
	streamName   = "flight.search.requested"
	resultStream = "flight.search.results"
	sseClients   = make(map[string]chan string) // key: search_id, value: channel
)

type SearchRequest struct {
	From       string `json:"from"`
	To         string `json:"to"`
	Date       string `json:"date"`
	Passengers int    `json:"passengers"`
}

type SearchResult struct {
	SearchID string        `json:"search_id"`
	Status   string        `json:"status"`
	Results  []interface{} `json:"results,omitempty"`
}

func HandleSearchRequest(c *fiber.Ctx) error {
	var req SearchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request",
		})
	}

	searchID := uuid.NewString()

	// Publish to Redis Stream
	values := map[string]interface{}{
		"search_id":  searchID,
		"from":       req.From,
		"to":         req.To,
		"date":       req.Date,
		"passengers": req.Passengers,
	}
	if err := service.PublishSearchRequest(values); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to send request",
		})
	}

	// Prepare SSE channel
	sseClients[searchID] = make(chan string, 10)
	log.Println("[SSE INIT] Registered search_id:", searchID)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Search request submitted",
		"data": fiber.Map{
			"search_id": searchID,
			"status":    "processing",
		},
	})
}

func StreamSearchResults(c *fiber.Ctx) error {
	searchID := c.Params("search_id")
	ch, ok := sseClients[searchID]
	if !ok {
		log.Println("[SSE ERROR] search_id not found in sseClients:", searchID)
		return c.Status(404).SendString("search_id not found or expired")
	}

	// Set SSE headers
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		for msg := range ch {
			fmt.Fprintf(w, "data: %s\n\n", msg)
			w.Flush()
		}
	})

	return nil
}

func ListenFlightResults() {
	group := "main-service-group"
	consumer := "consumer-1"

	// Create consumer group (ignore error if already exists)
	_ = rdb.XGroupCreateMkStream(ctx, resultStream, group, "0")

	for {
		streams, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: consumer,
			Streams:  []string{resultStream, ">"},
			Block:    5 * time.Second,
			Count:    10,
		}).Result()

		if err != nil && err != redis.Nil {
			log.Println("XREADGROUP error:", err)
			continue
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				data := make(map[string]interface{})
				for k, v := range msg.Values {
					data[k] = v
				}

				jsonData, _ := json.Marshal(data)
				searchID := fmt.Sprintf("%v", data["search_id"])
				log.Println("[SSE SEND] Sending result for:", searchID)

				if ch, ok := sseClients[searchID]; ok {
					ch <- string(jsonData)

					// Close SSE channel when complete
					if status := fmt.Sprintf("%v", data["status"]); status == "completed" {
						time.Sleep(50 * time.Second)
						close(ch)
						//delete(sseClients, searchID)
					}
				}

				// Acknowledge message
				_ = rdb.XAck(ctx, resultStream, group, msg.ID)
			}
		}
	}
}
