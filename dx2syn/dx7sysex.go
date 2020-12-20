package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

type Dx7Sysex struct {
	Voices []Dx7Voice
}

type Dx7Osc struct {
	EgRate                    [4]byte
	EgLevel                   [4]byte
	KeyLevelScalingBreakPoint byte
	KeyLevelScalingLeftDepth  byte
	KeyLevelScalingRightDepth byte
	KeyLevelScalingLeftCurve  byte
	KeyLevelScalingRightCurve byte
	KeyRateScaling            byte
	AmpModSensitivity         byte
	KeyVelocitySensitivity    byte
	OperatorOutputLevel       byte
	OscMode                   bool // true == fixed, false == ratio
	OscFreqCoarse             int8
	OscFreqFine               byte
	OscDetune                 byte
}

//
// Nicely documented at https://github.com/asb2m10/dexed/blob/master/Documentation/sysex-format.txt
// NOTE: the file structure is "packed" - we represent each param as a byte even if they are packed
// several params per byte in the file
type Dx7Voice struct {
	Osc                 [6]Dx7Osc // in logic order (not as in the file) - "Osc 1" is Osc[0], "Osc 6" is Osc[5]
	PitchEgRate         [4]byte   // PitchEgRate[0] = rate1, [1] = rate2, etc.
	PitchEgLevel        [4]byte   // PitchEgLevel[0] = level1, [1] = level2, etc
	Algorithm           byte
	Feedback            byte
	OscSync             bool
	LfoSpeed            byte
	LfoDelay            byte
	LfoPitchModDepth    byte
	LfoAmpModDepth      byte
	LfoSync             bool
	Waveform            byte
	PitchModSensitivity byte
	Transpose           byte
	VoiceName           string
}

func readDx7Osc(reader *bytes.Reader) (osc Dx7Osc, err error) {
	if err = binary.Read(reader, binary.LittleEndian, &osc.EgRate); err != nil {
		return
	}
	if err = binary.Read(reader, binary.LittleEndian, &osc.EgLevel); err != nil {
		return
	}
	if err = binary.Read(reader, binary.LittleEndian, &osc.KeyLevelScalingBreakPoint); err != nil {
		return
	}
	if err = binary.Read(reader, binary.LittleEndian, &osc.KeyLevelScalingLeftDepth); err != nil {
		return
	}
	if err = binary.Read(reader, binary.LittleEndian, &osc.KeyLevelScalingRightDepth); err != nil {
		return
	}
	var v byte
	if err = binary.Read(reader, binary.LittleEndian, &v); err != nil {
		return
	}
	osc.KeyLevelScalingLeftCurve = v & 0x03
	osc.KeyLevelScalingRightCurve = (v & 0x0C) >> 2

	if err = binary.Read(reader, binary.LittleEndian, &v); err != nil {
		return
	}
	osc.OscDetune = (v & 0x78) >> 3
	osc.KeyRateScaling = v & 0x07

	if err = binary.Read(reader, binary.LittleEndian, &v); err != nil {
		return
	}
	osc.AmpModSensitivity = v & 0x03
	osc.KeyVelocitySensitivity = (v & 0x1C) >> 3

	if err = binary.Read(reader, binary.LittleEndian, &osc.OperatorOutputLevel); err != nil {
		return
	}
	if err = binary.Read(reader, binary.LittleEndian, &v); err != nil {
		return
	}
	osc.OscMode = (v & 0x01) != 0
	osc.OscFreqCoarse = int8((v & 0x7E) >> 1)

	if err = binary.Read(reader, binary.LittleEndian, &osc.OscFreqFine); err != nil {
		return
	}
	return
}

func readDx7Voice(reader *bytes.Reader) (voice Dx7Voice, err error) {
	for i := 5; i >= 0; i-- {
		var osc Dx7Osc
		if osc, err = readDx7Osc(reader); err != nil {
			return
		}
		voice.Osc[i] = osc
	}
	if err = binary.Read(reader, binary.LittleEndian, &voice.PitchEgRate); err != nil {
		return
	}
	if err = binary.Read(reader, binary.LittleEndian, &voice.PitchEgLevel); err != nil {
		return
	}
	if err = binary.Read(reader, binary.LittleEndian, &voice.Algorithm); err != nil {
		return
	}
	var v byte
	if err = binary.Read(reader, binary.LittleEndian, &v); err != nil {
		return
	}
	voice.OscSync = (v & 0x08) != 0
	voice.Feedback = v & 0x07

	if err = binary.Read(reader, binary.LittleEndian, &voice.LfoSpeed); err != nil {
		return
	}
	if err = binary.Read(reader, binary.LittleEndian, &voice.LfoDelay); err != nil {
		return
	}
	if err = binary.Read(reader, binary.LittleEndian, &voice.LfoPitchModDepth); err != nil {
		return
	}
	if err = binary.Read(reader, binary.LittleEndian, &voice.LfoAmpModDepth); err != nil {
		return
	}
	if err = binary.Read(reader, binary.LittleEndian, &v); err != nil {
		return
	}
	voice.LfoSync = (v & 0x01) != 0
	voice.Waveform = (v & 0x1E) >> 1
	voice.PitchModSensitivity = (v & 0xC0) >> 6

	if err = binary.Read(reader, binary.LittleEndian, &voice.Transpose); err != nil {
		return
	}

	var rawName [10]byte
	if err = binary.Read(reader, binary.LittleEndian, &rawName); err != nil {
		return
	}
	voice.VoiceName = string(rawName[:])
	return
}

func ReadDx7Sysex(pathname string) (sysex Dx7Sysex, err error) {
	var b []byte
	if b, err = ioutil.ReadFile(pathname); err != nil {
		return
	}
	reader := bytes.NewReader(b)

	// validate that the header of the Sysex is a "bulk DX7 sysex":
	var header [6]byte
	expectHeader := []byte{0xF0, 0x43, 0x00, 0x09, 0x20, 0x00}
	if err = binary.Read(reader, binary.LittleEndian, &header); err != nil {
		return
	}
	for i := range expectHeader {
		if expectHeader[i] != header[i] {
			err = errors.Errorf("Invalid Sysex header byte[%d] - expected %2x, saw %2x", i, expectHeader[i], header[i])
		}
	}
	for i := 0; i < 32; i++ {
		var v Dx7Voice
		if v, err = readDx7Voice(reader); err != nil {
			err = errors.Wrapf(err, "Error reading voice[%d]", i)
			return
		}
		sysex.Voices = append(sysex.Voices, v)
	}
	return
}

func Dx7VoiceToJSON(v Dx7Voice) (result string) {
	b, _ := json.MarshalIndent(v, "", "\t")
	result = string(b)
	return
}
