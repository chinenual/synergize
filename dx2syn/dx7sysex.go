package dx2syn

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"io"
	"io/ioutil"

	"github.com/chinenual/synergize/logger"

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
	Osc [6]Dx7Osc // in file order (not logical order) - "Osc 6" is Osc[0], "Osc 1" is Osc[5]
	// This simplifies mapping to Synergy oscillators since Synergy patches are reversed from DX conventions
	// DX has higher numbered ops modulating lower numbered ones; Synergy has the opposite convention
	// Lower numbered syn osc modulates higher numbered oscs.

	PitchEgRate         [4]byte // PitchEgRate[0] = rate1, [1] = rate2, etc.
	PitchEgLevel        [4]byte // PitchEgLevel[0] = level1, [1] = level2, etc
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

func readDx7Osc(reader *bytes.Reader, packed bool) (osc Dx7Osc, err error) {
	var v byte
	if err = binary.Read(reader, binary.LittleEndian, &osc.EgRate); err != nil { //osc.EgRate
		return
	}
	for a := 0; a < 4; a++ {
		osc.EgRate[a] = osc.EgRate[a] & 0x7f
	}

	if err = binary.Read(reader, binary.LittleEndian, &osc.EgLevel); err != nil { //osc.EgLevel
		return
	}
	for a := 0; a < 4; a++ {
		osc.EgLevel[a] = osc.EgLevel[a] & 0x7f
	}

	if err = binary.Read(reader, binary.LittleEndian, &v); err != nil { //osc.KeyLevelScalingBreakPoint
		return
	}

	if err = binary.Read(reader, binary.LittleEndian, &v); err != nil { //osc.KeyLevelScalingLeftDepth
		return
	}

	if err = binary.Read(reader, binary.LittleEndian, &v); err != nil { //osc.KeyLevelScalingRightDepth
		return
	}

	if packed {
		if err = binary.Read(reader, binary.LittleEndian, &v); err != nil {
			return
		}
		osc.KeyLevelScalingLeftCurve = v & 0x03
		osc.KeyLevelScalingRightCurve = (v & 0x0C) >> 2
	} else {
		if err = binary.Read(reader, binary.LittleEndian, &osc.KeyLevelScalingLeftCurve); err != nil {
			return
		}

		if err = binary.Read(reader, binary.LittleEndian, &osc.KeyLevelScalingRightCurve); err != nil {
			return
		}
	}

	if packed {
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

		if err = binary.Read(reader, binary.LittleEndian, &v); err != nil { //&osc.OperatorOutputLevel
			return
		}

		if err = binary.Read(reader, binary.LittleEndian, &v); err != nil {
			return
		}
		osc.OscMode = (v & 0x01) != 0
		osc.OscFreqCoarse = int8((v & 0x7E) >> 1)

		if err = binary.Read(reader, binary.LittleEndian, &v); err != nil { //&osc.OscFreqFine
			return
		}
	} else {
		// Groan - in the single voice "unpacked" case, it's not just a
		// matter of not packing multiple things into same byte.   Values are
		// also in a different order...
		if err = binary.Read(reader, binary.LittleEndian, &osc.KeyRateScaling); err != nil {
			return
		}
		if err = binary.Read(reader, binary.LittleEndian, &osc.AmpModSensitivity); err != nil {
			return
		}
		if err = binary.Read(reader, binary.LittleEndian, &osc.KeyVelocitySensitivity); err != nil {
			return
		}
		if err = binary.Read(reader, binary.LittleEndian, &osc.OperatorOutputLevel); err != nil {
			return
		}
		if err = binary.Read(reader, binary.LittleEndian, &v); err != nil {
			return
		}
		osc.OscMode = (v & 0x01) != 0
		if err = binary.Read(reader, binary.LittleEndian, &osc.OscFreqCoarse); err != nil {
			return
		}
		if err = binary.Read(reader, binary.LittleEndian, &osc.OscFreqFine); err != nil {
			return
		}
		if err = binary.Read(reader, binary.LittleEndian, &osc.OscDetune); err != nil {
			return
		}
	}

	// mask off invalid values - common to both packed and unpacked
	osc.KeyLevelScalingBreakPoint &= 0x7f
	osc.KeyLevelScalingLeftDepth &= 0x7f
	osc.KeyLevelScalingRightDepth &= 0x7f
	osc.KeyLevelScalingLeftCurve &= 0x03
	osc.KeyLevelScalingRightCurve &= 0x03
	osc.KeyRateScaling &= 0x07
	osc.AmpModSensitivity &= 0x03
	osc.KeyVelocitySensitivity &= 0x07
	osc.OperatorOutputLevel &= 0x7f
	osc.OscFreqCoarse &= 0x1f
	osc.OscFreqFine &= 0x7f
	osc.OscDetune &= 0x0f

	return
}

func readDx7Voice(reader *bytes.Reader, packed bool) (voice Dx7Voice, err error) {
	var v byte
	for i := 0; i < 6; i++ {
		var osc Dx7Osc
		if osc, err = readDx7Osc(reader, packed); err != nil {
			return
		}
		voice.Osc[i] = osc
	}

	if err = binary.Read(reader, binary.LittleEndian, &voice.PitchEgRate); err != nil {
		return
	}
	for a := 0; a < 4; a++ {
		voice.PitchEgRate[a] &= 0x7f
	}

	if err = binary.Read(reader, binary.LittleEndian, &voice.PitchEgLevel); err != nil {
		return
	}
	for a := 0; a < 4; a++ {

		voice.PitchEgLevel[a] &= 0x7f
	}

	if err = binary.Read(reader, binary.LittleEndian, &voice.Algorithm); err != nil {
		return
	}
	// some files have bogus values for Algorithm.  Seems effective to just mask off
	// the upper order bits.  Do that rather than reject them (Dexed does something similar).
	voice.Algorithm &= 0x1f

	if packed {
		if err = binary.Read(reader, binary.LittleEndian, &v); err != nil {
			return
		}
		voice.OscSync = (v & 0x08) != 0
		voice.Feedback = v & 0x07
	} else {
		if err = binary.Read(reader, binary.LittleEndian, &voice.Feedback); err != nil {
			return
		}
		voice.Feedback &= 0x07
		if err = binary.Read(reader, binary.LittleEndian, &v); err != nil {
			return
		}
		voice.OscSync = (v & 0x08) != 0
	}

	if err = binary.Read(reader, binary.LittleEndian, &voice.LfoSpeed); err != nil { //voice.LfoSpeed
		return
	}
	voice.LfoSpeed &= 0x7f

	if err = binary.Read(reader, binary.LittleEndian, &voice.LfoDelay); err != nil { //voice.LfoDelay
		return
	}
	voice.LfoDelay &= 0x7f

	if err = binary.Read(reader, binary.LittleEndian, &voice.LfoPitchModDepth); err != nil { //voice.LfoPitchModDepth
		return
	}
	voice.LfoPitchModDepth &= 0x7f

	if err = binary.Read(reader, binary.LittleEndian, &voice.LfoAmpModDepth); err != nil { //voice.LfoAmpModDepth
		return
	}
	voice.LfoAmpModDepth &= 0x7f

	if packed {
		if err = binary.Read(reader, binary.LittleEndian, &v); err != nil {
			return
		}
		voice.LfoSync = (v & 0x01) != 0
		voice.Waveform = (v & 0x1E) >> 1
		voice.PitchModSensitivity = (v & 0xC0) >> 6
	} else {
		if err = binary.Read(reader, binary.LittleEndian, &v); err != nil {
			return
		}
		voice.LfoSync = (v & 0x01) != 0
		if err = binary.Read(reader, binary.LittleEndian, &voice.Waveform); err != nil {
			return
		}
		voice.Waveform &= 0x05
		if err = binary.Read(reader, binary.LittleEndian, &voice.PitchModSensitivity); err != nil {
			return
		}
		voice.PitchModSensitivity &= 0x07
	}

	if err = binary.Read(reader, binary.LittleEndian, &voice.Transpose); err != nil { //voice.Transpose
		return
	}
	voice.Transpose &= 0x7f

	var rawName [10]byte
	if err = binary.Read(reader, binary.LittleEndian, &rawName); err != nil {
		return
	}
	for a := 0; a < 10; a++ {
		rawName[a] &= 0x7f
	}
	voice.VoiceName = string(rawName[:])
	return
}

type sysexHeaderType int

const (
	noHeader sysexHeaderType = iota
	sysex32Voice
	sysex1Voice
)

func checkHeader(header [6]byte) (headerType sysexHeaderType) {
	multiHeaders := [][6]byte{
		{0xF0, 0x43, 0x00, 0x09, 0x20, 0x00},
		{0xF0, 0x43, 0x00, 0x09, 0x10, 0x00},
		{0xF0, 0x43, 0x00, 0x09, 0x00, 0x10},
		{0xF0, 0x43, 0x00, 0x09, 0x00, 0x20},
	}
	singleHeaders := [][6]byte{
		{0xF0, 0x43, 0x00, 0x00, 0x01, 0x1B},
	}

	for _, h := range multiHeaders {
		match := true
		for i := 0; i < 6; i++ {
			if header[i] != h[i] {
				match = false
				break // try next header
			}
		}
		if match {
			return sysex32Voice
		}
	}
	for _, h := range singleHeaders {
		match := true
		for i := 0; i < 6; i++ {
			if header[i] != h[i] {
				match = false
				break // try next header
			}
		}
		if match {
			return sysex1Voice
		}
	}
	return noHeader
}

func ReadDx7Sysex(pathname string) (sysex Dx7Sysex, err error) {
	var b []byte
	if b, err = ioutil.ReadFile(pathname); err != nil {
		return
	}
	reader := bytes.NewReader(b)

	// validate that the header of the Sysex is a "bulk DX7 sysex":
	var header [6]byte

	if err = binary.Read(reader, binary.LittleEndian, &header); err != nil {
		return
	}

	headerType := checkHeader(header)

	var packed bool
	var voiceCount int
	switch headerType {
	case noHeader:
		if _, err = reader.Seek(0, io.SeekStart); err != nil {
			err = errors.Wrapf(err, "Invalid Sysex header but failed to rewind to try to parse without header")
			return
		}
		voiceCount = 32
		packed = true
	case sysex1Voice:
		voiceCount = 1
		packed = false
	case sysex32Voice:
		voiceCount = 32
		packed = true
	}

	for i := 0; i < voiceCount; i++ {
		var v Dx7Voice
		if v, err = readDx7Voice(reader, packed); err != nil {
			err = errors.Wrapf(err, "Error reading voice[%d]", i)
			return
		}
		if v.VoiceName[0] != '\000' {
			ok := true

			// Data validation:
			if v.Algorithm > 31 {
				logger.Warnf("%s - Voice #%d \"%s\" DX Algorithm out of range: %d - must be between 0 and 31. Voice ignored",
					pathname, i, v.VoiceName, v.Algorithm)
				ok = false
			}
			if ok {
				sysex.Voices = append(sysex.Voices, v)
			}
		}
	}

	return
}

func Dx7VoiceToJSON(v Dx7Voice) (result string) {
	b, _ := json.MarshalIndent(v, "", "\t")
	result = string(b)
	return
}
