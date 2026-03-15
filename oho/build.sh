#!/usr/bin/env bash
# Build script for oho CLI

set -e

echo "Building oho CLI..."

# Get version from git tag
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)

echo "Version: $VERSION"
echo "Commit: $COMMIT"
echo "Date: $DATE"

# Build with version info
go build -ldflags "-s -w -X main.Version=$VERSION -X main.Commit=$COMMIT -X main.Date=$DATE" -o bin/oho ./cmd

echo "Build complete: bin/oho"

# Show version info
./bin/oho --version
