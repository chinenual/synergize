package midi

import (
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
)

var astilectronWindow *astilectron.Window

func initBridge(w *astilectron.Window) {
	astilectronWindow = w
}

// pass an incoming UI msg to MIDI (msg is JSON format)
func SendToMIDI(bytes []byte) (err error) {
	return
}

// pass an incoming MIDI msg to the UI (msg is JSON format)
func SendToUI(name string, payload interface{}) (err error) {
	if astilectronWindow == nil {
		return
	}
	if err = bootstrap.SendMessage(astilectronWindow, name, payload, func(m *bootstrap.MessageIn) { /* ignore response */ }); err != nil {
		return
	}
	return
}
