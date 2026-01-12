# Makefile
.PHONY: run test migrate docker

run:  ## Run the bot locally
	@echo "Starting bot"
	go run ./cmd/bot

test:  ## Run tests
	@echo "Running tests"
	go test ./... -v

migrate-up:  ## Run database migrations
	@echo "Applying migrations"
	# Здесь будет команда для миграций

docker-up:  ## Start services with Docker Compose
	@echo "Starting Docker services"
	docker-compose up -d

docker-down:  ## Stop Docker services
	@echo "Stopping Docker services"
	docker-compose down
