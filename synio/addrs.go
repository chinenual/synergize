package synio

import (
	"github.com/chinenual/synergize/data"
)

// VRAM offsets are managed in the data package since the parsers need to know about VRAM format.
// CMOS offsets are managed here since they are strictly runtime artifacts of where things live
// during the voicing process

const (
	Off_CMOS_VTSENS = 20  // 0x14
	Off_CMOS_VTCENT = 44  // 0x2c
	Off_CMOS_VASENS = 68  // 0x44
	Off_CMOS_VACENT = 92  // 0x5c
	Off_CMOS_VVBDLY = 116 // 0x74
	Off_CMOS_VVBDEP = 140 // 0x8c
	Off_CMOS_VVBRAT = 164 // 0xa4
	Off_CMOS_VTRANS = 188 // 0xbc
)

func CmosAddr(fieldOffset int) uint16 {
	return synAddrs.CMOS + uint16(fieldOffset)
}

func VramAddr(fieldOffset int) uint16 {
	return synAddrs.VRAM + uint16(fieldOffset)
}

func VoiceHeadAddr(voiceFieldOffset int) uint16 {
	return synAddrs.VRAM + data.VRAMVoiceHeadOffset(voiceFieldOffset)
}

func VoiceOscAddr(osc /*1-based*/ int, oscFieldOffset int) uint16 {
	// osc is 1-based
	return synAddrs.VRAM + data.VRAMVoiceOscOffset(osc, oscFieldOffset)
}
