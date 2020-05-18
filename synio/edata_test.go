package synio

import (
	"fmt"
	"log"
	"testing"
)


func TestEdataAddrs(t *testing.T) {
	// some selective fields to validate that the reflective offset
	// computation works

	// set EDATA to reflect the value we will get at runtime so its easy
	// to compare values of this unit test vs. the dynamic case
	synAddrs.EDATA=0x6300
	assertUint16(t, edataHeadAddr(off_EDATA_VOITAB),
		synAddrs.EDATA, "VOITAB")
	assertUint16(t, edataHeadAddr(off_EDATA_OSCPTR),
		synAddrs.EDATA+1, "OSCPTR")
	assertUint16(t, edataHeadAddr(off_EDATA_VTRANS),
		synAddrs.EDATA+33, "VTRANS")
	assertUint16(t, edataHeadAddr(off_EDATA_APVIB),
		synAddrs.EDATA+98, "APVIB")

	assertUint16(t, edataOscAddr(1,off_EOSC_OPTCH),
		synAddrs.EDATA+off_EDATA_EOSC+0, "OSC[0].OPTCH")
	assertUint16(t, edataOscAddr(1,off_EOSC_OHARM),
		synAddrs.EDATA+off_EDATA_EOSC+1, "OSC[0].OHARM")
	assertUint16(t, edataOscAddr(1,off_EOSC_FDETUN),
		synAddrs.EDATA+off_EDATA_EOSC+2, "OSC[0].FDETUN")

	assertUint16(t, edataOscAddr(2,off_EOSC_OPTCH),
		synAddrs.EDATA+off_EDATA_EOSC+sizeof_EOSC+0, "OSC[1].OPTCH")
	assertUint16(t, edataOscAddr(2,off_EOSC_OHARM),
		synAddrs.EDATA+off_EDATA_EOSC+sizeof_EOSC+1, "OSC[1].OHARM")
	assertUint16(t, edataOscAddr(2,off_EOSC_FDETUN),
		synAddrs.EDATA+off_EDATA_EOSC+sizeof_EOSC+2, "OSC[1].FDETUN")
}

func TestLocalEdata(t *testing.T) {
	log.Printf("edata_head_default: %v\n", edata_head_default)
	// spot check some data:
	assertByte(t, edata[0], 0, "VOITAB")
	assertByte(t, edata[off_EDATA_APVIB], 32, "APVIB")

	for osc := 1; osc <= 16; osc++ {
		off := edataLocalOscOffset(osc, 0)
		assertByte(t, edata[off], 4, fmt.Sprintf("osc %d OPTCH", osc))
		off = edataLocalOscOffset(osc, off_EOSC_OHARM)
		assertByte(t, edata[off], 1, fmt.Sprintf("osc %d OHARM", osc))
		off = edataLocalOscOffset(osc, off_EOSC_FDETUN)
		assertByte(t, edata[off], 0, fmt.Sprintf("osc %d FDETUN", osc))
	}
}
