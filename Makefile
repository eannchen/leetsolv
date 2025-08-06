.PHONY: help prod dev clean

# Default target
help:
	@echo "Available targets:"
	@echo "  make dev   - Run in development mode (uses test files)"
	@echo "  make prod  - Run in production mode (uses production files)"
	@echo "  make clean - Remove all testing data files"
	@echo "  make help  - Show this help"

# Development mode - uses test files
dev:
	@echo "Running in TEST mode..."
	@./run_dev.sh

# Production mode - uses production files
prod:
	@echo "Running in PRODUCTION mode..."
	@./run_prod.sh

# Clean all data files
clean:
	@echo "Removing all testing data files..."
	@rm -f questions.test.json deltas.test.json info.test.log error.test.log
	@echo "Clean complete!"