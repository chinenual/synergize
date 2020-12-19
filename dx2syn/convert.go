package main

import (
	"github.com/chinenual/synergize/data"
	"github.com/pkg/errors"
)

func TranslateDx7ToVce(dx7Voice Dx7Voice) (vce data.VCE, err error) {
	if vce, err = helperBlankVce(); err != nil {
		return
	}

	for i := 0; i < 8; i++ {
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

	// harmonics:
	transposedDown := false
	for _,o := range dx7Voice.Osc {
		if o.OscFreqCoarse < 0 {
			transposedDown = true
			vce.Head.VTRANS =- 12
			break
		}
	}
	for i,o := range dx7Voice.Osc {
		vce.Envelopes[i].FreqEnvelope.OHARM = o.OscFreqCoarse
		if transposedDown {
			vce.Envelopes[i].FreqEnvelope.OHARM += 1
		}

		// envelopes: DX amp envelopes always have 4 points
		vce.Envelopes[i].AmpEnvelope.NPOINTS = 4
		// Each Synergy oscillator is voice twice - for low and high key velocity response
		// Synergy envelopes are represented as quads of ValLow, ValHi, RateLow and RateHi
		// set both upper and lower envs the same
		// point1
		vce.Envelopes[i].AmpEnvelope.Table[0] = o.EgLevel[0] - 99
		vce.Envelopes[i].AmpEnvelope.Table[1] = o.EgLevel[0] - 99
		vce.Envelopes[i].AmpEnvelope.Table[2] = o.EgRate[0] - 99
		vce.Envelopes[i].AmpEnvelope.Table[3] = o.EgRate[0] - 99

		//point2
		vce.Envelopes[i].AmpEnvelope.Table[4] = (o.EgLevel[1] - 99)
		vce.Envelopes[i].AmpEnvelope.Table[5] = (o.EgLevel[1] - 99)
		vce.Envelopes[i].AmpEnvelope.Table[6] = (o.EgRate[0] - 99) +
			(o.EgRate[1] - 99)
		vce.Envelopes[i].AmpEnvelope.Table[7] = (o.EgRate[0] - 99) +
			(o.EgRate[1] - 99)

		//point3
		vce.Envelopes[i].AmpEnvelope.Table[8] = (o.EgLevel[2] - 99)
		vce.Envelopes[i].AmpEnvelope.Table[9] = (o.EgLevel[2] - 99)
		vce.Envelopes[i].AmpEnvelope.Table[10] = (o.EgRate[0] - 99) +
			(o.EgRate[1] - 99) + (o.EgRate[2] - 99)
		vce.Envelopes[i].AmpEnvelope.Table[11] = (o.EgRate[0] - 99) +
			(o.EgRate[1] - 99) + (o.EgRate[2] - 99)

		//point4
		vce.Envelopes[i].AmpEnvelope.Table[12] = (o.EgLevel[3] - 99)
		vce.Envelopes[i].AmpEnvelope.Table[13] = (o.EgLevel[3] - 99)
		vce.Envelopes[i].AmpEnvelope.Table[14] = (o.EgRate[0] - 99) +
			(o.EgRate[1] - 99) + (o.EgRate[2] - 99)+ (o.EgRate[3] - 99)
		vce.Envelopes[i].AmpEnvelope.Table[15] = (o.EgRate[0] - 99) +
			(o.EgRate[1] - 99) + (o.EgRate[2] - 99)+ (o.EgRate[3] - 99)

		// DX only has a single frequency envelope - replicate it on each Synergy osc:
		// NOTE the first point in the Synergy freq table is "special" - it stores a "freq.scale and wavetype" instead of rates
		// Like the amp table, the values are stored in quads, two values, two rates per point

		// envelopes: DX freq envelopes always have 4 points
		vce.Envelopes[i].FreqEnvelope.NPOINTS = 4

		// point1
		vce.Envelopes[i].FreqEnvelope.Table[0] = (dx7Voice.PitchEgLevel[0] - 99)
		vce.Envelopes[i].FreqEnvelope.Table[1] = (dx7Voice.PitchEgLevel[0] - 99)
		// special case for point1
		vce.Envelopes[i].FreqEnvelope.Table[2] = 0x80 // matches default from EDATA
		vce.Envelopes[i].FreqEnvelope.Table[3] = 0 // 0 == Sine, octave 0, freq int and amp int disabled

		// point2
		vce.Envelopes[i].FreqEnvelope.Table[4] = (dx7Voice.PitchEgLevel[1] - 99)
		vce.Envelopes[i].FreqEnvelope.Table[5] = (dx7Voice.PitchEgLevel[1] - 99)
		vce.Envelopes[i].FreqEnvelope.Table[6] = (dx7Voice.PitchEgRate[0] - 99) +
			(dx7Voice.PitchEgRate[1] - 99)
		vce.Envelopes[i].FreqEnvelope.Table[7] = (dx7Voice.PitchEgRate[0] - 99) +
			(dx7Voice.PitchEgRate[1] - 99)

		// point3
		vce.Envelopes[i].FreqEnvelope.Table[8] = (dx7Voice.PitchEgLevel[1] - 99)
		vce.Envelopes[i].FreqEnvelope.Table[9] = (dx7Voice.PitchEgLevel[1] - 99)
		vce.Envelopes[i].FreqEnvelope.Table[10] = (dx7Voice.PitchEgRate[0] - 99) +
			(dx7Voice.PitchEgRate[1] - 99) + (dx7Voice.PitchEgRate[2] - 99)
		vce.Envelopes[i].FreqEnvelope.Table[11] = (dx7Voice.PitchEgRate[0] - 99) +
			(dx7Voice.PitchEgRate[1] - 99) + (dx7Voice.PitchEgRate[2] - 99)

		// point4
		vce.Envelopes[i].FreqEnvelope.Table[12] = (dx7Voice.PitchEgLevel[1] - 99)
		vce.Envelopes[i].FreqEnvelope.Table[13] = (dx7Voice.PitchEgLevel[1] - 99)
		vce.Envelopes[i].FreqEnvelope.Table[14] = (dx7Voice.PitchEgRate[0] - 99) +
			(dx7Voice.PitchEgRate[1] - 99) + (dx7Voice.PitchEgRate[2] - 99) + (dx7Voice.PitchEgRate[3] - 99)
		vce.Envelopes[i].FreqEnvelope.Table[15] = (dx7Voice.PitchEgRate[0] - 99) +
			(dx7Voice.PitchEgRate[1] - 99) + (dx7Voice.PitchEgRate[2] - 99) + (dx7Voice.PitchEgRate[3] - 99)
	}
	// ... everything else ...

	// if you need to abort, use:
	//	err = errors.New("an error message")
	//  return

	return
}