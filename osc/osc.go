package osc

import (
	"fmt"
	"log"
	"net"

	goosc "github.com/hypebeast/go-osc/osc"
)

var verboseOscIn = false
var verboseOscOut = false

var client *goosc.Client
var server *goosc.Server
var listener net.PacketConn

func OscInit(port uint, csurfaceAddress string, csurfacePort uint, verboseIn bool, verboseOut bool) (err error) {
	verboseOscIn = verboseIn
	verboseOscOut = verboseOut

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	d := goosc.NewStandardDispatcher()
	if err = d.AddMsgHandler("*", func(msg *goosc.Message) {
		goosc.PrintMessage(msg)
	}); err != nil {
		return
	}
	log.Printf("OSC listen to %s...\n", addr)
	server = &goosc.Server{
		Addr:       addr,
		Dispatcher: d,
	}

	client = goosc.NewClient(csurfaceAddress, int(csurfacePort))

	go func() {
		if err := closeableListenAndServe(server); err != nil {
			log.Printf("ERROR: could not start OSC server %v\n", err)
		}
	}()

	return
}

func oscSendString(address string, arg string) (err error) {
	if verboseOscOut {
		log.Printf("  OSC send %s %v", address, arg)
	}
	message := goosc.NewMessage(address, arg)
	if err = client.Send(message); err != nil {
		return
	}
	return
}
func oscSendInt(address string, arg int32) (err error) {
	if verboseOscOut {
		log.Printf("  OSC send %s %v", address, arg)
	}
	message := goosc.NewMessage(address, arg)
	if err = client.Send(message); err != nil {
		return
	}
	return
}

// copied from osc.ListenAndServe, but exposing the listener socket so it can be closed prematurely
func closeableListenAndServe(s *goosc.Server) (err error) {
	if s.Dispatcher == nil {
		s.Dispatcher = goosc.NewStandardDispatcher()
	}

	listener, err = net.ListenPacket("udp", s.Addr)
	if err != nil {
		return err
	}

	return s.Serve(listener)
}

func OscQuit() (err error) {
	if err = listener.Close(); err != nil {
		return
	}
	return
}
