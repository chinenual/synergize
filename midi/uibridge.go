package midi

import (
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"log"
)

var astilectronWindow *astilectron.Window

func RegisterBridge(w *astilectron.Window) (err error) {
	astilectronWindow = w
	return
}

func SendToMIDI(field string, value int) (err error) {
	return csSendEvent(field, uint8(value))
}

type UIMsg struct {
	Field string
	Value int
}

func SendToUI(field string, value int) (err error) {
	log.Printf("SendToUI: %s %d\n", field, value)
	if astilectronWindow == nil {
		return
	}
	if err = bootstrap.SendMessage(astilectronWindow, "updateFromMIDI", UIMsg{field, value},
		func(m *bootstrap.MessageIn) { /* ignore response */ }); err != nil {
		return
	}
	return
}
