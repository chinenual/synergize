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

// Indexes are pretty confusing.
//   Synergy KEY value == MIDI Note - 28 (from MIDI.Z80, line 119)
//   FTAB index offset = 38 (from KYZ.Z80, line 59) (so SYNKEY+38 == FTAB index)
//   center frequency is computed from one-octave's worth of RFTAB value, starting at byte
//      index 96 (which is logic index 48)
// so... how do we convert from MIDI note to RFTAB index or vice versa?   We want to make sure
// we update the octave that the synergy will read for scale definition (at index 48) - and that needs to match
// in the end, I found the offsets empirically by a bit of brute force.

// from COMMON.Z80 "FTAB":
// NOTE: index is NOT by MIDI note: FTAB index 68 == Middle C (which would be MIDI 60).
// Index factoryFrequencyValues by MIDI note
// Index factoryRftabValues by [MIDInote + midiToRftabOffset]
const midiToRftabOffset = 8

var factoryRftabValues = []uint16{
	/*   0 */ 0, 2, 4, 6, 8, 10, 12, 14,
	/*   8 */ 15, 16, 17, 18, 19, 20, 21, 22,
	/*  16 */ 24, 25, 27, 28, 30, 32, 34, 36,
	/*  24 */ 38, 40, 43, 45, 48, 51, 54, 57,
	/*  32 */ 61, 64, 68, 72, 76, 81, 86, 91,
	/*  40 */ 96, 102, 108, 115, 122, 129, 137, 145,
	/*  48 */ 153, 163, 172, 183, 193, 205, 217, 230,
	/*  56 */ 244, 258, 274, 290, 307, 326, 345, 366,
	/*  64 */ 387, 411, 435, 461, 488, 517, 548, 581,
	/*  72 */ 615, 652, 691, 732, 775, 822, 870, 922,
	/*  80 */ 977, 1035, 1097, 1162, 1231, 1304, 1382, 1464,
	/*  88 */ 1551, 1644, 1741, 1845, 1955, 2071, 2194, 2325,
	/*  96 */ 2463, 2609, 2765, 2929, 3103, 3288, 3483, 3691,
	/* 104 */ 3910, 4143, 4389, 4650, 4926, 5219, 5530, 5859,
	/* 112 */ 6207, 6576, 6967, 7382, 7820, 8286, 8778, 9300,
	/* 120 */ 9853, 10439, 11060, 11718, 12414, 13153, 13935, 14764,
}

// standard 12-tone equal temperament frequencies:
// NOTE: 440.0 is at index 69 - so index is by MIDI note
var factoryFrequencyValues = []float64{
	/*   0 */ 8.17579891564371, 8.661957218027254, 9.177023997418992, 9.722718241315032,
	/*   4 */ 10.300861153527189, 10.913382232281378, 11.56232570973858, 12.249857374429675,
	/*   8 */ 12.978271799373285, 13.750000000000005, 14.567617547440317, 15.433853164253883,
	/*  12 */ 16.35159783128742, 17.32391443605451, 18.354047994837984, 19.445436482630065,
	/*  16 */ 20.601722307054377, 21.826764464562757, 23.124651419477157, 24.499714748859343,
	/*  20 */ 25.95654359874658, 27.50000000000001, 29.135235094880635, 30.867706328507765,
	/*  24 */ 32.70319566257484, 34.64782887210902, 36.70809598967597, 38.89087296526013,
	/*  28 */ 41.203444614108754, 43.653528929125514, 46.249302838954314, 48.99942949771869,
	/*  32 */ 51.91308719749316, 55.00000000000002, 58.27047018976127, 61.73541265701553,
	/*  36 */ 65.40639132514968, 69.29565774421803, 73.41619197935194, 77.78174593052026,
	/*  40 */ 82.40688922821751, 87.30705785825103, 92.49860567790863, 97.99885899543737,
	/*  44 */ 103.82617439498632, 110.00000000000004, 116.54094037952254, 123.47082531403106,
	/*  48 */ 130.81278265029937, 138.59131548843607, 146.83238395870387, 155.56349186104052,
	/*  52 */ 164.81377845643502, 174.61411571650206, 184.99721135581726, 195.9977179908748,
	/*  56 */ 207.65234878997256, 220.00000000000009, 233.08188075904516, 246.94165062806204,
	/*  60 */ 261.62556530059874, 277.18263097687213, 293.66476791740774, 311.12698372208104,
	/*  64 */ 329.62755691287003, 349.2282314330041, 369.9944227116345, 391.9954359817496,
	/*  68 */ 415.3046975799451, 440.00000000000017, 466.1637615180903, 493.88330125612407,
	/*  72 */ 523.2511306011975, 554.3652619537443, 587.3295358348155, 622.2539674441621,
	/*  76 */ 659.2551138257401, 698.4564628660082, 739.988845423269, 783.9908719634992,
	/*  80 */ 830.6093951598903, 880.0000000000003, 932.3275230361807, 987.7666025122481,
	/*  84 */ 1046.502261202395, 1108.7305239074885, 1174.659071669631, 1244.5079348883241,
	/*  88 */ 1318.5102276514801, 1396.9129257320164, 1479.977690846538, 1567.9817439269984,
	/*  92 */ 1661.2187903197805, 1760.0000000000007, 1864.6550460723613, 1975.5332050244963,
	/*  96 */ 2093.00452240479, 2217.4610478149757, 2349.3181433392633, 2489.0158697766483,
	/* 100 */ 2637.020455302958, 2793.8258514640347, 2959.955381693076, 3135.963487853997,
	/* 104 */ 3322.437580639561, 3520.0000000000014, 3729.3100921447226, 3951.0664100489926,
	/* 108 */ 4186.00904480958, 4434.922095629951, 4698.636286678527, 4978.031739553297,
	/* 112 */ 5274.040910605916, 5587.651702928069, 5919.910763386152, 6271.926975707994,
	/* 116 */ 6644.875161279122, 7040.000000000003, 7458.620184289445, 7902.132820097985,
	/* 120 */ 8372.01808961916, 8869.844191259903, 9397.272573357053, 9956.063479106593,
	/* 124 */ 10548.081821211832, 11175.303405856139, 11839.821526772304, 12543.853951415987}

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
		if k, err = scala.KeyboardMappingStartScaleOnAndTuneNoteTo(params.MiddleNote, params.ReferenceNote, params.ReferenceFrequency); err != nil {
			logger.Errorf("KeyboardMappingStartScaleOnAndTuneNoteTo err: %v\n", err)
			return
		}
	} else {
		if s, err = scala.ScaleFromSCLFile(params.SCLPath); err != nil {
			logger.Errorf("ScaleFromSCLFile err: %v\n", err)
			return
		}
		if params.UseStandardKeyboardMapping {
			if k, err = scala.KeyboardMappingStartScaleOnAndTuneNoteTo(params.MiddleNote, params.ReferenceNote, params.ReferenceFrequency); err != nil {
				logger.Errorf("KeyboardMappingStartScaleOnAndTuneNoteTo err: %v\n", err)
				return
			}
		} else {
			if k, err = scala.KeyboardMappingFromKBMFile(params.KBMPath); err != nil {
				logger.Errorf("KeyboardMappingFromKBMFile err: %v\n", err)
				return
			}
		}
	}
	if MIDI_SYN_OFFSET != 0 {
		k.MiddleNote -= MIDI_SYN_OFFSET
		k.Name += " (with MiddleNote override for Synergy)"
	}
	logger.Infof("Selected SCL %#v\n", s)
	logger.Infof("Selected KBM %#v\n", k)
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
	for i := 0; i < 128+MIDI_SYN_OFFSET; i++ {
		freqs = append(freqs, t.FrequencyForMidiNote(i))
	}
	tuningParams = params
	return
}

// Scala note 60 == middle C - that needs to map tp Synergy note 58, so offset by 2
//const MIDI_SYN_OFFSET = 2
const MIDI_SYN_OFFSET = 0

func scaleFrequencies(params TuningParams, freqs []float64) (rftabValues []uint16) {
	rftabValues = make([]uint16, 128)
	if params.UseStandardTuning {
		// special case "standard tuning" to mean "factory settings" - which are not identical
		// to what we compute via scaling due to different rounding errors - so just reset the
		// factory values
		_ = copy(rftabValues, factoryRftabValues)
	} else {
		// most of the table is identical to the factory settings we only override the 12 notes that
		// the synergy uses to compute the key center frequencies (reading the Z80, I would have thought 48..59,
		// but Hal remembers "middle C" and in fact, middle C (60..71) is what works
		const tuneStartMidi = 60
		_ = copy(rftabValues, factoryRftabValues)
		for i := tuneStartMidi; i < tuneStartMidi+12; i++ {
			rftab_i := i + midiToRftabOffset

			// determine the delta from the factory value as a ratio, then apply that to the ROM table value
			tgtFreq := freqs[i]
			romFreq := factoryFrequencyValues[i]
			ratio := tgtFreq / romFreq

			rftabValues[rftab_i] = uint16(math.Round(ratio * float64(factoryRftabValues[rftab_i])))
			logger.Infof("ROM delta [%v] midi:%v %v %v (%v)\n", rftab_i, i, rftabValues[rftab_i], factoryRftabValues[rftab_i], ratio)
		}
		/**
		//var scale = 2.097152 //Synergia
		var scale = 1.17671
		for i := 0; i < 128; i++ {
			var v uint16
			if i < 24 {
				// retain the factory linear scale at bottom range
				v = factoryRftabValues[i]
			} else {
				v = uint16(math.Round(scale * freqs[i]))
			}
			rftabValues = append(rftabValues, v)
		}
		*/
	}
	logger.Infof("Scala freq table: %v\n", freqs)
	logger.Infof("Synergy freq table: %v\n", rftabValues)
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
	rftabValues := scaleFrequencies(params, freqs)
	var b []byte
	for _, f := range rftabValues {
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
