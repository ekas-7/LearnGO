.PHONY: help build run test clean migrate-up migrate-down docker-build docker-run

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building application..."
	@go build -o bin/api cmd/api/main.go
	@echo "Build complete!"

run: ## Run the application
	@echo "Starting application..."
	@go run cmd/api/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -cover ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean

deps: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

lint: fmt vet ## Run linters

dev: ## Run in development mode with hot reload (requires air)
	@air

install-air: ## Install air for hot reload
	@go install github.com/cosmtrek/air@latest

# Database commands
db-start: ## Start PostgreSQL with Docker
	@echo "Starting PostgreSQL with Docker..."
	@docker-compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 3
	@echo "✅ PostgreSQL is ready on localhost:5432"

db-stop: ## Stop PostgreSQL
	@echo "Stopping PostgreSQL..."
	@docker-compose down

db-reset: ## Reset database (removes all data)
	@echo "Resetting database..."
	@docker-compose down -v
	@docker-compose up -d postgres
	@sleep 3
	@echo "✅ Database reset complete"

db-logs: ## View database logs
	@docker logs ecommerce_postgres -f

db-shell: ## Access PostgreSQL shell
	@docker exec -it ecommerce_postgres psql -U postgres -d ecommerce_db

db-status: ## Check database status
	@docker ps --filter "name=ecommerce_postgres" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# Docker commands
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t ecommerce-api .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env ecommerce-api

docker-compose-up: ## Start docker-compose services
	@docker-compose up -d

docker-compose-down: ## Stop docker-compose services
	@docker-compose down

# Production
prod-build: ## Build for production
	@echo "Building for production..."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api cmd/api/main.go

.DEFAULT_GOAL := help
