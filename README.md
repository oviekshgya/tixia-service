# Flight Search System

A microservice-based flight search system built using:

- **Language:** Go
- **Framework:** Fiber
- **Messaging:** Redis Streams
- **Streaming:** Server-Sent Events (SSE)

## ✨ Overview

This system allows users to search for flights via a REST endpoint. Results are asynchronously processed and streamed to the client in real time using SSE.

It consists of:

- `main-service`: Receives search requests, streams responses.
- `provider-service`: Consumes search requests and publishes mocked results.

---

## ⚙️ Setup Instructions

### Prerequisites

- Docker & Docker Compose

### Run the system:

```bash
docker-compose up --build
```

This will start:

- Redis on port `6379`
- main-service on port `8080`
- provider-service (no external port needed)

---

## 📉 API Usage

### 1. Submit Flight Search

```bash
curl -X POST http://localhost:8080/api/flights/search \
  -H 'Content-Type: application/json' \
  -d '{
    "from": "CGK",
    "to": "DPS",
    "date": "2025-07-10",
    "passengers": 2
}'
```

**Response:**

```json
{
  "success": true,
  "message": "Search request submitted",
  "data": {
    "search_id": "cef6b50b-6052-485a-83d6-224e388499af",
    "status": "processing"
  }
}
```

### 2. Stream Flight Results

```bash
curl http://localhost:8080/api/flights/search/<search_id>/stream
```

**Example SSE Output:**

```
data: {"search_id":"...","status":"completed","results":[...]}
```

---

## 📃 Architecture & Design Decisions

### Redis Streams

- `flight.search.requested`: Used by main-service to publish requests.
- `flight.search.results`: Used by provider-service to publish search results.

### Server-Sent Events

- Allows real-time push of flight results to clients.
- Efficient for many clients with minimal overhead.

### Design Principles

- UUIDs are used to uniquely track each search.
- Decoupled via Redis Streams for scalability.
- Dockerized for portability and easy testing.

### Trade-offs

- SSE is unidirectional (not bi-directional like WebSockets), but simpler.
- In-memory `sseClients` map can be replaced with pub-sub if scaling horizontally.

---

## ⚠️ Error Handling

- Invalid JSON input returns `400 Bad Request`.
- Redis connection failure during publish/consume is caught and logged. Retries continue in background loops.
- If publishing to `flight.search.requested` fails, the POST returns `500 Internal Server Error`.
- If consuming `flight.search.results` fails, the error is logged and stream resumes.
- SSE returns `404` if `search_id` is not found.
- Optional: Add alert/metrics in future if repeated Redis failures occur.

---

## 📊 Structured Logging

All services use structured logging via `log.Printf` with relevant context, e.g.:

```bash
[SSE SEND] Sending result for: <search_id>
[ERROR] XREADGROUP error: dial tcp ...
[REDIS PUBLISH ERROR] Failed to add search request
```

---

## ✨ Bonus Features

### Clean Architecture

- `handler/`: Fiber handlers
- `service/`: Redis logic
- `mockapi/`: Simulated flight provider

### Redis Consumer Groups

- Uses `XREADGROUP` and `XACK` to ensure reliability and delivery tracking.

### Unit & Integration Tests

- Unit tests can be added to `handler` and `service` layers using `testing` package.
- Example:

```go
func TestHandleSearchRequest(t *testing.T) {
  app := fiber.New()
  app.Post("/api/flights/search", handler.HandleSearchRequest)

  payload := `{"from":"CGK","to":"DPS","date":"2025-07-10","passengers":1}`
  req := httptest.NewRequest("POST", "/api/flights/search", strings.NewReader(payload))
  req.Header.Set("Content-Type", "application/json")
  resp, _ := app.Test(req)

  assert.Equal(t, 200, resp.StatusCode)
}
```

- Integration tests can simulate POST then SSE GET using test clients with Redis mocked.

---

## 💼 Project Structure

```
flight-search/
├── main-service/
│   ├── main.go
│   ├── handler/flight_handler.go
│   ├── service/stream.go
│   └── utils/sse.go (optional helper)
├── provider-service/
│   ├── main.go
│   ├── consumer/redis_consumer.go
│   └── mockapi/flight_provider.go
├── docker-compose.yml
└── README.md
```

---

## 🌟 Authors & License

Built with ❤️ by Oviek Shagya for Backend Engineering Technical Test. MIT License.

