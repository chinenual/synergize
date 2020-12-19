package main

import (
	"bytes"
	"github.com/chinenual/synergize/data"
)

// Helper routines that may find their way back into the synergize/data module.

func helperBlankVce() (vce data.VCE, err error) {
	rdr := bytes.NewReader(data.VRAM_EDATA[data.Off_VRAM_EDATA:])
	if vce, err = data.ReadVce(rdr, false); err != nil {
		return
	}
	// re-allocate the Envelopes and each env Table to allow us to control size
	//and #osc simply by writing to VOITAB and NPOINTS params
	for i := 1; i < 16; i++ {
		// make a copy of the first osc:
		vce.Envelopes = append(vce.Envelopes, vce.Envelopes[0])
	}
	for i := 0; i < 16; i++ {
		// now re-allocate each envelope to their max possible length:
		vce.Envelopes[i].AmpEnvelope.Table = make([]byte, 4 * 16)
		vce.Envelopes[i].FreqEnvelope.Table = make([]byte, 4 * 16)
	}
	return
}

func helperSetPatchType(vce *data.VCE, patchType int) {
	for i := range data.PatchTypePerOscTable[patchType-1] {
		vce.Envelopes[i].FreqEnvelope.OPTCH = data.PatchTypePerOscTable[patchType-1][i]
	}
}
