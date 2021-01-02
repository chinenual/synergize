package main

import (
	"math"
)

// The following code adapted from Javascript written by Jari Kleimola:

// -- nonlinear level mapping
var _outputLevel = []int{
	0, 5, 9, 13, 17, 20, 23, 25, 27, 29, 31, 33, 35, 37, 39,
	41, 42, 43, 45, 46, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61,
	62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80,
	81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99,
	100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114,
	115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127}

// -- pattern that defines when EG chip is active
// -- four array elements indexed by fractional qr (ie., mod 4)
// -- eight items in each array : indexed by sampleNumber mod 8
var _envMask = [][]int{
	{0, 1, 0, 1, 0, 1, 0, 1},
	{0, 1, 0, 1, 0, 1, 1, 1},
	{0, 1, 1, 1, 0, 1, 1, 1},
	{0, 1, 1, 1, 1, 1, 1, 1}}

// EG is not ticked every sample:
// this function returns false if computation should be skipped for current sample
func _egEnabled(shift int, qr int, sampleCounter *int) bool {
	*sampleCounter += 1
	istep := *sampleCounter
	if shift < 0 {
		sm := (1 << -shift) - 1
		if (istep & sm) != sm {
			return false
		}
		istep >>= -shift
	}
	return _envMask[qr&3][istep&7] != 0
}

// returns duration of isegment in samples
func _computeDurationSamples(isegment int, levels [4]byte, rates [4]byte, sampleCounter *int) (nsamples int) {
	// -- level
	var startLevel int
	var targetLevel int
	if isegment == 0 {
		startLevel = int(levels[3])
		targetLevel = int(levels[0])
	} else {
		startLevel = int(levels[isegment-1])
		targetLevel = int(levels[isegment])
	}
	startLevel = max(0, (_outputLevel[startLevel]<<5)-224)
	targetLevel = max(0, (_outputLevel[targetLevel]<<5)-224)
	rising := targetLevel > startLevel
	// -- rate
	rateScaling := 0
	qr := min(63, rateScaling+((int(rates[isegment])*41)>>6))
	shift := (qr >> 2) - 11
	// -- loop
	level := startLevel
	nsamples = 0
	if rising {
		for level < targetLevel {
			if _egEnabled(shift, qr, sampleCounter) {
				slope := 17 - (level >> 8)
				level += slope << max(shift, 0)
			}
			nsamples += 1
		}
	} else { // decaying
		for level >= targetLevel {
			if _egEnabled(shift, qr, sampleCounter) {
				level -= 1 << max(shift, 0)
			}
			nsamples += 1
		}
	}
	nsamples -= 1
	return
}

// computes envelope duration and returns each segment duration in msecs
func computeDurationsMs(levels [4]byte, rates [4]byte) (ms [4]int) {
	var sampleCounter int
	sampleCounter = 0

	sampleRate := 49096
	for isegment := 0; isegment < 4; isegment++ {
		nsamples := _computeDurationSamples(isegment, levels, rates, &sampleCounter)
		ms[isegment]  = int(math.Round(float64(nsamples) / float64(sampleRate) * 1000))
	}
	return
}
