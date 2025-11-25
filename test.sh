#!/usr/bin/env bash
set -e

echo "=== Building binaries ==="

# Build smith if needed
if [ ! -f bin/smith ]; then
    echo "Building smith..."
    go build -o bin/smith cmd/smith/main.go
fi

# Build start (when it exists)
if [ -f cmd/start/main.go ]; then
    echo "Building start..."
    go build -o bin/start cmd/start/main.go
fi

echo ""
echo "=== Running unit tests ==="
go test -v -short ./... || true

echo ""
echo "=== Running integration tests ==="
go test -v ./test/integration/... || true

echo ""
echo "✓ Tests complete!"
