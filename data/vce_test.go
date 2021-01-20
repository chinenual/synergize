package data

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/orcaman/writerseeker"
)

var testfilepath = flag.String("testfilepath", "testfiles", "path to VCE and CRT files")

func testOnePathname(t *testing.T, path string, expect string) {
	AssertString(t, vceNameFromPathname(path), expect, path)
}

func TestVceNameFromPathname(t *testing.T) {
	testOnePathname(t, "/foo/bar/foo.vce", "FOO     ")
	testOnePathname(t, "/foo/bar/foo.VcE", "FOO     ")
	testOnePathname(t, "/foo/bar/foo.VCE", "FOO     ")
	testOnePathname(t, "C:\\foo\\bar\\FOO.VCE", "FOO     ")
	testOnePathname(t, "C:\\foo\\bar\\foo.VcE", "FOO     ")
	testOnePathname(t, "/foo/bar/foo", "FOO     ")
	testOnePathname(t, "/foo/bar/foo.baz", "FOO.BAZ ")
	testOnePathname(t, "/foo/bar/f123456789012345", "F1234567")
	testOnePathname(t, "/foo/bar/f123456789012345.vce", "F1234567")
}

func testReadWriteVCE(t *testing.T, path string) {
	log.Println("test ", path)

	var read_bytes []byte
	var err error

	if read_bytes, err = ioutil.ReadFile(path); err != nil {
		t.Errorf("error reading %s: %v", path, err)
		return
	}

	var readbuf = bytes.NewReader(read_bytes)
	var vce VCE

	if vce, err = ReadVce(readbuf, false); err != nil {
		t.Errorf("error parsing %s: %v", path, err)
		return
	}

	testVCEName(t, path, vce)

	var writebuf = writerseeker.WriterSeeker{}

	if err = WriteVce(&writebuf, vce, vceNameFromPathname(path), false); err != nil {
		t.Errorf("error writing %s: %v", path, err)
		return
	}
	write_bytes, _ := ioutil.ReadAll(writebuf.Reader())

	if !reflect.DeepEqual(read_bytes, write_bytes) {
		// before we report an error, this might be the case that the
		// original file just has some extra bytes at the end (G7S.VCE
		// has a bunch of null, some files have ASCII EOF chars).
		//
		// compare the parsed data: if they are the same, we'll consider the
		// files identical

		var readbuf2 = bytes.NewReader(write_bytes)
		var vce2 VCE

		if vce2, err = ReadVce(readbuf2, false); err != nil {
			t.Errorf("error parsing generated stream: %v", err)
			return
		}
		if !diffVCE(vce, vce2) {
			t.Errorf("read/write data doesnt match")
			return
		}

		return
	}
}

func testVCEName(t *testing.T, path string, vce VCE) {

	base := filepath.Base(path)
	if base != VceName(vce.Head)+".VCE" {
		t.Errorf("name doesnt match file - expected %s, got %s", base, VceName(vce.Head))
	}
}

func TestAllVCE(t *testing.T) {
	fileList := []string{}

	_ = filepath.Walk(*testfilepath,
		func(path string, f os.FileInfo, err error) error {
			if filepath.Ext(path) == ".VCE" {
				fileList = append(fileList, path)
			}
			return nil
		})
	for _, path := range fileList {
		testReadWriteVCE(t, path)
	}
}

func testVceValidate(t *testing.T, path string) {
	var vce VCE
	var err error
	if vce, err = ReadVceFile(path); err != nil {
		t.Errorf("Could not read VCE file %s\n", path)
	}
	if err = VceValidate(vce); err != nil {
		t.Errorf("unexpected validation error for %s: %v\n", path, err)
	}
	if false {
		// now inject some errors into the fields and check that the validator catches them:
		vce.Head.FILTER[0] = 100
		if err = VceValidate(vce); err != nil {
			fmt.Printf("Got expected error: %v\n", err)
		} else {
			t.Errorf("expected validation error FILTER[0] out of range for %s\n", path)
		}
	}
}

func TestVceValidate(t *testing.T) {
	fileList := []string{}

	_ = filepath.Walk(*testfilepath,
		func(path string, f os.FileInfo, err error) error {
			if filepath.Ext(path) == ".VCE" {
				fileList = append(fileList, path)
			}
			return nil
		})
	for _, path := range fileList {
		testVceValidate(t, path)
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}
