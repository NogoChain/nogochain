#!/bin/bash

set -e

echo "Verifying dependencies..."

# Check Go version
go_version=$(go version | awk '{print $3}')
echo "Go version: $go_version"

# Check for CGO_ENABLED
echo "CGO_ENABLED: $(go env CGO_ENABLED)"
if [ "$(go env CGO_ENABLED)" != "0" ]; then
    echo "WARNING: CGO is enabled. Please set CGO_ENABLED=0"
fi

# Verify go.mod
echo "Verifying go.mod..."
go mod verify

# Check for prohibited dependencies
prohibited_deps=("github.com/ethereum/go-ethereum" "github.com/etclabscore/core-geth")

for dep in "${prohibited_deps[@]}"; do
    if go list -m "$dep" >/dev/null 2>&1; then
        echo "ERROR: Prohibited dependency found: $dep"
        exit 1
    fi
done

echo "Dependencies verified successfully!"
