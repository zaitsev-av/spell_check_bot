.PHONY: build run test clean deps lint

# Build variables
BINARY_NAME=spell_bot
BUILD_DIR=bin
GO_VERSION=1.24

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/bot

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	go run ./cmd/bot

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod verify

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Lint the code
lint:
	@echo "Linting..."
	@if command -v golangci-lint >/dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping..."; \
		echo "Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi


# Build for production with optimizations
build-prod:
	@echo "Building for production..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/bot

# Development build with debug info
build-dev:
	@echo "Building for development..."
	@mkdir -p $(BUILD_DIR)
	go build -race -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/bot

# Show help
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  deps         - Install dependencies"
	@echo "  test         - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean        - Clean build artifacts"
	@echo "  lint         - Lint the code"
	@echo "  fmt          - Format code"
	@echo "  build-prod   - Build optimized production binary"
	@echo "  build-dev    - Build development binary with race detector"
	@echo "  help         - Show this help message"
