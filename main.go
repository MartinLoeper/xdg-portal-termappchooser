package main

/*
#cgo pkg-config: gio-2.0
#include <gio/gio.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/esiqveland/notify"
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

const (
	dbusName           = "org.freedesktop.impl.portal.desktop.termappchooser"
	dbusPath           = "/org/freedesktop/portal/desktop"
	appChooserInterface = "org.freedesktop.impl.portal.AppChooser"
	openURIInterface    = "org.freedesktop.impl.portal.OpenURI"
)

// PortalBackend implements both AppChooser and OpenURI interfaces
type PortalBackend struct {
	conn     *dbus.Conn
	notifier notify.Notifier
}

// ChooseApplication handles the ChooseApplication D-Bus method
func (pb *PortalBackend) ChooseApplication(
	handle dbus.ObjectPath,
	appID string,
	parentWindow string,
	choices []string,
	options map[string]dbus.Variant,
) (uint32, map[string]dbus.Variant, *dbus.Error) {

	fmt.Println("=== AppChooser.ChooseApplication Called ===")
	fmt.Printf("Handle: %s\n", handle)
	fmt.Printf("App ID: %s\n", appID)
	fmt.Printf("Parent Window: %s\n", parentWindow)
	fmt.Printf("Choices: %v\n", choices)
	fmt.Printf("Options: %v\n", options)
	fmt.Println("============================================")

	// For now, just return the first choice if available
	results := make(map[string]dbus.Variant)
	if len(choices) > 0 {
		results["choice"] = dbus.MakeVariant(choices[0])
		fmt.Printf("Returning choice: %s\n", choices[0])
	}

	// Return success response (0) and results
	return 0, results, nil
}

// UpdateChoices handles the UpdateChoices D-Bus method
func (pb *PortalBackend) UpdateChoices(
	handle dbus.ObjectPath,
	choices []string,
) *dbus.Error {

	fmt.Println("=== AppChooser.UpdateChoices Called ===")
	fmt.Printf("Handle: %s\n", handle)
	fmt.Printf("Choices: %v\n", choices)
	fmt.Println("=====================================")

	return nil
}

// OpenURI handles the OpenURI D-Bus method
func (pb *PortalBackend) OpenURI(
	handle dbus.ObjectPath,
	appID string,
	parentWindow string,
	uri string,
	options map[string]dbus.Variant,
) (uint32, map[string]dbus.Variant, *dbus.Error) {

	fmt.Println("=== OpenURI.OpenURI Called ===")
	fmt.Printf("Handle: %s\n", handle)
	fmt.Printf("App ID: %s\n", appID)
	fmt.Printf("Parent Window: %s\n", parentWindow)
	fmt.Printf("URI: %s\n", uri)
	fmt.Printf("Options: %v\n", options)
	fmt.Println("===============================")

	// Get MIME type for the URI
	contentType := pb.getContentTypeForURI(uri)
	fmt.Printf("Content type: %s\n", contentType)

	// Find and launch default application
	appInfo, err := pb.getDefaultAppForContentType(contentType)
	if err != nil {
		fmt.Printf("Error finding app: %v\n", err)
		return 1, nil, dbus.NewError("org.freedesktop.portal.Error.Failed", []interface{}{err.Error()})
	}

	// Launch the application
	err = pb.launchAppWithURI(appInfo, uri)
	if err != nil {
		fmt.Printf("Error launching app: %v\n", err)
		return 1, nil, dbus.NewError("org.freedesktop.portal.Error.Failed", []interface{}{err.Error()})
	}

	// Show notification
	pb.showAppLaunchNotification(appInfo.GetDisplayName(), uri)

	results := make(map[string]dbus.Variant)
	return 0, results, nil
}

// OpenFile handles the OpenFile D-Bus method
func (pb *PortalBackend) OpenFile(
	handle dbus.ObjectPath,
	appID string,
	parentWindow string,
	fd dbus.UnixFD,
	options map[string]dbus.Variant,
) (uint32, map[string]dbus.Variant, *dbus.Error) {

	fmt.Println("=== OpenURI.OpenFile Called ===")
	fmt.Printf("Handle: %s\n", handle)
	fmt.Printf("App ID: %s\n", appID)
	fmt.Printf("Parent Window: %s\n", parentWindow)
	fmt.Printf("File Descriptor: %d\n", fd)
	fmt.Printf("Options: %v\n", options)
	fmt.Println("===============================")

	// Get file path from file descriptor
	filePath, err := pb.getFilePathFromFD(int(fd))
	if err != nil {
		fmt.Printf("Error getting file path: %v\n", err)
		return 1, nil, dbus.NewError("org.freedesktop.portal.Error.Failed", []interface{}{err.Error()})
	}

	// Get MIME type for the file
	contentType := pb.getContentTypeForFile(filePath)
	fmt.Printf("File: %s, Content type: %s\n", filePath, contentType)

	// Find and launch default application
	appInfo, err := pb.getDefaultAppForContentType(contentType)
	if err != nil {
		fmt.Printf("Error finding app: %v\n", err)
		return 1, nil, dbus.NewError("org.freedesktop.portal.Error.Failed", []interface{}{err.Error()})
	}

	// Launch the application
	fileURI := "file://" + filePath
	err = pb.launchAppWithURI(appInfo, fileURI)
	if err != nil {
		fmt.Printf("Error launching app: %v\n", err)
		return 1, nil, dbus.NewError("org.freedesktop.portal.Error.Failed", []interface{}{err.Error()})
	}

	// Show notification
	pb.showAppLaunchNotification(appInfo.GetDisplayName(), filepath.Base(filePath))

	results := make(map[string]dbus.Variant)
	return 0, results, nil
}

// OpenDirectory handles the OpenDirectory D-Bus method  
func (pb *PortalBackend) OpenDirectory(
	handle dbus.ObjectPath,
	appID string,
	parentWindow string,
	fd dbus.UnixFD,
	options map[string]dbus.Variant,
) (uint32, map[string]dbus.Variant, *dbus.Error) {

	fmt.Println("=== OpenURI.OpenDirectory Called ===")
	fmt.Printf("Handle: %s\n", handle)
	fmt.Printf("App ID: %s\n", appID)
	fmt.Printf("Parent Window: %s\n", parentWindow)
	fmt.Printf("File Descriptor: %d\n", fd)
	fmt.Printf("Options: %v\n", options)
	fmt.Println("===================================")

	// Get directory path from file descriptor
	dirPath, err := pb.getFilePathFromFD(int(fd))
	if err != nil {
		fmt.Printf("Error getting directory path: %v\n", err)
		return 1, nil, dbus.NewError("org.freedesktop.portal.Error.Failed", []interface{}{err.Error()})
	}

	// Get default file manager
	appInfo, err := pb.getDefaultAppForContentType("inode/directory")
	if err != nil {
		fmt.Printf("Error finding file manager: %v\n", err)
		return 1, nil, dbus.NewError("org.freedesktop.portal.Error.Failed", []interface{}{err.Error()})
	}

	// Launch file manager
	dirURI := "file://" + dirPath
	err = pb.launchAppWithURI(appInfo, dirURI)
	if err != nil {
		fmt.Printf("Error launching file manager: %v\n", err)
		return 1, nil, dbus.NewError("org.freedesktop.portal.Error.Failed", []interface{}{err.Error()})
	}

	// Show notification
	pb.showAppLaunchNotification(appInfo.GetDisplayName(), filepath.Base(dirPath))

	results := make(map[string]dbus.Variant)
	return 0, results, nil
}

// SchemeSupported handles the SchemeSupported D-Bus method
func (pb *PortalBackend) SchemeSupported(
	scheme string,
	options map[string]dbus.Variant,
) (bool, *dbus.Error) {

	fmt.Printf("=== OpenURI.SchemeSupported Called ===\n")
	fmt.Printf("Scheme: %s\n", scheme)
	fmt.Printf("Options: %v\n", options)
	fmt.Println("=====================================")

	// Check if scheme is commonly supported
	commonSchemes := map[string]bool{
		"http":    true,
		"https":   true,
		"ftp":     true,
		"mailto":  true,
		"file":    true,
		"magnet":  true,
	}
	
	supported := commonSchemes[scheme]
	fmt.Printf("Scheme %s supported: %t\n", scheme, supported)
	return supported, nil
}

// Helper methods for GIO integration

// getContentTypeForURI determines the MIME type for a URI
func (pb *PortalBackend) getContentTypeForURI(uri string) string {
	if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
		return "text/html"
	}
	if strings.HasPrefix(uri, "mailto:") {
		return "message/rfc822"
	}
	if strings.HasPrefix(uri, "ftp://") {
		return "application/octet-stream"
	}
	// Default for unknown schemes
	return "application/octet-stream"
}

// getContentTypeForFile determines the MIME type for a file using GIO
func (pb *PortalBackend) getContentTypeForFile(filePath string) string {
	cFilePath := C.CString(filePath)
	defer C.free(unsafe.Pointer(cFilePath))
	
	gFile := C.g_file_new_for_path(cFilePath)
	defer C.g_object_unref(C.gpointer(gFile))
	
	cAttrs := C.CString("standard::content-type")
	defer C.free(unsafe.Pointer(cAttrs))
	
	fileInfo := C.g_file_query_info(gFile, cAttrs, C.G_FILE_QUERY_INFO_NONE, nil, nil)
	if fileInfo == nil {
		return "application/octet-stream"
	}
	defer C.g_object_unref(C.gpointer(fileInfo))
	
	contentType := C.g_file_info_get_content_type(fileInfo)
	if contentType == nil {
		return "application/octet-stream"
	}
	
	return C.GoString(contentType)
}

// AppInfo represents a GIO AppInfo
type AppInfo struct {
	ptr *C.GAppInfo
}

// GetDisplayName returns the display name of the application
func (ai *AppInfo) GetDisplayName() string {
	displayName := C.g_app_info_get_display_name(ai.ptr)
	if displayName == nil {
		return "Unknown Application"
	}
	return C.GoString(displayName)
}

// getDefaultAppForContentType finds the default application for a MIME type using GIO
func (pb *PortalBackend) getDefaultAppForContentType(contentType string) (*AppInfo, error) {
	cContentType := C.CString(contentType)
	defer C.free(unsafe.Pointer(cContentType))
	
	appInfo := C.g_app_info_get_default_for_type(cContentType, C.gboolean(0))
	if appInfo == nil {
		return nil, fmt.Errorf("no application found for content type: %s", contentType)
	}
	
	return &AppInfo{ptr: appInfo}, nil
}

// launchAppWithURI launches an application with the given URI using GIO
func (pb *PortalBackend) launchAppWithURI(appInfo *AppInfo, uri string) error {
	cUri := C.CString(uri)
	defer C.free(unsafe.Pointer(cUri))
	
	// Create GList with the URI
	gList := C.g_list_append(nil, C.gpointer(unsafe.Pointer(cUri)))
	defer C.g_list_free(gList)
	
	context := C.g_app_launch_context_new()
	defer C.g_object_unref(C.gpointer(context))
	
	var gError *C.GError
	success := C.g_app_info_launch_uris(appInfo.ptr, gList, context, &gError)
	
	if gError != nil {
		errorMsg := C.GoString(gError.message)
		C.g_error_free(gError)
		return fmt.Errorf("failed to launch application: %s", errorMsg)
	}
	
	if success == C.FALSE {
		return fmt.Errorf("application launch returned false")
	}
	
	return nil
}

// getFilePathFromFD gets the file path from a file descriptor
func (pb *PortalBackend) getFilePathFromFD(fd int) (string, error) {
	// Read the symlink from /proc/self/fd/
	linkPath := fmt.Sprintf("/proc/self/fd/%d", fd)
	filePath, err := os.Readlink(linkPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file descriptor %d: %v", fd, err)
	}
	return filePath, nil
}

// showAppLaunchNotification shows a desktop notification about the launched app
func (pb *PortalBackend) showAppLaunchNotification(appName, target string) {
	if pb.notifier == nil {
		return
	}

	notification := notify.Notification{
		AppName: "XDG Portal",
		Summary: fmt.Sprintf("Opened with %s", appName),
		Body:    fmt.Sprintf("Target: %s", target),
	}

	// Send notification and handle action clicks
	id, err := pb.notifier.SendNotification(notification)
	if err != nil {
		fmt.Printf("Failed to send notification: %v\n", err)
		return
	}

	fmt.Printf("Notification sent with ID: %d\n", id)

	// TODO: Handle action clicks
	// For now, just log that we would handle it
	fmt.Println("Note: Click 'Change Default App' to set preferences (placeholder implementation)")
}

const appChooserIntrospectXML = `
<interface name="org.freedesktop.impl.portal.AppChooser">
	<method name="ChooseApplication">
		<arg type="o" name="handle" direction="in"/>
		<arg type="s" name="app_id" direction="in"/>
		<arg type="s" name="parent_window" direction="in"/>
		<arg type="as" name="choices" direction="in"/>
		<arg type="a{sv}" name="options" direction="in"/>
		<arg type="u" name="response" direction="out"/>
		<arg type="a{sv}" name="results" direction="out"/>
	</method>
	<method name="UpdateChoices">
		<arg type="o" name="handle" direction="in"/>
		<arg type="as" name="choices" direction="in"/>
	</method>
</interface>` + introspect.IntrospectDataString

const openURIIntrospectXML = `
<interface name="org.freedesktop.impl.portal.OpenURI">
	<method name="OpenURI">
		<arg type="o" name="handle" direction="in"/>
		<arg type="s" name="app_id" direction="in"/>
		<arg type="s" name="parent_window" direction="in"/>
		<arg type="s" name="uri" direction="in"/>
		<arg type="a{sv}" name="options" direction="in"/>
		<arg type="u" name="response" direction="out"/>
		<arg type="a{sv}" name="results" direction="out"/>
	</method>
	<method name="OpenFile">
		<arg type="o" name="handle" direction="in"/>
		<arg type="s" name="app_id" direction="in"/>
		<arg type="s" name="parent_window" direction="in"/>
		<arg type="h" name="fd" direction="in"/>
		<arg type="a{sv}" name="options" direction="in"/>
		<arg type="u" name="response" direction="out"/>
		<arg type="a{sv}" name="results" direction="out"/>
	</method>
	<method name="OpenDirectory">
		<arg type="o" name="handle" direction="in"/>
		<arg type="s" name="app_id" direction="in"/>
		<arg type="s" name="parent_window" direction="in"/>
		<arg type="h" name="fd" direction="in"/>
		<arg type="a{sv}" name="options" direction="in"/>
		<arg type="u" name="response" direction="out"/>
		<arg type="a{sv}" name="results" direction="out"/>
	</method>
	<method name="SchemeSupported">
		<arg type="s" name="scheme" direction="in"/>
		<arg type="a{sv}" name="options" direction="in"/>
		<arg type="b" name="supported" direction="out"/>
	</method>
</interface>` + introspect.IntrospectDataString

func main() {
	// GLib type system is automatically initialized in modern versions

	// Connect to the session bus
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		log.Fatalf("Failed to connect to session bus: %v", err)
	}
	defer conn.Close()

	// Initialize notification system
	notifier, err := notify.New(conn)
	if err != nil {
		fmt.Printf("Warning: Failed to initialize notifications: %v\n", err)
	}

	// Create portal backend instance
	backend := &PortalBackend{
		conn:     conn,
		notifier: notifier,
	}

	// Export AppChooser interface
	err = conn.Export(backend, dbusPath, appChooserInterface)
	if err != nil {
		log.Fatalf("Failed to export AppChooser interface: %v", err)
	}

	// Export OpenURI interface
	err = conn.Export(backend, dbusPath, openURIInterface)
	if err != nil {
		log.Fatalf("Failed to export OpenURI interface: %v", err)
	}

	// Export introspection data
	err = conn.Export(introspect.Introspectable(appChooserIntrospectXML), dbusPath, "org.freedesktop.DBus.Introspectable")
	if err != nil {
		log.Fatalf("Failed to export introspectable: %v", err)
	}

	// Request the well-known name
	reply, err := conn.RequestName(dbusName, dbus.NameFlagDoNotQueue)
	if err != nil {
		log.Fatalf("Failed to request name: %v", err)
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		log.Fatalf("Name already taken")
	}

	fmt.Printf("XDG Portal Backend started on bus name: %s\n", dbusName)
	fmt.Printf("Object path: %s\n", dbusPath)
	fmt.Printf("Interfaces: %s, %s\n", appChooserInterface, openURIInterface)
	fmt.Println("Waiting for requests...")

	// Wait for signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("Shutting down...")
	
	// Clean up
	if notifier != nil {
		notifier.Close()
	}
}