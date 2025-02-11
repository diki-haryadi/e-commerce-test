# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=server

# Docker parameters
DOCKER_COMPOSE=docker-compose
DOCKER=docker

# Services
SERVICES=auth-service product-service order-service shop-service warehouse-service

# Build parameters
BUILD_DIR=build
MIGRATIONS_DIR=migrations

.PHONY: all test clean run deps docker-build docker-push migrate-up migrate-down \
        seed test-integration proto swagger lint help $(SERVICES)

all: clean deps

# Run tests for all services
test:
	@for service in $(SERVICES); do \
		echo "Testing $$service..." ; \
		cd ./$$service && $(GOTEST) -v ./... ; \
		cd .. ; \
	done

# Run integration tests
test-integration:
	@for service in $(SERVICES); do \
		echo "Running integration tests for $$service..." ; \
		cd ./$$service && $(GOTEST) -v -tags=integration ./... ; \
		cd .. ; \
	done

# Clean build directory
clean:
	@rm -rf $(BUILD_DIR)
	@for service in $(SERVICES); do \
		echo "Cleaning $$service..." ; \
		cd ./$$service && rm -f $$service ; \
		cd .. ; \
	done

# Install dependencies
deps:
	@for service in $(SERVICES); do \
		echo "Installing dependencies for $$service..." ; \
		cd ./$$service && $(GOMOD) download && $(GOMOD) tidy ; \
		cd .. ; \
	done

# Docker operations
docker-build:
	@for service in $(SERVICES); do \
		echo "Building Docker image for $$service..." ; \
		$(DOCKER) build -t $$service -f ./$$service/deployments/Dockerfile ./$$service ; \
	done

docker-push:
	@for service in $(SERVICES); do \
		echo "Pushing Docker image for $$service..." ; \
		$(DOCKER) push $$service ; \
	done
# docker-compose  exec product_service ping auth_service
# Docker Compose operations
up:
	$(DOCKER_COMPOSE) up -d

down:
	$(DOCKER_COMPOSE) down

logs:
	$(DOCKER_COMPOSE) logs -f

ps:
	$(DOCKER_COMPOSE) ps

# Database migrations
migrate-up:
	@for service in $(SERVICES); do \
		echo "Running migrations up for $$service..." ; \
		$(DOCKER_COMPOSE) exec $$service go run cmd/migrate/main.go up ; \
	done

migrate-down:
	@for service in $(SERVICES); do \
		echo "Running migrations down for $$service..." ; \
		$(DOCKER_COMPOSE) exec $$service go run cmd/migrate/main.go down ; \
	done

# Database seeding
seed:
	@for service in $(SERVICES); do \
		echo "Seeding database for $$service..." ; \
		$(DOCKER_COMPOSE) exec $$service go run cmd/seeder/main.go ; \
	done

# Generate Swagger documentation
swagger:
	@for service in $(SERVICES); do \
		echo "Generating Swagger docs for $$service..." ; \
		cd ./$$service && swag init -g cmd/main.go ; \
		cd .. ; \
	done

# Code linting
lint:
	@for service in $(SERVICES); do \
		echo "Linting $$service..." ; \
		cd ./$$service && golangci-lint run ; \
		cd .. ; \
	done

# Individual service commands
$(SERVICES):
	@echo "Building $@..."
	@cd ./$@ && $(GOBUILD) -o ../$(BUILD_DIR)/$@ ./cmd/main.go

# Development helpers
dev-deps:
	$(GOCMD) install github.com/swaggo/swag/cmd/swag@latest
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOCMD) install github.com/golang/mock/mockgen@latest

mock:
	@for service in $(SERVICES); do \
		echo "Generating mocks for $$service..." ; \
		cd ./$$service && go generate ./... ; \
		cd .. ; \
	done

# Reset development environment
reset: down clean
	@echo "Removing all containers and volumes..."
	$(DOCKER_COMPOSE) down -v
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "Environment reset complete."

# Help target
help:
	@echo "Available targets:"
	@echo "  build            - Build all services"
	@echo "  test            - Run tests for all services"
	@echo "  test-integration - Run integration tests"
	@echo "  clean           - Clean build artifacts"
	@echo "  deps            - Install dependencies"
	@echo "  docker-build    - Build Docker images"
	@echo "  docker-push     - Push Docker images"
	@echo "  up              - Start all services with Docker Compose"
	@echo "  down            - Stop all services"
	@echo "  logs            - View service logs"
	@echo "  ps              - List running services"
	@echo "  migrate-up      - Run database migrations"
	@echo "  migrate-down    - Rollback database migrations"
	@echo "  seed            - Seed databases"
	@echo "  swagger         - Generate Swagger documentation"
	@echo "  lint            - Run linter"
	@echo "  dev-deps        - Install development dependencies"
	@echo "  mock            - Generate mocks"
	@echo "  reset           - Reset development environment"
	@echo "  help            - Show this help message"

# Default target
.DEFAULT_GOAL := help