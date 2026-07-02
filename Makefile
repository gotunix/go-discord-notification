# Variables
BINARY_NAME=bot-binary
BIN_DIR=bin
SRC_DIR=src
DOCKER_IMAGE=go-discord-notification

.PHONY: all build run test fmt tidy lint clean docker-build docker-up docker-down pre-commit-run help

# Default target
all: build

## build: Build the binary and place it in bin/
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BIN_DIR)
	@cd $(SRC_DIR) && go build -o ../$(BIN_DIR)/$(BINARY_NAME) main.go

## run: Run the application locally
run:
	@echo "Running application..."
	@cd $(SRC_DIR) && go run main.go

## test: Run tests
test:
	@echo "Running tests..."
	@cd $(SRC_DIR) && go test -v ./...

## fmt: Format codebase using gofmt and goimports
fmt:
	@echo "Formatting code..."
	@gofmt -s -w $(SRC_DIR)
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w $(SRC_DIR); \
	else \
		echo "goimports not found, skipping import formatting. Install with: go install golang.org/x/tools/cmd/goimports@latest"; \
	fi

## tidy: Run go mod tidy
tidy:
	@echo "Tidying go modules..."
	@cd $(SRC_DIR) && go mod tidy

## lint: Run golangci-lint
lint:
	@echo "Running golangci-lint..."
	@cd $(SRC_DIR) && golangci-lint run

## clean: Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)

## docker-build: Build the docker image
docker-build:
	@echo "Building docker image..."
	@docker build -t $(DOCKER_IMAGE) .

## docker-up: Start services using docker-compose
docker-up:
	@echo "Starting services..."
	@docker compose up -d

## docker-down: Stop services using docker-compose
docker-down:
	@echo "Stopping services..."
	@docker compose down

## pre-commit-run: Run pre-commit hooks against all files
pre-commit-run:
	@echo "Running pre-commit hooks..."
	@pre-commit run --all-files

## help: Show this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
