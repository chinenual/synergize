package synio

import (
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

var (
	edata_head_default = []byte {
		0, // VOITAB
		0xff,0xff, // OSCPTR[0] -- filled in dynamically since they are 16-bit values
		0xff,0xff, // OSCPTR[1]
		0xff,0xff, // OSCPTR[2]
		0xff,0xff, // OSCPTR[3]
		0xff,0xff, // OSCPTR[4]
		0xff,0xff, // OSCPTR[5]
		0xff,0xff, // OSCPTR[6]
		0xff,0xff, // OSCPTR[7]
		0xff,0xff, // OSCPTR[8]
		0xff,0xff, // OSCPTR[9]
		0xff,0xff, // OSCPTR[10]
		0xff,0xff, // OSCPTR[11]
		0xff,0xff, // OSCPTR[12]
		0xff,0xff, // OSCPTR[13]
		0xff,0xff, // OSCPTR[14]
		0xff,0xff, // OSCPTR[15]
		0, // VTRANS
		0, // VTCENT
		0, // VTSENSE
		0, // patchTableIndex
		0,0,0,0,0,0,0,0, // VEQ[]
		0,0,0,0,0,0,0,0,
		0,0,0,0,0,0,0,0,
		' ', ' ', ' ', ' ', // VNAME 
		' ', ' ', ' ', ' ',
		24, // VACENT 
		0, // VASENS 
		16, // VIBRAT 
		0, // VIBDEL 
		0, // VIBDEP 
		0,0,0,0,2,4,6,8, // KPROP
		10,12,14,16,18,20,22,24,
		26,28,30,32,32,32,32,32,
		32, // APVIB  
		0,0,0,0,0,0,0,0,  // FILTER
		0,0,0,0,0,0,0,0 }

	edata_osc_default = []byte {
		4, // OPTCH
		1, // OHARM     
		0, // FDETUN    
		68, // FENVL
		
		// frequency envelope:
		1, // ENVTYPE   
		1, // NPOINTS   
		30, // SUSTAINPT 
		30, // LOOPPT    
		0,0,0x80,0, 	// point1
		0,0,0,0, 	// point2
		0,0,0,0, 	// point3
		0,0,0,0, 	// point4
		0,0,0,0, 	// point5
		0,0,0,0, 	// point6
		0,0,0,0, 	// point7
		0,0,0,0, 	// point8
		0,0,0,0, 	// point9
		0,0,0,0, 	// point10
		0,0,0,0, 	// point11
		0,0,0,0, 	// point12
		0,0,0,0, 	// point13
		0,0,0,0, 	// point14
		0,0,0,0, 	// point15
		0,0,0,0, 	// point16
		
		// amplitude envelope:
		1, // ENVTYPE
		1, // NPOINTS   
		30, // SUSTAINPT 
		30, // LOOPPT    
		55,55,0,0, 	// point1
		55,55,0,0, 	// point2
		55,55,0,0, 	// point3
		55,55,0,0, 	// point4
		55,55,0,0, 	// point5
		55,55,0,0, 	// point6
		55,55,0,0, 	// point7
		55,55,0,0, 	// point8
		55,55,0,0, 	// point9
		55,55,0,0, 	// point10
		55,55,0,0, 	// point11
		55,55,0,0, 	// point12
		55,55,0,0, 	// point13
		55,55,0,0, 	// point14
		55,55,0,0, 	// point15
		55,55,0,0 }	// point16	

	edata [off_EDATA_EOSC + 16 * sizeof_EOSC]byte
)

func init() {
	// fixup the default data that can't be easily made with literal values
	for osc := 0; osc < 16; osc++ {
		offset := osc * sizeof_EOSC + off_EDATA_EOSC
		hob,lob := wordToBytes(uint16(offset))
		edata_head_default[1 + osc*2] = lob
		edata_head_default[2 + osc*2] = hob
	}
	ClearLocalEDATA()
}

func ClearLocalEDATA() {
	for i := 0; i < len(edata_head_default); i++ {
		edata[i] = edata_head_default[i]
	}
	for osc := 0; osc < 16; osc++ {
		for i := 0; i < len(edata_osc_default); i++ {
			offset := off_EDATA_EOSC + osc * sizeof_EOSC + i
			edata[offset] = edata_osc_default[i]
		}
	}
}

func edataLocalHeadOffset(fieldOffset int) uint16 {
	return uint16(fieldOffset)
}

func edataLocalOscOffset(osc int, fieldOffset int) uint16 {
	// osc is 1-based
	return uint16(off_EDATA_EOSC +
		((osc-1) * sizeof_EOSC) +
		fieldOffset)
}

func edataHeadAddr(fieldOffset int) uint16 {
	return synAddrs.EDATA + edataLocalHeadOffset(fieldOffset)
}


func edataOscAddr(osc int, fieldOffset int) uint16 {
	// osc is 1-based
	return synAddrs.EDATA + edataLocalOscOffset(osc, fieldOffset)	
}
