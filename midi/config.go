package midi

import (
	"github.com/pelletier/go-toml"
	"log"
	"reflect"
	"strconv"
)

func loadConfig() (err error) {
	var tree *toml.Tree
	if tree, err = toml.LoadFile("midi/midi.cfg"); err != nil {
		log.Printf("Failed to read config: %v\n", err)
		return
	}

	log.Printf("settings: %v\n", tree)

	channel := uint8(tree.Get("defaults.channel").(int64)) - 1 // config file uses 1-16, code uses 0-15

	var midiMap inboundMidiMap
	midiMap.ccMap = make(map[uint8]inboundField)
	midiMap.noteMap = make(map[uint8]inboundField)
	midiMap.polyMap = make(map[uint8]inboundField)
	inboundChannelMap[channel] = midiMap

	for _, k := range tree.Keys() {
		log.Printf("page " + k)
		if k == "defaults" {
			// ignore
		} else {
			// a page/tab:
			page := tree.Get(k).(*toml.Tree)
			for _, fieldname := range page.Keys() {
				if fieldname == "page-init" {
					// ignore for now...
				} else {
					// field definitions
					//log.Printf("  field " + fieldname)
					fielddef := page.Get(fieldname).(*toml.Tree)
					//log.Printf("     fielddef: %v\n", fielddef)
					scale := 1
					if fielddef.Has("scale") {
						scale = int(fielddef.Get("scale").(int64))
					}
					if fielddef.Has("cc") {
						val := fielddef.Get("cc")
						if reflect.ValueOf(val).Kind() == reflect.Int64 {
							v := uint8(val.(int64))
							midiMap.ccMap[v] = inboundField{name: fieldname, scale: scale}
							outboundMidiMap[fieldname] = outboundField{eventtype: Cc, index: v, scale: scale}
						} else {
							arr := val.([]interface{})
							for i, ele := range arr {
								v := uint8(ele.(int64))
								name := fieldname + "[" + strconv.Itoa(i+1) + "]"
								midiMap.ccMap[v] = inboundField{name: name, scale: scale}
								outboundMidiMap[name] = outboundField{eventtype: Cc, index: v, scale: scale}
							}
						}
					} else if fielddef.Has("note") {
						val := fielddef.Get("note")
						if reflect.ValueOf(val).Kind() == reflect.Int64 {
							v := uint8(val.(int64))
							midiMap.noteMap[v] = inboundField{name: fieldname, scale: scale}
							outboundMidiMap[fieldname] = outboundField{eventtype: Note, index: v, scale: scale}
						} else {
							arr := val.([]interface{})
							for i, ele := range arr {
								v := uint8(ele.(int64))
								name := fieldname + "[" + strconv.Itoa(i+1) + "]"
								midiMap.noteMap[v] = inboundField{name: name, scale: scale}
								outboundMidiMap[name] = outboundField{eventtype: Note, index: v, scale: scale}
							}
						}
					} else if fielddef.Has("poly") {
						val := fielddef.Get("poly")
						if reflect.ValueOf(val).Kind() == reflect.Int64 {
							v := uint8(val.(int64))
							midiMap.polyMap[v] = inboundField{name: fieldname, scale: scale}
							outboundMidiMap[fieldname] = outboundField{eventtype: Poly, index: v, scale: scale}
						} else {
							arr := val.([]interface{})
							for i, ele := range arr {
								v := uint8(ele.(int64))
								name := fieldname + "[" + strconv.Itoa(i+1) + "]"
								midiMap.polyMap[v] = inboundField{name: name, scale: scale}
								outboundMidiMap[name] = outboundField{eventtype: Poly, index: v, scale: scale}
							}
						}
					}
				}
			}
		}
	}

	log.Printf("in map: %v\n", inboundChannelMap)
	log.Printf("out map: %v\n", outboundMidiMap)

	return
}
