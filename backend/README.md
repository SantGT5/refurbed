# Refurbed Backend

Go API serving a products catalogue with filtering (search, color, price, bestseller) and in-memory caching.

**Prerequisites:** Go 1.22+

## Quick Start

```bash
make start
```

The server listens on `http://localhost:8080`. Main endpoints:

| Endpoint        | Description                              |
| --------------- | ---------------------------------------- |
| `GET /products` | Product list with optional query filters |
| `GET /health`   | Health check                             |

## Available Commands

| Command             | Description                     |
| ------------------- | ------------------------------- |
| `make start`        | Build and run the application   |
| `make build`        | Compile to `tmp/main`           |
| `make format`       | Format Go code                  |
| `make format-check` | Fail if code is not formatted   |
| `make urls`         | Print backend and frontend URLs |
| `make help`         | List all commands (default)     |

## Production Readiness Improvements

For a production-ready environment, the following improvements would be recommended:

### Configuration & Environment

- **Externalize configuration**: Move hardcoded values (port `8080`, file paths, cache TTL, CORS origins) to environment variables or a config file. Use something like `viper` or `envconfig` for structured config loading.
- **Configurable CORS**: Replace the current “allow all origins” with an explicit whitelist of production domains loaded from config.

### Security

- **Restrict CORS**: Set `allowedOrigins` to only trusted frontend domains instead of `[]string{}` (allow all).
- **Rate limiting**: Add rate limiting (e.g. `golang.org/x/time/rate` or middleware) to protect against abuse and DoS.
- **Input validation**: Validate and sanitize query parameters (search length, price ranges, etc.) to avoid abuse and injection.
- **HTTPS**: Serve over TLS and consider a reverse proxy (nginx, Cloudflare) for termination.

### Observability

- **Structured logging**: Use structured logging (e.g. `zerolog`, `zap`) with log levels, request IDs, and JSON output.
- **Metrics**: Add Prometheus (or similar) metrics for latency, error rates, cache hit/miss, and request counts.
- **Tracing**: Add distributed tracing (OpenTelemetry) for debugging and performance analysis.
- **Health checks**: Split `/health` into liveness (process alive) and readiness (ready to serve, e.g. data loaded, cache warm).

### Reliability

- **Graceful shutdown**: Handle `SIGTERM`/`SIGINT` to drain connections and shut down cleanly.
- **Distributed cache**: Replace in-memory cache with Redis, for example, so multiple instances share cache and survive restarts.
- **Absolute paths**: Use configurable, absolute paths for data files instead of relative paths like `data/products.json`.

### Data Layer

- **Database**: Replace JSON files with a proper database (e.g. PostgreSQL) for scalability, consistency, and querying.
- **Data loading**: Define a clear strategy for data ingestion (ETL, migrations, admin API) instead of merging at startup.

### Testing

- **Unit tests**: Add unit tests for API handlers, cache logic, and merge/filter logic.
- **Integration tests**: Add integration tests for HTTP endpoints with mocked data.
- **Load tests**: Use tools like `k6` or `wrk` to validate performance and stability.

### Deployment

- **Containerization**: Provide a `Dockerfile` and optionally `compose.yml` for consistent builds.
- **CI/CD**: Add a pipeline (e.g. GitHub Actions) for tests, linting, and deployment.
- **API versioning**: Introduce versioned routes (e.g. `/v1/products`) to support backward-compatible changes.

## Notes on Architecture, Decisions & Comments

### Architecture Overview

The application follows a simple layered structure:

- **Entry point** (`main.go`): Wires handlers, middleware, and initializers.
- **API layer** (`api/products.go`): Handles HTTP requests, filtering, and cache interaction.
- **Cache layer** (`cache/`): In-memory cache keyed by query parameters with configurable TTL.
- **Data layer**: JSON files read at runtime; an initializer merges `metadata.json` and `details.json` into `products.json` on first run.

### Design Decisions

- **JSON files over a database**: Keeps the assignment simple and avoids external dependencies. Data is merged once at startup and then read from disk.
- **In-memory cache per query**: Cache keys include all query parameters (search, color, bestseller, minPrice, maxPrice), so different filters hit different cache entries. Trade-off: higher memory use for many distinct queries vs. fast repeated lookups.
- **Standard library only**: Uses `net/http` and no web framework to keep dependencies minimal.
- **CORS allow-all**: Intended for local development and assignment flexibility; production should use an explicit origin whitelist.

### Request Flow

1. Request hits `/products` with optional query params.
2. Cache is checked using a key built from the query string.
3. On miss: read `data/products.json`, filter in memory, marshal response, store in cache.
4. Response is returned as JSON.

### Project Structure

```
api/           - HTTP handlers
cache/         - In-memory cache and cache key logic
data/          - JSON data files
initializers/  - Startup data merge (metadata + details → products)
middleware/    - CORS
```
