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
