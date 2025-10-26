#!/bin/bash
# Comprehensive test script for both AppChooser and OpenURI

set -e

echo "Comprehensive XDG Portal Test Suite"
echo "===================================="

# Function to check if service is running
check_service() {
    if ! busctl --user list | grep -q "termappchooser"; then
        echo "ERROR: Portal service is not running!"
        echo "Start it with: ./run.sh (in another terminal)"
        exit 1
    fi
    echo "✓ Portal service is running"
}

# Test 1: Service availability
echo "1. Checking service availability..."
check_service

# Test 2: AppChooser functionality  
echo ""
echo "2. Testing AppChooser interface..."
dbus-send \
    --session \
    --print-reply \
    --dest=org.freedesktop.impl.portal.desktop.termappchooser \
    /org/freedesktop/portal/desktop \
    org.freedesktop.impl.portal.AppChooser.ChooseApplication \
    objpath:/org/freedesktop/portal/desktop/request/test \
    string:"test.app" \
    string:"" \
    array:string:"firefox","chromium","gedit" \
    dict:string:variant:last_choice,variant:string:"firefox"

# Test 3: SchemeSupported
echo ""
echo "3. Testing SchemeSupported..."
dbus-send \
    --session \
    --print-reply \
    --dest=org.freedesktop.impl.portal.desktop.termappchooser \
    /org/freedesktop/portal/desktop \
    org.freedesktop.impl.portal.OpenURI.SchemeSupported \
    string:"https" \
    dict:string:variant:

# Test 4: OpenURI with web URL
echo ""
echo "4. Testing OpenURI with web URL..."
dbus-send \
    --session \
    --print-reply \
    --dest=org.freedesktop.impl.portal.desktop.termappchooser \
    /org/freedesktop/portal/desktop \
    org.freedesktop.impl.portal.OpenURI.OpenURI \
    objpath:/org/freedesktop/portal/desktop/request/web \
    string:"test.app" \
    string:"" \
    string:"https://www.github.com/MartinLoeper/xdg-portal-termappchooser" \
    dict:string:variant:ask,variant:boolean:false

echo ""
echo "5. Testing OpenURI with mailto..."
dbus-send \
    --session \
    --print-reply \
    --dest=org.freedesktop.impl.portal.desktop.termappchooser \
    /org/freedesktop/portal/desktop \
    org.freedesktop.impl.portal.OpenURI.OpenURI \
    objpath:/org/freedesktop/portal/desktop/request/mail \
    string:"test.app" \
    string:"" \
    string:"mailto:example@test.com?subject=Portal%20Test" \
    dict:string:variant:ask,variant:boolean:false

echo ""
echo "✓ All tests completed successfully!"
echo ""
echo "Check the portal console for detailed execution logs."
echo "You should have seen notifications if libnotify is working."
echo ""
echo "Expected behavior:"
echo "- AppChooser: Returns first choice from provided list"  
echo "- SchemeSupported: Returns true/false for URL scheme support"
echo "- OpenURI: Launches default application and shows notification"