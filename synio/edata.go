package synio

import (
	"log"
)

// layout of the EDATA, OSC tables on the synergy.

// NOTE: can't use fancy go reflection or unsafe pointer computation -
// can't directly memory map Z80 data structures (byte aligned) with
// whatever processor we're running on (which is going to be word and
// pointer aligned).  So we have to manually compute offsets. Sigh.

// OSC layout in VOICE.Z80, OSC1:
// EDHEAD: header of the EDATA data structure

// EDHEAD is identical to data.VCEHead. Use same names

const (
	off_EDATA_VOITAB 	= 0
	off_EDATA_OSCPTR 	= 1
	off_EDATA_VTRANS 	= 33
	off_EDATA_VTCENT 	= 34
	off_EDATA_VTSENS 	= 35
	off_EDATA_patchTableIndex = 36
	off_EDATA_VEQ 		= 37
	off_EDATA_VNAME 	= 61	
	off_EDATA_VACENT 	= 69
	off_EDATA_VASENS 	= 70
	off_EDATA_VIBRAT 	= 71
	off_EDATA_VIBDEL 	= 72
	off_EDATA_VIBDEP 	= 73
	off_EDATA_KPROP		= 97
	off_EDATA_APVIB  	= 98
	off_EDATA_FILTER_arr	= 99

	// offset of OSC[0] from EDATA:
	off_EDATA_EOSC		= 115
	// size of each EOSC array element
	sizeof_EOSC = 140

	off_EOSC_OPTCH		= 0
	off_EOSC_OHARM		= 1
	off_EOSC_FDETUN		= 2
	off_EOSC_FENVL		= 3
		// frequency env:
	off_EOSC_FreqENVTYPE	= 4
	off_EOSC_FreqNPOINTS	= 5
	off_EOSC_FreqSUSTAINPT	= 6
	off_EOSC_FreqLOOPPT	= 7
	off_EOSC_FreqPoints	= 8
		// amp env:
	off_EOSC_AmpENVTYPE	= 72
	off_EOSC_AmpNPOINTS	= 73
	off_EOSC_AmpSUSTAINPT	= 74
	off_EOSC_AmpLOOPPT	= 75
	off_EOSC_AmpPoints 	= 76
	
)

func edataHeadAddr(fieldOffset int) uint16 {
	return synAddrs.EDATA + uint16(fieldOffset)
}

func edataOscAddr(osc int, fieldOffset int) uint16 {
	// osc is 1-based
	return (synAddrs.EDATA +
		uint16(off_EDATA_EOSC +
		((osc-1) * sizeof_EOSC) +
		fieldOffset))
	
}
