package synio

import (
	"testing"
)


func TestEdata(t *testing.T) {
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
