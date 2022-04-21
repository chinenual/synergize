package seq

import (
	"fmt"
	"testing"

	"github.com/chinenual/synergize/data"
)

// first test that the firmware emulation used to sanity check the round trip velocities do what the firmware actually
// does:
func TestVelocitySEQtoNativeEmulation(t *testing.T) {
	tests := []struct {
		seq    uint8
		native uint8
	}{
		{0, 0},
		{1, 4},
		{2, 8},
		{3, 12},
		{4, 16},
		{5, 20},
		{6, 24},
		{7, 32},
	}
	for _, v := range tests {
		data.AssertInt(t, int(seqVelocityToSynergy(v.seq)), int(v.native), fmt.Sprintf("(for seq %d)", v.seq))
	}
}

var midiNativeTests = []struct {
	native uint8
	midi   uint8
}{
	{0, 1},
	{1, 3},
	{2, 7},
	{3, 11},
	{4, 15},
	{5, 19},
	{6, 23},
	{7, 27},
	{8, 31},
	{9, 35},
	{10, 39},
	{11, 43},
	{12, 47},
	{13, 51},
	{14, 55},
	{15, 59},
	{16, 63},
	{17, 67},
	{18, 71},
	{19, 75},
	{20, 79},
	{21, 83},
	{22, 87},
	{23, 91},
	{24, 95},
	{25, 99},
	{26, 103},
	{27, 107},
	{28, 111},
	{29, 115},
	{30, 119},
	{31, 123},
	{32, 127},
}

func TestVelocityNativeToMIDIEmulation(t *testing.T) {
	for _, v := range midiNativeTests {
		data.AssertInt(t, int(synergyVelocityToMIDI(v.native)), int(v.midi), fmt.Sprintf("(for native %d)", v.native))
	}
}

func TestVelocityMIDIToNativeEmulation(t *testing.T) {
	for _, v := range midiNativeTests {
		expect := v.native
		if v.midi == 1 {
			// special case: the synergy does not convert MIDI 1 to native 0
			expect = 1
		}
		data.AssertInt(t, int(midiVelocityToSynergy(v.midi)), int(expect), fmt.Sprintf("(for MIDI %d)", v.midi))
	}
}

func TestVelocityRoundTrip(t *testing.T) {
	// for any of the 8 "sequencer" velocities, compare how the synergy firmware converted them to native velocity
	// to the value we get if we convert to midi and back:
	for seq := uint8(0); seq <= 7; seq++ {
		expect := seqVelocityToSynergy(seq)
		if expect == 0 {
			// special case: the synergy does not convert MIDI 1 to native 0
			expect = 1
		}
		toMidi := seqVelocityToMIDI(seq)
		fromMidi := midiVelocityToSynergy(toMidi)
		data.AssertInt(t, int(fromMidi), int(expect), fmt.Sprintf("(for seq %d -> MIDI %d)", seq, toMidi))
	}
}
