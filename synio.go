package main

import (
	"io"
	"log"
	"strings"
	"github.com/pkg/errors"
	"github.com/snksoft/crc"
)

const VERBOSE = true
const TIMEOUT_MS = 5000

const OP_VRLOD      = byte(0x6b)
const OP_VCELOD     = byte(0x6e)
const OP_ENABLEVRAM = byte(0x70)
const OP_GETID      = byte(0x74)
const OP_STDUMP     = byte(0x79)

const ACK byte = 0x06
const DC1 byte = 0x11
const NAK byte = 0x15

var vramInitialized bool = false

var crcHash *crc.Hash

var (
	stream io.ReadWriteCloser
)

func SynioInit(port string) (err error) {
	stream,err = SerialInit(port, VERBOSE);
	if err != nil {
		return errors.Wrap(err, "Could not open serial port")
	}

	// From SYN-V322/CRCSET64.Z80:
	// ;       CYCLIC REDUNDANCY CHECK CHARACTER CHALCULATOR
	// ;       BASED ON X**16 + X**15 + X**2 +1 POLYNOMIAL
	//
	// which means "1100000000000101" (binary) or 0x8005.
	// In the Z80 code, I see left shifts which implies CRC16-BUYPASS rather than CRC16-ARC. 

	CRC16_BUYPASS := &crc.Parameters{Width: 16, Polynomial: 0x8005, Init: 0x0000, ReflectIn: false, ReflectOut: false, FinalXor: 0x0}
	
	crcHash = crc.NewHash(CRC16_BUYPASS)

//	crcHash = crc.NewHash(crc.CRC16)
	return 
}

func command(opcode byte, name string) (err error) {

	// check for pending input --
	//  silently read zero's
	//  if 1, 2 or 3 - treat it as a key or pot opcode (read 3 bytes incliuding the opcode)
	//  if >= 4 NAK it
	//  loop until no pending input
	// then send our opcode and loop until ACK'd

	if VERBOSE { log.Printf("send opcode %02x - %s\n", opcode, name); }
	
	var status byte
	
	// FIXME:
	// this drain the input loop doesnt work yet - can work on other opcode support
	// as long as the synergy hasnt queued up a bunch of output
//	var retry = true;
	var retry = false;
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
		// SYNHCS doesnt limit the number of retries, but it can lead to infinite loops/hangs.
		// We will only try N times
 log.Printf("aaa %v\n",err)
		err = SerialWriteByte(stream, TIMEOUT_MS, opcode, "write opcode")
		if err != nil {
			err = errors.Wrap(err, "error sending opcode")
			return 
		}
 log.Printf("bbb %v\n",err)
		status,err = SerialReadByte(stream, TIMEOUT_MS, "read opcode ACK/NAK")
 log.Printf("ccc %02x %v\n",status,err)
		if err != nil {
			err = errors.Wrap(err, "error reading opcode ACK/NAK")
 log.Printf("ddd %v\n",err)
			return 
		}
 log.Printf("eee %v\n",err)
	}
 log.Printf("fff %v\n",err)
	if status != ACK {
		err = errors.Errorf("com error sending opcode %02x - did not get ACK/NAK, got %02x",opcode,status)
	}
 log.Printf("ggg %v\n",err)
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

	crcHash.Reset()

	var length = uint16(len(crt))
	if verbose {log.Printf("length: %d (dec) %x (hex)\n", length, length)}

	lenHob,lenLob := wordToBytes(length)
	// LOB of the length
	calcCRCByte(lenLob)
	err = SerialWriteByte(stream, TIMEOUT_MS, lenLob, "write length LOB")
	if err != nil {
		err = errors.Wrap(err, "error sending length LOB")
		return 
	}
	// HOB of the length
	calcCRCByte(lenHob)
	err = SerialWriteByte(stream, TIMEOUT_MS, lenHob, "write length HOB")
	if err != nil {
		err = errors.Wrap(err, "error sending length HOB")
		return 
	}

	calcCRCBytes(crt)
	
	err = SerialWriteBytes(stream, TIMEOUT_MS, crt, "write CRT bytes")
	if err != nil {
		err = errors.Wrap(err, "error writing crt bytes ")
		return 
	}

	crc := crcHash.CRC16()
	if verbose {log.Printf("CRC: %d (dec) %x (hex) %x\n", crc, crc, crcHash.CRC())}
	
	crcHob,crcLob := wordToBytes(crc)
	// LOB of the crc
	err = SerialWriteByte(stream, TIMEOUT_MS, crcLob, "write CRC LOB")
	if err != nil {
		err = errors.Wrap(err, "error sending crc LOB")
		return 
	}
	// HOB of the crc
	err = SerialWriteByte(stream, TIMEOUT_MS, crcHob, "write CRC HOB")
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


// Retrieve Synergy "state" (STDUMP in the Z80 sources)
func SynioSaveSYN() (bytes []byte, err error) {
	err = command(OP_STDUMP, "STDUMP")
	if err != nil {
		err = errors.Wrap(err, "error sending opcode")
		return 
	}

	var len_buf []byte
	len_buf,err = SerialReadBytes(stream, TIMEOUT_MS, 2, "read CRC")
	if err != nil {
		err = errors.Wrap(err, "error reading CMOS length")
		return 
	}

	cmos_len := bytesToWord(len_buf[1],len_buf[0])
	
	// read two CMOS data banks and the length of the sequencer (2 more bytes);
	
	if verbose {log.Printf("CMOS LEN %d so read %d\n", cmos_len, cmos_len * 2 + 2)}
	
	var cmos_buf []byte
	cmos_buf,err = SerialReadBytes(stream, TIMEOUT_MS, cmos_len * 2 + 2, "read CMOS")
	if err != nil {
		err = errors.Wrap(err, "error reading CMOS length")
		return 
	}

	// FIXME: decode sequencer length and possibly grab more
	
	var crc_buf []byte
	crc_buf,err = SerialReadBytes(stream, TIMEOUT_MS, 2, "read CRC")
	if err != nil {
		err = errors.Wrap(err, "error reading CMOS length")
		return 
	}

	// FIXME: these bytes seem out of order vs the length HOB/LOB yet seem to be transmitted the same from INTF.Z80 firmware sourcecode - I dont understand something..
	crcFromSynergy := bytesToWord(crc_buf[0], crc_buf[1])

	crcHash.Reset();

	calcCRCBytes(len_buf)
	calcCRCBytes(cmos_buf)
//	calcCRCBytes(crc_buf)
	if verbose {log.Printf("CRC from synergy %x - our calculation %x\n", crcFromSynergy, crcHash.CRC16())}

	if crcFromSynergy != crcHash.CRC16() {
		err = errors.Errorf("STDUMP CRC does not match got %x, expected %x",
			crcFromSynergy, crcHash.CRC16())
		return
	}
	bytes = append(len_buf, cmos_buf...)
	bytes = append(bytes, crc_buf...)
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
		log.Printf("%d (%02x) ...\n", b, b)

		err = SerialWriteByte(stream, TIMEOUT_MS, b, "write test byte");
		if err != nil {
			return errors.Wrapf(err, "failed to write byte %02x", b)
		}
		var read_b byte
		read_b,err = SerialReadByte(stream, TIMEOUT_MS, "read test byte");
		if err != nil {
			return errors.Wrapf(err, "failed to read byte %02x", b)
		}
		if read_b != b {
			return errors.Errorf("read byte (%02x) does not match what we sent (%02x)", read_b, b)
		}
	}
	return nil
}

func SynioDiagLOOPTST() (err error) {	

	// infinite loop
	for true {

		var b byte
		b,err = SerialReadByte(stream, 1000 * 60 * 5, "read test byte");
		if err != nil {
			return errors.Wrapf(err, "failed to read byte %02x", b)
		}

		log.Printf("%d (%02x) ...\n", b, b)

		err = SerialWriteByte(stream, TIMEOUT_MS, b, "write test byte");
		if err != nil {
			return errors.Wrapf(err, "failed to write byte %02x", b)
		}
	}
	return nil
}

func bytesToWord(hob byte, lob byte) uint16 {
	return uint16(hob) << 8 + uint16(lob)
}

func wordToBytes(word uint16) (hob byte, lob byte) {
	hob = byte(word >> 8)
	lob = byte(word)
	return
}

func calcCRCByte(b byte)  {
	var arr []byte = make([]byte,1)
	arr[0] = b;
	calcCRCBytes(arr)
}

func calcCRCBytes(bytes []byte)  {
	crcHash.Update(bytes);
}
