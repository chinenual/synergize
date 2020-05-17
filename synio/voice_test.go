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

func TestLoadSyn(t *testing.T) {
	var syn_bytes []byte
	var err error

	if syn_bytes, err = ioutil.ReadFile("../data/testfiles/TEST.SYN"); err != nil {
		t.Fatalf("Error when reading test file: %v", err)
	}

	if err = LoadSYN(syn_bytes); err != nil {
		t.Fatalf("Error calling LoadSYN: %v", err)
	}
	return
}

func TestSaveSyn(t *testing.T) {
	var expect_bytes []byte
	var syn_bytes []byte
	var err error

	if expect_bytes, err = ioutil.ReadFile("../data/testfiles/TEST.SYN"); err != nil {
		t.Fatalf("Error when reading test file: %v", err)
	}

	if syn_bytes, err = SaveSYN(); err != nil {
		t.Fatalf("Error calling SaveSYN: %v", err)
	}

	if !reflect.DeepEqual(syn_bytes, expect_bytes) {
		t.Fatalf("downloaded SYN doesnt match what we uploaded\n%v\n\n\n %v",syn_bytes,expect_bytes)
	}
	return
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
