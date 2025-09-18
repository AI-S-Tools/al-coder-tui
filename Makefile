.PHONY: build run install clean test deps

APP_NAME := ai-cli-manager
MAIN_PATH := cmd/$(APP_NAME)/main.go
BUILD_DIR := build

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

# Run the application
run:
	@go run $(MAIN_PATH)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

# Install the application globally
install: build
	@echo "Installing $(APP_NAME) to /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(APP_NAME) /usr/local/bin/
	@echo "Installation complete"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Development mode with auto-reload (requires entr)
dev:
	@echo "Starting development mode..."
	@find . -name "*.go" | entr -r go run $(MAIN_PATH)

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Formatting complete"

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

# Check for vulnerabilities
vuln:
	@echo "Checking for vulnerabilities..."
	@go run golang.org/x/vuln/cmd/govulncheck@latest ./...

help:
	@echo "Available targets:"
	@echo "  make build    - Build the application"
	@echo "  make run      - Run the application"
	@echo "  make install  - Install to /usr/local/bin"
	@echo "  make deps     - Install dependencies"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make test     - Run tests"
	@echo "  make dev      - Run in development mode"
	@echo "  make fmt      - Format code"
	@echo "  make lint     - Run linter"
	@echo "  make vuln     - Check for vulnerabilities"