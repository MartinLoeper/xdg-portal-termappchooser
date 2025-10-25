package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	conn *dbus.Conn
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

	// TODO: Use fuzzel to let user choose application for this URI
	// For now, just return success
	results := make(map[string]dbus.Variant)
	
	// Return success response (0) and results
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

	// TODO: Use fuzzel to let user choose application for this file
	// For now, just return success
	results := make(map[string]dbus.Variant)
	
	// Return success response (0) and results
	return 0, results, nil
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
</interface>` + introspect.IntrospectDataString

func main() {
	// Connect to the session bus
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		log.Fatalf("Failed to connect to session bus: %v", err)
	}
	defer conn.Close()

	// Create portal backend instance
	backend := &PortalBackend{conn: conn}

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

	// Export introspection data for AppChooser
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
}