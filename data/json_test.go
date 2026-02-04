package data

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/chinenual/synergize/logger"
	"github.com/rogpeppe/go-internal/diff"
)

func testOneString(t *testing.T, s string) bool {
	var ses SpaceEncodedString
	var err error
	var b []byte

	if err = ses.UnmarshalJSON([]byte(s)); err != nil {
		t.Error("Error unmarshalling", err)
	}
	if b, err = ses.MarshalJSON(); err != nil {
		t.Error("Error marshalling", err)
	}

	if string(b) != s {
		t.Errorf("Round trip unmarshal/marshal got \"%s\" expected \"%s\"\n", string(b), s)
		return false
	}
	return true
}

func TestSpaceEncodedString(t *testing.T) {
	testOneString(t, "\"12345678\"")
	testOneString(t, "\"A      Z\"")
	testOneString(t, "\"A       \"")
}

func TestVceJsonWrite(t *testing.T) {
	logger.Error("in testVceJsonWrite")
	var vce VCE
	var err error
	const VCE_PATH = "testfiles/G7S.VCE"
	const JSON_PATH = "testfiles/G7S.json"
	const TMP_JSON_PATH = "testfiles/G7S.json.test"
	if vce, err = ReadVceFile(VCE_PATH); err != nil {
		t.Errorf("Could not read VCE file %s\n", VCE_PATH)
	}

	var b []byte
	if b, err = json.MarshalIndent(vce, "", " "); err != nil {
		t.Error("Error formatting json", err)
	}
	if err = ioutil.WriteFile(TMP_JSON_PATH, b, 0644); err != nil {
		t.Error("Error saving json", err)
	}
	var expect_b []byte
	if expect_b, err = ioutil.ReadFile(JSON_PATH); err != nil {
		t.Errorf("Error loading json %#v: %v", JSON_PATH, err)
		return
	}

	if d := diff.Diff("generated", b, "expected", expect_b); d != nil {
		t.Errorf("JSON write data doesnt match:\n\n%s", string(d))
		return
	}
	return
}
func TestVceJsonRead(t *testing.T) {
	var vce VCE
	var vce_fromJson VCE
	var err error
	const VCE_PATH = "testfiles/G7S.VCE"
	const JSON_PATH = "testfiles/G7S.json"
	if vce, err = ReadVceFile(VCE_PATH); err != nil {
		t.Errorf("Could not read VCE file %s\n", VCE_PATH)
	}
	var b []byte
	if b, err = ioutil.ReadFile(JSON_PATH); err != nil {
		logger.Errorf("Error loading json %#v: %v", JSON_PATH, err)
		return
	}
	if err = json.Unmarshal(b, &vce_fromJson); err != nil {
		logger.Errorf("Error parsing json. %#v: %v", JSON_PATH, err)
		return
	}
	if !diffVCE(vce, vce_fromJson) {
		t.Errorf("VCE and JSON read data doesnt match")
		return
	}
	return
}

func TestCrtJsonWrite(t *testing.T) {
	var crt CRT
	var err error
	const CRT_PATH = "testfiles/INTERNAL.CRT"
	const JSON_PATH = "testfiles/INTERNAL.json"
	const TMP_JSON_PATH = "testfiles/INTERNAL.json.test"
	if crt, err = ReadCrtFile(CRT_PATH); err != nil {
		t.Errorf("Could not read CRT file %s\n", CRT_PATH)
	}

	var b []byte
	if b, err = json.MarshalIndent(crt, "", " "); err != nil {
		t.Error("Error formatting json", err)
	}
	if err = ioutil.WriteFile(TMP_JSON_PATH, b, 0644); err != nil {
		t.Error("Error saving json", err)
	}
	var expect_b []byte
	if expect_b, err = ioutil.ReadFile(JSON_PATH); err != nil {
		t.Errorf("Error loading json %#v: %v", JSON_PATH, err)
		return
	}

	if d := diff.Diff("generated", b, "expected", expect_b); d != nil {
		t.Errorf("JSON write data doesnt match:\n\n%s", string(d))
		return
	}
	return
}
func TestCrtJsonRead(t *testing.T) {
	var crt CRT
	var crt_fromJson CRT
	var err error
	const CRT_PATH = "testfiles/INTERNAL.CRT"
	const JSON_PATH = "testfiles/INTERNAL.json"
	if crt, err = ReadCrtFile(CRT_PATH); err != nil {
		t.Errorf("Could not read CRT file %s\n", CRT_PATH)
	}
	var b []byte
	if b, err = ioutil.ReadFile(JSON_PATH); err != nil {
		t.Errorf("Error loading json %#v: %v", JSON_PATH, err)
		return
	}
	if err = json.Unmarshal(b, &crt_fromJson); err != nil {
		t.Errorf("Error parsing json. %#v: %v", JSON_PATH, err)
		return
	}
	if !diffCRT(crt, crt_fromJson) {
		t.Error("CRT and JSON read data doesnt match")
		return
	}
	return
}
