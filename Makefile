.PHONY: help run test unit docker-up docker-down docker-restart logs ps db


run:
	go run ./cmd/api

test:
	go test ./...

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

docker-restart:
	docker compose down
	docker compose up -d --build

db:
	docker compose exec db psql -U postgres -d go_kbtu

logs:
	docker compose logs -f

ps:
	docker ps -a
