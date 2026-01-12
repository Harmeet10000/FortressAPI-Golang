# Makefile for Production-grade Auth Golang Project

# Load environment variables from .env.dev file if it exists
ifneq (,$(wildcard .env.dev))
	include .env.dev
	export
endif

# Variables
APP_NAME=fortress_api
MAIN_PATH=./src/cmd/api
BINARY_PATH=./bin/$(APP_NAME)
MIGRATION_DIR=./src/internal/db/migrations
SEED_DIR=./src/internal/db/seeds

# Database Configuration
# DATABASE_URL is loaded from .env file (see above)
# If not set, you can override with: make target DATABASE_URL="your-url"
DB_URL=$(DATABASE_URL)

# Colors for output
# Using printf with octal escapes (most compatible)
COLOR_RESET=$(shell printf '\033[0m')
COLOR_BLUE=$(shell printf '\033[34m')
COLOR_GREEN=$(shell printf '\033[32m')
COLOR_YELLOW=$(shell printf '\033[33m')

.PHONY: help
help: ## Show this help message
	@echo '$(COLOR_BLUE)Available commands:$(COLOR_RESET)'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(COLOR_GREEN)%-20s$(COLOR_RESET) %s\n", $$1, $$2}'

# ==============================================================================
# Development
# ==============================================================================

.PHONY: run
run: ## Run the application
	@echo "$(COLOR_BLUE)Running application...$(COLOR_RESET)"
	go run $(MAIN_PATH)/main.go

.PHONY: dev
dev: ## Run the application with live reload (requires air)
	@echo "$(COLOR_BLUE)Running with live reload...$(COLOR_RESET)"
	air

.PHONY: build
build: ## Build the application binary
	@echo "$(COLOR_BLUE)Building application...$(COLOR_RESET)"
	@mkdir -p bin
	go build -o $(BINARY_PATH) $(MAIN_PATH)/main.go
	@echo "$(COLOR_GREEN)Binary created at $(BINARY_PATH)$(COLOR_RESET)"

.PHONY: clean
clean: ## Remove build artifacts
	@echo "$(COLOR_BLUE)Cleaning build artifacts...$(COLOR_RESET)"
	rm -rf bin/
	rm -rf tmp/
	@echo "$(COLOR_GREEN)Clean complete$(COLOR_RESET)"

# ==============================================================================
# Dependencies
# ==============================================================================

.PHONY: deps
deps: ## Download Go module dependencies
	@echo "$(COLOR_BLUE)Downloading dependencies...$(COLOR_RESET)"
	go mod download

.PHONY: deps-tidy
deps-tidy: ## Tidy Go module dependencies
	@echo "$(COLOR_BLUE)Tidying dependencies...$(COLOR_RESET)"
	go mod tidy

.PHONY: deps-verify
deps-verify: ## Verify Go module dependencies
	@echo "$(COLOR_BLUE)Verifying dependencies...$(COLOR_RESET)"
	go mod verify

.PHONY: deps-upgrade
deps-upgrade: ## Upgrade all dependencies to latest versions
	@echo "$(COLOR_BLUE)Upgrading dependencies...$(COLOR_RESET)"
	go get -u ./...
	go mod tidy

# ==============================================================================
# Code Quality
# ==============================================================================

.PHONY: fmt
fmt: ## Format Go code
	@echo "$(COLOR_BLUE)Formatting code...$(COLOR_RESET)"
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo "$(COLOR_BLUE)Running go vet...$(COLOR_RESET)"
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint
	@echo "$(COLOR_BLUE)Running golangci-lint...$(COLOR_RESET)"
	golangci-lint run ./...

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with auto-fix
	@echo "$(COLOR_BLUE)Running golangci-lint with auto-fix...$(COLOR_RESET)"
	golangci-lint run --fix ./...

# ==============================================================================
# Testing
# ==============================================================================

.PHONY: test
test: ## Run all tests
	@echo "$(COLOR_BLUE)Running tests...$(COLOR_RESET)"
	go test -v -race ./...

.PHONY: test-cover
test-cover: ## Run tests with coverage
	@echo "$(COLOR_BLUE)Running tests with coverage...$(COLOR_RESET)"
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_GREEN)Coverage report generated: coverage.html$(COLOR_RESET)"

.PHONY: test-integration
test-integration: ## Run integration tests
	@echo "$(COLOR_BLUE)Running integration tests...$(COLOR_RESET)"
	go test -v -race -tags=integration ./...

.PHONY: test-unit
test-unit: ## Run unit tests only
	@echo "$(COLOR_BLUE)Running unit tests...$(COLOR_RESET)"
	go test -v -race -short ./...

.PHONY: benchmark
benchmark: ## Run benchmarks
	@echo "$(COLOR_BLUE)Running benchmarks...$(COLOR_RESET)"
	go test -bench=. -benchmem ./...

# ==============================================================================
# Database - Docker
# ==============================================================================

.PHONY: db-up
db-up: ## Start database container
	@echo "$(COLOR_BLUE)Starting database container...$(COLOR_RESET)"
	docker-compose up -d postgres
	@echo "$(COLOR_GREEN)Database container started$(COLOR_RESET)"

.PHONY: db-down
db-down: ## Stop database container
	@echo "$(COLOR_BLUE)Stopping database container...$(COLOR_RESET)"
	docker-compose down
	@echo "$(COLOR_GREEN)Database container stopped$(COLOR_RESET)"

.PHONY: db-logs
db-logs: ## View database container logs
	docker-compose logs -f postgres

.PHONY: db-shell
db-shell: ## Connect to database shell (psql)
	@echo "$(COLOR_BLUE)Connecting to database...$(COLOR_RESET)"
	docker-compose exec postgres psql -U $(DB_USER) -d $(DB_NAME)

.PHONY: db-reset
db-reset: db-down db-up ## Reset database (stop, start, migrate)
	@echo "$(COLOR_YELLOW)Waiting for database to be ready...$(COLOR_RESET)"
	@sleep 3
	@$(MAKE) migrate-up

# ==============================================================================
# Migrations - Goose
# ==============================================================================

.PHONY: migrate-install
migrate-install: ## Install goose migration tool
	@echo "$(COLOR_BLUE)Installing goose...$(COLOR_RESET)"
	go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "$(COLOR_GREEN)Goose installed$(COLOR_RESET)"

.PHONY: migrate-create
migrate-create: ## Create a new migration (usage: make migrate-create NAME=create_users_table)
	@if [ -z "$(NAME)" ]; then \
		echo "$(COLOR_YELLOW)Usage: make migrate-create NAME=migration_name$(COLOR_RESET)"; \
		exit 1; \
	fi
	@mkdir -p $(MIGRATION_DIR)
	@echo "$(COLOR_BLUE)Creating migration: $(NAME)...$(COLOR_RESET)"
	goose -dir $(MIGRATION_DIR) create $(NAME) sql
	@echo "$(COLOR_GREEN)Migration created$(COLOR_RESET)"

.PHONY: migrate-up
migrate-up: ## Run all pending migrations
	@echo "$(COLOR_BLUE)Running migrations...$(COLOR_RESET)"
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" up
	@echo "$(COLOR_GREEN)Migrations complete$(COLOR_RESET)"

.PHONY: migrate-up-by-one
migrate-up-by-one: ## Run next pending migration
	@echo "$(COLOR_BLUE)Running next migration...$(COLOR_RESET)"
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" up-by-one

.PHONY: migrate-down
migrate-down: ## Rollback last migration
	@echo "$(COLOR_YELLOW)Rolling back last migration...$(COLOR_RESET)"
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" down
	@echo "$(COLOR_GREEN)Rollback complete$(COLOR_RESET)"

.PHONY: migrate-down-to
migrate-down-to: ## Rollback to specific version (usage: make migrate-down-to VERSION=20240101000000)
	@if [ -z "$(VERSION)" ]; then \
		echo "$(COLOR_YELLOW)Usage: make migrate-down-to VERSION=version_number$(COLOR_RESET)"; \
		exit 1; \
	fi
	@echo "$(COLOR_YELLOW)Rolling back to version $(VERSION)...$(COLOR_RESET)"
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" down-to $(VERSION)

.PHONY: migrate-reset
migrate-reset: ## Rollback all migrations
	@echo "$(COLOR_YELLOW)Resetting all migrations...$(COLOR_RESET)"
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" reset
	@echo "$(COLOR_GREEN)Reset complete$(COLOR_RESET)"

.PHONY: migrate-status
migrate-status: ## Show migration status
	@echo "$(COLOR_BLUE)Migration status:$(COLOR_RESET)"
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" status

.PHONY: migrate-version
migrate-version: ## Show current migration version
	@echo "$(COLOR_BLUE)Current version:$(COLOR_RESET)"
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" version

.PHONY: migrate-redo
migrate-redo: ## Rollback and re-run last migration
	@echo "$(COLOR_BLUE)Redoing last migration...$(COLOR_RESET)"
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" redo

.PHONY: migrate-validate
migrate-validate: ## Validate migration files
	@echo "$(COLOR_BLUE)Validating migrations...$(COLOR_RESET)"
	goose -dir $(MIGRATION_DIR) validate

# ==============================================================================
# Database Seeding
# ==============================================================================

.PHONY: seed-create
seed-create: ## Create a new seed file (usage: make seed-create NAME=users)
	@if [ -z "$(NAME)" ]; then \
		echo "$(COLOR_YELLOW)Usage: make seed-create NAME=seed_name$(COLOR_RESET)"; \
		exit 1; \
	fi
	@mkdir -p $(SEED_DIR)
	@echo "$(COLOR_BLUE)Creating seed file: $(NAME).sql...$(COLOR_RESET)"
	@touch $(SEED_DIR)/$(NAME).sql
	@echo "$(COLOR_GREEN)Seed file created at $(SEED_DIR)/$(NAME).sql$(COLOR_RESET)"

.PHONY: seed-run
seed-run: ## Run all seed files
	@echo "$(COLOR_BLUE)Running seed files...$(COLOR_RESET)"
	@for file in $(SEED_DIR)/*.sql; do \
		if [ -f "$$file" ]; then \
			echo "$(COLOR_BLUE)Seeding $$file...$(COLOR_RESET)"; \
			PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -f $$file; \
		fi \
	done
	@echo "$(COLOR_GREEN)Seeding complete$(COLOR_RESET)"

.PHONY: seed-run-file
seed-run-file: ## Run specific seed file (usage: make seed-run-file FILE=users.sql)
	@if [ -z "$(FILE)" ]; then \
		echo "$(COLOR_YELLOW)Usage: make seed-run-file FILE=filename.sql$(COLOR_RESET)"; \
		exit 1; \
	fi
	@echo "$(COLOR_BLUE)Running seed file: $(FILE)...$(COLOR_RESET)"
	PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -f $(SEED_DIR)/$(FILE)
	@echo "$(COLOR_GREEN)Seed file executed$(COLOR_RESET)"

# ==============================================================================
# Database Utilities
# ==============================================================================

.PHONY: db-create
db-create: ## Create database
	@echo "$(COLOR_BLUE)Creating database $(DB_NAME)...$(COLOR_RESET)"
	PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres -c "CREATE DATABASE $(DB_NAME);"
	@echo "$(COLOR_GREEN)Database created$(COLOR_RESET)"

.PHONY: db-drop
db-drop: ## Drop database
	@echo "$(COLOR_YELLOW)Dropping database $(DB_NAME)...$(COLOR_RESET)"
	PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres -c "DROP DATABASE IF EXISTS $(DB_NAME);"
	@echo "$(COLOR_GREEN)Database dropped$(COLOR_RESET)"

.PHONY: db-recreate
db-recreate: db-drop db-create migrate-up ## Drop, create, and migrate database
	@echo "$(COLOR_GREEN)Database recreated and migrated$(COLOR_RESET)"

.PHONY: db-backup
db-backup: ## Backup database to file
	@echo "$(COLOR_BLUE)Backing up database...$(COLOR_RESET)"
	@mkdir -p backups
	PGPASSWORD=$(DB_PASSWORD) pg_dump -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "$(COLOR_GREEN)Backup created in backups/$(COLOR_RESET)"

.PHONY: db-restore
db-restore: ## Restore database from file (usage: make db-restore FILE=backup.sql)
	@if [ -z "$(FILE)" ]; then \
		echo "$(COLOR_YELLOW)Usage: make db-restore FILE=backup_file.sql$(COLOR_RESET)"; \
		exit 1; \
	fi
	@echo "$(COLOR_BLUE)Restoring database from $(FILE)...$(COLOR_RESET)"
	PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) < $(FILE)
	@echo "$(COLOR_GREEN)Database restored$(COLOR_RESET)"

# ==============================================================================
# Docker
# ==============================================================================

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "$(COLOR_BLUE)Building Docker image...$(COLOR_RESET)"
	docker build -t $(APP_NAME):latest .

.PHONY: docker-up
docker-up: ## Start all containers
	@echo "$(COLOR_BLUE)Starting containers...$(COLOR_RESET)"
	docker-compose up -d

.PHONY: docker-down
docker-down: ## Stop all containers
	@echo "$(COLOR_BLUE)Stopping containers...$(COLOR_RESET)"
	docker-compose down

.PHONY: docker-logs
docker-logs: ## View container logs
	docker-compose logs -f

.PHONY: docker-clean
docker-clean: ## Remove all containers, images, and volumes
	@echo "$(COLOR_YELLOW)Cleaning Docker resources...$(COLOR_RESET)"
	docker-compose down -v --rmi all
	@echo "$(COLOR_GREEN)Docker cleanup complete$(COLOR_RESET)"

# ==============================================================================
# Utilities
# ==============================================================================

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "$(COLOR_BLUE)Installing development tools...$(COLOR_RESET)"
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/air-verse/air@latest
	@echo "$(COLOR_GREEN)Tools installed$(COLOR_RESET)"

.PHONY: swagger-gen
swagger-gen: ## Generate Swagger documentation
	@echo "$(COLOR_BLUE)Generating Swagger docs...$(COLOR_RESET)"
	swag init -g $(MAIN_PATH)/main.go -o ./docs
	@echo "$(COLOR_GREEN)Swagger docs generated$(COLOR_RESET)"

.PHONY: mock-gen
mock-gen: ## Generate mocks for testing
	@echo "$(COLOR_BLUE)Generating mocks...$(COLOR_RESET)"
	go generate ./...
	@echo "$(COLOR_GREEN)Mocks generated$(COLOR_RESET)"

.PHONY: check
check: fmt vet lint test ## Run all checks (format, vet, lint, test)
	@echo "$(COLOR_GREEN)All checks passed$(COLOR_RESET)"

.PHONY: setup
setup: deps install-tools migrate-up ## Setup project (deps, tools, migrations)
	@echo "$(COLOR_GREEN)Project setup complete$(COLOR_RESET)"

.DEFAULT_GOAL := help
