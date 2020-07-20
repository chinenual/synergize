package synio

import (
	"bytes"
	"github.com/chinenual/synergize/data"
	"log"
	//	"github.com/orcaman/writerseeker"
	"testing"
)

func TestReloadNoteGenerators(t *testing.T) {
	if !*synio {
		t.Skip()
	}
	var err error
	if err = ReloadNoteGenerators(); err != nil {
		t.Fatalf("Error reloading note generators: %v\n", err)
	}
}

func TestSetAPVIB(t *testing.T) {
	if !*synio {
		t.Skip()
	}
	var err error
	for v := 0; v <= 0xff; v++ {
		if err = SetVoiceHeadDataByte("APVIB", byte(v)); err != nil {
			t.Fatalf("Error setting APVIB value %v", v)
		}
	}
}

func TestSetOHARM(t *testing.T) {
	if !*synio {
		t.Skip()
	}
	var err error
	for osc := 0; osc < 2; osc++ {
		for v := -64; v <= 64; v++ {
			if err = SetVoiceOscDataByte(osc, "OHARM", byte(v)); err != nil {
				t.Fatalf("Error setting OHARM osc %v value %v", osc, v)
			}
		}
	}
}

func TestSetFDETUN(t *testing.T) {
	if !*synio {
		t.Skip()
	}
	var err error
	for osc := 0; osc < 2; osc++ {
		for v := -64; v <= 64; v++ {
			if err = SetVoiceOscDataByte(osc, "FDETUN", byte(v)); err != nil {
				t.Fatalf("Error setting FDETUN osc %v value %v", osc, v)
			}
		}
	}
}

func TestRoundTrip(t *testing.T) {
	if !*synio {
		t.Skip()
	}

	// Load a known VCE
	// change every param that can be changed
	// dump VRAM and create a new VCE
	// compare each changed param

	var vce data.VCE
	var err error
	if _, err = EnableVoicingMode(); err != nil {
		t.Fatalf("Failed to init voicing mode")
	}

	if vce, err = data.ReadVceFile("../data/testfiles/G7S.VCE"); err != nil {
		t.Fatalf("error reading G7S.VCE: %v", err)
	}

	if err = LoadVceVoicingMode(vce); err != nil {
		t.Fatalf("error LoadVceVoicingMode: %v", err)
	}

	type namedByte struct {
		Name  string
		Value byte
	}
	var headValues = []namedByte{
		{Name: "VTRANS", Value: 1}, // G7S: 0
		{Name: "VTCENT", Value: 1}, // G7S: 20
		{Name: "VTSENS", Value: 1}, // G7S: 3
		{Name: "VACENT", Value: 1}, // G7S: 15
		{Name: "VASENS", Value: 1}, // G7S: 27
		{Name: "VIBRAT", Value: 1}, // G7S: 0
		{Name: "VIBDEL", Value: 1}, // G7S:0
		{Name: "VIBDEP", Value: 1}, // G7S: 0
		{Name: "APVIB", Value: 1},  // G7S: 0
	}
	//	{Name: "VEQ    [24]int8
	//	{Name: "KPROP  [24]byte
	if err = SetVNAME("ROUNDTRIP"); err != nil { // GS7: GS7
		t.Fatalf("Failed to set SetVNAME")
	}
	if err = SetVoiceVEQEle(1, 1); err != nil { // GS7: -4
		t.Fatalf("Failed to set SetVoiceVEQEle")
	}
	if err = SetVoiceKPROPEle(1, 1); err != nil { // GS7: 0
		t.Fatalf("Failed to set SetVoiceKPROPEle")
	}
	for _, v := range headValues {
		if err = SetVoiceHeadDataByte(v.Name, v.Value); err != nil {
			t.Fatalf("Failed to set VoiceHeadDataByte %s to %v", v.Name, v.Value)
		}
	}
	if err = SetFilterEle(1, 1, 1); err != nil {
		t.Fatalf("Failed to set filter ele val 1 1 1 ")
	}
	if err = SetFilterEle(2, 1, 1); err != nil {
		t.Fatalf("Failed to set filter ele val 2 1 1 ")
	}

	var dumpedBytes []byte

	if dumpedBytes, err = DumpVRAM(); err != nil {
		t.Fatalf("DumpVRAM failed %v", err)
	}

	var readbuf = bytes.NewReader(dumpedBytes)
	var dumpedCrt data.CRT
	var dumpedVce data.VCE

	if dumpedCrt, err = data.ReadCrt(readbuf); err != nil {
		t.Errorf("error parsing dumpedVRAM %v", err)
		return
	}
	dumpedVce = *dumpedCrt.Voices[0]
	log.Printf("dumped CRT: %s\n", data.CrtToJson(dumpedCrt))
	log.Printf("dumped VCE: %s\n", data.VceToJson(dumpedVce))

	if data.VceName(dumpedVce.Head) != "ROUNDTRI" {
		t.Errorf("Round trip value failed %s to %v - got %v", "VEQ[0]", "ROUNDTRI", data.VceName(dumpedVce.Head))
	}
	if dumpedVce.Head.VEQ[0] != 1 {
		t.Errorf("Round trip value failed %s to %v - got %v", "VEQ[0]", 1, dumpedVce.Head.VEQ[0])
	}
	if dumpedVce.Head.KPROP[0] != 1 {
		t.Errorf("Round trip value failed %s to %v - got %v", "KPROP[0]", 1, dumpedVce.Head.KPROP[0])
	}
	for _, v := range headValues {
		var value byte
		switch v.Name {
		case "VTRANS":
			value = byte(dumpedVce.Head.VTRANS)
		case "VTCENT":
			value = dumpedVce.Head.VTCENT
		case "VTSENS":
			value = dumpedVce.Head.VTSENS
		case "VACENT":
			value = dumpedVce.Head.VACENT
		case "VASENS":
			value = dumpedVce.Head.VASENS
		case "VIBRAT":
			value = dumpedVce.Head.VIBRAT
		case "VIBDEL":
			value = dumpedVce.Head.VIBDEL
		case "VIBDEP":
			value = dumpedVce.Head.VIBDEP
		case "APVIB":
			value = dumpedVce.Head.APVIB
		default:
			t.Errorf("Unhandled field %s", v.Name)
		}
		if v.Value != value {
			t.Errorf("Round trip value failed %s to %v - got %v", v.Name, v.Value, value)
		}
	}
	// #osc and filters should be unchanged except for the first bytes
	if dumpedVce.Head.VOITAB != 3 {
		t.Errorf("Round trip value failed %s - expected %v - got %v", "VOITAB", 3, dumpedVce.Head.VOITAB)
	}
	if dumpedVce.Head.FILTER[0] != 1 {
		t.Errorf("Round trip value failed %s - expected %v - got %v", "FILTER[0]", 1, dumpedVce.Head.FILTER[0])
	}
	if dumpedVce.Head.FILTER[1] != 2 {
		t.Errorf("Round trip value failed %s - expected %v - got %v", "FILTER[1]", 2, dumpedVce.Head.FILTER[1])
	}
	if dumpedVce.Head.FILTER[2] != 0 {
		t.Errorf("Round trip value failed %s - expected %v - got %v", "FILTER[2]", 0, dumpedVce.Head.FILTER[2])
	}
	if dumpedVce.Head.FILTER[3] != 0 {
		t.Errorf("Round trip value failed %s - expected %v - got %v", "FILTER[3]", 0, dumpedVce.Head.FILTER[3])
	}
	if dumpedVce.Filters[0][0] != 1 {
		t.Errorf("Round trip value failed %s to %v - got %v", "Filters[0][0]", 1, dumpedVce.Filters[0][0])
	}
	if dumpedVce.Filters[0][1] != 0 {
		t.Errorf("Round trip value failed %s to %v - got %v", "Filters[0][1]", 0, dumpedVce.Filters[0][0])
	}
	if dumpedVce.Filters[1][0] != 1 {
		t.Errorf("Round trip value failed %s to %v - got %v", "Filters[0][0]", 1, dumpedVce.Filters[1][0])
	}
	if dumpedVce.Filters[1][1] != -24 {
		t.Errorf("Round trip value failed %s to %v - got %v", "Filters[0][1]", -24, dumpedVce.Filters[1][1])
	}

}
