# XDG Desktop Portal Architecture

## Overview

This application integrates into the XDG Desktop Portal ecosystem as a specialized backend implementation for the AppChooser interface. It provides application selection functionality through the standardized D-Bus portal framework.

## XDG Desktop Portal Ecosystem

```
┌─────────────────┐    ┌──────────────────────┐    ┌─────────────────────┐
│   Application   │    │  xdg-desktop-portal  │    │ Portal Backends     │
│                 │    │     (Frontend)       │    │                     │
│ ┌─────────────┐ │    │                      │    │ ┌─────────────────┐ │
│ │   Flatpak   │ │◄──►│ ┌──────────────────┐ │◄──►│ │ gtk/kde/wlroots │ │
│ │   Snap      │ │    │ │ Interface Router │ │    │ │                 │ │
│ │   Native    │ │    │ │                  │ │    │ │ File Chooser    │ │
│ └─────────────┘ │    │ └──────────────────┘ │    │ │ Screenshot      │ │
└─────────────────┘    │                      │    │ │ Notification    │ │
                       │ ┌──────────────────┐ │    │ │ ...             │ │
                       │ │   Permission     │ │    │ └─────────────────┘ │
                       │ │   Management     │ │    │                     │
                       │ └──────────────────┘ │    │ ┌─────────────────┐ │
                       └──────────────────────┘    │ │ termappchooser  │ │
                                                   │ │                 │ │
                                                   │ │ AppChooser Only │ │
                                                   │ └─────────────────┘ │
                                                   └─────────────────────┘
```

## D-Bus Communication Flow

### 1. Service Registration
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────────┐
│    systemd      │    │    D-Bus Daemon  │    │  termappchooser     │
│                 │    │                  │    │                     │
│ User Session    │───►│ Session Bus      │◄───│ Claims Bus Name:    │
│ Starts Services │    │                  │    │ org.freedesktop.   │
│                 │    │                  │    │ impl.portal.desktop │
│                 │    │                  │    │ .termappchooser     │
└─────────────────┘    └──────────────────┘    └─────────────────────┘
```

### 2. Portal Configuration
```
/usr/share/xdg-desktop-portal/portals/termappchooser.portal
┌─────────────────────────────────────────────────────────┐
│ [portal]                                                │
│ DBusName=org.freedesktop.impl.portal.desktop.termapp... │
│ Interfaces=org.freedesktop.impl.portal.AppChooser      │
│ UseIn=hyprland;sway;river                               │
└─────────────────────────────────────────────────────────┘
                               │
                               ▼
        xdg-desktop-portal reads configuration and routes
        AppChooser requests to termappchooser backend
```

### 3. Request Flow
```
┌─────────────┐ D-Bus Call   ┌──────────────────┐ Routes to  ┌─────────────────┐
│ Application │─────────────►│ xdg-desktop-     │───────────►│ termappchooser  │
│             │              │ portal           │            │                 │
│ Wants to    │              │                  │            │ Shows fuzzel    │
│ open file   │              │ Checks config:   │            │ for app         │
│ with app    │              │ AppChooser →     │            │ selection       │
│             │              │ termappchooser   │            │                 │
└─────────────┘              └──────────────────┘            └─────────────────┘
        ▲                                                              │
        │                    ┌──────────────────┐                     │
        └────────────────────│ Response with    │◄────────────────────┘
                             │ selected app ID  │
                             └──────────────────┘
```

## D-Bus Interface Implementation

### Interface: `org.freedesktop.impl.portal.AppChooser`

#### Method: ChooseApplication
```
ChooseApplication(
    handle: ObjectPath,           // Request handle for cancellation
    app_id: String,              // Calling application ID  
    parent_window: String,        // Window identifier for modal dialogs
    choices: Array<String>,       // Available application IDs
    options: Map<String,Variant>  // Additional options (content_type, etc.)
) → (response: UInt32, results: Map<String,Variant>)
```

**Our Implementation Flow:**
1. Receive D-Bus method call
2. Parse application choices and options
3. Launch fuzzel with formatted app list
4. Return selected application ID
5. Handle cancellation via Request interface

#### Method: UpdateChoices
```
UpdateChoices(
    handle: ObjectPath,     // Active request handle
    choices: Array<String>  // Updated application list
)
```

## Portal Backend Requirements

### systemd User Service (Primary Method)
Modern portal backends use systemd with D-Bus integration for activation:
```ini
# ~/.config/systemd/user/xdg-desktop-portal-termappchooser.service
[Unit]
Description=Portal service (termappchooser implementation)
PartOf=graphical-session.target

[Service]
Type=dbus
BusName=org.freedesktop.impl.portal.desktop.termappchooser
ExecStart=/path/to/xdg-portal-termappchooser
Restart=on-failure
```

With `Type=dbus`, systemd handles D-Bus activation automatically - no separate D-Bus service file needed!

### Legacy D-Bus Service File (Optional)
Only required for systems without systemd D-Bus integration:

```ini
# ~/.local/share/dbus-1/services/org.freedesktop.impl.portal.desktop.termappchooser.service
[D-BUS Service]
Name=org.freedesktop.impl.portal.desktop.termappchooser
Exec=/path/to/xdg-portal-termappchooser
SystemdService=xdg-desktop-portal-termappchooser.service
```

## Integration Points

### 1. Desktop Environment Detection
- Portal config specifies `UseIn=hyprland;sway;river`
- xdg-desktop-portal selects backend based on `XDG_CURRENT_DESKTOP`
- Falls back to other backends for unsupported environments

### 2. Application Discovery
- Reads `.desktop` files from standard locations
- Filters by MIME type associations
- Respects application priorities and defaults

### 3. User Interface
- Spawns fuzzel as external process
- Formats application list for display
- Handles user selection and cancellation
- Returns result via D-Bus response

## Security Model

### Sandboxing
- Portal runs outside application sandbox
- Can access system applications and configuration
- Mediates between sandboxed apps and system

### Permission Model
- No explicit permissions required for AppChooser
- Inherits desktop environment's application access
- User controls selection through interactive dialog

## Error Handling

### D-Bus Error Responses
- `org.freedesktop.portal.Error.Cancelled` - User cancelled
- `org.freedesktop.portal.Error.InvalidArgument` - Bad parameters
- `org.freedesktop.portal.Error.Failed` - General failure

### Graceful Degradation
- Falls back to system default if fuzzel unavailable
- Logs errors for debugging
- Maintains portal availability for other interfaces