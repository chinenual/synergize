package main

import (
	"io"
	"log"
	"github.com/pkg/errors"
)

const TIMEOUT_MS = 5000

const OP_GETID = byte(0x74)
const OP_VCELOD = byte(0x6e)
const OP_ENABLEVRAM = byte(0x70)

const ACK byte = 0x06
const DC1 byte = 0x11
const NAK byte = 0x15

var vramInitialized bool = false

var (
	stream io.ReadWriteCloser
)

func SynioInit(port string) (err error) {
	stream,err = SerialInit(port);
	if err != nil {
		return errors.Wrap(err, "Could not open serial port")
	}
	return 
}

func command(opcode byte) (err error) {

	status := NAK
	for status == NAK {
		// FIXME: this seemss to be what SYNHCS does, but it can lead to infinite loops/hangs. 
		err = SerialWriteByte(stream, TIMEOUT_MS, opcode)
		if err != nil {
			err = errors.Wrap(err, "error sending opcode")
			return 
		}
		status,err = SerialReadByte(stream, TIMEOUT_MS)
		if err != nil {
			err = errors.Wrap(err, "error reading ACK/NAK")
			return 
		}
	}
	if status != ACK {
		err = errors.Errorf("com error sending opcode %x - did not get ACK/NAK, got %x",opcode,status)
	}
	return
}

func SynioInitVRAM() (err error) {
	if vramInitialized {
		return
	}
	err = command(OP_ENABLEVRAM)
	if err != nil {
		err = errors.Wrap(err, "error sending ENABLEVRAM opcode")
		return 
	}
	return
}

func SynioLoadVCE(slotnum byte, vce []byte) (err error) {
	err = SynioInitVRAM()
	if err != nil {
		err = errors.Wrap(err, "Failed to initialize Synergy VRAM")
		return 
	}
	
	err = command(OP_VCELOD)
	if err != nil {
		err = errors.Wrap(err, "error sending VCELOD opcode")
		return 
	}
	err = SerialWriteByte(stream, TIMEOUT_MS, slotnum)
	if err != nil {
		err = errors.Wrap(err, "error sending slotnum")
		return 
	}
	var status byte
	status,err = SerialReadByte(stream, TIMEOUT_MS)
	if err != nil {
		err = errors.Wrap(err, "error reading slotnum ack")
		return 
	}
	if status != DC1 {
		// slot error
		err = errors.Errorf("Invalid slotnum error")
		return
	}
	err = SerialWriteBytes(stream, TIMEOUT_MS, vce)
	if err != nil {
		err = errors.Wrap(err, "error writing vce ")
		return 
	}
	status,err = SerialReadByte(stream, TIMEOUT_MS)
	if err != nil {
		err = errors.Wrap(err, "error reading vce ack")
		return 
	}
	if status == ACK {
		// done - no filters
		return
	}
	err = errors.Errorf("Cant handle filters upload yet")
	return
}


func SynioGetID() (versionID [2]byte, err error) {
	err = command(OP_GETID)
	if err != nil {
		err = errors.Wrap(err, "error sending opcode")
		return 
	}
	versionID[0],err = SerialReadByte(stream, TIMEOUT_MS)
	if err != nil {
		err = errors.Wrap(err, "error reading HB")
		return 
	}
	versionID[1],err = SerialReadByte(stream, TIMEOUT_MS)
	if err != nil {
		err = errors.Wrap(err, "error reading LB")
		return 
	}
	return 
}

func SynioDiagCOMTST() (err error) {	

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
