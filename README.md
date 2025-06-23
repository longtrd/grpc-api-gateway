# gRPC API Gateway + Load Balancer

A high-performance API Gateway that provides REST endpoints for gRPC services with built-in load balancing, built with Go and following Clean Architecture principles.

## 🚀 Overview

This project implements a production-ready API Gateway that:
- **Translates REST to gRPC**: Exposes gRPC services as REST APIs
- **Load Balancing**: Distributes requests across multiple gRPC server instances
- **High Performance**: Built with Go for optimal performance and concurrency
- **Clean Architecture**: Follows SOLID principles and clean architecture patterns
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
│   │   ├── main.go         # CLI entry point
│   │   └── commands/       # Cobra commands
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
│   │   ├── repository/     # Data access implementations
│   │   ├── mocks/          # Test mocks with call tracking
│   │   ├── testutil/       # Test utilities and helpers
│   │   └── config/         # Configuration & Dependency Injection
│   ├── proto/
│   │   └── user.proto      # Service definitions
│   ├── tests/
│   │   └── integration/    # Integration tests
│   └── go.mod
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

```yaml
# config/server.yaml
server:
  port: "50051"

database:
  type: "memory"  # memory, postgres, mysql

features:
  enable_reflection: true
  enable_health_check: true
  enable_metrics: true
```

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
make test-domain     # Domain layer tests
make test-usecase    # Business logic tests
make test-unit       # All unit tests
make test-integration # Integration tests (build tag required)

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
- ✅ **Mock Isolation**: Clean mock implementations with call tracking
- ✅ **Race Detection**: All tests run with `-race` flag
- ✅ **Integration Tests**: Full workflow testing
- ✅ **Benchmarks**: Performance testing for critical operations
- ✅ **Test Utilities**: Reusable test helpers and data factories
- ✅ **Error Validation**: Proper error type checking
- ✅ **Concurrency Testing**: Multi-goroutine scenarios

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

## 🚀 Deployment

### Docker

```dockerfile
# api-gateway/Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o gateway cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/gateway .
CMD ["./gateway"]
```

### Docker Compose

```yaml
version: '3.8'
services:
  api-gateway:
    build: ./api-gateway
    ports:
      - "8080:8080"
    depends_on:
      - grpc-server-1
      - grpc-server-2

  grpc-server-1:
    build: ./gRPC-server
    command: ["./server", "start", "--port", "50051"]

  grpc-server-2:
    build: ./gRPC-server
    command: ["./server", "start", "--port", "50052"]
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
