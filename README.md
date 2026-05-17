# go-music-streamer

A Go API service for a music streaming backend.

This service integrates with a companion semantic service project, SemantiTag, for AI tag generation and natural-language song search.

## Stack

- Go 1.25+
- Gin (HTTP router)
- GORM (ORM)
- PostgreSQL 16 (via Docker Compose)
- SemantiTag (FastAPI + ChromaDB) for async tagging and semantic search

## Companion Repository

- SemantiTag repository: https://github.com/shawon-kanji/SematiTag
- SemantiTag docs: ../SemantiTag/README.md

## Project Structure

```text
cmd/api/main.go                 # App entrypoint
internal/api/router/router.go   # Route registration
internal/api/handlers/health.go # HTTP handlers
internal/config/config.go       # Env config loading
internal/database/database.go   # DB connection setup
```

## Prerequisites

- Go installed
- Docker and Docker Compose installed

## Environment Variables

Create a `.env` file in the project root. You can copy values from `.env.example`:

```bash
cp .env.example .env
```

Variables used by the app:

- `DB_HOST` (default: `localhost`)
- `DB_PORT` (default: `5432`)
- `DB_USER` (required)
- `DB_PASSWORD` (required)
- `DB_NAME` (required)
- `DB_SSLMODE` (default: `disable`)
- `APP_PORT` (default: `8080`)

Semantic integration variables used by this service:

- `SEMANTITAG_GENERATE_URL` (default: `http://localhost:8000/generate`)
- `SEMANTITAG_SEARCH_URL` (default: `http://localhost:8000/search`)
- `SEMANTITAG_TENANT_ID` (default: `1234`)
- `SEMANTITAG_SEARCH_TIMEOUT_SEC` (default: `20`)

Optional async tag worker tuning:

- `TAG_GENERATOR_QUEUE_SIZE` (default: `1000`)
- `TAG_GENERATOR_WORKERS` (default: `10`)
- `TAG_GENERATOR_HTTP_TIMEOUT_SEC` (default: `30`)
- `TAG_GENERATOR_MAX_RETRIES` (default: `3`)
- `TAG_GENERATOR_RETRY_BACKOFF_MS` (default: `500`)

## Run with Docker Postgres + Local API

1. Start PostgreSQL:

```bash
make docker-run
```

2. Run the API:

```bash
make run
```

3. Health check:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{"status":"ok"}
```

## Run with SemantiTag (Recommended)

1. Start PostgreSQL for go-music-streamer:

```bash
make docker-run
```

2. Start SemantiTag from its repository:

```bash
cd ../SemantiTag
uv sync
uv run uvicorn app.main:app --reload
```

3. Start go-music-streamer API:

```bash
make run
```

4. Use API endpoints:

- Upload song (triggers async tag generation via SemantiTag)
- Semantic search endpoint: `POST /api/v1/songs/search`

## Common Commands

```bash
make run         # Run API
make build       # Build binary at bin/server
make test        # Run tests
make clean       # Remove bin/
make docker-run  # Start PostgreSQL container
make docker-down # Stop PostgreSQL container
```

## Direct Go Commands

```bash
go run cmd/api/main.go
go build -o bin/server cmd/api/main.go
go test ./...
```

## Notes

- On startup, the app loads variables from `.env` automatically if present.
- The API will fail fast if `DB_USER`, `DB_PASSWORD`, or `DB_NAME` are missing.
