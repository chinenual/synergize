package data

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)


func testParseVCE(t *testing.T, path string) {
	log.Println("test ", path);
	
	vce, err := VceReadFile(path);
	if err != nil {
		t.Errorf("error parsing %s: %v", path, err)
		return
	}
	base := filepath.Base(path)
	if base != vceName(vce.Head)+".VCE" {
		t.Errorf("name doesnt match file - expected %s, got %s", base, vceName(vce.Head))
	}
}


func TestAllVCE(t *testing.T) {
	fileList := []string{}
	filepath.Walk("testfiles",
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
