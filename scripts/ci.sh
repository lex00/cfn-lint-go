#!/bin/bash
set -euo pipefail

echo "=== Running CI checks ==="

echo ""
echo "--- go mod tidy ---"
go mod tidy

echo ""
echo "--- go vet ---"
go vet ./...

echo ""
echo "--- go build ---"
go build -v ./...

echo ""
echo "--- go test ---"
go test -v -race ./...

echo ""
echo "--- golangci-lint ---"
if command -v golangci-lint &> /dev/null; then
    golangci-lint run
else
    echo "golangci-lint not installed, skipping"
fi

echo ""
echo "=== All CI checks passed ==="
