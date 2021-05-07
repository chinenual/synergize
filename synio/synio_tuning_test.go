package synio

import (
	"fmt"
	"math"
	"testing"

	"gotest.tools/assert"
)

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

func TestScaleFrequencies(t *testing.T) {
	freqs, err := GetTuningFrequencies(tuningParams)
	assert.NilError(t, err)
	intFreqs := scaleFrequencies(freqs)
	for i, f := range intFreqs {
		//assert.Equal(t, f, factoryROMTableValues[i], "i:%d", i)
		if f != factoryROMTableValues[i] {
			//fmt.Printf("    logs [%d] %v - expect %v (diff: %v)\n", i, math.Log2(float64(f)), math.Log2(float64(factoryROMTableValues[i])), math.Log2(float64(f))-math.Log2(float64(factoryROMTableValues[i])))
			fmt.Printf("    [%d] %v - expect %v (diff: %v)\n", i, f, factoryROMTableValues[i], int(f)-int(factoryROMTableValues[i]))

			// our scaling formula doesnt _exactly_ reproduce the factory table -- the lowest octave of the factory
			// table isnt even part of the exponential curve -- and the upper range veers off by up to 4 Hz.
			// Sanity check that the mapping is tight enough to what we have managed to reverse engineer:
			if i <= 6 {
				// no check - this part of the factory ROM is not exponential
			} else if i < 115 {
				assert.Check(t, math.Abs(float64(f)-float64(factoryROMTableValues[i])) <= 1.0, "i:%d", i)
			} else if i < 123 {
				assert.Check(t, math.Abs(float64(f)-float64(factoryROMTableValues[i])) <= 2.0, "i:%d", i)
			} else {
				assert.Check(t, math.Abs(float64(f)-float64(factoryROMTableValues[i])) <= 4.0, "i:%d", i)
			}
		}
	}
}
