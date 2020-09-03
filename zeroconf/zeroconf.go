package zeroconf

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/grandcat/zeroconf"
)

const browseTimeout = time.Second * 15

var server *zeroconf.Server

type Service struct {
	InstanceName string
	Address      string
	HostName     string
	Port         uint
	Text         []string
}

var OscServices []Service
var VstServices []Service

func newService(se *zeroconf.ServiceEntry) (s Service) {
	s.InstanceName = se.Instance
	s.Address = se.AddrIPv4[0].String()
	s.HostName = se.HostName
	s.Port = uint(se.Port)
	s.Text = se.Text
	return
}

func StartServer(oscListenPort uint, synergyName string) (err error) {
	CloseServer()
	serviceName := synergyName + " (Synergize)"
	serviceName = strings.ReplaceAll(serviceName, ".", ",")
	log.Printf("ZEROCONF: Starting Zeroconf registration server... for service %s (%s) on port %d\n", serviceName, synergyName, oscListenPort)
	if server, err = zeroconf.Register(serviceName, "_osc._udp", "local.", int(oscListenPort), []string{fmt.Sprintf("Synergy=%s", synergyName)}, nil); err != nil {
		log.Printf("ERROR: ZEROCONF: Zeroconf registration failed: %v\n", err)
		return
	}
	return
}

func CloseServer() {
	if server != nil {
		log.Printf("ZEROCONF: Stopping Zeroconf registration server...\n")
		server.Shutdown()
	}
}

var browsing = false

func Browse() {
	if browsing {
		return
	}
	browsing = true

	// Discover services on the network
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Printf("ERROR: ZEROCONF: Failed to initialize resolver: %v\n", err.Error())
	}
	log.Printf("ZEROCONF: start browse...\n")

	OscServices = make([]Service, 0)
	VstServices = make([]Service, 0)

	oscEntries := make(chan *zeroconf.ServiceEntry)
	vstEntries := make(chan *zeroconf.ServiceEntry)

	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			//log.Printf("ZEROCONF: ... OSC service %s\n", entry.Instance)
			// ignore other OSC services - only those on TouchOSC might be of interest
			if strings.Contains(entry.Instance, "TouchOSC") {
				log.Printf("ZEROCONF: Found OSC service %s\n", entry.Instance)
				var s = newService(entry)
				OscServices = append(OscServices, s)
				log.Printf("ZEROCONF: add OSC svcs: %#v\n", OscServices)
			} else if strings.Contains(entry.Instance, "Synergize") {
				// silently ignore
			} else {
				log.Printf("ZEROCONF: Ignoring OSC service %s\n", entry.Instance)
			}
		}
	}(oscEntries)

	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			log.Printf("ZEROCONF: Found VST service %s\n", entry.Instance)
			var s = newService(entry)
			VstServices = append(VstServices, s)
		}
	}(vstEntries)

	ctx1, cancel1 := context.WithTimeout(context.Background(), browseTimeout)
	defer cancel1()
	err = resolver.Browse(ctx1, "_osc._udp", "local.", oscEntries)
	if err != nil {
		log.Printf("ERROR: ZEROCONF: Failed to browse OSC: %v\n", err.Error())
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), browseTimeout)
	defer cancel2()
	err = resolver.Browse(ctx2, "_synergyvst._udp", "local.", vstEntries)
	if err != nil {
		log.Printf("ERROR: ZEROCONF: Failed to browse VST: %v\n", err.Error())
	}

	<-ctx1.Done()
	<-ctx2.Done()

	log.Printf("ZEROCONF: end Browse OSC svcs: %#v   VST svcs: %#v\n", OscServices, VstServices)
	browsing = false
}