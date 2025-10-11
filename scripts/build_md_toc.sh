#!/bin/bash

# Build script for md-toc

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building md-toc...${NC}"

# Get the project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Build directory
BUILD_DIR="${PROJECT_ROOT}/build"
mkdir -p "${BUILD_DIR}"

# Build the binary
cd "${PROJECT_ROOT}"

echo -e "${YELLOW}Running go build...${NC}"
go build -o "${BUILD_DIR}/md-toc" ./cmd/md-toc

if [ $? -eq 0 ]; then
    echo -e "${GREEN}Build successful!${NC}"
    echo -e "Binary location: ${BUILD_DIR}/md-toc"
    echo ""
    echo -e "Run with: ${BUILD_DIR}/md-toc [file/directory]"
else
    echo -e "${RED}Build failed!${NC}"
    exit 1
fi
