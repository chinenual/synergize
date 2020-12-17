package main

import (
	"github.com/chinenual/synergize/data"
)

func TranslateDx7ToVce(dx7Voice Dx7Voice) (vce data.VCE, err error) {
	if vce, err = BlankVce(); err != nil {
		return
	}

	for i := 0; i < 8; i++ {
		vce.Head.VNAME[i] = dx7Voice.VoiceName[i]
	}

	// ... everything else ...

	// if you need to abort, use:
	//	err = errors.New("an error message")
	//  return

	return
}