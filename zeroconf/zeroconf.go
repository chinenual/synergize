package zeroconf

import (
	"fmt"
	"log"

	"github.com/grandcat/zeroconf"
)

var server *zeroconf.Server

func StartServer(oscListenPort uint, synergyName string) (err error) {
	CloseServer()
	serviceName := "Synergize"
	log.Printf("Starting Zeroconf registration server... for service %s (%s) on port %d\n", serviceName, synergyName, oscListenPort)
	if server, err = zeroconf.Register(serviceName, "_osc._udp", "local.", int(oscListenPort), []string{fmt.Sprintf("Synergy=%s", synergyName)}, nil); err != nil {
		log.Printf("ERROR: Zeroconf registration failed: %v\n", err)
		return
	}
	return
}

func CloseServer() {
	if server != nil {
		log.Printf("Stopping Zeroconf registration server...\n")
		server.Shutdown()
	}
}
