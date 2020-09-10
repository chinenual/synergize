package zeroconf

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/brutella/dnssd"
)

type Service struct {
	InstanceName string
	HostName     string
	Port         uint
	Text         map[string]string
}

type syncMap struct {
	sync.RWMutex
	m map[string]Service
}

var oscServices syncMap
var vstServices syncMap

func GetOscServices() (result []Service) {
	oscServices.RLock()
	defer oscServices.RUnlock()

	for _, v := range oscServices.m {
		result = append(result, v)
	}
	return
}

func GetVstServices() (result []Service) {
	vstServices.RLock()
	defer vstServices.RUnlock()

	for _, v := range vstServices.m {
		result = append(result, v)
	}
	return
}

func newService(se *dnssd.Service) (s Service) {
	s.InstanceName = se.Name
	//s.Address = se.IPs[0].String()
	s.HostName = se.Host
	s.Port = uint(se.Port)
	s.Text = se.Text
	return
}

var responderCancel context.CancelFunc

func CloseServer() {
	if responderCancel != nil {
		responderCancel()
		responderCancel = nil
	}
}
func StartServer(oscListenPort uint, synergyName string) (err error) {
	CloseServer()
	serviceName := synergyName + " (Synergize)"
	serviceName = strings.ReplaceAll(serviceName, ".", ",")
	log.Printf("ZEROCONF: Starting Zeroconf registration server... for service %s (%s) on port %d\n", serviceName, synergyName, oscListenPort)

	cfg := dnssd.Config{
		Name: serviceName,
		Type: "_osc._udp",
		Port: int(oscListenPort),
	}
	var service dnssd.Service
	if service, err = dnssd.NewService(cfg); err != nil {
		return
	}
	var responder dnssd.Responder
	if responder, err = dnssd.NewResponder(); err != nil {
		return
	}
	//	var handle dnssd.ServiceHandle
	if _, err = responder.Add(service); err != nil {
		return
	}

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		responderCancel = cancel
		defer cancel()

		if err = responder.Respond(ctx); err != nil {
			return
		}
	}()
	return
}

func add(list *syncMap, entry *dnssd.Service) {
	list.Lock()
	defer list.Unlock()

	_, exists := list.m[entry.Name]
	if !exists {
		s := newService(entry)
		list.m[entry.Name] = s
	}
}

func remove(list *syncMap, entry *dnssd.Service) {
	list.Lock()
	defer list.Unlock()

	delete(list.m, entry.Name)
}

func StartListener() (err error) {
	oscServices.Lock()
	oscServices.m = make(map[string]Service)
	oscServices.Unlock()
	vstServices.Lock()
	vstServices.m = make(map[string]Service)
	vstServices.Unlock()

	touchOscName := func(name string) bool {
		return strings.Contains(name, "TouchOSC")
	}
	anyName := func(name string) bool {
		return true
	}

	go func() {
		if err = listenFor(&oscServices, "_osc._udp.local.", touchOscName); err != nil {
			log.Printf("ZEROCONF: ListenFor %s failed %v\n", "_osc._udp.local.", err)
		}
	}()
	go func() {
		if err = listenFor(&vstServices, "_synergia._tcp.local.", anyName); err != nil {
			log.Printf("ZEROCONF: ListenFor %s failed %v\n", "_synergia._tcp.local.", err)
		}
	}()
	return
}

func listenFor(list *syncMap, serviceType string, validName func(string) bool) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addFn := func(srv dnssd.Service) {
		if validName(srv.Name) {
			log.Printf("ZEROCONF: Add	%s	%s	%s\n", srv.Domain, srv.Type, srv.Name)
			add(list, &srv)
		} else {
			log.Printf("ZEROCONF: IGNORING: Add	%s	%s	%s\n", srv.Domain, srv.Type, srv.Name)
		}
	}

	rmvFn := func(srv dnssd.Service) {
		if validName(srv.Name) {
			log.Printf("ZEROCONF: Rmv	%s	%s	%s\n", srv.Domain, srv.Type, srv.Name)
			remove(list, &srv)
		} else {
			log.Printf("ZEROCONF: IGNORING: Rmv	%s	%s	%s\n", srv.Domain, srv.Type, srv.Name)
		}
	}

	log.Printf("ZEROCONF: ListenFor %s\n", serviceType)
	if err = dnssd.LookupType(ctx, serviceType, addFn, rmvFn); err != nil {
		log.Printf("ZEROCONF: ListenFor %s failed %v\n", serviceType, err)
		return
	}
	log.Printf("ZEROCONF: ListenFor %s RETURNS\n", serviceType)
	return
}
