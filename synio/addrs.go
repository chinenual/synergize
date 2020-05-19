package synio

import (
	"github.com/chinenual/synergize/data"
)

func EDATAHeadAddr(fieldOffset int) uint16 {
	return synAddrs.EDATA + data.EDATALocalHeadOffset(fieldOffset)
}


func EDATAOscAddr(osc int, fieldOffset int) uint16 {
	// osc is 1-based
	return synAddrs.EDATA + data.EDATALocalOscOffset(osc, fieldOffset)	
}
