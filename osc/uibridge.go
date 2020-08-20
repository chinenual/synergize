package osc

import (
	"encoding/json"
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
	var strval string
	if err = bootstrap.SendMessage(astilectronWindow, "updateFromCSurface", UIMsg{field, value},
		func(m *bootstrap.MessageIn) {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &strval); err != nil {
				log.Printf(" SendToUI failed to decode json response : %s %d: %v\n", field, value, err)
				return
			}
			// feedback to CSurface:
			if err = oscSendString("/stringval", strval); err != nil {
				return
			}

		}); err != nil {
		return
	}
	return
}
