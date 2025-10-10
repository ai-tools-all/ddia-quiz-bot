#!/bin/bash

# Build script for quiz-tui

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building quiz-tui...${NC}"

# Get the project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Build directory
BUILD_DIR="${PROJECT_ROOT}/build"
mkdir -p "${BUILD_DIR}"

# Build the binary
cd "${PROJECT_ROOT}"

echo -e "${YELLOW}Running go build...${NC}"
go build -o "${BUILD_DIR}/quiz-tui" ./cmd/quiz-tui

if [ $? -eq 0 ]; then
    echo -e "${GREEN}Build successful!${NC}"
    echo -e "Binary location: ${BUILD_DIR}/quiz-tui"
    echo ""
    echo -e "Run with: ${BUILD_DIR}/quiz-tui --user <username>"
else
    echo -e "${RED}Build failed!${NC}"
    exit 1
fi
