#!/usr/bin/env bash
# Test script to trigger the AppChooser interface

set -e

echo "Testing AppChooser D-Bus interface..."

# Test the ChooseApplication method
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

echo "Test completed!"