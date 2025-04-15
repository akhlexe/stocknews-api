.PHONY: run test test-verbose clean build deps tidy lint coverage

# Default binary output path
BINARY_NAME=stocknews-api
BUILD_DIR=./bin
# Main application entry point
run:
	@echo "🚀 Running the application..."
	go run ./cmd/server/main.go

# Build the application
build:
	@echo ">> Building the application..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BINARY_NAME) ./cmd/server
	@echo "[OK] Build complete!"

# Install dependencies
deps:
	@echo ">> Installing dependencies..."
	go mod download
	@echo "[OK] Dependencies installed!"

# Clean up binaries and test cache
clean:
	@echo ">> Cleaning up..."
	rm -rf $(BUILD_DIR) $(BINARY_NAME)
	go clean -testcache
	@echo "[OK] Clean complete!"	

# Run tests
test:
	@echo ">> Running tests..."
	go test ./...
	@echo "[OK] Tests complete!"


# Run tests with verbose output
test-verbose:
	@echo ">> Running tests with verbose output..."
	go test -v ./...
	@echo "[OK] Tests complete!"

# Format code
fmt:
	@echo ">> Formatting code..."
	go fmt ./...
	@echo "[OK] Code formatted!"

# Run go mod tidy
tidy:
	@echo ">> Running go mod tidy..."
	go mod tidy
	@echo "[OK] go mod tidy complete!"

# Run code coverage
coverage:
	@echo ">> Running code coverage..."
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo ">> Code coverage report generated: coverage.html"
	@echo "[OK] Code coverage complete!"

# Run specific test
test-file:
	@echo ">> Running tests in $(FILE)..."
	go test -v $(FILE)

# Help command
help:
	@echo "Available commands:"
	@echo "  run              Run the application"
	@echo "  build            Build the application"
	@echo "  deps             Install dependencies"
	@echo "  clean            Clean up binaries and test cache"
	@echo "  test             Run all tests"
	@echo "  test-verbose     Run all tests with verbose output"
	@echo "  test-file FILE=path/to/file  Run tests in specific file"
	@echo "  fmt              Format code"
	@echo "  tidy             Run go mod tidy"
	@echo "  coverage         Generate test coverage report"
