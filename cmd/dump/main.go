package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/chinenual/synergize/synio"
)

func main() {
	var start, length uint64
	var b []byte
	var err error

	err = synio.SetSynergySerialPort("/dev/tty.usbserial-AL05OC8S", 9600, true, false, false)
	if err != nil {
		log.Panic(err)
	}

	if os.Args[1] == "--addrs" {
		if _, err = synio.EnableVoicingMode(); err != nil {
			log.Panic(err)
		}
		return
	}
	start, err = strconv.ParseUint(os.Args[1], 10, 16)
	if err != nil {
		log.Panic(err)
	}
	length, err = strconv.ParseUint(os.Args[2], 10, 16)
	if err != nil {
		log.Panic(err)
	}
	path := os.Args[3]

	b, err = synio.blockDump(uint16(start), uint16(length), "dump")
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		log.Panic(err)
	}
}
