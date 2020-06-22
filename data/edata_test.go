package data

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"

	//	"log"
	"testing"

	"github.com/orcaman/writerseeker"
)

func TestLocalEDATAOffsets(t *testing.T) {
	//	log.Printf("edata_head_default: %v\n", edata_head_default)
	// spot check some data:
	AssertByte(t, VRAM_EDATA[Off_VRAM_EDATA+0], 0, "VOITAB")
	AssertByte(t, VRAM_EDATA[Off_VRAM_EDATA+Off_EDATA_APVIB], 32, "APVIB")

	for osc := 1; osc <= 16; osc++ {
		off := VRAMVoiceOscOffset(osc, 0)
		AssertByte(t, VRAM_EDATA[off], 4, fmt.Sprintf("osc %d OPTCH", osc))
		off = VRAMVoiceOscOffset(osc, Off_EOSC_OHARM)
		AssertByte(t, VRAM_EDATA[off], 1, fmt.Sprintf("osc %d OHARM", osc))
		off = VRAMVoiceOscOffset(osc, Off_EOSC_FDETUN)
		AssertByte(t, VRAM_EDATA[off], 0, fmt.Sprintf("osc %d FDETUN", osc))
	}
}

func TestInitEDATA(t *testing.T) {
	var read_bytes []byte
	var err error

	ClearLocalEDATA()

	var path = "testfiles/VRAMRAW.bin"
	if read_bytes, err = ioutil.ReadFile(path); err != nil {
		t.Errorf("error reading %s: %v", path, err)
		return
	}

	for i := Off_VRAM_VCHK; i < Off_VRAM_VCHK+5; i++ {
		if VRAM_EDATA[i] != 0xaa {
			t.Errorf("Bad check byte at %d: %2x", i, VRAM_EDATA[i])
		}
	}

	// only compare the portion of the EDATA actually returned from the
	// synergy (it doesnt send unused voices)
	read_EDATA := read_bytes[Off_VRAM_EDATA:]
	EDATA_subset := VRAM_EDATA[Off_VRAM_EDATA : Off_VRAM_EDATA+len(read_EDATA)]
	if !reflect.DeepEqual(read_EDATA, EDATA_subset) {
		t.Errorf("initialized EDATA data doesnt match. read:\n%v\n\nEDATA:\n %v", read_EDATA, EDATA_subset)
	}
}

//func TestDiffG7S(t *testing.T) {
//	vce1,_:=ReadVceFile("testfiles/G7S.VCE")
//	vce2,_:=ReadVceFile("testfiles/VRAMG7S.VCE")
//	_=diffVCE(vce1,vce2)
//}

func TestReadVceFromVRAM(t *testing.T) {
	var read_bytes []byte
	var err error

	var vram_path = "testfiles/VRAMG7S.bin"
	// don't get with real G7S.VCE - something about loading it into VRAM
	// and then back alters it -- changes VIBRAT. So compare against the
	// VCE that SYNCS created by the same VRAM dump
	var vce_path = "testfiles/VRAMG7S.VCE"
	if read_bytes, err = ioutil.ReadFile(vram_path); err != nil {
		t.Errorf("error reading %s: %v", vram_path, err)
		return
	}

	var vce_from_vram VCE
	var vce_from_file VCE

	if vce_from_vram, err = ReadVceFromVRAM(read_bytes); err != nil {
		t.Errorf("error reading vce: %v", err)
		t.Errorf("  parsed so far: %s\n", VceToJson(vce_from_vram))
		return
	}

	if vce_from_file, err = ReadVceFile(vce_path); err != nil {
		t.Errorf("error reading %s: %v", vce_path, err)
		return
	}
	// don't diff the raw vce -- it hasnt been re-compressed = write it
	// via WriteVce to get something comparable
	var writebuf = writerseeker.WriterSeeker{}
	if err = WriteVce(&writebuf, vce_from_vram, "VRAMG7S", false); err != nil {
		t.Errorf("error normalizing VRAM %v", err)
	}
	// now read in the compressed vce:
	write_bytes, _ := ioutil.ReadAll(writebuf.Reader())
	var readbuf = bytes.NewReader(write_bytes)

	var vce_from_vram2 VCE
	if vce_from_vram2, err = ReadVce(readbuf, false); err != nil {
		t.Errorf("error parsing generated stream: %v", err)
		t.Errorf("  parsed so far: %s\n", VceToJson(vce_from_vram2))
		return
	}

	if !diffVCE(vce_from_vram2, vce_from_file) {
		t.Errorf("VCE do not match %s %s", vram_path, vce_path)
	}
}
