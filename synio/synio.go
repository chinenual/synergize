package synio

import (
	"log"
	"strings"
	"github.com/pkg/errors"
	"github.com/snksoft/crc"
)

const RT_TIMEOUT_MS = 100 // "realtime" events need a shorter timeout
const LONG_TIMEOUT_MS = 10000 // after large amounts of IO, give the synergy more time to ack
const TIMEOUT_MS = 5000

const OP_KEYDWN       = byte(0x01)
const OP_KEYUP        = byte(0x02)
const OP_POT          = byte(0x03)

const OP_VRLOD        = byte(0x6b)
const OP_VCELOD       = byte(0x6e)
const OP_DISABLEVRAM  = byte(0x6f)
const OP_ENABLEVRAM   = byte(0x70)
const OP_BLOCKLOAD    = byte(0x71)
const OP_BLOCKDUMP    = byte(0x72)
const OP_GETID        = byte(0x74)
const OP_EXECUTE      = byte(0x75)
const OP_IMODE        = byte(0x76)
const OP_ASSIGNED_KEY = byte(0x77)
const OP_SELECT       = byte(0x78)
const OP_STDUMP       = byte(0x79)
const OP_STLOAD       = byte(0x7a)
const OP_SLOW_BLOCKDUMP = byte(0x7c)


const ACK byte = 0x06
const DC1 byte = 0x11
const NAK byte = 0x15

var vramInitialized bool = false
var synioVerbose bool = false
var crcHash *crc.Hash

func Init(port string, baud uint, synVerbose bool, serialVerbose bool) (err error) {
	synioVerbose = synVerbose
	err = serialInit(port, baud, serialVerbose);
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

	return 
}

func command(opcode byte, name string) (err error) {

	// check for pending input --
	//  silently read zero's
	//  if 1, 2 or 3 - treat it as a key or pot opcode (read 3 bytes incliuding the opcode)
	//  if >= 4 NAK it
	//  loop until no pending input
	// then send our opcode and loop until ACK'd

	if synioVerbose { log.Printf("send command opcode %02x - %s\n", opcode, name); }
	
	var status byte
	
	var retry = false;//true;

	for retry {
		// use the short timeout for reads that may or may not have any data
		const SHORT_TIMEOUT_MS = 1000
		status,err = serialReadByte(SHORT_TIMEOUT_MS, "test for avail bytes")
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
					_,err = serialReadByte(TIMEOUT_MS, "read key/pot data")
					if err != nil {
						err = errors.Wrap(err, "error syncing command key/pot comm")
					}
				}
			default:
				// otherwise, we need to send a NAK
				err = serialWriteByte(TIMEOUT_MS, NAK, "write NAK")
				if err != nil {
					err = errors.Wrap(err, "error sending NAK during command sync")
					return 
				}
			}
		}
	}
	
	status = NAK
	var countdown = 10
	for status == NAK && countdown > 0 {
		countdown = countdown-1
		// SYNHCS doesnt limit the number of retries, but it can lead to infinite loops/hangs.
		// We will only try N times
		err = serialWriteByte(TIMEOUT_MS, opcode, "write opcode")
		if err != nil {
			err = errors.Wrap(err, "error sending opcode")
			return 
		}
		status,err = serialReadByte(TIMEOUT_MS, "read opcode ACK/NAK")
		if err != nil {
			err = errors.Wrap(err, "error reading opcode ACK/NAK")
			return 
		}
	}
	if status != ACK {
		for {
			// TEMP: DRAIN
			status,err = serialReadByte(TIMEOUT_MS, "DRAIN")
			if err != nil {
				log.Println("error while draining",err)
				break;
			} else if status == ACK {
				return
			}
			log.Printf("DRAIN: %x\n",status)
		}


		err = errors.Errorf("com error sending opcode %02x - did not get ACK/NAK, got %02x",opcode,status)

		
	}
	return
}

func writeU16(v uint16, purpose string) (err error) {
	
	hob,lob := wordToBytes(v)

	err = serialWriteByte(TIMEOUT_MS, hob, "write HOB " + purpose)
	if err != nil {
		err = errors.Wrap(err, "error sending HOB " + purpose)
		return 
	}
	err = serialWriteByte(TIMEOUT_MS, lob, "write LOB " + purpose)
	if err != nil {
		err = errors.Wrap(err, "error sending LOB " + purpose)
		return 
	}
	return
}

func BlockDump(startAddress uint16, length uint16) (bytes []byte, err error) {
	command(OP_BLOCKDUMP, "OP_BLOCKDUMP")
	err = writeU16(startAddress, "blockdump start address")
	err = writeU16(startAddress, "blockdump len")
	if err != nil {
		return 
	}
	bytes,err = serialReadBytes(LONG_TIMEOUT_MS, length, "block dump" )
	if err != nil {
		return 
	}
	return
}

func BlockLoad(startAddress uint16, length uint16, bytes []byte) (err error) {
	command(OP_BLOCKLOAD, "OP_BLOCKLOAD")
	err = writeU16(startAddress, "blockload start address")
	err = writeU16(startAddress, "blockload len")
	if err != nil {
		return 
	}
	err = serialWriteBytes(LONG_TIMEOUT_MS, bytes, "block load" )
	if err != nil {
		return 
	}
	return
}

func InitVRAM() (err error) {
	if vramInitialized {
		return
	}
	err = command(OP_ENABLEVRAM, "ENABLEVRAM")
	if err != nil {
		err = errors.Wrap(err, "error sending ENABLEVRAM opcode")
		return 
	}
	// errors will implicitly show  up in the log but we need to explicitly log success
	if synioVerbose { log.Printf("ENABLEVRAM Success\n"); }	
	return
}

func LoadVCE(slotnum byte, vce []byte) (err error) {
	err = InitVRAM()
	if err != nil {
		err = errors.Wrap(err, "Failed to initialize Synergy VRAM")
		return 
	}
	
	err = command(OP_VCELOD, "VCELOD")
	if err != nil {
		err = errors.Wrap(err, "error sending VCELOD opcode")
		return 
	}

	err = serialWriteByte(TIMEOUT_MS, slotnum, "write slotnum")
	if err != nil {
		err = errors.Wrap(err, "error sending slotnum")
		return 
	}
	var status byte
	status,err = serialReadByte(TIMEOUT_MS, "read slotnum ACK")
	if err != nil {
		err = errors.Wrap(err, "error reading slotnum ack")
		return 
	}
	if status != DC1 {
		// slot error
		err = errors.Errorf("Invalid slotnum error")
		return
	}
	err = serialWriteBytes(TIMEOUT_MS, vce, "write VCE")
	if err != nil {
		err = errors.Wrap(err, "error writing vce ")
		return 
	}
	status,err = serialReadByte(LONG_TIMEOUT_MS, "read VCE ACK")
	if err != nil {
		err = errors.Wrap(err, "error reading vce ack")
		return 
	}
	if status == ACK {
		// done - no filters
		// errors will implicitly show  up in the log but we need to explicitly log success
		if synioVerbose { log.Printf("VCELOD Success\n"); }	
		return
	}
	err = errors.Errorf("VCELOD incomplete. Can't handle filters upload yet")
	return
}

func LoadCRT(crt []byte) (err error) {
	err = InitVRAM()
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
	if synioVerbose {log.Printf("length: %d (dec) %x (hex)\n", length, length)}

	lenHob,lenLob := wordToBytes(length)
	// LOB of the length
	calcCRCByte(lenLob)
	err = serialWriteByte(TIMEOUT_MS, lenLob, "write length LOB")
	if err != nil {
		err = errors.Wrap(err, "error sending length LOB")
		return 
	}
	// HOB of the length
	calcCRCByte(lenHob)
	err = serialWriteByte(TIMEOUT_MS, lenHob, "write length HOB")
	if err != nil {
		err = errors.Wrap(err, "error sending length HOB")
		return 
	}

	calcCRCBytes(crt)
	
	err = serialWriteBytes(TIMEOUT_MS, crt, "write CRT bytes")
	if err != nil {
		err = errors.Wrap(err, "error writing crt bytes ")
		return 
	}

	crc := crcHash.CRC16()
	if synioVerbose {log.Printf("CRC: %d (dec) %x (hex) %x\n", crc, crc, crcHash.CRC())}
	
	crcHob,crcLob := wordToBytes(crc)
	// HOB of the crc
	err = serialWriteByte(TIMEOUT_MS, crcHob, "write CRC HOB")
	if err != nil {
		err = errors.Wrap(err, "error sending crc HOB")
		return 
	}
	// LOB of the crc
	err = serialWriteByte(TIMEOUT_MS, crcLob, "write CRC LOB")
	if err != nil {
		err = errors.Wrap(err, "error sending crc LOB")
		return 
	}

	var status byte
	status,err = serialReadByte(LONG_TIMEOUT_MS, "read CRT ACK")
	if err != nil {
		err = errors.Wrap(err, "error reading crt ack")
		return 
	}
	if status == ACK {
		// errors will implicitly show  up in the log but we need to explicitly log success
		if synioVerbose { log.Printf("VRLOD Success\n"); }	
		return
	}
	err = errors.Errorf("Invalid CRC ACK from CRT upload")
	return
}


// Send Synergy "state" (STLOAD in the Z80 sources)
func LoadSYN(bytes []byte) (err error) {
	err = command(OP_STLOAD, "STLOAD")
	if err != nil {
		err = errors.Wrap(err, "error sending opcode")
		return 
	}
	// the SYN file actually has everything we need to send to the Synergy:
	// the initial byte count, the SEQ byte count and buffer and the final CRC.
	// Just send it as a block 
	err = serialWriteBytes(TIMEOUT_MS, bytes, "SYN bytes")
	if err != nil {
		err = errors.Wrap(err, "error sending SYN bytes")
		return 
	}
	// expect an ACK:
	var status byte
	status,err = serialReadByte(LONG_TIMEOUT_MS, "read SYN ACK")
	if err != nil {
		err = errors.Wrap(err, "error reading SYN ack")
		return 
	}
	if status == ACK {
		// errors will implicitly show  up in the log but we need to explicitly log success
		if synioVerbose { log.Printf("STLOD Success\n"); }	
		return
	}
	err = errors.Errorf("Invalid CRC ACK from SYN upload")
	return
}


// Retrieve Synergy "state" (STDUMP in the Z80 sources)
func SaveSYN() (bytes []byte, err error) {
	err = command(OP_STDUMP, "STDUMP")
	if err != nil {
		err = errors.Wrap(err, "error sending opcode")
		return 
	}

	var len_buf []byte
	len_buf,err = serialReadBytes(TIMEOUT_MS, 2, "read CRC")
	if err != nil {
		err = errors.Wrap(err, "error reading CMOS length")
		return 
	}

	cmos_len := bytesToWord(len_buf[1],len_buf[0])
	
	// read two CMOS data banks and the length of the sequencer (2 more bytes);
	
	if synioVerbose {log.Printf("CMOS LEN %d so read %d\n", cmos_len, cmos_len * 2 + 2)}
	
	var cmos_buf []byte
	cmos_buf,err = serialReadBytes(TIMEOUT_MS, cmos_len * 2 + 2, "read CMOS")
	if err != nil {
		err = errors.Wrap(err, "error reading CMOS length")
		return 
	}

	// decode sequencer length and possibly grab more
	seq_len := bytesToWord(cmos_buf[len(cmos_buf)-1], cmos_buf[len(cmos_buf)-2])
	if synioVerbose {log.Printf("SEQ LEN from synergy %d\n", seq_len)}

	// empty buf unless we have non-zero length to read
	seq_buf := []byte{}
	
	if seq_len != 0 {
		seq_buf,err = serialReadBytes(TIMEOUT_MS, seq_len, "read SEQ")
		if err != nil {
			err = errors.Wrap(err, "error reading SEQ")
			return 
		} 
	}
	var crc_buf []byte
	crc_buf,err = serialReadBytes(TIMEOUT_MS, 2, "read CRC")
	if err != nil {
		err = errors.Wrap(err, "error reading CMOS length")
		return 
	}

	// FIXME: these bytes seem out of order vs the length HOB/LOB yet seem to be transmitted the same from INTF.Z80 firmware sourcecode - I dont understand something..
	crcFromSynergy := bytesToWord(crc_buf[0], crc_buf[1])

	crcHash.Reset();

	calcCRCBytes(len_buf)
	calcCRCBytes(cmos_buf)
	calcCRCBytes(seq_buf)
	if synioVerbose {log.Printf("CRC from synergy %x - our calculation %x\n", crcFromSynergy, crcHash.CRC16())}

	if crcFromSynergy != crcHash.CRC16() {
		err = errors.Errorf("STDUMP CRC does not match got %x, expected %x",
			crcFromSynergy, crcHash.CRC16())
		return
	}
	// errors will implicitly show  up in the log but we need to explicitly log success
	if synioVerbose { log.Printf("STDUMP Success\n"); }	

	bytes = append(len_buf, cmos_buf...)
	bytes = append(bytes, seq_buf...)
	bytes = append(bytes, crc_buf...)
	return
}

func GetID() (versionID [2]byte, err error) {
	err = command(OP_GETID, "GETID")
	if err != nil {
		err = errors.Wrap(err, "error sending opcode")
		return 
	}
	versionID[0],err = serialReadByte(TIMEOUT_MS, "read id HB")
	if err != nil {
		err = errors.Wrap(err, "error reading HB")
		return 
	}
	versionID[1],err = serialReadByte(TIMEOUT_MS, "read id LB")
	if err != nil {
		err = errors.Wrap(err, "error reading LB")
		return 
	}
	
	// errors will implicitly show  up in the log but we need to explicitly log success
	if synioVerbose { log.Printf("GETID Success\n"); }	
	return 
}

func DisableVRAM() (err error) {
	err = command(OP_DISABLEVRAM, "DISABLEVRAM")
	if err == nil {
		// errors will implicitly show  up in the log but we need to explicitly log success
		if synioVerbose { log.Printf("DISABLEVRAM Success\n"); }	
	}
	return
}

func DiagCOMTST() (err error) {	

	var i int
	for i = 0; i < 256; i++ {
		b := byte(i)
		log.Printf("%d (%02x) ...\n", b, b)

		err = serialWriteByte(TIMEOUT_MS, b, "write test byte");
		if err != nil {
			return errors.Wrapf(err, "failed to write byte %02x", b)
		}
		var read_b byte
		read_b,err = serialReadByte(TIMEOUT_MS, "read test byte");
		if err != nil {
			return errors.Wrapf(err, "failed to read byte %02x", b)
		}
		if read_b != b {
			return errors.Errorf("read byte (%02x) does not match what we sent (%02x)", read_b, b)
		}
	}
	// errors will implicitly show  up in the log but we need to explicitly log success
	if synioVerbose { log.Printf("COMTST Success\n"); }	
	return nil
}

func DiagLOOPTST() (err error) {	

	if synioVerbose { log.Printf("WARNING: LOOPTST causes Synergize to enter an infinte loop supporting the Synergy based test.  You must explicitly kill the Synergize process to stop the test.\n"); }
	for true {

		var b byte
		b,err = serialReadByte(1000 * 60 * 5, "read test byte");
		if err != nil {
			return errors.Wrapf(err, "failed to read byte %02x", b)
		}

		log.Printf("%d (%02x) ...\n", b, b)

		err = serialWriteByte(TIMEOUT_MS, b, "write test byte");
		if err != nil {
			return errors.Wrapf(err, "failed to write byte %02x", b)
		}
	}
	return nil
}

func SelectVoiceMapping(v1, v2, v3 ,v4 byte) (err error) {
	if err = command(OP_SELECT, "OP_SELECT"); err != nil {
		return errors.Wrapf(err, "failed to OP_SELECT")
	}
	if err = serialWriteByte(RT_TIMEOUT_MS, v1, "voice1"); err != nil {
		return errors.Wrapf(err, "failed to voice1 mapping")
	}
	if err = serialWriteByte(RT_TIMEOUT_MS, v2, "voice2"); err != nil {
		return errors.Wrapf(err, "failed to voice2 mapping")
	}
	if err = serialWriteByte(RT_TIMEOUT_MS, v3, "voice3"); err != nil {
		return errors.Wrapf(err, "failed to voice3 mapping")
	}
	if err = serialWriteByte(RT_TIMEOUT_MS, v4, "voice4"); err != nil {
		return errors.Wrapf(err, "failed to voice4 mapping")
	}
	return
}

// voice      1..4
// key        0..73
// velocity   0..32
func KeyDown(voice, key, velocity byte) (err error) {
	if err = command(OP_ASSIGNED_KEY, "OP_ASSIGNED_KEY"); err != nil {
		return errors.Wrapf(err, "failed to OP_ASSIGNED_KEY")
	}
//	if err = serialWriteByte(RT_TIMEOUT_MS, OP_KEYDWN, "OP_KEYDWN"); err != nil {
//		return errors.Wrapf(err, "failed to OP_KEYDWN")
//	}
	if err = serialWriteByte(RT_TIMEOUT_MS, voice, "voice"); err != nil {
		return errors.Wrapf(err, "failed to send notedown voice")
	}
	if err = serialWriteByte(RT_TIMEOUT_MS, key, "key"); err != nil {
		return errors.Wrapf(err, "failed to send notedown key")
	}
	if err = serialWriteByte(RT_TIMEOUT_MS, velocity, "velocity"); err != nil {
		return errors.Wrapf(err, "failed to send notedown velocity")
	}
	return
}


// Synergy can't turn off voice-specific key - we're in rolling voice assign mode
// key        0..73
// velocity   0..32
func KeyUp(key, velocity byte) (err error) {
	if err = command(OP_KEYUP, "OP_KEYUP"); err != nil {
		return errors.Wrapf(err, "failed to OP_KEYUP")
	}
	if err = serialWriteByte(RT_TIMEOUT_MS, key, "key"); err != nil {
		return errors.Wrapf(err, "failed to send noteup key")
	}
	if err = serialWriteByte(RT_TIMEOUT_MS, velocity, "velocity"); err != nil {
		return errors.Wrapf(err, "failed to send noteup velocity")
	}
	return
}

func Pedal(up bool) (err error) {
	const OPERAND_PEDAL_SUSTAIN = byte(64)
	const OPERAND_PEDAL_LATCH = byte(65)

	if err = command(OP_POT, "OP_POT"); err != nil {
		return errors.Wrapf(err, "failed to send pedal OP")
	}
	if err = serialWriteByte(RT_TIMEOUT_MS, OPERAND_PEDAL_SUSTAIN, "OPERAND_PEDAL_SUSTAIN"); err != nil {
		return errors.Wrapf(err, "failed to send pedal SUSTAIN operand")
	}
	var value = byte(0) // down
	if up {
		value = 127
	}
	if err = serialWriteByte(RT_TIMEOUT_MS, value, "pedal value"); err != nil {
		return errors.Wrapf(err, "failed to send pedal value")
	}
	return
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
