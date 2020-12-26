package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/chinenual/synergize/data"
	"github.com/orcaman/writerseeker"
)

// Helper routines that may find their way back into the synergize/data module.

func helperBlankVce() (vce data.VCE, err error) {
	rdr := bytes.NewReader(data.VRAM_EDATA[data.Off_VRAM_EDATA:])
	if vce, err = data.ReadVce(rdr, false); err != nil {
		return
	}
	// re-allocate the Envelopes and each env Table to allow us to control size
	//and #osc simply by writing to VOITAB and NPOINTS params
	for i := 1; i < 16; i++ {
		// make a copy of the first osc:
		vce.Envelopes = append(vce.Envelopes, vce.Envelopes[0])
	}
	for i := 0; i < 16; i++ {
		// now re-allocate each envelope to their max possible length:
		vce.Envelopes[i].AmpEnvelope.Table = make([]byte, 4*16)
		vce.Envelopes[i].FreqEnvelope.Table = make([]byte, 4*16)
	}
	return
}

func helperSetPatchType(vce *data.VCE, patchType int) {
	for i := range data.PatchTypePerOscTable[patchType-1] {
		vce.Envelopes[i].FreqEnvelope.OPTCH = data.PatchTypePerOscTable[patchType-1][i]
	}
}

var dxAlgoNoFeedbackPatchTypePerOscTable = [32][16]byte{
	{100, 97, 97, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 1 (algo 0)
	{100, 97, 97, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 2 (algo 1)
	{100, 97, 33, 100, 97, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},  // DX 3 (algo 2)
	{100, 97, 33, 100, 97, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},  // DX 4 (algo 3)
	{100, 1, 100, 1, 100, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 5 (algo 4)
	{100, 1, 100, 1, 100, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 6 (algo 5)
	{100, 97, 76, 33, 100, 33, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 7 (algo 6)
	{100, 97, 76, 33, 100, 33, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 8 (algo 7)
	{100, 97, 76, 33, 100, 33, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 9 (algo 8)
	{100, 76, 33, 100, 97, 33, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 10 (algo 9)
	{100, 76, 33, 100, 97, 33, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 11 (algo 10)
	{100, 76, 76, 33, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},  // DX 12 (algo 11)
	{100, 76, 76, 33, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},  // DX 13 (algo 12)
	{100, 76, 97, 33, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},  // DX 14 (algo 13)
	{100, 76, 97, 33, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},  // DX 15 (algo 14)
	{100, 33, 100, 1, 4, 32, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 16 (algo 15)
	{100, 33, 100, 1, 4, 32, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 17 (algo 16)
	{100, 97, 33, 4, 4, 32, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},    // DX 18 (algo 17)
	{100, 33, 1, 100, 97, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 19 (algo 18)
	{100, 76, 33, 100, 1, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},   // DX 20 (algo 19)
	{100, 33, 1, 100, 1, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},    // DX 21 (algo 20)
	{100, 1, 1, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 22 (algo 21)
	{100, 33, 1, 100, 33, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 23 (algo 22)
	{100, 33, 1, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 24 (algo 23)
	{100, 33, 1, 36, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 25 (algo 24)
	{100, 76, 33, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 26 (algo 25)
	{100, 76, 33, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 27 (algo 26)
	{36, 100, 97, 1, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 28 (algo 27)
	{100, 33, 100, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 29 (algo 28)
	{36, 100, 97, 33, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 30 (algo 29)
	{68, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 31 (algo 30)
	{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, // DX 32 (algo 31)
}

func helperSetAlgorithmPatchType(vce *data.VCE, dxAlgo byte, dxFeedback byte) (err error) {
	if dxFeedback != 0 {
		log.Printf("WARNING: Limitation: unhandled DX feedback: %d", dxFeedback)
	}
	if len(dxAlgoNoFeedbackPatchTypePerOscTable[dxAlgo]) < 16 {
		log.Printf("ERROR: Limitation: unhandled DX algorithm: %d", dxAlgo)
	}

	for i := range dxAlgoNoFeedbackPatchTypePerOscTable[dxAlgo] {
		vce.Envelopes[i].FreqEnvelope.OPTCH = dxAlgoNoFeedbackPatchTypePerOscTable[dxAlgo][i]
	}

	return
}

func helperCompactVCE(vce data.VCE) (compacted data.VCE, err error) {
	var writebuf = writerseeker.WriterSeeker{}

	if err = data.WriteVce(&writebuf, vce, data.VceName(vce.Head), false); err != nil {
		return
	}
	write_bytes, _ := ioutil.ReadAll(writebuf.Reader())

	var readbuf2 = bytes.NewReader(write_bytes)

	if compacted, err = data.ReadVce(readbuf2, false); err != nil {
		return
	}
	return
}

func helperVCEToJSON(vce data.VCE) (result string) {
	// compact the vce before printing it:
	var err error
	var compacted data.VCE
	if compacted, err = helperCompactVCE(vce); err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}

	b, _ := json.MarshalIndent(compacted, "", "\t")
	result = string(b)
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
