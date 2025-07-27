# GoPortScanner Makefile
# Provides common development and build tasks

# Variables
BINARY_NAME=goportscanner
MAIN_PATH=cmd/goportscanner/main.go
VERSION?=1.0.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD)
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe
BINARY_DARWIN=$(BINARY_NAME)_darwin

# Default target
.DEFAULT_GOAL := build

# Build the application
.PHONY: build
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PATH)

# Build for current platform with version info
.PHONY: build-version
build-version:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PATH)

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_WINDOWS)
	rm -f $(BINARY_DARWIN)

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -v -cover ./...

# Run tests with race detection
.PHONY: test-race
test-race:
	$(GOTEST) -v -race ./...

# Run benchmarks
.PHONY: benchmark
benchmark:
	$(GOTEST) -bench=. ./...

# Install dependencies
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Update dependencies
.PHONY: deps-update
deps-update:
	$(GOMOD) get -u ./...
	$(GOMOD) tidy

# Run linter
.PHONY: lint
lint:
	golangci-lint run

# Run security checks
.PHONY: security
security:
	gosec ./...

# Format code
.PHONY: fmt
fmt:
	go fmt ./...
	goimports -w .

# Vet code
.PHONY: vet
vet:
	go vet ./...

# Code quality checks
.PHONY: quality
quality: fmt vet lint security

# Build for multiple platforms
.PHONY: build-all
build-all: build-linux build-windows build-darwin

# Build for Linux
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_UNIX) $(MAIN_PATH)

# Build for Windows
.PHONY: build-windows
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_WINDOWS) $(MAIN_PATH)

# Build for macOS
.PHONY: build-darwin
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DARWIN) $(MAIN_PATH)

# Install the binary
.PHONY: install
install:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PATH)
	cp $(BINARY_NAME) /usr/local/bin/

# Uninstall the binary
.PHONY: uninstall
uninstall:
	rm -f /usr/local/bin/$(BINARY_NAME)

# Run the application
.PHONY: run
run:
	$(GOCMD) run $(MAIN_PATH)

# Run with example arguments
.PHONY: run-example
run-example:
	$(GOCMD) run $(MAIN_PATH) -h localhost -s 1 -e 1024

# Generate documentation
.PHONY: docs
docs:
	godoc -http=:6060

# Create release artifacts
.PHONY: release
release: clean build-all
	mkdir -p release
	cp $(BINARY_UNIX) release/$(BINARY_NAME)-linux-amd64
	cp $(BINARY_WINDOWS) release/$(BINARY_NAME)-windows-amd64.exe
	cp $(BINARY_DARWIN) release/$(BINARY_NAME)-darwin-amd64
	cp README.md LICENSE CHANGELOG.md release/

# Create Docker image
.PHONY: docker-build
docker-build:
	docker build -t $(BINARY_NAME):$(VERSION) .

# Run in Docker
.PHONY: docker-run
docker-run:
	docker run --rm $(BINARY_NAME):$(VERSION)

# Development setup
.PHONY: dev-setup
dev-setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install golang.org/x/tools/cmd/godoc@latest

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  lint           - Run linter"
	@echo "  security       - Run security checks"
	@echo "  fmt            - Format code"
	@echo "  quality        - Run all quality checks"
	@echo "  build-all      - Build for all platforms"
	@echo "  install        - Install binary"
	@echo "  run            - Run the application"
	@echo "  release        - Create release artifacts"
	@echo "  help           - Show this help"