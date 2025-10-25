#!/usr/bin/env bash
set -e

echo "Running xdg-portal-termappchooser..."

# Build first if binary doesn't exist
if [ ! -f ./xdg-portal-termappchooser ]; then
    echo "Binary not found, building first..."
    ./build.sh
fi

echo "Starting D-Bus service..."
echo "Press Ctrl+C to stop"

# Run the binary
./xdg-portal-termappchooser