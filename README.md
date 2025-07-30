# ✈️ Flight Search System (Go + Fiber + Redis Streams)

## Overview
Sistem ini mensimulasikan pencarian penerbangan menggunakan arsitektur event-driven dengan Redis Streams dan komunikasi real-time dengan SSE (Server-Sent Events).

## Fitur
- REST API untuk memulai pencarian
- Real-time streaming hasil pencarian via SSE
- Redis Stream sebagai message broker
- Simulasi call ke 3rd-party API

## Endpoints
- `POST /api/flights/search`
- `GET /api/flights/search/:search_id/stream`

## Arsitektur
Lihat [Architecture Diagram](./docs/architecture.png)

## Cara Menjalankan
```bash
docker-compose up --build
