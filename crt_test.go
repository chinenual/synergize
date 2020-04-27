package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func loadVoiceList(path string) (list []string, err error) {
	var file *os.File;
	file, err = os.Open(path)
	if err != nil {
		return
	}
	
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	
	snl := bufio.NewScanner(file)
	for snl.Scan() {
		// handle blank line at end of file:
		name := strings.TrimSpace(snl.Text())		
		if "" != name {
			list = append(list, name)
		}
	}

	err = snl.Err()
	return
}

func compareVoiceList(t *testing.T, context string, crt CRT, list[]string) {
	if len(crt.Voices) != len(list) {
		t.Errorf("%s: Voice list length mismatch - got %d, expected %d",
			context, len(crt.Voices), len(list))
		return
	}
	for i,voicename := range(list) {
		if voicename != vceName(crt.Voices[i]) {
		t.Errorf("%s: Voice name mismatch [%d] - got '%s', expected '%s'",
			context, i, vceName(crt.Voices[i]), voicename)
		}
	}
}

func testParseCRT(t *testing.T, path string) {
	log.Println("test ", path);
	
	crt, err := crtReadFile(path);
	if err != nil {
		t.Errorf("error parsing %s: %v", path, err)
		return
	}
	
	voicelistPath := path + ".voice-list.txt";
	var list []string
	list, err = loadVoiceList(voicelistPath)
	if err != nil {
		t.Errorf("error parsing %s: %v", voicelistPath, err)
		return
	}

	compareVoiceList(t, path, crt, list)
}


func TestAllCRT(t *testing.T) {
	fileList := []string{}
	filepath.Walk("testfiles",
		func(path string, f os.FileInfo, err error) error {
			if (filepath.Ext(path) == ".CRT") {
				fileList = append(fileList, path)
			}
			return nil
		})
	for _,path := range fileList {
		testParseCRT(t, path);
	}
	return
}
