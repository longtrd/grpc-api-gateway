# gRPC API Gateway + Load Balancer

A high-performance API Gateway that provides REST endpoints for gRPC services with built-in load balancing, built with Go and following Clean Architecture principles.

## 🚀 Overview

This project implements a production-ready API Gateway that:
- **Translates REST to gRPC**: Exposes gRPC services as REST APIs
- **Load Balancing**: Distributes requests across multiple gRPC server instances
- **High Performance**: Built with Go for optimal performance and concurrency
- **Clean Architecture**: Follows SOLID principles and clean architecture patterns
- **Dependency Injection**: Complete DI container with interface-based design
- **Structured Logging**: JSON/text logging with configurable levels
- **gRPC Interceptors**: Comprehensive middleware chain (recovery, logging, validation, metrics)
- **Environment Configuration**: Feature flags and environment-based config management
- **Docker Ready**: Multi-stage builds with production optimizations
- **Comprehensive Testing**: 91%+ test coverage with unit, integration, and benchmark tests
- **Developer Experience**: Complete linting, formatting, and development tools
- **Production Ready**: Includes monitoring, health checks, and graceful shutdown

## 🏗️ Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   REST Client   │───▶│   API Gateway    │───▶│  gRPC Server 1  │
│                 │    │  + Load Balancer │    │                 │
└─────────────────┘    │                  │───▶│  gRPC Server 2  │
                       │                  │    │                 │
                       └──────────────────┘───▶│  gRPC Server N  │
                                               └─────────────────┘
```

### Components

- **API Gateway**: HTTP/REST interface that translates requests to gRPC
- **Load Balancer**: Distributes gRPC calls across multiple server instances
- **gRPC Servers**: Backend services implementing business logic
- **Service Discovery**: Automatic detection and health monitoring of gRPC services

### Clean Architecture Implementation

The project follows Clean Architecture principles with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────────┐
│                          Infrastructure                         │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                   Interface Adapters                     │  │
│  │  ┌─────────────────────────────────────────────────────┐  │  │
│  │  │                Application Layer                    │  │  │
│  │  │  ┌───────────────────────────────────────────────┐  │  │  │
│  │  │  │             Domain Layer (Core)              │  │  │  │
│  │  │  │  • Entities (User)                           │  │  │  │
│  │  │  │  • Business Rules                            │  │  │  │
│  │  │  │  • Domain Interfaces                         │  │  │  │
│  │  │  └───────────────────────────────────────────────┘  │  │  │
│  │  │  • Use Cases (UserUseCase)                         │  │  │
│  │  │  • Application Services                             │  │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  │  • gRPC Handlers                                         │  │
│  │  • REST Controllers                                      │  │
│  │  • Data Mappers                                          │  │
│  └───────────────────────────────────────────────────────────┘  │
│  • Database Implementations                                     │
│  • External APIs                                                │
│  • Framework Bindings                                           │
└─────────────────────────────────────────────────────────────────┘
```

**Dependency Rule**: Dependencies point inward only. Inner layers never depend on outer layers.

## 🆕 Recent Updates

### ✅ Enhanced Clean Architecture Implementation
- **Dependency Injection Container**: Complete DI system in [`internal/config/container.go`](gRPC-server/internal/config/container.go)
- **Environment Configuration**: Comprehensive config management in [`internal/config/config.go`](gRPC-server/internal/config/config.go)
- **Structured Logging**: JSON/text logging system in [`pkg/logger/logger.go`](gRPC-server/pkg/logger/logger.go)
- **Domain Validation**: Business rule validation in [`pkg/validator/user_validator.go`](gRPC-server/pkg/validator/user_validator.go)

### ✅ gRPC Interceptor Chain
- **Recovery Interceptor**: Panic recovery with stack traces
- **Logging Interceptor**: Structured request/response logging
- **Validation Interceptor**: Context and request validation
- **Metrics Interceptor**: Performance monitoring (placeholder)
- **Auth Interceptor**: Authentication middleware (extensible)
- **Timeout Interceptor**: Request timeout handling

### ✅ Production-Ready Docker Setup
- **Multi-stage Builds**: Optimized Docker images with protobuf compilation
- **Health Checks**: Built-in TCP and gRPC health monitoring
- **Environment Configuration**: 25+ configurable options via environment variables
- **Production Overrides**: Resource limits, security, and monitoring
- **Development Tools**: Docker helper scripts and automation

### ✅ Enhanced Testing Framework
- **Architecture Compliance**: Tests for clean architecture patterns
- **Configuration Testing**: Environment-based configuration validation
- **Mock Infrastructure**: Thread-safe mocks with call tracking
- **Integration Tests**: End-to-end workflow validation
- **Performance Benchmarks**: Critical operation benchmarking

### ✅ Developer Experience Improvements
- **35+ Makefile Targets**: Comprehensive development automation
- **Docker Commands**: Build, start, test, and manage containers
- **Environment Template**: Complete configuration example in [`env.example`](env.example)
- **Documentation**: Comprehensive [`DOCKER.md`](DOCKER.md) with best practices

## 📁 Project Structure

```
grpc-api-gateway/
├── api-gateway/              # API Gateway service
│   ├── cmd/
│   │   └── main.go          # Gateway entry point
│   ├── internal/
│   │   ├── config/          # Configuration management
│   │   ├── gateway/         # Gateway implementation
│   │   ├── loadbalancer/    # Load balancing logic
│   │   ├── middleware/      # HTTP middleware
│   │   └── handlers/        # REST endpoint handlers
│   └── go.mod
├── gRPC-server/             # gRPC backend services
│   ├── cmd/
│   │   ├── main.go         # CLI entry point with clean architecture
│   │   └── commands/       # Cobra commands
│   │       ├── root.go     # Root command setup
│   │       └── start.go    # Enhanced server startup with DI container
│   ├── internal/
│   │   ├── domain/         # Business entities (Clean Architecture core)
│   │   │   ├── user.go     # Domain entities and business rules
│   │   │   ├── interfaces.go # Repository and service contracts
│   │   │   └── user_test.go  # Domain tests (91.7% coverage)
│   │   ├── usecase/        # Business logic / Application layer
│   │   │   ├── user.go     # Use case implementations
│   │   │   ├── user_test.go # Use case tests (91.9% coverage)
│   │   │   └── user_bench_test.go # Performance benchmarks
│   │   ├── handler/        # gRPC handlers (Interface adapters)
│   │   │   └── middleware.go # gRPC interceptors (logging, recovery, validation)
│   │   ├── repository/     # Data access implementations
│   │   ├── mocks/          # Test mocks with call tracking
│   │   │   ├── user_repository_mock.go   # Thread-safe repository mock
│   │   │   ├── user_validator_mock.go    # Validator mock
│   │   │   ├── logger_mock.go            # Logger mock
│   │   │   └── event_publisher_mock.go   # Event publisher mock
│   │   ├── testutil/       # Test utilities and helpers
│   │   │   └── helpers.go  # Test factories and assertions
│   │   └── config/         # Configuration & Dependency Injection
│   │       ├── config.go   # Environment-based configuration management
│   │       └── container.go # Dependency injection container
│   ├── pkg/                # Public libraries (Clean Architecture)
│   │   ├── logger/         # Structured logging implementation
│   │   │   └── logger.go   # JSON/text logging with levels
│   │   └── validator/      # Domain validation utilities
│   │       └── user_validator.go # Business rule validation
│   ├── proto/
│   │   └── user.proto      # Service definitions
│   ├── tests/
│   │   └── integration/    # Integration tests
│   │       └── user_integration_test.go # End-to-end workflow tests
│   ├── Dockerfile          # Multi-stage Docker build
│   ├── .dockerignore       # Docker build optimization
│   └── go.mod
├── scripts/                # Development automation
│   ├── setup-dev.sh       # Development environment setup
│   ├── setup-git-hooks.sh # Git hooks for code quality
│   └── docker-dev.sh      # Docker development helper
├── docker-compose.yml     # Development Docker setup
├── docker-compose.prod.yml # Production Docker configuration
├── env.example            # Environment configuration template
├── DOCKER.md              # Comprehensive Docker documentation
├── Makefile               # Enhanced development automation (35+ targets)
└── README.md
```

## 🚦 Getting Started

### Prerequisites

- Go 1.21 or higher
- Protocol Buffers compiler (`protoc`)
- Make (optional, for automation)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/longtrd/grpc-api-gateway.git
   cd grpc-api-gateway
   ```

2. **Install dependencies**
   ```bash
   # Install gRPC server dependencies
   cd gRPC-server
   go mod tidy
   go get github.com/spf13/cobra
   go get google.golang.org/grpc
   go get google.golang.org/protobuf

   # Install API gateway dependencies
   cd ../api-gateway
   go mod tidy
   go get github.com/gin-gonic/gin
   go get google.golang.org/grpc
   ```

3. **Generate Protocol Buffer code**
   ```bash
   cd gRPC-server
   protoc --go_out=. --go-grpc_out=. proto/*.proto
   ```

### Quick Start

#### Option 1: Docker Development (Recommended)

```bash
# 1. Copy environment template
cp env.example .env

# 2. Build and start with Docker
make docker-build
make docker-start

# 3. Check status
make docker-status
make docker-logs

# 4. Test the gRPC server
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check

# 5. Stop services
make docker-stop
```

#### Option 2: Manual Development

1. **Start gRPC servers (multiple instances for load balancing)**
   ```bash
   # Terminal 1 - First gRPC server instance
   cd gRPC-server
   go run cmd/main.go start --port 50051

   # Terminal 2 - Second gRPC server instance
   cd gRPC-server
   go run cmd/main.go start --port 50052

   # Terminal 3 - Third gRPC server instance
   cd gRPC-server
   go run cmd/main.go start --port 50053
   ```

2. **Start API Gateway**
   ```bash
   # Terminal 4 - API Gateway
   cd api-gateway
   go run cmd/main.go
   ```

3. **Test the API**
   ```bash
   # Create a user via REST API
   curl -X POST http://localhost:8080/api/v1/users \
     -H "Content-Type: application/json" \
     -d '{"name": "John Doe", "email": "john@example.com"}'

   # Get user by ID
   curl http://localhost:8080/api/v1/users/user-123
   ```

## 📚 API Documentation

### User Service Endpoints

| Method | Endpoint | Description | gRPC Method |
|--------|----------|-------------|-------------|
| POST | `/api/v1/users` | Create a new user | `CreateUser` |
| GET | `/api/v1/users/{id}` | Get user by ID | `GetUser` |
| PUT | `/api/v1/users/{id}` | Update user | `UpdateUser` |
| DELETE | `/api/v1/users/{id}` | Delete user | `DeleteUser` |

### Request/Response Examples

#### Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john.doe@example.com"
  }'
```

**Response:**
```json
{
  "id": "user-123",
  "name": "John Doe",
  "email": "john.doe@example.com"
}
```

#### Get User
```bash
curl http://localhost:8080/api/v1/users/user-123
```

**Response:**
```json
{
  "id": "user-123",
  "name": "John Doe",
  "email": "john.doe@example.com"
}
```

## ⚙️ Configuration

### API Gateway Configuration

```yaml
# config/gateway.yaml
server:
  port: 8080
  host: "0.0.0.0"

grpc:
  servers:
    - "localhost:50051"
    - "localhost:50052"
    - "localhost:50053"
  timeout: "30s"

loadbalancer:
  strategy: "round_robin"  # round_robin, least_connections, random
  health_check_interval: "10s"

logging:
  level: "info"
  format: "json"
```

### gRPC Server Configuration

The gRPC server uses environment-based configuration with comprehensive feature flags:

```bash
# Copy and customize the environment template
cp env.example .env

# Key configuration options
GRPC_PORT=50051                    # Server port
GRPC_HOST=0.0.0.0                 # Server host
GRPC_READ_TIMEOUT=30s             # Request timeout
LOG_LEVEL=info                    # Logging level (debug, info, warn, error)
LOG_FORMAT=json                   # Log format (json, text)
DATABASE_TYPE=memory              # Database type (memory, postgres)
ENABLE_HEALTH_CHECK=true          # Health check service
ENABLE_METRICS=false              # Prometheus metrics
ENABLE_AUTH=false                 # Authentication
ENABLE_TRACING=false              # Distributed tracing
```

See [`env.example`](env.example) for complete configuration options.

### Dependency Injection Configuration

The server uses a comprehensive DI container that automatically wires dependencies:

```go
// Load configuration from environment
cfg := config.Load()

// Build dependency container with all services
container, err := config.NewContainer(cfg)
if err != nil {
    log.Fatal("Failed to build container:", err)
}

// Access services through container
logger := container.GetLogger()
userUseCase := container.GetUserUseCase()
```

### Feature Flags

Control functionality through environment variables:

| Feature | Environment Variable | Description |
|---------|---------------------|-------------|
| Health Checks | `ENABLE_HEALTH_CHECK=true` | gRPC health checking protocol |
| Metrics | `ENABLE_METRICS=false` | Prometheus metrics collection |
| Authentication | `ENABLE_AUTH=false` | JWT/API key authentication |
| Tracing | `ENABLE_TRACING=false` | Distributed tracing |
| Reflection | Auto-enabled in debug mode | gRPC service reflection |

## 🏗️ Development

### Development Setup

The project includes comprehensive development tools for code quality and testing:

```bash
# Complete development environment setup
./scripts/setup-dev.sh

# Install development tools only
make install-tools

# Setup Git hooks for code quality
./scripts/setup-git-hooks.sh
```

### Code Quality Tools

#### Linting & Formatting
```bash
# Format all Go code
make fmt

# Run all linters
make lint

# Fast linting for CI
make lint-fast

# Pre-commit checks (format + lint + test)
make pre-commit
```

#### Available Development Commands
```bash
# Development workflow
make setup           # Complete project setup
make build           # Build all binaries
make clean           # Clean generated files
make deps-check      # Check for outdated dependencies
make deps-update     # Update all dependencies
make security        # Run security analysis

# Docker development (NEW)
make docker-build    # Build Docker images with clean architecture
make docker-start    # Start services with Docker Compose
make docker-stop     # Stop Docker services
make docker-restart  # Restart Docker services
make docker-logs     # Show Docker service logs
make docker-status   # Show Docker service status
make docker-clean    # Clean Docker resources
make docker-rebuild  # Rebuild and restart Docker services
make docker-test     # Build and test with Docker

# Architecture and testing (NEW)
make test-architecture # Test clean architecture compliance
make test-config      # Test configuration management
make test-complete    # Complete test suite including architecture

# Continuous development
make watch-grpc      # Auto-restart gRPC server on changes
make watch-gateway   # Auto-restart API Gateway on changes
```

#### Git Hooks
- **Pre-commit**: Runs formatting, linting, and tests before each commit
- **Pre-push**: Comprehensive checks before pushing to remote
- **Commit-msg**: Enforces conventional commit message format

#### VS Code Integration
- Automatic code formatting on save
- Real-time linting and error detection
- Integrated test runner with coverage visualization
- Recommended extensions for Go development

### Adding New Services

1. **Define Protocol Buffer**
   ```protobuf
   // proto/order.proto
   service OrderService {
     rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);
     rpc GetOrder(GetOrderRequest) returns (GetOrderResponse);
   }
   ```

2. **Implement gRPC Handler**
   ```go
   // internal/handler/order_handler.go
   type OrderHandler struct {
     usecase domain.OrderUseCase
   }

   func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
     // Implementation
   }
   ```

3. **Add REST Endpoints**
   ```go
   // api-gateway/internal/handlers/order.go
   func (h *OrderHandler) CreateOrder(c *gin.Context) {
     // REST to gRPC translation
   }
   ```

### Testing

The project includes a comprehensive testing framework following clean architecture principles with **91%+ test coverage**.

#### Test Structure
```
gRPC-server/
├── internal/
│   ├── domain/
│   │   └── user_test.go        # Domain entity tests
│   ├── usecase/
│   │   ├── user_test.go        # Business logic tests
│   │   └── user_bench_test.go  # Performance benchmarks
│   ├── mocks/                  # Mock implementations
│   │   ├── user_repository_mock.go
│   │   ├── user_validator_mock.go
│   │   ├── logger_mock.go
│   │   └── event_publisher_mock.go
│   └── testutil/
│       └── helpers.go          # Test utilities and factories
└── tests/
    └── integration/
        └── user_integration_test.go  # End-to-end tests
```

#### Running Tests

```bash
# Run all unit tests
make test

# Run specific test suites
make test-domain     # Domain layer tests (91.7% coverage)
make test-usecase    # Business logic tests (91.9% coverage)
make test-unit       # All unit tests
make test-integration # Integration tests (build tag required)

# Architecture and compliance testing (NEW)
make test-architecture # Test clean architecture compliance
make test-config      # Test configuration scenarios
make test-complete    # Complete test suite with architecture validation

# Test analysis
make test-coverage   # Generate HTML coverage report
make test-verbose    # Verbose output
make test-bench      # Performance benchmarks

# Direct Go commands
cd gRPC-server
go test ./internal/...                    # All unit tests
go test -v ./internal/domain/...          # Domain tests with verbose output
go test -race -cover ./internal/...       # With race detection and coverage
go test -tags=integration ./tests/...     # Integration tests
go test -bench=. ./internal/usecase/...   # Benchmarks
```

#### Test Coverage Results
- **Domain Layer**: 91.7% coverage
- **UseCase Layer**: 91.9% coverage
- **Overall**: 91%+ coverage with race detection

#### Test Features
- ✅ **Table-Driven Tests**: Comprehensive test cases for all scenarios
- ✅ **Mock Isolation**: Thread-safe mock implementations with call tracking
- ✅ **Race Detection**: All tests run with `-race` flag
- ✅ **Integration Tests**: Full workflow testing with concurrent operations
- ✅ **Benchmarks**: Performance testing for critical operations
- ✅ **Test Utilities**: Reusable test helpers and data factories
- ✅ **Error Validation**: Proper error type checking with `errors.Is()`
- ✅ **Concurrency Testing**: Multi-goroutine scenarios
- ✅ **Architecture Compliance**: Tests for clean architecture patterns
- ✅ **Configuration Testing**: Environment-based configuration validation
- ✅ **Docker Testing**: Containerized testing with health checks

#### Test Development Workflow
```bash
# 1. Write failing test
go test -v ./internal/usecase/...

# 2. Implement feature
# ... code changes ...

# 3. Verify test passes
go test -v ./internal/usecase/...

# 4. Check coverage
make test-coverage

# 5. Run full suite
make test
```

#### Load Testing
```bash
# API Gateway performance testing
hey -n 1000 -c 10 http://localhost:8080/api/v1/users
```

## 🔧 Load Balancing Strategies

### Round Robin (Default)
Distributes requests evenly across all healthy servers.

### Least Connections
Routes requests to the server with the fewest active connections.

### Random
Randomly selects a healthy server for each request.

### Weighted Round Robin
Allows assigning different weights to servers based on capacity.

## 📊 Monitoring & Observability

### Health Checks
- **API Gateway**: `GET /health`
- **gRPC Servers**: gRPC Health Checking Protocol

### Metrics (Prometheus)
- Request latency histograms
- Request count by status code
- Active connections per backend
- Error rates

### Logging
- Structured JSON logging
- Request/response logging with correlation IDs
- Error tracking and alerting

### Tracing
- Distributed tracing with OpenTelemetry
- Request flow across gateway and backends

## 🐳 Docker Development

The project includes a comprehensive Docker setup with clean architecture integration. See [`DOCKER.md`](DOCKER.md) for complete documentation.

### Quick Docker Start

```bash
# Build and start gRPC server
make docker-build
make docker-start

# Check status and logs
make docker-status
make docker-logs

# Test the setup
make docker-test

# Stop services
make docker-stop
```

### Docker Features

- **Multi-stage builds** with protobuf compilation
- **Clean architecture** integration with DI container
- **Environment configuration** with 25+ options
- **Health checks** for service monitoring
- **Production optimizations** with resource limits
- **Development tools** with helper scripts

### Docker Commands

| Command | Description |
|---------|-------------|
| `./scripts/docker-dev.sh build` | Build Docker images |
| `./scripts/docker-dev.sh start` | Start services |
| `./scripts/docker-dev.sh logs` | Follow logs |
| `./scripts/docker-dev.sh status` | Show service status |
| `./scripts/docker-dev.sh clean` | Clean resources |

### Environment Configuration

```bash
# Copy template and customize
cp env.example .env

# Key Docker environment variables
GRPC_PORT=50051
LOG_LEVEL=info
LOG_FORMAT=json
ENABLE_HEALTH_CHECK=true
DATABASE_TYPE=memory
```

### Production Docker

```bash
# Production deployment with optimizations
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Features: resource limits, security, monitoring
```

## 🚀 Deployment

### Docker

The project uses multi-stage Docker builds with clean architecture integration:

```dockerfile
# gRPC-server/Dockerfile
FROM golang:1.24.2-alpine AS builder

# Install build dependencies and protobuf compiler
RUN apk add --no-cache git protobuf protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Copy source and generate protobuf
COPY . .
RUN protoc --go_out=. --go-grpc_out=. proto/*.proto

# Build with clean architecture
RUN CGO_ENABLED=0 GOOS=linux go build -o grpc-server ./cmd/main.go

# Runtime stage with security hardening
FROM alpine:latest
RUN apk --no-cache add ca-certificates netcat-openbsd
RUN addgroup -g 1000 -S appgroup && adduser -u 1000 -S appuser -G appgroup

WORKDIR /root/
COPY --from=builder /app/grpc-server .
RUN chown appuser:appgroup grpc-server

USER appuser
EXPOSE 50051
HEALTHCHECK --interval=30s --timeout=3s CMD nc -z localhost 50051 || exit 1
CMD ["./grpc-server", "start"]
```

### Docker Compose with Clean Architecture

```yaml
version: '3.8'
services:
  grpc-server:
    build:
      context: ./gRPC-server
      dockerfile: Dockerfile
    ports:
      - "50051:50051"
    environment:
      - GRPC_PORT=50051
      - LOG_LEVEL=info
      - LOG_FORMAT=json
      - DATABASE_TYPE=memory
      - ENABLE_HEALTH_CHECK=true
    healthcheck:
      test: ["CMD", "nc", "-z", "localhost", "50051"]
      interval: 30s
      timeout: 10s
      retries: 3
    networks:
      - grpc-network

networks:
  grpc-network:
    driver: bridge
```

### Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-gateway
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api-gateway
  template:
    metadata:
      labels:
        app: api-gateway
    spec:
      containers:
      - name: gateway
        image: your-registry/api-gateway:latest
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: api-gateway-service
spec:
  selector:
    app: api-gateway
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

## 🔒 Security

### Authentication & Authorization
- JWT token validation
- API key authentication
- Rate limiting per client
- CORS configuration

### TLS/SSL
- TLS termination at gateway
- mTLS for backend communication
- Certificate management

## 🛠️ Performance Tuning

### Connection Pooling
```go
// gRPC connection pool configuration
grpc.WithDefaultCallOptions(
    grpc.MaxCallRecvMsgSize(4*1024*1024),
    grpc.MaxCallSendMsgSize(4*1024*1024),
)
```

### Caching
- Response caching for read operations
- Redis integration for distributed caching
- Cache invalidation strategies

## 📖 Additional Documentation

For comprehensive information on specific topics:

- **[DOCKER.md](DOCKER.md)**: Complete Docker setup, development workflow, production deployment, troubleshooting, and best practices
- **[env.example](env.example)**: Environment configuration template with all available options and descriptions
- **[Makefile](Makefile)**: 35+ development targets for building, testing, Docker operations, and quality assurance

### Architecture Documentation

The project follows clean architecture with extensive documentation embedded in the cursor rules:

- **Domain Layer**: [`gRPC-server/internal/domain/`](gRPC-server/internal/domain/) - Business entities and rules
- **Use Case Layer**: [`gRPC-server/internal/usecase/`](gRPC-server/internal/usecase/) - Application services
- **Interface Layer**: [`gRPC-server/internal/handler/`](gRPC-server/internal/handler/) - gRPC handlers and middleware
- **Infrastructure**: [`gRPC-server/pkg/`](gRPC-server/pkg/) - External concerns (logging, validation)
- **Configuration**: [`gRPC-server/internal/config/`](gRPC-server/internal/config/) - DI container and config management

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [gRPC](https://grpc.io/) for the excellent RPC framework
- [Gin](https://github.com/gin-gonic/gin) for the HTTP web framework
- [Cobra](https://github.com/spf13/cobra) for the CLI interface
- Clean Architecture principles by Robert C. Martin
