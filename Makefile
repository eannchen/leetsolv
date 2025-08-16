.PHONY: help prod dev clean build build-all test lint install

# Default target
help:
	@echo "Available targets:"
	@echo "  make dev       - Run in development mode (uses test files)"
	@echo "  make prod      - Run in production mode (uses production files)"
	@echo "  make clean     - Remove all testing data files and build artifacts"
	@echo "  make build     - Build the application for current platform"
	@echo "  make build-all - Build for all supported platforms"
	@echo "  make test      - Run tests with coverage and race detection"
	@echo "  make test-no-race - Run tests without race detection"
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
	@rm -f questions.dev.json deltas.dev.json info.dev.log error.dev.log coverage.html coverage.out
	@rm -rf dist/
	@rm -f leetsolv
	@echo "Clean complete!"

# Build the application for current platform
build:
	@echo "Building leetsolv for current platform..."
	@BUILD_TIME=$$(date -u +"%Y-%m-%dT%H:%M:%SZ"); \
	GIT_COMMIT=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown"); \
	VERSION=$$(git describe --tags --always --dirty 2>/dev/null || echo "dev"); \
	go build -ldflags="-s -w -X main.Version=$$VERSION -X main.BuildTime=$$BUILD_TIME -X main.GitCommit=$$GIT_COMMIT" -o leetsolv
	@echo "Build complete! Binary: leetsolv"

# Build for all supported platforms
build-all:
	@echo "Building leetsolv for all platforms..."
	@mkdir -p dist
	@BUILD_TIME=$$(date -u +"%Y-%m-%dT%H:%M:%SZ"); \
	GIT_COMMIT=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown"); \
	VERSION=$$(git describe --tags --always --dirty 2>/dev/null || echo "dev"); \
	LDFLAGS="-s -w -X main.Version=$$VERSION -X main.BuildTime=$$BUILD_TIME -X main.GitCommit=$$GIT_COMMIT"; \
	GOOS=linux GOARCH=amd64 go build -ldflags="$$LDFLAGS" -o dist/leetsolv-linux-amd64; \
	GOOS=linux GOARCH=arm64 go build -ldflags="$$LDFLAGS" -o dist/leetsolv-linux-arm64; \
	GOOS=darwin GOARCH=amd64 go build -ldflags="$$LDFLAGS" -o dist/leetsolv-darwin-amd64; \
	GOOS=darwin GOARCH=arm64 go build -ldflags="$$LDFLAGS" -o dist/leetsolv-darwin-arm64; \
	GOOS=windows GOARCH=amd64 go build -ldflags="$$LDFLAGS" -o dist/leetsolv-windows-amd64.exe; \
	GOOS=windows GOARCH=arm64 go build -ldflags="$$LDFLAGS" -o dist/leetsolv-windows-arm64.exe
	@echo "Build complete! Check dist/ directory for binaries"

# Run tests with coverage
test:
	@echo "Running tests with coverage..."
	@CGO_ENABLED=1 go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Test coverage report: coverage.html"

# Run tests without race detection (for platforms that don't support it)
test-no-race:
	@echo "Running tests without race detection..."
	@go test -v -coverprofile=coverage.out ./...
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