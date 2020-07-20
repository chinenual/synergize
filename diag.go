package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"

	"github.com/chinenual/synergize/data"
	"github.com/chinenual/synergize/synio"
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
		for _, f := range files {
			return f
		}
	} else if runtime.GOOS == "linux" {
		files, _ := filepath.Glob("/dev/ttyUSB*")
		for _, f := range files {
			return f
		}
		// if no USB serial, assume /dev/ttyS0
		return "/dev/ttyS0"

	} else {
		// windows
		return "COM1"
	}
	return ""
}

func diagCOMTST() {
	flag.Parse()

	log.Printf("%s at %d baud\n",
		prefsUserPreferences.SerialPort, prefsUserPreferences.SerialBaud)

	if err := synio.Init(prefsUserPreferences.SerialPort,
		prefsUserPreferences.SerialBaud, true, *serialVerboseFlag); err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}
	if err := synio.DiagCOMTST(); err != nil {
		log.Printf("ERROR: %s\n", err)
		log.Printf("Note:\n\tThe Synergy must be running in COMTST mode before executing this test.\n\tPress RESTORE + PROGRAM 4 on the Synergy then rerun this program.\n")
	} else {
		log.Printf("SUCCESS!\n")
	}
}

func diagLOOPTST() {
	flag.Parse()

	log.Printf("%s at %d baud\n",
		prefsUserPreferences.SerialPort, prefsUserPreferences.SerialBaud)

	if err := synio.Init(prefsUserPreferences.SerialPort, prefsUserPreferences.SerialBaud, true, *serialVerboseFlag); err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	fmt.Printf("\nEntering LOOPBACK mode - any byte recieved from the Synergy is echo'd back\n")
	fmt.Printf("Start the test by pressing RESTORE + RESTORE + PROGRAM 1 on the Synergy\n")

	if err := synio.DiagLOOPTST(); err != nil {
		log.Printf("ERROR: %s\n", err)
	}
}

func diagLINKTST() {
	flag.Parse()

	log.Printf("%s at %d baud\n",
		prefsUserPreferences.SerialPort, prefsUserPreferences.SerialBaud)

	if err := synio.Init(prefsUserPreferences.SerialPort, prefsUserPreferences.SerialBaud, true, *serialVerboseFlag); err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	fmt.Printf("\nEntering LINK TEST mode - any byte you type is sent to the Synergy and is \n")
	fmt.Printf("echo'd back.  LED's on the Synergy show the bytes recieved and various status\n")
	fmt.Printf("registers of the Synergy's serial connection.\n")
	fmt.Printf("\nStart the test by pressing RESTORE + PROGRAM 4 on the Synergy\n")

	if err := synio.DiagLINKTST(); err != nil {
		log.Printf("ERROR: %s\n", err)
	}
}

func diagInitAndPrintFirmwareID() (err error) {
	flag.Parse()

	log.Printf("%s at %d baud\n",
		prefsUserPreferences.SerialPort, prefsUserPreferences.SerialBaud)

	if err = synio.Init(prefsUserPreferences.SerialPort, prefsUserPreferences.SerialBaud, true, *serialVerboseFlag); err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	var version [2]byte
	if version, err = synio.GetID(); err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}
	log.Printf("Synergy is running firmware version %d.%d\n", version[0], version[1])
	return
}

var slotnum byte = 1

func diagLoadVCE(path string) (err error) {
	err = connectToSynergyIfNotConnected()
	if err != nil {
		return
	}
	var vce_bytes []byte
	if vce_bytes, err = ioutil.ReadFile(path); err != nil {
		return
	}

	log.Printf("VCE %s -- %d bytes into slotnum %d\n", path, len(vce_bytes), slotnum)

	if err = synio.LoadVCE(slotnum, vce_bytes); err != nil {
		return
	}

	// load the next one in the next slot
	slotnum++
	return

}

func diagLoadCRT(path string) (err error) {
	if err = connectToSynergyIfNotConnected(); err != nil {
		return
	}
	var crt_bytes []byte
	if crt_bytes, err = ioutil.ReadFile(path); err != nil {
		return
	}

	log.Printf("CRT %s -- %d bytes \n", path, len(crt_bytes))

	err = synio.LoadCRT(crt_bytes)
	return
}

func diagSaveSYN(path string) (err error) {
	if err = connectToSynergyIfNotConnected(); err != nil {
		return
	}
	var syn_bytes []byte

	if syn_bytes, err = synio.SaveSYN(); err != nil {
		return
	}

	log.Printf("SYN %s -- %d bytes \n", path, len(syn_bytes))

	err = ioutil.WriteFile(path, syn_bytes, 0644)
	return
}
func diagSaveVCE(path string) (err error) {
	if err = connectToSynergyIfNotConnected(); err != nil {
		return
	}
	var dumpedBytes []byte

	if dumpedBytes, err = synio.DumpVRAM(); err != nil {
		return
	}

	var readbuf = bytes.NewReader(dumpedBytes)
	var dumpedCrt data.CRT
	var dumpedVce data.VCE

	if dumpedCrt, err = data.ReadCrt(readbuf); err != nil {
		log.Printf("error parsing dumpedVRAM %v", err)
		return
	}
	dumpedVce = *dumpedCrt.Voices[0]
	log.Printf("VCE %s -- %d bytes: %s\n", path, len(dumpedBytes), data.VceToJson(dumpedVce))

	err = data.WriteVceFile(path, dumpedVce)
	return
}

func diagLoadSYN(path string) (err error) {
	if err = connectToSynergyIfNotConnected(); err != nil {
		return
	}
	var syn_bytes []byte

	if syn_bytes, err = ioutil.ReadFile(path); err != nil {
		return
	}

	err = synio.LoadSYN(syn_bytes)
	return
}
