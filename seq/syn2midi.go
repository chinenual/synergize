package seq

import (
	"fmt"
	"io/ioutil"
	"time"

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

func ConvertSYNToMIDI(path string, trackMode TrackMode, tempoBPM float64) (err error) {
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
		// 960 ppm is a widely supported format; SMPTE silently fails to load in Logic, so just stick with 960 ppm for now
		tempo := smf.MetricTicks(960)
		s.TimeFormat = tempo

		// convert ms to "ticks" in the current tempo
		msToTick := func(ms uint32) uint32 {
			return tempo.Ticks(tempoBPM, time.Duration(ms)*time.Millisecond)
		}

		for i, t := range tracks {
			var tr smf.Track

			if i == 0 {
				// first track gets tempo metadata:
				tr.Add(0, smf.MetaTempo(tempoBPM))
			}
			absTime := uint32(0)
			for _, e := range t {
				deltaT := msToTick(e.timeMS - absTime)
				tr.Add(deltaT, e.msg)
				absTime = e.timeMS
			}
			tr.Close(0)

			if err = s.Add(tr); err != nil {
				return
			}
		}

		if err = s.WriteFile(midiPath); err != nil {
			return
		}
		logger.Infof("Converted SYN %s to MIDI %s at %v BPM : %d MIDI tracks\n", path, midiPath, tempoBPM, len(tracks))
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
		// HACK: this +1 and +2 is adhoc and seems to work, but I don't understand why based on what I see in the Z80 code...
		trackStartStop[i].start = seqcon[i] + 1
		trackStartStop[i].stop = seqcon[NUMTRACKS+i] + 2
		logger.Debugf("TRACK %d SEQCON START %d STOP %d\n", i, trackStartStop[i].start, trackStartStop[i].stop)
	}
	seqtabStart := 12 + 40 + 5
	for i := seqtabStart; i < len(bytes); i++ {
		logger.Debugf("SEQ TAB[%d]: %x (%d\t%d)\n", i, bytes[i], bytes[i], int8(bytes[i]))
	}

	var trackState [4]trackState
	evenMode := false // not yet supported

	for i := 0; i < NUMTRACKS; i++ {
		if trackStartStop[i].start != 0 {
			logger.Debugf("TRACK %d START %d STOP %d\n", i, uint16(seqtabStart)+trackStartStop[i].start, uint16(seqtabStart)+trackStartStop[i].stop)
			trackBytes := bytes[uint16(seqtabStart)+trackStartStop[i].start : uint16(seqtabStart)+trackStartStop[i].stop]
			// even mode only applies to track1 (i == 0):
			trackState[i].Init(i+1, trackBytes, trackMode, evenMode && i == 0)
			for i, b := range trackBytes {
				logger.Debugf("TRACK TAB[%d]: %x (%d\t%d)\n", i, b, b, int8(b))
			}
			// first pass over each track does not create MIDI events; it is used to calculate the start and end time
			// of the first events on the track (which will be used to determine the "longest currently running track"
			// which drives repeat playback behavior
			if err = processTrack(&trackState[i], 0); err != nil {
				return
			}
			// Reset the track to start at the beginning when we do the real parsing
			trackState[i].ResetClock()
			logger.Debugf("TRACK START/END[%d]: START: %d  END:%d)\n", i, trackState[i].StartRelTime(), trackState[i].EndRelTime())
		}
	}
	for i := 0; i < NUMTRACKS; i++ {
		if trackStartStop[i].start != 0 {

			var t [][]timestampedMessage

			if err = processTrack(&trackState[i], 0); err != nil {
				return
			}
			if t, err = trackState[i].GetResult(); err != nil {
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

// Track playback algorithm:
// track buttons are OFF/ON/REPEAT
// From the user manual:
//
//   If you are playing back one track (whether set to Repeat or Playback mode) the initial rest (up to 8 seconds) is
//   omitted; rests at the end of the track are retained.
//
//   If you are playing several tracks in Playback or Repeat mode, the first note across all of the tracks is played as
//   soon as the Sequencer On/Of switch is pressed. Te initia rests of the other tracks are shorrted accordingly
//
//   If there are any Repeat tracks, the longest Repeat track governs the repeat. Rests are inserted after the end of
//   Playback tracks and other (shorter) Repeat tracks until the longest Repeat track is finished.  Then all of the
//   tracks repeat together. Initial rests are omitted in repeasts, just like at the beginning of Playback with one track.
//
//   If there is a Playback track longer than the longest Repeat track, then it continues playing even though other
//   tracks have started to repeat.  Rests are inserted at the end of this longer Playback track until the other tracks
//   repeat again

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

func processTrack(ts *trackState, dTimeAdjust uint32) (err error) {
	var next uint32
	// time adjustment applies only to the first event (to trim initial rests):
	ts.ArmTrack()
	if next, err = processNextEvent(ts, dTimeAdjust); err != nil || next == 0 {
		return
	}
	for {
		if next, err = processNextEvent(ts, 0); err != nil || next == 0 {
			return
		}
	}
	return
}

func processNextEvent(ts *trackState, dTimeAdjust uint32) (nextEventTime uint32, err error) {

	// FROM arithmetic in SEQ.Z80 (near L330B),  "time" is a millisecond clock.
	// it is a delta time from the previous event

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

	if ts.byteIndex >= len(ts.trackBytes)-1 {
		nextEventTime = 0
		return
	}
	var dTime uint16
	if ts.evenMode {
		// from User Manual:
		//
		// When you press the Even switch ... When you play back Track No. 1, all of the notes in this track,
		// including the initial rests, are evened out so that the notes are evenly spaced. With the Speed knob at
		// its normal setting, this works out to 1/4 second per note.
		// ...
		// There is a built-in chord detector which looks for chords in the track. A "chord" in this case is any
		// group of notes which have start times very close to each other.  When the Even feature finds a "chord"
		// in Track 1, it plays the entire chord in one beat , and plays the next note (or chord) on the following
		// beat.
		//
		// From SEQ.Z80, EVEN MODE TIMING (in L310:) uses lookahead and a 40ms threshold to determine if a note is
		// part of a "chord"

		// FIXME: what about non-note events? (pedals, pitchbend, modulation)... Will need to reverse engineer
		// the playback code in SEQ.Z80...
		dTime = 250 // 250ms = 1/4 second
	} else {
		dTime = data.BytesToWord(ts.trackBytes[ts.byteIndex+1], ts.trackBytes[ts.byteIndex+0])
	}
	ts.absTime += uint32(dTime) - dTimeAdjust

	if ts.byteIndex == 0 {
		// 'first event'
		ts.MarkStartOfTrack()
	}
	if ts.byteIndex+2 >= len(ts.trackBytes)-1 {
		// this was the last timestamp in the sequence - no device data - just "end of sequence" dTime
		logger.Debugf("t:%d END OF TRACK [%d] time:%d/%d   \n", ts.trackID, ts.byteIndex, dTime, ts.absTime)
		{
			// HACK: send key-up for any playing keys (Synergy sometimes records a key-down and never records the
			// key-up (or pedal up) at end of track).
			for k := int8(0); k < 127; k++ {
				if !ts.activeKeyTracks[k].Empty() {
					ts.AddKeyUp(k)
				}
			}
			// similarly for pedals:
			ts.AddPedalUp()
		}
		{
			// Now explicitly mark the end of track in case user wants to create a repeat:
			ts.MarkEndOfTrack()
		}
		// 'end of loop'
		nextEventTime = 0
		return
	}
	// Synergy firmware uses term "device" for this value, so we do too
	device := int8(ts.trackBytes[ts.byteIndex+2])
	if device > 0 {
		if device == 124 {
			// "any" pedal up
			logger.Debugf("t:%d EVENT [%d] time:%d/%d  ANY PEDAL UP\n", ts.trackID, ts.byteIndex, dTime, ts.absTime)
			ts.AddPedalUp()
		} else if device == 125 {
			// "middle" pedal down
			logger.Debugf("t:%d EVENT [%d] time:%d/%d  LEFT PEDAL DOWN\n", ts.trackID, ts.absTime, ts.byteIndex, dTime)
			ts.AddLeftPedalDown()
		} else if device == 126 {
			// "regular" pedal down
			logger.Debugf("t:%d EVENT [%d] time:%d/%d  RIGHT PEDAL DOWN\n", ts.trackID, ts.byteIndex, dTime, ts.absTime)
			ts.AddRightPedalDown()
		} else if device == 127 {
			// "time extend code"
			// nop for us
		} else if device <= 74 && device >= 1 {
			// key up
			key := device + 27 // MIDI key is offset from internal SYNERGY key code
			logger.Debugf("t:%d EVENT [%d] time:%d/%d  KEYUP k:%d \n", ts.trackID, ts.byteIndex, dTime, ts.absTime, key)
			ts.AddKeyUp(key)

		} else {
			logger.Errorf("t:%d INVALID +EVENT [%d] time:%d/%d  device:%d)\n", ts.trackID, ts.byteIndex, dTime, ts.absTime, device)
		}
		// no data byte
		ts.byteIndex += 3
	} else {
		// negative device codes have 1-3 data bytes depending on device code
		if device >= -74 && device <= -1 {
			// key down
			key := -device + 27 // MIDI key is offset from internal SYNERGY key code
			v := ts.trackBytes[ts.byteIndex+3]
			velocity := v >> 5 // top 3 bits == velocity
			voice := v & 0x1f  // bottom 5 bits = voice
			if velocity == 0 {
				// MIDI velocity 0 means note off
				velocity = 1
			}
			// Synergy native velocity is 1..32
			// the values stored in the sequencer are compressed to 3-bits (so 0..7 !)
			// MIDI is 0..127 - so multiply by 18 to scale
			logger.Debugf("t:%d EVENT [%d] time:%d/%d  KEYDOWN k:%d vel:%d voice:%d\n", ts.trackID, ts.byteIndex, dTime, ts.absTime, key, velocity, voice)
			ts.AddKeyDown(key, 18*uint8(velocity), voice)

			ts.byteIndex += 4
		} else if device == -116 {
			v := ts.trackBytes[ts.byteIndex+3]
			logger.Debugf("t:%d EVENT [%d] time:%d/%d  MOD device:%d (%d\t%d)\n", ts.trackID, ts.byteIndex, dTime, ts.absTime, device, v, v)
			// SYNERGY value is 0..48 - need to scale to 0..127 MIDI value
			midiVal := int(float32(v) * 127.0 / 48.0)
			if midiVal < 0 {
				midiVal = 0
			} else if midiVal > 127 {
				midiVal = 127
			}
			ts.AddModulation(uint8(midiVal))
			ts.byteIndex += 4
		} else if device == -115 {
			v := int8(ts.trackBytes[ts.byteIndex+3])
			logger.Debugf("t:%d EVENT [%d] time:%d/%d  BEND device:%d (%d\t%d)\n", ts.trackID, ts.byteIndex, dTime, ts.absTime, device, v, v)
			// SYNERGY value is -127..127 - need to scale to -8192..8191 MIDI value
			midiVal := int(float32(v) * 8191.0 / 48.0)
			if midiVal < -8192 {
				midiVal = -8192
			} else if midiVal > 8191 {
				midiVal = 8191
			}
			ts.AddPitchbend(int16(midiVal))
			ts.byteIndex += 4
		} else if device == -126 {
			v := []byte{ts.trackBytes[ts.byteIndex+3]}
			ts.AddExternalMidi(v)
			logger.Debugf("t:%d EVENT [%d] time:%d/%d  MIDI 1-byte: device:%d (%d\t%d)\n", ts.trackID, ts.byteIndex, dTime, ts.absTime, device, v, v)

			ts.byteIndex += 4
		} else if device == -127 {
			v := []byte{ts.trackBytes[ts.byteIndex+3], ts.trackBytes[ts.byteIndex+4]}
			ts.AddExternalMidi(v)
			logger.Debugf("t:%d EVENT [%d] time:%d/%d  MIDI 2-byte: device:%d (%d\t%d)\t(%d\t%d) \n", ts.trackID, ts.byteIndex, dTime, ts.absTime, device, v[0], v[0], v[1], v[1])
			ts.byteIndex += 5
		} else if device == -128 {
			v := []byte{ts.trackBytes[ts.byteIndex+3], ts.trackBytes[ts.byteIndex+4], ts.trackBytes[ts.byteIndex+5]}
			ts.AddExternalMidi(v)
			logger.Debugf("t:%d EVENT [%d] time:%d/%d  MIDI 3-byte: device:%d (%d\t%d)\t(%d\t%d) \t(%d\t%d) \n", ts.trackID, ts.byteIndex, dTime, ts.absTime, device, v[0], v[0], v[1], v[1], v[2], v[2])
			ts.byteIndex += 6
		} else {
			v := ts.trackBytes[ts.byteIndex+3]
			logger.Errorf("t:%d INVALID -EVENT [%d] time:%d/%d  device:%d (%d\t%d)\n", ts.trackID, ts.byteIndex, dTime, ts.absTime, device, v, v)
			ts.byteIndex += 4
		}
	}

	if ts.byteIndex >= len(ts.trackBytes)-1 {
		nextEventTime = 0
	} else {
		var dTime uint16
		dTime = data.BytesToWord(ts.trackBytes[ts.byteIndex+1], ts.trackBytes[ts.byteIndex+0])
		nextEventTime = ts.absTime + uint32(dTime)
	}
	return
}
