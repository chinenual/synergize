package synio

import (
	"github.com/chinenual/synergize/data"
	"github.com/pkg/errors"
)

func VoicingMode() (err error) {
	if err = getSynergyAddrs(); err != nil {
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

func SetVoiceHeadDataArray(offset int, value []byte, purpose string, reloadGen bool) (err error) {
	addr := EDATAHeadAddr(offset)
	if err = BlockLoad(addr, value, purpose); err != nil {
		return
	}
	if reloadGen {
		if err = ReloadNoteGenerators(); err != nil {
			return
		}
	}
	return
}

func SetVoiceHeadDataByte(offset int, value byte, purpose string, reloadGen bool) (err error) {
	addr := EDATAHeadAddr(offset)
	if err = LoadByte(addr, value, purpose); err != nil {
		return
	}
	if reloadGen {
		if err = ReloadNoteGenerators(); err != nil {
			return
		}
	}
	return
}

// osc is 1-based
func SetVoiceOscDataByte(osc int, offset int, value byte, purpose string, reloadGen bool) (err error) {
	addr := EDATAOscAddr(osc, offset)
	if err = LoadByte(addr, value, purpose); err != nil {
		return
	}
	if reloadGen {
		if err = ReloadNoteGenerators(); err != nil {
			return
		}
	}
	return
}

// osc is 1-based
func GetVoiceOscDataByte(osc int, offset int, purpose string) (value byte, err error) {
	addr := EDATAOscAddr(osc, offset)
	if value, err = DumpByte(addr, purpose); err != nil {
		return
	}
	return
}

func SetVoiceVEQValue(index, value byte) (err error) {
	if err = SetVoiceHeadDataByte(data.Off_EDATA_VEQ+int(index), value, "set VEQ", false); err != nil {
		return
	}
	return
}

func SetVoiceVEQArray(value []byte) (err error) {
	if err = SetVoiceHeadDataArray(data.Off_EDATA_VEQ, value, "set VEQ", false); err != nil {
		return
	}
	return
}

func SetVoiceKPROPValue(index, value byte) (err error) {
	if err = SetVoiceHeadDataByte(data.Off_EDATA_KPROP+int(index), value, "set Kprop", false); err != nil {
		return
	}
	return
}

func SetVoiceKPROPArray(value []byte) (err error) {
	if err = SetVoiceHeadDataArray(data.Off_EDATA_KPROP, value, "set Kprop", false); err != nil {
		return
	}
	return
}

func SetVoiceAPVIB(value byte) (err error) {
	if err = SetVoiceHeadDataByte(data.Off_EDATA_APVIB, value, "set APVIB", true); err != nil {
		return
	}
	return
}

func SetVoiceOscOHARM(osc int, value int8) (err error) {
	if err = SetVoiceOscDataByte(osc, data.Off_EOSC_OHARM, byte(value), "set OHARM", true); err != nil {
		return
	}
	return
}

func SetVoiceOscFDETUN(osc int, value int8) (err error) {
	if err = SetVoiceOscDataByte(osc, data.Off_EOSC_FDETUN, byte(value), "set FDETUN", true); err != nil {
		return
	}
	return
}

// emulate the SYNHCS GEDPTR subroutine: get OSC specific offset into the EDATA array
func gedptr(osc int) uint16 {
	return uint16(2*osc) + synAddrs.EDATA + 1
}

func EncodePatchControl(outputDSR byte, inhibitAddr byte,
	adderInputDSR byte, inhibitF0 byte, f0InputDSR byte) (control byte) {

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

func GetOscPATCHControl(osc int) (value byte, err error) {
	if value, err = GetVoiceOscDataByte(osc, data.Off_EOSC_OPTCH, "get OPATCH"); err != nil {
		return
	}
	return
}

func SetOscPATCHControl(osc int, value byte) (err error) {
	if err = SetVoiceOscDataByte(osc, data.Off_EOSC_OPTCH, value, "set OPATCH", true); err != nil {
		return
	}
	return
}

func SetOscFREQControl(osc int, value byte) (err error) {
	err = errors.New("not yet implemented")
	return
}

func GetOscWAVEControl(osc int) (value byte, err error) {
	if value, err = GetVoiceOscDataByte(osc, data.Off_EOSC_FreqPoints+3, "get fenv[3]"); err != nil {
		return
	}
	return
}

func SetOscWAVEControl(osc int, value byte) (err error) {
	err = errors.New("not yet implemented")
	// wave is stored in 3 bits the 4th entry in the freq envelope. eesh.
	// 0x633e, 0x01 == sine, 0x00 == triangle
	if err = SetVoiceOscDataByte(osc, data.Off_EOSC_FreqPoints+3, value, "set fenv[3]", true); err != nil {
		return
	}
	return
}

func SetOscFILTER(osc int, filter int8) (err error) {
	err = errors.New("not yet implemented")
	return
}

func SetOscWAVE(osc int, triangle bool) (err error) {
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

/*****
func SetVoiceVTCENT(val byte) (err error) {
	if err = SetVoiceOscDataByte(data.Off_EDATA_VTCENT, byte(value), "set VTCENT", false); err != nil {
		return
	}
	return
}

func SetVoiceVTSENS(val byte) (err error) {
	if err = SetVoiceOscDataByte(data.Off_EDATA_VTSENS, byte(value), "set VTSENS", false); err != nil {
		return
	}
	return
}

func SetVoiceVACENT(val byte) (err error) {
	if err = SetVoiceOscDataByte(data.Off_EDATA_VACENT, byte(value), "set VACENT", false); err != nil {
		return
	}
	return
}

func SetVoiceVASENS(val byte) (err error) {
	if err = SetVoiceOscDataByte(data.Off_EDATA_VASENS, byte(value), "set VASENS", false); err != nil {
		return
	}
	return
}

func SetVoiceVIBRAT(val byte) (err error) {
	if err = SetVoiceOscDataByte(data.Off_EDATA_VIBRAT, byte(value), "set VIBRAT", false); err != nil {
		return
	}
	return
}

func SetVoiceVIBDEL(val byte) (err error) {
	if err = SetVoiceOscDataByte(data.Off_EDATA_VIBDEL, byte(value), "set VIBDEL", false); err != nil {
		return
	}
	return
}

func SetVoiceVIBDEP(val byte) (err error) {
	if err = SetVoiceOscDataByte(data.Off_EDATA_VIBDEP, byte(value), "set VIBDEP", false); err != nil {
		return
	}
	return
}

******/
