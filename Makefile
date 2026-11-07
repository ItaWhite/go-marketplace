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
	docker compose run --rm postgres-migrate -path /migrations -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable up

migrate-down:
	docker compose run --rm postgres-migrate -path /migrations -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable down

protoc-gen:
	protoc -I proto proto/users.proto --go_out=./proto/gen --go_opt=paths=source_relative --go-grpc_out=./proto/gen --go-grpc_opt=paths=source_relative

swagger-api:
	docker compose run --rm swagger init -g main.go -d ./cmd/api,./internal/features/users/handler,./internal/features/products/handler -o docs/api --parseInternal --parseDependency

swagger-auth:
	docker compose run --rm swagger init -g main.go -d ./cmd/auth,./internal/features/auth/handler -o docs/auth --parseInternal --parseDependency