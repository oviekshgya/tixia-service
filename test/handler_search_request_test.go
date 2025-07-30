package test

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strings"
	"testing"
	"tixia-service/main-service/handler"
)

func TestHandleSearchRequest(t *testing.T) {
	app := fiber.New()
	app.Post("/api/flights/search", handler.HandleSearchRequest)

	payload := `{"from":"CGK","to":"DPS","date":"2025-07-10","passengers":1}`
	req := httptest.NewRequest("POST", "/api/flights/search", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}
