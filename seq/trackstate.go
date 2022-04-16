package seq

import (
	"fmt"

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

type trackState struct {
	trackID     int
	trackBytes  []byte
	trackMode   TrackMode
	midiChannel uint8
	absTime     uint32

	allTracks []*[]timestampedMessage

	// key is a voice 1..24, or 0 for "everything on one track" or -1 for "external MIDI"
	trackMap map[int]*[]timestampedMessage

	lPedalDown bool
	rPedalDown bool

	extMidiTrackIndex int
	modulationTrack   []timestampedMessage
	activeKeyTracks   [130]trackset
}

func (ts *trackState) Init(trackID /*one based*/ int, trackBytes []byte, trackMode TrackMode) {
	ts.trackID = trackID
	ts.trackMode = trackMode
	ts.trackBytes = trackBytes
	ts.midiChannel = 0
	ts.absTime = 0
	ts.lPedalDown = false
	ts.rPedalDown = false
	ts.extMidiTrackIndex = -1 // init to a non-usable index - set during GetTrack() if there's any external midi data

	ts.trackMap = make(map[int]*[]timestampedMessage)

	for i := range ts.activeKeyTracks {
		ts.activeKeyTracks[i].Init()
	}
}

func copyMessages(source []timestampedMessage, dest *[]timestampedMessage) {
	// skip the first event (the TrackSequenceName)
	for _, e := range source[1:] {
		*dest = append(*dest, e)
	}
}

func (ts *trackState) GetTrack(trackKey int) (midiTrack *[]timestampedMessage) {
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

func (ts *trackState) AddActiveKeyEvent(tm timestampedMessage, voice int, device int) {
	var trackKey int
	if ts.trackMode == AllVoicesSameTrack {
		trackKey = comboTrackKey
	} else {
		trackKey = voice
	}
	midiTrack := ts.GetTrack(trackKey)
	*midiTrack = append(*midiTrack, tm)
	ts.activeKeyTracks[device].Add(trackKey)
}

func (ts *trackState) ClearActiveKeyEvent(tm timestampedMessage, device int) {
	// for any track this event may have been written to:
	for _, k := range ts.activeKeyTracks[device].Contents() {
		logger.Debugf("Clearing event %d from track_key %d", device, k)
		midiTrack := ts.GetTrack(k)
		*midiTrack = append(*midiTrack, tm)
	}
	ts.activeKeyTracks[device].Clear()
}

func (ts *trackState) AddToAllActiveTracks(tm timestampedMessage) {
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

func (ts *trackState) AddModulation(val uint8) {
	m := midi.ControlChange(ts.midiChannel, midi.ModulationWheelMSB, val)
	tm := timestampedMessage{ts.absTime, m}
	ts.AddToAllActiveTracks(tm)
}

func (ts *trackState) AddPitchbend(val int16) {
	m := midi.Pitchbend(ts.midiChannel, val)
	tm := timestampedMessage{ts.absTime, m}
	ts.AddToAllActiveTracks(tm)
}

func (ts *trackState) AddKeyDown(key int8, velocity uint8, synVoice uint8) {
	m := midi.NoteOn(ts.midiChannel, uint8(key), 18*uint8(velocity))
	tm := timestampedMessage{ts.absTime, m}
	ts.AddActiveKeyEvent(tm, int(synVoice), int(key))
	logger.Debugf("ADD ACTIVE - map %d now %v\n", synVoice, ts.activeKeyTracks[key])

}

func (ts *trackState) AddKeyUp(key int8) {
	m := midi.NoteOff(ts.midiChannel, uint8(key))
	tm := timestampedMessage{ts.absTime, m}
	// keyup needs to apply to all the keydown's - there may be more than one when multiple voices
	// use same msg
	ts.ClearActiveKeyEvent(tm, int(key))
}

func (ts *trackState) AddLeftPedalDown() {
	m := midi.ControlChange(ts.midiChannel, midi.PortamentoSwitch, 127)
	tm := timestampedMessage{ts.absTime, m}
	ts.lPedalDown = true
	ts.AddToAllActiveTracks(tm)
}

func (ts *trackState) AddRightPedalDown() {
	m := midi.ControlChange(ts.midiChannel, midi.HoldPedalSwitch, 127)
	tm := timestampedMessage{ts.absTime, m}
	ts.rPedalDown = true
	ts.AddToAllActiveTracks(tm)
}

func (ts *trackState) AddPedalUp() {
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

func (ts *trackState) AddExternalMidi(v []byte) {
	m := midi.Message(v)
	tm := timestampedMessage{ts.absTime, m}
	midiTrack := ts.GetTrack(extMidiTrackKey)
	*midiTrack = append(*midiTrack, tm)
}

func (ts *trackState) MarkStartOfTrack() {
	// Explicitly mark the start of track in case user wants to create a repeat:
	// FIXME: investigate other SMF meta events?
	m := midi.ControlChange(ts.midiChannel, fakeTrackExtentCC, fakeTrackExtentStartVal)
	tm := timestampedMessage{ts.absTime, m}
	ts.AddToAllActiveTracks(tm)
}

func (ts *trackState) MarkEndOfTrack() {

	// Can't directly add MetaEndOfTrack since that's done behind the scenes in the midi library's Close()
	// function.  So we add a CC event:

	m := midi.ControlChange(ts.midiChannel, fakeTrackExtentCC, fakeTrackExtentEndVal)
	tm := timestampedMessage{ts.absTime, m}
	ts.AddToAllActiveTracks(tm)

}
