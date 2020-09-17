package zeroconf

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/chinenual/synergize/logger"

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

var vstServiceType = "_synergia._tcp"

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
	s.InstanceName = strings.ReplaceAll(se.Name, "\\", "") // zeroconf escapes spaces and parens with \
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
	logger.Infof("ZEROCONF: Starting Zeroconf registration server... for service %s (%s) on port %d\n", serviceName, synergyName, oscListenPort)

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

var listenerRunning = false

func ListenerRunning() bool {
	return listenerRunning
}

func StartListener(vstServiceTypePrefix string) (err error) {
	vstServiceType = vstServiceTypePrefix

	// once we start listening we never stop or restart

	listenerRunning = true

	logger.Infof("ZEROCONF: Starting Zeroconf listener for service %s and %s\n", "_osc._udp", vstServiceType)

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

	// we begin with short lived queries since (on MacOS at least), the OS might block the initial responses
	// (until a user agrees to allow the application to connect to the network).  So we loop with 5s timeouts
	// allowing the listen to resend the query each time -- until we get a response for one of the listens.
	// Then we allow the listeners to run "forever"
	go func() {
		for {
			var timeout time.Duration
			if len(oscServices.m) == 0 && len(vstServices.m) == 0 {
				timeout = time.Second * 5
				logger.Infof("ZEROCONF: no results yet - sending queries\n")
			} else {
				logger.Infof("ZEROCONF: got first response - starting daemon\n")
			}
			var wg sync.WaitGroup
			wg.Add(2)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				if err = listenFor(timeout, &oscServices, "_osc._udp.local.", touchOscName); err != nil {
					return
				}
			}(&wg)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				if err = listenFor(timeout, &vstServices, vstServiceType+".local.", anyName); err != nil {
					return
				}
			}(&wg)
			wg.Wait()
		}
	}()
	return
}

func listenFor(timeout time.Duration, list *syncMap, serviceType string, validName func(string) bool) (err error) {
	var ctx context.Context
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(timeout))
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	addFn := func(srv dnssd.Service) {
		if validName(srv.Name) {
			logger.Infof("ZEROCONF: Add	%s	%s	%s\n", srv.Domain, srv.Type, srv.Name)
			add(list, &srv)
		} else {
			logger.Infof("ZEROCONF: IGNORING: Add	%s	%s	%s\n", srv.Domain, srv.Type, srv.Name)
		}
	}

	rmvFn := func(srv dnssd.Service) {
		if validName(srv.Name) {
			logger.Infof("ZEROCONF: Rmv	%s	%s	%s\n", srv.Domain, srv.Type, srv.Name)
			remove(list, &srv)
		} else {
			logger.Infof("ZEROCONF: IGNORING: Rmv	%s	%s	%s\n", srv.Domain, srv.Type, srv.Name)
		}
	}

	logger.Debugf("ZEROCONF: ListenFor %s\n", serviceType)
	if err = dnssd.LookupType(ctx, serviceType, addFn, rmvFn); err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			logger.Debugf("ZEROCONF: ListenFor %s %v\n", serviceType, err)
		} else {
			logger.Errorf("ZEROCONF: ListenFor %s %v\n", serviceType, err)
		}
		return
	}
	logger.Infof("ZEROCONF: ListenFor %s RETURNS\n", serviceType)
	return
}
