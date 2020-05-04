package main

import (
	"log"
	"flag"
	"io/ioutil"
	"runtime"
	"path/filepath"

)

var (
	port = flag.String("port", getDefaultPort(), "the serial device")
	baud = flag.Uint("baud", getDefaultBaud(), "the serial baud rate")
)

func getDefaultBaud() uint {
	// FIXME: loads the prefs twice - harmless, but annoying
	prefsLoadPreferences()
	
	if prefsUserPreferences.SerialBaud != 0 {
		return prefsUserPreferences.SerialBaud
	}
	return 9600
}

func getDefaultPort() string {
	// FIXME: loads the prefs twice - harmless, but annoying
	prefsLoadPreferences()
	
	if prefsUserPreferences.SerialPort != "" {
		return prefsUserPreferences.SerialPort
	}
	if runtime.GOOS == "darwin" {
		files, _ := filepath.Glob("/dev/tty.usbserial*")
		for _,f := range(files) {
			return f
		}
	} else if runtime.GOOS == "linux" {
		files, _ := filepath.Glob("/dev/ttyUSB*")
		for _,f := range(files) {
			return f
		}
		
	} else {
		// windows
		return "COM1"
	}
	return ""
}

func diagCOMTST() {
	flag.Parse()

	log.Printf("%s at %d baud\n",
		prefsUserPreferences.SerialPort, prefsUserPreferences.SerialBaud);

	err := synioInit(prefsUserPreferences.SerialPort,
		prefsUserPreferences.SerialBaud, true, *serialVerboseFlag)
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

	log.Printf("%s at %d baud\n",
		prefsUserPreferences.SerialPort, prefsUserPreferences.SerialBaud);

	err := synioInit(prefsUserPreferences.SerialPort, prefsUserPreferences.SerialBaud, true, *serialVerboseFlag)
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

func diagInitAndPrintFirmwareID() (err error) {
	flag.Parse()
	
	log.Printf("%s at %d baud\n",
		prefsUserPreferences.SerialPort, prefsUserPreferences.SerialBaud);
	
	err = synioInit(prefsUserPreferences.SerialPort, prefsUserPreferences.SerialBaud, true, *serialVerboseFlag)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	var version [2]byte
	version,err = synioGetID()
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}
	log.Printf("Synergy is running firmware version %d.%d\n", version[0],version[1]);
	return
}

var slotnum byte = 1

func diagLoadVCE(path string) (err error) {
	var vce_bytes []byte
	vce_bytes,err = ioutil.ReadFile(path)
	if err != nil {
		return
	}

	log.Printf("VCE %s -- %d bytes into slotnum %d\n", path, len(vce_bytes), slotnum)

	err = synioLoadVCE(slotnum, vce_bytes)
	if err != nil {
		return
	}
	
	// load the next one in the next slot
	slotnum++
	return
	
}

func diagLoadCRT(path string) (err error) {
	var crt_bytes []byte
	crt_bytes,err = ioutil.ReadFile(path)
	if err != nil {
		return
	}

	log.Printf("CRT %s -- %d bytes \n", path, len(crt_bytes))

	err = synioLoadCRT(crt_bytes)
	return
}

func diagSaveSYN(path string) (err error) {
	var syn_bytes []byte
	syn_bytes, err = synioSaveSYN()
	
	if err != nil {
		return
	}

	log.Printf("SYN %s -- %d bytes \n", path, len(syn_bytes))

	err = ioutil.WriteFile(path, syn_bytes, 0644)
	return
}

func diagLoadSYN(path string) (err error) {	
	var syn_bytes []byte
	syn_bytes, err = ioutil.ReadFile(path)

	if err != nil {
		return
	}


	err = synioLoadSYN(syn_bytes)
	return
}
