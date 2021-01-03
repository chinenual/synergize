package main

import (
	"fmt"
	"math"

	"github.com/chinenual/synergize/data"
)

func TranslateDx7ToVce(dx7Voice Dx7Voice) (vce data.VCE, err error) {
	if vce, err = helperBlankVce(); err != nil {
		return
	}

	for i := 0; i < 7; i++ {
		vce.Head.VNAME[i] = dx7Voice.VoiceName[i]
	}

	if err = helperSetAlgorithmPatchType(&vce, dx7Voice.Algorithm, dx7Voice.Feedback); err != nil {
		return
	}
	// DX7 always uses 6 oscillators
	vce.Head.VOITAB = 5

	vce.Head.VIBRAT = dx7Voice.LfoSpeed
	vce.Head.VIBDEL = dx7Voice.LfoDelay
	vce.Head.VIBDEP = dx7Voice.LfoPitchModDepth

	// Transpose
	// Set for OHARM +1 of any of the OPs OscFreqCoarse are 0 (1/2 octave)

	vce.Head.VTRANS = int8(dx7Voice.Transpose - 24)

	transposedDown := false
	attkR := 0
	decyR := 0
	sustR := 0
	relsR := 0
	var OSClevelPercent float64

	var ms [4]int

	for _, o := range dx7Voice.Osc {
		if o.OscFreqCoarse == 0 {
			transposedDown = true
			vce.Head.VTRANS = -12
			break
		}
	}

	for i, o := range dx7Voice.Osc {
		vce.Envelopes[i].FreqEnvelope.OHARM = o.OscFreqCoarse
		//At least one OSC Marh = 0
		if transposedDown {
			vce.Envelopes[i].FreqEnvelope.OHARM++
		}
		// Set OSC detune

		//vce.Envelopes[i].FreqEnvelope.FDETUN = int8(o.OscDetune)

		// type = 1  : no loop (and LOOPPT and SUSTAINPT are accelleration rates not point positions)
		// type = 2  : S only
		// type = 3  : L and S - L must be before S
		// type = 4  : R and S - R must be before S
		// WARNING: when type1, the LOOPPT and SUSTAINPT values are _acceleration_ rates, not point positions. What a pain.
		// Set for Sustain poiint only.
		vce.Envelopes[i].AmpEnvelope.ENVTYPE = 2
		//  Always DX7 Sustain Point
		vce.Envelopes[i].AmpEnvelope.SUSTAINPT = 2
		// envelopes: DX amp envelopes always have 4 points
		vce.Envelopes[i].AmpEnvelope.NPOINTS = 4

		OSClevelPercent = float64((o.OperatorOutputLevel) / 99.00)
		for k := 0; k < 4; k++ {
			o.EgLevel[k] = byte(float64(o.EgLevel[k]) * OSClevelPercent)
		}

		// Each Synergy oscillator is voice twice - for low and high key velocity response
		// Synergy envelopes are represented as quads of ValLow, ValHi, RateLow and RateHi
		// set both upper and lower envs the same

		ms = computeDurationsMs(o.EgLevel, o.EgRate)

		attkR = ms[0]
		decyR = ms[1]
		sustR = ms[2]
		relsR = ms[3]

		fmt.Printf(" %s %f \n", " OSC % = ", OSClevelPercent)
		fmt.Printf(" %s %d \n", " attkR = ", attkR)
		fmt.Printf(" %s %d \n", " decyR = ", decyR)
		fmt.Printf(" %s %d \n", " sustR = ", sustR)
		fmt.Printf(" %s %d \n \n", " relsR = ", relsR)

		// point1
		vce.Envelopes[i].AmpEnvelope.Table[0] = byte((math.Round(float64(o.EgLevel[0]) * 0.727)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[1] = byte((math.Round(float64(o.EgLevel[0]) * 0.727)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[2] = byte(helperNearestAmpTimeIndex(attkR))
		vce.Envelopes[i].AmpEnvelope.Table[3] = byte(helperNearestAmpTimeIndex(attkR))

		//point2
		vce.Envelopes[i].AmpEnvelope.Table[4] = byte((math.Round(float64(o.EgLevel[1]) * 0.727)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[5] = byte((math.Round(float64(o.EgLevel[1]) * 0.727)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[6] = byte(helperNearestAmpTimeIndex(attkR + decyR))
		vce.Envelopes[i].AmpEnvelope.Table[7] = byte(helperNearestAmpTimeIndex(attkR + decyR))

		//point3
		vce.Envelopes[i].AmpEnvelope.Table[8] = byte((math.Round(float64(o.EgLevel[2]) * 0.727)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[9] = byte((math.Round(float64(o.EgLevel[2]) * 0.727)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[10] = byte(helperNearestAmpTimeIndex(attkR + decyR + sustR))
		vce.Envelopes[i].AmpEnvelope.Table[11] = byte(helperNearestAmpTimeIndex(attkR + decyR + sustR))

		//point4
		vce.Envelopes[i].AmpEnvelope.Table[12] = byte((math.Round(float64(o.EgLevel[3]) * 0.727)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[13] = byte((math.Round(float64(o.EgLevel[3]) * 0.727)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[14] = byte(helperNearestAmpTimeIndex(attkR + decyR + sustR + relsR))
		vce.Envelopes[i].AmpEnvelope.Table[15] = byte(helperNearestAmpTimeIndex(attkR + decyR + sustR + relsR))

		// DX only has a single frequency envelope - replicate it on each Synergy osc:
		// NOTE the first point in the Synergy freq table is "special" - it stores a "freq.scale and wavetype" instead of rates
		// Like the amp table, the values are stored in quads, two values, two rates per point
		/* ****************************************
		// envelopes: DX freq envelopes always have 4 points
		vce.Envelopes[i].FreqEnvelope.NPOINTS = 4

		// point1
		vce.Envelopes[i].FreqEnvelope.Table[0] = byte((math.Round(float64(dx7Voice.PitchEgLevel[0]) * .727)) + 55)
		vce.Envelopes[i].FreqEnvelope.Table[1] = byte((math.Round(float64(dx7Voice.PitchEgLevel[0]) * .727)) + 55)
		// special case for point1
		vce.Envelopes[i].FreqEnvelope.Table[2] = 0x80 // matches default from EDATA
		vce.Envelopes[i].FreqEnvelope.Table[3] = 0 // 0 == Sine, octave 0, freq int and amp int disabled

		// point2
		vce.Envelopes[i].FreqEnvelope.Table[4] = byte((math.Round(float64(dx7Voice.PitchEgLevel[1]) * .727)) + 55)
		vce.Envelopes[i].FreqEnvelope.Table[5] = byte((math.Round(float64(dx7Voice.PitchEgLevel[1]) * .727)) + 55)
		vce.Envelopes[i].FreqEnvelope.Table[6] = (99 - dx7Voice.PitchEgRate[0]) +
			(99 - dx7Voice.PitchEgRate[1])
		vce.Envelopes[i].FreqEnvelope.Table[7] = (99 - dx7Voice.PitchEgRate[0]) +
			(99 - dx7Voice.PitchEgRate[1])

		// point3
		vce.Envelopes[i].FreqEnvelope.Table[8] = byte((math.Round(float64(dx7Voice.PitchEgLevel[2]) * .727)) + 55)
		vce.Envelopes[i].FreqEnvelope.Table[9] = byte((math.Round(float64(dx7Voice.PitchEgLevel[2]) * .727)) + 55)
		vce.Envelopes[i].FreqEnvelope.Table[10] = (99 - dx7Voice.PitchEgRate[0]) +
			(99 - dx7Voice.PitchEgRate[1]) + (99 - dx7Voice.PitchEgRate[2])
		vce.Envelopes[i].FreqEnvelope.Table[11] = (99 - dx7Voice.PitchEgRate[0]) +
			(99 - dx7Voice.PitchEgRate[1]) + (99 - dx7Voice.PitchEgRate[2])

		// point4
		vce.Envelopes[i].FreqEnvelope.Table[12] = byte((math.Round(float64(dx7Voice.PitchEgLevel[3]) * .727)) + 55)
		vce.Envelopes[i].FreqEnvelope.Table[13] = byte((math.Round(float64(dx7Voice.PitchEgLevel[3]) * .727)) + 55)
		vce.Envelopes[i].FreqEnvelope.Table[14] = (99 - dx7Voice.PitchEgRate[0]) +
			(99 - dx7Voice.PitchEgRate[1]) + (99 - dx7Voice.PitchEgRate[2]) + (99 - dx7Voice.PitchEgRate[3])
		vce.Envelopes[i].FreqEnvelope.Table[15] = (99 - dx7Voice.PitchEgRate[0]) +
			(99 - dx7Voice.PitchEgRate[1]) + (99 - dx7Voice.PitchEgRate[2]) + (99 - dx7Voice.PitchEgRate[3])
		*/
	}
	// ... everything else ...

	// if you need to abort, use:
	//	err = errors.New("an error message")
	//  return

	return
}

var DXRisetoSYN = [100]byte{53, 50, 47, 44, 41, 38, 36, 34, 31, 29, 27, 26, 25, 24, 23, 22, 22,
	21, 21, 20, 20, 19, 19, 18, 18, 18, 17, 17, 17, 17, 17, 17, 16, 16, 16, 16, 16, 16, 16, 16,
	16, 16, 16, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

var DXDecaytoSYN = [100]byte{73, 72, 71, 70, 68, 67, 67, 66, 65, 64, 63, 62, 61, 60, 59, 58, 57, 57,
	56, 54, 50, 46, 42, 39, 35, 33, 31, 30, 28, 26, 25, 24, 23, 23, 22, 22, 21, 21, 21, 21, 21,
	20, 20, 20, 20, 19, 19, 18, 18, 17, 17, 17, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16,
	16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0}

var EGrateDecay = [100]int{318, 284, 249, 215, 181, 167, 155, 141, 128, 115, 105, 94, 84, 73, 63,
	58, 54, 50, 44, 4, 35, 32, 28, 24, 20, 18, 16, 15, 13, 11, 10, 9, 8, 7, 7, 7, 7, 67, 6, 6, 6,
	6, 6, 5, 4, 4, 4, 3, 3, 3, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

var EGrateRise = [100]int{38, 35, 32, 29, 26, 23, 21, 18, 16, 14, 12, 11, 10, 9, 8,
	8, 7, 6, 6, 5, 5, 4, 4, 4, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
