package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/chinenual/synergize/data"
	"github.com/chinenual/synergize/logger"
	"github.com/chinenual/synergize/synio"
)

func diagCOMTST() {
	if err := synio.DiagCOMTST(); err != nil {
		logger.Errorf("%s\n", err)
		logger.Errorf("Note:\n\tThe Synergy must be running in COMTST mode before executing this test.\n\tPress RESTORE + PROGRAM 4 on the Synergy then rerun this program.\n")
	} else {
		logger.Infof("SUCCESS!\n")
	}
}

func diagLOOPTST() {
	fmt.Printf("\nEntering LOOPBACK mode - any byte received from the Synergy is echo'd back\n")
	fmt.Printf("Start the test by pressing RESTORE + RESTORE + PROGRAM 1 on the Synergy\n")

	if err := synio.DiagLOOPTST(); err != nil {
		logger.Errorf("%s\n", err)
	}
}

func diagLINKTST() {
	fmt.Printf("\nEntering LINK TEST mode - any byte you type is sent to the Synergy and is \n")
	fmt.Printf("echo'd back.  LED's on the Synergy show the bytes received and various status\n")
	fmt.Printf("registers of the Synergy's serial connection.\n")
	fmt.Printf("\nStart the test by pressing RESTORE + PROGRAM 4 on the Synergy\n")

	if err := synio.DiagLINKTST(); err != nil {
		logger.Errorf("%s\n", err)
	}
}

func diagInitAndPrintFirmwareID() (err error) {
	var version [2]byte
	if version, err = synio.GetID(); err != nil {
		logger.Errorf("%s\n", err)
		return
	}
	logger.Infof("Synergy is running firmware version %d.%d\n", version[0], version[1])
	return
}

var slotnum byte = 1

func diagLoadVCE(path string) (err error) {
	var vce_bytes []byte
	if vce_bytes, err = ioutil.ReadFile(path); err != nil {
		return
	}

	logger.Infof("VCE %s -- %d bytes into slotnum %d\n", path, len(vce_bytes), slotnum)

	if err = synio.LoadVCE(slotnum, vce_bytes); err != nil {
		return
	}

	// load the next one in the next slot
	slotnum++
	return

}

func diagLoadCRT(path string) (err error) {
	var crt_bytes []byte
	if crt_bytes, err = ioutil.ReadFile(path); err != nil {
		return
	}

	logger.Infof("CRT %s -- %d bytes \n", path, len(crt_bytes))

	err = synio.LoadCRTBytes(crt_bytes)
	return
}

func diagSaveSYN(path string) (err error) {
	var syn_bytes []byte

	if syn_bytes, err = synio.SaveSYN(); err != nil {
		return
	}

	logger.Infof("SYN %s -- %d bytes \n", path, len(syn_bytes))

	err = ioutil.WriteFile(path, syn_bytes, 0644)
	return
}
func diagSaveVCE(path string) (err error) {
	var dumpedBytes []byte

	if dumpedBytes, err = synio.DumpVRAM(); err != nil {
		return
	}

	var readbuf = bytes.NewReader(dumpedBytes)
	var dumpedCrt data.CRT
	var dumpedVce data.VCE

	if dumpedCrt, err = data.ReadCrt(readbuf); err != nil {
		logger.Errorf("error parsing dumpedVRAM %v", err)
		return
	}
	dumpedVce = *dumpedCrt.Voices[0]
	logger.Infof("VCE %s -- %d bytes: %s\n", path, len(dumpedBytes), data.VceToJson(dumpedVce))

	err = data.WriteVceFile(path, dumpedVce, false)
	return
}

func diagLoadSYN(path string) (err error) {
	var syn_bytes []byte

	if syn_bytes, err = ioutil.ReadFile(path); err != nil {
		return
	}

	err = synio.LoadSYN(syn_bytes)

	return
}
