#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="telegram-audio-bot"
BUILD_DIR="build"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

# Build flags
LDFLAGS="-s -w -X telegram-audio-bot/internal/buildinfo.Version=${VERSION} -X telegram-audio-bot/internal/buildinfo.Commit=${COMMIT} -X telegram-audio-bot/internal/buildinfo.BuildTime=${BUILD_TIME}"

print_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -h, --help     Show this help message"
    echo "  -d, --dev      Development build (with debug info)"
    echo "  -r, --release  Release build (optimized, stripped)"
    echo "  -s, --static   Static build (no external dependencies)"
    echo "  -c, --compress Compress with UPX (requires UPX installed)"
    echo "  -a, --all      Build for all platforms"
    echo "  -t, --test     Run tests before building"
    echo "  --clean        Clean build directory"
    echo ""
    echo "Examples:"
    echo "  $0 --release --compress    # Optimized compressed build"
    echo "  $0 --static --all          # Static builds for all platforms"
    echo "  $0 --dev --test            # Development build with tests"
}

log() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

check_dependencies() {
    log "Checking dependencies..."
    
    if ! command -v go &> /dev/null; then
        error "Go is not installed or not in PATH"
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    log "Using Go version: $GO_VERSION"
    
    if [[ "$COMPRESS" == "true" ]] && ! command -v upx &> /dev/null; then
        warn "UPX not found, compression will be skipped"
        COMPRESS="false"
    fi
}

clean_build() {
    log "Cleaning build directory..."
    rm -rf "$BUILD_DIR"
    mkdir -p "$BUILD_DIR"
}

run_tests() {
    log "Running tests..."
    go test ./...
    success "All tests passed"
}

build_binary() {
    local os="$1"
    local arch="$2"
    local output_name="$3"
    
    log "Building for $os/$arch..."
    
    local build_flags=""
    if [[ "$STATIC" == "true" ]]; then
        build_flags="CGO_ENABLED=0"
    fi
    
    local ldflags="$LDFLAGS"
    if [[ "$STATIC" == "true" ]]; then
        ldflags="$ldflags -extldflags '-static'"
    fi
    
    env GOOS="$os" GOARCH="$arch" $build_flags go build \
        -ldflags="$ldflags" \
        -o "$BUILD_DIR/$output_name" \
        .
    
    if [[ $? -eq 0 ]]; then
        local size=$(ls -lh "$BUILD_DIR/$output_name" | awk '{print $5}')
        success "Built $output_name ($size)"
    else
        error "Failed to build $output_name"
    fi
}

compress_binary() {
    local binary="$1"
    
    if [[ "$COMPRESS" == "true" ]] && command -v upx &> /dev/null; then
        log "Compressing $binary with UPX..."
        local original_size=$(ls -lh "$BUILD_DIR/$binary" | awk '{print $5}')
        
        upx --best "$BUILD_DIR/$binary" &> /dev/null
        
        if [[ $? -eq 0 ]]; then
            local compressed_size=$(ls -lh "$BUILD_DIR/$binary" | awk '{print $5}')
            success "Compressed $binary: $original_size → $compressed_size"
        else
            warn "Failed to compress $binary"
        fi
    fi
}

build_all_platforms() {
    log "Building for all platforms..."
    
    # Common platforms
    declare -A platforms=(
        ["linux/amd64"]="telegram-audio-bot-linux-amd64"
        ["linux/arm64"]="telegram-audio-bot-linux-arm64"
        ["darwin/amd64"]="telegram-audio-bot-darwin-amd64"
        ["darwin/arm64"]="telegram-audio-bot-darwin-arm64"
        ["windows/amd64"]="telegram-audio-bot-windows-amd64.exe"
    )
    
    for platform in "${!platforms[@]}"; do
        IFS='/' read -r os arch <<< "$platform"
        build_binary "$os" "$arch" "${platforms[$platform]}"
        compress_binary "${platforms[$platform]}"
    done
}

build_single() {
    local os=$(go env GOOS)
    local arch=$(go env GOARCH)
    local ext=""
    
    if [[ "$os" == "windows" ]]; then
        ext=".exe"
    fi
    
    build_binary "$os" "$arch" "$BINARY_NAME$ext"
    compress_binary "$BINARY_NAME$ext"
}

show_build_info() {
    log "Build Information:"
    echo "  Version: $VERSION"
    echo "  Commit:  $COMMIT"
    echo "  Time:    $BUILD_TIME"
    echo "  Mode:    $([ "$DEV_MODE" == "true" ] && echo "Development" || echo "Release")"
    echo "  Static:  $([ "$STATIC" == "true" ] && echo "Yes" || echo "No")"
    echo "  Compress: $([ "$COMPRESS" == "true" ] && echo "Yes" || echo "No")"
    echo ""
}

# Parse command line arguments
DEV_MODE="false"
RELEASE_MODE="false"
STATIC="false"
COMPRESS="false"
BUILD_ALL="false"
RUN_TESTS="false"
CLEAN="false"

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            print_usage
            exit 0
            ;;
        -d|--dev)
            DEV_MODE="true"
            ;;
        -r|--release)
            RELEASE_MODE="true"
            ;;
        -s|--static)
            STATIC="true"
            ;;
        -c|--compress)
            COMPRESS="true"
            ;;
        -a|--all)
            BUILD_ALL="true"
            ;;
        -t|--test)
            RUN_TESTS="true"
            ;;
        --clean)
            CLEAN="true"
            ;;
        *)
            error "Unknown option: $1"
            ;;
    esac
    shift
done

# Set default mode if none specified
if [[ "$DEV_MODE" == "false" && "$RELEASE_MODE" == "false" ]]; then
    RELEASE_MODE="true"
fi

# Adjust LDFLAGS for dev mode
if [[ "$DEV_MODE" == "true" ]]; then
    LDFLAGS="-X telegram-audio-bot/internal/buildinfo.Version=${VERSION} -X telegram-audio-bot/internal/buildinfo.Commit=${COMMIT} -X telegram-audio-bot/internal/buildinfo.BuildTime=${BUILD_TIME}"
fi

# Main execution
main() {
    log "Starting build process..."
    
    check_dependencies
    show_build_info
    
    if [[ "$CLEAN" == "true" ]]; then
        clean_build
    else
        mkdir -p "$BUILD_DIR"
    fi
    
    if [[ "$RUN_TESTS" == "true" ]]; then
        run_tests
    fi
    
    if [[ "$BUILD_ALL" == "true" ]]; then
        build_all_platforms
    else
        build_single
    fi
    
    log "Build artifacts:"
    ls -lh "$BUILD_DIR"
    
    success "Build completed successfully!"
}

# Execute main function
main "$@"