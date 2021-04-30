package dx2syn

import (
	"math"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/chinenual/synergize/logger"
)

// Helper routines that may find their way back into the synergize/data module.

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func sanitizeFilename(v string) (result string) {
	r, _ := regexp.Compile("[^A-Za-z0-9+!@# _-]")
	result = r.ReplaceAllString(v, "_")
	return
}

func MakeVCEFilename(sysexPath string, synVoiceName string) (pathname string, err error) {
	sysexExt := path.Ext(sysexPath)
	sysexDir := strings.TrimRight((sysexPath)[0:len(sysexPath)-len(sysexExt)], " ")
	if err = os.MkdirAll(sysexDir, 0777); err != nil {
		logger.Errorf("Could not create output directory %s: %v\n", sysexDir, err)
		return
	}
	base := path.Join(sysexDir, sanitizeFilename(synVoiceName))
	pathname = base + ".VCE"
	return
}

// translations of the javascript functions in viewVCE_envs.js

var _freqTimeScale = []int{0, 1, 2, 3, 4, 5, 6, 7,
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
var _ampTimeScale = []int{0, 1, 2, 3, 4, 5, 6, 7,
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

var _freqValues = []int{1, 2, 3, 4, 5, 6, 7, 7, 7, 8, 8, 9, 9, 10, 10, 11, 12, 12,
	13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 25, 26, 28, 30, 31, 33, 35, 37, 40,
	42, 44, 47, 50, 53, 56, 60, 63, 67, 71, 75, 80, 84, 89, 94, 100, 106, 112, 119,
	126, 134, 142, 150, 159, 169, 179, 189, 201, 212, 225, 238, 253, 268, 284, 300,
	318, 337, 357, 378, 401, 425, 450, 477, 505, 536, 568, 601, 637, 675, 715, 757,
	803, 850, 901, 955, 1011, 1072, 1135, 1203, 1274, 1350, 1430, 1515, 1605,
	1701, 1802, 1909, 2023, 2143, 2271, 2405, 2548, 2700, 2861, 3031, 3211,
	3402, 3605, 3818, 4046, 4286, 4541, 4811, 5097, 5400, 5722, 6061, 6422, 6804}

func _indexOfNearestValue(val int, array []int) (index int) {
	bestDiff := math.MaxInt32
	index = -1
	// brute force - we look at the whole array rather than return as soon as we find a minima
	// this is fine since the array is known to be short.
	for i, v := range array {
		diff := val - v
		if diff < 0 {
			// diffs are absolute values
			diff *= -1
		}
		if diff < bestDiff {
			bestDiff = diff
			index = i
		}
	}
	return
}

func helperNearestFreqTimeIndex(val int) (index int) {
	return _indexOfNearestValue(val, _freqTimeScale)
}
func helperNearestAmpTimeIndex(val int) (index int) {
	return _indexOfNearestValue(val, _ampTimeScale)
}

func helperNearestFreqValueIndex(val int) (index int) {
	return _indexOfNearestValue(val, _freqValues)
}

// translate a frequency time "as displayed" to "byte value as stored"
func helperUnscaleFreqTimeValue(time int) byte {
	// fixme: linear search is brute force - but the list is short - performance is "ok" as is...
	for i, v := range _freqTimeScale {
		if v >= time {
			return byte(i)
		}
	}
	// shouldnt happen!
	return byte(len(_freqTimeScale))
}

// translate a amplitude time "as displayed" to "byte value as stored"
func helperUnscaleAmpTimeValue(time int) byte {
	// fixme: linear search is brute force - but the list is short - performance is "ok" as is...
	for i, v := range _ampTimeScale {
		if v >= time {
			return byte(i)
		}
	}
	// shouldnt happen!
	return byte(len(_ampTimeScale))
}

// translate a frequency value "as displayed" to "byte value as stored"
func helperUnscaleFreqEnvValue(val byte) byte {
	return val
}

// translate a amplitude value "as displayed" to "byte value as stored"
func UnscaleAmpEnvValue(val byte) byte {
	return val + 55
}

// translate -- does not handle the "RAND" values, only the numeric ones
func helperUnscaleDetune(val int) int8 {
	// XREF: original source in TextToFDETUN() - javascript in viewVCE_voice.js

	// See FDETUNToText.  This "reverses" that attrocity

	if val >= (-32*3) && val <= (32*3) {
		// CASE B
		val /= 3
	} else if val > 0 {
		// CASE C
		val = ((val / 3) + 32) / 2
	} else {
		// CASE D
		val = ((val / 3) - 32) / 2
	}
	return int8(val)
}
