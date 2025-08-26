# PIXELZX POS EVM Chain Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
BINARY_NAME=pixelzx
BINARY_PATH=./bin/$(BINARY_NAME)

# Build target
.PHONY: all build clean test deps run init

# Default target
all: test build

# Build the binary
build:
	@echo "Building PIXELZX POS EVM Chain..."
	@mkdir -p bin
	$(GOBUILD) -o $(BINARY_PATH) -v ./cmd/pixelzx

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf bin/

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Run the node
run: build
	@echo "Running PIXELZX node..."
	$(BINARY_PATH) start

# Initialize genesis
init: build
	@echo "Initializing genesis..."
	$(BINARY_PATH) init

# Install dependencies and build
install: deps build

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Generate documentation
docs:
	@echo "Generating documentation..."
	godoc -http=:6060

# Development commands
dev-deps:
	@echo "Installing development dependencies..."
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t pixelzx/pos:latest .

docker-run:
	@echo "Running Docker container..."
	docker run -p 8545:8545 -p 8546:8546 pixelzx/pos:latest

# Help
help:
	@echo "Available commands:"
	@echo "  build      - Build the binary"
	@echo "  clean      - Clean build artifacts"
	@echo "  test       - Run tests"
	@echo "  deps       - Download dependencies"
	@echo "  run        - Build and run the node"
	@echo "  init       - Initialize genesis"
	@echo "  install    - Install dependencies and build"
	@echo "  fmt        - Format code"
	@echo "  lint       - Lint code"
	@echo "  docs       - Generate documentation"
	@echo "  dev-deps   - Install development dependencies"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run - Run Docker container"
	@echo "  help       - Show this help"