# Makefile for PureKernels

.PHONY: all build test clean fmt vet lint help

# Default target
all: fmt vet test

# Build verifies that all packages compile
build:
	@echo "Building packages..."
	@go build ./pkg/...

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./pkg/...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./pkg/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	@go test -race -v ./pkg/...

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./pkg/...

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./pkg/...

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1)
	@golangci-lint run ./pkg/...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f coverage.out coverage.html
	@go clean

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# Verify dependencies
verify:
	@echo "Verifying dependencies..."
	@go mod verify

# Check everything before commit
check: fmt vet lint test
	@echo "All checks passed!"

# Install development dependencies
dev-deps:
	@echo "Installing development dependencies..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Help
help:
	@echo "PureKernels Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  make          - Run fmt, vet, and test"
	@echo "  make build    - Verify all packages compile"
	@echo "  make test     - Run all tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make test-race - Run tests with race detector"
	@echo "  make fmt      - Format code"
	@echo "  make vet      - Run go vet"
	@echo "  make lint     - Run linter (requires golangci-lint)"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make tidy     - Tidy go.mod"
	@echo "  make verify   - Verify dependencies"
	@echo "  make check    - Run all checks before commit"
	@echo "  make dev-deps - Install development dependencies"
	@echo "  make help     - Show this help"


