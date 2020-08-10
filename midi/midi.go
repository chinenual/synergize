package midi

import (
	"log"

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

func QuitMidi() (err error) {
	log.Printf("Closing midi streams...\n")
	_ = in.StopListening()
	in.Close()
	out.Close()
	drv.Close()
	return
}

func InitMidi() (err error) {
	// you would take a real driver here e.g. rtmididrv.New()

	loadConfig()

	if drv, err = driver.New(); err != nil {
		return
	}

	// make sure to close all open ports at the end
	//	defer drv.Close()

	ins, err := drv.Ins()
	if err != nil {
		return
	}

	outs, err := drv.Outs()
	if err != nil {
		return
	}

	printInPorts(ins)
	printOutPorts(outs)

	in, out = ins[1], outs[1] // FIXME: hardcoded

	log.Printf("IN PORT:  [%v] %s\n", in.Number(), in.String())
	log.Printf("OUT PORT: [%v] %s\n", out.Number(), out.String())

	if err = in.Open(); err != nil {
		return
	}
	if err = out.Open(); err != nil {
		return
	}

	// the writer we are writing to
	wr = writer.New(out)

	// to disable logging, pass mid.NoLogger() as option
	rd := reader.New(
		//reader.NoLogger(),
		reader.ControlChange(handleCC),
		reader.NoteOff(handleNoteOff),
		reader.NoteOn(handleNoteOn),
		reader.PolyAftertouch(handlePolyAftertouch),
	)

	wr.SetChannel(15)

	log.Printf("listen to %s...\n", in)
	// listen for MIDI
	if err = rd.ListenTo(in); err != nil {
		log.Printf("ListenTo failed: %s\n", err)
		midiChannelQuit <- true
		return
	}

	// Output: got channel.NoteOn channel 0 key 60 velocity 100
	// got channel.NoteOff channel 0 key 60
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
