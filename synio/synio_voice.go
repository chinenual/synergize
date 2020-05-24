package synio
import (
	"github.com/chinenual/synergize/data"
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

func SetVoiceVeqValue(index, value byte) (err error) {
	if err = SetVoiceHeadDataByte(data.Off_EDATA_VEQ+int(index), value, "set VEQ", false); err != nil {
		return
	}
	return
}

func SetVoiceVeqArray(value []byte) (err error) {
	if err = SetVoiceHeadDataArray(data.Off_EDATA_VEQ, value, "set VEQ", false); err != nil {
		return
	}
	return
}

func SetVoiceKpropValue(index, value byte) (err error) {
	if err = SetVoiceHeadDataByte(data.Off_EDATA_KPROP+int(index), value, "set Kprop", false); err != nil {
		return
	}
	return
}

func SetVoiceKpropArray(value []byte) (err error) {
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

func SetOscHarmonic(osc int, value byte) (err error) {
	addr := gedptr(osc)
	if err = LoadByte(addr, value, "Osc harmonic value"); err != nil {
		return
	}
	if err = ReloadNoteGenerators(); err != nil {
		return
	}
	return
}
