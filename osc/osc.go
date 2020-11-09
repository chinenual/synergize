package osc

import (
	"fmt"
	"net"

	"github.com/pkg/errors"

	"github.com/chinenual/synergize/logger"
	"github.com/chinenual/synergize/zeroconf"
	goosc "github.com/hypebeast/go-osc/osc"
)

var verboseOscIn = false
var verboseOscOut = false

var client *goosc.Client
var server *goosc.Server
var listener net.PacketConn
var started = false

var csurfaceName string
var csurfaceAddress string
var csurfacePort uint

func SetControlSurface(name string, addr string, port uint) (err error) {
	{
		// Eek -- if csurfaceAddress is mis-configured -- eg.  "10.0.6" and missing a fourth quad, the net resolver code
		// gets confused and creates an address 10.0.0.6.  Huh?   Surely this is a bug in standard net the library.
		// But until that's fixed, we need to defend and give the user a useful error.
		// not quite as easy as it sounds since the address might be an IPV4 quad or might be a DNS name
		ip := net.ParseIP(addr)
		if ip == nil {
			// not a valid numeric quad; if can resolve the name continue, else it's probably a mis-typed numeric address
			if _, err = net.LookupAddr(addr); err != nil {
				return errors.Errorf("Invalid address for Control Surface: %s", addr)
			}
		}
	}
	csurfaceName = name
	csurfaceAddress = addr
	csurfacePort = port
	return
}

func ControlSurfaceName() string {
	if ControlSurfaceConfigured() {
		if csurfaceName == "" {
			return fmt.Sprintf("%s:%d", csurfaceAddress, csurfacePort)
		} else {
			return fmt.Sprintf("%s (%s:%d)", csurfaceName, csurfaceAddress, csurfacePort)
		}
	} else {
		return ""
	}
}

func ControlSurfaceConfigured() bool {
	return csurfaceAddress != ""
}

func Init(port uint, verboseIn bool, verboseOut bool, synergyName string) (err error) {
	verboseOscIn = verboseIn
	verboseOscOut = verboseOut

	started = false

	if ControlSurfaceConfigured() {
		client = goosc.NewClient(csurfaceAddress, int(csurfacePort))
	}

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	d := goosc.NewStandardDispatcher()
	if err = d.AddMsgHandler("*", func(msg *goosc.Message) {
		if verboseOscIn {
			logger.Infof("  OSC handle %v", msg)
		}
		if len(msg.Arguments) > 0 {
			// the touchosc can send empty args (eg for /ping) - filter them
			if err := OscHandleFromCSurface(msg.Address, msg.Arguments[0]); err != nil {
				logger.Errorf("Error handling OSC message: %v\n", err)
			}
		}
	}); err != nil {
		return
	}
	logger.Infof("OSC listen to %s...\n", addr)
	server = &goosc.Server{
		Addr:       addr,
		Dispatcher: d,
	}

	if err := zeroconf.StartServer(port, synergyName); err != nil {
		logger.Errorf("could not start zeroconf: %v\n", err)
	}

	go func() {
		if err := closeableListenAndServe(server); err != nil {
			logger.Errorf("could not start OSC server %v\n", err)
		}
	}()

	if err = csurfaceInit(); err != nil {
		return
	}

	return
}

func oscSendString(address string, arg string) (err error) {
	if client == nil {
		return
	}
	if verboseOscOut {
		logger.Infof("  OSC send %s %v", address, arg)
	}
	message := goosc.NewMessage(address, arg)
	if err = client.Send(message); err != nil {
		return
	}
	return
}
func oscSendInt(address string, arg int32) (err error) {
	if client == nil {
		return
	}
	if verboseOscOut {
		logger.Infof("  OSC send %s %v", address, arg)
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
	started = true

	return s.Serve(listener)
}

func Quit() (err error) {
	if started {
		if err = listener.Close(); err != nil {
			return
		}
		started = false
		client = nil
		server = nil
	}

	csurfaceName = ""
	csurfaceAddress = ""
	return
}
