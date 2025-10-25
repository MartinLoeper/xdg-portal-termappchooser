#!/bin/bash
set -e

echo "Unregistering development D-Bus service..."

DBUS_SERVICES_DIR="$HOME/.local/share/dbus-1/services"
SERVICE_FILE="$DBUS_SERVICES_DIR/org.freedesktop.impl.portal.desktop.termappchooser.service"

if [ -f "$SERVICE_FILE" ]; then
    rm "$SERVICE_FILE"
    echo "Removed: $SERVICE_FILE"
else
    echo "Service file not found: $SERVICE_FILE"
fi

# Kill any running instances
pkill -f "xdg-portal-termappchooser" || true

echo "Development service unregistered!"