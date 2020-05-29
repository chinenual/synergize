package synio

import (
	"github.com/chinenual/synergize/data"
)

func VoiceHeadAddr(voiceFieldOffset int) uint16 {
	return synAddrs.VRAM + data.VRAMVoiceHeadOffset(voiceFieldOffset)
}

func VoiceOscAddr(osc /*1-based*/ int, oscFieldOffset int) uint16 {
	// osc is 1-based
	return synAddrs.VRAM + data.VRAMVoiceOscOffset(osc, oscFieldOffset)
}
