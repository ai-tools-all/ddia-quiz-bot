#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

print_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help        Show this help message"
    echo "  -v, --verbose     Verbose output"
    echo "  -r, --race        Enable race detector"
    echo "  -c, --coverage    Generate coverage report"
    echo "  -s, --short       Run short tests only"
    echo "  -p, --package     Run tests for specific package (e.g., ./internal/markdown/...)"
    echo "  -t, --timeout     Set timeout (default: 5m)"
    echo "  -w, --watch       Watch mode (requires gotestsum)"
    echo "  --parallel N      Number of parallel tests (default: $(nproc))"
    echo ""
    echo "Examples:"
    echo "  $0                      # Run all tests"
    echo "  $0 -v -r                # Run with verbose and race detection"
    echo "  $0 -c                   # Run with coverage"
    echo "  $0 -p ./internal/tui/...  # Run TUI tests only"
}

log() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Default values
VERBOSE=""
RACE=""
COVERAGE=""
SHORT=""
PACKAGE="./..."
TIMEOUT="5m"
WATCH=false
PARALLEL=$(nproc)

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            print_usage
            exit 0
            ;;
        -v|--verbose)
            VERBOSE="-v"
            shift
            ;;
        -r|--race)
            RACE="-race"
            shift
            ;;
        -c|--coverage)
            COVERAGE="-cover -coverprofile=coverage.out"
            shift
            ;;
        -s|--short)
            SHORT="-short"
            shift
            ;;
        -p|--package)
            PACKAGE="$2"
            shift 2
            ;;
        -t|--timeout)
            TIMEOUT="$2"
            shift 2
            ;;
        -w|--watch)
            WATCH=true
            shift
            ;;
        --parallel)
            PARALLEL="$2"
            shift 2
            ;;
        *)
            error "Unknown option: $1"
            ;;
    esac
done

# Change to project root
cd "$PROJECT_ROOT"

log "Starting test suite..."
log "Project root: $PROJECT_ROOT"
log "Testing package: $PACKAGE"

# Check if go is available
if ! command -v go &> /dev/null; then
    error "Go is not installed or not in PATH"
fi

GO_VERSION=$(go version | awk '{print $3}')
log "Using Go version: $GO_VERSION"

# Build test command
TEST_CMD="go test $PACKAGE -count=1 -timeout=$TIMEOUT -p $PARALLEL -parallel 8 $VERBOSE $RACE $COVERAGE $SHORT"

if [ "$WATCH" = true ]; then
    if ! command -v gotestsum &> /dev/null; then
        error "gotestsum is not installed. Install with: go install gotest.tools/gotestsum@latest"
    fi
    log "Running in watch mode..."
    gotestsum --watch -- $PACKAGE -count=1 -timeout=$TIMEOUT $VERBOSE $RACE $SHORT
    exit 0
fi

log "Running command: $TEST_CMD"
echo ""

# Run tests
if eval $TEST_CMD; then
    echo ""
    success "All tests passed!"
    
    # Show coverage if enabled
    if [ -n "$COVERAGE" ]; then
        echo ""
        log "Coverage summary:"
        go tool cover -func=coverage.out | tail -1
        log "To view detailed coverage: go tool cover -html=coverage.out"
    fi
    
    exit 0
else
    echo ""
    error "Tests failed!"
    exit 1
fi
