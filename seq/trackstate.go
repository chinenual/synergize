package seq

import (
	"fmt"
	"math"

	"github.com/chinenual/synergize/logger"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

const comboTrackKey = 0
const extMidiTrackKey = -1
const modulationTrackKey = -2

const fakeTrackExtentCC = uint8(105) // a CC value that is documented as "unused" - so unlikely to get misinterpreted by DAW
const fakeTrackExtentStartVal = uint8(64)
const fakeTrackExtentEndVal = uint8(127)

type TrackPlayMode int

const (
	PlayModeOff TrackPlayMode = iota
	PlayModeOn
	PlayModeRepeat
)

type globalStateType struct {
	trackPlayMode [4]TrackPlayMode
	transpose     [4]int8
}

var globalState globalStateType

type trackStateType struct {
	trackID int

	hasEvents   bool
	byteIndex   int
	trackBytes  []byte
	trackMode   TrackMode
	evenMode    bool
	midiChannel uint8

	// absolute time (in ms) sync'd to the virtual clock
	absTime uint32 // 0 .. 4294967295 ms (4294967.295sec or ~71582 minutes or ~ 1193 hours - don't worry about overflow :))

	trackStartRelTime uint32 // first event (in track-relative ms)
	trackEndRelTime   uint32 // end of track time (in track-relative ms)

	allTracks []*[]timestampedMessage

	// key is a voice 1..24, or 0 for "everything on one track" or -1 for "external MIDI"
	trackMap map[int]*[]timestampedMessage

	lPedalDown bool
	rPedalDown bool

	extMidiTrackIndex int
	modulationTrack   []timestampedMessage
	// indexed by UNTRANSPOSED MIDI key code:
	activeKeyTracks [130]trackset
	// indexed by UNTRANSPOSED MIDI key code, value is the actual MIDI key code sent.  Value will be MIDI_NOKEY if not active
	// this allows us to handle the race condition where a key down is sent, then a transpose event, later the key up
	activeKeyTransposed [130]int
}

// used to store a key code to indicate nothing was active at this point
const MIDI_NOKEY = -1

func (ts *trackStateType) HasEvents() bool {
	return ts.hasEvents
}

func (ts *trackStateType) IsFirstEvent() bool {
	return ts.byteIndex == 0
}

const NoNextEvent = uint32(math.MaxUint32)

func (ts *trackStateType) InitNoData() {
	ts.hasEvents = false
}

func (ts *trackStateType) Init(trackID /*one based*/ int, trackBytes []byte, trackMode TrackMode, evenMode bool) {
	ts.hasEvents = true
	ts.trackID = trackID
	ts.trackMode = trackMode
	ts.evenMode = evenMode
	ts.trackBytes = trackBytes
	ts.byteIndex = 0

	// timestamp in (track-relative ms) of the first event
	ts.trackStartRelTime = 0
	// timestamp in (track-relative ms) of the end of the track
	ts.trackEndRelTime = 0

	ts.midiChannel = 0
	ts.absTime = 0
	ts.lPedalDown = false
	ts.rPedalDown = false
	ts.extMidiTrackIndex = -1 // init to a non-usable index - set during GetTrack() if there's any external midi data

	ts.trackMap = make(map[int]*[]timestampedMessage)

	for i := range ts.activeKeyTracks {
		ts.activeKeyTracks[i].Init()
		ts.activeKeyTransposed[i] = MIDI_NOKEY
	}
}

func (ts *trackStateType) ResetClock() {
	ts.absTime = 0
}

func (ts *trackStateType) ArmTrack() {
	ts.byteIndex = 0
}

func (ts *trackStateType) IsCalculatingTrackExtent() bool {
	return ts.trackEndRelTime == 0
}

func (ts *trackStateType) StartRelTime() uint32 {
	return ts.trackStartRelTime
}

func (ts *trackStateType) EndRelTime() uint32 {
	return ts.trackEndRelTime
}

func copyMessages(source []timestampedMessage, dest *[]timestampedMessage) {
	// skip the first event (the TrackSequenceName)
	for _, e := range source[1:] {
		*dest = append(*dest, e)
	}
}

func (ts *trackStateType) GetTrack(trackKey int) (midiTrack *[]timestampedMessage) {
	midiTrack = ts.trackMap[trackKey]
	if midiTrack == nil {
		midiTrack = new([]timestampedMessage)
		// add meta info to the track to name it:
		var name string
		switch trackKey {
		case extMidiTrackKey:
			name = fmt.Sprintf("SYN TRK %d EXTMIDI", ts.trackID)
		case modulationTrackKey:
			// the pseudo track for non-key events recorded before first note
			name = fmt.Sprintf("SYN TRK %d non-key", ts.trackID)
		case comboTrackKey:
			name = fmt.Sprintf("SYN TRK %d", ts.trackID)
		default:
			name = fmt.Sprintf("SYN TRK %d VOICE %d", ts.trackID, trackKey)
		}
		midievent := midi.Message(smf.MetaTrackSequenceName(name))
		e := timestampedMessage{0, midievent}
		*midiTrack = append(*midiTrack, e)
		ts.trackMap[trackKey] = midiTrack

		if trackKey != modulationTrackKey {
			if ts.modulationTrack != nil && trackKey != extMidiTrackKey {
				// copy non-key events into this track (all tracks except the extMidiTrack get copies of pb, mod and pedals)
				copyMessages(ts.modulationTrack, midiTrack)
			}
			// allTracks does not include the pseudo track
			ts.allTracks = append(ts.allTracks, midiTrack)
			if trackKey == extMidiTrackKey {
				ts.extMidiTrackIndex = len(ts.allTracks) - 1
			}
		}
	}
	return
}

func (ts *trackStateType) AddActiveKeyEvent(tm timestampedMessage, voice int, untransposedDevice int, device int) {
	var trackKey int
	if ts.trackMode == AllVoicesSameTrack {
		trackKey = comboTrackKey
	} else {
		trackKey = voice
	}
	midiTrack := ts.GetTrack(trackKey)
	*midiTrack = append(*midiTrack, tm)
	ts.activeKeyTracks[untransposedDevice].Add(trackKey)
	ts.activeKeyTransposed[untransposedDevice] = device
}

func (ts *trackStateType) ClearActiveKeyEvent(tm timestampedMessage, untransposedDevice int) {
	// for any track this event may have been written to:
	for _, k := range ts.activeKeyTracks[untransposedDevice].Contents() {
		logger.Debugf("Clearing event %d from track_key %d", untransposedDevice, k)
		midiTrack := ts.GetTrack(k)
		*midiTrack = append(*midiTrack, tm)
	}
	ts.activeKeyTracks[untransposedDevice].Clear()
	ts.activeKeyTransposed[untransposedDevice] = MIDI_NOKEY
}

func (ts *trackStateType) AddToAllActiveTracks(tm timestampedMessage) {
	// this is for non-key events (pb. mod, pedals).  If no track already allocated by a note event,
	// allocate a placeholder
	if ts.modulationTrack == nil {
		ts.modulationTrack = *ts.GetTrack(modulationTrackKey)
	}
	ts.modulationTrack = append(ts.modulationTrack, tm)
	for i, t := range ts.allTracks {
		if i != ts.extMidiTrackIndex {
			*t = append(*t, tm)
		}
	}
}

func (ts *trackStateType) AddModulation(val uint8) {
	if !ts.IsCalculatingTrackExtent() {
		m := midi.ControlChange(ts.midiChannel, midi.ModulationWheelMSB, val)
		tm := timestampedMessage{ts.absTime, m}
		ts.AddToAllActiveTracks(tm)
	}
}

func (ts *trackStateType) AddPitchbend(val int16) {
	if !ts.IsCalculatingTrackExtent() {
		m := midi.Pitchbend(ts.midiChannel, val)
		tm := timestampedMessage{ts.absTime, m}
		ts.AddToAllActiveTracks(tm)
	}
}

func (ts *trackStateType) AddKeyDown(untransposedKey int8, midiKey int8, midiVelocity uint8, synVoice uint8) {
	if !ts.IsCalculatingTrackExtent() {
		m := midi.NoteOn(ts.midiChannel, uint8(midiKey), midiVelocity)
		tm := timestampedMessage{ts.absTime, m}
		ts.AddActiveKeyEvent(tm, int(synVoice), int(untransposedKey), int(midiKey))
		logger.Debugf("ADD ACTIVE - map %d now %v\n", synVoice, ts.activeKeyTracks[untransposedKey])
	}
}

func (ts *trackStateType) AddKeyUp(untransposedKey int8, key int8) {
	if !ts.IsCalculatingTrackExtent() {
		// what was the transposed value of the key when originally pressed?
		k := ts.activeKeyTransposed[untransposedKey]

		m := midi.NoteOff(ts.midiChannel, uint8(k))
		tm := timestampedMessage{ts.absTime, m}
		// keyup needs to apply to all the keydown's - there may be more than one when multiple voices
		// use same msg
		ts.ClearActiveKeyEvent(tm, int(untransposedKey))
	}
}

func (ts *trackStateType) AddLeftPedalDown() {
	if !ts.IsCalculatingTrackExtent() {
		m := midi.ControlChange(ts.midiChannel, midi.PortamentoSwitch, 127)
		tm := timestampedMessage{ts.absTime, m}
		ts.lPedalDown = true
		ts.AddToAllActiveTracks(tm)
	}
}

func (ts *trackStateType) AddRightPedalDown() {
	if !ts.IsCalculatingTrackExtent() {
		m := midi.ControlChange(ts.midiChannel, midi.HoldPedalSwitch, 127)
		tm := timestampedMessage{ts.absTime, m}
		ts.rPedalDown = true
		ts.AddToAllActiveTracks(tm)
	}
}

func (ts *trackStateType) AddPedalUp() {
	if !ts.IsCalculatingTrackExtent() {
		if ts.lPedalDown {
			m := midi.ControlChange(ts.midiChannel, midi.PortamentoSwitch, 0)
			tm := timestampedMessage{ts.absTime, m}
			ts.AddToAllActiveTracks(tm)
			ts.lPedalDown = false
		}
		if ts.rPedalDown {
			m := midi.ControlChange(ts.midiChannel, midi.HoldPedalSwitch, 0)
			tm := timestampedMessage{ts.absTime, m}
			ts.AddToAllActiveTracks(tm)
			ts.rPedalDown = false
		}
	}
}

func (ts *trackStateType) AddExternalMidi(v []byte) {
	if !ts.IsCalculatingTrackExtent() {
		m := midi.Message(v)
		tm := timestampedMessage{ts.absTime, m}
		midiTrack := ts.GetTrack(extMidiTrackKey)
		*midiTrack = append(*midiTrack, tm)
	}
}

func (ts *trackStateType) MarkStartOfTrack() {
	if !ts.IsCalculatingTrackExtent() {

		// Explicitly mark the start of track in case user wants to create a repeat:
		// FIXME: investigate other SMF meta events?
		m := midi.ControlChange(ts.midiChannel, fakeTrackExtentCC, fakeTrackExtentStartVal)
		tm := timestampedMessage{ts.absTime, m}
		ts.AddToAllActiveTracks(tm)
	} else {
		ts.trackStartRelTime = ts.absTime
	}
}

func (ts *trackStateType) MarkEndOfTrack() {
	if !ts.IsCalculatingTrackExtent() {
		// Can't directly add MetaEndOfTrack since that's done behind the scenes in the midi library's Close()
		// function.  So we add a CC event:

		m := midi.ControlChange(ts.midiChannel, fakeTrackExtentCC, fakeTrackExtentEndVal)
		tm := timestampedMessage{ts.absTime, m}
		ts.AddToAllActiveTracks(tm)
	} else {
		ts.trackEndRelTime = ts.absTime
	}
}

func (ts *trackStateType) GetResult() (tracks [][]timestampedMessage, err error) {
	if len(ts.allTracks) == 0 {
		// only modulation events - no notes. So use the modulation track
		tracks = append(tracks, ts.modulationTrack)
	} else {
		// otherwise, modulation events have already been copied into all the note tracks - just use those
		for _, t := range ts.allTracks {
			tracks = append(tracks, *t)
		}
	}
	return
}
