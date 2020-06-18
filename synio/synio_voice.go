package synio

import (
	"bytes"
	"strconv"

	"github.com/chinenual/synergize/data"
)

type offsetMapEle struct {
	Offset    int
	ReloadGen bool
}

var oscOffsetMap map[string]offsetMapEle
var voiceOffsetMap map[string]offsetMapEle
var cmosOffsetMap map[string]offsetMapEle

func initMaps() {
	if oscOffsetMap != nil {
		return
	}
	oscOffsetMap = make(map[string]offsetMapEle)
	voiceOffsetMap = make(map[string]offsetMapEle)
	cmosOffsetMap = make(map[string]offsetMapEle)

	voiceOffsetMap["VOITAB"] = offsetMapEle{data.Off_EDATA_VOITAB, false}
	voiceOffsetMap["VTRANS"] = offsetMapEle{data.Off_EDATA_VTRANS, false}
	voiceOffsetMap["APVIB"] = offsetMapEle{data.Off_EDATA_APVIB, false}
	voiceOffsetMap["KPROP"] = offsetMapEle{data.Off_EDATA_KPROP, false}
	voiceOffsetMap["VEQ"] = offsetMapEle{data.Off_EDATA_VEQ, false}
	voiceOffsetMap["FILTER"] = offsetMapEle{data.Off_EDATA_FILTER_arr, false}

	cmosOffsetMap["VTCENT"] = offsetMapEle{Off_CMOS_VTCENT, false}
	cmosOffsetMap["VTSENS"] = offsetMapEle{Off_CMOS_VTSENS, false}
	cmosOffsetMap["VACENT"] = offsetMapEle{Off_CMOS_VACENT, false}
	cmosOffsetMap["VASENS"] = offsetMapEle{Off_CMOS_VASENS, false}
	cmosOffsetMap["VIBRAT"] = offsetMapEle{Off_CMOS_VVBRAT, false}
	cmosOffsetMap["VIBDEL"] = offsetMapEle{Off_CMOS_VVBDLY, false}
	cmosOffsetMap["VIBDEP"] = offsetMapEle{Off_CMOS_VVBDEP, false}

	oscOffsetMap["OPTCH"] = offsetMapEle{data.Off_EOSC_OPTCH, false} // does require reload, but we do it after a batch
	oscOffsetMap["OHARM"] = offsetMapEle{data.Off_EOSC_OHARM, true}
	oscOffsetMap["FDETUN"] = offsetMapEle{data.Off_EOSC_FDETUN, true}

	oscOffsetMap["FreqENVTYPE"] = offsetMapEle{data.Off_EOSC_FreqENVTYPE, true}
	oscOffsetMap["FreqNPOINTS"] = offsetMapEle{data.Off_EOSC_FreqNPOINTS, true}
	oscOffsetMap["FreqSUSTAINPT"] = offsetMapEle{data.Off_EOSC_FreqSUSTAINPT, true}
	oscOffsetMap["FreqLOOPPT"] = offsetMapEle{data.Off_EOSC_FreqLOOPPT, true}
	oscOffsetMap["FreqPoints"] = offsetMapEle{data.Off_EOSC_FreqPoints, true}
	oscOffsetMap["FreqPoints_WAVE_KEYPROP"] = offsetMapEle{data.Off_EOSC_FreqPoints + 3, true}

	oscOffsetMap["AmpENVTYPE"] = offsetMapEle{data.Off_EOSC_AmpENVTYPE, true}
	oscOffsetMap["AmpNPOINTS"] = offsetMapEle{data.Off_EOSC_AmpNPOINTS, true}
	oscOffsetMap["AmpSUSTAINPT"] = offsetMapEle{data.Off_EOSC_AmpSUSTAINPT, true}
	oscOffsetMap["AmpLOOPPT"] = offsetMapEle{data.Off_EOSC_AmpLOOPPT, true}
	oscOffsetMap["AmpPoints"] = offsetMapEle{data.Off_EOSC_AmpPoints, true}
}

func EnableVoicingMode() (vce data.VCE, err error) {
	initMaps()

	if err = getSynergyAddrs(); err != nil {
		return
	}
	if err = InitVRAM(); err != nil {
		return
	}
	data.ClearLocalEDATA()
	if err = LoadCRT(data.VRAM_EDATA[:]); err != nil {
		return
	}
	rdr := bytes.NewReader(data.VRAM_EDATA[data.Off_VRAM_EDATA:])
	if vce, err = data.ReadVce(rdr, false); err != nil {
		return
	}

	if err = rawSetOscSolo([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}); err != nil {
		return
	}
	return
}

func DisableVoicingMode() (err error) {
	if err = rawSetOscSolo([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}); err != nil {
		return
	}
	return
}

func rawSetOscSolo(oscStatus [16]byte) (err error) {
	if err = BlockLoad(synAddrs.SOLOSC, oscStatus[:], "set SOLOSC"); err != nil {
		return
	}
	if err = ReloadNoteGenerators(); err != nil {
		return
	}
	return
}

func SetOscSolo(mute, solo []bool) (oscStatus [16]bool, err error) {
	// 0 = on, 1 = off
	var state = [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	oscStatus = [16]bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true}

	// solo takes precedence. If any soloed, then ignore mutes
	for _, soloed := range solo {
		if soloed {
			for i, v := range solo {
				state[i] = 1
				oscStatus[i] = false
				if v {
					state[i] = 0
					oscStatus[i] = true
				}
			}
			if err = rawSetOscSolo(state); err != nil {
				return
			}
			return
		}
	}
	// if no solo, then just mute the ones selected (if any):
	for i, muted := range mute {
		state[i] = 0
		oscStatus[i] = true
		if muted {
			state[i] = 1
			oscStatus[i] = false
		}
	}
	if err = rawSetOscSolo(state); err != nil {
		return
	}
	return
}

func SetFilterEle(uiFilterIndex /*0 for Af, one-based osc# for Bf */ int, index /* one-based */ int, value int) (err error) {
	// ASSUMES we're only editing voice #1.
	// AFilter is always at 0 in the FILTAB;
	// Bfilters start at 2, so osc #1's filter is at zero-based index 1 of the FILTAB
	// Bfilter value is the one-based osc#

	addr := VramAddr(data.Off_VRAM_FILTAB) + uint16((uiFilterIndex*data.VRAM_FILTR_length)+(index-1))
	if err = LoadByte(addr, byte(value), "set BFILTER["+strconv.Itoa(index)+"]"); err != nil {
		return
	}
	if err = RecalcFilters(); err != nil {
		return
	}
	return
}

func SetOscFILTER(osc /*1-based*/ int, value int) (err error) {
	addr := VoiceHeadAddr(data.Off_EDATA_FILTER_arr) + uint16(osc-1)
	if err = LoadByte(addr, byte(value), "set FILTER["+strconv.Itoa(osc)+"]"); err != nil {
		return
	}
	if err = ReloadNoteGenerators(); err != nil {
		return
	}
	return
}

func SetPatchType(index int) (patchBytes [16]byte, err error) {
	// write all 16 oscillators whether they're in use or not
	for osc := 1; osc <= 16; osc++ {
		if err = SetVoiceOscDataByte(osc, "OPTCH", data.PatchTypePerOscTable[index][osc-1]); err != nil {
			return
		}
	}
	if err = ReloadNoteGenerators(); err != nil {
		return
	}
	patchBytes = data.PatchTypePerOscTable[index]
	return
}

func SetNumOscillators(newNumOsc int, patchType int) (patchBytes [16]byte, err error) {
	// assumes that we dont need to reininitialize freq and amp envelopes (they were initialized
	// when we started voicing mode and if we are reusing one partially edited, we get those edits back)
	if err = SetVoiceHeadDataByte("VOITAB", byte(newNumOsc-1)); err != nil {
		return
	}
	if patchBytes, err = SetPatchType(patchType); err != nil {
		return
	}
	return
}

func LoadVceVoicingMode(vce data.VCE) (err error) {
	if err = data.LoadVceIntoEDATA(vce); err != nil {
		return
	}
	if err = LoadCRT(data.VRAM_EDATA[:]); err != nil {
		return
	}
	return
}

func RecalcEq() (err error) {
	if err = command(OP_EXECUTE, "OP_EXECUTE"); err != nil {
		return
	}
	if err = writeU16(synAddrs.exec_REAEQ, "REAEQ addr"); err != nil {
		return
	}
	if err = writeU16(0, "REAEQ args"); err != nil {
		return
	}
	return
}

func RecalcFilters() (err error) {
	if err = command(OP_EXECUTE, "OP_EXECUTE"); err != nil {
		return
	}
	if err = writeU16(synAddrs.exec_REFIL, "REFIL addr"); err != nil {
		return
	}
	if err = writeU16(0, "REFIL args"); err != nil {
		return
	}
	return
}

func ReloadNoteGenerators() (err error) {
	if err = command(OP_EXECUTE, "OP_EXECUTE"); err != nil {
		return
	}
	if err = writeU16(synAddrs.exec_LDGENR, "LDGENR addr"); err != nil {
		return
	}
	if err = writeU16(0, "LDGENR args"); err != nil {
		return
	}
	return
}

// Sets the value in the Synergy address space and then reloads the note
// generators

func SetVoiceHeadDataArray(fieldName string, value []byte) (err error) {
	// some things that are stored in the head are actually stored in a different location in CMOS
	// at runtime.  Deal with that here:
	offsetMap := cmosOffsetMap
	var addr uint16
	if _, ok := voiceOffsetMap[fieldName]; ok {
		offsetMap = voiceOffsetMap
		addr = VoiceHeadAddr(offsetMap[fieldName].Offset)
	} else {
		addr = CmosAddr(offsetMap[fieldName].Offset)
	}

	if err = BlockLoad(addr, value, "set array "+fieldName); err != nil {
		return
	}
	if offsetMap[fieldName].ReloadGen {
		if err = ReloadNoteGenerators(); err != nil {
			return
		}
	}
	return
}

func SetVoiceHeadDataByte(fieldName string, value byte) (err error) {
	offsetMap := cmosOffsetMap
	var addr uint16
	if _, ok := voiceOffsetMap[fieldName]; ok {
		offsetMap = voiceOffsetMap
		addr = VoiceHeadAddr(offsetMap[fieldName].Offset)
	} else {
		addr = CmosAddr(offsetMap[fieldName].Offset)
	}
	if err = LoadByte(addr, value, "set "+fieldName); err != nil {
		return
	}
	if offsetMap[fieldName].ReloadGen {
		if err = ReloadNoteGenerators(); err != nil {
			return
		}
	}
	return
}

// osc is 1-based
func SetVoiceOscDataByte(osc /*1-based*/ int, fieldName string, value byte) (err error) {
	addr := VoiceOscAddr(osc, oscOffsetMap[fieldName].Offset)
	if err = LoadByte(addr, value, "set "+fieldName+"["+strconv.Itoa(osc)+"]"); err != nil {
		return
	}
	if oscOffsetMap[fieldName].ReloadGen {
		if err = ReloadNoteGenerators(); err != nil {
			return
		}
	}
	return
}

// osc is 1-based
func GetVoiceOscDataByte(osc /*1-based*/ int, fieldName string) (value byte, err error) {
	addr := VoiceOscAddr(osc, oscOffsetMap[fieldName].Offset)
	if value, err = DumpByte(addr, "get "+fieldName+"["+strconv.Itoa(osc)+"]"); err != nil {
		return
	}
	return
}

func SetVoiceVEQEle(index /* 1-based */ int, value int) (err error) {
	addr := VoiceHeadAddr(data.Off_EDATA_VEQ) + uint16(index-1)
	if err = LoadByte(addr, byte(value), "set VEQ["+strconv.Itoa(index)+"]"); err != nil {
		return
	}
	if err = RecalcEq(); err != nil {
		return
	}
	return
}

func SetVoiceKPROPEle(index /* 1-based */ int, value int) (err error) {
	addr := VoiceHeadAddr(data.Off_EDATA_KPROP) + uint16(index-1)
	if err = LoadByte(addr, byte(value), "set KPROP["+strconv.Itoa(index)+"]"); err != nil {
		return
	}
	return
}

func EncodePatchControl(outputDSR byte, inhibitAddr byte,
	adderInputDSR byte, inhibitF0 byte, f0InputDSR byte) (control byte) {

	//Patch Control Byte (k=0):
	//
	// 7  6    5     4  3     2      1  0
	//-------------------------------------
	//| OUT | ENAB |  DSR  | ENAB |  DSR  |
	//-------------------------------------
	//   ^     ^       ^      ^       ^
	//   ^     ^       ^      ^       +++++ FO Input DSR
	//   ^     ^       ^      ^
	//   ^     ^       ^      +++++++++++++ 1 = Inhibit FO
	//   ^     ^       ^
	//   ^     ^       ++++++++++++++++++++ Adder Input DSR
	//   ^     ^
	//   ^     ++++++++++++++++++++++++++++ 1 = Inhibit Adder
	//   ^
	//   ++++++++++++++++++++++++++++++++++ Output DSR

	control = 0
	control = control | (0x3 & f0InputDSR)
	control = control | ((0x1 & inhibitF0) << 2)
	control = control | ((0x3 & adderInputDSR) << 3)
	control = control | ((0x1 & inhibitAddr) << 5)
	control = control | ((0x3 & outputDSR) << 6)
	return
}

func DecodePatchControl(control byte) (outputDSR byte, inhibitAddr byte,
	adderInputDSR byte, inhibitF0 byte, f0InputDSR byte) {
	f0InputDSR = control & 0x3
	inhibitF0 = (control >> 2) & 0x1
	adderInputDSR = (control >> 3) & 0x3
	inhibitAddr = (control >> 5) & 0x1
	outputDSR = (control >> 6)
	return
}

func GetOscWAVEControl(osc int) (value byte, err error) {
	if value, err = GetVoiceOscDataByte(osc, "FreqPoints_WAVE_KEYPROP"); err != nil {
		return
	}
	return
}

func SetOscWAVEControl(osc int, value byte) (err error) {
	// based on snooping the serial line, wave is stored in 3 bits the
	// 4th entry in the freq envelope. eesh.
	// 0x633e, 0x01 == sine, 0x00 == triangle
	//
	// this does not match the osc bit descriptor in the manual:
	//     7         6      5 4 3    2 1 0
	//---------------------------------------
	//| AMP INT | FRQ INT | OCTAVE |  WAVE  |
	//---------------------------------------
	//     ^         ^        ^        ++++++ 000 = Sine
	//     ^         ^        ^		001 = Triangle
	//     ^         ^        ^
	//     ^         ^        +++++++++++++++	000 = No reduction
	//     ^         ^			001 = Freq./2
	//     ^         ^			010 = Freq./4
	//     ^         ^			011 = Freq./8
	//     ^         ^			100 = Freq./16
	//     ^         ^			101 = Freq./32
	//     ^         ^			110 = Freq./64
	//     ^         ^			111 = Shut Down
	//     ^         ^
	//     ^         ++++++++++++++++++++++++	Freq. Ramp Interrupt
	//     ^					1 = Enabled
	//     ^					0 = Disabled
	//     ^
	//     ++++++++++++++++++++++++++++++++++	Amp. Ramp Interrupt
	//					1 = Enabled
	//					0 = Disabled

	if err = SetVoiceOscDataByte(osc, "FreqPoints_WAVE_KEYPROP", value); err != nil {
		return
	}
	return
}

func SetOscWAVE(osc /*1-based*/ int, triangle bool) (err error) {
	// FIXME: this requires a fetch before the set -- could avoid this if we keep a copy of each value
	// on our side - like SYNHCS does.  For now, I'm trying to avoid that bookkeeping - as long as performance
	// is ok...
	var value byte
	if value, err = GetOscWAVEControl(osc); err != nil {
		return
	}
	if triangle {
		value = 0x1 | value
	} else {
		value = 0x1 &^ value
	}
	if err = SetOscWAVEControl(osc, value); err != nil {
		return
	}
	return
}

func SetOscKEYPROP(osc int, usesKeypro bool) (err error) {
	// keyprop is in the same control byte as the waveform (4th entry in the osc freq env)
	// 0x10 == 'k', 0x00 == ' '
	// NOTE: This seems to contradict the documentation -- this byte must not be exactly like the "osc control byte"
	// there, bit 0x10 is part of the OCTAVE setting not keyprop.
	var value byte
	if value, err = GetOscWAVEControl(osc); err != nil {
		return
	}
	if usesKeypro {
		value = 0x10 | value
	} else {
		value = 0x10 &^ value
	}
	if err = SetOscWAVEControl(osc, value); err != nil {
		return
	}
	return
}

// Each point in the Freq or Amp table has 4 values:  valLow, valUp, timeLow, TimeUp

func SetEnvFreqLowVal(osc /* 1-based */ int, pointIndex /* 1-based */ int, value byte) (err error) {
	var addr = VoiceOscAddr(osc, data.Off_EOSC_FreqPoints) + uint16(4*(pointIndex-1)+0)
	if err = LoadByte(addr, byte(value), "set EnvFreqLowVal["+strconv.Itoa(osc)+"]["+strconv.Itoa(pointIndex)+"]"); err != nil {
		return
	}
	return
}

func SetEnvFreqUpVal(osc /* 1-based */ int, pointIndex /* 1-based */ int, value byte) (err error) {
	var addr = VoiceOscAddr(osc, data.Off_EOSC_FreqPoints) + uint16(4*(pointIndex-1)+1)
	if err = LoadByte(addr, byte(value), "set EnvFreqUpVal["+strconv.Itoa(osc)+"]["+strconv.Itoa(pointIndex)+"]"); err != nil {
		return
	}
	return
}

func SetEnvFreqLowTime(osc /* 1-based */ int, pointIndex /* 1-based */ int, value byte) (err error) {
	var addr = VoiceOscAddr(osc, data.Off_EOSC_FreqPoints) + uint16(4*(pointIndex-1)+2)
	if err = LoadByte(addr, byte(value), "set EnvFreqLowTime["+strconv.Itoa(osc)+"]["+strconv.Itoa(pointIndex)+"]"); err != nil {
		return
	}
	return
}

func SetEnvFreqUpTime(osc /* 1-based */ int, pointIndex /* 1-based */ int, value byte) (err error) {
	var addr = VoiceOscAddr(osc, data.Off_EOSC_FreqPoints) + uint16(4*(pointIndex-1)+3)
	if err = LoadByte(addr, byte(value), "set EnvFreqUpTime["+strconv.Itoa(osc)+"]["+strconv.Itoa(pointIndex)+"]"); err != nil {
		return
	}
	return
}

func SetEnvAmpLowVal(osc /* 1-based */ int, pointIndex /* 1-based */ int, value byte) (err error) {
	var addr = VoiceOscAddr(osc, data.Off_EOSC_AmpPoints) + uint16(4*(pointIndex-1)+0)
	if err = LoadByte(addr, byte(value), "set EnvAmpLowVal["+strconv.Itoa(osc)+"]["+strconv.Itoa(pointIndex)+"]"); err != nil {
		return
	}
	return
}

func SetEnvAmpUpVal(osc /* 1-based */ int, pointIndex /* 1-based */ int, value byte) (err error) {
	var addr = VoiceOscAddr(osc, data.Off_EOSC_AmpPoints) + uint16(4*(pointIndex-1)+1)
	if err = LoadByte(addr, byte(value), "set EnvAmpUpVal["+strconv.Itoa(osc)+"]["+strconv.Itoa(pointIndex)+"]"); err != nil {
		return
	}
	return
}

func SetEnvAmpLowTime(osc /* 1-based */ int, pointIndex /* 1-based */ int, value byte) (err error) {
	var addr = VoiceOscAddr(osc, data.Off_EOSC_AmpPoints) + uint16(4*(pointIndex-1)+2)
	if err = LoadByte(addr, byte(value), "set EnvAmpLowTime["+strconv.Itoa(osc)+"]["+strconv.Itoa(pointIndex)+"]"); err != nil {
		return
	}
	return
}

func SetEnvAmpUpTime(osc /* 1-based */ int, pointIndex /* 1-based */ int, value byte) (err error) {
	var addr = VoiceOscAddr(osc, data.Off_EOSC_AmpPoints) + uint16(4*(pointIndex-1)+3)
	if err = LoadByte(addr, byte(value), "set EnvAmpUpTime["+strconv.Itoa(osc)+"]["+strconv.Itoa(pointIndex)+"]"); err != nil {
		return
	}
	return
}

func SetOscEnvLengths(osc /* 1-based */ int, freqLength int, ampLength int) (err error) {
	var addr = VoiceOscAddr(osc, data.Off_EOSC_FreqNPOINTS)
	if err = LoadByte(addr, byte(freqLength), "set EnvFreq NPOINTS["+strconv.Itoa(osc)+"]"); err != nil {
		return
	}
	addr = VoiceOscAddr(osc, data.Off_EOSC_AmpNPOINTS)
	if err = LoadByte(addr, byte(ampLength), "set EnvAmp NPOINTS["+strconv.Itoa(osc)+"]"); err != nil {
		return
	}
	return

}
