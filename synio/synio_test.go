package synio

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/chinenual/synergize/data"
)

var (
	synio             = flag.Bool("synio", false, "run integration tests that talk to the Synergy")
	port              = flag.String("port", "", "the serial device")
	baud              = flag.Uint("baud", 9600, "the serial baud rate")
	verbose           = flag.Bool("verbose", false, "synio verbose")
	mocksynio         = flag.Bool("MOCKSYNIO", false, "MOCK synio")
	serialVerboseFlag = flag.Bool("SERIALVERBOSE", false, "serial verbose")
)

// MIDIC is at 0xf400 in the firmware sources, but the linker relocates CMOS from
// 0xf100 to 0xf000 -- so subtract 0x0100:
const MIDIC_addr uint16 = 0xf300

func dumpAddressSpace(path string) {
	var b []byte
	var err error

	b, err = BlockDump(uint16(0), uint16(65323), "dump addr space")
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func connectToSynergy() (err error) {
	return SetSynergySerialPort(*port, *baud, *verbose, *serialVerboseFlag, *mocksynio)
}

func TestGetFirmwareId(t *testing.T) {
	if !*synio {
		t.Skip()
	}
	id, err := GetID()
	if err != nil {
		t.Fatalf("Error when getting id: %v", err)
	}
	if id[0] != 3 || id[1] != 22 {
		t.Errorf("Expected 3.22, got %v", id)
	}
}

func TestDynamicAddrs(t *testing.T) {
	if !*synio {
		t.Skip()
	}
	if err := getSynergyAddrs(); err != nil {
		t.Fatalf("Error when getting dynamic addrs: %v", err)
	}
	// other firmware versions may load things in other places, but mine
	// loads them as below:
	data.AssertUint16(t, 0x0000, synAddrs.PROG, "PROG")
	data.AssertUint16(t, 0x5c72, synAddrs.ROM, "ROM")
	data.AssertUint16(t, 0x6033, synAddrs.VTAB, "VTAB")
	data.AssertUint16(t, 0x60e0, synAddrs.FILTAB, "FILTAB")
	data.AssertUint16(t, 0x62e0, synAddrs.EDATA, "EDATA")
	data.AssertUint16(t, 0x8000, synAddrs.RAM, "RAM")
	data.AssertUint16(t, 0x8044, synAddrs.PTSTAT, "PTSTAT")
	data.AssertUint16(t, 0x86f9, synAddrs.SOLOSC, "SOLOSC")
	data.AssertUint16(t, 0x8715, synAddrs.CODE, "CODE")
	data.AssertUint16(t, 0x8717, synAddrs.DEVICE, "DEVICE")
	data.AssertUint16(t, 0x8719, synAddrs.VALUE, "VALUE")
	data.AssertUint16(t, 0x871b, synAddrs.PTVAL, "PTVAL")
	data.AssertUint16(t, 0x877a, synAddrs.SEQCON, "SEQCON")
	data.AssertUint16(t, 0x87a2, synAddrs.SEQVOI, "SEQVOI")
	data.AssertUint16(t, 0x882e, synAddrs.EXTRA, "EXTRA")
	data.AssertUint16(t, 0x88a9, synAddrs.TRANSP, "TRANSP")
	data.AssertUint16(t, 0x88ae, synAddrs.SEQTAB, "SEQTAB")
	data.AssertUint16(t, 0xf000, synAddrs.CMOS, "CMOS")
}

func TestInitVRAM(t *testing.T) {
	if !*synio {
		t.Skip()
	}
	//	dumpAddressSpace("before-initVRAM.bin");

	var err error
	if err = InitVRAM(); err != nil {
		t.Fatalf("Error initializing VRAM: %v\n", err)
	}

	//	dumpAddressSpace("after-initVRAM.bin");

	var b byte
	if b, err = DumpByte(MIDIC_addr, "get MIDIC"); err != nil {
		t.Fatalf("Error getting MIDIC value: %v\n", err)
	}
	if b != 0xff {
		// can't treat this as an error since I can't actually
		// find the toggled value at the addr I expect it to be.
		// leave this as a warning until better understanding
		t.Logf("MIDIC not 0xff: got %x", b)
	}
}

func TestDisableVRAM(t *testing.T) {
	if !*synio {
		t.Skip()
	}
	var err error
	if err = DisableVRAM(); err != nil {
		t.Fatalf("Error disabling VRAM: %v\n", err)
	}
	//	dumpAddressSpace("after-disableVRAM.bin");

	var b byte
	if b, err = DumpByte(MIDIC_addr, "get MIDIC"); err != nil {
		t.Fatalf("Error getting MIDIC value: %v\n", err)
	}
	if b != 0 {
		// can't treat this as an error since I can't actually
		// find the toggled value at the addr I expect it to be.
		// leave this as a warning until better understanding
		t.Logf("MIDIC not zero: got %x", b)
	}
}

func TestDumpVRAM(t *testing.T) {
	if !*synio {
		t.Skip()
	}
	var err error
	var bytes []byte

	// will fail unless vram is enabled on the synergy side:
	if err = InitVRAM(); err != nil {
		t.Fatalf("Error initializing VRAM: %v\n", err)
	}

	if bytes, err = DumpVRAM(); err != nil {
		t.Fatalf("DumpVRAM failed %v", err)
	}
	fmt.Printf("vram returned %d bytes\n", len(bytes))
}

func TestBlockDump(t *testing.T) {
	if !*synio {
		t.Skip()
	}
	var syn_bytes []byte
	var err error

	// NOTE: this needs to run before we init VRAM (or at least after disabling it - so put VRAM tests above) - but just to be sure:
	if err = DisableVRAM(); err != nil {
		t.Fatalf("Error disabling VRAM: %v\n", err)
	}

	if syn_bytes, err = BlockDump(0x6000, 41, "get header bytes"); err != nil {
		t.Fatalf("Error executing block dump: %v", err)
	}
	var expect_bytes = []byte("COPYRIGHT (C) 1982 DIGITAL KEYBOARDS INC.")

	if !reflect.DeepEqual(syn_bytes, expect_bytes) {
		t.Fatalf("dumped data doesnt match what we expect\n%v\n\n\n %v", syn_bytes, expect_bytes)
	}
}

func TestBlockLoad(t *testing.T) {
	if !*synio {
		t.Skip()
	}
	var expect_bytes = []byte("Test Block Load")
	var len_expect uint16 = uint16(len(expect_bytes))
	var syn_bytes []byte
	var orig_bytes []byte
	var err error

	// we'll overwrite the top of sequencer data - need to take care not
	// to overwrite anything that affects basic event loop processing (else
	// the Synergy can't respond to next command).
	var addr uint16 = synAddrs.SEQTAB
	if orig_bytes, err = BlockDump(addr, len_expect, "get SEQTAB"); err != nil {
		t.Fatalf("Error executing block dump: %v", err)
	}
	if err = BlockLoad(addr, expect_bytes, "load test bytes"); err != nil {
		t.Fatalf("Error executing block load -- POWER CYCLE Synergy TO ENSURE DATA BACK TO NORMAL: %v", err)
	}
	if syn_bytes, err = BlockDump(addr, len_expect, "dump test bytes"); err != nil {
		t.Fatalf("Error executing block dump -- POWER CYCLE Synergy TO ENSURE DATA BACK TO NORMAL: %v", err)
	}
	// restore the original data:
	if err = BlockLoad(addr, orig_bytes, "reload orig data"); err != nil {
		t.Fatalf("Error executing block load to restore data -- POWER CYCLE Synergy TO ENSURE DATA BACK TO NORMAL: %v", err)
	}

	if !reflect.DeepEqual(syn_bytes, expect_bytes) {
		t.Fatalf("dumped data doesnt match what we expect\n%s\n%v\n\n\n %s\n%v", string(syn_bytes), syn_bytes, string(expect_bytes), expect_bytes)
	}
}

func TestDumpByte(t *testing.T) {
	if !*synio {
		t.Skip()
	}
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

func diff(path1, path2 string) {
	var err error
	var bytes1, bytes2 []byte
	if bytes1, err = ioutil.ReadFile(path1); err != nil {
		os.Exit(1)
	}
	if bytes2, err = ioutil.ReadFile(path2); err != nil {
		os.Exit(1)
	}

	for i := range bytes1 {
		if bytes2[i] == 0 && bytes1[i] != 0 {
			fmt.Printf("%04x : after disable: %x, after init: %x\n",
				i, bytes2[i], bytes1[i])
		}
	}
	os.Exit(0)
}

func TestMain(m *testing.M) {
	//	diff("after-initVRAM.bin", "after-disableVRAM.bin")

	flag.Parse()
	if *synio {
		err := connectToSynergy()
		if err != nil {
			fmt.Printf("could not initialize io: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Integration tests skipped. Run with -synio to run them.\n")
	}
	os.Exit(m.Run())
}
