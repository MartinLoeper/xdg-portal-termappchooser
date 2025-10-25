#!/bin/bash
set -e

INSTALL_PREFIX="${1:-$HOME/.local}"
BINDIR="$INSTALL_PREFIX/bin"
DBUS_SERVICES_DIR="$HOME/.local/share/dbus-1/services"
SYSTEMD_USER_DIR="$HOME/.config/systemd/user"
PORTAL_DIR="$INSTALL_PREFIX/share/xdg-desktop-portal/portals"

echo "Installing xdg-portal-termappchooser to $INSTALL_PREFIX"

# Create directories
mkdir -p "$BINDIR"
mkdir -p "$DBUS_SERVICES_DIR"
mkdir -p "$SYSTEMD_USER_DIR" 
mkdir -p "$PORTAL_DIR"

# Build the binary
./build.sh

# Install binary
echo "Installing binary to $BINDIR"
cp xdg-portal-termappchooser "$BINDIR/"

# Install D-Bus service file
echo "Installing D-Bus service file"
sed "s|@BINDIR@|$BINDIR|g" data/org.freedesktop.impl.portal.desktop.termappchooser.service.in > "$DBUS_SERVICES_DIR/org.freedesktop.impl.portal.desktop.termappchooser.service"

# Install systemd user service
echo "Installing systemd user service"
sed "s|@BINDIR@|$BINDIR|g" data/xdg-desktop-portal-termappchooser.service.in > "$SYSTEMD_USER_DIR/xdg-desktop-portal-termappchooser.service"

# Install portal configuration
echo "Installing portal configuration"
cp data/termappchooser.portal "$PORTAL_DIR/"

# Reload systemd and D-Bus
echo "Reloading systemd user daemon"
systemctl --user daemon-reload

echo "Restarting xdg-desktop-portal to pick up new backend"
systemctl --user restart xdg-desktop-portal.service || true

echo ""
echo "Installation complete!"
echo "The portal should now be available for AppChooser requests."
echo ""
echo "To test:"
echo "  ./test.sh"
echo ""
echo "To uninstall:"
echo "  ./uninstall.sh"