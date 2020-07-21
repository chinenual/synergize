package main

import (
	"github.com/chinenual/synergize/synio"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func main() {
	var start, length uint64
	var b []byte
	var err error

	err = synio.Init("/dev/tty.usbserial-AL05OC8S", 9600, true, false, false)
	if err != nil {
		log.Panic(err)
	}

	if os.Args[1] == "--addrs" {
		synio.VoicingMode()
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

	b, err = synio.BlockDump(uint16(start), uint16(length))
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		log.Panic(err)
	}
}
