package data

import (
	"bytes"
	"flag"
	"log"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var testfilepath = flag.String("testfilepath", "testfiles", "path to VCE and CRT files")


func testReadWriteVCE(t *testing.T, path string) {
	log.Println("test ", path);

	var read_bytes []byte
	var err error
	
	if read_bytes,err = ioutil.ReadFile(path); err != nil {
		t.Errorf("error reading %s: %v", path, err)
		return 
	}

	var readbuf = bytes.NewBuffer(read_bytes)
	var vce VCE
	
	if vce, err = ReadVce(readbuf, false); err != nil {
		t.Errorf("error parsing %s: %v", path, err)
		return
	}

	testVCEName(t, path, vce)
	
	var writebuf bytes.Buffer

	if err = WriteVce(&writebuf, vce, false); err != nil {
		t.Errorf("error writing %s: %v", path, err)
		return
	}
	write_bytes := writebuf.Bytes()
	
	if !reflect.DeepEqual(read_bytes, write_bytes) {
		// before we report an error, this might be the case thata the
		// original file just has some extra bytes at the end (G7S.VCE
		// has a bunch of null, some files have ASCII EOF chars).
		//
		// compare the filters if they are the same, we'll consider the
		// files identical

		var readbuf2 = bytes.NewBuffer(write_bytes)
		var vce2 VCE

		if vce2, err = ReadVce(readbuf2, false); err != nil {
			t.Errorf("error parsing generated stream: %v", err)
			return
		}
		if !reflect.DeepEqual(vce.Filters, vce2.Filters) {
			t.Errorf("read/write data doesnt match. read:\n%v\nfilters:%v\n\nwrote:\n %v\nfilters:%v",read_bytes,vce.Filters,write_bytes,vce2.Filters)
			return
		}		

		return
	}	
}

func testVCEName(t *testing.T, path string, vce VCE) {
	
	base := filepath.Base(path)
	if base != vceName(vce.Head)+".VCE" {
		t.Errorf("name doesnt match file - expected %s, got %s", base, vceName(vce.Head))
	}
}

func TestAllVCE(t *testing.T) {
	fileList := []string{}
	
	filepath.Walk(*testfilepath,
		func(path string, f os.FileInfo, err error) error {
			if (filepath.Ext(path) == ".VCE") {
				fileList = append(fileList, path)
			}
			return nil
		})
	for _,path := range fileList {
		testReadWriteVCE(t, path)
	}
	return
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
	return
}
