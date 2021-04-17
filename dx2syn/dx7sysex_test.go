package dx2syn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

func diffOsc(voiceNum, oscNum int, v, expected Dx7Osc) bool {
	result := true
	vv := reflect.ValueOf(v)
	tv := reflect.TypeOf(v)
	vexpected := reflect.ValueOf(expected)
	for i := 0; i < vv.NumField(); i++ {
		sf := tv.Field(i)
		// Osc special cased above
		fv := vv.Field(i)
		fexpected := vexpected.FieldByName(sf.Name)
		if fv.String() != fexpected.String() {
			fmt.Printf("Voice %d, Osc %d, field %s: Saw %v expected %v\n", voiceNum, oscNum, sf.Name, fv, fexpected)
			result = false
		}
	}
	return result
}

func diffVoice(voiceNum int, v, expected Dx7Voice) bool {
	result := true
	for i := 0; i < len(v.Osc); i++ {
		ok := diffOsc(voiceNum, i, v.Osc[i], expected.Osc[i])
		if !ok {
			result = false
		}
	}
	vv := reflect.ValueOf(v)
	tv := reflect.TypeOf(v)
	vexpected := reflect.ValueOf(expected)
	for i := 0; i < vv.NumField(); i++ {
		sf := tv.Field(i)
		if sf.Name != "Osc" {
			// Osc special cased above
			fv := vv.Field(i)
			fexpected := vexpected.FieldByName(sf.Name)
			if fv.String() != fexpected.String() {
				fmt.Printf("Voice %d, field %s: Saw %v expected %v\n", voiceNum, sf.Name, fv, fexpected)
				result = false
			}
		}
	}
	return result
}

func diffSysex(sysex Dx7Sysex, expected Dx7Sysex) bool {
	if len(sysex.Voices) != len(expected.Voices) {
		fmt.Printf("Saw %d voices; expected %d\n", len(sysex.Voices), len(expected.Voices))
		return false
	}
	result := true
	for i := 0; i < len(sysex.Voices); i++ {
		ok := diffVoice(i, sysex.Voices[i], expected.Voices[i])
		if !ok {
			result = false
		}
	}
	return result
}

func Test32(t *testing.T) {
	var err error
	var sysex Dx7Sysex
	if sysex, err = ReadDx7Sysex("testfiles/DX7IIFDVoice32.SYX"); err != nil {
		t.Errorf("Failed to parse SYX: %v\n", err)
		return
	}
	var b []byte
	if b, err = ioutil.ReadFile("testfiles/DX7IIFDVoice32.json"); err != nil {
		t.Errorf("Failed to read json: %v\n", err)
		return
	}
	var expected Dx7Sysex
	if err = json.Unmarshal(b, &expected); err != nil {
		t.Errorf("Failed to parse json: %v\n", err)
		return
	}
	if !diffSysex(sysex, expected) {
		t.Fatalf("Sysex is different\n")
	}
}

func Test1(t *testing.T) {

}
