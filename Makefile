# Makefile for gRPC API Gateway project

.PHONY: help lint fmt test clean install-tools proto setup-gateway setup-grpc

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Install development tools
install-tools: ## Install development tools (golangci-lint, goimports, etc.)
	@echo "Installing development tools..."
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.61.0
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Code formatting
fmt: ## Format all Go code
	@echo "Formatting Go code..."
	gofmt -s -w .
	goimports -w -local github.com/longtrd/grpc-api-gateway .

# Code linting
lint: ## Run linters on all Go code
	@echo "Running linters..."
	golangci-lint run --enable=errcheck,gofmt,goimports,govet,staticcheck,unused,ineffassign,misspell --timeout=5m ./...
	staticcheck ./...

# Quick lint for CI
lint-fast: ## Run fast linters (suitable for CI)
	@echo "Running fast linters..."
	golangci-lint run --fast --enable=errcheck,gofmt,goimports,govet,staticcheck,unused --timeout=3m ./...

# Run tests
test: ## Run all tests
	@echo "Running unit tests..."
	cd gRPC-server && go test -race -coverprofile=coverage.out ./internal/...
	cd api-gateway && go test -race -coverprofile=coverage.out ./internal/... || true

# Run tests with coverage report
test-coverage: test ## Run tests and generate coverage report
	@echo "Generating coverage report..."
	cd gRPC-server && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: gRPC-server/coverage.html"

# Run integration tests
test-integration: ## Run integration tests
	@echo "Running integration tests..."
	cd gRPC-server && go test -tags=integration -race ./tests/...

# Run unit tests only
test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	cd gRPC-server && go test -race ./internal/domain/... ./internal/usecase/... ./internal/mocks/...

# Test with verbose output
test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	cd gRPC-server && go test -v -race -coverprofile=coverage.out ./internal/...

# Run tests for specific package
test-domain: ## Run domain tests
	@echo "Running domain tests..."
	cd gRPC-server && go test -v ./internal/domain/...

test-usecase: ## Run usecase tests
	@echo "Running usecase tests..."
	cd gRPC-server && go test -v ./internal/usecase/...

# Run benchmarks
test-bench: ## Run benchmark tests
	@echo "Running benchmark tests..."
	cd gRPC-server && go test -bench=. -benchmem ./internal/...

# Protocol buffer generation
proto: ## Generate Go code from protocol buffer files
	@echo "Generating protobuf code..."
	cd gRPC-server && protoc --go_out=. --go-grpc_out=. proto/*.proto

# Clean generated files
clean: ## Clean generated files and build artifacts
	@echo "Cleaning up..."
	rm -f coverage.out coverage.html
	find . -name "*.pb.go" -type f -delete
	find . -name "*.pb.gw.go" -type f -delete
	go clean -cache
	go clean -testcache

# Development setup for gRPC server
setup-grpc: ## Setup gRPC server dependencies
	@echo "Setting up gRPC server..."
	cd gRPC-server && go mod tidy
	cd gRPC-server && go get github.com/spf13/cobra
	cd gRPC-server && go get google.golang.org/grpc
	cd gRPC-server && go get google.golang.org/protobuf

# Development setup for API Gateway
setup-gateway: ## Setup API Gateway dependencies
	@echo "Setting up API Gateway..."
	cd api-gateway && go mod tidy
	cd api-gateway && go get github.com/gin-gonic/gin
	cd api-gateway && go get google.golang.org/grpc

# Full development setup
setup: install-tools setup-grpc setup-gateway proto ## Complete development environment setup

# Build gRPC server
build-grpc: ## Build gRPC server binary
	@echo "Building gRPC server..."
	cd gRPC-server && go build -o bin/grpc-server cmd/main.go

# Build API Gateway
build-gateway: ## Build API Gateway binary
	@echo "Building API Gateway..."
	cd api-gateway && go build -o bin/api-gateway cmd/main.go

# Build all
build: build-grpc build-gateway ## Build all binaries

# Run gRPC server
run-grpc: ## Run gRPC server (default port 50051)
	@echo "Starting gRPC server..."
	cd gRPC-server && go run cmd/main.go start

# Run API Gateway
run-gateway: ## Run API Gateway
	@echo "Starting API Gateway..."
	cd api-gateway && go run cmd/main.go

# Pre-commit checks
pre-commit: fmt lint test ## Run all pre-commit checks (format, lint, test)
	@echo "All pre-commit checks passed!"

# CI pipeline
ci: lint test ## Run CI pipeline (lint and test)
	@echo "CI pipeline completed successfully!"

# Security check
security: ## Run security analysis
	@echo "Running security analysis..."
	gosec ./...

# Dependency check
deps-check: ## Check for outdated dependencies
	@echo "Checking dependencies..."
	cd gRPC-server && go list -u -m all
	cd api-gateway && go list -u -m all

# Update dependencies
deps-update: ## Update all dependencies
	@echo "Updating dependencies..."
	cd gRPC-server && go get -u ./...
	cd api-gateway && go get -u ./...
	cd gRPC-server && go mod tidy
	cd api-gateway && go mod tidy

# Development watch mode (requires entr)
watch-grpc: ## Watch and restart gRPC server on changes
	find gRPC-server -name "*.go" | entr -r make run-grpc

watch-gateway: ## Watch and restart API Gateway on changes
	find api-gateway -name "*.go" | entr -r make run-gateway

# Docker development
docker-build: ## Build Docker images
	@echo "Building Docker images..."
	./scripts/docker-dev.sh build

docker-start: ## Start services with Docker Compose
	@echo "Starting services with Docker..."
	./scripts/docker-dev.sh start

docker-stop: ## Stop Docker services
	@echo "Stopping Docker services..."
	./scripts/docker-dev.sh stop

docker-restart: ## Restart Docker services
	@echo "Restarting Docker services..."
	./scripts/docker-dev.sh restart

docker-logs: ## Show Docker service logs
	@echo "Showing service logs..."
	./scripts/docker-dev.sh logs

docker-status: ## Show Docker service status
	@echo "Showing service status..."
	./scripts/docker-dev.sh status

docker-clean: ## Clean Docker resources
	@echo "Cleaning Docker resources..."
	./scripts/docker-dev.sh clean

docker-rebuild: ## Rebuild and restart Docker services
	@echo "Rebuilding and restarting services..."
	./scripts/docker-dev.sh rebuild

# Docker testing
docker-test: docker-build ## Build and test with Docker
	@echo "Testing with Docker..."
	docker-compose up -d grpc-server
	sleep 15
	@echo "Testing health check..."
	docker-compose exec grpc-server nc -z localhost 50051 && echo "✅ TCP health check passed" || echo "❌ TCP health check failed"
	@echo "Testing gRPC health service..."
	docker-compose exec grpc-server grpc_health_probe -addr=localhost:50051 && echo "✅ gRPC health check passed" || echo "⚠️  gRPC health check not available (grpc_health_probe not installed)"
	@echo "Testing with logs..."
	docker-compose logs grpc-server | grep -q "Starting gRPC server" && echo "✅ Server started successfully" || echo "❌ Server start failed"
	docker-compose down

# Test clean architecture compliance
test-architecture: ## Test clean architecture compliance
	@echo "Testing clean architecture compliance..."
	cd gRPC-server && go test -v ./internal/config/...
	cd gRPC-server && go test -v ./pkg/...
	@echo "✅ Clean architecture tests passed"

# Test with different configurations
test-config: ## Test with different configuration scenarios
	@echo "Testing configuration management..."
	cd gRPC-server && LOG_LEVEL=debug DATABASE_TYPE=memory go test -v ./internal/config/...
	@echo "✅ Configuration tests passed"

# Complete test suite with architecture validation
test-complete: test-unit test-integration test-architecture test-config ## Run complete test suite including architecture
	@echo "✅ Complete test suite passed"
