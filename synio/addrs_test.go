package synio

import (
	"testing"

	"github.com/chinenual/synergize/data"
)

func TestEDATAAddrs(t *testing.T) {
	// some selective fields to validate that the reflective offset
	// computation works

	// set VRAM base address to reflect the value we will get at runtime so
	// it's easy to compare values of this unit test vs. the dynamic case
	synAddrs.VRAM = 0x6000

	// expected values are derived from SYNHCS wire logging
	data.AssertUint16(t, VoiceHeadAddr(data.Off_EDATA_VOITAB),
		0x62c0, "VOITAB")
	data.AssertUint16(t, VoiceHeadAddr(data.Off_EDATA_OSCPTR),
		0x62c1, "OSCPTR")
	data.AssertUint16(t, VoiceHeadAddr(data.Off_EDATA_VTRANS),
		0x62e1, "VTRANS")
	data.AssertUint16(t, VoiceHeadAddr(data.Off_EDATA_APVIB),
		0x6322, "APVIB")

	data.AssertUint16(t, VoiceOscAddr(1, data.Off_EOSC_OPTCH),
		0x6333, "OSC[1].OPTCH")
	data.AssertUint16(t, VoiceOscAddr(1, data.Off_EOSC_OHARM),
		0x6334, "OSC[1].OHARM")
	data.AssertUint16(t, VoiceOscAddr(1, data.Off_EOSC_FDETUN),
		0x6335, "OSC[1].FDETUN")

	data.AssertUint16(t, VoiceOscAddr(2, data.Off_EOSC_OPTCH),
		0x63bf, "OSC[2].OPTCH")
	data.AssertUint16(t, VoiceOscAddr(2, data.Off_EOSC_OHARM),
		0x63c0, "OSC[2].OHARM")
	data.AssertUint16(t, VoiceOscAddr(2, data.Off_EOSC_FDETUN),
		0x63c1, "OSC[2].FDETUN")
}
