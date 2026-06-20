# URL Shortener

A production-style URL shortener built with Go, Fiber v3, and Redis.

## Architecture

The project follows Clean Architecture principles with clear separation of concerns:

```
┌─────────────────────────────────────────────────┐
│                   Handler (HTTP)                 │
├─────────────────────────────────────────────────┤
│                   Service (Logic)                │
├─────────────────────────────────────────────────┤
│                  Repository (Data)               │
├─────────────────────────────────────────────────┤
│               Redis (Persistence)                │
└─────────────────────────────────────────────────┘
```

Each layer communicates through interfaces, enabling testability and maintainability.

## Folder Structure

```
.
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── config/              # Environment configuration
│   ├── handler/             # HTTP request handlers (Fiber)
│   ├── middleware/           # Rate limiter & logger
│   ├── model/               # Request/response structs
│   ├── repository/          # Data access interface + Redis impl
│   └── service/             # Business logic
├── pkg/
│   └── redis/               # Redis client factory
├── Dockerfile               # Multi-stage Docker build
├── docker-compose.yml       # App + Redis orchestration
├── .env                     # Environment variables
├── .env.example             # Environment variable template
└── README.md
```

## Tech Stack

| Component | Technology   |
| --------- | ------------ |
| Language  | Go 1.26      |
| Web       | Fiber v3     |
| Cache/DB  | Redis 7      |
| Container | Docker       |

## Environment Variables

| Variable             | Default                | Description                     |
| -------------------- | ---------------------- | ------------------------------- |
| `SERVER_PORT`        | `5000`                 | Application listen port         |
| `REDIS_HOST`         | `localhost`            | Redis hostname                  |
| `REDIS_PORT`         | `6379`                 | Redis port                      |
| `BASE_URL`           | `http://localhost:5000` | Base URL for shortened links    |
| `RATE_LIMIT_REQUESTS`| `100`                  | Max requests per IP per window  |
| `RATE_LIMIT_WINDOW`  | `1h`                   | Rate limit window duration      |
| `URL_TTL`            | `24h`                  | Short URL expiration duration   |

## Running Locally

### Prerequisites

- Go 1.24+
- Redis 7+ (running on `localhost:6379`)

### Steps

```bash
# Clone and enter the project
cd myfirstproject

# Copy environment file
cp .env.example .env
# Update REDIS_HOST to "localhost" if running Redis locally

# Install dependencies
go mod tidy

# Run the server
go run ./cmd/main.go
```

### With local Redis via Docker

```bash
# Start only Redis
docker run -d --name redis -p 6379:6379 redis:7-alpine

# Run the app
REDIS_HOST=localhost go run ./cmd/main.go
```

## Running with Docker

```bash
# Build and start both services
docker compose up --build

# Run in detached mode
docker compose up --build -d

# View logs
docker compose logs -f

# Stop everything
docker compose down

# Stop and remove volumes
docker compose down -v
```

The application will be available at `http://localhost:5000`.

## API Examples

### Create Short URL

```bash
curl -X POST http://localhost:5000/api/v1/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'
```

**Response** (201 Created):
```json
{
  "short_code": "aB3xK9m",
  "short_url": "http://localhost:5000/aB3xK9m",
  "expires_in": "24h0m0s"
}
```

### Follow Redirect

```bash
curl -v http://localhost:5000/aB3xK9m
```

**Response**: HTTP 302 redirect to `https://example.com`.

### Non-existent Code

```bash
curl http://localhost:5000/nonexist
```

**Response** (404 Not Found):
```json
{
  "success": false,
  "message": "short code not found"
}
```

### Invalid URL

```bash
curl -X POST http://localhost:5000/api/v1/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "not-a-url"}'
```

**Response** (400 Bad Request):
```json
{
  "success": false,
  "message": "invalid url"
}
```

### Rate Limited

After 100 requests from the same IP within an hour:

```bash
curl -X POST http://localhost:5000/api/v1/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'
```

**Response** (429 Too Many Requests):
```json
{
  "error": "rate limit exceeded"
}
```

## Rate Limiting

The rate limiter applies **only** to `POST /api/v1/shorten` and uses a sliding-window approach backed by Redis:

1. The client IP is used as the rate limit key: `rate_limit:{ip}`
2. `INCR` increments the request counter for that IP
3. On the first request, `EXPIRE` sets the TTL to the configured window (default: 1 hour)
4. If the counter exceeds the limit (default: 100), requests are rejected with HTTP 429
5. The counter automatically resets when the window expires

Configuration via environment variables:

| Variable             | Default |
| -------------------- | ------- |
| `RATE_LIMIT_REQUESTS`| 100     |
| `RATE_LIMIT_WINDOW`  | 1h      |

## Graceful Shutdown

The server handles `SIGINT` and `SIGTERM` signals:
1. Stops accepting new requests
2. Drains in-flight requests (10-second timeout)
3. Closes the Redis connection cleanly

## Design Decisions

- **Short code generation**: 7-character cryptographically random alphanumeric string (62⁷ ≈ 3.5 trillion combinations).
- **No global state**: Dependencies are injected explicitly at startup in `cmd/main.go`.
- **TTL-based expiration**: Short URLs auto-expire after 24 hours via Redis TTL.
- **Minimal Docker image**: Multi-stage build produces a `scratch`-based image (< 15 MB).

## License

MIT
