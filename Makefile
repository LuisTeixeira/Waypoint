.PHONY: up down build logs migrate-up migrate-down migrate-reset

up:
	docker compose up -d

down:
	docker compose down

build:
	docker compose build

logs:
	docker compose logs -f api

migrate-up:
	docker compose run --rm migrate up

migrate-down:
	docker compose run --rm migrate down 1

migrate-reset:
	docker compose run --rm migrate down --all