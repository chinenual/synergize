package zeroconf

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/grandcat/zeroconf"
)

// Connection/Zeroconf Lifecycle:
//  OSC server (re)started when:
//       voicing mode starts
//  OSC client (re)started when:
//       voicing mode starts
//
//  VST client started when:
//       first time IO requiring synergy connection
//          user explictly connects
//          load CRT for editing
//          load CRT
//          load SYN
//          save SYN
//          disable VRAM
//          run COMTest
//          toggle voicing mode
//
//  zeroconf browses when:
//       at program startup
//       user explicitly asks for a re-scan
//
//  zeroconf publishes OSC server address when:
//       at program startup
//       whenever server restarted

const numQueries = 5
const shortTimeout = time.Second * 3

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

func addIfNew(list *[]Service, entry *zeroconf.ServiceEntry) {
	for _, v := range *list {
		if v.InstanceName == entry.Instance {
			return
		}
	}
	s := newService(entry)
	*list = append(*list, s)
	return
}

var browsing = false

func Browse() {
	if browsing {
		return
	}
	defer func() {
		browsing = false
	}()
	browsing = true

	OscServices = make([]Service, 0)
	VstServices = make([]Service, 0)

	// HACK: some devices don't respond on the first query - we run several short queries and accumulate all the unique responses.
	for i := 0; i < numQueries; i++ {
		browse(shortTimeout)
	}

	log.Printf("ZEROCONF: end Browse OSC svcs: %#v   VST svcs: %#v\n", OscServices, VstServices)

}

func browse(timeout time.Duration) {
	// Discover services on the network
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Printf("ERROR: ZEROCONF: Failed to initialize resolver: %v\n", err.Error())
	}
	log.Printf("ZEROCONF: start browse...\n")

	oscEntries := make(chan *zeroconf.ServiceEntry)
	vstEntries := make(chan *zeroconf.ServiceEntry)

	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			//log.Printf("ZEROCONF: ... OSC service %s\n", entry.Instance)
			// ignore other OSC services - only those on TouchOSC might be of interest
			if strings.Contains(entry.Instance, "TouchOSC") {
				log.Printf("ZEROCONF: Found OSC service %s\n", entry.Instance)
				addIfNew(&OscServices, entry)
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
			addIfNew(&VstServices, entry)
		}
	}(vstEntries)

	ctx1, cancel1 := context.WithTimeout(context.Background(), timeout)
	defer cancel1()
	err = resolver.Browse(ctx1, "_osc._udp", "local.", oscEntries)
	if err != nil {
		log.Printf("ERROR: ZEROCONF: Failed to browse OSC: %v\n", err.Error())
	}
	ctx2, cancel2 := context.WithTimeout(context.Background(), timeout)
	defer cancel2()
	err = resolver.Browse(ctx2, "_synergy-vst._udp", "local.", vstEntries)
	if err != nil {
		log.Printf("ERROR: ZEROCONF: Failed to browse VST: %v\n", err.Error())
	}

	<-ctx1.Done()
	<-ctx2.Done()

}
