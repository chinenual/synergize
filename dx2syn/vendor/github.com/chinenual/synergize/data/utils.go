package data

func BytesToWord(hob byte, lob byte) uint16 {
	return uint16(hob)<<8 + uint16(lob)
}

func WordToBytes(word uint16) (hob byte, lob byte) {
	hob = byte(word >> 8)
	lob = byte(word)
	return
}

var appVersion string

func SetAppVersion(v string) {
	appVersion = v
}
