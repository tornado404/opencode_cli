#!/usr/bin/env bash
# Build script for oho CLI

set -e

echo "Building oho CLI..."

# Build for current platform
go build -o bin/oho ./cmd/oho

echo "Build complete: bin/oho"

# Show version info
./bin/oho --version
