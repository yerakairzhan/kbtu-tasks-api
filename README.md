# Practice 4 — Tasks API (Go, PostgreSQL, Docker)

This project is a Go REST API for tasks with PostgreSQL, Docker Compose orchestration, Swagger, and CI for unit tests.

## What was completed

### Core backend
- Layered architecture: `Handler -> Usecase -> Repository`
- PostgreSQL integration via `sqlx`
- Auto migrations on startup
- CRUD endpoints for tasks
- Middleware:
  - API key auth (`X-API-KEY`)
  - request logging
- Healthcheck endpoint (`GET /healthz`)

### Additional implementation
- `.env` configuration loading (`godotenv`)
- Swagger docs:
  - `GET /swagger`
  - `GET /swagger.yaml`
- Unit tests with mocks for handler/usecase
- GitHub Actions CI for unit tests
- Dockerization:
  - multi-stage `Dockerfile`
  - `docker-compose.yml` with `app + db`

## Practice 4 criteria mapping

- Multi-stage build: implemented (`Dockerfile`)
- Compose orchestration: implemented (`docker-compose.yml`)
- Healthchecks + depends_on: implemented (`db` healthcheck + `service_healthy`)
- Named volume persistence: implemented (`pgdata`)
- No hardcoded DB credentials in compose: implemented via `${DB_USER}`, `${DB_PASSWORD}`, `${DB_NAME}`
- DB schema initialization from compose: implemented via `docker/init.sql` mapped to `/docker-entrypoint-initdb.d/init.sql`

## Project structure

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
internal/middleware/auth.go
internal/models/task.go
pkg/modules/configs.go
database/migrations/000001_init.down.sql
database/migrations/000002_init.up.sql
docker/init.sql
docs/swagger.yaml
.github/workflows/unit-tests.yml
Dockerfile
docker-compose.yml
```

## Environment config

Use `.env` (or copy from `.env.example`):

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

## Run locally (without compose)

1. Start local DB container:

```bash
docker start db_kbtu
```

2. Run API:

```bash
go run ./cmd/api
```

## Run with Docker Compose

```bash
docker compose up -d --build
```

Check containers:

```bash
docker ps -a
```

Check DB tables:

```bash
docker compose exec db psql -U postgres -d go_kbtu -c "\dt"
```

## API endpoints

- Public:
  - `GET /healthz`
  - `GET /swagger`
  - `GET /swagger.yaml`
- Protected (`X-API-KEY` required):
  - `GET /v1/tasks`
  - `GET /v1/tasks?id={id}`
  - `POST /v1/tasks`
  - `PATCH /v1/tasks?id={id}`
  - `DELETE /v1/tasks?id={id}`
  - `GET /v1/external-tasks`

## Quick test curls

```bash
curl -i http://localhost:8080/healthz
curl -i http://localhost:8080/swagger
curl -i -H "X-API-KEY: secret12345" http://localhost:8080/v1/tasks
```

## Unit tests

Run all:

```bash
go test ./...
```

Focused unit tests:

```bash
go test ./internal/handlers ./internal/usecase ./internal/middleware -v
```

## CI

GitHub Actions workflow:

- `.github/workflows/unit-tests.yml`

It runs unit tests on push and pull request.

## Docker image size check

```bash
docker build -t tasks-api .
docker images | grep tasks-api
```

(Practice 4 requires showing this in demo.)
