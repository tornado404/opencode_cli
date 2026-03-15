.PHONY: build build-dev build-linux build-linux-version test clean run

# Variables
BINARY_NAME=oho
BUILD_DIR=bin
CMD_PATH=./cmd

# Get version from git tag
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# Build for current platform (with version info)
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@go build -ldflags "-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE)" -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"
	@$(BUILD_DIR)/$(BINARY_NAME) --version

# Build without version info (for development)
build-dev:
	@echo "Building $(BINARY_NAME) (dev mode)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for Linux (Ubuntu)
build-linux:
	@echo "Building $(BINARY_NAME) for Linux..."
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(CMD_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)-linux"

# Build for Linux with version info
build-linux-version:
	@echo "Building $(BINARY_NAME) for Linux $(VERSION)..."
	@GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(CMD_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)-linux"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "Clean complete"

# Run the CLI
run:
	@go run $(CMD_PATH) $(ARGS)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Lint
lint:
	@golangci-lint run ./...

# Format
fmt:
	@go fmt ./...

# Help
help:
	@echo "Available targets:"
	@echo "  build        - Build for current platform"
	@echo "  build-linux  - Build for Linux (Ubuntu)"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  run          - Run the CLI (use ARGS=*)"
	@echo "  deps         - Install dependencies"
	@echo "  lint         - Run linter"
	@echo "  fmt          - Format code"
