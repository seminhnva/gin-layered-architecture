# gin-layered-architecture

A reusable Go + Gin backend starter organized around layered architecture. The goal of this repository is to provide a clean base project that can grow into real-world APIs without needing an early rewrite.

## Current Scope

- Clear layering: `handler -> service -> repository`
- Environment-based configuration with validation
- Shared API response and application error handling
- DTO-based request validation
- Core middleware: request ID, logging, recovery, CORS
- CORS allow-list configured from environment variables
- Graceful HTTP server shutdown
- Health endpoints: `/health`, `/ready`
- Rotating file logs with `lumberjack`
- In-memory user CRUD module as a reference implementation

## Project Structure

- `cmd/api`: application entrypoint
- `internal/app`: application bootstrap and module wiring
- `internal/api/handler`: HTTP handlers
- `internal/service`: business logic layer
- `internal/repository`: data access layer
- `internal/routes`: route registration
- `internal/config`: environment configuration
- `internal/common`: shared response and error packages
- `internal/dto`: request and response DTOs
- `pkg/logger`: logger setup and rotating file log support
- `docs`: roadmap and project notes

## Run Locally

1. Copy `.env.example` to `.env` if you want to customize local configuration.
2. Run:

```bash
go run ./cmd/api
```

The application runs on `http://localhost:8080` by default.

## Environment Variables

Important variables from `.env.example`:

- `JWT_SECRET`: required secret used for future auth-related work
- `CORS_ALLOWED_ORIGINS`: comma-separated allow-list such as `http://localhost:3000,http://localhost:5173`
- `LOG_FILE_PATH`: log file path, defaults to `logs/app.log`
- `LOG_MAX_SIZE_MB`, `LOG_MAX_BACKUPS`, `LOG_MAX_AGE_DAYS`, `LOG_COMPRESS`: log rotation settings

## Example Endpoints

- `GET /health`
- `GET /ready`
- `GET /api/v1/users`
- `POST /api/v1/users`
- `GET /api/v1/users/:uuid`
- `PUT /api/v1/users/:uuid`
- `DELETE /api/v1/users/:uuid`

## Logging

Logs are written to both:

- `stdout` for local development and container-friendly output
- a rotating file managed by `lumberjack`

The default log file is `logs/app.log`.

## Next Milestone

The next major step is Phase 2 of the roadmap:

- add a real database layer
- introduce migrations
- replace the in-memory repository with persistent storage
