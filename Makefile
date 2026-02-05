.PHONY: help build run test test-integration coverage swagger lint clean install-tools

# Default target
.DEFAULT_GOAL := help

# Binary name
BINARY_NAME=pokemon-api

# Build variables
BUILD_DIR=bin
MAIN_PATH=cmd/api/main.go

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOCLEAN=$(GOCMD) clean

# Swagger
SWAG=swag

help: ## Display this help screen
	@echo "Pokemon REST API - Makefile Commands"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

run: ## Run the application
	@echo "Running $(BINARY_NAME)..."
	@$(GORUN) $(MAIN_PATH)

test: ## Run unit tests
	@echo "Running unit tests..."
	@$(GOTEST) -v -short ./...

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	@$(GOTEST) -v ./test/integration/...

test-all: ## Run all tests (unit + integration)
	@echo "Running all tests..."
	@$(GOTEST) -v ./...

coverage: ## Generate test coverage report
	@echo "Generating coverage report..."
	@$(GOTEST) -v -coverprofile=coverage.out -covermode=atomic ./...
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"
	@$(GOCMD) tool cover -func=coverage.out | grep total | awk '{print "Total Coverage: " $$3}'

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@$(SWAG) init -g $(MAIN_PATH) -o docs/swagger --parseDependency --parseInternal
	@echo "✓ Swagger documentation generated in docs/swagger/"

lint: ## Run golangci-lint
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint is not installed. Install it with:"; \
		echo "  make install-tools"; \
	fi

fmt: ## Format code with gofmt
	@echo "Formatting code..."
	@gofmt -s -w .
	@echo "✓ Code formatted"

vet: ## Run go vet
	@echo "Running go vet..."
	@$(GOCMD) vet ./...
	@echo "✓ go vet passed"

tidy: ## Tidy go.mod dependencies
	@echo "Tidying dependencies..."
	@$(GOMOD) tidy
	@echo "✓ Dependencies tidied"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@$(GOMOD) download
	@echo "✓ Dependencies downloaded"

clean: ## Clean build artifacts and caches
	@echo "Cleaning build artifacts..."
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "✓ Clean complete"

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@echo "Installing swag..."
	@$(GOCMD) install github.com/swaggo/swag/cmd/swag@latest
	@echo "Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin
	@echo "✓ Tools installed"

dev: ## Run in development mode with live reload (requires air)
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "air is not installed. Install it with:"; \
		echo "  go install github.com/cosmtrek/air@latest"; \
		echo ""; \
		echo "Running without live reload..."; \
		$(MAKE) run; \
	fi

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):latest .
	@echo "✓ Docker image built: $(BINARY_NAME):latest"

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env $(BINARY_NAME):latest

check: fmt vet lint test ## Run all checks (format, vet, lint, test)
	@echo "✓ All checks passed"

ci: deps check test-all coverage ## Run CI pipeline
	@echo "✓ CI pipeline complete"

# Development workflow targets
init: deps swagger ## Initialize project (download deps, generate swagger)
	@echo "✓ Project initialized"

all: clean build test ## Clean, build, and test
	@echo "✓ Build and test complete"
