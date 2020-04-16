package main

import (
	"log"
	"flag"
	"io/ioutil"
)

var (
	port = flag.String("port", "/dev/tty.usbserial-AL05OC8S", "the serial port")
)

func DiagCOMTST() {
	flag.Parse()

	log.Printf("%s\n", *port);

	err := SynioInit(*port)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}
	err = SynioDiagCOMTST()
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		log.Printf("Note:\n\tThe Synergy must be running in COMTST mode before executing this test.\n\tPress RESTORE + PROGRAM 4 on the Synergy then rerun this program.\n");
	} else {
		log.Printf("SUCCESS!\n")
	}
}

func DiagLOOPTST() {
	flag.Parse()

	log.Printf("%s\n", *port);

	err := SynioInit(*port)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	log.Printf("Entering LOOPBACK mode - any byte recieved from the Synergy is echo'd back\n")
	log.Printf("Start the test by pressing RESTORE + RESTORE + PROGRAM 1 on the Synergy\n")

	err = SynioDiagLOOPTST()
	if err != nil {
		log.Printf("ERROR: %s\n", err)
	}
}

func DiagPrintFirmwareID() {

	version,err := SynioGetID()
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}
	log.Printf("Synergy is running firmware version %d.%d\n", version[0],version[1]);
}

var slotnum byte = 1

func DiagLoadVCE(path string) {
	flag.Parse()

	log.Printf("%s\n", *port);
	err := SynioInit(*port)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	DiagPrintFirmwareID()

	var vce_bytes []byte
	vce_bytes,err = ioutil.ReadFile(path)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	log.Printf("VCE %s -- %d bytes into slotnum %d\n", path, len(vce_bytes), slotnum)

	err = SynioLoadVCE(slotnum, vce_bytes)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}
	
	// load the next one in the next slot
	slotnum++
	
	
}

func DiagLoadCRT(path string) {
	flag.Parse()

	log.Printf("%s\n", *port);
	err := SynioInit(*port)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	DiagPrintFirmwareID()

	var crt_bytes []byte
	crt_bytes,err = ioutil.ReadFile(path)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	log.Printf("CRT %s -- %d bytes \n", path, len(crt_bytes))

	err = SynioLoadCRT(crt_bytes)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}
}
