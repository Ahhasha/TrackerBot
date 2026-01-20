# Makefile
.PHONY: run test docker-up docker-down migrate-up migrate-down migrate-create migrate-version

-include .env
export

DB_HOST_LOCAL ?= localhost
DB_PORT ?= 5432
DB_SSLMODE ?= disable

DB_URL_LOCAL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST_LOCAL):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

run:
	go run ./cmd/bot

test:
	go test ./... -v

docker-up:
	docker compose up -d

docker-down:
	docker compose down

migrate-up:
	migrate -path migrations -database "$(DB_URL_LOCAL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL_LOCAL)" down 1

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "usage: make migrate-create name=your_migration_name"; \
		exit 1; \
	fi
	@mkdir -p migrations
	migrate create -ext sql -dir migrations -seq $(name)

migrate-version:
	migrate -path migrations -database "$(DB_URL_LOCAL)" version
