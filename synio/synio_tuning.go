package synio

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"

	"github.com/chinenual/synergize/io"

	"github.com/chinenual/go-scala"
	"github.com/chinenual/synergize/data"
	"github.com/chinenual/synergize/logger"
)

const FREQTAB_LEN = 127 * 2

type TuningParams struct {
	UseStandardTuning          bool
	UseStandardKeyboardMapping bool
	SCLPath                    string
	KBMPath                    string
	MiddleNote                 int
	ReferenceNote              int
	ReferenceFrequency         float64
}

var tuningParams = TuningParams{
	UseStandardTuning:          true,
	UseStandardKeyboardMapping: true,
	SCLPath:                    "",
	KBMPath:                    "",
	MiddleNote:                 60,
	ReferenceNote:              69,
	ReferenceFrequency:         440.0,
}

const fixedReferenceNote = 69
const fixedReferenceFrequency = 440.0

//const midi0Freq = 8.17579891564371 // or 440.0 * pow( 2.0, - (69.0/12.0 ) )

// from COMMON.Z80 "FTAB":
var factoryROMTableValues = []uint16{0, 2, 4, 6, 8, 10, 12, 14,
	15, 16, 17, 18, 19, 20, 21, 22,
	24, 25, 27, 28, 30, 32, 34, 36,
	38, 40, 43, 45, 48, 51, 54, 57,
	61, 64, 68, 72, 76, 81, 86, 91,
	96, 102, 108, 115, 122, 129, 137, 145,
	153, 163, 172, 183, 193, 205, 217, 230,
	244, 258, 274, 290, 307, 326, 345, 366,
	387, 411, 435, 461, 488, 517, 548, 581,
	615, 652, 691, 732, 775, 822, 870, 922,
	977, 1035, 1097, 1162, 1231, 1304, 1382, 1464,
	1551, 1644, 1741, 1845, 1955, 2071, 2194, 2325,
	2463, 2609, 2765, 2929, 3103, 3288, 3483, 3691,
	3910, 4143, 4389, 4650, 4926, 5219, 5530, 5859,
	6207, 6576, 6967, 7382, 7820, 8286, 8778, 9300,
	9853, 10439, 11060, 11718, 12414, 13153, 13935, 14764,
}

func GetTuningParams() TuningParams {
	logger.Infof("GetTuningParams: %#v\n", tuningParams)
	return tuningParams
}

func GetTuningFrequencies(params TuningParams) (freqs []float64, tones []scala.Tone, scalePos []int, err error) {
	var t scala.Tuning
	var s scala.Scale
	var k scala.KeyboardMapping
	logger.Infof("GetTuningFrequencies for %#v\n", params)
	if params.UseStandardTuning {
		if s, err = scala.ScaleEvenTemperment12NoteScale(); err != nil {
			logger.Errorf("ScaleEvenTemperment12NoteScale err: %v\n", err)
			return
		}
		if k, err = scala.KeyboardMappingStartScaleOnAndTuneNoteTo(params.MiddleNote, fixedReferenceNote, fixedReferenceFrequency); err != nil {
			logger.Errorf("KeyboardMappingStartScaleOnAndTuneNoteTo err: %v\n", err)
			return
		}
	} else {
		if s, err = scala.ScaleFromSCLFile(params.SCLPath); err != nil {
			logger.Errorf("ScaleFromSCLFile err: %v\n", err)
			return
		}
		if params.UseStandardKeyboardMapping {
			if k, err = scala.KeyboardMappingStartScaleOnAndTuneNoteTo(params.MiddleNote, fixedReferenceNote, fixedReferenceFrequency); err != nil {
				logger.Errorf("KeyboardMappingStartScaleOnAndTuneNoteTo err: %v\n", err)
				return
			}
		} else {
			if k, err = scala.KeyboardMappingFromKBMFile(params.KBMPath); err != nil {
				logger.Errorf("KeyboardMappingFromKBMFile err: %v\n", err)
				return
			}
			/*
				// override the reference values since Synergy tables require A440 reference
				k.TuningFrequency = fixedReferenceFrequency
				k.TuningConstantNote = fixedReferenceNote
				k.TuningPitch = fixedReferenceFrequency / midi0Freq
			*/
		}
	}
	if t, err = scala.TuningFromSCLAndKBM(s, k); err != nil {
		logger.Errorf("TuningFromSCLAndKBM err: %v\n", err)
		return
	}
	tones = s.Tones
	for i := 0; i < 128; i++ {
		scalePos = append(scalePos, t.ScalePositionForMidiNote(i))
	}
	logger.Infof("Scale tones: %v\n", tones)
	logger.Infof("Scale scalePos: %v\n", scalePos)
	for i := 0; i < 128; i++ {
		freqs = append(freqs, t.FrequencyForMidiNote(i))
	}
	tuningParams = params
	return
}

func scaleFrequencies(params TuningParams, freqs []float64) (intFreqs []uint16) {
	if params.UseStandardTuning {
		// special case "standard tuning" to mean "factory settings" - which are not identical
		// to what we compute via scaling.
		intFreqs = factoryROMTableValues
	} else {
		var scale = 1.17671
		for i, f := range freqs {
			var v uint16
			if i < 24 {
				v = factoryROMTableValues[i]
			} else {
				v = uint16(math.Round(scale * f))
			}
			intFreqs = append(intFreqs, v)
		}
	}
	logger.Infof("Scala freq table: %v\n", freqs)
	logger.Infof("Synergy freq table: %v\n", intFreqs)
	return
}

func SendTuningToSynergy(params TuningParams) (freqs []float64, tones []scala.Tone, scalePos []int, err error) {
	c.Lock()
	defer c.Unlock()

	if err = getSynergyAddrs(); err != nil {
		return
	}

	if freqs, tones, scalePos, err = GetTuningFrequencies(params); err != nil {
		return
	}
	if synioVerbose {
		logger.Infof("SYNIO: ** SendTuningToSynergy\n")
	}
	// Adjust the frequencies to Synergy format
	intFreqs := scaleFrequencies(params, freqs)
	var b []byte
	for _, f := range intFreqs {
		hob, lob := data.WordToBytes(f)
		b = append(b, lob, hob)
	}

	addr := synAddrs.WENDY_FREQTAB
	if io.SynergyConnectionType() == "vst" {
		addr = synAddrs.ROM_FREQTAB
	}
	//dumpAddressSpace("DUMP.bin")
	if err = blockLoad(addr, b, "setFreqTable(ROM)"); err != nil {
		return
	}
	//if err = reloadNoteGenerators(); err != nil {
	//	return
	//}

	//PrintFreqTable()
	return
}

func dumpSynergyAddressSpace(path string) {
	var b []byte
	var err error

	total := uint16(65323)
	chunk := total / 4
	for i := 0; i < int(total); i += int(chunk) {
		var bchunk []byte
		n := chunk
		if i+int(n) > int(total) {
			n = total - uint16(i)
		}

		//	b, err = blockDump(uint16(0), uint16(65323), "dump addr space")
		bchunk, err = blockDump(uint16(i), uint16(n), "dump addr space")
		if err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(1)
		}
		b = append(b, bchunk...)
	}
	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func PrintFreqTable() (err error) {
	var b []byte
	if err = getSynergyAddrs(); err != nil {
		return
	}

	//dumpSynergyAddressSpace("DUMP.bin")
	if b, err = blockDump(synAddrs.ROM_FREQTAB, FREQTAB_LEN, "getFreqTable(ROM)"); err != nil {
		return
	}

	for i := 0; i < FREQTAB_LEN; i += 2 {
		w := data.BytesToWord(b[i+1], b[i])
		fmt.Printf("ROM FTAB[%d]: %v\n", i/2, w)
	}
	/*
		if b, err = blockDump(synAddrs.WENDY_FREQTAB, FREQTAB_LEN, "getFreqTable(WENDY)"); err != nil {
			return
		}

		for i := 0; i < FREQTAB_LEN; i += 2 {
			w := data.BytesToWord(b[i+1], b[i])
			fmt.Printf("WENDY[%d]: %v\n", i/2, w)
		}
	*/
	return
}
