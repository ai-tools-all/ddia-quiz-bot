#!/bin/bash

# Build script for quiz-evaluator

set -e

echo "🔨 Building Quiz Evaluator..."

# Clean previous builds
rm -f quiz-eval quiz-eval-*

# Get dependencies
echo "📦 Getting dependencies..."
go mod download

# Run tests
echo "🧪 Running tests..."
go test ./... || true

# Build for current platform
echo "🏗️  Building binary..."
go build -o quiz-eval -ldflags="-s -w" .

# Build for multiple platforms (optional)
if [ "$1" == "--all" ]; then
    echo "🌍 Building for multiple platforms..."
    
    # Linux AMD64
    GOOS=linux GOARCH=amd64 go build -o quiz-eval-linux-amd64 -ldflags="-s -w" .
    
    # Mac AMD64
    GOOS=darwin GOARCH=amd64 go build -o quiz-eval-darwin-amd64 -ldflags="-s -w" .
    
    # Mac ARM64 (M1/M2)
    GOOS=darwin GOARCH=arm64 go build -o quiz-eval-darwin-arm64 -ldflags="-s -w" .
    
    # Windows AMD64
    GOOS=windows GOARCH=amd64 go build -o quiz-eval-windows-amd64.exe -ldflags="-s -w" .
    
    echo "✅ Multi-platform build complete!"
    ls -lh quiz-eval*
else
    echo "✅ Build complete!"
    echo "📍 Binary: ./quiz-eval"
    echo ""
    echo "To build for all platforms, run: ./build.sh --all"
fi

# Make the binary executable
chmod +x quiz-eval*

echo ""
echo "🚀 Ready to use! Run './quiz-eval --help' to get started."
