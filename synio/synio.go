package synio

import (
	"log"
	"strings"

	"github.com/chinenual/synergize/data"

	"github.com/pkg/errors"
	"github.com/snksoft/crc"
)

const RT_TIMEOUT_MS = 100     // "realtime" events need a shorter timeout
const LONG_TIMEOUT_MS = 20000 // after large amounts of IO, give the synergy more time to ack
const TIMEOUT_MS = 10000

const OP_KEYDWN = byte(0x01)
const OP_KEYUP = byte(0x02)
const OP_POT = byte(0x03)

const OP_VRLOD = byte(0x6b)
const OP_VRDUMP = byte(0x6c)
const OP_VCELOD = byte(0x6e)
const OP_DISABLEVRAM = byte(0x6f)
const OP_ENABLEVRAM = byte(0x70)
const OP_BLOCKLOAD = byte(0x71)

const OP_BLOCKDUMP = byte(0x72)
const OP_GETID = byte(0x74)
const OP_EXECUTE = byte(0x75)
const OP_IMODE = byte(0x76)
const OP_ASSIGNED_KEY = byte(0x77)
const OP_SELECT = byte(0x78)
const OP_STDUMP = byte(0x79)
const OP_STLOAD = byte(0x7a)
const OP_SLOW_BLOCKDUMP = byte(0x7c)

const ACK byte = 0x06
const DC1 byte = 0x11
const NAK byte = 0x15

var vramInitialized bool = false
var synioVerbose bool = false
var crcHash *crc.Hash

func Init(port string, baud uint, synVerbose bool, serialVerbose bool) (err error) {
	synioVerbose = synVerbose
	if err = serialInit(port, baud, serialVerbose); err != nil {
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

	return
}

func command(opcode byte, name string) (err error) {

	// check for pending input --
	//  silently read zero's
	//  if 1, 2 or 3 - treat it as a key or pot opcode (read 3 bytes incliuding the opcode)
	//  if >= 4 NAK it
	//  loop until no pending input
	// then send our opcode and loop until ACK'd

	if synioVerbose {
		log.Printf("synio: Send command opcode %02x - %s\n", opcode, name)
	}

	var status byte

	var retry = false //true;

	for retry {
		// use the short timeout for reads that may or may not have any data
		const SHORT_TIMEOUT_MS = 1000
		status, err = serialReadByte(SHORT_TIMEOUT_MS, "test for avail bytes")
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
					_, err = serialReadByte(TIMEOUT_MS, "read key/pot data")
					if err != nil {
						err = errors.Wrap(err, "error syncing command key/pot comm")
						return
					}
				}
			default:
				// otherwise, we need to send a NAK
				if err = serialWriteByte(TIMEOUT_MS, NAK, "write command NAK"); err != nil {
					return
				}
			}
		}
	}

	status = NAK
	var countdown = 10
	for status == NAK && countdown > 0 {
		countdown = countdown - 1
		// SYNHCS doesnt limit the number of retries, but it can lead to infinite loops/hangs.
		// We will only try N times
		err = serialWriteByte(TIMEOUT_MS, opcode, "write opcode")
		if err != nil {
			err = errors.Wrap(err, "error sending opcode")
			return
		}
		status, err = serialReadByte(TIMEOUT_MS, "read opcode ACK/NAK")
		if err != nil {
			err = errors.Wrap(err, "error reading opcode ACK/NAK")
			return
		}
	}
	if status != ACK {
		for {
			// TEMP: DRAIN
			status, err = serialReadByte(TIMEOUT_MS, "DRAIN")
			if err != nil {
				log.Println("error while draining", err)
				break
			} else if status == ACK {
				return
			}
			log.Printf("DRAIN: %x\n", status)
		}

		err = errors.Errorf("com error sending opcode %02x - did not get ACK/NAK, got %02x", opcode, status)

	}
	return
}

func LoadByte(addr uint16, value byte, purpose string) (err error) {
	var arr = []byte{value}
	if synioVerbose {
		log.Printf("synio: SetByte, addr: %04x, val: %02x (%s)", addr, value, purpose)
	}
	if err = BlockLoad(addr, arr, purpose); err != nil {
		return
	}
	return
}

func DumpByte(addr uint16, purpose string) (value byte, err error) {
	var arr []byte
	if arr, err = BlockDump(addr, 1, purpose); err != nil {
		return
	}
	value = arr[0]
	if synioVerbose {
		log.Printf("synio: GetByte, addr: %04x, val: %02x (%s)", addr, value, purpose)
	}
	return
}

func writeU16(v uint16, purpose string) (err error) {

	hob, lob := data.WordToBytes(v)

	if err = serialWriteByte(TIMEOUT_MS, hob, "write HOB "+purpose); err != nil {
		err = errors.Wrap(err, "error sending HOB "+purpose)
		return
	}
	if err = serialWriteByte(TIMEOUT_MS, lob, "write LOB "+purpose); err != nil {
		err = errors.Wrap(err, "error sending LOB "+purpose)
		return
	}
	return
}

func BlockDump(startAddress uint16, length uint16, purpose string) (bytes []byte, err error) {
	if err = command(OP_BLOCKDUMP, "OP_BLOCKDUMP"); err != nil {
		return
	}
	if err = writeU16(startAddress, "blockdump start address "+purpose); err != nil {
		return
	}
	if err = writeU16(length, "blockdump len "+purpose); err != nil {
		return
	}
	if bytes, err = serialReadBytes(LONG_TIMEOUT_MS, length, "block dump "+purpose); err != nil {
		return
	}
	return
}

func BlockLoad(startAddress uint16, bytes []byte, purpose string) (err error) {
	if err = command(OP_BLOCKLOAD, "OP_BLOCKLOAD"); err != nil {
		return
	}
	if err = writeU16(startAddress, "blockload start address "+purpose); err != nil {
		return
	}
	if err = writeU16(uint16(len(bytes)), "blockload len "+purpose); err != nil {
		return
	}
	if err = serialWriteBytes(LONG_TIMEOUT_MS, bytes, "block load "+purpose); err != nil {
		return
	}
	return
}

var synAddrs struct {
	SEQTAB uint16
	SEQCON uint16
	SEQVOI uint16
	CODE   uint16
	PTVAL  uint16
	PTSTAT uint16
	SOLOSC uint16
	// and the implied addresses relative to the above:
	EXTRA  uint16
	DEVICE uint16
	VALUE  uint16
	TRANSP uint16

	// Used in many SYNHCS address calculations:
	FILTAB uint16
	EDATA  uint16

	// Fixed addresses (from SYN-322.LNK)
	PROG uint16
	VRAM uint16
	VTAB uint16
	ROM  uint16
	RAM  uint16 // AKA DATA
	CMOS uint16

	exec_LDGENR uint16 // SUBROUTINE LDGENR - reload note generators
	exec_REAEQ  uint16 // SUBROUTINE REAEQ - alter amp scale for sounding notes
	exec_REFIL  uint16 // SUBROUTINE REFIL - recalc filter values
	exec_QUIET  uint16 // SUBROUTINE QUIET - stop all sounding notes
	exec_SETCON uint16 // SUBROUTINE SETCON - force immediate use of CMOS ram voice storage values
	exec_CKCMOS uint16 // SUBROUTINE CKCMOS - force reload note generators.
}

func getSynergyAddrs() (err error) {
	var b []byte
	if b, err = BlockDump(0x00c5, 14, "getSynergyAddrs"); err != nil {
		return
	}
	synAddrs.SEQTAB = data.BytesToWord(b[1], b[0])
	synAddrs.SEQCON = data.BytesToWord(b[3], b[2])
	synAddrs.SEQVOI = data.BytesToWord(b[5], b[4])
	synAddrs.CODE = data.BytesToWord(b[7], b[6])
	synAddrs.PTVAL = data.BytesToWord(b[9], b[8])
	synAddrs.PTSTAT = data.BytesToWord(b[11], b[10])
	synAddrs.SOLOSC = data.BytesToWord(b[13], b[12])
	// and the implied addresses relative to the above:
	synAddrs.EXTRA = synAddrs.SEQTAB - 128
	synAddrs.DEVICE = synAddrs.CODE + 2
	synAddrs.VALUE = synAddrs.CODE + 4
	synAddrs.TRANSP = synAddrs.SEQTAB - 5

	// Fixed addresses:
	synAddrs.PROG = 0x0000
	synAddrs.ROM = 0x5c72
	synAddrs.CMOS = 0xf000
	synAddrs.VRAM = 0x6000
	synAddrs.VTAB = 0x6033 // voice table ROM
	synAddrs.RAM = 0x8000  // aka DATA

	// Used in many SYNHCS address calculations:
	synAddrs.FILTAB = synAddrs.VTAB + 173
	synAddrs.EDATA = synAddrs.FILTAB + (16 * 32)

	synAddrs.exec_CKCMOS = uint16(0x00c2)
	synAddrs.exec_LDGENR = uint16(0x007a)
	synAddrs.exec_QUIET = uint16(0x009b)
	synAddrs.exec_REAEQ = uint16(0x00bc)
	synAddrs.exec_REFIL = uint16(0x00bf)
	synAddrs.exec_SETCON = uint16(0x00b9)

	if synioVerbose {
		log.Printf("synio: Synergy Addrs: %#v\n", synAddrs)
	}
	return
}

func InitVRAM() (err error) {
	if vramInitialized {
		return
	}
	if err = command(OP_ENABLEVRAM, "ENABLEVRAM"); err != nil {
		return
	}
	// errors will implicitly show  up in the log but we need to explicitly log success
	if synioVerbose {
		log.Printf("synio: ENABLEVRAM Success\n")
	}
	return
}

func DumpVRAM() (bytes []byte, err error) {
	if err = command(OP_VRDUMP, "VRDUMP"); err != nil {
		return
	}

	var len_buf []byte
	if len_buf, err = serialReadBytes(TIMEOUT_MS, 2, "read VRAM length"); err != nil {
		return
	}

	vram_len := data.BytesToWord(len_buf[1], len_buf[0])

	if synioVerbose {
		log.Printf("synio: DumpVRAM: len: %d bytes\n", vram_len)
	}

	if bytes, err = serialReadBytes(LONG_TIMEOUT_MS, vram_len, "read VRAM"); err != nil {
		return
	}

	var crc_buf []byte
	if crc_buf, err = serialReadBytes(TIMEOUT_MS, 2, "read CRC"); err != nil {
		return
	}

	crcFromSynergy := data.BytesToWord(crc_buf[0], crc_buf[1])

	crcHash.Reset()

	calcCRCBytes(len_buf)
	calcCRCBytes(bytes)
	//	calcCRCBytes(crc_buf)

	if synioVerbose {
		log.Printf("synio: CRC from synergy %04x - our calculation %04x\n", crcFromSynergy, crcHash.CRC16())
	}

	if crcFromSynergy != crcHash.CRC16() {
		err = errors.Errorf("VRDUMP CRC does not match got %04x, expected %04x",
			crcFromSynergy, crcHash.CRC16())
		return
	}
	// errors will implicitly show  up in the log but we need to explicitly log success
	if synioVerbose {
		log.Printf("synio: VRDUMP Success\n")
	}

	return

}

func GetID() (versionID [2]byte, err error) {
	if err = command(OP_GETID, "GETID"); err != nil {
		return
	}
	if versionID[0], err = serialReadByte(TIMEOUT_MS, "read id HB"); err != nil {
		return
	}
	if versionID[1], err = serialReadByte(TIMEOUT_MS, "read id LB"); err != nil {
		return
	}

	// errors will implicitly show  up in the log but we need to explicitly log success
	if synioVerbose {
		log.Printf("synio: GETID Success\n")
	}
	return
}

func DisableVRAM() (err error) {
	if err = command(OP_DISABLEVRAM, "DISABLEVRAM"); err == nil {
		// errors will implicitly show  up in the log but we need to explicitly log success
		if synioVerbose {
			log.Printf("synio: DISABLEVRAM Success\n")
		}
	}
	return
}

func calcCRCByte(b byte) {
	var arr []byte = make([]byte, 1)
	arr[0] = b
	calcCRCBytes(arr)
}

func calcCRCBytes(bytes []byte) {
	crcHash.Update(bytes)
}
