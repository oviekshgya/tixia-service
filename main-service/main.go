package main

import (
	"github.com/gofiber/fiber/v2"
	"tixia-service/main-service/handler"
)

func main() {
	app := fiber.New()
	api := app.Group("/api/flights")

	api.Post("/search", handler.HandleSearchRequest)
	api.Get("/search/:search_id/stream", handler.StreamSearchResults)

	go handler.ListenFlightResults() // Listen Redis Stream in background

	app.Listen(":8880")
}
