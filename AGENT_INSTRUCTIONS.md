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

## Specifications

### AppChooser Interface
https://flatpak.github.io/xdg-desktop-portal/docs/doc-org.freedesktop.impl.portal.AppChooser.html

### OpenURI Interface  
https://flatpak.github.io/xdg-desktop-portal/docs/doc-org.freedesktop.portal.OpenURI.html

Agents may use the fetch_webpage tool to retrieve the latest specification details.

## Implementation Goals

### Phase 1: AppChooser (Current)
The current implementation spawns a fuzzel app launcher to allow users to select which application should fulfill `org.freedesktop.impl.portal.AppChooser.ChooseApplication` requests.

### Phase 2: OpenURI Integration (New Plan)
Extend the implementation to handle OpenURI portal requests:

#### OpenURI Methods to Implement:
- `OpenURI` - Open web URLs and non-file URIs
- `OpenFile` - Open local files via file descriptors  
- `OpenDirectory` - Open directories in file manager
- `SchemeSupported` - Check if URI schemes are supported

#### Implementation Strategy:
1. **MIME Type Detection**: Use GIO to determine appropriate applications for URIs/files
2. **Default App Selection**: Choose system default app if available, otherwise first in list
3. **Application Launch**: Use `g_app_info_launch_uris` from GIO package to launch selected app
4. **User Notification**: Show libnotify notification with chosen app name
5. **User Interaction**: Provide notification action to change default app (placeholder for now)

#### Key Dependencies:
- **GIO Package**: `github.com/linuxdeepin/go-gir/gio-2.0` for MIME handling and app launching
- **LibNotify**: For desktop notifications
- **Reference Implementation**: Follow patterns from linuxdeepin/go-lib mime.go

#### User Experience:
- Automatic app selection without user intervention (unless `ask=true`)
- Notification shows chosen app with action to change defaults
- Future: Guide users on setting preferred applications
- Seamless integration with existing desktop workflow

#### No Fuzzel for OpenURI:
Unlike AppChooser, OpenURI methods should work automatically using system defaults, only showing notifications to inform users of the chosen application.