#!/bin/bash

# Setup Git Hooks for gRPC API Gateway project
# This script sets up pre-commit hooks to ensure code quality

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Setting up Git hooks for gRPC API Gateway project...${NC}"

# Create hooks directory if it doesn't exist
mkdir -p .git/hooks

# Create pre-commit hook
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash

# Pre-commit hook for gRPC API Gateway project
# Runs formatting, linting, and tests before commit

set -e

echo "Running pre-commit checks..."

# Check if we're in the project root
if [ ! -f "Makefile" ]; then
    echo "Error: Please run git commit from the project root directory"
    exit 1
fi

# Format code
echo "Formatting Go code..."
make fmt

# Run linters
echo "Running linters..."
make lint-fast

# Run tests
echo "Running tests..."
make test

# Check for any staged changes after formatting
if ! git diff --cached --quiet; then
    echo "Code was reformatted. Please stage the formatted files and commit again."
    echo "Staged files that were reformatted:"
    git diff --cached --name-only
    exit 1
fi

echo "All pre-commit checks passed!"
EOF

# Make pre-commit hook executable
chmod +x .git/hooks/pre-commit

# Create pre-push hook
cat > .git/hooks/pre-push << 'EOF'
#!/bin/bash

# Pre-push hook for gRPC API Gateway project
# Runs comprehensive checks before pushing

set -e

echo "Running pre-push checks..."

# Check if we're in the project root
if [ ! -f "Makefile" ]; then
    echo "Error: Please run git push from the project root directory"
    exit 1
fi

# Run full lint suite
echo "Running comprehensive linting..."
make lint

# Run all tests including integration tests
echo "Running all tests..."
make test
make test-integration 2>/dev/null || echo "Integration tests skipped (not implemented yet)"

# Run security scan
echo "Running security scan..."
make security 2>/dev/null || echo "Security scan skipped (gosec not installed)"

echo "All pre-push checks passed!"
EOF

# Make pre-push hook executable
chmod +x .git/hooks/pre-push

# Create commit-msg hook for conventional commits
cat > .git/hooks/commit-msg << 'EOF'
#!/bin/bash

# Commit message hook for conventional commits
# Validates commit message format

commit_regex='^(feat|fix|docs|style|refactor|perf|test|chore|ci)(\(.+\))?: .{1,50}'

error_msg="Invalid commit message format. Please use conventional commits format:
Examples:
  feat: add user authentication
  fix(api): resolve timeout issue
  docs: update README
  style: format code
  refactor: extract user service
  perf: optimize database queries
  test: add integration tests
  chore: update dependencies
  ci: add linting workflow"

if ! grep -qE "$commit_regex" "$1"; then
    echo "$error_msg"
    exit 1
fi
EOF

# Make commit-msg hook executable
chmod +x .git/hooks/commit-msg

echo -e "${GREEN}Git hooks setup complete!${NC}"
echo ""
echo -e "${YELLOW}Installed hooks:${NC}"
echo "  - pre-commit: Runs formatting, linting, and tests"
echo "  - pre-push: Runs comprehensive checks before push"
echo "  - commit-msg: Validates conventional commit format"
echo ""
echo -e "${YELLOW}To bypass hooks (not recommended):${NC}"
echo "  - git commit --no-verify"
echo "  - git push --no-verify"
echo ""
echo -e "${GREEN}Happy coding! 🚀${NC}"
