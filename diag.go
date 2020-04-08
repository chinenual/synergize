package main

import (
	"log"
	"io"
	"flag"
	"github.com/pkg/errors"
)

var (
	port = flag.String("port", "/dev/tty.usbserial-AL05OC8S", "the serial port")
)

const TIMEOUT_MS = 5000

func DiagMain() {
	flag.Parse()

	log.Printf("%s\n", *port);

	stream,err := SerialInit(*port);
	if err != nil {
		log.Fatal("Could not open serial port: ",err)
	}

	err = DiagCOMTST(stream);
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		log.Printf("Note:\n\tThe Synergy must be running in COMTST mode before executing this test.\n\tPress RESTORE + PROGRAM 4 on the Synergy then rerun this program.\n");
	} else {
		log.Printf("SUCCESS!\n")
	}
}

func DiagCOMTST(stream io.ReadWriteCloser) (err error) {	
	var i int
	for i = 0; i < 256; i++ {
		b := byte(i)
		log.Printf("%d ...\n", b)

		err = SerialWriteByte(stream, TIMEOUT_MS, b);
		if err != nil {
			return errors.Wrapf(err, "failed to write byte %d", b)
		}
		var read_b byte
		read_b,err = SerialReadByte(stream, TIMEOUT_MS);
		if err != nil {
			return errors.Wrapf(err, "failed to read byte %d", b)
		}
		if read_b != b {
			return errors.Errorf("read byte (%d) does not match what we sent (%d)", read_b, b)
		}
	}
	return nil
}
