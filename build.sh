#!/usr/bin/env bash
set -e

echo "Building xdg-portal-termappchooser..."

# Initialize go module if go.sum doesn't exist
if [ ! -f go.sum ]; then
    echo "Downloading dependencies..."
    go mod tidy
fi

# Build the binary
echo "Compiling..."
go build -o xdg-portal-termappchooser .

echo "Build complete! Binary: ./xdg-portal-termappchooser"