# Makefile for sort_utility

# Variables
BINARY_NAME=sort_utility
MAIN_PATH=./cmd/main.go
BUILD_DIR=./bin
COVERAGE_DIR=./coverage
COVERAGE_THRESHOLD=80

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet
GOLINT=golint

.PHONY: all build clean test coverage coverage-check deps fmt vet lint check install run help

# Default target
all: check build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -rf $(COVERAGE_DIR)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	$(GOCMD) tool cover -func=$(COVERAGE_DIR)/coverage.out
	@echo "Coverage report generated: $(COVERAGE_DIR)/coverage.html"

# Check coverage threshold
coverage-check: coverage
	@echo "Checking coverage threshold ($(COVERAGE_THRESHOLD)%)..."
	@COVERAGE=$$($(GOCMD) tool cover -func=$(COVERAGE_DIR)/coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$COVERAGE < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage $$COVERAGE% is below threshold $(COVERAGE_THRESHOLD)%"; \
		exit 1; \
	else \
		echo "✅ Coverage $$COVERAGE% meets threshold $(COVERAGE_THRESHOLD)%"; \
	fi

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@which $(GOLINT) > /dev/null || (echo "Installing golint..." && $(GOGET) -u golang.org/x/lint/golint)

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

# Run golint
lint:
	@echo "Running golint..."
	@which $(GOLINT) > /dev/null || (echo "Installing golint..." && $(GOGET) -u golang.org/x/lint/golint)
	$(GOLINT) ./...

# Run all checks (fmt, vet, lint, test)
check: fmt vet lint test

# Run all checks with coverage verification
check-coverage: fmt vet lint coverage-check

# Install the binary
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GOCMD) install $(MAIN_PATH)

# Create example test file
example:
	@echo "Creating example.txt..."
	@echo -e "cherry\napple\nbanana\n3\n1\n2" > example.txt

# Run with example
run: build example
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME) example.txt

# Run with custom arguments
run-args: build
	@echo "Running $(BINARY_NAME) with arguments: $(ARGS)"
	./$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

# Benchmark tests
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Race condition detection
race:
	@echo "Running tests with race detection..."
	$(GOTEST) -race ./...

# Show coverage in terminal
coverage-summary:
	@echo "Coverage Summary:"
	@$(GOCMD) tool cover -func=$(COVERAGE_DIR)/coverage.out 2>/dev/null || echo "Run 'make coverage' first"

# Show help
help:
	@echo "Available targets:"
	@echo "  all             - Run checks and build (default)"
	@echo "  build           - Build the application"
	@echo "  clean           - Clean build artifacts"
	@echo "  test            - Run tests"
	@echo "  coverage        - Run tests with coverage report"
	@echo "  coverage-check  - Verify coverage meets threshold"
	@echo "  deps            - Install dependencies"
	@echo "  fmt             - Format code"
	@echo "  vet             - Run go vet"
	@echo "  lint            - Run golint"
	@echo "  check           - Run all checks (fmt, vet, lint, test)"
	@echo "  check-coverage  - Run all checks with coverage verification"
	@echo "  install         - Install binary"
	@echo "  run             - Build and run with example.txt"
	@echo "  example         - Create example.txt"
	@echo "  bench           - Run benchmark tests"
	@echo "  race            - Run tests with race detection"
	@echo "  coverage-summary - Show coverage summary"
	@echo "  help            - Show this help"
