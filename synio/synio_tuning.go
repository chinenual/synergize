package synio

import (
	"fmt"
	"io/ioutil"
	"os"

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
	ReferenceNote              int
	ReferenceFrequency         float64
}

var tuningParams = TuningParams{
	UseStandardTuning:          true,
	UseStandardKeyboardMapping: true,
	SCLPath:                    "",
	KBMPath:                    "",
	ReferenceNote:              69,
	ReferenceFrequency:         440.0,
}

func GetTuningParams() TuningParams {
	logger.Infof("GetTuningParams: %#v\n", tuningParams)
	return tuningParams
}

func GetTuningFrequencies(params TuningParams) (freqs []float64, err error) {
	var t scala.Tuning
	var s scala.Scale
	var k scala.KeyboardMapping
	logger.Infof("GetTuningFrequencies for %#v\n", params)
	if params.UseStandardTuning {
		if s, err = scala.ScaleEvenTemperment12NoteScale(); err != nil {
			logger.Errorf("ScaleEvenTemperment12NoteScale err: %v\n", err)
			return
		}
		if k, err = scala.KeyboardMappingTuneNoteTo(params.ReferenceNote, params.ReferenceFrequency); err != nil {
			logger.Errorf("KeyboardMappingTuneNoteTo err: %v\n", err)
			return
		}
	} else {
		if s, err = scala.ScaleFromSCLFile(params.SCLPath); err != nil {
			logger.Errorf("ScaleFromSCLFile err: %v\n", err)
			return
		}
		if params.UseStandardKeyboardMapping {
			if k, err = scala.KeyboardMappingTuneNoteTo(params.ReferenceNote, params.ReferenceFrequency); err != nil {
				logger.Errorf("KeyboardMappingTuneNoteTo err: %v\n", err)
				return
			}
		} else {
			if k, err = scala.KeyboardMappingFromKBMFile(params.KBMPath); err != nil {
				logger.Errorf("KeyboardMappingFromKBMFile err: %v\n", err)
				return
			}
		}
	}
	if t, err = scala.TuningFromSCLAndKBM(s, k); err != nil {
		logger.Errorf("TuningFromSCLAndKBM err: %v\n", err)
		return
	}

	for i := 0; i < 128; i++ {
		freqs = append(freqs, t.FrequencyForMidiNote(i))
	}
	tuningParams = params
	return
}

func SendTuningToSynergy(params TuningParams) (freqs []float64, err error) {
	c.Lock()
	defer c.Unlock()

	if freqs, err = GetTuningFrequencies(params); err != nil {
		return
	}
	if synioVerbose {
		logger.Infof("SYNIO: ** SendTuningToSynergy\n")
	}
	// Adjust the frequencies to Synergy format
	var scale = 1.175
	var intFreqs []uint16
	for _, f := range freqs {
		intFreqs = append(intFreqs, uint16(scale*f))
	}
	logger.Debug("Scala freq table: %v\n", freqs)
	logger.Debug("Synergy freq table: %v\n", intFreqs)
	var b []byte
	for _, f := range intFreqs {
		hob, lob := data.WordToBytes(f)
		b = append(b, lob, hob)
	}

	//dumpAddressSpace("DUMP.bin")
	if err = blockLoad(synAddrs.ROM_FREQTAB, b, "setFreqTable(ROM)"); err != nil {
		return
	}

	PrintFreqTable()
	return
}

func dumpAddressSpace(path string) {
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

	//dumpAddressSpace("DUMP.bin")
	if b, err = blockDump(synAddrs.ROM_FREQTAB, FREQTAB_LEN, "getFreqTable(ROM)"); err != nil {
		return
	}

	for i := 0; i < FREQTAB_LEN; i += 2 {
		w := data.BytesToWord(b[i+1], b[i])
		fmt.Printf("factory[%d]: %v\n", i/2, w)
	}
	if b, err = blockDump(synAddrs.WENDY_FREQTAB, FREQTAB_LEN, "getFreqTable(WENDY)"); err != nil {
		return
	}

	for i := 0; i < FREQTAB_LEN; i += 2 {
		w := data.BytesToWord(b[i+1], b[i])
		fmt.Printf("WENDY[%d]: %v\n", i/2, w)
	}
	return
}
