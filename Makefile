# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Create SQLite DB container with Docker
docker-run:
	@echo "Running with SQLite and Docker..."
	@docker compose -f docker-compose-sqlite.yml up --build || \
		(echo "Falling back to Docker Compose V1" && \
		docker-compose -f docker-compose-sqlite.yml up --build)

# Shutdown SQLite DB container with Docker
docker-down:
	@docker compose -f docker-compose-sqlite.yml down || \
		(echo "Falling back to Docker Compose V1" && \
		docker-compose -f docker-compose-sqlite.yml down)

# Create SQLite DB container with Podman
podman-run:
	@echo "Running with SQLite and Podman..."
	@podman-compose -f docker-compose-sqlite.yml --env-file ./.env up --build || \
	(echo "podman not found" && exit 1)

# Shutdown SQLite DB container with Podman
podman-down:
	@echo "Stopping Postgres Podman container..."
	@podman-compose -f docker-compose-sqlite.yml down || \
	(echo "podman not found" && exit 1)

# Create Postgres DB container with Docker
docker-run-postgres:
	@echo "Running Postgres with Docker..."
	@docker compose -f docker-compose-postgres.yml up --build || \
		(echo "Falling back to Docker Compose V1" && \
		docker-compose -f docker-compose-postgres.yml up --build)

# Shutdown Postgres DB container with Docker
docker-down-postgres:
	@echo "Stopping Postgres Docker container..."
	@docker compose -f docker-compose-postgres.yml down || \
		(echo "Falling back to Docker Compose V1" && \
		docker-compose -f docker-compose-postgres.yml down)

# Create Postgres DB container with Podman
podman-run-postgres:
	@echo "Running Postgres with Podman..."
	@podman-compose -f docker-compose-postgres.yml --env-file ./.env up --build || \
	(echo "podman not found" && exit 1)

# Shutdown Postgres DB container with Podman
podman-down-postgres:
	@echo "Stopping Postgres Podman container..."
	@podman-compose -f docker-compose-postgres.yml down || \
	(echo "podman not found" && exit 1)

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main *.db
	@rm -rf tmp

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: all build run test clean watch docker-run docker-down podman-run podman-down docker-run-postgres docker-down-postgres podman-run-postgres podman-down-postgres itest
