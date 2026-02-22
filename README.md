# Tasks Assignment API (Go + PostgreSQL)

A layered REST API for managing tasks using Go, PostgreSQL, and SQL migrations.

## What Was Done

### Core requirements implemented
- HTTP server on `:8080`
- JSON-only responses
- Correct status codes and `Content-Type: application/json`
- API key middleware (`X-API-KEY`)
- Logging middleware (timestamp + method + endpoint)
- Healthcheck endpoint (`GET /healthz`)
- PostgreSQL connection and auto-migrations on startup
- CRUD for tasks with real DB repository
- Layered architecture:
  - Handler -> Usecase -> Repository

### Additional work implemented
- `.env`-based configuration loading
- Swagger/OpenAPI documentation endpoint
- Unit tests with mocks (handler and usecase layers)
- Dockerization:
  - `Dockerfile`
  - `docker-compose.yml` (app + postgres)
- GitHub Actions CI for unit tests

## Project Structure

```text
cmd/api/main.go
internal/app/main.go
internal/handlers/task.go
internal/handlers/task_test.go
internal/usecase/task.go
internal/usecase/task_test.go
internal/repository/repository.go
internal/repository/_postgres/postgres.go
internal/repository/_postgres/tasks/tasks.go
database/migrations/000001_init.down.sql
database/migrations/000002_init.up.sql
pkg/modules/configs.go
docs/swagger.yaml
.github/workflows/unit-tests.yml
Dockerfile
docker-compose.yml
```

## API Endpoints

- `GET /healthz`
- `GET /v1/tasks`
- `GET /v1/tasks?id={id}`
- `POST /v1/tasks`
- `PATCH /v1/tasks?id={id}`
- `DELETE /v1/tasks?id={id}`
- `GET /v1/external-tasks`
- `GET /swagger`
- `GET /swagger.yaml`

All protected endpoints require:

```http
X-API-KEY: secret12345
```

## Configuration (`.env`)

Copy `.env.example` values into `.env` (already prepared in this project):

```env
APP_PORT=8080
API_KEY=secret12345
DB_HOST=localhost
DB_PORT=5434
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_kbtu
DB_SSLMODE=disable
DB_EXEC_TIMEOUT_SEC=5
```

## Run Locally

### 1) Start PostgreSQL (existing local container)

```bash
docker start db_kbtu
```

### 2) Run app

```bash
go run ./cmd/api
```

## Run With Docker Compose

```bash
docker compose up -d --build
```

Then test:

```bash
curl -i -H "X-API-KEY: secret12345" http://localhost:8080/healthz
```

## Swagger

- UI: [http://localhost:8080/swagger](http://localhost:8080/swagger)
- OpenAPI YAML: [http://localhost:8080/swagger.yaml](http://localhost:8080/swagger.yaml)

## Testing

### Unit tests

```bash
go test ./internal/handlers ./internal/usecase ./internal/middleware -v
```

### Full local test run

```bash
go test ./...
go test ./internal/repository/_postgres -v
```

## CI (GitHub Actions)

Unit tests run automatically on:
- every push
- every pull request

Workflow file:
- `.github/workflows/unit-tests.yml`

