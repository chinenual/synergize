package seq

import (
	"fmt"
	"io/ioutil"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/writer"

	"gitlab.com/gomidi/midi/midimessage/channel"

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
		var tracks [][]timestampedEvent
		// last two bytes of the file are the CRC
		if tracks, err = parseSEQTAB(syn_bytes[seqtab_start+2 : len(syn_bytes)-2]); err != nil {
			return
		}

		midiPath := path + ".mid"
		if err = writer.WriteSMF(midiPath, uint16(len(tracks)), func(wr *writer.SMF) (err error) {
			time := uint64(0)
			for _, t := range tracks {
				for _, e := range t {
					if e.timeMS >= time {
						ms := e.timeMS - time
						ticks := uint32(ms) // FIXME:  ms->ticks conversion needs to be tempo aware
						wr.SetDelta(ticks)
					}
					if err = wr.Write(e.event); err != nil {
						return
					}
				}
				if err = writer.EndOfTrack(wr); err != nil {
					return
				}
			}
			return
		}); err != nil {
			return
		}

	}
	return
}

func parseSEQTAB(bytes []byte) (tracks [][]timestampedEvent, err error) {
	//PTVAL:	DS	64			;Current active processed pot value]
	const NUMTRACKS = 4
	for i := 0; i < NUMTRACKS; i++ {
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
	var seqcon [20]uint16
	for i := 0; i < 20; i++ {
		offset := 5 /* the 5 TRANSP bytes */ + 2*i
		seqcon[i] = data.BytesToWord(bytes[offset+1], bytes[offset])
	}
	logger.Debugf("SEQCON: %v\n", seqcon)
	var trackStartStop [NUMTRACKS]struct {
		start uint16
		stop  uint16
	}
	for i := 0; i < NUMTRACKS; i++ {
		trackStartStop[i].start = seqcon[i]
		trackStartStop[i].stop = seqcon[NUMTRACKS+i]
		logger.Debugf("TRACK %d SEQCON START %d STOP %d\n", i, trackStartStop[i].start, trackStartStop[i].stop)
	}
	start_of_seqtab := 12 + 40 + 5
	for i := start_of_seqtab; i < len(bytes); i++ {
		logger.Debugf("SEQ TAB[%d]: %x (%d\t%d)\n", i, bytes[i], bytes[i], int8(bytes[i]))
	}

	for i := 0; i < NUMTRACKS; i++ {
		if trackStartStop[i].start != 0 {
			logger.Debugf("TRACK %d START %d STOP %d\n", i, uint16(start_of_seqtab)+trackStartStop[i].start-1, uint16(start_of_seqtab)+trackStartStop[i].stop)
			track_bytes := bytes[uint16(start_of_seqtab)+trackStartStop[i].start-1 : uint16(start_of_seqtab)+trackStartStop[i].stop]

			var events []timestampedEvent
			if events, err = processTrack(i, track_bytes); err != nil {
				return
			}
			tracks = append(tracks, events)
			logger.Debugf("MIDI EVENTS: %v\n", events)
		}
	}
	return
}

type timestampedEvent struct {
	timeMS uint64
	event  midi.Message
}

var midiChannels = []channel.Channel{
	channel.Channel0,
	channel.Channel1,
	channel.Channel2,
	channel.Channel3,
	channel.Channel4,
	channel.Channel5,
	channel.Channel6,
	channel.Channel7,
	channel.Channel8,
	channel.Channel9,
	channel.Channel10,
	channel.Channel11,
	channel.Channel12,
	channel.Channel13,
	channel.Channel14,
	channel.Channel15,
}

func (e timestampedEvent) String() string {
	return fmt.Sprintf("{%d: %s}", e.timeMS, e.event)
}

func processTrack(track int, track_bytes []byte) (events []timestampedEvent, err error) {
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

	// (SEQCON is in the VOIUSE-TRANSP data above)
	//
	// SEQCN1	EQU	0			-- track start offsets
	// SEQCN2	EQU	2*NUMTRK	-- track end offsets
	// SEQCN3	EQU	4*NUMTRK	-- current track time (during play or record)
	// SEQCN4	EQU	6*NUMTRK	-- current track offset (during play or record)
	// SEQCN5	EQU	8*NUMTRK	-- current pedal state (during play or record)
	//
	// Track start is at SEQCON+SEQCN1+<track>
	// Track end   is at SEQCON+SEQCN2+<track>

	// HACK: this extra +2 is adhoc - doesnt seem to match the firmware comments

	// FIXME: how to map track/voice usage to MIDI channel?
	midiChannel := midiChannels[0]

	timeAccumulator := uint64(0)
	for i := 2; i < len(track_bytes); {
		time := data.BytesToWord(track_bytes[i+1], track_bytes[i+0])
		timeAccumulator += uint64(time)
		if i+2 > len(track_bytes)-1 {
			// this was the last timestamp in the sequence - no device data - just "end of sequence" time
			logger.Debugf("t:%d END OF SEQUENCE [%d] time:%d   \n", track, i, time)
			break
		}
		device := int8(track_bytes[i+2])
		if device > 0 {
			if device <= 74 && device >= 1 {
				// key up
				logger.Debugf("t:%d EVENT [%d] time:%d  KEYUP k:%d \n", track, i, time, device)
				midievent := midiChannel.NoteOff(uint8(device))
				e := timestampedEvent{timeAccumulator, midievent}
				events = append(events, e)
			} else {
				logger.Debugf("t:%d UNKNOWN +EVENT [%d] time:%d  device:%d)\n", track, i, time, device)
			}
			// no data byte
			i += 3
		} else {
			// negative device codes have 1-3 data bytes depending on device code
			if device >= -74 && device <= -1 {
				// key down
				v := track_bytes[i+3]
				velocity := v >> 5 // top 3 bits == velocity
				voice := v & 0x1f  // bottom 5 bits = voice
				// Synergy native velocity is 1..32
				// the values stored in the sequencer are compressed to 3-bits (so 0..7 !)
				// MIDI is 0..127 - so multiply by 18 to scale
				logger.Debugf("t:%d EVENT [%d] time:%d  KEYDOWN k:%d vel:%d voice:%d\n", track, i, time, -device, velocity, voice)
				midievent := midiChannel.NoteOn(uint8(-device), 18*uint8(velocity))
				e := timestampedEvent{timeAccumulator, midievent}
				events = append(events, e)
				i += 4
			} else if device == -116 {
				v := track_bytes[i+3]
				logger.Debugf("t:%d EVENT [%d] time:%d  MOD device:%d (%d\t%d)\n", track, i, time, device, v, v)
				// FIXME: need to scale this
				midievent := midiChannel.ControlChange(1, uint8(v))
				e := timestampedEvent{timeAccumulator, midievent}
				events = append(events, e)
				i += 4
			} else if device == -115 {
				v := track_bytes[i+3]
				logger.Debugf("t:%d EVENT [%d] time:%d  BEND device:%d (%d\t%d)\n", track, i, time, device, v, v)
				// FIXME: need to scale this
				midievent := midiChannel.Pitchbend(int16(v))
				e := timestampedEvent{timeAccumulator, midievent}
				events = append(events, e)
				i += 4
			} else if device == -126 {
				v := track_bytes[i+3]
				logger.Debugf("t:%d EVENT [%d] time:%d  MIDI 1-byte: device:%d (%d\t%d)\n", track, i, time, device, v, v)
				i += 4
			} else if device == -127 {
				v := []byte{track_bytes[i+3], track_bytes[i+4]}
				logger.Debugf("t:%d EVENT [%d] time:%d  MIDI 2-byte: device:%d (%d\t%d)\t(%d\t%d) \n", track, i, time, device, v[0], v[0], v[1], v[1])
				i += 5
			} else if device == -128 {
				v := []byte{track_bytes[i+3], track_bytes[i+4], track_bytes[i+5]}
				logger.Debugf("t:%d EVENT [%d] time:%d  MIDI 3-byte: device:%d (%d\t%d)\t(%d\t%d) \t(%d\t%d) \n", track, i, time, device, v[0], v[0], v[1], v[1], v[2], v[2])
				i += 6
			} else {
				v := track_bytes[i+3]
				logger.Debugf("t:%d UNKNOWN -EVENT [%d] time:%d  device:%d (%d\t%d)\n", track, i, time, device, v, v)
				i += 4
			}
		}
	}
	return
}
