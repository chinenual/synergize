package data

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

func dumpTestBytes(path string, bytes []byte) (err error) {
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return
	}
	return
}

func AssertString(t *testing.T, b string, expected string, context string) {
	if b != expected {
		t.Errorf("expected %s '%s', got '%s'\n", context, expected, b)
	}
}

func AssertByte(t *testing.T, b byte, expected byte, context string) {
	if b != expected {
		t.Errorf("expected %s %04x, got %04x\n", context, expected, b)
	}
}

func AssertUint16(t *testing.T, b uint16, expected uint16, context string) {
	if b != expected {
		t.Errorf("expected %s %04x, got %04x\n", context, expected, b)
	}
}

func diffByte(v1, v2 byte, description string) (same bool) {
	if v1 != v2 {
		log.Printf("DIFF: %s: %02x %02x\n", description, v1, v2)
		return false
	}
	return true
}

func diffUint16(v1, v2 uint16, description string) (same bool) {
	if v1 != v2 {
		log.Printf("DIFF: %s: %04x %04x\n", description, v1, v2)
		return false
	}
	return true
}

func diffInt8(v1, v2 int8, description string) (same bool) {
	if v1 != v2 {
		log.Printf("DIFF: %s: %d %d\n", description, v1, v2)
		return false
	}
	return true
}

func diffObject(v1, v2 interface{}, description string) (same bool) {
	if !reflect.DeepEqual(v1, v2) {
		log.Printf("DIFF: %s: \nobj1: %#v\nobj2: %#v\n", description, v1, v2)
		return false
	}
	return true
}

// return true if same; false if different.  Log any diffs
func diffVCE(vce1, vce2 VCE) (same bool) {
	same = true

	// can't just DeepEqual two VCE's -- SYNCHS and the firmware leave
	// unused elements in various arrays with undefined values.  So we
	// need a custom diff function that ignores the unused values
	same = diffByte(vce1.Head.VOITAB, vce2.Head.VOITAB, "Head.VOITAB") && same
	for i := byte(0); i <= vce1.Head.VOITAB; i++ {
		same = diffUint16(vce1.Head.OSCPTR[i], vce2.Head.OSCPTR[i],
			fmt.Sprintf("Head.OSCPTR[%d]", i)) && same
	}
	same = diffInt8(vce1.Head.VTRANS, vce2.Head.VTRANS, "Head.VTRANS") && same
	same = diffByte(vce1.Head.VTCENT, vce2.Head.VTCENT, "Head.VTCENT") && same
	same = diffByte(vce1.Head.VTSENS, vce2.Head.VTSENS, "Head.VTSENS") && same
	same = diffObject(vce1.Head.VEQ, vce2.Head.VEQ, "vce1.Head.VEQ") && same
	same = diffObject(vce1.Head.VNAME, vce2.Head.VNAME, "Head.VNAME") && same
	same = diffByte(vce1.Head.VACENT, vce2.Head.VACENT, "Head.VACENT") && same
	same = diffByte(vce1.Head.VASENS, vce2.Head.VASENS, "Head.VASENS") && same
	same = diffByte(vce1.Head.VIBRAT, vce2.Head.VIBRAT, "Head.VIBRAT") && same
	same = diffByte(vce1.Head.VIBDEL, vce2.Head.VIBDEL, "Head.VIBDEL") && same
	same = diffInt8(vce1.Head.VIBDEP, vce2.Head.VIBDEP, "Head.VIBDEP") && same
	same = diffObject(vce1.Head.KPROP, vce2.Head.KPROP, "Head.KPROP") && same
	same = diffInt8(vce1.Head.APVIB, vce2.Head.APVIB, "Head.APVIB") && same
	for i := byte(0); i <= vce1.Head.VOITAB; i++ {
		same = diffInt8(vce1.Head.FILTER[i], vce2.Head.FILTER[i],
			fmt.Sprintf("Head.FILTER[%d]", i)) && same
	}
	// technically, Envelopes need to be byte for byte identical. But when
	// they differ, its nice to get field specific diffs, so we do it the
	// tedious way:
	same = diffInt8(int8(len(vce1.Envelopes)),
		int8(len(vce2.Envelopes)), "Envelopes length") && same
	for i, e1 := range vce1.Envelopes {
		e2 := vce2.Envelopes[i]

		same = diffByte(e1.FreqEnvelope.OPTCH, e2.FreqEnvelope.OPTCH,
			fmt.Sprintf("env[%d].FreqEnvelope.OPTCH", i)) && same
		same = diffInt8(e1.FreqEnvelope.OHARM, e2.FreqEnvelope.OHARM,
			fmt.Sprintf("env[%d].FreqEnvelope.OHARM", i)) && same
		same = diffInt8(e1.FreqEnvelope.FDETUN, e2.FreqEnvelope.FDETUN,
			fmt.Sprintf("env[%d].FreqEnvelope.FDETUN", i)) && same
		same = diffByte(e1.FreqEnvelope.FENVL, e2.FreqEnvelope.FENVL,
			fmt.Sprintf("env[%d].FreqEnvelope.FENVL", i)) && same
		same = diffByte(e1.FreqEnvelope.ENVTYPE, e2.FreqEnvelope.ENVTYPE,
			fmt.Sprintf("env[%d].FreqEnvelope.ENVTYPE", i)) && same
		same = diffByte(e1.FreqEnvelope.NPOINTS, e2.FreqEnvelope.NPOINTS,
			fmt.Sprintf("env[%d].FreqEnvelope.NPOINTS", i)) && same
		same = diffByte(e1.FreqEnvelope.SUSTAINPT, e2.FreqEnvelope.SUSTAINPT,
			fmt.Sprintf("env[%d].FreqEnvelope.SUSTAINPT", i)) && same
		same = diffByte(e1.FreqEnvelope.LOOPPT, e2.FreqEnvelope.LOOPPT,
			fmt.Sprintf("env[%d].FreqEnvelope.LOOPPT", i)) && same
		same = diffObject(e1.FreqEnvelope.Table, e2.FreqEnvelope.Table,
			fmt.Sprintf("env[%d].FreqEnvelope.Table", i)) && same

		same = diffByte(e1.AmpEnvelope.ENVTYPE, e2.AmpEnvelope.ENVTYPE,
			fmt.Sprintf("env[%d].AmpEnvelope.ENVTYPE", i)) && same
		same = diffByte(e1.AmpEnvelope.NPOINTS, e2.AmpEnvelope.NPOINTS,
			fmt.Sprintf("env[%d].AmpEnvelope.NPOINTS", i)) && same
		same = diffByte(e1.AmpEnvelope.SUSTAINPT, e2.AmpEnvelope.SUSTAINPT,
			fmt.Sprintf("env[%d].AmpEnvelope.SUSTAINPT", i)) && same
		same = diffByte(e1.AmpEnvelope.LOOPPT, e2.AmpEnvelope.LOOPPT,
			fmt.Sprintf("env[%d].AmpEnvelope.LOOPPT", i)) && same
		same = diffObject(e1.AmpEnvelope.Table, e2.AmpEnvelope.Table,
			fmt.Sprintf("env[%d].AmpEnvelope.Table", i)) && same
	}

	same = diffObject(vce1.Filters, vce2.Filters, "Filters") && same
	return same
}
