package synio

import (
	"testing"
	"github.com/chinenual/synergize/data"
)


func TestEDATAAddrs(t *testing.T) {
	// some selective fields to validate that the reflective offset
	// computation works

	// set EDATA to reflect the value we will get at runtime so its easy
	// to compare values of this unit test vs. the dynamic case
	synAddrs.EDATA=0x6300
	
	data.AssertUint16(t, EDATAHeadAddr(data.Off_EDATA_VOITAB),
		synAddrs.EDATA, "VOITAB")
	data.AssertUint16(t, EDATAHeadAddr(data.Off_EDATA_OSCPTR),
		synAddrs.EDATA+1, "OSCPTR")
	data.AssertUint16(t, EDATAHeadAddr(data.Off_EDATA_VTRANS),
		synAddrs.EDATA+33, "VTRANS")
	data.AssertUint16(t, EDATAHeadAddr(data.Off_EDATA_APVIB),
		synAddrs.EDATA+98, "APVIB")

	data.AssertUint16(t, EDATAOscAddr(1,data.Off_EOSC_OPTCH),
		synAddrs.EDATA+data.Off_EDATA_EOSC+0, "OSC[0].OPTCH")
	data.AssertUint16(t, EDATAOscAddr(1,data.Off_EOSC_OHARM),
		synAddrs.EDATA+data.Off_EDATA_EOSC+1, "OSC[0].OHARM")
	data.AssertUint16(t, EDATAOscAddr(1,data.Off_EOSC_FDETUN),
		synAddrs.EDATA+data.Off_EDATA_EOSC+2, "OSC[0].FDETUN")

	data.AssertUint16(t, EDATAOscAddr(2,data.Off_EOSC_OPTCH),
		synAddrs.EDATA+data.Off_EDATA_EOSC+data.Sizeof_EOSC+0,
		"OSC[1].OPTCH")
	data.AssertUint16(t, EDATAOscAddr(2,data.Off_EOSC_OHARM),
		synAddrs.EDATA+data.Off_EDATA_EOSC+data.Sizeof_EOSC+1,
		"OSC[1].OHARM")
	data.AssertUint16(t, EDATAOscAddr(2,data.Off_EOSC_FDETUN),
		synAddrs.EDATA+data.Off_EDATA_EOSC+data.Sizeof_EOSC+2,
		"OSC[1].FDETUN")
}
