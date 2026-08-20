# URL Shortener

URL shortener with semantic search, built with Go, Redis/Valkey, and Ollama.

## Features

- Shorten URLs with optional TTL and tags
- Redirect with click counting
- Semantic vector search (Ollama + RediSearch)
- Leaderboard, recent links, stats
- React frontend

## Redis Features Used

Strings, Hashes, Lists, Sets, Sorted Sets, Pub/Sub, Streams, HyperLogLog, Bitmaps, Geospatial, Transactions, Pipelining, Lua Scripting, Vector Search

## Prerequisites

- Go 1.21+
- Docker
- Ollama
- Node.js (for frontend)

## Setup

```bash
# Start Redis Stack and pull embedding model
docker run -d --name redis-stack -p 6379:6379 redis/redis-stack-server
ollama pull nomic-embed-text

# Install Go dependencies
go mod tidy

# Start backend
go run .

# Start frontend (separate terminal)
cd url-shortener-frontend
npm install
npm run dev
```

Backend: `http://localhost:3000`
Frontend: `http://localhost:5173`

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /shorten | Create short link |
| GET | /:code | Redirect to original URL |
| GET | /stats/:code | Link statistics |
| GET | /recent | Recent links |
| GET | /top | Most clicked (leaderboard) |
| GET | /tags/:tag | Links by tag |
| GET | /search?q= | Semantic search |
| GET | /active/:code | Active days (Bitmaps) |
| GET | /locations/:code | Click locations (Geo) |
| GET | /events | Event log (Streams) |
