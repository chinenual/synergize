package synio

import (
	"flag"
	"fmt"
	"io/ioutil"
	"testing"
	"os"
	"reflect"
)

var (
	synio = flag.Bool("synio", false, "run integration tests that talk to the Synergy")
	port = flag.String("port", "", "the serial device")
	baud = flag.Uint("baud", 9600, "the serial baud rate")
)

func connectToSynergy() (err error) {
	return Init(*port, *baud, true, false)
}

func TestGetFirmwareId(t *testing.T) {
	id, err := GetID();
	if err != nil {
		t.Fatalf("Error when getting id: %v", err)
	}
	if id[0] != 3 || id[1] != 22 {
		t.Errorf("Expected 3.22, got %v", id)
	}
}

func assertUint16(t *testing.T, b uint16, expected uint16, context string) {
	if b != expected {
		t.Errorf("expected %s %04x, got %04x\n", context, expected, b)
	}
}

func TestDynamicAddrs(t *testing.T) {
	if err := getSynergyAddrs(); err != nil {
		t.Fatalf("Error when getting dynamic addrs: %v", err)
	}
	// other firmware versions may load things in other places, but mine
	// loads them as below:
	assertUint16(t, 0x0000, synAddrs.PROG,   "PROG")
	assertUint16(t, 0x5c72, synAddrs.ROM,    "ROM")
	assertUint16(t, 0x6033, synAddrs.VTAB,   "VTAB")
	assertUint16(t, 0x60e0, synAddrs.FILTAB, "FILTAB")
	assertUint16(t, 0x6300, synAddrs.EDATA,  "EDATA")
	assertUint16(t, 0x8000, synAddrs.RAM,    "RAM")
	assertUint16(t, 0x8044, synAddrs.PTSTAT, "PTSTAT")
	assertUint16(t, 0x86f9, synAddrs.SOLOSC, "SOLOSC")
	assertUint16(t, 0x8715, synAddrs.CODE,   "CODE")
	assertUint16(t, 0x8717, synAddrs.DEVICE, "DEVICE")
	assertUint16(t, 0x8719, synAddrs.VALUE,  "VALUE")
	assertUint16(t, 0x871b, synAddrs.PTVAL,  "PTVAL")
	assertUint16(t, 0x877a, synAddrs.SEQCON, "SEQCON")
	assertUint16(t, 0x87a2, synAddrs.SEQVOI, "SEQVOI")
	assertUint16(t, 0x882e, synAddrs.EXTRA,  "EXTRA")
	assertUint16(t, 0x88a9, synAddrs.TRANSP, "TRANSP")
	assertUint16(t, 0x88ae, synAddrs.SEQTAB, "SEQTAB")
	assertUint16(t, 0xf000, synAddrs.CMOS,   "CMOS")
}

func TestBlockDump(t *testing.T) {
	var syn_bytes []byte
	var err error

	if syn_bytes, err = BlockDump(0x6000, 41); err != nil {
		t.Fatalf("Error executing block dump: %v", err)
	}
	var expect_bytes = []byte("COPYRIGHT (C) 1982 DIGITAL KEYBOARDS INC.")
	
	if !reflect.DeepEqual(syn_bytes, expect_bytes) {
		t.Fatalf("dumped data doesnt match what we expect\n%v\n\n\n %v",syn_bytes,expect_bytes)
	}
}

func TestBlockLoad(t *testing.T) {
	var expect_bytes = []byte("Test Block Load")
	var len_expect uint16 = uint16(len(expect_bytes))
	var syn_bytes []byte
	var orig_bytes []byte
	var err error

	// we'll overwrite the top of sequencer data - need to take care not
	// to overwrite anything that affects basic event loop processing (else
	// the Synergy can't respond to next command).  
	var addr uint16 = synAddrs.SEQTAB 
	if orig_bytes, err = BlockDump(addr, len_expect); err != nil {
		t.Fatalf("Error executing block dump: %v", err)
	}
	if err = BlockLoad(addr, expect_bytes); err != nil {
		t.Fatalf("Error executing block load -- POWER CYCLE Synergy TO ENSURE DATA BACK TO NORMAL: %v", err)
	}
	if syn_bytes, err = BlockDump(addr, len_expect); err != nil {
		t.Fatalf("Error executing block dump -- POWER CYCLE Synergy TO ENSURE DATA BACK TO NORMAL: %v", err)
	}
	// restore the original data:
	if err = BlockLoad(addr, orig_bytes); err != nil {
		t.Fatalf("Error executing block load to restore data -- POWER CYCLE Synergy TO ENSURE DATA BACK TO NORMAL: %v", err)
	}

	if !reflect.DeepEqual(syn_bytes, expect_bytes) {
		t.Fatalf("dumped data doesnt match what we expect\n%s\n%v\n\n\n %s\n%v",string(syn_bytes),syn_bytes,string(expect_bytes),expect_bytes)
	}
}

func TestDumpByte(t *testing.T) {
	// first few bytes of the the copyright header
	var b byte
	var err error
	if b, err = DumpByte(0x6000, "get test byte0"); err != nil {
		t.Fatalf("Error dumping byte: %v", err)
	}
	if b != byte('C') {
		t.Fatalf("Dumped byte doesnt match expected value got %v expected %v", b, 'C')
	}
	if b, err = DumpByte(0x6001, "get test byte1"); err != nil {
		t.Fatalf("Error dumping byte: %v", err)
	}
	if b != byte('O') {
		t.Fatalf("Dumped byte doesnt match expected value got %v expected %v", b, 'O')
	}
	if b, err = DumpByte(0x6002, "get test byte2"); err != nil {
		t.Fatalf("Error dumping byte: %v", err)
	}
	if b != byte('P') {
		t.Fatalf("Dumped byte doesnt match expected value got %v expected %v", b, 'P')
	}
}

func TestLoadSaveSyn(t *testing.T) {
	var expect_bytes []byte
	var syn_bytes []byte
	var err error

	if expect_bytes, err = ioutil.ReadFile("../data/testfiles/TEST.SYN"); err != nil {
		t.Fatalf("Error when reading test file: %v", err)
	}

	if err = LoadSYN(expect_bytes); err != nil {
		t.Fatalf("Error calling LoadSYN: %v", err)
	}

	if syn_bytes, err = SaveSYN(); err != nil {
		t.Fatalf("Error calling SaveSYN: %v", err)
	}
	
	if !reflect.DeepEqual(syn_bytes, expect_bytes) {
		t.Fatalf("downloaded SYN doesnt match what we uploaded\n%v\n\n\n %v",syn_bytes,expect_bytes)
	}
}

func TestInitVRAM(t *testing.T) {
	var err error
	if err = InitVRAM(); err != nil {
		t.Fatalf("Error initializing VRAM: %v\n", err)
	}

	/*
	var midic_addr uint16 = 0xf400
	var b byte
	if b, err = DumpByte(midic_addr, "get MIDIC"); err != nil {
		t.Fatalf("Error getting MIDIC value: %v\n", err);
	}
	if b != 0xff {
		t.Fatalf("MIDIC not 0xff: got %x", b);
	}	
*/
}

func TestDisableVRAM(t *testing.T) {
	var err error
	if err = DisableVRAM(); err != nil {
		t.Fatalf("Error disabling VRAM: %v\n", err);
	}
/*
	var midic_addr uint16 = 0xf400
	var b byte
	if b, err = DumpByte(midic_addr, "get MIDIC"); err != nil {
		t.Fatalf("Error getting MIDIC value: %v\n", err);
	}
	if b != 0 {
		t.Fatalf("MIDIC not zero: got %x", b);
	}	
*/
}

func TestLoadCRT(t *testing.T) {
}


func TestLoadVCE(t *testing.T) {
}


func TestMain(m *testing.M) {
	flag.Parse()
	if *synio {
		err := connectToSynergy(); if err != nil {
			fmt.Printf("could not initialize io: %v\n", err)
			os.Exit(1)
		}
		os.Exit(m.Run())
	}	
	os.Exit(0)
}

































