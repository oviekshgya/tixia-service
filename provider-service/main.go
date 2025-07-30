package main

import (
	"log"
	"tixia-service/provider-service/consumer"
)

func main() {
	log.Println("Starting provider service...")
	consumer.StartFlightSearchConsumer()
}
