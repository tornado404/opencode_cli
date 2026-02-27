.PHONY: build test clean run

# Variables
BINARY_NAME=oho
BUILD_DIR=bin
CMD_PATH=./cmd

# Build for current platform
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for Linux (Ubuntu)
build-linux:
	@echo "Building $(BINARY_NAME) for Linux..."
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(CMD_PATH)
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
