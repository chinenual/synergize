package data

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/orcaman/writerseeker"

	"github.com/chinenual/synergize/logger"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

// XREF: needs to match the array in viewVCE.html
// FIXME: see viewVCE_voices.patchTable() -- for commented out code to dynamically update
// the options array from this string.
var PatchTypeNames = []string{
	"1 + 2 + 3 + 4 + 5 + 6 + 7 + 8",
	"(1~2) + (3~4) + (5~6) + (7~8)",
	"((1+2+3)~4) + ((5+6+7)~8)",
	"(1~2+3)~4) + ((5~6+7)~8)",
	"(1~2) + 3 + 4 + (5~6) + 7 + 8",
	"((1+2)~3) + ((1+2)~4) + ((5+6)~7) + ((5+6)~8)",
	"(1~2) + (1~3) + (1~4) + (1~5) + (1~6) + (1~7) + (1~8)",
	"(1~2~3) + (1~2~4) + (1~2~5) + (1~2~6) + (1~2~7) + (1~2~8)",
	"(1~2~3~4) + (1~2~3~5) + (1~2~3~6) + (1~2~3~7) + (1~2~3~8)",
	"((1~2+3)~4) + ((1~2+3)~5) + ((1~2+3)~6) + ((1~2+3)~7) + ((1~2+3)~8)",

	// DX patches. Some share same register layout - see PatchTypePerOscTable
	"(1~2~3~4) + (5~6)", // DX 1 & 2
	"(1~2~3) + (4~5~6)", // DX 3 & 4
	//                      DX 5 & 6 are same as Syn #2
	"((1~2)+3)~4) + (5~6)",  // DX 7 & 8 & 9
	"((1+2)~3) + (4~5~6)",   // DX 10 & 11
	"((1+2+3)~4) + (5~6)",   // DX 12 & 13
	"(((1+2)~3)~4) + (5~6)", // DX 14 & 15

	"((1~2)+(3~4)+5)~6", // DX 16 & 17
	"((1~2)+3+4)~5",     // DX 16 & 17 - 5-osc variant

	"((1~2~3)+4+5)~6",               // DX 18
	"(1~2) + (1~3) + (4~5~6)",       // DX 19
	"(1+2)~3 + ((4~5) + (4~6))",     // DX 20
	"(1~2) + (1~3) + (4~5) + (4~6)", // DX 21
	"(1~2) + (1~3) + (1~4) + (5~6)", // DX 22
	"(1~2) + (1~3) + (4~5) + 6",     // DX 23
	"(1~2) + (1~3) + (1~4) + 5 + 6", // DX 24
	"(1~2) + (1~3) + 4 + 5 + 6",     // DX 25
	"(1~3) + (2~3) + (4~5) + 6",     // DX 26 & 27
	"1 + (2~3~4) + (5~6)",           // DX 28
	"(1~2) + (3~4) + 5 + 6",         // DX 29
	"1 + (2~3~4) + 5 + 6",           // DX 30
	"(1~2) + 3 + 4 + 5 + 6",         // DX 31
}

// the first 8 oscillators OPTCH values use this pattern based on which patch type
// is chosen.   The second 8 oscillators are always '4' (aka additive).
// From SYNHCS.Z80, PTCHTB:
var PatchTypePerOscTable = [][16]byte{
	{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // SYN type 1, also DX 32
	{100, 1, 100, 1, 100, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4},
	{100, 76, 76, 1, 100, 76, 76, 1, 4, 4, 4, 4, 4, 4, 4, 4},
	{100, 97, 76, 1, 100, 97, 76, 1, 4, 4, 4, 4, 4, 4, 4, 4},
	{100, 1, 4, 4, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},
	{100, 76, 1, 1, 100, 76, 1, 1, 4, 4, 4, 4, 4, 4, 4, 4},
	{100, 1, 1, 1, 1, 1, 1, 1, 4, 4, 4, 4, 4, 4, 4, 4},
	{100, 97, 1, 1, 1, 1, 1, 1, 4, 4, 4, 4, 4, 4, 4, 4},
	{100, 97, 97, 1, 1, 1, 1, 1, 4, 4, 4, 4, 4, 4, 4, 4},
	{100, 97, 76, 1, 1, 1, 1, 1, 4, 4, 4, 4, 4, 4, 4, 4},

	// XREF: keep in sync with dx2syn/algo.go: dxAlgoNoFeedbackPatchTypePerOscTable
	{100, 97, 97, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 1 & 2
	{100, 97, 1, 100, 97, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 3 & 4
	//                                                         DX 5 & 6 same as SYN #2
	{100, 97, 76, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 7 & 8 & 9
	{100, 76, 1, 100, 97, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 10 & 11
	{100, 76, 76, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 12 & 13
	{100, 76, 97, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 14 & 15

	{100, 97, 100, 73, 76, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 16 & 17 // can't reproduce the algo due to lack of register in SYN - this is approximation
	{100, 97, 76, 76, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 16 & 17 // 5-osc variant

	{100, 97, 97, 76, 76, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 18
	{100, 1, 1, 100, 97, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},  // DX 19
	{100, 76, 1, 100, 1, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},  // DX 20
	{100, 1, 1, 100, 1, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 21
	{100, 1, 1, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 22
	{100, 1, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 23
	{100, 1, 1, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},     // DX 24
	{100, 1, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},     // DX 25
	{100, 76, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},  // DX 26 & 27
	{4, 100, 97, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},  // DX 28
	{100, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 29
	{4, 100, 97, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},    // DX 30
	{100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},     // DX 31
}

func inferPatchType(vce VCE) (patchType int) {
	for patchType = 0; patchType < len(PatchTypePerOscTable); patchType++ {
		var allmatch = true
		for osc := 0; osc <= int(vce.Head.VOITAB) && osc < 8; osc++ {
			//log.Printf("patchType %d, osc %d, OPTCH %d, tbl %d\n", patchType, osc, vce.Envelopes[osc].FreqEnvelope.OPTCH, PatchTypePerOscTable[patchType][osc])
			if vce.Envelopes[osc].FreqEnvelope.OPTCH != PatchTypePerOscTable[patchType][osc] {
				// not a match
				allmatch = false
				break
			}
		}
		if allmatch {
			return patchType
		}
	}
	return patchType
}

type VCEExtra struct {
	// extra stuff we can compute from the raw voice data and want to have
	// available in the UI, but which is not serialized to/from the byte stream

	PatchType int    // inferred based on each oscilators OPTCH byte
	Doc       string // contents of the .DOC or .txt file (if any)
}

type FreqEnvelopeTable struct {
	OPTCH     byte
	OHARM     int8 `validate:"min=-12,max=30"`
	FDETUN    int8 `validate:"min=-63,max=63"`
	FENVL     byte
	ENVTYPE   byte `validate:"min=1,max=4"`
	NPOINTS   byte
	SUSTAINPT byte
	LOOPPT    byte
	Table     ArrayOfByte // force proper JSON encoding
	// entries are quads:  valLow, valUp, timeLow, timeUp
	// valid ranges for val: -61 .. 63
	// valid ranges for time: 0 .. 84
	// the first quad is an exception: timeLow and timeUp are used for keyprop/waveform
}

type AmpEnvelopeTable struct {
	ENVTYPE   byte `validate:"min=1,max=4"`
	NPOINTS   byte
	SUSTAINPT byte
	LOOPPT    byte
	Table     ArrayOfByte // force proper JSON encoding
	// entries are quads:  valLow, valUp, timeLow, timeUp
	// valid ranges for val: 55 .. 127
	// valid ranges for time: 0 .. 84
}

type Envelope struct {
	FreqEnvelope FreqEnvelopeTable
	AmpEnvelope  AmpEnvelopeTable
}

type VCE struct {
	Head      VCEHead
	Envelopes []Envelope
	Filters   [][32]int8 // values are -64..63 - validated by the extraValidation function rather than the tag
	Extra     VCEExtra
}

type VCEHead struct {
	VOITAB byte
	OSCPTR [16]uint16
	VTRANS int8               `validate:"min=-127,max=128"`
	VTCENT byte               `validate:"min=0,max=32"`
	VTSENS byte               `validate:"min=0,max=31"`
	UNUSED byte               //
	VEQ    [24]int8           `validate:"dive,min=-24,max=7"`
	VNAME  SpaceEncodedString // force string encoding for the name
	VACENT byte               `validate:"min=0,max=32"`
	VASENS byte               `validate:"min=0,max=31"`
	VIBRAT byte               `validate:"min=0,max=127"`
	VIBDEL byte               `validate:"min=0,max=127"`
	VIBDEP int8               `validate:"min=-128,max=128"`
	KPROP  [24]byte           `validate:"dive,min=0,max=32"`
	APVIB  int8               `validate:"min=-128,max=128"`
	FILTER [16]int8
}

func extraVceValidation(sl validator.StructLevel) {

	vce := sl.Current().Interface().(VCE)
	hasAFilter := false
	for osc := byte(0); osc <= vce.Head.VOITAB; osc++ {
		flt := vce.Head.FILTER[osc]
		if flt < 0 {
			hasAFilter = true
		}
	}
	for osc := byte(0); osc <= vce.Head.VOITAB; osc++ {
		flt := vce.Head.FILTER[osc]
		if flt < -1 || flt > int8(vce.Head.VOITAB+1) {
			sl.ReportError(vce.Head, fmt.Sprintf("FILTER[%d]", osc), "filterId", "", "")
		} else {
			// if this osc has a filter, determine its index into the Filters array and then check all the values
			findex := -1
			if flt == 0 {
				// skip
			} else if flt < 0 {
				findex = 0
			} else if hasAFilter {
				findex = int(flt)
			} else {
				findex = int(flt) - 1
			}
			if findex >= 0 {
				// validate all the entries:
				for i, f := range vce.Filters[findex] {
					if f < -64 || f > 63 {
						sl.ReportError(vce, fmt.Sprintf("Filter[%d][%d]", findex, i), "filterValue", "", strconv.Itoa(int(f)))
					}
				}
			}
		}
		// now check the osc's envelopes:
		if vce.Envelopes[osc].AmpEnvelope.ENVTYPE == 1 {
			if vce.Envelopes[osc].AmpEnvelope.LOOPPT > 127 {
				sl.ReportError(vce.Envelopes[osc].AmpEnvelope, "LOOPPT", "typeOneAccel", "", strconv.Itoa(int(vce.Envelopes[osc].AmpEnvelope.LOOPPT)))
			}
			if vce.Envelopes[osc].AmpEnvelope.SUSTAINPT > 127 {
				sl.ReportError(vce.Envelopes[osc].AmpEnvelope, "SUSTAINPT", "typeOneAccel", "", strconv.Itoa(int(vce.Envelopes[osc].AmpEnvelope.SUSTAINPT)))
			}
		} else {
			// TODO: add checks that REPEAT and SUSTAIN and LOOP points are in the right order
		}
		for i := byte(0); i < vce.Envelopes[osc].AmpEnvelope.NPOINTS; i++ {
			valLow := vce.Envelopes[osc].AmpEnvelope.Table[(4*i)+0]
			valUp := vce.Envelopes[osc].AmpEnvelope.Table[(4*i)+1]
			timeLow := vce.Envelopes[osc].AmpEnvelope.Table[(4*i)+2]
			timeUp := vce.Envelopes[osc].AmpEnvelope.Table[(4*i)+3]
			if valLow < 55 || valLow > 127 {
				sl.ReportError(vce.Envelopes[osc].AmpEnvelope.Table[(4*i)+0], fmt.Sprintf("vce.Envelopes[%d].AmpEnvelope.<point>[%d].valLow", osc, i), "ampvalue", strconv.Itoa(int(valLow)), "")
			}
			if valUp < 55 || valUp > 127 {
				sl.ReportError(vce.Envelopes[osc].AmpEnvelope.Table[(4*i)+0], fmt.Sprintf("vce.Envelopes[%d].AmpEnvelope.<point>[%d].valUp", osc, i), "ampvalue", strconv.Itoa(int(valUp)), "")
			}
			if timeLow > 84 {
				sl.ReportError(vce.Envelopes[osc].AmpEnvelope.Table[(4*i)+2], fmt.Sprintf("vce.Envelopes[%d].AmpEnvelope.<point>[%d].timeLow", osc, i), "amptime", strconv.Itoa(int(timeLow)), "")
			}
			if timeUp > 84 {
				sl.ReportError(vce.Envelopes[osc].AmpEnvelope.Table[(4*i)+3], fmt.Sprintf("vce.Envelopes[%d].AmpEnvelope.<point>[%d].timeUp", osc, i), "amptime", strconv.Itoa(int(timeUp)), "")
			}
		}
		if vce.Envelopes[osc].FreqEnvelope.ENVTYPE == 1 {
			if vce.Envelopes[osc].FreqEnvelope.LOOPPT > 127 {
				sl.ReportError(vce.Envelopes[osc].FreqEnvelope, "LOOPPT", "typeOneAccel", "", strconv.Itoa(int(vce.Envelopes[osc].FreqEnvelope.LOOPPT)))
			}
			if vce.Envelopes[osc].FreqEnvelope.SUSTAINPT > 127 {
				sl.ReportError(vce.Envelopes[osc].FreqEnvelope, "SUSTAINPT", "typeOneAccel", "", strconv.Itoa(int(vce.Envelopes[osc].FreqEnvelope.SUSTAINPT)))
			}
		} else {
			// TODO: add checks that REPEAT and SUSTAIN and LOOP points are in the right order
		}
		for i := byte(0); i < vce.Envelopes[osc].FreqEnvelope.NPOINTS; i++ {
			timeLow := vce.Envelopes[osc].FreqEnvelope.Table[(4*i)+2]
			timeUp := vce.Envelopes[osc].FreqEnvelope.Table[(4*i)+3]
			// freq vals use the full range -127..127 so nothing to validate
			if i != 0 {
				// TODO: add validation for the wave/keyprop bytes?
				if timeLow > 84 {
					sl.ReportError(vce.Envelopes[osc].FreqEnvelope.Table[(4*i)+2], fmt.Sprintf("vce.Envelopes[%d].FreqEnvelope.<point>[%d].timeLow", osc, i), "freqtime", strconv.Itoa(int(timeLow)), "")
				}
				if timeUp > 84 {
					sl.ReportError(vce.Envelopes[osc].FreqEnvelope.Table[(4*i)+3], fmt.Sprintf("vce.Envelopes[%d].FreqEnvelope.<point>[%d].timeUp", osc, i), "freqtime", strconv.Itoa(int(timeUp)), "")
				}
			}
		}
	}
}

func VceValidate(vce VCE) (err error) {
	// simple field level validation first:
	v := validator.New()
	v.RegisterStructValidation(extraVceValidation, VCE{})

	err = v.Struct(vce)
	return
}

func VceName(vceHead VCEHead) (name string) {
	name = ""
	for i := 0; i < 8; i++ {
		name = name + string(vceHead.VNAME[i])
	}
	name = strings.Trim(name, " ")
	return
}

func VceAFilterCount(vce VCE) (count int) {
	for i := byte(0); i <= vce.Head.VOITAB; i++ {
		v := vce.Head.FILTER[i]
		if v < 0 {
			return 1
		}
	}
	return 0
}

func VceBFilterCount(vce VCE) (count int) {
	count = 0
	for i := byte(0); i <= vce.Head.VOITAB; i++ {
		v := vce.Head.FILTER[i]
		if v > 0 {
			count = count + 1
		}
	}
	return
}

func ReadVceFile(filename string) (vce VCE, err error) {
	var b []byte

	if b, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	buf := bytes.NewReader(b)
	if vce, err = ReadVce(buf, false); err != nil {
		return
	}
	vce.Extra.Doc = readDocFile(filename)
	return
}

func readDocFile(vcePath string) (result string) {
	vceExt := path.Ext(vcePath)
	base := (vcePath)[0 : len(vcePath)-len(vceExt)]
	// tolerate expected case differences on case-sensitive file systems:
	extensions := []string{".DOC", ".doc", ".Doc"}
	for _, ext := range extensions {
		pathname := base + ext
		if b, err := ioutil.ReadFile(pathname); err == nil {
			result = string(b)
			return
		}
	}
	result = ""
	return
}

func vceNameFromPathname(filename string) (name string) {
	// filepath.Base() is OS dependent = looks for \ on windows, else /.
	// We want to be platform independent (at least for the unit tests :)) so roll our own
	name = filename
	if i := strings.LastIndexByte(name, '/'); i >= 0 {
		name = name[i+1:]
	}
	if i := strings.LastIndexByte(name, '\\'); i >= 0 {
		name = name[i+1:]
	}

	// remove .vce suffix
	name = strings.ToUpper(name)
	name = strings.TrimSuffix(name, ".VCE")
	name = VcePaddedName(name)
	return name
}

func VcePaddedName(name string) (padded string) {
	padded = name
	if len(padded) > 8 {
		padded = padded[:8]
	} else if len(padded) < 8 {
		add := 8 - len(padded)
		for i := 0; i < add; i++ {
			padded = padded + " "
		}
	}
	return padded
}

func WriteVceFile(filename string, vce VCE, overrideVNAME bool) (err error) {
	var file *os.File
	if file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0755); err != nil {
		return
	}
	defer file.Close()

	name := VceName(vce.Head)
	if overrideVNAME {
		name = vceNameFromPathname(filename)
	}

	if err = WriteVce(file, vce, name, false); err != nil {
		return
	}
	return
}

func VceReadAFilters(buf io.Reader, vce *VCE) (err error) {
	// a voice has at most 1 A filter. These are at the head of the filter array
	// so we can unconditionally put it in slot 0 if there is one
	for i := byte(0); i <= vce.Head.VOITAB; i++ {
		if vce.Head.FILTER[i] < 0 {
			for j := 0; j < 32; j++ {
				if err = binary.Read(buf, binary.LittleEndian, &vce.Filters[0][j]); err != nil {
					err = errors.Wrapf(err, "Failed to read A filter[%d]", j)
					return
				}
			}
			return
		}
	}
	return
}

func VceWriteAFilters(buf io.Writer, vce VCE) (err error) {
	for i := byte(0); i <= vce.Head.VOITAB; i++ {
		if vce.Head.FILTER[i] < 0 {
			if verboseWriting {
				logger.Infof("WRITE A filter %d %v\n", i, vce.Filters[0])
			}
			for j := 0; j < 32; j++ {
				if err = binary.Write(buf, binary.LittleEndian, vce.Filters[0][j]); err != nil {
					err = errors.Wrapf(err, "Failed to write A filter[%d]", j)
					return
				}
			}
			return
		}
	}
	return
}

func VceReadBFilters(buf io.Reader, vce *VCE) (err error) {
	var filterCount = 0
	var hasAFilter = false
	for i := byte(0); i <= vce.Head.VOITAB; i++ {
		f := vce.Head.FILTER[i]
		if f != 0 {
			filterCount = filterCount + 1
		}
		if f < 0 {
			hasAFilter = true
		}
	}

	var offset = 0
	if hasAFilter {
		offset = 1
	}
	for i := byte(0); i <= vce.Head.VOITAB; i++ {
		f := vce.Head.FILTER[i]
		if f > 0 {
			// filters are one-based
			var index = int(f) - 1 + offset
			for j := 0; j < 32; j++ {
				if err = binary.Read(buf, binary.LittleEndian, &vce.Filters[index][j]); err != nil {
					err = errors.Wrapf(err, "Failed to read B filter[%d][%d]", index, j)
					return
				}
			}
		}
	}
	return
}
func VceWriteBFilters(buf io.Writer, vce VCE) (err error) {
	var filterCount = 0
	var hasAFilter = false
	for i := byte(0); i <= vce.Head.VOITAB; i++ {
		f := vce.Head.FILTER[i]
		if f != 0 {
			filterCount = filterCount + 1
		}
		if f < 0 {
			hasAFilter = true
		}
	}

	var offset = 0
	if hasAFilter {
		offset = 1
	}
	for i := byte(0); i <= vce.Head.VOITAB; i++ {
		f := vce.Head.FILTER[i]
		if f > 0 {
			// filters are one-based
			var index = int(f) - 1 + offset
			if verboseWriting {
				logger.Infof("WRITE B filter %d (index: %d)\n", f, index)
			}
			for j := 0; j < 32; j++ {
				if err = binary.Write(buf, binary.LittleEndian, vce.Filters[index][j]); err != nil {
					err = errors.Wrapf(err, "Failed to write B filter[%d][%d]", index, j)
					return
				}
			}
		}
	}
	return
}

func ReadVce(buf io.ReadSeeker, skipFilters bool) (vce VCE, err error) {
	var startVceOffset int64
	if startVceOffset, err = buf.Seek(0, io.SeekCurrent); err != nil {
		return
	}

	if err = binary.Read(buf, binary.LittleEndian, &vce.Head); err != nil {
		err = errors.Wrapf(err, "Failed to read voice header")
		return
	}

	if verboseParsing {
		logger.Infof("voice head: %s\n", vceHeadToJson(vce.Head))
	}

	vce.Envelopes = make([]Envelope, vce.Head.VOITAB+1)
	for i := byte(0); i <= vce.Head.VOITAB; i++ {
		var e Envelope

		offset := int64(vce.Head.OSCPTR[i]) + startVceOffset

		if _, err = buf.Seek(offset, io.SeekStart); err != nil {
			err = errors.Wrapf(err, "failed to seek to osc #%d start: %04x", i, offset)
			return
		}

		if err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.OPTCH); err != nil {
			err = errors.Wrapf(err, "Failed to read voice freq env [%d] OPTCH", i)
			return
		}
		if err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.OHARM); err != nil {
			err = errors.Wrapf(err, "Failed to read voice freq env [%d] OHARM", i)
			return
		}
		if err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.FDETUN); err != nil {
			err = errors.Wrapf(err, "Failed to read voice freq env [%d] FDETUN", i)
			return
		}
		if err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.FENVL); err != nil {
			err = errors.Wrapf(err, "Failed to read voice freq env [%d] FENVL", i)
			return
		}

		// XREF: amp start
		// Danger Will Robinson. Data coming directly from the Synergy
		// via VRAM DUMP will have the envelopes spaced out with unused
		// points in between.  Need to trust the FENVL value and not
		// just use the NPOINTS value to determine where to find the
		// beginning of the Amp envelope.  So determine where that is now:
		var startAmpEnv int64
		if startAmpEnv, err = buf.Seek(0, io.SeekCurrent); err != nil {
			return
		}
		startAmpEnv = startAmpEnv + int64(e.FreqEnvelope.FENVL)

		if err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.ENVTYPE); err != nil {
			err = errors.Wrapf(err, "Failed to read voice freq env [%d] ENVTYPE", i)
			return
		}
		if err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.NPOINTS); err != nil {
			err = errors.Wrapf(err, "Failed to read voice freq env [%d] NPOINTS", i)
			return
		}
		if err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.SUSTAINPT); err != nil {
			err = errors.Wrapf(err, "Failed to read voice freq env [%d] SUSTAINPT", i)
			return
		}
		if err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.LOOPPT); err != nil {
			err = errors.Wrapf(err, "Failed to read voice freq env [%d] LOOPPT", i)
			return
		}

		// 4 values per point:

		e.FreqEnvelope.Table = make([]byte, e.FreqEnvelope.NPOINTS*4)
		for k := byte(0); k < e.FreqEnvelope.NPOINTS*4; k++ {
			if err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.Table[k]); err != nil {
				err = errors.Wrapf(err, "Failed to read voice freq env [%d] table[%d]", i, k)
				return
			}
		}

		// XREF: amp start
		if _, err = buf.Seek(startAmpEnv, io.SeekStart); err != nil {
			err = errors.Wrapf(err, "failed to seek to amp env #%d start: %04x", i, startAmpEnv)
			return
		}

		if err = binary.Read(buf, binary.LittleEndian, &e.AmpEnvelope.ENVTYPE); err != nil {
			err = errors.Wrapf(err, "Failed to read voice amp env [%d] ENVTYPE", i)
			return
		}
		if err = binary.Read(buf, binary.LittleEndian, &e.AmpEnvelope.NPOINTS); err != nil {
			err = errors.Wrapf(err, "Failed to read voice amp env [%d] NPOINTS", i)
			return
		}
		if err = binary.Read(buf, binary.LittleEndian, &e.AmpEnvelope.SUSTAINPT); err != nil {
			err = errors.Wrapf(err, "Failed to read voice amp env [%d] SUSTAINPT", i)
			return
		}
		if err = binary.Read(buf, binary.LittleEndian, &e.AmpEnvelope.LOOPPT); err != nil {
			err = errors.Wrapf(err, "Failed to read voice amp env [%d] LOOPPT", i)
			return
		}
		// 4 values per point:
		e.AmpEnvelope.Table = make([]byte, e.AmpEnvelope.NPOINTS*4)
		for k := byte(0); k < e.AmpEnvelope.NPOINTS*4; k++ {
			if err = binary.Read(buf, binary.LittleEndian, &e.AmpEnvelope.Table[k]); err != nil {
				err = errors.Wrapf(err, "Failed to read amp freq env [%d] table[%d]", i, k)
				return
			}
		}
		vce.Envelopes[i] = e
	}

	vce.Extra.PatchType = inferPatchType(vce)

	var filterCount = 0
	var hasAFilter = false
	for _, f := range vce.Head.FILTER {
		if f > 0 {
			filterCount++
		} else if f < 0 {
			hasAFilter = true
		}
	}
	if hasAFilter {
		filterCount = filterCount + 1
	}
	vce.Filters = make([][32]int8, filterCount)

	if !skipFilters {
		if err = VceReadAFilters(buf, &vce); err != nil {
			err = errors.Wrapf(err, "Failed to read A filters")
			return
		}
		if err = VceReadBFilters(buf, &vce); err != nil {
			err = errors.Wrapf(err, "Failed to read B filters")
			return
		}
	}
	return
}

func updateOscPtr(buf io.WriteSeeker, headOffset int64, osc byte, val uint16) (err error) {
	var curr int64
	if curr, err = buf.Seek(0, io.SeekCurrent); err != nil {
		return
	}
	// OSC[n] offset is Head+1 + (osc)*2  [2-byte word]
	if _, err = buf.Seek(headOffset+1+(int64(osc)*2), io.SeekStart); err != nil {
		return
	}

	lob, hob := WordToBytes(val)

	if err = binary.Write(buf, binary.LittleEndian, hob); err != nil {
		err = errors.Wrapf(err, "Failed to write OSCPTR[%d].hob", osc)
		return
	}
	if err = binary.Write(buf, binary.LittleEndian, lob); err != nil {
		err = errors.Wrapf(err, "Failed to write OSCPTR[%d].lob", osc)
		return
	}

	// restore the position
	if _, err = buf.Seek(curr, io.SeekStart); err != nil {
		return
	}
	return
}

// sets the VNAME, rewrites the OSCPTR array and compresses the FENVL values
func WriteVce(buf io.WriteSeeker, vce VCE, name string, skipFilters bool) (err error) {
	if err = writeVce(buf, vce, name, skipFilters, false); err != nil {
		return
	}
	return
}

// sets the VNAME, but does not overwrite OSCPTR array or compress envelopes.
func WriteVcePreserveOffsets(buf io.WriteSeeker, vce VCE, name string, skipFilters bool) (err error) {
	if err = writeVce(buf, vce, name, skipFilters, true); err != nil {
		return
	}
	return
}

func VceWriteOscillator(buf io.WriteSeeker, e Envelope, osc /*1-based*/ byte, preserveOffsets bool) (err error) {
	var oscOffset int64
	if oscOffset, err = buf.Seek(0, io.SeekCurrent); err != nil {
		return
	}
	if err = binary.Write(buf, binary.LittleEndian, e.FreqEnvelope.OPTCH); err != nil {
		err = errors.Wrapf(err, "Failed to write voice freq env [%d] OPTCH", osc)
		return
	}
	if err = binary.Write(buf, binary.LittleEndian, e.FreqEnvelope.OHARM); err != nil {
		err = errors.Wrapf(err, "Failed to write voice freq env [%d] OHARM", osc)
		return
	}
	if err = binary.Write(buf, binary.LittleEndian, e.FreqEnvelope.FDETUN); err != nil {
		err = errors.Wrapf(err, "Failed to write voice freq env [%d] FDETUN", osc)
		return
	}

	if !preserveOffsets {
		e.FreqEnvelope.FENVL = 4 + 4*e.FreqEnvelope.NPOINTS
	}
	if err = binary.Write(buf, binary.LittleEndian, e.FreqEnvelope.FENVL); err != nil {
		err = errors.Wrapf(err, "Failed to write voice freq env [%d] FENVL", osc)
		return
	}
	if err = binary.Write(buf, binary.LittleEndian, e.FreqEnvelope.ENVTYPE); err != nil {
		err = errors.Wrapf(err, "Failed to write voice freq env [%d] ENVTYPE", osc)
		return
	}
	if err = binary.Write(buf, binary.LittleEndian, e.FreqEnvelope.NPOINTS); err != nil {
		err = errors.Wrapf(err, "Failed to write voice freq env [%d] NPOINTS", osc)
		return
	}
	if err = binary.Write(buf, binary.LittleEndian, e.FreqEnvelope.SUSTAINPT); err != nil {
		err = errors.Wrapf(err, "Failed to write voice freq env [%d] SUSTAINPT", osc)
		return
	}
	if err = binary.Write(buf, binary.LittleEndian, e.FreqEnvelope.LOOPPT); err != nil {
		err = errors.Wrapf(err, "Failed to write voice freq env [%d] LOOPPT", osc)
		return
	}
	// 4 values per point:

	for k := byte(0); k < e.FreqEnvelope.NPOINTS*4; k++ {
		if err = binary.Write(buf, binary.LittleEndian, e.FreqEnvelope.Table[k]); err != nil {
			err = errors.Wrapf(err, "Failed to write freq env [%d] table[%d]", osc, k)
			return
		}
	}

	if preserveOffsets {
		// skip over the unused freq envelope entries
		if _, err = buf.Seek(oscOffset+72, io.SeekStart); err != nil {
			return
		}
	}
	if err = binary.Write(buf, binary.LittleEndian, e.AmpEnvelope.ENVTYPE); err != nil {
		err = errors.Wrapf(err, "Failed to write voice amp env [%d] ENVTYPE", osc)
		return
	}
	if err = binary.Write(buf, binary.LittleEndian, e.AmpEnvelope.NPOINTS); err != nil {
		err = errors.Wrapf(err, "Failed to write voice amp env [%d] NPOINTS", osc)
		return
	}
	if err = binary.Write(buf, binary.LittleEndian, e.AmpEnvelope.SUSTAINPT); err != nil {
		err = errors.Wrapf(err, "Failed to write voice amp env [%d] SUSTAINPT", osc)
		return
	}
	if err = binary.Write(buf, binary.LittleEndian, e.AmpEnvelope.LOOPPT); err != nil {
		err = errors.Wrapf(err, "Failed to write voice amp env [%d] LOOPPT", osc)
		return
	}
	// 4 values per point:
	for k := byte(0); k < e.AmpEnvelope.NPOINTS*4; k++ {
		if err = binary.Write(buf, binary.LittleEndian, e.AmpEnvelope.Table[k]); err != nil {
			err = errors.Wrapf(err, "Failed to write amp freq env [%d] table[%d]", osc, k)
			return
		}
	}
	return
}

func stringToSpaceEncodedString(s string) (u SpaceEncodedString) {
	for i := range u {
		if i < len(s) {
			u[i] = s[i]
		} else {
			u[i] = ' '
		}
	}
	return u
}

// sets the VNAME, rewrites the OSCPTR array and compresses the FENVL values
func writeVce(buf io.WriteSeeker, vce VCE, name string, skipFilters bool, preserveOffsets bool) (err error) {
	var headOffset int64
	if headOffset, err = buf.Seek(0, io.SeekCurrent); err != nil {
		return
	}
	if verboseWriting {
		logger.Infof("SEEK - top of voice at 0x%04x", headOffset)
	}
	vce.Head.VNAME = stringToSpaceEncodedString(name)
	if !preserveOffsets {
		for i := range vce.Head.OSCPTR {
			vce.Head.OSCPTR[i] = 0
		}
	}

	if err = binary.Write(buf, binary.LittleEndian, &vce.Head); err != nil {
		err = errors.Wrapf(err, "Failed to write voice header")
		return
	}

	for i := byte(0); i <= vce.Head.VOITAB; i++ {
		e := vce.Envelopes[i]

		var oscOffset int64
		if preserveOffsets {
			// trust the OSCPTR value
			oscOffset = headOffset + int64(vce.Head.OSCPTR[i])
			if _, err = buf.Seek(oscOffset, io.SeekStart); err != nil {
				return
			}
			//logger.Infof("SEEK - top of osc[%d] at 0x%04x", i, headOffset+int64(vce.Head.OSCPTR[i]))
		} else {
			if oscOffset, err = buf.Seek(0, io.SeekCurrent); err != nil {
				return
			}
			// fixup the oscptr based on where we are in the bytestream
			if err = updateOscPtr(buf, headOffset, i, uint16(oscOffset-headOffset)); err != nil {
				return
			}
		}

		if err = VceWriteOscillator(buf, e, i, preserveOffsets); err != nil {
			return
		}
	}

	if !skipFilters {
		if err = VceWriteAFilters(buf, vce); err != nil {
			err = errors.Wrapf(err, "Failed to write A filters")
			return
		}
		if err = VceWriteBFilters(buf, vce); err != nil {
			err = errors.Wrapf(err, "Failed to write B filters")
			return
		}
	}
	return
}

/*
func vceToString(vce VCE) (result string) {
	b, _ := json.MarshalIndent(vce, "", " ")
	result = string(b)

	return
}*/

func VceToJson(vce VCE) (result string) {
	b, _ := json.Marshal(vce)
	result = string(b)

	return
}
func vceHeadToJson(head VCEHead) (result string) {
	b, _ := json.Marshal(head)
	result = string(b)

	return
}

func CompactVCE(vce VCE) (compacted VCE, err error) {
	var writebuf = writerseeker.WriterSeeker{}

	if err = WriteVce(&writebuf, vce, VceName(vce.Head), false); err != nil {
		return
	}
	writeBytes, _ := ioutil.ReadAll(writebuf.Reader())

	var readbuf2 = bytes.NewReader(writeBytes)

	if compacted, err = ReadVce(readbuf2, false); err != nil {
		return
	}
	return
}

func CompactVceToJson(vce VCE) (result string) {
	// compact the vce before printing it:
	var err error
	var compacted VCE
	if compacted, err = CompactVCE(vce); err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}

	b, _ := json.MarshalIndent(compacted, "", "\t")
	result = string(b)
	return
}
