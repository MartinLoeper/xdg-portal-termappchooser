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
	dbusName      = "org.freedesktop.impl.portal.desktop.termappchooser"
	dbusPath      = "/org/freedesktop/portal/desktop"
	dbusInterface = "org.freedesktop.impl.portal.AppChooser"
)

// AppChooser implements the org.freedesktop.impl.portal.AppChooser interface
type AppChooser struct {
	conn *dbus.Conn
}

// ChooseApplication handles the ChooseApplication D-Bus method
func (ac *AppChooser) ChooseApplication(
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
func (ac *AppChooser) UpdateChoices(
	handle dbus.ObjectPath,
	choices []string,
) *dbus.Error {

	fmt.Println("=== AppChooser.UpdateChoices Called ===")
	fmt.Printf("Handle: %s\n", handle)
	fmt.Printf("Choices: %v\n", choices)
	fmt.Println("=====================================")

	return nil
}

const introspectXML = `
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

func main() {
	// Connect to the session bus
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		log.Fatalf("Failed to connect to session bus: %v", err)
	}
	defer conn.Close()

	// Create AppChooser instance
	appChooser := &AppChooser{conn: conn}

	// Export the object
	err = conn.Expasdfsdafort(appChooser, dbusPath, dbusInterface)
	if err != nil {
		log.Fatalf("Failed to export object: %v", err)
	}

	// Export introspection data
	err = conn.Export(introspect.Introspectable(introspectXML), dbusPath, "org.freedesktop.DBus.Introspectable")
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

	fmt.Printf("XDG Portal AppChooser started on bus name: %s\n", dbusName)
	fmt.Printf("Object path: %s\n", dbusPath)
	fmt.Printf("Interface: %s\n", dbusInterface)
	fmt.Println("Waiting for requests...")

	// Wait for signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("Shutting down...")
}