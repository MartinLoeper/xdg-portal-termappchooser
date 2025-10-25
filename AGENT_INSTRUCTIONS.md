# Agent Instructions: XDG Portal Terminal App Chooser

## Project Overview

We are writing a Go application that implements the `org.freedesktop.impl.portal.AppChooser` XDG Desktop Portal interface using the `godbus/dbus` library. The application should intercept app chooser requests and print them out on the console.

## Technical Implementation

- **Language**: Go
- **D-Bus Library**: `godbus/dbus`
- **Interfaces**: 
  - `org.freedesktop.impl.portal.AppChooser`
  - `org.freedesktop.impl.portal.OpenURI`
- **Initial Functionality**: Console logging of intercepted requests

## Specification

The XDG Desktop Portal AppChooser specification can be found at:
https://flatpak.github.io/xdg-desktop-portal/docs/doc-org.freedesktop.impl.portal.AppChooser.html

Agents may use the fetch_webpage tool to retrieve the latest specification details.

## End Goal

The final implementation should spawn a fuzzel app launcher to allow users to select which application should fulfill:
- `org.freedesktop.impl.portal.AppChooser.ChooseApplication` requests
- `org.freedesktop.impl.portal.OpenURI.OpenURI` requests
- `org.freedesktop.impl.portal.OpenURI.OpenFile` requests

This ensures that both explicit app chooser requests AND URI/file opening requests get routed through our fuzzel-based selection interface.