package main

import (
	"io"
	"log"
	"strings"
	"github.com/pkg/errors"
)

const VERBOSE = true
const TIMEOUT_MS = 5000

const OP_GETID = byte(0x74)
const OP_VCELOD = byte(0x6e)
const OP_VRLOD = byte(0x68)
const OP_ENABLEVRAM = byte(0x70)

const ACK byte = 0x06
const DC1 byte = 0x11
const NAK byte = 0x15

var vramInitialized bool = false

var (
	stream io.ReadWriteCloser
)

func SynioInit(port string) (err error) {
	stream,err = SerialInit(port, VERBOSE);
	if err != nil {
		return errors.Wrap(err, "Could not open serial port")
	}
	return 
}

func command(opcode byte, name string) (err error) {

	// check for pending input --
	//  silently read zero's
	//  if 1, 2 or 3 - treat it as a key or pot opcode (read 3 bytes incliuding the opcode)
	//  if >= 4 NAK it
	//  loop until no pending input
	// then send our opcode and loop until ACK'd

	if VERBOSE { log.Printf("send opcode %2x - %s\n", opcode, name); }
	
	var status byte
	var retry = false;//true;
	for retry {
		// use the short timeout for reads that may or may not have any data
		const SHORT_TIMEOUT_MS = 5000
		status,err = SerialReadByte(stream, SHORT_TIMEOUT_MS, "test for avail bytes")
		if err != nil && (!strings.Contains(err.Error(), "TIMEOUT:")) {
			err = errors.Wrap(err, "error syncing command comm")
			return 
		}
		if err != nil {
			// it timed out -- exit the loop
			err = nil
			retry = false
			
		} else {
			// if it didnt timeout, process the command:
			
			switch status {
			case 0:
				// ignore
			case 1, 2, 3:
				// KEY OR POT msg; consume 2 more bytes
				for i := 0; i < 3; i++ {
					_,err = SerialReadByte(stream, TIMEOUT_MS, "read key/pot data")
					if err != nil {
						err = errors.Wrap(err, "error syncing command key/pot comm")
					}
				}
			default:
				// otherwise, we need to send a NAK
				err = SerialWriteByte(stream, TIMEOUT_MS, NAK, "write NAK")
				if err != nil {
					err = errors.Wrap(err, "error sending NAK during command sync")
					return 
				}
			}
		}
	}
	
	status = NAK
	var countdown = 3
	for status == NAK && countdown > 0 {
		countdown = countdown-1
		// FIXME: this seems to be what SYNHCS does, but it can lead to infinite loops/hangs.  We will only try N times
 log.Printf("aaa %x\n",err)
		err = SerialWriteByte(stream, TIMEOUT_MS, opcode, "write opcode")
		if err != nil {
			err = errors.Wrap(err, "error sending opcode")
			return 
		}
 log.Printf("bbb %x\n",err)
		status,err = SerialReadByte(stream, TIMEOUT_MS, "read opcode ACK/NAK")
 log.Printf("ccc %x\n",err)
		if err != nil {
			err = errors.Wrap(err, "error reading opcode ACK/NAK")
 log.Printf("ddd %x\n",err)
			return 
		}
 log.Printf("eee %x\n",err)
	}
 log.Printf("fff %x\n",err)
	if status != ACK {
		err = errors.Errorf("com error sending opcode %2x - did not get ACK/NAK, got %2x",opcode,status)
	}
 log.Printf("ggg %x\n",err)
	return
}

func SynioInitVRAM() (err error) {
	if vramInitialized {
		return
	}
	err = command(OP_ENABLEVRAM, "ENABLEVRAM")
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
	
	err = command(OP_VCELOD, "VCELOD")
	if err != nil {
		err = errors.Wrap(err, "error sending VCELOD opcode")
		return 
	}

	err = SerialWriteByte(stream, TIMEOUT_MS, slotnum, "write slotnum")
	if err != nil {
		err = errors.Wrap(err, "error sending slotnum")
		return 
	}
	var status byte
	status,err = SerialReadByte(stream, TIMEOUT_MS, "read slotnum ACK")
	if err != nil {
		err = errors.Wrap(err, "error reading slotnum ack")
		return 
	}
	if status != DC1 {
		// slot error
		err = errors.Errorf("Invalid slotnum error")
		return
	}
	err = SerialWriteBytes(stream, TIMEOUT_MS, vce, "write VCE")
	if err != nil {
		err = errors.Wrap(err, "error writing vce ")
		return 
	}
	status,err = SerialReadByte(stream, TIMEOUT_MS, "read VCE ACK")
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

func SynioLoadCRT(crt []byte) (err error) {
	err = SynioInitVRAM()
	if err != nil {
		err = errors.Wrap(err, "Failed to initialize Synergy VRAM")
		return 
	}
	
	err = command(OP_VRLOD, "VRLOD")
	if err != nil {
		err = errors.Wrap(err, "error sending VRLOD opcode")
		return 
	}

	var crc uint16 = 0
	var length = uint16(len(crt))
	
	// LOB of the length
	calcCRC(&crc, byte(0xff & length))
	err = SerialWriteByte(stream, TIMEOUT_MS, byte(0xff & length), "write length LOB")
	if err != nil {
		err = errors.Wrap(err, "error sending length LOB")
		return 
	}
	// HOB of the length
	calcCRC(&crc, byte(0xff & (length>>8)))
	err = SerialWriteByte(stream, TIMEOUT_MS, byte(0xff & (length>>8)), "write length HOB")
	if err != nil {
		err = errors.Wrap(err, "error sending length HOB")
		return 
	}
	for _,b := range(crt) {
		calcCRC(&crc, b)
	}
	
	err = SerialWriteBytes(stream, TIMEOUT_MS, crt, "write CRT bytes")
	if err != nil {
		err = errors.Wrap(err, "error writing crt bytes ")
		return 
	}

	// LOB of the crc
	err = SerialWriteByte(stream, TIMEOUT_MS, byte(0xff & crc), "write CRC LOB")
	if err != nil {
		err = errors.Wrap(err, "error sending crc LOB")
		return 
	}
	// HOB of the crc
	err = SerialWriteByte(stream, TIMEOUT_MS, byte(0xff & (crc>>8)), "write CRC HOB")
	if err != nil {
		err = errors.Wrap(err, "error sending crc HOB")
		return 
	}

	var status byte
	status,err = SerialReadByte(stream, TIMEOUT_MS, "read CRT ACK")
	if err != nil {
		err = errors.Wrap(err, "error reading crt ack")
		return 
	}
	if status == ACK {
		return
	}
	err = errors.Errorf("Invalid CRC ACK from CRT upload")
	return
}


func SynioGetID() (versionID [2]byte, err error) {
	err = command(OP_GETID, "GETID")
	if err != nil {
		err = errors.Wrap(err, "error sending opcode")
		return 
	}
	versionID[0],err = SerialReadByte(stream, TIMEOUT_MS, "read id HB")
	if err != nil {
		err = errors.Wrap(err, "error reading HB")
		return 
	}
	versionID[1],err = SerialReadByte(stream, TIMEOUT_MS, "read id LB")
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
		log.Printf("%d (%2x) ...\n", b, b)

		err = SerialWriteByte(stream, TIMEOUT_MS, b, "write test byte");
		if err != nil {
			return errors.Wrapf(err, "failed to write byte %2x", b)
		}
		var read_b byte
		read_b,err = SerialReadByte(stream, TIMEOUT_MS, "read test byte");
		if err != nil {
			return errors.Wrapf(err, "failed to read byte %2x", b)
		}
		if read_b != b {
			return errors.Errorf("read byte (%2x) does not match what we sent (%2x)", read_b, b)
		}
	}
	return nil
}

func calcCRC(crc *uint16, val byte)  {
	// TBD
	*crc += uint16(val)
}
