# Build Info Embedding Summary

## Mechanism
The codebase embeds build information (version, commit hash, build time) into the final binary using Go's `-ldflags` with variable injection.

## Implementation
- **Package**: `internal/buildinfo/buildinfo.go` defines variables (`Version`, `Commit`, `BuildTime`) with default values
- **Build Script**: `scripts/build.sh` uses `go build -ldflags="-X telegram-audio-bot/internal/buildinfo.Version=${VERSION}..."` to inject values at compile time
- **Runtime Info**: Go version and platform are captured at runtime using `runtime` package

## Usage
- Version: Git tag/describe (fallback: "dev")
- Commit: Short git hash (fallback: "unknown")  
- Build Time: UTC timestamp at build time
- Displayed on startup via `buildinfo.Print()` in `main.go`