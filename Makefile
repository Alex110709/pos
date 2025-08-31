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

# Docker Hub configuration
DOCKER_REGISTRY=yuchanshin
DOCKER_IMAGE=pixelzx-evm
DOCKER_TAG=latest

# Get version from git
VERSION=$(shell git describe --tags --always --dirty)

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

# Docker Hub commands
docker-build-hub:
	@echo "Building Docker image for Docker Hub..."
	docker build -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG) .
	docker build -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(VERSION) .

docker-push-hub:
	@echo "Pushing Docker image to Docker Hub..."
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(VERSION)

docker-login:
	@echo "Logging into Docker Hub..."
	docker login

docker-deploy-hub: docker-build-hub docker-push-hub
	@echo "Docker Hub deployment completed!"
	docker images | grep $(DOCKER_REGISTRY)/$(DOCKER_IMAGE)

docker-test-hub:
	@echo "Testing Docker Hub image..."
	docker run --rm $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG) pixelzx version

docker-build-dev:
	@echo "Building development Docker image..."
	docker build -f Dockerfile.dev -t pixelzx/pos:dev .

docker-run:
	@echo "Running Docker container..."
	docker run -d --name pixelzx-node \
		-p 8545:8545 -p 8546:8546 -p 30303:30303 -p 6060:6060 \
		-v pixelzx-data:/app/data \
		-v pixelzx-keystore:/app/keystore \
		-v pixelzx-logs:/app/logs \
		pixelzx/pos:latest

docker-stop:
	@echo "Stopping Docker container..."
	docker stop pixelzx-node || true
	docker rm pixelzx-node || true

docker-clean:
	@echo "Cleaning Docker resources..."
	docker stop pixelzx-node || true
	docker rm pixelzx-node || true
	docker rmi pixelzx/pos:latest || true
	docker rmi pixelzx/pos:dev || true

# Docker Compose commands
compose-up:
	@echo "Starting production environment with Docker Compose..."
	docker-compose up -d

compose-down:
	@echo "Stopping production environment..."
	docker-compose down

compose-logs:
	@echo "Showing Docker Compose logs..."
	docker-compose logs -f

compose-dev-up:
	@echo "Starting development environment with Docker Compose..."
	docker-compose -f docker-compose.dev.yml up -d

compose-dev-down:
	@echo "Stopping development environment..."
	docker-compose -f docker-compose.dev.yml down

compose-dev-logs:
	@echo "Showing development environment logs..."
	docker-compose -f docker-compose.dev.yml logs -f

# Production deployment
deploy-prod: docker-build
	@echo "Deploying to production..."
	docker-compose down || true
	docker-compose up -d
	@echo "Production deployment completed!"

# Development environment
dev-setup: docker-build-dev
	@echo "Setting up development environment..."
	mkdir -p data keystore logs
	docker-compose -f docker-compose.dev.yml up -d
	@echo "Development environment ready!"

# Initialize production chain
init-prod:
	@echo "Initializing production chain..."
	docker run --rm \
		-v pixelzx-data:/app/data \
		-v pixelzx-keystore:/app/keystore \
		pixelzx/pos:latest \
		pixelzx init pixelzx-mainnet --chain-id 8888 --datadir /app/data

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
	@echo ""
	@echo "Docker commands:"
	@echo "  docker-build     - Build Docker image"
	@echo "  docker-build-dev - Build development Docker image"
	@echo "  docker-run       - Run Docker container"
	@echo "  docker-stop      - Stop Docker container"
	@echo "  docker-clean     - Clean Docker resources"
	@echo "  docker-build-hub - Build Docker image for Docker Hub"
	@echo "  docker-push-hub  - Push Docker image to Docker Hub"
	@echo "  docker-login     - Login to Docker Hub"
	@echo "  docker-deploy-hub- Build and push to Docker Hub"
	@echo "  docker-test-hub  - Test Docker Hub image"
	@echo ""
	@echo "Docker Compose commands:"
	@echo "  compose-up       - Start production environment"
	@echo "  compose-down     - Stop production environment"
	@echo "  compose-logs     - Show production logs"
	@echo "  compose-dev-up   - Start development environment"
	@echo "  compose-dev-down - Stop development environment"
	@echo "  compose-dev-logs - Show development logs"
	@echo ""
	@echo "Deployment commands:"
	@echo "  deploy-prod      - Deploy to production"
	@echo "  dev-setup        - Setup development environment"
	@echo "  init-prod        - Initialize production chain"
	@echo "  help             - Show this help"