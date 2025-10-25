#!/bin/bash
# Test script to trigger the OpenURI interface

set -e

echo "Testing OpenURI D-Bus interface..."

# Test the OpenURI method
echo "Testing OpenURI with a web URL..."
dbus-send \
    --session \
    --print-reply \
    --dest=org.freedesktop.impl.portal.desktop.termappchooser \
    /org/freedesktop/portal/desktop \
    org.freedesktop.impl.portal.OpenURI.OpenURI \
    objpath:/org/freedesktop/portal/desktop/request/test \
    string:"test.app" \
    string:"" \
    string:"https://www.example.com" \
    dict:string:variant:ask,variant:boolean:true

echo ""
echo "Testing with a file URI..."
dbus-send \
    --session \
    --print-reply \
    --dest=org.freedesktop.impl.portal.desktop.termappchooser \
    /org/freedesktop/portal/desktop \
    org.freedesktop.impl.portal.OpenURI.OpenURI \
    objpath:/org/freedesktop/portal/desktop/request/test2 \
    string:"test.app" \
    string:"" \
    string:"file:///etc/passwd" \
    dict:string:variant:ask,variant:boolean:true

echo "OpenURI tests completed!"