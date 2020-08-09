package midi

import (
	"fmt"
	"strconv"

	"github.com/asticode/go-astilectron"
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
	in.Close()
	out.Close()
	drv.Close()
	return
}

func InitMidi(w *astilectron.Window) (err error) {
	// you would take a real driver here e.g. rtmididrv.New()
	initBridge(w)

	if drv, err = driver.New(); err != nil {
		return
	}

	// make sure to close all open ports at the end
	//	defer drv.Close()

	ins, err := drv.Ins()

	outs, err := drv.Outs()
	if err != nil {
		return
	}

	printInPorts(ins)
	printOutPorts(outs)

	in, out = ins[1], outs[1] // FIXME: hardcoded

	fmt.Printf("IN PORT:  [%v] %s\n", in.Number(), in.String())
	fmt.Printf("OUT PORT: [%v] %s\n", out.Number(), out.String())

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
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			fmt.Printf("got %s\n", msg)
		}),

		reader.ControlChange(handleCC),
		reader.NoteOff(handleNoteOff),
		reader.NoteOn(handleNoteOn),
	)

	wr.SetChannel(15)

	fmt.Printf("send CC1 20...\n")
	if err = writer.ControlChange(wr, 1, 20); err != nil {
		return
	}
	fmt.Printf("send CC2 42...\n")
	if err = writer.ControlChange(wr, 2, 42); err != nil {
		return
	}

	fmt.Printf("send NoteOn 1 127...\n")
	if err = writer.NoteOn(wr, 1, 127); err != nil {
		return
	}
	fmt.Printf("send NoteOn 3 127...\n")
	if err = writer.NoteOn(wr, 3, 127); err != nil {
		return
	}
	fmt.Printf("send NoteOn 33 127...\n")
	if err = writer.NoteOn(wr, 33, 127); err != nil {
		return
	}
	fmt.Printf("send NoteOn 35 127...\n")
	if err = writer.NoteOn(wr, 35, 127); err != nil {
		return
	}
	fmt.Printf("send NoteOn 82 1...\n")
	if err = writer.NoteOn(wr, 82, 1); err != nil {
		return
	}
	fmt.Printf("send NoteOn 83 2...\n")
	if err = writer.NoteOn(wr, 83, 2); err != nil {
		return
	}
	fmt.Printf("send NoteOn 84 3...\n")
	if err = writer.NoteOn(wr, 84, 3); err != nil {
		return
	}

	//	go func(in midi.In) {
	//defer in.Close()
	//defer out.Close()

	fmt.Printf("listen to %s...\n", in)
	// listen for MIDI
	if err = rd.ListenTo(in); err != nil {
		fmt.Printf("ListenTo failed: %s\n", err)
		midiChannelQuit <- true
		return
	}
	//	}(in)

	// Output: got channel.NoteOn channel 0 key 60 velocity 100
	// got channel.NoteOff channel 0 key 60
	return
}

func WaitMidi() {
	// wait (forever) for the goroutine to exit
	<-midiChannelQuit
}

func printInPorts(ports []midi.In) {
	fmt.Printf("MIDI IN Ports\n")
	for _, port := range ports {
		fmt.Printf("[%v] %s\n", port.Number(), port.String())
	}
	fmt.Printf("\n\n")

}
func printOutPorts(ports []midi.Out) {
	fmt.Printf("MIDI OUT Ports\n")
	for _, port := range ports {
		fmt.Printf("[%v] %s\n", port.Number(), port.String())
	}
	fmt.Printf("\n\n")

}

type UIMsg struct {
	Field string
	Value int
}

func handleCC(p *reader.Position, channel, cc, val uint8) {
	fmt.Printf("Handle CC: %d %d %d\n", channel, cc, val)
	if channel == 15 && cc >= 1 && cc <= 16 {
		v := -11 + int(val)
		SendToUI("updateFromMIDI", UIMsg{"OHARM[" + strconv.Itoa(int(cc)) + "]", v})
	}
}
func handleNoteOn(p *reader.Position, channel, key, vel uint8) {
	fmt.Printf("Handle NoteOn: %d %d %d\n", channel, key, vel)
}
func handleNoteOff(p *reader.Position, channel, key, vel uint8) {
	fmt.Printf("Handle NoteOff: %d %d %d\n", channel, key, vel)
}

func sendCC(cc uint8, val uint8) {
	fmt.Printf("  send MIDI CC %d %d\n", cc, val)
	writer.ControlChange(wr, cc, val)
}
