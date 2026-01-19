# Makefile
.PHONY: run test docker-up docker-down migrate-up migrate-down migrate-create migrate-version

# Подгружаем .env если он есть
-include .env
export

# --- Defaults (на случай если кто-то не указал в .env) ---
DB_HOST_LOCAL ?= localhost
DB_PORT ?= 5432
DB_SSLMODE ?= disable

# DB_URL для миграций, которые запускаются НА ХОСТЕ (поэтому localhost)
DB_URL_LOCAL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST_LOCAL):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

run:  ## Run the bot locally
	@echo "Starting bot locally"
	go run ./cmd/bot

test:  ## Run tests
	@echo "Running tests"
	go test ./... -v

docker-up:  ## Start services with Docker Compose
	@echo "Starting Docker services"
	docker compose up -d

docker-down:  ## Stop Docker services
	@echo "Stopping Docker services"
	docker compose down

migrate-up: ## Apply all up migrations (local migrate)
	@echo "Applying migrations"
	migrate -path migrations -database "$(DB_URL_LOCAL)" up

migrate-down: ## Rollback last migration (local migrate)
	@echo "Rollback last migration"
	migrate -path migrations -database "$(DB_URL_LOCAL)" down 1

migrate-create: ## Create new migration: make migrate-create name=init
	@if [ -z "$(name)" ]; then \
		echo "usage: make migrate-create name=your_migration_name"; \
		exit 1; \
	fi
	@mkdir -p migrations
	migrate create -ext sql -dir migrations -seq $(name)

migrate-version: ## Show migration version
	@echo "Migration version:"
	migrate -path migrations -database "$(DB_URL_LOCAL)" version
