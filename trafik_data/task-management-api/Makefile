.PHONY: help build run test clean migrate-up migrate-down sqlc swagger docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building application..."
	@go build -o bin/api cmd/api/main.go

run: ## Run the application
	@echo "Running application..."
	@go run cmd/api/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf docs/

migrate-up: ## Run database migrations up
	@./scripts/migrate.sh up

migrate-down: ## Rollback last migration
	@./scripts/migrate.sh down

migrate-reset: ## Reset all migrations
	@./scripts/migrate.sh reset

sqlc: ## Generate SQLC code
	@echo "Generating SQLC code..."
	@sqlc generate

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger docs..."
	@swag init -g cmd/api/main.go -o docs

docker-up: ## Start Docker containers
	@echo "Starting Docker containers..."
	@docker-compose up -d

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	@docker-compose down

install-tools: ## Install required tools
	@echo "Installing tools..."
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

setup: install-tools deps docker-up sqlc swagger ## Complete setup
	@echo "Setup complete!"

dev: docker-up migrate-up sqlc swagger run ## Start development environment
