package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/chinenual/synergize/data"
)

// compress the DX7 10-character name to something that fits in the 8-character VNAME.
// Use the supplied map to keep track of names we've already used so we can ensure they are all unique
func convertName(nameMap *map[string]bool, dxName string, vce *data.VCE) {
	newName := _convertName(dxName, 8)
	_, exists := (*nameMap)[strings.ToUpper(newName)]
	if exists {
		prefix := _convertName(dxName, 7)
		suffixes := []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
			'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W'}
		for _, v := range suffixes {
			candidate := prefix + string(v)
			_, exists = (*nameMap)[strings.ToUpper(candidate)]
			if !exists {
				newName = candidate
				break
			}
		}
	}
	(*nameMap)[strings.ToUpper(newName)] = true
	//finally, pad with spaces:
	newName = data.VcePaddedName(newName)
	for i := 0; i < 8; i++ {
		vce.Head.VNAME[i] = newName[i]
	}
	fmt.Printf("DX7 VoiceName: '%s' Synergy VNAME: '%s'\n", dxName, newName)
}

func _convertName(dxName string, length int) string {
	newName := dxName
	// remove leading or trailing spaces:
	newName = strings.Trim(newName, " ")
	//fmt.Printf("AFTER TRIM: '%s'\n", newName)

	for len(newName) > length && strings.Contains(newName, "  ") {
		// compress internal spaces:
		newName = strings.ReplaceAll(newName, "  ", " ")
	}
	//fmt.Printf("AFTER COMPRESS SPACE: '%s'\n", newName)
	for len(newName) > length && strings.Contains(newName, " ") {
		// remove spaces - starting with the last one:
		i := strings.LastIndex(newName, " ")
		newName = newName[:i] + newName[i+1:]
	}
	//fmt.Printf("AFTER INTERNAL SPACE: '%s'\n", newName)
	const punctuation = "_-\\/:;,.<>?!@#$%^&*()=[]{}'\""
	for len(newName) > length && strings.ContainsAny(newName, punctuation) {
		// remove punctuation - starting with the last one:
		i := strings.LastIndexAny(newName, punctuation)
		newName = newName[:i] + newName[i+1:]
	}
	//fmt.Printf("AFTER PUNCTUATION: '%s'\n", newName)
	const vowels = "aeiouAEIOU"
	for len(newName) > length && strings.ContainsAny(newName, vowels) {
		// remove vowels - starting with the last one:
		i := strings.LastIndexAny(newName, vowels)
		fmt.Printf("last idx: '%s' %d\n", newName, i)
		newName = newName[:i] + newName[i+1:]
	}
	//fmt.Printf("AFTER VOWELS: '%s'\n", newName)

	if len(newName) > length {
		newName = newName[0:length]
	}

	return newName
}

func TranslateDx7ToVce(nameMap *map[string]bool, dx7Voice Dx7Voice) (vce data.VCE, err error) {
	if vce, err = helperBlankVce(); err != nil {
		return
	}

	convertName(nameMap, dx7Voice.VoiceName, &vce)

	if err = helperSetAlgorithmPatchType(&vce, dx7Voice.Algorithm, dx7Voice.Feedback); err != nil {
		return
	}
	// DX7 always uses 6 oscillators
	vce.Head.VOITAB = 5

	vce.Head.VASENS = 31
	vce.Head.VTSENS = 31
	vce.Head.VTCENT = 24
	vce.Head.VIBRAT = (dx7Voice.LfoSpeed + 1) / 3
	vce.Head.VIBDEL = dx7Voice.LfoDelay
	vce.Head.VIBDEP = int8(dx7Voice.LfoPitchModDepth)
	vce.Head.VTRANS = int8(dx7Voice.Transpose - 24)

	attkR := 0
	decyR := 0
	sustR := 0
	relsR := 0
	var freqValueByte byte // Synergy byte encoding for freq value
	var freqValueInt int   // scaled freq value
	var harmonic float64
	var OSClevelPercent float64
	var VelocityPercent float64
	var PMfix float64
	var LevComp float64
	var patchOutputDSR byte

	var ms [4]int

	filterIndex := int8(-1) // increment as we allocate each filter

	// transposedown code -- test to be sure in 'ratio' mode
	// Do not transpose if OP is in fixed Freq Mode.
	transposedDown := false
	for _, o := range dx7Voice.Osc {
		if o.OscFreqCoarse == 0 && o.OscMode == false {
			transposedDown = true
			vce.Head.VTRANS = vce.Head.VTRANS - 12
			//fmt.Printf(" %s %d %d \n", " transpose - Ratio mode  ", o, o.OscMode)
			break
		}
	}

	//Lower carriers for volume compnstaon.
	//Find # of carriers  ---  patchOutputDSR = 0
	carrier := 0
	LevComp = 1.0
	for i := range dx7Voice.Osc {
		patchOutputDSR = ((vce.Envelopes[i].FreqEnvelope.OPTCH & 0xc0) >> 6)
		if patchOutputDSR == 0 {
			carrier++
			fmt.Printf(" %s %d \n", " # Carriers = ", carrier)
		}
	}

	//***************************************************************************

	for oscIndex, dxOsc := range dx7Voice.Osc {

		////  fm: y = sin(phase); phase += phaseInc + mod;
		////  pm: y = sin(phase + mod); phase += phaseInc;
		//OPTCH
		//patchOutputDSR = ((patchByte & 0xc0) >> 6)
		patchOutputDSR = ((vce.Envelopes[oscIndex].FreqEnvelope.OPTCH & 0xc0) >> 6)

		//
		//  ********************  Check if OSC is MODULATOR
		//  ********************  If yes, raise "OSClevelPercent" by 5 %  Trying to fake FM/PM difference
		PMfix = 1.0

		if patchOutputDSR > 0 {
			PMfix = 1.05
			//fmt.Printf(" %s %f \n", " PMfix = ", PMfix)
		} else {
			PMfix = 1.0
			//fmt.Printf(" %s %f \n", " PMfix = ", PMfix)
		}
		//fmt.Printf(" %s %d %d %f \n", " OSC - POdsr= ", oscIndex+1, patchOutputDSR, PMfix)

		//OPTCH
		//patchOutputDSR = ((patchByte & 0xc0) >> 6)
		if patchOutputDSR == 0 {
			switch carrier {
			case 6:
				LevComp = .70
			case 5:
				LevComp = .75
			case 4:
				LevComp = .80
			case 3:
				LevComp = .85
			case 2:
				LevComp = .90
			case 1:
				LevComp = 1.0
			}
			//fmt.Printf(" %s %f \n", " Levcomp =  ", LevComp)

		}
		// Fix Over Values where max is 99...
		if dxOsc.KeyLevelScalingBreakPoint > 99 {
			fmt.Printf(" %s %d \n", " BP before =  ", dxOsc.KeyLevelScalingBreakPoint)
			dxOsc.KeyLevelScalingBreakPoint = byte(math.Round(float64(dxOsc.KeyLevelScalingBreakPoint) * 0.727))
			fmt.Printf(" %s %d \n", " BP after =  ", dxOsc.KeyLevelScalingBreakPoint)
		}
		if dxOsc.KeyLevelScalingRightDepth > 99 {
			fmt.Printf(" %s %d \n", " RT before =  ", dxOsc.KeyLevelScalingRightDepth)
			dxOsc.KeyLevelScalingRightDepth = byte(math.Round(float64(dxOsc.KeyLevelScalingRightDepth) * 0.727))
			fmt.Printf(" %s %d \n", " RT after =  ", dxOsc.KeyLevelScalingRightDepth)
		}
		if dxOsc.KeyLevelScalingLeftDepth > 99 {
			fmt.Printf(" %s %d \n", " LT before =  ", dxOsc.KeyLevelScalingLeftDepth)
			dxOsc.KeyLevelScalingLeftDepth = byte(math.Round(float64(dxOsc.KeyLevelScalingLeftDepth) * 0.727))
			fmt.Printf(" %s %d \n", " LT after =  ", dxOsc.KeyLevelScalingLeftDepth)
		}
		if dxOsc.OscFreqFine > 99 {
			fmt.Printf(" %s %d \n", " Fine before=  ", dxOsc.OscFreqFine)
			dxOsc.OscFreqFine = byte(math.Round(float64(dxOsc.OscFreqFine) * 0.727))
			fmt.Printf(" %s %d \n", " Fine after =  ", dxOsc.OscFreqFine)
		}

		// ******************************************************************************************
		// *************** put key scaling in filter B  filter B[0 - 31] for each OSC  **************
		// ********************* pretty low res but better than nothing  ****************************
		// ******************************************************************************************
		//Activate FILTER B above per voice above (in Header)

		if dxOsc.KeyLevelScalingLeftDepth == 0 && dxOsc.KeyLevelScalingRightDepth == 0 {
			// optimization: if key scaling depths are zero, don't waste space for a filter
			vce.Head.FILTER[oscIndex] = 0
		} else {
			filterIndex++
			if err = convertKeyLevelScalingToFilters(oscIndex, dxOsc, filterIndex, &vce); err != nil {
				return
			}
		}

		// Set OSC mode     false = ratio   true = Fixed
		if dxOsc.OscMode == false { //Ratio Mode
			//vce.Envelopes[oscIndex].FreqEnvelope.OHARM = dxOsc.OscFreqCoarse

			// ***************************************************************
			// TO DO  :::::  Add more code to take DX7 Fine into consideration
			// Fine is 100 steps, including first 0 step.
			// Each step is defined by DX7 OscFreqCoarse number.
			//      DX OscFreqCoarse 0 = 1/2 octave below 1 so each step is 1/2 cent.
			//      DX OscFreqCoarse 1 means each FINE is 1 cent,
			//      Dx OscFreqCoarse 2 means each FINE is 2 cents, and so on
			//
			//
			// FINE is a frequency that is added to OscFreqCoarse Freq in Fixed Mode
			// ***************************************************************
			//DX7 OP OscFreqCoarse == 0 means .5 1 octave below 1, which synergy does not have
			// Transpose down takes care of that

			if transposedDown == false && dxOsc.OscFreqFine == 0 { //No harmonic changes
				fmt.Printf(" %d %s  \n", oscIndex, " transdowndown false == 0 Standard Harm")
				vce.Envelopes[oscIndex].FreqEnvelope.OHARM = dxOsc.OscFreqCoarse
				//freqValueInt = int(math.Round(fineValues[dxOsc.OscFreqFine]))
				//freqValueByte = byte(helperNearestFreqValueIndex(freqValueInt))

			} else if transposedDown == false && dxOsc.OscFreqFine != 0 {
				fmt.Printf(" %d %s  \n", oscIndex, " transdown false !=0      Harm + Fine")
				vce.Envelopes[oscIndex].FreqEnvelope.OHARM = dxOsc.OscFreqCoarse
				freqValueInt = int(math.Round(fineValues[dxOsc.OscFreqFine] * float64(dxOsc.OscFreqCoarse)))
				freqValueByte = byte(helperNearestFreqValueIndex(freqValueInt))

			} else if transposedDown == true && dxOsc.OscFreqCoarse == 0 && dxOsc.OscFreqFine == 0 {
				vce.Envelopes[oscIndex].FreqEnvelope.OHARM = 1
				//fmt.Printf(" %d %s  \n", oscIndex, " Coarse = 0  Fine 0 Harm = 1")
			} else if transposedDown == true && dxOsc.OscFreqCoarse == 0 && dxOsc.OscFreqFine == 50 {
				//fmt.Printf(" %d %s  \n", oscIndex, " in true 0 50")
				harmonic = float64(dxOsc.OscFreqCoarse)
				harmonic = harmonic + .5
				harmonic = harmonic * 2
				vce.Envelopes[oscIndex].FreqEnvelope.OHARM = int8(harmonic)

			} else if transposedDown == true && dxOsc.OscFreqCoarse == 0 && dxOsc.OscFreqFine != 50 {
				// TO DO
				// Calculate combined freq of coarse and fine, then double for octave.
				// Find closest harm to the OCT freq, then find remainder,
				// and add the (remainder * 2) in as freqValueInt
				//if dxOsc.OscFreqFine < 50 {
				fmt.Printf(" %d %s  \n", oscIndex, " in true 0 !=50")
				vce.Envelopes[oscIndex].FreqEnvelope.OHARM = 1
				freqValueInt = 0
				freqValueInt = int(math.Round(1 + fineValues[dxOsc.OscFreqFine]*2))
				freqValueByte = byte(helperNearestFreqValueIndex(freqValueInt))

			} else if transposedDown == true && dxOsc.OscFreqCoarse != 0 && dxOsc.OscFreqFine == 50 {
				fmt.Printf(" %d %s  \n", oscIndex, " in true !0 = 50")
				harmonic = float64(dxOsc.OscFreqCoarse)
				harmonic = harmonic + .5
				harmonic = harmonic * 2
				vce.Envelopes[oscIndex].FreqEnvelope.OHARM = int8(harmonic)

			} else if transposedDown == true && dxOsc.OscFreqCoarse != 0 && dxOsc.OscFreqFine != 50 {
				fmt.Printf(" %d %s  \n", oscIndex, " in true !0 !50  ")
				vce.Envelopes[oscIndex].FreqEnvelope.OHARM = dxOsc.OscFreqCoarse * 2
				freqValueInt = int(math.Round(fineValues[dxOsc.OscFreqFine] * float64(dxOsc.OscFreqCoarse*2)))
				freqValueByte = byte(helperNearestFreqValueIndex(freqValueInt))
				//  have to add in FINE
			} else {
				fmt.Printf(" %d %s  \n", oscIndex, " in else")
				vce.Envelopes[oscIndex].FreqEnvelope.OHARM = dxOsc.OscFreqCoarse * 2
			} // Add Coarse <>0  fine = 50

		} else { //Fixed Mode
			vce.Envelopes[oscIndex].FreqEnvelope.OHARM = -12

			// test freq code
			// for k := -100; k >= -127; k-- {
			//	 freqValueInt = int(14764 * math.Pow(2, float64(k-127)/12))
			//	 fmt.Printf(" %s %d %d  \n", "******* m  FRQACC = ", k, freqValueInt)
			// }

			switch dxOsc.OscFreqCoarse {
			case 0, 4, 8, 12, 16, 20, 24, 28:
				// base freq = 1
				{
					freqValueInt = 0
					freqValueInt = int(math.Round(1 + fineValues[dxOsc.OscFreqFine]))
					fmt.Printf(" %s %d  \n", " Fixed freq = ", freqValueInt)
				}

			case 1, 5, 9, 13, 17, 21, 25, 29:
				// freq = 10
				{
					freqValueInt = 0
					freqValueInt = int(math.Round(10 + (fineValues[dxOsc.OscFreqFine] * 10)))
					fmt.Printf(" %s %d  \n", " Fixed freq = ", freqValueInt)
				}
			case 2, 6, 10, 14, 18, 22, 26, 30:
				// freq = 100
				{
					freqValueInt = 0
					freqValueInt = int(math.Round(100 + (fineValues[dxOsc.OscFreqFine] * 100)))
					fmt.Printf(" %s %d  \n", " Fixed freq = ", freqValueInt)
				}
			case 3, 7, 11, 15, 19, 23, 27, 31:
				// freq = 1000
				{
					freqValueInt = 0
					freqValueInt = int(math.Round(1000 + (fineValues[dxOsc.OscFreqFine] * 1000)))
					fmt.Printf(" %s %d  \n", " Fixed freq = ", freqValueInt)
				}
			}
			freqValueByte = byte(helperNearestFreqValueIndex(freqValueInt))
			fmt.Printf(" %d %s %d %d %d  \n \n", oscIndex+1, "  Fixed freq = ", dxOsc.OscFreqFine, freqValueInt, freqValueByte)
		}

		// Set OSC detune
		vce.Envelopes[oscIndex].FreqEnvelope.FDETUN = helperUnscaleDetune(int(dTune[int(dxOsc.OscDetune)]))
		//fmt.Printf(" %s %d %d \n", " Detune = ", dxOsc.OscDetune, vce.Envelopes[oscIndex].FreqEnvelope.FDETUN)

		// type = 1  : no loop (and LOOPPT and SUSTAINPT are accelleration rates not point positions)
		// type = 2  : S only
		// type = 3  : L and S - L must be before S
		// type = 4  : R and S - R must be before S
		// WARNING: when type1, the LOOPPT and SUSTAINPT values are _acceleration_ rates, not point positions. What a pain.
		// Set for Sustain point only.
		vce.Envelopes[oscIndex].AmpEnvelope.ENVTYPE = 2
		//  Always DX7 Sustain Point
		vce.Envelopes[oscIndex].AmpEnvelope.SUSTAINPT = 3
		// envelopes: DX amp envelopes always have 4 points
		vce.Envelopes[oscIndex].AmpEnvelope.NPOINTS = 4
		// Freq Env Sustain point for DX& FINE
		vce.Envelopes[oscIndex].FreqEnvelope.ENVTYPE = 2
		vce.Envelopes[oscIndex].FreqEnvelope.SUSTAINPT = 1
		vce.Envelopes[oscIndex].FreqEnvelope.NPOINTS = 2
		// set lower Env levels for velocity sensitivity
		VelocityPercent = 1.0 - (float64(dxOsc.KeyVelocitySensitivity) / 10.0)
		//fmt.Printf(" %s %f %d \n", " Vel% = ", VelocityPercent, dxOsc.KeyVelocitySensitivity)

		//  change levels for velocity sensitivity, PM Fix, and level comp
		OSClevelPercent = float64(float64(dxOsc.OperatorOutputLevel) / 99.00)

		for k := 0; k < 4; k++ {
			dxOsc.EgLevel[k] = byte(float64(dxOsc.EgLevel[k]) * OSClevelPercent * PMfix * LevComp * 0.727)
		}

		// Each Synergy oscillator is voice twice - for low and high key velocity response
		// Synergy envelopes are represented as quads of ValLow, ValHi, RateLow and RateHi
		// set both upper and lower envs the same

		ms = computeDurationsMs(dxOsc.EgLevel, dxOsc.EgRate)
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
		vce.Envelopes[oscIndex].AmpEnvelope.Table[0] = byte((math.Round(float64(dxOsc.EgLevel[0]) * VelocityPercent)) + 55)
		vce.Envelopes[oscIndex].AmpEnvelope.Table[1] = byte((math.Round(float64(dxOsc.EgLevel[0]))) + 55)
		vce.Envelopes[oscIndex].AmpEnvelope.Table[2] = byte(helperNearestAmpTimeIndex(attkR))
		vce.Envelopes[oscIndex].AmpEnvelope.Table[3] = byte(helperNearestAmpTimeIndex(attkR))

		//point2
		vce.Envelopes[oscIndex].AmpEnvelope.Table[4] = byte((math.Round(float64(dxOsc.EgLevel[1]) * VelocityPercent)) + 55)
		vce.Envelopes[oscIndex].AmpEnvelope.Table[5] = byte((math.Round(float64(dxOsc.EgLevel[1]))) + 55)
		vce.Envelopes[oscIndex].AmpEnvelope.Table[6] = byte(helperNearestAmpTimeIndex(decyR))
		vce.Envelopes[oscIndex].AmpEnvelope.Table[7] = byte(helperNearestAmpTimeIndex(decyR))

		//point3
		vce.Envelopes[oscIndex].AmpEnvelope.Table[8] = byte((math.Round(float64(dxOsc.EgLevel[2]) * VelocityPercent)) + 55)
		vce.Envelopes[oscIndex].AmpEnvelope.Table[9] = byte((math.Round(float64(dxOsc.EgLevel[2]))) + 55)
		vce.Envelopes[oscIndex].AmpEnvelope.Table[10] = byte(helperNearestAmpTimeIndex(sustR))
		vce.Envelopes[oscIndex].AmpEnvelope.Table[11] = byte(helperNearestAmpTimeIndex(sustR))

		//point4
		vce.Envelopes[oscIndex].AmpEnvelope.Table[12] = byte((math.Round(float64(dxOsc.EgLevel[3]) * VelocityPercent)) + 55)
		vce.Envelopes[oscIndex].AmpEnvelope.Table[13] = byte((math.Round(float64(dxOsc.EgLevel[3]))) + 55)
		vce.Envelopes[oscIndex].AmpEnvelope.Table[14] = byte(helperNearestAmpTimeIndex(relsR))
		vce.Envelopes[oscIndex].AmpEnvelope.Table[15] = byte(helperNearestAmpTimeIndex(relsR))

		//TEMPORARY
		// Freq envelope is commented out for now -- but we need the two control bytes at very minimum: adding that here
		// point1

		//vce.Envelopes[oscIndex].FreqEnvelope.Table[0] = 93
		//vce.Envelopes[oscIndex].FreqEnvelope.Table[1] = 93
		//vce.Envelopes[oscIndex].FreqEnvelope.Table[0] = freqValueByte
		//vce.Envelopes[oscIndex].FreqEnvelope.Table[1] = freqValueByte

		// special case for point1
		//vce.Envelopes[oscIndex].FreqEnvelope.Table[2] = 0x80 // matches default from EDATA
		//vce.Envelopes[oscIndex].FreqEnvelope.Table[3] = 0    // 0 == Sine, octave 0, freq int and amp int disabled

		// DX only has a single frequency envelope - replicate it on each Synergy osc:
		// NOTE the first point in the Synergy freq table is "special" - it stores a "freq.scale and wavetype" instead of rates
		// Like the amp table, the values are stored in quads, two values, two rates per point
		/* ****************************************
		// envelopes: DX freq envelopes always have 4 points
		vce.Envelopes[oscIndex].FreqEnvelope.NPOINTS = 4
		*/
		// point1
		vce.Envelopes[oscIndex].FreqEnvelope.Table[0] = freqValueByte
		vce.Envelopes[oscIndex].FreqEnvelope.Table[1] = freqValueByte
		// special case for point1
		vce.Envelopes[oscIndex].FreqEnvelope.Table[2] = 0x80 // matches default from EDATA
		vce.Envelopes[oscIndex].FreqEnvelope.Table[3] = 0    // 0 == Sine, octave 0, freq int and amp int disabled

		// point2
		vce.Envelopes[oscIndex].FreqEnvelope.Table[4] = freqValueByte
		vce.Envelopes[oscIndex].FreqEnvelope.Table[5] = freqValueByte
		vce.Envelopes[oscIndex].FreqEnvelope.Table[6] = 30
		vce.Envelopes[oscIndex].FreqEnvelope.Table[7] = 30
		/*
			// point3
			vce.Envelopes[oscIndex].FreqEnvelope.Table[8] = byte((math.Round(float64(dx7Voice.PitchEgLevel[2]))) + 55)
			vce.Envelopes[oscIndex].FreqEnvelope.Table[9] = byte((math.Round(float64(dx7Voice.PitchEgLevel[2]))) + 55)
			vce.Envelopes[oscIndex].FreqEnvelope.Table[10] = (99 - dx7Voice.PitchEgRate[0]) +
				(99 - dx7Voice.PitchEgRate[1]) + (99 - dx7Voice.PitchEgRate[2])
			vce.Envelopes[oscIndex].FreqEnvelope.Table[11] = (99 - dx7Voice.PitchEgRate[0]) +
				(99 - dx7Voice.PitchEgRate[1]) + (99 - dx7Voice.PitchEgRate[2])

			// point4
			vce.Envelopes[oscIndex].FreqEnvelope.Table[12] = byte((math.Round(float64(dx7Voice.PitchEgLevel[3]))) + 55)
			vce.Envelopes[oscIndex].FreqEnvelope.Table[13] = byte((math.Round(float64(dx7Voice.PitchEgLevel[3]))) + 55)
			vce.Envelopes[oscIndex].FreqEnvelope.Table[14] = (99 - dx7Voice.PitchEgRate[0]) +
				(99 - dx7Voice.PitchEgRate[1]) + (99 - dx7Voice.PitchEgRate[2]) + (99 - dx7Voice.PitchEgRate[3])
			vce.Envelopes[oscIndex].FreqEnvelope.Table[15] = (99 - dx7Voice.PitchEgRate[0]) +
				(99 - dx7Voice.PitchEgRate[1]) + (99 - dx7Voice.PitchEgRate[2]) + (99 - dx7Voice.PitchEgRate[3])
		*/
	}

	// if you need to abort, use:
	//	err = errors.New("an error message")
	//  return

	return
}

func convertKeyLevelScalingToFilters(oscIndex int, dxOsc Dx7Osc, filterIndex int8, vce *data.VCE) (err error) {
	//set filter B on for voice, b-filters are indicated by the 1-based osc index
	vce.Head.FILTER[oscIndex] = filterIndex + 1

	// Assumes no A-filter - so B filter for osc 1 (index 0) is always stored at 0:
	vce.Filters[filterIndex][(breakPoint[dxOsc.KeyLevelScalingBreakPoint])] = 0 //KEY to FREQ Array is breakPoint[] (below)

	// Scale from DX7 0 to 99 to Syn -64 to 63    //using DX 50 = 0

	lMax := float64(dxOsc.KeyLevelScalingLeftDepth) * 0.63 //Assuming the DX7 is in Db also)
	rMax := float64(dxOsc.KeyLevelScalingRightDepth) * 0.63

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

	switch dxOsc.KeyLevelScalingLeftCurve { //0=-LIN, -EXP, +EXP, +LIN
	case 0:
		//-linear from -lMax to 0
		slope := lMax / float64(breakPoint[dxOsc.KeyLevelScalingBreakPoint]-0)
		for k := byte(0); k < breakPoint[dxOsc.KeyLevelScalingBreakPoint]; k++ {
			vce.Filters[filterIndex][k] = int8(math.Round(-lMax + slope*float64(k)))
		}
	case 1:
		//-EXP from -lMax to 0
		for k := byte(0); k < breakPoint[dxOsc.KeyLevelScalingBreakPoint]; k++ {
			x := float64(breakPoint[dxOsc.KeyLevelScalingBreakPoint] - k)
			vce.Filters[filterIndex][k] = int8(math.Pow(expBase, x/expDivisor) / expScale * -lMax)
		}
	case 2:
		//EXP from lMax to 0
		for k := byte(0); k < breakPoint[dxOsc.KeyLevelScalingBreakPoint]; k++ {
			x := float64(breakPoint[dxOsc.KeyLevelScalingBreakPoint] - k)
			vce.Filters[filterIndex][k] = int8(math.Pow(expBase, x/expDivisor) / expScale * lMax)
		}
	case 3:
		//linear from lMax to 0
		slope := -lMax / float64(breakPoint[dxOsc.KeyLevelScalingBreakPoint]-0)
		for k := byte(0); k < breakPoint[dxOsc.KeyLevelScalingBreakPoint]; k++ {
			vce.Filters[filterIndex][k] = int8(math.Round(lMax + slope*float64(k)))
		}
	}

	switch dxOsc.KeyLevelScalingRightCurve { //0=-LIN, -EXP, +EXP, +LIN
	case 0:
		// -Linear from 0 to -rMax
		slope := -rMax / float64(32-breakPoint[dxOsc.KeyLevelScalingBreakPoint])
		for k := breakPoint[dxOsc.KeyLevelScalingBreakPoint] + 1; k < 32; k++ {
			vce.Filters[filterIndex][k] = int8(math.Round(slope * float64(k-breakPoint[dxOsc.KeyLevelScalingBreakPoint])))
		}
	case 1:
		// -EXP from 0 to -rMax
		for k := breakPoint[dxOsc.KeyLevelScalingBreakPoint] + 1; k < 32; k++ {
			x := float64(k - breakPoint[dxOsc.KeyLevelScalingBreakPoint])
			vce.Filters[filterIndex][k] = int8(math.Pow(expBase, x/expDivisor) / expScale * -rMax)
		}
	case 2:
		// EXP from 0 to rMax
		for k := breakPoint[dxOsc.KeyLevelScalingBreakPoint] + 1; k < 32; k++ {
			x := float64(k - breakPoint[dxOsc.KeyLevelScalingBreakPoint])
			vce.Filters[filterIndex][k] = int8(math.Pow(expBase, x/expDivisor) / expScale * rMax)
		}
	case 3:
		// Linear from 0 to rMax
		slope := rMax / float64(32-breakPoint[dxOsc.KeyLevelScalingBreakPoint])
		for k := breakPoint[dxOsc.KeyLevelScalingBreakPoint] + 1; k < 32; k++ {
			vce.Filters[filterIndex][k] = int8(math.Round(slope * float64(k-breakPoint[dxOsc.KeyLevelScalingBreakPoint])))
		}
	}
	return
}

var fineValues = [100]float64{0.00000, 0.02329, 0.04713, 0.07152, 0.09648, 0.12202, 0.14815,
	0.17490, 0.20226, 0.23027, 0.25893, 0.28825, 0.31826, 0.34896,
	0.38038, 0.41254, 0.44544, 0.47911, 0.51356, 0.54882, 0.58489,
	0.62181, 0.65959, 0.69824, 0.73780, 0.77828, 0.81970, 0.86209,
	0.90546, 0.94984, 0.99526, 1.04174, 1.08930, 1.13796, 1.18776,
	1.23872, 1.29087, 1.34423, 1.39882, 1.45471, 1.51189, 1.57040,
	1.63027, 1.69153, 1.75423, 1.81838, 1.88403, 1.95121, 2.01995,
	2.09030, 2.16228, 2.23594, 2.31131, 2.38844, 2.46737, 2.54813,
	2.63078, 2.71536, 2.80189, 2.89045, 2.88107, 3.07380, 3.16869,
	3.26580, 3.36516, 3.46684, 3.57088, 3.67735, 3.78630, 3.89779,
	4.01187, 4.12861, 4.24808, 4.37032, 4.49541, 4.62341, 4.75440,
	4.88844, 5.02560, 5.16595, 5.30957, 5.45654, 5.60693, 5.76083,
	5.91831, 6.07946, 6.24436, 6.41310, 6.58578, 6.76247, 6.94328,
	7.12831, 7.31764, 7.51138, 7.70964, 7.91251, 8.12011, 8.33254,
	8.54993, 8.77237}

var dTune = [16]int8{-81, -75, -69, -60, -42, -30, -24, 0, 9, 21, 24, 33, 45, 66, 72, 81}

var breakPoint = [100]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 9, 9, 9,
	9, 9, 10, 10, 10, 10, 10, 10, 11, 11, 11, 11, 11, 11, 12, 12, 12, 12, 13, 13, 13, 13, 14, 14, 14, 15, 15, 15, 15,
	16, 16, 16, 16, 17, 17, 17, 18, 18, 18, 18, 19, 19, 19, 20, 20, 20, 20, 21, 21, 21, 21, 22, 22, 22, 22,
	23, 23, 23, 23, 24, 24, 24, 24, 25, 25, 25, 25, 26, 26, 26, 26, 27}

var dxRisetoSYN = [100]byte{53, 50, 47, 44, 41, 38, 36, 34, 31, 29, 27, 26, 25, 24, 23, 22, 22,
	21, 21, 20, 20, 19, 19, 18, 18, 18, 17, 17, 17, 17, 17, 17, 16, 16, 16, 16, 16, 16, 16, 16,
	16, 16, 16, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

var dxDecaytoSYN = [100]byte{73, 72, 71, 70, 68, 67, 67, 66, 65, 64, 63, 62, 61, 60, 59, 58, 57, 57,
	56, 54, 50, 46, 42, 39, 35, 33, 31, 30, 28, 26, 25, 24, 23, 23, 22, 22, 21, 21, 21, 21, 21,
	20, 20, 20, 20, 19, 19, 18, 18, 17, 17, 17, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16,
	16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0}

var egRateDecay = [100]int{318, 284, 249, 215, 181, 167, 155, 141, 128, 115, 105, 94, 84, 73, 63,
	58, 54, 50, 44, 4, 35, 32, 28, 24, 20, 18, 16, 15, 13, 11, 10, 9, 8, 7, 7, 7, 7, 67, 6, 6, 6,
	6, 6, 5, 4, 4, 4, 3, 3, 3, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

var egRateRise = [100]int{38, 35, 32, 29, 26, 23, 21, 18, 16, 14, 12, 11, 10, 9, 8,
	8, 7, 6, 6, 5, 5, 4, 4, 4, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
