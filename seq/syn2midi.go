package seq

import (
	"io/ioutil"

	"github.com/chinenual/synergize/logger"

	"github.com/chinenual/synergize/data"
)

// See INTF.Z80 TSTATE   and SEQREQ.Z80
//
// SYN File contains:
//  8) 2-bytes "CMOS length"
// 	1) internal CMOS ram data
//;	2) external CMOS ram data
//;	3) 2-bytes Length of seq. data (will send 0 if no sequence and stop)
//;	4) PTVAL+55, +56, +57, +58
//;	5) seq control data (TRANSP - VOIUSE)
//;	6) seq. data table
//; 7) CRC

func ConvertSYNToMIDI(path string) (err error) {
	var syn_bytes []byte

	if syn_bytes, err = ioutil.ReadFile(path); err != nil {
		return
	}

	cmos_len := data.BytesToWord(syn_bytes[1], syn_bytes[0])
	logger.Debugf("CMOS LEN: %v\n", cmos_len)

	seqtab_start := cmos_len*2 + 4
	logger.Debugf("SEQTAB start: 0x%x\n", seqtab_start)

	// seqlen is immediately after the cmos bytes - +2 for the header (cmos_len) bytes:
	seq_len := data.BytesToWord(syn_bytes[seqtab_start-1], syn_bytes[seqtab_start-2])
	logger.Debugf("SEQTAB LEN: %v\n", seq_len)

	if seq_len > 0 {
		// last two bytes of the file are the CRC
		if err = parseSEQTAB(syn_bytes[seqtab_start+2 : len(syn_bytes)-2]); err != nil {
			return
		}
	}
	return
}

func parseSEQTAB(bytes []byte) (err error) {
	//PTVAL:	DS	64			;Current active processed pot value
	for i := 0; i < 4; i++ {
		logger.Debugf("SEQ TAB PTVAL[%d]: %x (%d\t%d)\n", i, bytes[i], bytes[i], int8(bytes[i]))

	}
	//TRANSP:	DS	5			;Sequencer playback transpose factor
	//SEQCON:	DS	40			;Sequence control table; 20 integers
	//SEQVOI:	DS	12			;(trk,3) fortran bit array for tracking
	//							;Voices 1-24 in seq tracks 1-4
	//VOIUSE:	DS	24			;Count of how many notes are currently
	//							;Using each voice - for finding first
	//							;Available voice in assignment mode
	for i := 4; i < (12 + 40 + 5); i++ {
		logger.Debugf("SEQ TAB VOIUSE-TRANSP[%d]: %x (%d\t%d)\n", i, bytes[i], bytes[i], int8(bytes[i]))

	}
	// FROM arithmetic in SEQ.Z80 (near L330B),  "time" appears to be a millisecond clock.
	// it appears to be a relative time from the previous event

	// (from SEQREQ.Z80
	//
	//;	The code sequence that is placed in "Seqtab" is as follows.
	//;	Time code (2 bytes), device code, value code.
	//;
	//;	Positive device codes have no value code following.
	//;	Negative device codes have different numbers of value codes.
	//
	//;	Device Code	Indication
	//;	---------------------------------------
	//;	 127	     = Time extend code
	//;	 126	     = Regular pedal down
	//;	 125	     = Middle pedal down
	//;	 124	     = Any pedal up
	//;	 123	     = Track #4 switch
	//;	 122	     = Track #3 switch
	//;	 121	     = Track #2 switch
	//;	 120	     = Track #1 switch
	//;	 119	     = RECOUT call
	//;	 118 to   75 = not used
	//;	  74 to    1 = Key up
	//;	   0	     = not used
	//;	  -1 to  -74 = Key down (key value byte follows)
	//;	 -75 to -114 = Transpose key (tracks affected byte follows)
	//;	-115	     = Pitchbend (bend value byte follows)
	//;	-116	     = Modulation (mod value follows)
	//;	-117 to -125 = not used
	//;	-126	     = 1 byte of Midi data follows
	//;	-127	     = 2 bytes of Midi data follows
	//;	-128	     = 3 bytes of Midi data follows
	for i := (12 + 40 + 5); i < len(bytes); i++ {
		logger.Debugf("SEQ TAB[%d]: %x (%d\t%d)\n", i, bytes[i], bytes[i], int8(bytes[i]))
	}
	for i := (2 + 12 + 40 + 5); i < len(bytes); {
		time := data.BytesToWord(bytes[i+1], bytes[i+0])
		device := int8(bytes[i+2])
		if device > 0 {
			if device <= 74 && device >= 1 {
				// key up
				logger.Debugf("EVENT [%d] time:%d  KEYUP k:%d \n", i, time, device)
			} else {
				logger.Debugf("UNKNOWN +EVENT [%d] time:%d  device:%d)\n", i, time, device)
			}
			// no data byte
			i += 3
		} else {
			// negative device codes have 1-3 data bytes depending on device code
			if device >= -74 && device <= -1 {
				// key down
				v := bytes[i+3]
				velocity := v >> 5 // top 3 bits == velocity
				voice := v & 0x1f  // bottom 5 bits = voice
				logger.Debugf("EVENT [%d] time:%d  KEYDOWN k:%d vel:%d voice:%d\n", i, time, -device, velocity, voice)
				i += 4
			} else if device == -116 {
				v := bytes[i+3]
				logger.Debugf("EVENT [%d] time:%d  MOD device:%d (%d\t%d)\n", i, time, device, v, v)
				i += 4
			} else if device == -115 {
				v := bytes[i+3]
				logger.Debugf("EVENT [%d] time:%d  BEND device:%d (%d\t%d)\n", i, time, device, v, v)
				i += 4
			} else if device == -126 {
				v := bytes[i+3]
				logger.Debugf("EVENT [%d] time:%d  MIDI 1-byte: device:%d (%d\t%d)\n", i, time, device, v, v)
				i += 4
			} else if device == -127 {
				v := []byte{bytes[i+3], bytes[i+4]}
				logger.Debugf("EVENT [%d] time:%d  MIDI 2-byte: device:%d (%d\t%d)\t(%d\t%d) \n", i, time, device, v[0], v[0], v[1], v[1])
				i += 5
			} else if device == -128 {
				v := []byte{bytes[i+3], bytes[i+4], bytes[i+5]}
				logger.Debugf("EVENT [%d] time:%d  MIDI 3-byte: device:%d (%d\t%d)\t(%d\t%d) \t(%d\t%d) \n", i, time, device, v[0], v[0], v[1], v[1], v[2], v[2])
				i += 6
			} else {
				v := bytes[i+3]
				logger.Debugf("UNKNOWN -EVENT [%d] time:%d  device:%d (%d\t%d)\n", i, time, device, v, v)
				i += 4
			}
		}
	}
	return
}
