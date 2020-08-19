package midi

import (
	"log"
	"strings"
)

type inboundField struct {
	name  string
	scale int
}

type inboundMidiMap struct {
	// ccMap[CC#] -> fieldname
	ccMap map[uint8]inboundField
	// noteMap[note#] -> fieldname
	noteMap map[uint8]inboundField
	// polyMap[note#] -> fieldname
	polyMap map[uint8]inboundField
}

// inboundChannelMidiMap[channel] -> inboundMidiMap
var inboundChannelMap = make(map[uint8]inboundMidiMap)

type EventType int

const (
	Cc EventType = iota
	Note
	Poly
)

type outboundField struct {
	eventtype EventType
	index     uint8
	scale     int
}

// outboundMidiMap[fieldname]->outboundField
var outboundMidiMap = make(map[string]outboundField)
var outboundChannelMap = make(map[string]uint8)

//type outboundFieldMap map[string]outboundField

func csSendEvent(field string, val uint8) (err error) {
	if !open {
		return
	}
	var channel uint8
	var found bool
	var fieldinfo outboundField

	channel, found = outboundChannelMap[field]
	if !found {
		log.Printf("WARN: Unknown channel for outbound event for control surface - %s\n", field)
		return
	}
	fieldinfo, found = outboundMidiMap[field]
	if !found {
		log.Printf("WARN: Unknown outbound event for control surface - %s\n", field)
		return
	}
	if fieldinfo.eventtype == Cc {
		v := uint8(int(val) / fieldinfo.scale)
		if strings.HasSuffix(field, "-tab") {
			v = uint8(fieldinfo.scale)
		}
		if err = sendCC(channel, fieldinfo.index, v); err != nil {
			return
		}
	} else if fieldinfo.eventtype == Poly {
		v := uint8(int(val) / fieldinfo.scale)
		if strings.HasSuffix(field, "-tab") {
			v = uint8(fieldinfo.scale)
		}
		if err = sendPolyAftertouch(channel, fieldinfo.index, v); err != nil {
			return
		}
	} else if fieldinfo.eventtype == Note {
		if val == 0 {
			if err = sendNoteOff(channel, fieldinfo.index, 0); err != nil {
				return
			}
		} else {
			if err = sendNoteOn(channel, fieldinfo.index, val); err != nil {
				return
			}
		}
	}
	return
}

func csHandleCC(channel uint8, cc uint8, val uint8) {
	var found bool
	var m inboundMidiMap
	var field inboundField
	m, found = inboundChannelMap[channel]
	if !found {
		log.Printf("ERROR no channel\n")
		return // we don't process this channel. Silently ignore
	}
	field, found = m.ccMap[cc]
	if !found {
		log.Printf("WARN: Unknown CC on control surface channel %d (CC%d)\n", channel+1, cc)
		return
	}
	_ = SendToUI(field.name, field.scale*int(val))
}

func csHandlePolyAftertouch(channel uint8, note uint8, val uint8) {
	var found bool
	var m inboundMidiMap
	var field inboundField
	m, found = inboundChannelMap[channel]
	if !found {
		return // we don't process this channel. Silently ignore
	}
	field, found = m.polyMap[note]
	if !found {
		log.Printf("WARN: Unknown Poly on control surface channel %d (CC%d)\n", channel+1, note)
		return
	}
	_ = SendToUI(field.name, field.scale*int(val))
}

func csHandleNoteOn(channel uint8, note uint8, velocity uint8) {
	var found bool
	var m inboundMidiMap
	var field inboundField
	m, found = inboundChannelMap[channel]
	if !found {
		return // we don't process this channel. Silently ignore
	}
	field, found = m.noteMap[note]
	if !found {
		log.Printf("WARN: Unknown Note event on control surface channel %d (note: %d)\n", channel+1, note)
		return
	}
	_ = SendToUI(field.name, int(velocity))
}

func csHandleNoteOff(channel uint8, note uint8, velocity uint8) {
	var found bool
	var m inboundMidiMap
	var field inboundField
	m, found = inboundChannelMap[channel]
	if !found {
		return // we don't process this channel. Silently ignore
	}
	field, found = m.noteMap[note]
	if !found {
		log.Printf("WARN: Unknown Note event on control surface channel %d (note: %d)\n", channel+1, note)
		return
	}
	_ = SendToUI(field.name, 0)
}
