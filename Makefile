include .env
export

export PROJECT_ROOT=$(shell pwd)

pg-up:
	docker compose up postgres -d

pg-down:
	docker compose down postgres

pg-erase:
	docker compose down postgres
	rm -rf out/pgdata

migrate-create:
	docker compose run --rm postgres-migrate create -ext sql -dir /migrations -seq init

migrate-up:
	docker compose run --rm postgres-migrate -path /migrations -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:${POSTGRES_DB}?sslmode=disable up

migrate-down:
	docker compose run --rm postgres-migrate -path /migrations -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:${POSTGRES_DB}?sslmode=disable down
