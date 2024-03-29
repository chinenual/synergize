package dx2syn

import (
	"github.com/chinenual/synergize/data"
	"github.com/chinenual/synergize/logger"
	"github.com/pkg/errors"
)

// XREF: keep in sync with data/vce.go: PatchTypePerOscTable
var dxAlgoNoFeedbackPatchTypePerOscTable = [32][16]byte{
	{100, 97, 97, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 1 (algo 0)
	{100, 97, 97, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 2 (algo 1)
	{100, 97, 1, 100, 97, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 3 (algo 2)
	{100, 97, 1, 100, 97, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 4 (algo 3)
	{100, 1, 100, 1, 100, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 5 (algo 4)
	{100, 1, 100, 1, 100, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 6 (algo 5)
	{100, 97, 76, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 7 (algo 6)
	{100, 97, 76, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 8 (algo 7)
	{100, 97, 76, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 9 (algo 8)
	{100, 76, 1, 100, 97, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 10 (algo 9)
	{100, 76, 1, 100, 97, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 11 (algo 10)
	{100, 76, 76, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 12 (algo 11)
	{100, 76, 76, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 13 (algo 12)
	{100, 76, 97, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 14 (algo 13)
	{100, 76, 97, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 15 (algo 14)
	{100, 97, 76, 76, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},    // DX 16 (algo 15) // fixed ....can't reproduce the algo due to lack of register in SYN - this is approximation
	{100, 97, 76, 76, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},    // DX 16 (algo 16) // fixed ....can't reproduce the algo due to lack of register in SYN - this is approximation
	{100, 97, 97, 76, 76, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 18 (algo 17)
	{100, 1, 1, 100, 97, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},    // DX 19 (algo 18)
	{100, 76, 1, 100, 1, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},    // DX 20 (algo 19)
	{100, 1, 1, 100, 1, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},     // DX 21 (algo 20)
	{100, 1, 1, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},     // DX 22 (algo 21)
	{100, 1, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},     // DX 23 (algo 22)
	{100, 1, 1, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},       // DX 24 (algo 23)
	{100, 1, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},       // DX 25 (algo 24)
	{100, 76, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},    // DX 26 (algo 25)
	{100, 76, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},    // DX 27 (algo 26)
	{4, 100, 97, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},    // DX 28 (algo 27)
	{100, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},     // DX 29 (algo 28)
	{4, 100, 97, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},      // DX 30 (algo 29)
	{100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},       // DX 31 (algo 30)
	{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},         // DX 32 (algo 31)
}

func SetAlgorithmPatchType(vce *data.VCE, dxAlgo byte, dxFeedback byte) (err error) {
	if dxFeedback != 0 {
		logger.Debugf("WARNING: Limitation: unhandled DX feedback: %d", dxFeedback)
	}
	if dxAlgo < 0 || dxAlgo > byte(len(dxAlgoNoFeedbackPatchTypePerOscTable)) {
		return errors.Errorf("Invalid Algorithm value %d - expected 0 .. 31", dxAlgo)
	}

	for i := range dxAlgoNoFeedbackPatchTypePerOscTable[dxAlgo] {
		vce.Envelopes[i].FreqEnvelope.OPTCH = dxAlgoNoFeedbackPatchTypePerOscTable[dxAlgo][i]
	}

	return
}
