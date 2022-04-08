package seq

import (
	"fmt"
	"io/ioutil"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"

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

func ConvertSYNToMIDI(path string, trackMode TrackMode) (err error) {
	var synBytes []byte

	if synBytes, err = ioutil.ReadFile(path); err != nil {
		return
	}

	cmosLen := data.BytesToWord(synBytes[1], synBytes[0])
	logger.Debugf("CMOS LEN: %v\n", cmosLen)

	seqtabStart := cmosLen*2 + 4
	logger.Debugf("SEQTAB start: 0x%x\n", seqtabStart)

	// seqlen is immediately after the cmos bytes - +2 for the header (cmosLen) bytes:
	seqLen := data.BytesToWord(synBytes[seqtabStart-1], synBytes[seqtabStart-2])
	logger.Debugf("SEQTAB LEN: %v\n", seqLen)

	if seqLen > 0 {
		var tracks [][]timestampedMessage
		// last two bytes of the file are the CRC
		if tracks, err = parseSEQTAB(synBytes[seqtabStart+2:len(synBytes)-2], trackMode); err != nil {
			return
		}

		midiPath := path + ".mid"

		s := smf.New()
		//s.TimeFormat = smf.SMPTE25(40)
		s.TimeFormat = smf.MetricTicks(960)

		for _, t := range tracks {
			var tr smf.Track

			time := uint32(0)
			for _, e := range t {
				deltaT := e.timeMS - time
				tr.Add(deltaT, e.msg)
				time = e.timeMS
			}
			tr.Close(0)

			if err = s.Add(tr); err != nil {
				return
			}
		}

		if err = s.WriteFile(midiPath); err != nil {
			return
		}
	}
	return
}

func parseSEQTAB(bytes []byte, trackMode TrackMode) (tracks [][]timestampedMessage, err error) {
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
		offset := 5 /* the 5 TRANSP bytes */ + 2*i + 2 /* unsure why seqcon seemes to be offset by 2 bytes */
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
	seqtabStart := 12 + 40 + 5
	for i := seqtabStart; i < len(bytes); i++ {
		logger.Debugf("SEQ TAB[%d]: %x (%d\t%d)\n", i, bytes[i], bytes[i], int8(bytes[i]))
	}

	for i := 0; i < NUMTRACKS; i++ {
		if trackStartStop[i].start != 0 {
			logger.Debugf("TRACK %d START %d STOP %d\n", i, uint16(seqtabStart)+trackStartStop[i].start-1, uint16(seqtabStart)+trackStartStop[i].stop)
			trackBytes := bytes[uint16(seqtabStart)+trackStartStop[i].start-1 : uint16(seqtabStart)+trackStartStop[i].stop]

			for i, b := range trackBytes {
				logger.Debugf("TRACK TAB[%d]: %x (%d\t%d)\n", i, b, b, int8(b))
			}

			var t [][]timestampedMessage
			if t, err = processTrack(i, trackBytes, trackMode); err != nil {
				return
			}
			for _, track := range t {
				tracks = append(tracks, track)
				logger.Debugf("MIDI EVENTS: %v\n", track)
			}
		}
	}
	return
}

type TrackMode int

const (
	AllVoicesSameTrack TrackMode = iota
	TrackPerVoice
)

type timestampedMessage struct {
	timeMS uint32
	msg    midi.Message
}

func (e timestampedMessage) String() string {
	return fmt.Sprintf("{%d: %s}", e.timeMS, e.msg)
}

func processTrack(track int, trackBytes []byte, trackMode TrackMode) (tracks [][]timestampedMessage, err error) {
	// FROM arithmetic in SEQ.Z80 (near L330B),  "time" appears to be a millisecond clock.
	// it appears to be a relative time from the previous msg

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

	// FIXME: how to map track/voice usage to MIDI channel?
	midiChannel := uint8(1)
	timeAccumulator := uint32(0)

	var allTracks []*[]timestampedMessage

	// key is a voice 1..24, or 0 for "everything on one track" or -1 for "external MIDI"
	var trackMap = make(map[int]*[]timestampedMessage)

	var lPedalDown = false
	var rPedalDown = false

	const comboTrackKey = 0
	const extMidiTrackKey = -1
	const modulationTrackKey = -2
	var modulationTrack []timestampedMessage
	var activeKeyTracks [130]trackset
	for i := range activeKeyTracks {
		activeKeyTracks[i].Init()
	}

	copyMessages := func(source []timestampedMessage, dest *[]timestampedMessage) {
		// skip the first event (the TrackSequenceName)
		for _, e := range source[1:] {
			*dest = append(*dest, e)
		}
	}

	getTrack := func(trackKey int) (midiTrack *[]timestampedMessage) {
		midiTrack = trackMap[trackKey]
		if midiTrack == nil {
			midiTrack = new([]timestampedMessage)
			// add meta info to the track to name it:
			var name string
			switch trackKey {
			case extMidiTrackKey:
				name = fmt.Sprintf("SYN TRK %d EXTMIDI", track+1)
			case modulationTrackKey:
				// the pseudo track for non-key events recorded before first note
				name = fmt.Sprintf("SYN TRK %d non-key", track+1)
			case comboTrackKey:
				name = fmt.Sprintf("SYN TRK %d", track+1)
			default:
				name = fmt.Sprintf("SYN TRK %d VOICE %d", track+1, trackKey)
			}
			midievent := midi.Message(smf.MetaTrackSequenceName(name))
			e := timestampedMessage{0, midievent}
			*midiTrack = append(*midiTrack, e)
			trackMap[trackKey] = midiTrack

			if trackKey != modulationTrackKey {
				if modulationTrack != nil {
					// copy non-key events into this track (all tracks get copies of pb, mod and pedals)
					copyMessages(modulationTrack, midiTrack)
				}
				// allTracks does not include the pseudo track
				allTracks = append(allTracks, midiTrack)
			}
		}
		return
	}

	addActiveKeyEvent := func(tm timestampedMessage, voice int, device int) {
		var trackKey int
		if trackMode == AllVoicesSameTrack {
			trackKey = comboTrackKey
		} else {
			trackKey = voice
		}
		midiTrack := getTrack(trackKey)
		*midiTrack = append(*midiTrack, tm)
		activeKeyTracks[device].Add(trackKey)
	}

	clearActiveKeyEvent := func(tm timestampedMessage, device int) {
		// for any track this event may have been written to:
		for _, k := range activeKeyTracks[device].Contents() {
			logger.Debugf("Clearing event %d from track_key %d", device, k)
			midiTrack := getTrack(k)
			*midiTrack = append(*midiTrack, tm)
		}
		activeKeyTracks[device].Clear()
	}

	addToAllActiveTracks := func(tm timestampedMessage) {
		// this is for non-key events (pb. mod, pedals).  If no track already allocated by a note event,
		// allocate a placeholder
		if modulationTrack == nil {
			modulationTrack = *getTrack(modulationTrackKey)
		}
		modulationTrack = append(modulationTrack, tm)
		for _, t := range allTracks {
			*t = append(*t, tm)
		}
	}

	// HACK: this extra +2 is adhoc - doesnt seem to match the firmware comments
	for i := 2; i < len(trackBytes); {
		logger.Debugf("TOP OF LOOP %d < %d\n", i, len(trackBytes))
		time := data.BytesToWord(trackBytes[i+1], trackBytes[i+0])
		timeAccumulator += uint32(time)
		if i+2 > len(trackBytes)-1 {
			// this was the last timestamp in the sequence - no device data - just "end of sequence" time
			logger.Debugf("t:%d END OF TRACK [%d] time:%d   \n", track, i, time)
			break
		}
		device := int8(trackBytes[i+2])
		if device > 0 {
			if device == 124 {
				// "any" pedal up
				logger.Debugf("t:%d EVENT [%d] time:%d  ANY PEDAL UP\n", track, i, time)
				if lPedalDown {
					m := midi.ControlChange(midiChannel, midi.PortamentoSwitch, 0)
					tm := timestampedMessage{timeAccumulator, m}
					addToAllActiveTracks(tm)
					lPedalDown = false
				}
				if rPedalDown {
					m := midi.ControlChange(midiChannel, midi.HoldPedalSwitch, 0)
					tm := timestampedMessage{timeAccumulator, m}
					addToAllActiveTracks(tm)
					rPedalDown = false
				}
			} else if device == 125 {
				// "middle" pedal down
				logger.Debugf("t:%d EVENT [%d] time:%d  LEFT PEDAL DOWN\n", track, i, time)
				m := midi.ControlChange(midiChannel, midi.PortamentoSwitch, 127)
				tm := timestampedMessage{timeAccumulator, m}
				lPedalDown = true
				addToAllActiveTracks(tm)
			} else if device == 126 {
				// "regular" pedal down
				logger.Debugf("t:%d EVENT [%d] time:%d  RIGHT PEDAL DOWN\n", track, i, time)
				m := midi.ControlChange(midiChannel, midi.HoldPedalSwitch, 127)
				tm := timestampedMessage{timeAccumulator, m}
				rPedalDown = true
				addToAllActiveTracks(tm)
			} else if device <= 74 && device >= 1 {
				// key up
				logger.Debugf("t:%d EVENT [%d] time:%d  KEYUP k:%d \n", track, i, time, device)
				m := midi.NoteOff(midiChannel, uint8(device))
				tm := timestampedMessage{timeAccumulator, m}
				// keyup needs to apply to all the keydown's - there may be more than one when multiple voices
				// use same msg
				clearActiveKeyEvent(tm, int(device))
			} else {
				logger.Errorf("t:%d INVALID +EVENT [%d] time:%d  device:%d)\n", track, i, time, device)
			}
			// no data byte
			i += 3
		} else {
			// negative device codes have 1-3 data bytes depending on device code
			if device >= -74 && device <= -1 {
				// key down
				v := trackBytes[i+3]
				velocity := v >> 5 // top 3 bits == velocity
				voice := v & 0x1f  // bottom 5 bits = voice
				if velocity == 0 {
					// MIDI velocity 0 means note off
					velocity = 1
				}
				// Synergy native velocity is 1..32
				// the values stored in the sequencer are compressed to 3-bits (so 0..7 !)
				// MIDI is 0..127 - so multiply by 18 to scale
				logger.Debugf("t:%d EVENT [%d] time:%d  KEYDOWN k:%d vel:%d voice:%d\n", track, i, time, -device, velocity, voice)
				m := midi.NoteOn(midiChannel, uint8(-device), 18*uint8(velocity))
				tm := timestampedMessage{timeAccumulator, m}
				addActiveKeyEvent(tm, int(voice), int(-device))
				logger.Debugf("ADD ACTIVE - map %d now %v\n", voice, activeKeyTracks[-device])

				i += 4
			} else if device == -116 {
				v := trackBytes[i+3]
				logger.Debugf("t:%d EVENT [%d] time:%d  MOD device:%d (%d\t%d)\n", track, i, time, device, v, v)
				// SYNERGY value is 0..48 - need to scale to 0..127 MIDI value
				midiVal := int(float32(v) * 127.0 / 48.0)
				if midiVal < 0 {
					midiVal = 0
				} else if midiVal > 127 {
					midiVal = 127
				}
				m := midi.ControlChange(midiChannel, midi.ModulationWheelMSB, uint8(midiVal))
				tm := timestampedMessage{timeAccumulator, m}
				addToAllActiveTracks(tm)
				i += 4
			} else if device == -115 {
				v := int8(trackBytes[i+3])
				logger.Debugf("t:%d EVENT [%d] time:%d  BEND device:%d (%d\t%d)\n", track, i, time, device, v, v)
				// SYNERGY value is -127..127 - need to scale to -8192..8191 MIDI value
				midiVal := int(float32(v) * 8191.0 / 48.0)
				if midiVal < -8192 {
					midiVal = -8192
				} else if midiVal > 8191 {
					midiVal = 8191
				}
				m := midi.Pitchbend(midiChannel, int16(midiVal))
				tm := timestampedMessage{timeAccumulator, m}
				addToAllActiveTracks(tm)
				i += 4
			} else if device == -126 {
				v := trackBytes[i+3]
				logger.Debugf("t:%d EVENT [%d] time:%d  MIDI 1-byte: device:%d (%d\t%d)\n", track, i, time, device, v, v)
				i += 4
			} else if device == -127 {
				v := []byte{trackBytes[i+3], trackBytes[i+4]}
				logger.Debugf("t:%d EVENT [%d] time:%d  MIDI 2-byte: device:%d (%d\t%d)\t(%d\t%d) \n", track, i, time, device, v[0], v[0], v[1], v[1])
				i += 5
			} else if device == -128 {
				v := []byte{trackBytes[i+3], trackBytes[i+4], trackBytes[i+5]}
				logger.Debugf("t:%d EVENT [%d] time:%d  MIDI 3-byte: device:%d (%d\t%d)\t(%d\t%d) \t(%d\t%d) \n", track, i, time, device, v[0], v[0], v[1], v[1], v[2], v[2])
				i += 6
			} else {
				v := trackBytes[i+3]
				logger.Errorf("t:%d INVALID -EVENT [%d] time:%d  device:%d (%d\t%d)\n", track, i, time, device, v, v)
				i += 4
			}
		}
	}
	if len(allTracks) == 0 {
		// only modulation events - no notes. So use the modulation track
		tracks = append(tracks, modulationTrack)
	} else {
		// otherwise, modulation events have already been copied into all the note tracks - just use those
		for _, t := range allTracks {
			tracks = append(tracks, *t)
		}
	}
	return
}
