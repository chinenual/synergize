package osc

import (
	"log"
	"strings"
)

// convert the field to an OSC address.  OHARM[3] -> /OHARM/3   VASENS -> /VASENS
func fieldToAddr(field string) string {
	return "/" + strings.ReplaceAll(strings.ReplaceAll(field, "[", "/"), "]", "")
}

// convert an OSC address to a field.  /OHARM/3 -> OHARM[3],   /VASENS -> VASENS
func addrToField(addr string) string {
	// assume first character is "/"
	slash := strings.IndexByte(addr[1:], '/')
	if slash < 0 {
		return addr[1:]
	}
	return addr[1:slash+1] + "[" + addr[slash+2:] + "]"
}

// field looks like "OHARM[3]" or VASENS  -- i.e an identifier followed by an optional index
func OscSendToCSurface(field string, val int) (err error) {
	addr := fieldToAddr(field)
	reverse := addrToField(addr)
	log.Printf("  field: " + field + " OSC addr: " + addr + "  reversed: " + reverse)

	if err = oscSendInt(addr, int32(val)); err != nil {
		return
	}
	return
}
