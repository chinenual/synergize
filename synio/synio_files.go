package synio

import (
	"log"
	"github.com/chinenual/synergize/data"
	"github.com/pkg/errors"
)


func LoadVCE(slotnum byte, vce []byte) (err error) {
	if err = InitVRAM(); err != nil {
		err = errors.Wrap(err, "Failed to initialize Synergy VRAM")
		return
	}

	if err = command(OP_VCELOD, "VCELOD"); err != nil {
		return
	}

	if err = serialWriteByte(TIMEOUT_MS, slotnum, "write slotnum"); err != nil {
		return
	}
	var status byte
	if status, err = serialReadByte(TIMEOUT_MS, "read slotnum ACK"); err != nil {
		return
	}
	if status != DC1 {
		// slot error
		err = errors.Errorf("Invalid slotnum error")
		return
	}
	if err = serialWriteBytes(LONG_TIMEOUT_MS, vce, "write VCE"); err != nil {
		return
	}
	if status, err = serialReadByte(LONG_TIMEOUT_MS, "read VCE ACK"); err != nil {
		return
	}
	if status == ACK {
		// done - no filters
		// errors will implicitly show  up in the log but we need to explicitly log success
		if synioVerbose {
			log.Printf("VCELOD Success\n")
		}
		return
	}
	err = errors.Errorf("VCELOD incomplete. Can't handle filters upload yet")
	return
}

func LoadCRT(crt []byte) (err error) {
	if err = InitVRAM(); err != nil {
		return
	}

	if err = command(OP_VRLOD, "VRLOD"); err != nil {
		return
	}

	crcHash.Reset()

	var length = uint16(len(crt))
	if synioVerbose {
		log.Printf("length: %d (dec) %x (hex)\n", length, length)
	}

	lenHob, lenLob := data.WordToBytes(length)
	// LOB of the length
	calcCRCByte(lenLob)
	if err = serialWriteByte(TIMEOUT_MS, lenLob, "write length LOB"); err != nil {
		return
	}
	// HOB of the length
	calcCRCByte(lenHob)
	if err = serialWriteByte(TIMEOUT_MS, lenHob, "write length HOB"); err != nil {
		return
	}

	calcCRCBytes(crt)

	if err = serialWriteBytes(LONG_TIMEOUT_MS, crt, "write CRT bytes"); err != nil {
		return
	}

	crc := crcHash.CRC16()
	if synioVerbose {
		log.Printf("CRC: %d (dec) %x (hex) %x\n", crc, crc, crcHash.CRC())
	}

	crcHob, crcLob := data.WordToBytes(crc)
	// HOB of the crc
	if err = serialWriteByte(TIMEOUT_MS, crcHob, "write CRC HOB"); err != nil {
		return
	}
	// LOB of the crc
	if err = serialWriteByte(TIMEOUT_MS, crcLob, "write CRC LOB"); err != nil {
		return
	}

	var status byte
	if status, err = serialReadByte(LONG_TIMEOUT_MS, "read CRT ACK"); err != nil {
		return
	}
	if status == ACK {
		// errors will implicitly show  up in the log but we need to explicitly log success
		if synioVerbose {
			log.Printf("VRLOD Success\n")
		}
		return
	}
	err = errors.Errorf("Invalid CRC ACK from CRT upload")
	return
}


// Send Synergy "state" (STLOAD in the Z80 sources)
func LoadSYN(bytes []byte) (err error) {
	if err = command(OP_STLOAD, "STLOAD"); err != nil {
		return
	}
	// the SYN file actually has everything we need to send to the Synergy:
	// the initial byte count, the SEQ byte count and buffer and the final CRC.
	// Just send it as a block
	if err = serialWriteBytes(LONG_TIMEOUT_MS, bytes, "SYN bytes"); err != nil {
		return
	}
	// expect an ACK:
	var status byte
	if status, err = serialReadByte(LONG_TIMEOUT_MS, "read SYN ACK"); err != nil {
		return
	}
	if status == ACK {
		// errors will implicitly show  up in the log but we need to explicitly log success
		if synioVerbose {
			log.Printf("STLOD Success\n")
		}
		return
	}
	err = errors.Errorf("Invalid CRC ACK from SYN upload")
	return
}

// Retrieve Synergy "state" (STDUMP in the Z80 sources)
func SaveSYN() (bytes []byte, err error) {
	if err = command(OP_STDUMP, "STDUMP"); err != nil {
		return
	}

	var len_buf []byte
	if len_buf, err = serialReadBytes(TIMEOUT_MS, 2, "read CMOS length"); err != nil {
		return
	}

	cmos_len := data.BytesToWord(len_buf[1], len_buf[0])

	// read two CMOS data banks and the length of the sequencer (2 more bytes);

	if synioVerbose {
		log.Printf("CMOS LEN %d so read %d\n", cmos_len, cmos_len*2+2)
	}

	var cmos_buf []byte
	if cmos_buf, err = serialReadBytes(LONG_TIMEOUT_MS, cmos_len*2+2, "read CMOS"); err != nil {
		return
	}

	// decode sequencer length and possibly grab more
	seq_len := data.BytesToWord(cmos_buf[len(cmos_buf)-1], cmos_buf[len(cmos_buf)-2])
	if synioVerbose {
		log.Printf("SEQ LEN from synergy %d\n", seq_len)
	}

	// empty buf unless we have non-zero length to read
	seq_buf := []byte{}

	if seq_len != 0 {
		if seq_buf, err = serialReadBytes(LONG_TIMEOUT_MS, seq_len, "read SEQ"); err != nil {
			return
		}
	}
	var crc_buf []byte
	if crc_buf, err = serialReadBytes(TIMEOUT_MS, 2, "read CRC"); err != nil {
		return
	}

	// FIXME: these bytes seem out of order vs the length HOB/LOB yet seem to be transmitted the same from INTF.Z80 firmware sourcecode - I dont understand something..
	crcFromSynergy := data.BytesToWord(crc_buf[0], crc_buf[1])

	crcHash.Reset()

	calcCRCBytes(len_buf)
	calcCRCBytes(cmos_buf)
	calcCRCBytes(seq_buf)
	if synioVerbose {
		log.Printf("CRC from synergy %04x - our calculation %04x\n", crcFromSynergy, crcHash.CRC16())
	}

	if crcFromSynergy != crcHash.CRC16() {
		err = errors.Errorf("STDUMP CRC does not match got %04x, expected %04x",
			crcFromSynergy, crcHash.CRC16())
		return
	}
	// errors will implicitly show  up in the log but we need to explicitly log success
	if synioVerbose {
		log.Printf("STDUMP Success\n")
	}

	bytes = append(len_buf, cmos_buf...)
	bytes = append(bytes, seq_buf...)
	bytes = append(bytes, crc_buf...)
	return
}