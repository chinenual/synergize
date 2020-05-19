package data

import (
	"fmt"
//	"log"
	"testing"
)

func TestLocalEDATA(t *testing.T) {
//	log.Printf("edata_head_default: %v\n", edata_head_default)
	// spot check some data:
	AssertByte(t, EDATA[0], 0, "VOITAB")
	AssertByte(t, EDATA[Off_EDATA_APVIB], 32, "APVIB")

	for osc := 1; osc <= 16; osc++ {
		off := EDATALocalOscOffset(osc, 0)
		AssertByte(t, EDATA[off], 4, fmt.Sprintf("osc %d OPTCH", osc))
		off = EDATALocalOscOffset(osc, Off_EOSC_OHARM)
		AssertByte(t, EDATA[off], 1, fmt.Sprintf("osc %d OHARM", osc))
		off = EDATALocalOscOffset(osc, Off_EOSC_FDETUN)
		AssertByte(t, EDATA[off], 0, fmt.Sprintf("osc %d FDETUN", osc))
	}
}
