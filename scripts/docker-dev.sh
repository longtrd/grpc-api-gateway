#!/bin/bash

# Docker development helper script for gRPC API Gateway

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker first."
        exit 1
    fi
}

# Build the gRPC server image
build() {
    print_status "Building gRPC server Docker image..."
    docker-compose build grpc-server
    print_status "Build completed successfully!"
}

# Start the services
start() {
    print_status "Starting gRPC server..."
    docker-compose up -d grpc-server
    print_status "gRPC server started on port 50051"

    # Wait for health check
    print_status "Waiting for server to be healthy..."
    timeout 60 bash -c 'until docker-compose ps grpc-server | grep -q healthy; do sleep 2; done' || {
        print_warning "Health check timeout. Server might still be starting..."
    }

    show_status
}

# Stop the services
stop() {
    print_status "Stopping gRPC server..."
    docker-compose down
    print_status "Services stopped"
}

# Restart the services
restart() {
    print_status "Restarting gRPC server..."
    stop
    start
}

# Show logs
logs() {
    print_status "Showing gRPC server logs..."
    docker-compose logs -f grpc-server
}

# Show service status
status() {
    print_status "Service status:"
    docker-compose ps
}

show_status() {
    echo ""
    print_status "=== gRPC Server Status ==="
    docker-compose ps grpc-server
    echo ""
    print_status "Server is accessible at: localhost:50051"
    echo ""
}

# Clean up Docker resources
clean() {
    print_status "Cleaning up Docker resources..."
    docker-compose down -v --remove-orphans
    docker system prune -f
    print_status "Cleanup completed"
}

# Rebuild and restart
rebuild() {
    print_status "Rebuilding and restarting services..."
    stop
    build
    start
}

# Show help
help() {
    echo "Docker development helper for gRPC API Gateway"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  build     Build the gRPC server Docker image"
    echo "  start     Start the gRPC server"
    echo "  stop      Stop all services"
    echo "  restart   Restart the gRPC server"
    echo "  logs      Show server logs (follow mode)"
    echo "  status    Show service status"
    echo "  clean     Clean up Docker resources"
    echo "  rebuild   Rebuild and restart services"
    echo "  help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 build && $0 start    # Build and start"
    echo "  $0 logs                 # Follow logs"
    echo "  $0 rebuild              # Full rebuild and restart"
}

# Main script logic
check_docker

case "${1:-help}" in
    build)
        build
        ;;
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    logs)
        logs
        ;;
    status)
        status
        ;;
    clean)
        clean
        ;;
    rebuild)
        rebuild
        ;;
    help|--help|-h)
        help
        ;;
    *)
        print_error "Unknown command: $1"
        echo ""
        help
        exit 1
        ;;
esac
