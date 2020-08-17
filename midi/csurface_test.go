package midi

import (
	"flag"
	"fmt"
	"os"
	"testing"
	"time"
)

var (
	midiio = flag.Bool("midiio", false, "run integration tests that talk to MIDI")
)

func sendEventsArr(field string, fieldRange []int, valRange []int, delay time.Duration) (err error) {
	for i := fieldRange[0]; i <= fieldRange[1]; i++ {
		f := fmt.Sprintf("%s[%d]", field, i)
		for v := valRange[0]; v <= valRange[1]; v++ {
			fmt.Printf("  %s %v\n", f, v)
			if err = csSendEvent(f, uint8(v)); err != nil {
				return
			}
			time.Sleep(delay)
		}
	}
	return
}

func sendEvents(f string, valRange []int, delay time.Duration) (err error) {
	for v := valRange[0]; v <= valRange[1]; v++ {
		fmt.Printf("  %s %v\n", f, v)
		if err = csSendEvent(f, uint8(v)); err != nil {
			return
		}
		time.Sleep(delay)
	}

	return
}

func testEventsArr(t *testing.T, field string, fieldRange []int, valRange []int, delay time.Duration) {
	if err := sendEventsArr(field, fieldRange, valRange, delay); err != nil {
		fmt.Printf("%s failed %v\n", field, err)
		t.Fail()
	}
}

func testEvents(t *testing.T, field string, valRange []int, delay time.Duration) {
	if err := sendEvents(field, valRange, delay); err != nil {
		fmt.Printf("%s failed %v\n", field, err)
		t.Fail()
	}
}

func TestVoiceTab(t *testing.T) {
	if !*midiio {
		t.Skip()
	}
	testEvents(t, "voice-tab", []int{1, 1}, 10*time.Millisecond)

	// turn on:
	testEventsArr(t, "osc-enabled", []int{1, 16}, []int{1, 1}, 75*time.Millisecond)
	testEventsArr(t, "MUTE", []int{1, 16}, []int{1, 1}, 75*time.Millisecond)
	testEventsArr(t, "SOLO", []int{1, 16}, []int{1, 1}, 75*time.Millisecond)
	testEventsArr(t, "wkWAVE", []int{1, 16}, []int{1, 1}, 75*time.Millisecond)
	testEventsArr(t, "wkKEYPROP", []int{1, 16}, []int{1, 1}, 75*time.Millisecond)
	testEventsArr(t, "FILTER", []int{1, 16}, []int{0, 0}, 75*time.Millisecond)
	testEventsArr(t, "FILTER", []int{1, 16}, []int{1, 1}, 75*time.Millisecond)
	testEventsArr(t, "FILTER", []int{1, 16}, []int{2, 2}, 75*time.Millisecond)

	testEventsArr(t, "OHARM", []int{1, 16}, []int{0, 42}, 1*time.Millisecond)
	testEventsArr(t, "FDETUN", []int{1, 16}, []int{0, 127}, 1*time.Millisecond)

	// turn off
	testEventsArr(t, "osc-enabled", []int{1, 16}, []int{0, 0}, 75*time.Millisecond)
	testEventsArr(t, "SOLO", []int{1, 16}, []int{0, 0}, 75*time.Millisecond)
	testEventsArr(t, "MUTE", []int{1, 16}, []int{0, 0}, 75*time.Millisecond)
	testEventsArr(t, "wkWAVE", []int{1, 16}, []int{0, 0}, 75*time.Millisecond)
	testEventsArr(t, "wkKEYPROP", []int{1, 16}, []int{0, 0}, 75*time.Millisecond)

	testEventsArr(t, "FILTER", []int{1, 16}, []int{0, 0}, 75*time.Millisecond)

	testEvents(t, "VIBDEL", []int{0, 127}, 1*time.Millisecond)
	testEvents(t, "VIBRAT", []int{0, 127}, 1*time.Millisecond)
	testEvents(t, "VIBDEP", []int{0, 255}, 1*time.Millisecond)
	testEvents(t, "APVIB", []int{0, 255}, 1*time.Millisecond)
	testEvents(t, "VACENT", []int{0, 32}, 1*time.Millisecond)
	testEvents(t, "VASENS", []int{0, 31}, 1*time.Millisecond)
	testEvents(t, "VTCENT", []int{0, 32}, 1*time.Millisecond)
	testEvents(t, "VTSENS", []int{0, 31}, 1*time.Millisecond)
	testEvents(t, "VTRANS", []int{0, 255}, 1*time.Millisecond)
}
func TestVoiceFreqsTab(t *testing.T) {
	if !*midiio {
		t.Skip()
	}
	testEvents(t, "voice-freqs-tab", []int{1, 1}, 10*time.Millisecond)

	// turn on:
	testEventsArr(t, "osc-enabled", []int{1, 16}, []int{1, 1}, 75*time.Millisecond)

	testEventsArr(t, "OHARM", []int{1, 16}, []int{0, 42}, 1*time.Millisecond)
	testEventsArr(t, "FDETUN", []int{1, 16}, []int{0, 127}, 1*time.Millisecond)

	// turn off
	testEventsArr(t, "osc-enabled", []int{1, 16}, []int{0, 0}, 75*time.Millisecond)
}
func TestFreqEnvelopeTab(t *testing.T) {
	if !*midiio {
		t.Skip()
	}
	testEvents(t, "freq-envelopes-tab", []int{1, 1}, 10*time.Millisecond)

	testEventsArr(t, "envFreqLowVal", []int{1, 16}, []int{0, 124}, 1*time.Millisecond)
	testEventsArr(t, "envFreqUpVal", []int{1, 16}, []int{0, 124}, 1*time.Millisecond)
	testEventsArr(t, "envFreqLowTime", []int{1, 16}, []int{0, 84}, 1*time.Millisecond)
	testEventsArr(t, "envFreqUpTime", []int{1, 16}, []int{0, 84}, 1*time.Millisecond)
	testEvents(t, "accelFreqLow", []int{0, 127}, 1*time.Millisecond)
	testEvents(t, "accelFreqUp", []int{0, 127}, 1*time.Millisecond)

}
func TestAmpEnvelopeTab(t *testing.T) {
	if !*midiio {
		t.Skip()
	}
	testEvents(t, "amp-envelopes-tab", []int{1, 1}, 10*time.Millisecond)

	testEventsArr(t, "envAmpLowVal", []int{1, 16}, []int{0, 72}, 1*time.Millisecond)
	testEventsArr(t, "envAmpUpVal", []int{1, 16}, []int{0, 72}, 1*time.Millisecond)
	testEventsArr(t, "envAmpLowTime", []int{1, 16}, []int{0, 84}, 1*time.Millisecond)
	testEventsArr(t, "envAmpUpTime", []int{1, 16}, []int{0, 84}, 1*time.Millisecond)
	testEvents(t, "accelAmpLow", []int{0, 127}, 1*time.Millisecond)
	testEvents(t, "accelAmpUp", []int{0, 127}, 1*time.Millisecond)

}

func TestFiltersTab(t *testing.T) {
	if !*midiio {
		t.Skip()
	}
	testEvents(t, "filters-tab", []int{1, 1}, 10*time.Millisecond)

	testEventsArr(t, "flt", []int{1, 32}, []int{0, 127}, 1*time.Millisecond)
}
func TestKeyeqTab(t *testing.T) {
	if !*midiio {
		t.Skip()
	}
	testEvents(t, "keyeq-tab", []int{1, 1}, 10*time.Millisecond)

	testEventsArr(t, "keyeq", []int{1, 24}, []int{0, 31}, 1*time.Millisecond)
}

func TestKeypropTab(t *testing.T) {
	if !*midiio {
		t.Skip()
	}
	testEvents(t, "keyprop-tab", []int{1, 1}, 10*time.Millisecond)

	testEventsArr(t, "keyprop", []int{1, 24}, []int{0, 32}, 1*time.Millisecond)
}

func TestMain(m *testing.M) {

	flag.Parse()
	if *midiio {
		defer func() {
			fmt.Printf("Close Event.\n")
			if err := QuitMidi(); err != nil {
				fmt.Println(err)
			}
		}()
		err := InitMidi("Network Session 1", "./synergize-touchosc.cfg", false, false)
		if err != nil {
			fmt.Printf("could not initialize io: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Integration tests skipped. Run with -midiio to run them.\n")
	}
	os.Exit(m.Run())
}
