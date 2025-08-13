.PHONY: help prod dev clean build build-all test lint install

# Default target
help:
	@echo "Available targets:"
	@echo "  make dev       - Run in development mode (uses test files)"
	@echo "  make prod      - Run in production mode (uses production files)"
	@echo "  make clean     - Remove all testing data files and build artifacts"
	@echo "  make build     - Build the application for current platform"
	@echo "  make build-all - Build for all supported platforms"
	@echo "  make test      - Run tests with coverage"
	@echo "  make lint      - Run linting and formatting"
	@echo "  make install   - Install the application locally"
	@echo "  make help      - Show this help"

# Development mode - uses test files
dev:
	@echo "Running in TEST mode..."
	@./run_dev.sh

# Production mode - uses production files
prod:
	@echo "Running in PRODUCTION mode..."
	@./run_prod.sh

# Clean all data files and build artifacts
clean:
	@echo "Removing all testing data files and build artifacts..."
	@rm -f questions.test.json deltas.test.json info.test.log error.test.log
	@rm -rf dist/
	@rm -f leetsolv
	@echo "Clean complete!"

# Build the application for current platform
build:
	@echo "Building leetsolv for current platform..."
	@go build -ldflags="-s -w" -o leetsolv
	@echo "Build complete! Binary: leetsolv"

# Build for all supported platforms
build-all:
	@echo "Building leetsolv for all platforms..."
	@mkdir -p dist
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/leetsolv-linux-amd64
	@GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/leetsolv-linux-arm64
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/leetsolv-darwin-amd64
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/leetsolv-darwin-arm64
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/leetsolv-windows-amd64.exe
	@GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o dist/leetsolv-windows-arm64.exe
	@echo "Build complete! Check dist/ directory for binaries"

# Run tests with coverage
test:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Test coverage report: coverage.html"

# Run linting and formatting
lint:
	@echo "Running linting and formatting..."
	@go fmt ./...
	@go vet ./...
	@echo "Linting complete!"

# Install the application locally
install:
	@echo "Installing leetsolv locally..."
	@go install
	@echo "Installation complete! Run 'leetsolv' from anywhere"