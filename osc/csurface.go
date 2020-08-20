package osc

import (
	"fmt"
	"github.com/pkg/errors"
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

func changeOscRowVisibility(row int, onoff int) (err error) {
	for _, field := range []string{"MUTE", "SOLO", "OHARM", "FDETUN", "wkWAVE", "wkKEYPROP", "FILTER"} {
		addr := fmt.Sprintf("/%s/%d/visible", field, row)

		if err = oscSendInt(addr, int32(onoff)); err != nil {
			return
		}
	}
	return
}

func changeAmpEnvRowVisibility(row int, onoff int) (err error) {
	for _, field := range []string{"envAmpLowVal", "envAmpUpVal", "envAmpLowTime", "envAmpUpTime"} {
		addr := fmt.Sprintf("/%s/%d/visible", field, row)

		if err = oscSendInt(addr, int32(onoff)); err != nil {
			return
		}
	}
	return
}

func changeFreqEnvRowVisibility(row int, onoff int) (err error) {
	for _, field := range []string{"envFreqLowVal", "envFreqUpVal", "envFreqLowTime", "envFreqUpTime"} {
		addr := fmt.Sprintf("/%s/%d/visible", field, row)

		if err = oscSendInt(addr, int32(onoff)); err != nil {
			return
		}
	}
	return
}

// field looks like "OHARM[3]" or VASENS  -- i.e an identifier followed by an optional index
func OscSendToCSurface(field string, val int) (err error) {
	if field == "num-osc" {
		// special case for hiding unused controls:
		for i := 1; i <= 16; i++ {
			var onoff = 0
			if i <= val {
				onoff = 1
			}
			if err = changeOscRowVisibility(i, onoff); err != nil {
				return
			}
		}
		return
	} else if field == "num-freq-env-points" {
		// special case for hiding unused controls:
		for i := 1; i <= 16; i++ {
			var onoff = 0
			if i <= val {
				onoff = 1
			}
			if err = changeFreqEnvRowVisibility(i, onoff); err != nil {
				return
			}
		}
		return
	} else if field == "num-amp-env-points" {
		// special case for hiding unused controls:
		for i := 1; i <= 16; i++ {
			var onoff = 0
			if i <= val {
				onoff = 1
			}
			if err = changeAmpEnvRowVisibility(i, onoff); err != nil {
				return
			}
		}
		return
	}
	addr := fieldToAddr(field)
	//reverse := addrToField(addr)
	//	log.Printf("  field: " + field + " OSC addr: " + addr + "  reversed: " + reverse)

	if strings.HasPrefix(addr, "/FILTER") {
		var filterColorMap = []string{"gray", "red", "green"}
		// special case the tri-state filter values to also set color
		if err = oscSendString(addr+"/color", filterColorMap[val]); err != nil {
			return
		}
	}
	if err = oscSendInt(addr, int32(val)); err != nil {
		return
	}
	return
}

func OscHandleFromCSurface(addr string, args ...interface{}) (err error) {
	field := addrToField(addr)
	var arg int
	if len(args) >= 1 {
		switch args[0].(type) {
		case int32:
			arg = int(args[0].(int32))
		case float32:
			arg = int(args[0].(float32))
		default:
			err = errors.Errorf("Unhandled OSC arg type: %v\n", args[0])
			return
		}
	}
	if err = sendToUI(field, arg); err != nil {
		return
	}
	return
}