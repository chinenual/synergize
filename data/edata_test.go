package data

import (
	"fmt"
	"io/ioutil"
	"reflect"
//	"log"
	"testing"
)

func TestLocalEDATAOffsets(t *testing.T) {
//	log.Printf("edata_head_default: %v\n", edata_head_default)
	// spot check some data:
	AssertByte(t, EDATA[0], 0, "VOITAB")
	AssertByte(t, EDATA[Off_EDATA_APVIB], 32, "APVIB")

	for osc := 1; osc <= 16; osc++ {
		off := EDATALocalOscOffset(osc, 0)
		AssertByte(t, EDATA[off], 4, fmt.Sprintf("osc %d OPTCH", osc))
		off = EDATALocalOscOffset(osc, Off_EOSC_OHARM)
		AssertByte(t, EDATA[off], 1, fmt.Sprintf("osc %d OHARM", osc))
		off = EDATALocalOscOffset(osc, Off_EOSC_FDETUN)
		AssertByte(t, EDATA[off], 0, fmt.Sprintf("osc %d FDETUN", osc))
	}
}

func TestInitEDATA(t *testing.T) {
	var read_bytes []byte
	var err error

	ClearLocalEDATA()

	var path = "testfiles/VRAMRAW.bin"
	if read_bytes,err = ioutil.ReadFile(path); err != nil {
		t.Errorf("error reading %s: %v", path, err)
		return 
	}
	// only compare the portion of the EDATA actually returned from the
	// synergy (it doesnt send unused voices)
	read_EDATA   := read_bytes[Off_VRAM_EDATA:]
	EDATA_subset := EDATA[:len(read_EDATA)]
	if !reflect.DeepEqual(read_EDATA, EDATA_subset) {
		t.Errorf("initialized EDATA data doesnt match. read:\n%v\n\nEDATA:\n %v", read_EDATA, EDATA_subset)
	}
}

