# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=awards-management
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe

.PHONY: all build clean test deps run help

# Default target
all: test build

# Build the binary
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_WINDOWS)

# Run tests
test:
	$(GOTEST) -v ./...

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run the application
run:
	$(GOCMD) run main.go

# Run with live reload (requires air)
dev:
	air

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

# Build for Windows
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WINDOWS) -v

# Cross compilation
build-all: build-linux build-windows

# Install air for live reload
install-air:
	$(GOCMD) install github.com/cosmtrek/air@latest

# Format code
fmt:
	$(GOCMD) fmt ./...

# Vet code
vet:
	$(GOCMD) vet ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Install linter
install-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2

# Database operations
db-create:
	mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS awards_management CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

db-drop:
	mysql -u root -p -e "DROP DATABASE IF EXISTS awards_management;"

db-reset: db-drop db-create

# Docker operations
docker-build:
	docker build -t awards-management .

docker-run:
	docker run -p 8000:8000 awards-management

# Help
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  clean        - Clean build files"
	@echo "  test         - Run tests"
	@echo "  deps         - Download dependencies"
	@echo "  run          - Run the application"
	@echo "  dev          - Run with live reload (requires air)"
	@echo "  build-linux  - Build for Linux"
	@echo "  build-windows- Build for Windows"
	@echo "  build-all    - Cross compile for all platforms"
	@echo "  install-air  - Install air for live reload"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code"
	@echo "  lint         - Lint code (requires golangci-lint)"
	@echo "  install-lint - Install golangci-lint"
	@echo "  db-create    - Create database"
	@echo "  db-drop      - Drop database"
	@echo "  db-reset     - Reset database"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  help         - Show this help"