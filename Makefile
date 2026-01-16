.PHONY: build build-publisher build-subscriber run test clean fmt lint help docker-up docker-down docker-logs

# Build output directory
BIN_DIR := bin

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOFMT := $(GOCMD) fmt
GOVET := $(GOCMD) vet
GOMOD := $(GOCMD) mod

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## build: Build both publisher and subscriber
build: build-publisher build-subscriber
	@echo "Built publisher and subscriber"

## build-publisher: Build the publisher application
build-publisher:
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) -o $(BIN_DIR)/publisher ./cmd/publisher
	@echo "Built: $(BIN_DIR)/publisher"

## build-subscriber: Build the subscriber application
build-subscriber:
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) -o $(BIN_DIR)/subscriber ./cmd/subscriber
	@echo "Built: $(BIN_DIR)/subscriber"

## test: Run tests
test:
	$(GOTEST) -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## fmt: Format code
fmt:
	$(GOFMT) ./...

## lint: Run linter
lint:
	$(GOVET) ./...

## tidy: Tidy go modules
tidy:
	$(GOMOD) tidy

## clean: Clean build artifacts
clean:
	rm -rf $(BIN_DIR)
	rm -f coverage.out coverage.html

## deps: Download dependencies
deps:
	$(GOMOD) download

## demo: Quick demo (increment counter from both publishers)
demo:
	@echo "=== Publisher A: Incrementing counter ==="
	curl -s -X POST http://localhost:8081/api/counter/increment \
		-H "Content-Type: application/json" \
		-d '{"amount": 5}' | jq .
	@echo ""
	@echo "=== Publisher B: Incrementing counter ==="
	curl -s -X POST http://localhost:8082/api/counter/increment \
		-H "Content-Type: application/json" \
		-d '{"amount": 10}' | jq .
	@echo ""
	@echo "=== Check subscriber logs: docker logs subscriber-app ==="

## docker-up: Start all services with Docker Compose
docker-up:
	docker compose up --build -d
	@echo "Services started. Use 'make docker-logs' to view logs"
	@echo "Publisher A API: http://localhost:8081"
	@echo "Publisher B API: http://localhost:8082"

## docker-down: Stop all Docker services
docker-down:
	docker compose down -v

## docker-logs: View Docker logs
docker-logs:
	docker compose logs -f

## docker-logs-subscriber: View only subscriber logs
docker-logs-subscriber:
	docker logs -f subscriber-app

## docker-build: Build Docker images
docker-build:
	docker compose build
