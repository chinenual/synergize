package osc

import (
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"log"
)

var astilectronWindow *astilectron.Window

func OscRegisterBridge(w *astilectron.Window) (err error) {
	astilectronWindow = w
	return
}

type UIMsg struct {
	Field string
	Value int
}

func sendToUI(field string, value int) (err error) {
	log.Printf("SendToUI: %s %d\n", field, value)
	if astilectronWindow == nil {
		return
	}
	if err = bootstrap.SendMessage(astilectronWindow, "updateFromCSurface", UIMsg{field, value},
		func(m *bootstrap.MessageIn) { /* ignore response */ }); err != nil {
		return
	}
	return
}
