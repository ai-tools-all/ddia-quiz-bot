#!/bin/bash

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

BUILD_DIR="build"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

log() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Create build directory
mkdir -p "$BUILD_DIR"

log "Building DDIA Clicker binaries..."
log "Version: $VERSION"
log "Commit:  $COMMIT"
log "Time:    $BUILD_TIME"
echo ""

# Build quiz TUI
log "Building quiz-tui..."
go build -o "$BUILD_DIR/quiz-tui" ./cmd/quiz-tui
if [ -f "$BUILD_DIR/quiz-tui" ]; then
    size=$(ls -lh "$BUILD_DIR/quiz-tui" | awk '{print $5}')
    success "Built quiz-tui ($size)"
else
    echo "Failed to build quiz-tui"
    exit 1
fi

# Build validate-quiz
log "Building validate-quiz..."
go build -o "$BUILD_DIR/validate-quiz" ./cmd/validate-quiz
if [ -f "$BUILD_DIR/validate-quiz" ]; then
    size=$(ls -lh "$BUILD_DIR/validate-quiz" | awk '{print $5}')
    success "Built validate-quiz ($size)"
else
    echo "Failed to build validate-quiz"
    exit 1
fi

# Build md-toc
log "Building md-toc..."
go build -o "$BUILD_DIR/md-toc" ./cmd/md-toc
if [ -f "$BUILD_DIR/md-toc" ]; then
    size=$(ls -lh "$BUILD_DIR/md-toc" | awk '{print $5}')
    success "Built md-toc ($size)"
else
    echo "Failed to build md-toc"
    exit 1
fi

echo ""
log "Build artifacts:"
ls -lh "$BUILD_DIR"

echo ""
success "Build completed successfully!"
