# Docker Setup for gRPC API Gateway

This document provides comprehensive instructions for running the gRPC server with Docker.

## Quick Start

### 1. Build and Run with Docker Compose

```bash
# Build and start the gRPC server
make docker-build
make docker-start

# Or use the helper script directly
./scripts/docker-dev.sh build
./scripts/docker-dev.sh start
```

### 2. Check Server Status

```bash
# Check service status
make docker-status

# Follow logs
make docker-logs

# Test connectivity
nc -z localhost 50051
```

### 3. Stop Services

```bash
make docker-stop
# or
./scripts/docker-dev.sh stop
```

## Docker Files Overview

### Core Files

- **`gRPC-server/Dockerfile`** - Multi-stage build for gRPC server
- **`docker-compose.yml`** - Development setup with health checks
- **`docker-compose.prod.yml`** - Production overrides with resource limits
- **`gRPC-server/.dockerignore`** - Optimizes build context
- **`scripts/docker-dev.sh`** - Development helper script

### Build Process

The Docker build process includes:

1. **Protocol Buffer Compilation** - Generates Go code from `.proto` files
2. **Dependency Download** - Downloads Go modules
3. **Clean Architecture Setup** - Builds with dependency injection container
4. **Binary Build** - Creates optimized static binary with all layers
5. **Multi-stage** - Minimal runtime image with only the binary

### Architecture Integration

The Docker setup integrates with the clean architecture:

- **Configuration Management** - Environment-based configuration loading
- **Dependency Injection** - Automated container setup with all services
- **Structured Logging** - JSON/text logging with proper levels
- **gRPC Interceptors** - Logging, recovery, validation, metrics
- **Health Checks** - Built-in health monitoring
- **Graceful Shutdown** - Proper resource cleanup

## Available Commands

### Make Targets

```bash
make docker-build      # Build Docker images
make docker-start      # Start services
make docker-stop       # Stop services
make docker-restart    # Restart services
make docker-logs       # Show service logs
make docker-status     # Show service status
make docker-clean      # Clean Docker resources
make docker-rebuild    # Rebuild and restart
make docker-test       # Build and test connectivity
```

### Helper Script Commands

```bash
./scripts/docker-dev.sh build      # Build images
./scripts/docker-dev.sh start      # Start services
./scripts/docker-dev.sh stop       # Stop services
./scripts/docker-dev.sh restart    # Restart services
./scripts/docker-dev.sh logs       # Follow logs
./scripts/docker-dev.sh status     # Show status
./scripts/docker-dev.sh clean      # Clean resources
./scripts/docker-dev.sh rebuild    # Rebuild and restart
./scripts/docker-dev.sh help       # Show help
```

## Configuration

### Environment Variables

Copy the example environment file and customize as needed:

```bash
cp env.example .env
```

Key configuration options:

**Server Settings:**
- `GRPC_PORT=50051` - Server port
- `GRPC_HOST=0.0.0.0` - Server host
- `GRPC_READ_TIMEOUT=30s` - Request timeout
- `LOG_LEVEL=info` - Logging level (debug, info, warn, error)
- `LOG_FORMAT=json` - Log format (json, text)

**Feature Flags:**
- `ENABLE_HEALTH_CHECK=true` - Health check service
- `ENABLE_METRICS=false` - Prometheus metrics
- `ENABLE_AUTH=false` - Authentication
- `ENABLE_TRACING=false` - Distributed tracing

**Database:**
- `DATABASE_TYPE=memory` - Database type (memory, postgres)

See `env.example` for complete configuration options.

### Port Configuration

The gRPC server runs on port `50051` by default. To change:

```bash
# In docker-compose.yml
ports:
  - "8080:50051"  # Map to different host port

# Or set environment variable
GRPC_PORT=8080
```

## Development Workflow

### 1. Development Mode

```bash
# Start development environment
make docker-start

# Make code changes...

# Rebuild and restart
make docker-rebuild

# Follow logs during development
make docker-logs
```

### 2. Testing Changes

```bash
# Test with Docker
make docker-test

# Or manual testing
docker-compose up -d grpc-server
grpcurl -plaintext localhost:50051 list
docker-compose down
```

### 3. Debugging

```bash
# Access container shell
docker-compose exec grpc-server sh

# Check health
docker-compose exec grpc-server nc -z localhost 50051

# View detailed logs
docker-compose logs -f grpc-server
```

## Production Deployment

### 1. Using Production Override

```bash
# Build for production
docker-compose -f docker-compose.yml -f docker-compose.prod.yml build

# Start production services
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Check status
docker-compose -f docker-compose.yml -f docker-compose.prod.yml ps
```

### 2. Production Features

- **Resource Limits** - CPU and memory constraints
- **Health Checks** - Longer intervals for stability
- **Logging** - Structured JSON logs with rotation
- **Security** - Read-only filesystem, non-root user
- **Restart Policies** - Automatic restart on failure

### 3. Monitoring (Optional)

Uncomment monitoring services in `docker-compose.prod.yml`:

```bash
# Enable Prometheus and Grafana
# Edit docker-compose.prod.yml to uncomment monitoring services

# Start with monitoring
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Access monitoring
# Prometheus: http://localhost:9090
# Grafana: http://localhost:3000 (admin/admin)
```

## Docker Image Details

### Base Images

- **Build Stage**: `golang:1.24.2-alpine` - Includes Go toolchain and protoc
- **Runtime Stage**: `alpine:latest` - Minimal runtime environment

### Security Features

- **Non-root user** - Runs as `appuser:appgroup` (UID/GID 1000)
- **Read-only filesystem** - Prevents runtime modifications
- **Minimal attack surface** - Only includes necessary components

### Build Optimizations

- **Multi-stage build** - Reduces final image size
- **Layer caching** - Optimized for fast rebuilds
- **Static binary** - No external dependencies
- **.dockerignore** - Excludes unnecessary files

## Networking

### Service Communication

```bash
# Internal service network
grpc-network (bridge)

# Service discovery
grpc-server:50051  # gRPC server
# api-gateway:8080  # API Gateway (when added)
# postgres:5432     # Database (when added)
```

### External Access

```bash
# gRPC server
localhost:50051

# Health check endpoint
nc -z localhost 50051

# Future services
# localhost:8080  # API Gateway
# localhost:3000  # Grafana
# localhost:9090  # Prometheus
```

## Troubleshooting

### Common Issues

#### Build Failures

```bash
# Clear Docker cache
make docker-clean
docker system prune -a

# Rebuild from scratch
make docker-rebuild
```

#### Connection Issues

```bash
# Check if port is available
lsof -i :50051

# Check container health
docker-compose ps
docker-compose logs grpc-server

# Test connectivity
nc -z localhost 50051
```

#### Performance Issues

```bash
# Check resource usage
docker stats

# Check logs for errors
docker-compose logs grpc-server | grep ERROR

# Monitor health checks
docker-compose ps | grep healthy
```

### Debug Mode

```bash
# Run with debug logging
LOG_LEVEL=debug make docker-start

# Access container for debugging
docker-compose exec grpc-server sh

# Check running processes
docker-compose exec grpc-server ps aux
```

## Integration with Existing Tools

### With Makefile

All Docker commands are integrated with the existing Makefile:

```bash
make help  # Shows all available targets including Docker
```

### With Testing

```bash
# Test Docker setup
make docker-test

# Run tests inside container
docker-compose exec grpc-server go test ./internal/...
```

### With Development Scripts

The Docker setup works alongside existing development tools:

```bash
# Setup development environment
make setup

# Run tests locally
make test

# Run with Docker
make docker-start
```

## Next Steps

### Adding API Gateway

When the API Gateway is ready:

1. Uncomment API Gateway service in `docker-compose.yml`
2. Create `api-gateway/Dockerfile`
3. Update networking configuration
4. Add service discovery between components

### Adding Database

When database integration is needed:

1. Uncomment Postgres service in `docker-compose.yml`
2. Add database migrations
3. Update application configuration
4. Add data persistence volumes

### Adding Monitoring

For production monitoring:

1. Uncomment monitoring services in `docker-compose.prod.yml`
2. Add Prometheus metrics to application
3. Create Grafana dashboards
4. Set up alerting rules

## Clean Architecture Integration

This Docker setup fully implements the clean architecture patterns:

### ✅ Applied Cursor Rules

**Clean Architecture:**
- [x] Domain layer with business entities and interfaces
- [x] Use case layer with application-specific business logic
- [x] Handler layer with gRPC interceptors and adapters
- [x] Infrastructure layer with external dependencies
- [x] Dependency injection container for all services

**SOLID Principles:**
- [x] Single Responsibility - Each service has one purpose
- [x] Open/Closed - Extensible through interfaces
- [x] Liskov Substitution - All implementations are interchangeable
- [x] Interface Segregation - Small, focused interfaces
- [x] Dependency Inversion - Depend on abstractions

**gRPC Best Practices:**
- [x] Error handling with proper gRPC status codes
- [x] Context handling with cancellation support
- [x] Input validation at handler level
- [x] Logging, recovery, and metrics interceptors
- [x] Health checks and reflection
- [x] Graceful shutdown with timeout

**Testing Patterns:**
- [x] Table-driven tests for domain validation
- [x] Mock-based testing for use cases
- [x] Integration tests for complete workflows
- [x] Benchmark tests for performance
- [x] Test utilities and factories
- [x] 91%+ test coverage maintained

**Configuration Management:**
- [x] Environment-based configuration loading
- [x] Feature flags for optional functionality
- [x] Validation at startup
- [x] Multiple environment support

## Best Practices

1. **Use specific image tags** - Avoid `latest` in production
2. **Health checks** - Always implement proper health checks
3. **Resource limits** - Set appropriate CPU and memory limits
4. **Logging** - Use structured logging with proper levels
5. **Security** - Run as non-root user with minimal permissions
6. **Monitoring** - Implement metrics and observability
7. **Backups** - Regular backup of data volumes
8. **Updates** - Keep base images and dependencies updated
9. **Clean Architecture** - Follow dependency inversion principle
10. **Documentation** - Keep documentation in sync with code changes
