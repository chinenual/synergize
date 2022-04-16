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
	ts.extMidiTrackIndex = -1 // init to a non-usable index - set during getTrack() if there's any external midi data

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

func (ts *trackState) getTrack(trackKey int) (midiTrack *[]timestampedMessage) {
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

func (ts *trackState) addActiveKeyEvent(tm timestampedMessage, voice int, device int) {
	var trackKey int
	if ts.trackMode == AllVoicesSameTrack {
		trackKey = comboTrackKey
	} else {
		trackKey = voice
	}
	midiTrack := ts.getTrack(trackKey)
	*midiTrack = append(*midiTrack, tm)
	ts.activeKeyTracks[device].Add(trackKey)
}

func (ts *trackState) clearActiveKeyEvent(tm timestampedMessage, device int) {
	// for any track this event may have been written to:
	for _, k := range ts.activeKeyTracks[device].Contents() {
		logger.Debugf("Clearing event %d from track_key %d", device, k)
		midiTrack := ts.getTrack(k)
		*midiTrack = append(*midiTrack, tm)
	}
	ts.activeKeyTracks[device].Clear()
}

func (ts *trackState) addToAllActiveTracks(tm timestampedMessage) {
	// this is for non-key events (pb. mod, pedals).  If no track already allocated by a note event,
	// allocate a placeholder
	if ts.modulationTrack == nil {
		ts.modulationTrack = *ts.getTrack(modulationTrackKey)
	}
	ts.modulationTrack = append(ts.modulationTrack, tm)
	for i, t := range ts.allTracks {
		if i != ts.extMidiTrackIndex {
			*t = append(*t, tm)
		}
	}
}

func (ts *trackState) keyUp(absTime uint32, key int8) {
	m := midi.NoteOff(ts.midiChannel, uint8(key))
	tm := timestampedMessage{absTime, m}
	// keyup needs to apply to all the keydown's - there may be more than one when multiple voices
	// use same msg
	ts.clearActiveKeyEvent(tm, int(key))
}

func (ts *trackState) pedalUp(absTime uint32) {
	if ts.lPedalDown {
		m := midi.ControlChange(ts.midiChannel, midi.PortamentoSwitch, 0)
		tm := timestampedMessage{absTime, m}
		ts.addToAllActiveTracks(tm)
		ts.lPedalDown = false
	}
	if ts.rPedalDown {
		m := midi.ControlChange(ts.midiChannel, midi.HoldPedalSwitch, 0)
		tm := timestampedMessage{absTime, m}
		ts.addToAllActiveTracks(tm)
		ts.rPedalDown = false
	}
}
