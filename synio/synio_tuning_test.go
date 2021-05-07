package synio

import (
	"testing"

	"gotest.tools/assert"
)

func TestScaleFrequenciesStandardTuning(t *testing.T) {
	freqs, tones, _, err := GetTuningFrequencies(tuningParams)
	assert.NilError(t, err)
	assert.Check(t, tuningParams.UseStandardTuning)
	intFreqs := scaleFrequencies(tuningParams, freqs)
	for i, f := range intFreqs {
		assert.Equal(t, f, factoryROMTableValues[i], "i:%d", i)
	}
	for i, tone := range tones {
		assert.Equal(t, tone.Cents, float64(i+1)*100.0, "i:%d", i)
	}
}
