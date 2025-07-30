package mockapi

import (
	"fmt"
	"math/rand"
	"time"
)

type FlightResult struct {
	FlightNumber string  `json:"flight_number"`
	Airline      string  `json:"airline"`
	From         string  `json:"from"`
	To           string  `json:"to"`
	Departure    string  `json:"departure_time"`
	Arrival      string  `json:"arrival_time"`
	Price        float64 `json:"price"`
}

func MockSearchFlights(from, to, date string) []FlightResult {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(1 * time.Second)

	airlines := []string{"Garuda Indonesia", "Lion Air", "AirAsia", "Citilink"}
	results := make([]FlightResult, 0)

	// Return 2–4 dummy flights
	for i := 0; i < rand.Intn(3)+2; i++ {
		depart := randomTime()
		arrive := depart.Add(time.Duration(rand.Intn(3)+1) * time.Hour)

		result := FlightResult{
			FlightNumber: fmt.Sprintf("GA%d", rand.Intn(900)+100),
			Airline:      airlines[rand.Intn(len(airlines))],
			From:         from,
			To:           to,
			Departure:    depart.Format("15:04"),
			Arrival:      arrive.Format("15:04"),
			Price:        float64(rand.Intn(500)+500) * 1000,
		}
		results = append(results, result)
	}

	return results
}

func randomTime() time.Time {
	hour := rand.Intn(12) + 6 // 6–17
	min := rand.Intn(60)
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location())
}
