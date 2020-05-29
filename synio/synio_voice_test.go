package synio

import (
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

/***
func TestSetVEQ(t *testing.T) {
	if !*synio {
		t.Skip()
	}
	var err error
	for v := -24; v <= 6; v++ {
		if err = SetVoiceVEValue("VEQ", byte(v)); err != nil {
			t.Fatalf("Error setting Veq value %v", v)
		}
	}
	var arr []byte = []byte{0, 1, 2, 3, 4, 5, 6, 5, 4, 3, 2, 1, 0, 1, 2, 3, 4, 5, 6, 5, 4, 3, 2, 1}
	if err = SetVoiceVEQArray(arr); err != nil {
		t.Fatalf("Error setting Veq array %v", arr)
	}
}
func TestSetKPROP(t *testing.T) {
	if !*synio {
		t.Skip()
	}
	var err error
	for v := 0; v <= 32; v++ {
		if err = SetVoiceKPROPValue(0, byte(v)); err != nil {
			t.Fatalf("Error setting Veq value %v", v)
		}
	}
	var arr []byte = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}
	if err = SetVoiceKPROPArray(arr); err != nil {
		t.Fatalf("Error setting Veq array %v", arr)
	}
}
***/

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
