package main

import (
	"log"

	"github.com/godbus/dbus/v5"
)

const (
	busName        = "org.freedesktop.portal.Desktop"
	objectPath     = "/org/freedesktop/portal/desktop"
	methodName     = "org.freedesktop.portal.FileChooser.OpenFile"
	requestIFace   = "org.freedesktop.portal.Request"
	responseSignal = "Response"
	handleToken    = "revelation"
	dialogTitle    = "Choose file"
)

func SelectFile() string {
	conn := connectDBus()
	defer conn.Close()

	responsePath := openFileDialog(conn)
	setupSignalHandler(conn, responsePath)

	return processSignal(<-waitForSignal(conn), responsePath)
}

func connectDBus() *dbus.Conn {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		log.Fatalf("Failed to connect to session bus: %v", err)
	}
	return conn
}

func openFileDialog(conn *dbus.Conn) dbus.ObjectPath {
	options := map[string]dbus.Variant{
		"handle_token": dbus.MakeVariant(handleToken),
		"title":        dbus.MakeVariant(dialogTitle),
	}

	call := conn.Object(busName, objectPath).Call(methodName, 0, "", "", options)
	if call.Err != nil {
		log.Fatalf("Failed to trigger file picker: %v", call.Err)
	}

	return call.Body[0].(dbus.ObjectPath)
}

func setupSignalHandler(conn *dbus.Conn, path dbus.ObjectPath) {
	err := conn.AddMatchSignal(
		dbus.WithMatchInterface(requestIFace),
		dbus.WithMatchMember(responseSignal),
		dbus.WithMatchPathNamespace(path),
	)
	if err != nil {
		log.Fatalf("Failed to add signal match: %v", err)
	}
}

func waitForSignal(conn *dbus.Conn) <-chan *dbus.Signal {
	ch := make(chan *dbus.Signal, 1)
	conn.Signal(ch)
	return ch
}

func processSignal(signal *dbus.Signal, expectedPath dbus.ObjectPath) string {
	if signal.Path != expectedPath || signal.Name != requestIFace+"."+responseSignal {
		return ""
	}

	if len(signal.Body) < 2 {
		// nothing selected
		return ""
	}

	results, ok := signal.Body[1].(map[string]dbus.Variant)
	if !ok {
		// invalid response
		return ""
	}

	urisVariant, exists := results["uris"]
	if !exists {
		// nothing selected
		return ""
	}

	uris, ok := urisVariant.Value().([]string)
	if ok && len(uris) > 0 {
		return uris[0]
	}
	return ""
}
