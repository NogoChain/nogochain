#!/bin/bash

set -e

# Set CGO_ENABLED to 0
export CGO_ENABLED=0

# Build directory
build_dir="./build"
mkdir -p "$build_dir"

echo "Building NogoChain..."

# Build the main binary
go build -o "$build_dir/nogochain" ./cmd/nogochain

# Build the pool binary
go build -o "$build_dir/nogopool" ./cmd/nogopool

echo "Build completed successfully!"
echo "Binaries available at: $build_dir/"
echo "- nogochain: Main blockchain node"
echo "- nogopool: Mining pool server"
