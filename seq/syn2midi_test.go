package seq

import (
	"testing"
)

func TestNoteOn(t *testing.T) {

	tests := []struct {
		trackBytes    []byte
		expectedTime  uint32
		expectedKey   uint8
		expectedVel   uint8
		expectedVoice uint8
	}{
		{
			[]byte{0x0, 0x0, 0x0, 0x0, 0xce /*-50*/, uint8(5<<5 | 22), 0x0, 0x0},
			0,
			50,
			90, // scaled 5
			22,
		},
		{
			[]byte{0x0, 0x0, 0x0, 0x17, 0xba /*-70*/, uint8(1<<5 | 1), 0x0, 0x0},
			5888,
			70,
			18, // scaled 1
			1,
		},
	}
	for _, test := range tests {

		//logger.Init("", logger.LevelDebug)

		var tracks [][]timestampedMessage
		var err error
		if tracks, err = processTrack(0, test.trackBytes, AllVoicesSameTrack); err != nil {
			t.Errorf("Unexpected error: %v\n", err)
		}
		if len(tracks) != 1 {
			t.Errorf("Expected one track, got %d\n", len(tracks))
		}
		if len(tracks[0]) != 2 {
			t.Errorf("Expected two events on track, got %d\n", len(tracks[0]))
		}
		// skip the header event:
		tm := tracks[0][1]
		var channel, key, vel uint8
		if is := tm.msg.GetNoteOn(&channel, &key, &vel); !is {
			t.Errorf("Expected NoteOn msg - got: %v\n", tm.msg)
		}
		if channel != 0 {
			t.Errorf("Expected NoteOn channel == 0 - got: %v\n", channel)
		}
		if test.expectedKey != key {
			t.Errorf("Expected NoteOn key == %d - got: %v\n", test.expectedKey, key)
		}
		if test.expectedVel != vel {
			t.Errorf("Expected NoteOn velocity == %d - got: %v\n", test.expectedVel, vel)
		}
	}
}
