package main

import (
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

	// transposedDown := false
	attkR := 0
	decyR := 0
	sustR := 0
	relsR := 0
	var OSClevelPercent float64
	var VeloctiyPercent float64

	var ms [4]int
	/*
		for _, o := range dx7Voice.Osc {
			if o.OscFreqCoarse == 0 {
				transposedDown = true
				//vce.Head.VTRANS = -12
				break
			}
		}
	*/
	for i, o := range dx7Voice.Osc {
		/*
			// Set OSC mode
			if o.OscMode == true {
				vce.Envelopes[i].FreqEnvelope.OHARM = 1
			}
			else
			{
				vce.Envelopes[i].FreqEnvelope.OHARM = 0
			}
		*/
		// ******************************************************************************************
		// *************** put key scaling in filter B  filter B[0 - 31] for each OSC  **************
		// ******************************************************************************************

		//Activate FILTER B above per voice above (in Header)

		vce.Head.FILTER[i] = int8(i + 1) //set filter B on for voice, b-filters are indicated by the 1-based osc index

		// set "0" freq to match Synergy freq.

		// Assumes no A-filter - so B filter for osc 1 (index 0) is always stored at 0:
		vce.Filters[i][(BreakPoint[o.KeyLevelScalingBreakPoint])] = 0 //KEY to FREQ Array is BreakPoint[] (below)

		// Scale from DX7 0 to 99 to Syn -64 to 63    //using DX 50 = 0

		lMax := float64(o.KeyLevelScalingLeftDepth) * 0.63 //Trusting the DX7 is in Db also)
		rMax := float64(o.KeyLevelScalingRightDepth) * 0.63

		//  set Key Scaling curve below and above break point

		// for linear, we compute via linear function y = slope*x + b
		// b is the y value at "0" where "0" is the breakpoint, -- where y is by definition 0. So b is always 0
		//
		// for exponential, we base the curve on the array from the Dexed soft synth:
		// const uint8_t exp_scale_data[] = {
		//    0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 14, 16, 19, 23, 27, 33, 39, 47, 56, 66,
		//    80, 94, 110, 126, 142, 158, 174, 190, 206, 222, 238, 250
		//};
		// this can be modeled with the following equation:
		//   y = pow(2.0, x / 3.5) * scale
		//   where
		//    x = abs(offset from the breakpoint)
		//    scale = lMax / 256.0  (or rMax)
		const expBase = 2.0
		const expDivisor = 3.5
		const expScale = 256.0

		switch o.KeyLevelScalingLeftCurve { //0=-LIN, -EXP, +EXP, +LIN
		case 0:
			//-linear from -lMax to 0
			slope := lMax / float64(BreakPoint[o.KeyLevelScalingBreakPoint]-0)
			for k := byte(0); k < BreakPoint[o.KeyLevelScalingBreakPoint]; k++ {
				vce.Filters[i][k] = int8(math.Round(-lMax + slope*float64(k)))
			}
		case 1:
			//-EXP from -lMax to 0
			for k := byte(0); k < BreakPoint[o.KeyLevelScalingBreakPoint]; k++ {
				x := float64(BreakPoint[o.KeyLevelScalingBreakPoint] - k)
				vce.Filters[i][k] = int8(math.Pow(expBase, x/expDivisor) / expScale * -lMax)
			}
		case 2:
			//EXP from lMax to 0
			for k := byte(0); k < BreakPoint[o.KeyLevelScalingBreakPoint]; k++ {
				x := float64(BreakPoint[o.KeyLevelScalingBreakPoint] - k)
				vce.Filters[i][k] = int8(math.Pow(expBase, x/expDivisor) / expScale * lMax)
			}
		case 3:
			//linear from lMax to 0
			slope := -lMax / float64(BreakPoint[o.KeyLevelScalingBreakPoint]-0)
			for k := byte(0); k < BreakPoint[o.KeyLevelScalingBreakPoint]; k++ {
				vce.Filters[i][k] = int8(math.Round(lMax + slope*float64(k)))
			}
		}

		switch o.KeyLevelScalingRightCurve { //0=-LIN, -EXP, +EXP, +LIN
		case 0:
			// -Linear from 0 to -rMax
			slope := -rMax / float64(32-BreakPoint[o.KeyLevelScalingBreakPoint])
			for k := BreakPoint[o.KeyLevelScalingBreakPoint] + 1; k < 32; k++ {
				vce.Filters[i][k] = int8(math.Round(slope * float64(k-BreakPoint[o.KeyLevelScalingBreakPoint])))
			}
		case 1:
			// -EXP from 0 to -rMax
			for k := BreakPoint[o.KeyLevelScalingBreakPoint] + 1; k < 32; k++ {
				x := float64(k - BreakPoint[o.KeyLevelScalingBreakPoint])
				vce.Filters[i][k] = int8(math.Pow(expBase, x/expDivisor) / expScale * -rMax)
			}
		case 2:
			// EXP from 0 to rMax
			for k := BreakPoint[o.KeyLevelScalingBreakPoint] + 1; k < 32; k++ {
				x := float64(k - BreakPoint[o.KeyLevelScalingBreakPoint])
				vce.Filters[i][k] = int8(math.Pow(expBase, x/expDivisor) / expScale * rMax)
			}
		case 3:
			// Linear from 0 to rMax
			slope := rMax / float64(32-BreakPoint[o.KeyLevelScalingBreakPoint])
			for k := BreakPoint[o.KeyLevelScalingBreakPoint] + 1; k < 32; k++ {
				vce.Filters[i][k] = int8(math.Round(slope * float64(k-BreakPoint[o.KeyLevelScalingBreakPoint])))
			}
		}

		// *****************************************************************

		vce.Envelopes[i].FreqEnvelope.OHARM = o.OscFreqCoarse
		//At least one OSC Oharm = 0
		if o.OscFreqCoarse == 0 {
			vce.Envelopes[i].FreqEnvelope.OHARM = 1
		}
		// Set OSC detune

		vce.Envelopes[i].FreqEnvelope.FDETUN = int8(o.OscDetune - 7)

		// type = 1  : no loop (and LOOPPT and SUSTAINPT are accelleration rates not point positions)
		// type = 2  : S only
		// type = 3  : L and S - L must be before S
		// type = 4  : R and S - R must be before S
		// WARNING: when type1, the LOOPPT and SUSTAINPT values are _acceleration_ rates, not point positions. What a pain.
		// Set for Sustain poiint only.
		vce.Envelopes[i].AmpEnvelope.ENVTYPE = 2
		//  Always DX7 Sustain Point
		vce.Envelopes[i].AmpEnvelope.SUSTAINPT = 3
		// envelopes: DX amp envelopes always have 4 points
		vce.Envelopes[i].AmpEnvelope.NPOINTS = 4

		// set lower Env levels for velocity sensitivity
		switch o.KeyVelocitySensitivity {
		case 0:
			{
				VeloctiyPercent = 1.0
			}
		case 1:
			{
				VeloctiyPercent = .90
			}
		case 2:
			{
				VeloctiyPercent = .80
			}
		case 3:
			{
				VeloctiyPercent = .70
			}
		case 4:
			{
				VeloctiyPercent = .60
			}
		case 5:
			{
				VeloctiyPercent = .50
			}
		case 6:
			{
				VeloctiyPercent = .40
			}
		case 7:
			{
				VeloctiyPercent = .30
			}
		}

		OSClevelPercent = float64(float64(o.OperatorOutputLevel) / 99.00)
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
		/*
			fmt.Printf(" %s %f \n", " OSC % = ", OSClevelPercent)
			fmt.Printf(" %s %d \n", " attkR = ", attkR)
			fmt.Printf(" %s %d \n", " decyR = ", decyR)
			fmt.Printf(" %s %d \n", " sustR = ", sustR)
			fmt.Printf(" %s %d \n \n", " relsR = ", relsR)
		*/

		// point1
		vce.Envelopes[i].AmpEnvelope.Table[0] = byte((math.Round(float64(o.EgLevel[0]) * 0.727 * VeloctiyPercent)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[1] = byte((math.Round(float64(o.EgLevel[0]) * 0.727)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[2] = byte(helperNearestAmpTimeIndex(attkR))
		vce.Envelopes[i].AmpEnvelope.Table[3] = byte(helperNearestAmpTimeIndex(attkR))

		//point2
		vce.Envelopes[i].AmpEnvelope.Table[4] = byte((math.Round(float64(o.EgLevel[1]) * 0.727 * VeloctiyPercent)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[5] = byte((math.Round(float64(o.EgLevel[1]) * 0.727)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[6] = byte(helperNearestAmpTimeIndex(decyR))
		vce.Envelopes[i].AmpEnvelope.Table[7] = byte(helperNearestAmpTimeIndex(decyR))

		//point3
		vce.Envelopes[i].AmpEnvelope.Table[8] = byte((math.Round(float64(o.EgLevel[2]) * 0.727 * VeloctiyPercent)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[9] = byte((math.Round(float64(o.EgLevel[2]) * 0.727)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[10] = byte(helperNearestAmpTimeIndex(sustR))
		vce.Envelopes[i].AmpEnvelope.Table[11] = byte(helperNearestAmpTimeIndex(sustR))

		//point4
		vce.Envelopes[i].AmpEnvelope.Table[12] = byte((math.Round(float64(o.EgLevel[3]) * 0.727 * VeloctiyPercent)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[13] = byte((math.Round(float64(o.EgLevel[3]) * 0.727)) + 55)
		vce.Envelopes[i].AmpEnvelope.Table[14] = byte(helperNearestAmpTimeIndex(relsR))
		vce.Envelopes[i].AmpEnvelope.Table[15] = byte(helperNearestAmpTimeIndex(relsR))

		//TEMPORARY
		// Freq envelope is commented out for now -- but we need the two control bytes at very minimum: adding that here
		// point1
		vce.Envelopes[i].FreqEnvelope.Table[0] = 0
		vce.Envelopes[i].FreqEnvelope.Table[1] = 0
		// special case for point1
		vce.Envelopes[i].FreqEnvelope.Table[2] = 0x80 // matches default from EDATA
		vce.Envelopes[i].FreqEnvelope.Table[3] = 0    // 0 == Sine, octave 0, freq int and amp int disabled

		// END TEMPORARY

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

var BreakPoint = [100]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 9, 9, 9,
	9, 9, 10, 10, 10, 10, 10, 10, 11, 11, 11, 11, 11, 11, 12, 12, 12, 12, 13, 13, 13, 13, 14, 14, 14, 15, 15, 15, 15,
	16, 16, 16, 16, 17, 17, 17, 18, 18, 18, 18, 19, 19, 19, 20, 20, 20, 20, 21, 21, 21, 21, 22, 22, 22, 22,
	23, 23, 23, 23, 24, 24, 24, 24, 25, 25, 25, 25, 26, 26, 26, 26, 27}

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
