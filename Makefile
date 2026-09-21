.PHONY: help build build-linux run dev swagger migrate seed clean test

APP_NAME=prangibar-go
MAIN_FILE=main.go
ENV_FILE=.env

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the application binary (Windows)
	go build -o $(APP_NAME).exe $(MAIN_FILE)

build-linux: ## Build the application binary for Linux (Ubuntu server)
	set GOOS=linux&& set GOARCH=amd64&& go build -o $(APP_NAME) $(MAIN_FILE)

run: ## Run the application
	go run $(MAIN_FILE)

dev: ## Run the application with hot reload (requires air)
	air

swagger: ## Generate Swagger documentation
	swag init

migrate: ## Run database migrations (AutoMigrate)
	go run $(MAIN_FILE) -migrate

seed: ## Seed initial data (admin & wilayah)
	go run $(MAIN_FILE) -seed

clean: ## Remove build artifacts
	rm -f $(APP_NAME).exe
	rm -f $(APP_NAME)
	rm -rf docs/
	go clean

test: ## Run all tests
	go test ./...

deps: ## Download and tidy dependencies
	go mod tidy

install-tools: ## Install development tools (swag, air)
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/cosmtrek/air@latest

docker-build: ## Build Docker image
	docker build -t $(APP_NAME) .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file $(ENV_FILE) $(APP_NAME)

lint: ## Run linter
	golangci-lint run

fmt: ## Format Go code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

.DEFAULT_GOAL := help
