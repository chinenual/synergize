package main

import (
	"log"
	"flag"
	"io/ioutil"
)

var (
	port = flag.String("port", "/dev/tty.usbserial-AL05OC8S", "the serial port")
)

func diagCOMTST() {
	flag.Parse()

	log.Printf("%s\n", *port);

	err := synioInit(*port)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}
	err = synioDiagCOMTST()
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		log.Printf("Note:\n\tThe Synergy must be running in COMTST mode before executing this test.\n\tPress RESTORE + PROGRAM 4 on the Synergy then rerun this program.\n");
	} else {
		log.Printf("SUCCESS!\n")
	}
}

func diagLOOPTST() {
	flag.Parse()

	log.Printf("%s\n", *port);

	err := synioInit(*port)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	log.Printf("Entering LOOPBACK mode - any byte recieved from the Synergy is echo'd back\n")
	log.Printf("Start the test by pressing RESTORE + RESTORE + PROGRAM 1 on the Synergy\n")

	err = synioDiagLOOPTST()
	if err != nil {
		log.Printf("ERROR: %s\n", err)
	}
}

func diagPrintFirmwareID() {

	version,err := synioGetID()
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}
	log.Printf("Synergy is running firmware version %d.%d\n", version[0],version[1]);
}

var slotnum byte = 1

func diagLoadVCE(path string) {
	flag.Parse()

	log.Printf("%s\n", *port);
	err := synioInit(*port)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	diagPrintFirmwareID()

	var vce_bytes []byte
	vce_bytes,err = ioutil.ReadFile(path)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	log.Printf("VCE %s -- %d bytes into slotnum %d\n", path, len(vce_bytes), slotnum)

	err = synioLoadVCE(slotnum, vce_bytes)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}
	
	// load the next one in the next slot
	slotnum++
	
	
}

func diagLoadCRT(path string) {
	flag.Parse()

	log.Printf("%s\n", *port);
	err := synioInit(*port)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	diagPrintFirmwareID()

	var crt_bytes []byte
	crt_bytes,err = ioutil.ReadFile(path)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	log.Printf("CRT %s -- %d bytes \n", path, len(crt_bytes))

	err = synioLoadCRT(crt_bytes)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}
}

func diagSaveSYN(path string) {
	flag.Parse()

	log.Printf("%s\n", *port);
	err := synioInit(*port)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	diagPrintFirmwareID()

	var syn_bytes []byte
	syn_bytes, err = synioSaveSYN()
	
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	log.Printf("SYN %s -- %d bytes \n", path, len(syn_bytes))

	err = ioutil.WriteFile(path, syn_bytes, 0644)

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}
}
