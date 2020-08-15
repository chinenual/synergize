package midi

import (
	"log"

	"github.com/pkg/errors"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/midi/writer"
	driver "gitlab.com/gomidi/portmididrv"
)

var midiChannelQuit = make(chan bool)
var drv midi.Driver
var in midi.In
var out midi.Out
var wr *writer.Writer
var rd *reader.Reader
var open = false

func QuitMidi() (err error) {
	if open {
		log.Printf("Closing midi streams...\n")
		_ = in.StopListening()
		in.Close()
		out.Close()
		drv.Close()
	}
	return
}

func InitMidi(midiInterface string, midiDeviceConfig string) (err error) {
	loadConfig(midiDeviceConfig)

	defer func(err *error) {
		if *err != nil {
			if quiterr := QuitMidi(); quiterr != nil {
				log.Printf("Errors from QuitMidi after init midi failed %v\n", quiterr)
			}
		}
	}(&err)

	if drv, err = driver.New(); err != nil {
		return
	}

	ins, err := drv.Ins()
	if err != nil {
		return
	}

	outs, err := drv.Outs()
	if err != nil {
		return
	}

	open = true

	printInPorts(ins)
	printOutPorts(outs)

	var found = false
	for _, port := range ins {
		if port.String() == midiInterface {
			in = port
			found = true
			break
		}
	}
	if !found {
		err = errors.Errorf("MIDI Interface %s not found as in inbound interface", midiInterface)
		return
	}
	found = false
	for _, port := range outs {
		if port.String() == midiInterface {
			out = port
			found = true
			break
		}
	}
	if !found {
		err = errors.Errorf("MIDI Interface %s not found as in outbound interface", midiInterface)
		return
	}

	log.Printf("IN PORT:  [%v] %s\n", in.Number(), in.String())
	log.Printf("OUT PORT: [%v] %s\n", out.Number(), out.String())

	if err = in.Open(); err != nil {
		return
	}
	if err = out.Open(); err != nil {
		return
	}

	wr = writer.New(out)

	rd = reader.New(
		//reader.NoLogger(),
		reader.ControlChange(handleCC),
		reader.NoteOff(handleNoteOff),
		reader.NoteOn(handleNoteOn),
		reader.PolyAftertouch(handlePolyAftertouch),
	)

	wr.SetChannel(15)
	return
}

func ListenMidi() (err error) {

	log.Printf("MIDI listen to %s...\n", in)
	// listen for MIDI
	if err = rd.ListenTo(in); err != nil {
		log.Printf("ListenTo failed: %s\n", err)
		midiChannelQuit <- true
		return
	}
	return
}

func WaitMidi() {
	// wait (forever) for the goroutine to exit
	<-midiChannelQuit
}

func printInPorts(ports []midi.In) {
	log.Printf("MIDI IN Ports\n")
	for _, port := range ports {
		log.Printf("    [%v] %s\n", port.Number(), port.String())
	}
}
func printOutPorts(ports []midi.Out) {
	log.Printf("MIDI OUT Ports\n")
	for _, port := range ports {
		log.Printf("    [%v] %s\n", port.Number(), port.String())
	}
}

func sendCC(channel, cc, val uint8) (err error) {
	log.Printf("Send CC: %d %d %d\n", channel, cc, val)
	err = writer.ControlChange(wr, cc, val)
	return
}

func sendPolyAftertouch(channel, note, velocity uint8) (err error) {
	log.Printf("Send Poly: %d %d %d\n", channel, note, velocity)
	err = writer.PolyAftertouch(wr, note, velocity)
	return
}

func sendNoteOn(channel, note, velocity uint8) (err error) {
	log.Printf("Send NoteOn: %d %d %d\n", channel, note, velocity)
	err = writer.NoteOn(wr, note, velocity)
	return
}

func sendNoteOff(channel, note, velocity uint8) (err error) {
	log.Printf("Send NoteOff: %d %d %d\n", channel, note, velocity)
	err = writer.NoteOffVelocity(wr, note, velocity)
	return
}

func handleCC(p *reader.Position, channel, cc, val uint8) {
	//	log.Printf("Handle CC: %d %d %d\n", channel, cc, val)
	csHandleCC(channel, cc, val)
}
func handlePolyAftertouch(p *reader.Position, channel, key, vel uint8) {
	//	log.Printf("Handle Poly: %d %d %d\n", channel, key, vel)
	csHandlePolyAftertouch(channel, key, vel)
}
func handleNoteOn(p *reader.Position, channel, key, vel uint8) {
	//	log.Printf("Handle NoteOn: %d %d %d\n", channel, key, vel)
	csHandleNoteOn(channel, key, vel)
}
func handleNoteOff(p *reader.Position, channel, key, vel uint8) {
	//	log.Printf("Handle NoteOff: %d %d %d\n", channel, key, vel)
	csHandleNoteOff(channel, key, vel)
}
