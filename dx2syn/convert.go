package main

import (
	"math"

	"github.com/chinenual/synergize/data"
	"github.com/pkg/errors"
)

func TranslateDx7ToVce(dx7Voice Dx7Voice) (vce data.VCE, err error) {
	if vce, err = helperBlankVce(); err != nil {
		return
	}

	for i := 0; i < 7; i++ {
		vce.Head.VNAME[i] = dx7Voice.VoiceName[i]
	}

	if dx7Voice.Algorithm == 4 {
		helperSetPatchType(&vce, 2)
	} else {
		err = errors.New("Limitation: currently only handle DX algorithm #4")
		return
	}
	// DX7 always uses 6 oscillators
	vce.Head.VOITAB = 5

	vce.Head.VIBRAT = dx7Voice.LfoSpeed
	vce.Head.VIBDEL = dx7Voice.LfoDelay
	vce.Head.VIBDEP = dx7Voice.LfoPitchModDepth

	// Transpose
	// May be modified if Coarse < 1
	vce.Head.VTRANS = int8(dx7Voice.Transpose - 24)

	for i, o := range dx7Voice.Osc {

		vce.Envelopes[i].FreqEnvelope.FDETUN = int8(o.OscDetune)

		vce.Envelopes[i].FreqEnvelope.OHARM = o.OscFreqCoarse
		// envelopes: DX amp envelopes always have 4 points
		vce.Envelopes[i].AmpEnvelope.NPOINTS = 4
		// Each Synergy oscillator is voice twice - for low and high key velocity response
		// Synergy envelopes are represented as quads of ValLow, ValHi, RateLow and RateHi
		// set both upper and lower envs the same
		// point1
		vce.Envelopes[i].AmpEnvelope.Table[0] = byte((math.Round(float64(o.EgLevel[0]) * 0.727)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[1] = byte((math.Round(float64(o.EgLevel[0]) * 0.727)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[2] = 99 - o.EgRate[0]
		vce.Envelopes[i].AmpEnvelope.Table[3] = 99 - o.EgRate[0]

		//point2
		vce.Envelopes[i].AmpEnvelope.Table[4] = byte((math.Round(float64(o.EgLevel[1]) * 0.727)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[5] = byte((math.Round(float64(o.EgLevel[1]) * 0.727)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[6] = (99 - o.EgRate[0]) +
			(99 - o.EgRate[1])
		vce.Envelopes[i].AmpEnvelope.Table[7] = (99 - o.EgRate[0]) +
			(99 - o.EgRate[1])

		//point3
		vce.Envelopes[i].AmpEnvelope.Table[8] = byte((math.Round(float64(o.EgLevel[2]) * 0.727)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[9] = byte((math.Round(float64(o.EgLevel[2]) * 0.727)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[10] = (99 - o.EgRate[0]) +
			(99 - o.EgRate[1]) + (99 - o.EgRate[2])
		vce.Envelopes[i].AmpEnvelope.Table[11] = (99 - o.EgRate[0]) +
			(99 - o.EgRate[1]) + (99 - o.EgRate[2])

		//point4
		vce.Envelopes[i].AmpEnvelope.Table[12] = byte((math.Round(float64(o.EgLevel[3]) * 0.727)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[13] = byte((math.Round(float64(o.EgLevel[3]) * 0.727)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[14] = (99 - o.EgRate[0]) +
			(99 - o.EgRate[1]) + (99 - o.EgRate[2]) + (99 - o.EgRate[3])
		vce.Envelopes[i].AmpEnvelope.Table[15] = (99 - o.EgRate[0]) +
			(99 - o.EgRate[1]) + (99 - o.EgRate[2]) + (99 - o.EgRate[3])

		// DX only has a single frequency envelope - replicate it on each Synergy osc:
		// NOTE the first point in the Synergy freq table is "special" - it stores a "freq.scale and wavetype" instead of rates
		// Like the amp table, the values are stored in quads, two values, two rates per point

		// envelopes: DX freq envelopes always have 4 points
		vce.Envelopes[i].FreqEnvelope.NPOINTS = 4

		// point1
		vce.Envelopes[i].FreqEnvelope.Table[0] = byte((math.Round(float64(dx7Voice.PitchEgLevel[0]) * .727)) + 55)
		vce.Envelopes[i].FreqEnvelope.Table[1] = byte((math.Round(float64(dx7Voice.PitchEgLevel[0]) * .727)) + 55)
		// special case for point1
		vce.Envelopes[i].FreqEnvelope.Table[2] = 0x80 // matches default from EDATA
		vce.Envelopes[i].FreqEnvelope.Table[3] = 0    // 0 == Sine, octave 0, freq int and amp int disabled

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
	}
	// ... everything else ...

	// if you need to abort, use:
	//	err = errors.New("an error message")
	//  return

	return

}

/*
int16 AmpTimeScale  := [100]int{0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7,
		8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23,
		24, 25, 26, 27, 28, 29, 30, 31,
		32, 33, 34, 35, 36, 37, 38, 39,
		40, 45, 51, 57, 64, 72, 81, 91,
		102, 115, 129, 145, 163, 183, 205, 230,
		258, 290, 326, 366, 411, 461, 517, 581,
		652, 732, 822, 922, 1035, 1162, 1304, 1464,
		1644, 1845, 2071, 2325, 2609, 2929, 3288, 3691,
		4143, 4650, 5219, 5859, 6576}

int16 FreqTimeScale  := [100]int{0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7,
		8, 9, 10, 11, 12, 13, 14, 15,
		25, 28, 32, 36, 40, 45, 51, 57,
		64, 72, 81, 91, 102, 115, 129, 145,
		163, 183, 205, 230, 258, 290, 326, 366,
		411, 461, 517, 581, 652, 732, 822, 922,
		1035, 1162, 1304, 1464, 1644, 1845, 2071, 2325,
		2609, 2929, 3288, 3691, 4143, 4650, 5219, 5859,
		6576, 7382, 8286, 9300, 10439, 11718, 13153, 14764,
		16572, 18600, 20078, 23436, 26306, 29528, 29529, 29530,
		29531, 29532, 29533, 29534, 29535}
*/
