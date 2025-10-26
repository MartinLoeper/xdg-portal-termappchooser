#!/bin/bash
# Test script to verify OpenURI interface functionality

set -e

echo "Testing OpenURI D-Bus interface..."
echo "=================================="

# Test SchemeSupported method
echo "1. Testing SchemeSupported..."
dbus-send \
    --session \
    --print-reply \
    --dest=org.freedesktop.impl.portal.desktop.termappchooser \
    /org/freedesktop/portal/desktop \
    org.freedesktop.impl.portal.OpenURI.SchemeSupported \
    string:"https" \
    dict:string:variant:

echo ""

# Test OpenURI method with HTTP URL
echo "2. Testing OpenURI with HTTP URL..."
dbus-send \
    --session \
    --print-reply \
    --dest=org.freedesktop.impl.portal.desktop.termappchooser \
    /org/freedesktop/portal/desktop \
    org.freedesktop.impl.portal.OpenURI.OpenURI \
    objpath:/org/freedesktop/portal/desktop/request/test1 \
    string:"test.app" \
    string:"" \
    string:"https://www.example.com" \
    dict:string:variant:ask,variant:boolean:false

echo ""

# Test OpenURI method with mailto
echo "3. Testing OpenURI with mailto..."
dbus-send \
    --session \
    --print-reply \
    --dest=org.freedesktop.impl.portal.desktop.termappchooser \
    /org/freedesktop/portal/desktop \
    org.freedesktop.impl.portal.OpenURI.OpenURI \
    objpath:/org/freedesktop/portal/desktop/request/test2 \
    string:"test.app" \
    string:"" \
    string:"mailto:test@example.com" \
    dict:string:variant:ask,variant:boolean:false

echo ""

# Create a temporary test file for OpenFile
TEST_FILE="/tmp/test-portal-file.txt"
echo "This is a test file for the OpenURI portal" > "$TEST_FILE"

echo "4. Testing OpenFile with temporary file..."
echo "Created test file: $TEST_FILE"

# Note: OpenFile requires a file descriptor, which is complex to create via dbus-send
# This would normally be called by applications that already have the file open
echo "OpenFile test requires file descriptor handling - normally called by applications"

echo ""
echo "5. Testing with a directory..."
# Similar note for OpenDirectory - needs proper file descriptor
echo "OpenDirectory test requires file descriptor handling - normally called by applications" 

echo ""
echo "OpenURI tests completed!"
echo "Check the portal application console for detailed logs."

# Clean up
rm -f "$TEST_FILE"