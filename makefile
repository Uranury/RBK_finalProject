.PHONY: help test test-coverage build run clean docker-build docker-up docker-down migrate-up migrate-down wire swagger

# Default target
help:
	@echo "Available commands:"
	@echo "  test          - Run all tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  build         - Build the application"
	@echo "  run           - Run the API server locally"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-build  - Build Docker images"
	@echo "  docker-up     - Start all services with Docker Compose"
	@echo "  docker-down   - Stop all services"
	@echo "  migrate-up    - Apply database migrations"
	@echo "  migrate-down  - Rollback database migrations"
	@echo "  wire          - Generate Wire dependency injection"
	@echo "  swagger       - Generate Swagger documentation"

# Testing
test:
	@echo "Running tests..."
	go test ./... -v

test-coverage:
	@echo "Running tests with coverage..."
	go test ./... -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Building
build:
	@echo "Building application..."
	go build -o bin/api cmd/api/main.go
	go build -o bin/worker cmd/worker/main.go

# Running
run:
	@echo "Starting API server..."
	go run cmd/api/main.go

# Cleaning
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Docker commands
docker-build:
	@echo "Building Docker images..."
	docker-compose build

docker-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d

docker-down:
	@echo "Stopping services..."
	docker-compose down

# Database migrations
migrate-up:
	@echo "Applying database migrations..."
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5436/postgres?sslmode=disable" up

migrate-down:
	@echo "Rolling back database migrations..."
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5436/postgres?sslmode=disable" down

create-mig:
	@echo "Creating new migration..."
	migrate create -ext sql -dir migrations -seq $(name)

swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/api/main.go -o docs

# Run tests in CI environment
ci-test:
	@echo "Running CI tests..."
	go test ./... -v -race -coverprofile=coverage.out
	go tool cover -func=coverage.out