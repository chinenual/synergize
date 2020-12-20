package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chinenual/synergize/data"
	"github.com/orcaman/writerseeker"
	"io/ioutil"
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

func helperCompactVCE(vce data.VCE) (compacted data.VCE, err error) {
	var writebuf = writerseeker.WriterSeeker{}

	if err = data.WriteVce(&writebuf, vce, data.VceName(vce.Head), false); err != nil {
		return
	}
	write_bytes, _ := ioutil.ReadAll(writebuf.Reader())

	var readbuf2 = bytes.NewReader(write_bytes)

	if compacted, err = data.ReadVce(readbuf2, false); err != nil {
		return
	}
	return
}

func helperVCEToJSON(vce data.VCE) (result string) {
	// compact the vce before printing it:
	var err error
	var compacted data.VCE
	if compacted,err = helperCompactVCE(vce); err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}

	b, _ := json.MarshalIndent(compacted, "", "\t")
	result = string(b)
	return
}
