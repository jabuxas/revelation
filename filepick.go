package main

import (
	"fmt"
	"log"
	"os"

	"github.com/godbus/dbus/v5"
)

func SelectFile() string {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		os.Exit(1)
	}
	defer conn.Close()

	token := "revelation"
	options := map[string]dbus.Variant{
		"handle_token": dbus.MakeVariant(token),
		"title":        dbus.MakeVariant("Choose file"),
	}

	obj := conn.Object("org.freedesktop.portal.Desktop", "/org/freedesktop/portal/desktop")
	call := obj.Call("org.freedesktop.portal.FileChooser.OpenFile", 0, "revelation", "", options)
	if call.Err != nil {
		log.Fatalf("Failed to trigger file picker: %v", call.Err)
	}

	replyPath := call.Body[0].(dbus.ObjectPath)

	err = conn.AddMatchSignal(
		dbus.WithMatchOption("interface", "org.freedesktop.portal.Request"),
		dbus.WithMatchOption("member", "Response"),
		dbus.WithMatchOption("path", string(replyPath)),
	)
	if err != nil {
		log.Fatalf("Failed to add match signal: %v", err)
	}

	c := make(chan *dbus.Signal, 1)
	conn.Signal(c)

	for signal := range c {
		if signal.Path == replyPath && signal.Name == "org.freedesktop.portal.Request.Response" {
			reply := signal.Body
			if len(reply) > 1 {
				results, ok := reply[1].(map[string]dbus.Variant)
				if ok {
					if urisVariant, exists := results["uris"]; exists {
						uris, ok := urisVariant.Value().([]string)
						if ok && len(uris) > 0 {
							return uris[0]
						} else {
							return ""
						}
					}
				}
			}
			break
		}
	}
	return ""
}
