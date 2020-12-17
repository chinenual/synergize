package main

import (
	"bytes"
	"github.com/chinenual/synergize/data"
)

// Helper routines that may find their way back into the synergize/data module.

func BlankVce() (vce data.VCE, err error) {
	rdr := bytes.NewReader(data.VRAM_EDATA[data.Off_VRAM_EDATA:])
	if vce, err = data.ReadVce(rdr, false); err != nil {
		return
	}
	return
}
