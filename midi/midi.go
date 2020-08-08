package midi

import (
	"fmt"
	//	"time"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/midi/writer"
	driver "gitlab.com/gomidi/portmididrv"
)

func InitMidi() (err error) {
	// you would take a real driver here e.g. rtmididrv.New()
	var drv midi.Driver

	if drv, err = driver.New(); err != nil {
		return
	}

	// make sure to close all open ports at the end
	defer drv.Close()

	ins, err := drv.Ins()

	outs, err := drv.Outs()
	if err != nil {
		return
	}

	printInPorts(ins)
	printOutPorts(outs)

	in, out := ins[1], outs[1] // FIXME: hardcoded

	fmt.Printf("IN PORT:  [%v] %s\n", in.Number(), in.String())
	fmt.Printf("OUT PORT: [%v] %s\n", out.Number(), out.String())

	if err = in.Open(); err != nil {
		return
	}
	if err = out.Open(); err != nil {
		return
	}

	defer in.Close()
	defer out.Close()

	// the writer we are writing to
	wr := writer.New(out)

	// to disable logging, pass mid.NoLogger() as option
	rd := reader.New(
		//reader.NoLogger(),
		// write every message to the out port
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			fmt.Printf("got %s\n", msg)
		}),
	)

	wr.SetChannel(15)

	fmt.Printf("send CC1 100...\n")
	if err = writer.ControlChange(wr, 1, 100); err != nil {
		return
	}
	fmt.Printf("send CC2 50...\n")
	if err = writer.ControlChange(wr, 2, 50); err != nil {
		return
	}

	fmt.Printf("send NoteOn 0 0...\n")
	if err = writer.NoteOn(wr, 0, 1); err != nil {
		return
	}
	fmt.Printf("send NoteOn 2 0...\n")
	if err = writer.NoteOn(wr, 2, 1); err != nil {
		return
	}

	fmt.Printf("listen to...\n")
	// listen for MIDI
	if err = rd.ListenTo(in); err != nil {
		return
	}

	// Output: got channel.NoteOn channel 0 key 60 velocity 100
	// got channel.NoteOff channel 0 key 60
	return
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
