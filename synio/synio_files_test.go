package synio

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestLoadSaveSyn(t *testing.T) {
	if !*synio {
		t.Skip()
	}
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
		t.Fatalf("downloaded SYN doesnt match what we uploaded\n%v\n\n\n %v", syn_bytes, expect_bytes)
	}
}

func TestLoadCRT(t *testing.T) {
	if !*synio {
		t.Skip()
	}
	var err error
	var bytes []byte

	// FIXME: probably shoudnt be using the test files from the data package
	if bytes, err = ioutil.ReadFile("../data/testfiles/INTERNAL.CRT"); err != nil {
		t.Fatalf("Can't load test data %v", err)
	}

	if err = LoadCRTBytes(bytes); err != nil {
		t.Fatalf("LoadCRTBytes failed %v", err)
	}
}
