#!/bin/bash

# Development Environment Setup Script for gRPC API Gateway
# This script sets up everything needed for development

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."

    local missing_tools=()

    if ! command_exists go; then
        missing_tools+=("go")
    else
        local go_version=$(go version | grep -o 'go[0-9]\+\.[0-9]\+' | head -1)
        log_success "Go detected: $go_version"
    fi

    if ! command_exists git; then
        missing_tools+=("git")
    else
        log_success "Git detected: $(git --version)"
    fi

    if ! command_exists make; then
        missing_tools+=("make")
    else
        log_success "Make detected: $(make --version | head -1)"
    fi

    if ! command_exists protoc; then
        missing_tools+=("protoc")
        log_warning "Protocol Buffers compiler (protoc) not found"
        log_info "Install protoc from: https://grpc.io/docs/protoc-installation/"
    else
        log_success "Protoc detected: $(protoc --version)"
    fi

    if [ ${#missing_tools[@]} -ne 0 ]; then
        log_error "Missing required tools: ${missing_tools[*]}"
        log_info "Please install the missing tools and run this script again"
        exit 1
    fi
}

# Setup development tools
setup_tools() {
    log_info "Installing development tools..."

    if make install-tools; then
        log_success "Development tools installed successfully"
    else
        log_error "Failed to install development tools"
        exit 1
    fi
}

# Setup project dependencies
setup_dependencies() {
    log_info "Setting up project dependencies..."

    if make setup; then
        log_success "Project dependencies set up successfully"
    else
        log_error "Failed to set up project dependencies"
        exit 1
    fi
}

# Setup git hooks
setup_git_hooks() {
    log_info "Setting up Git hooks..."

    if [ -f "scripts/setup-git-hooks.sh" ]; then
        if bash scripts/setup-git-hooks.sh; then
            log_success "Git hooks set up successfully"
        else
            log_error "Failed to set up Git hooks"
            exit 1
        fi
    else
        log_warning "Git hooks setup script not found"
    fi
}

# Generate protobuf code
generate_proto() {
    log_info "Generating Protocol Buffer code..."

    if make proto; then
        log_success "Protocol Buffer code generated successfully"
    else
        log_error "Failed to generate Protocol Buffer code"
        exit 1
    fi
}

# Run initial tests
run_tests() {
    log_info "Running initial tests..."

    if make test; then
        log_success "All tests passed"
    else
        log_warning "Some tests failed, but continuing setup"
    fi
}

# Print setup summary
print_summary() {
    echo ""
    echo -e "${GREEN}🎉 Development environment setup complete!${NC}"
    echo ""
    echo -e "${YELLOW}Available Make targets:${NC}"
    echo "  make help              - Show all available targets"
    echo "  make fmt               - Format Go code"
    echo "  make lint              - Run linters"
    echo "  make test              - Run tests"
    echo "  make build             - Build all binaries"
    echo "  make run-grpc          - Run gRPC server"
    echo "  make run-gateway       - Run API Gateway"
    echo "  make pre-commit        - Run all pre-commit checks"
    echo ""
    echo -e "${YELLOW}VS Code users:${NC}"
    echo "  - Install recommended extensions from .vscode/extensions.json"
    echo "  - Use Ctrl+Shift+P -> 'Tasks: Run Task' to access build tasks"
    echo ""
    echo -e "${YELLOW}Git hooks installed:${NC}"
    echo "  - pre-commit: Runs formatting, linting, and tests"
    echo "  - pre-push: Runs comprehensive checks"
    echo "  - commit-msg: Validates conventional commit format"
    echo ""
    echo -e "${YELLOW}Getting started:${NC}"
    echo "  1. Run: make run-grpc (in one terminal)"
    echo "  2. Run: make run-gateway (in another terminal)"
    echo "  3. Test: curl http://localhost:8080/health"
    echo ""
    echo -e "${GREEN}Happy coding! 🚀${NC}"
}

# Main setup function
main() {
    echo -e "${GREEN}gRPC API Gateway - Development Environment Setup${NC}"
    echo "=================================================="
    echo ""

    # Check if we're in the right directory
    if [ ! -f "Makefile" ]; then
        log_error "Please run this script from the project root directory"
        exit 1
    fi

    check_prerequisites
    setup_tools
    setup_dependencies
    setup_git_hooks
    generate_proto
    run_tests
    print_summary
}

# Handle interruption
trap 'log_error "Setup interrupted"; exit 1' INT TERM

# Run main function
main "$@"
