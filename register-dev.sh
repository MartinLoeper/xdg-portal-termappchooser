#!/usr/bin/env bash
set -e

# Development D-Bus registration script
CURRENT_DIR="$(pwd)"
BINARY_PATH="$CURRENT_DIR/xdg-portal-termappchooser"
DBUS_SERVICES_DIR="$HOME/.local/share/dbus-1/services"

echo "Registering xdg-portal-termappchooser for development..."

# Create D-Bus services directory if it doesn't exist
mkdir -p "$DBUS_SERVICES_DIR"

# Create temporary D-Bus service file pointing to current directory
cat > "$DBUS_SERVICES_DIR/org.freedesktop.impl.portal.desktop.termappchooser.service" << EOF
[D-BUS Service]
Name=org.freedesktop.impl.portal.desktop.termappchooser
Exec=$BINARY_PATH
EOF

echo "D-Bus service registered at: $DBUS_SERVICES_DIR/org.freedesktop.impl.portal.desktop.termappchooser.service"
echo "Binary path: $BINARY_PATH"

# Build the binary if it doesn't exist
if [ ! -f "$BINARY_PATH" ]; then
    echo "Building binary..."
    ./build.sh
fi

echo ""
echo "Development registration complete!"
echo "Now you can:"
echo "1. Run './run.sh' to start the service manually, OR"
echo "2. Let D-Bus auto-start it when needed"
echo "3. Test with './test.sh'"