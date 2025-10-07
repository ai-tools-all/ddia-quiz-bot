#!/bin/bash

# DDIA Quiz Bot Build Script
# This script builds the Go application with proper error handling

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo -e "${GREEN}Building DDIA Quiz Bot...${NC}"
echo "Project root: $PROJECT_ROOT"

# Change to project root
cd "$PROJECT_ROOT"

# Clean previous build
echo -e "${YELLOW}Cleaning previous build...${NC}"
rm -f quiz-daemon

# Run go mod tidy
echo -e "${YELLOW}Running go mod tidy...${NC}"
go mod tidy

# Run tests if they exist
if [ -d "tests" ] || find . -name "*_test.go" | grep -q .; then
    echo -e "${YELLOW}Running tests...${NC}"
    go test ./...
else
    echo -e "${YELLOW}No tests found, skipping...${NC}"
fi

# Build the application
echo -e "${YELLOW}Building application...${NC}"
go build -o quiz-daemon ./cmd/quiz-daemon

# Check if build was successful
if [ $? -eq 0 ]; then
    echo -e "${GREEN}Build successful!${NC}"
    echo "Binary created: $PROJECT_ROOT/quiz-daemon"

    # Show binary info
    if command -v file >/dev/null 2>&1; then
        echo "Binary info:"
        file quiz-daemon
    fi

    # Show size
    if command -v du >/dev/null 2>&1; then
        echo "Binary size:"
        du -h quiz-daemon
    fi
else
    echo -e "${RED}Build failed!${NC}"
    exit 1
fi

echo -e "${GREEN}Build completed successfully!${NC}"