package osc

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// convert the field to an OSC address.  OHARM[3] -> /OHARM/3   VASENS -> /VASENS
func fieldToAddr(field string) string {
	return "/" + strings.ReplaceAll(strings.ReplaceAll(field, "[", "/"), "]", "")
}

// preemptively filter accelerometer, ping and touch messages from touchosc
var oscBlacklist = map[string]bool{
	"/accxyz": true,
	"/ping":   true,
}

// convert an OSC address to a field.  /OHARM/3 -> OHARM[3],   /VASENS -> VASENS
func addrToField(addr string) (field string, ok bool) {
	// assume first character is "/"
	ok = true
	if _, blacklisted := oscBlacklist[addr]; blacklisted {
		ok = false
		return
	}
	if strings.LastIndex(addr, "/z") == len(addr)-2 {
		// touchOSC can send extra "touch" messages for any control /theaddr/z
		ok = false
		return
	}
	slash := strings.IndexByte(addr[1:], '/')
	if slash < 0 {
		field = addr[1:]
		return
	}
	field = addr[1:slash+1] + "[" + addr[slash+2:] + "]"
	return
}

func csurfaceInit() (err error) {
	if err = oscSendString("/stringval", ""); err != nil {
		return
	}
	return
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

func changeAmpEnvRowVisibility(row int, isLast bool, onoff int) (err error) {
	for _, field := range []string{"envAmpLowTime", "envAmpUpTime"} {
		addr := fmt.Sprintf("/%s/%d/visible", field, row)

		if err = oscSendInt(addr, int32(onoff)); err != nil {
			return
		}
	}
	if isLast {
		// force last amp vals elements to be disabled
		onoff = 0
	}
	for _, field := range []string{"envAmpLowVal", "envAmpUpVal"} {
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
			if err = changeAmpEnvRowVisibility(i, i == val, onoff); err != nil {
				return
			}
		}
		return
	} else if field == "freq-env-accel-visible" {
		// special case for hiding unused controls:
		if err = oscSendInt("/accelFreqLow/visible", int32(val)); err != nil {
			return
		}
		if err = oscSendInt("/accelFreqUp/visible", int32(val)); err != nil {
			return
		}
		return
	} else if field == "amp-env-accel-visible" {
		// special case for hiding unused controls:
		if err = oscSendInt("/accelAmpLow/visible", int32(val)); err != nil {
			return
		}
		if err = oscSendInt("/accelAmpUp/visible", int32(val)); err != nil {
			return
		}
		return
	}
	addr := fieldToAddr(field)
	//reverse := addrToField(addr)
	//	logger.Infof("  field: " + field + " OSC addr: " + addr + "  reversed: " + reverse)

	if strings.HasPrefix(addr, "/FILTER") {
		var filterColorMap = []string{"gray", "red", "green"}
		// special case the tri-state filter values to also set color
		var color = filterColorMap[0]
		if val < 0 {
			val = 1
			color = filterColorMap[1]
		} else if val > 0 {
			val = 1
			color = filterColorMap[2]
		}
		//fmt.Printf("\n   %s %s %d\n\n", addr, color, val)
		if err = oscSendString(addr+"/color", color); err != nil {
			return
		}
	}
	if err = oscSendInt(addr, int32(val)); err != nil {
		return
	}
	return
}

func OscHandleFromCSurface(addr string, args ...interface{}) (err error) {
	field, ok := addrToField(addr)
	if !ok {
		// blacklisted message
		return
	}
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
