# Ship Shape - Makefile
# Build, test, and development automation

.PHONY: help build test lint fmt vet coverage clean install run security actionlint

# Binary name
BINARY_NAME=shipshape
BINARY_PATH=./bin/$(BINARY_NAME)

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOINSTALL=$(GOCMD) install

# Build parameters
BUILD_DIR=./bin
CMD_DIR=./cmd/shipshape
PKG_LIST=$(shell go list ./... | grep -v /vendor/)

# Version info (for embedding)
VERSION ?= dev
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

## help: Display this help message
help:
	@echo "Ship Shape - Available Make Targets:"
	@echo ""
	@echo "  make build      - Build the shipshape binary"
	@echo "  make test       - Run all tests with coverage"
	@echo "  make lint       - Run golangci-lint"
	@echo "  make fmt        - Format all Go code"
	@echo "  make vet        - Run go vet"
	@echo "  make coverage   - Generate HTML coverage report"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make install    - Install binary to \$$GOPATH/bin"
	@echo "  make run        - Build and run (usage: make run ARGS='analyze .')"
	@echo "  make deps       - Download and verify dependencies"
	@echo "  make tidy       - Tidy and verify module dependencies"
	@echo "  make security   - Run gosec security scanner"
	@echo "  make actionlint - Validate GitHub Actions workflows"
	@echo "  make check      - Run all quality checks (fmt, vet, lint, test)"
	@echo ""

## build: Build the shipshape binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) $(CMD_DIR)
	@echo "Binary created at $(BINARY_PATH)"

## test: Run all tests with coverage
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@echo "Tests completed"

## lint: Run golangci-lint
lint:
	@echo "Running golangci-lint..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install: https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run --timeout=5m ./...
	@echo "Linting completed"

## fmt: Format all Go code
fmt:
	@echo "Formatting Go code..."
	$(GOFMT) ./...
	@echo "Formatting completed"

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...
	@echo "Vet completed"

## coverage: Generate HTML coverage report
coverage: test
	@echo "Generating coverage report..."
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

## clean: Remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean completed"

## install: Install binary to $GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME) to \$$GOPATH/bin..."
	$(GOINSTALL) $(CMD_DIR)
	@echo "Installation completed"

## run: Build and run (usage: make run ARGS='analyze .')
run: build
	@echo "Running $(BINARY_NAME) $(ARGS)..."
	$(BINARY_PATH) $(ARGS)

## deps: Download and verify dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOGET) -v ./...
	$(GOMOD) verify
	@echo "Dependencies downloaded and verified"

## tidy: Tidy and verify module dependencies
tidy:
	@echo "Tidying module dependencies..."
	$(GOMOD) tidy
	$(GOMOD) verify
	@echo "Module tidy completed"

## security: Run gosec security scanner
security:
	@echo "Running gosec security scanner..."
	@which gosec > /dev/null || (echo "gosec not installed. Install: go install github.com/securego/gosec/v2/cmd/gosec@latest" && exit 1)
	gosec ./...
	@echo "Security scan completed"

## actionlint: Validate GitHub Actions workflows
actionlint:
	@echo "Validating GitHub Actions workflows..."
	@which actionlint > /dev/null || (echo "actionlint not installed. Install: go install github.com/rhysd/actionlint/cmd/actionlint@latest" && exit 1)
	actionlint
	@echo "Workflow validation completed"

## check: Run all quality checks (fmt, vet, lint, test)
check: fmt vet lint test
	@echo "All quality checks passed!"

# Default target
.DEFAULT_GOAL := help
