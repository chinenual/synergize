package main

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)


func testParseVCE(t *testing.T, path string) {
	log.Println("test ", path);
	
	vce, err := ReadVCEFile(path);
	if err != nil {
		t.Errorf("error parsing %s: %v", path, err)
		return
	}
	base := filepath.Base(path)
	if base != Name(vce)+".VCE" {
		t.Errorf("name doesnt match file - expected %s, got %s", base, Name(vce))
	}
}


func TestAllVCE(t *testing.T) {
	fileList := []string{}
	filepath.Walk("VOICES",
		func(path string, f os.FileInfo, err error) error {
			if (filepath.Ext(path) == ".VCE") {
				fileList = append(fileList, path)
			}
			return nil
		})
	for _,path := range fileList {
		testParseVCE(t, path);
	}
	return
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
	return
}
