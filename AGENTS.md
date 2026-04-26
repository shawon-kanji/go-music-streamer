# Go Music Streamer - Copilot Instructions

## Tech Stack
- **Language**: Go 1.25.x
- **Web Framework**: Gin
- **Database**: PostgreSQL via GORM, MongoDB via mongo-driver
- **Validation**: go-playground/validator
- **Configuration**: godotenv

## Architecture Overview
The project follows a standard Clean Architecture / Layered pattern approach in Go:
- [cmd/api/main.go](cmd/api/main.go): Application entrypoint.
- [cmd/seed/main.go](cmd/seed/main.go): Database seeding utility.
- `internal/api/handlers/`: HTTP request handlers and controllers.
- `internal/bootstrap/`: Dependency Injection setup and application wiring (`AppHandlers`).
- `internal/usecase/`: Pure business logic bridging handlers and repositories.
- `internal/repository/`: Data layer, managing interactions with Postgres.
- `internal/domain/`: Core entities (models) and DTOs (Data Transfer Objects).
- `internal/database/`: Database connection singletons and setup.

## Development & Build Commands
Use these standard commands provided by the Makefile:
- Start infrastructure (e.g., DBs): `make docker-run`
- Run the server locally: `make run`
- Run all tests: `make test`
- Build binary: `make build`

## Project Conventions
- **Routing & Handlers**: Define routes in `internal/api/router/`. Handlers should purely map requests to DTOs, invoke usecases, and format responses. Do not put business logic in handlers.
- **Responses & Error Handling**:
  - Always use `internal/framework` helpers (`framework.SendSuccess`, `framework.BadRequest`, etc.) to format JSON responses cleanly in handlers.
  - Return errors from the usecase or repository layer using `internal/domain/apperror` (`apperror.New("CODE", "msg")`) to pass specific error codes up to the client.
  - For struct validation errors, translate them using `framework.FormatValidationError(err)` before responding.
- **Data parsing**: Always use the models in `internal/domain/dto` when parsing API payloads. Follow validation tags for strict boundary checks.
- **Dependency Injection**: Services (repositories, usecases) are initialized and injected into handlers inside the `internal/bootstrap/app.go` container. This `AppHandlers` registry is then passed to the router instance. Do not pollute `main.go` with handler instantiation.

## Documentation
- **Roadmap & Goals**: Read [agent.md](agent.md) for the active phase and feature roadmap.
- **Development**: Reference [README.md](README.md) for standard onboarding details.

## Custom Workflows
- A PR description generator template exists at [agents/pr-description.md](agents/pr-description.md). Follow its rules and format strictly when requested to write a pull request description.
