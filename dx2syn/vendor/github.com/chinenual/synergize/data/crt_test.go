package data

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/orcaman/writerseeker"
)

func loadVoiceList(path string) (list []string, err error) {
	var file *os.File
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

func compareVoiceList(t *testing.T, context string, crt CRT, list []string) {
	if len(crt.Voices) != len(list) {
		t.Errorf("%s: Voice list length mismatch - got %d, expected %d",
			context, len(crt.Voices), len(list))
		return
	}
	for i, voicename := range list {
		if voicename != VceName(crt.Voices[i].Head) {
			t.Errorf("%s: Voice name mismatch [%d] - got '%s', expected '%s'",
				context, i, VceName(crt.Voices[i].Head), voicename)
		}
	}
}

func testWriteCRT(t *testing.T, crt CRT) {
	// test creating a new CRT and then deep comparing it to the original.
	var writebuf = writerseeker.WriterSeeker{}

	var err error
	if err = WriteCrt(&writebuf, crt.Voices); err != nil {
		t.Errorf("error writing CRT: %v", err)
		return
	}

	write_bytes, _ := ioutil.ReadAll(writebuf.Reader())
	_ = dumpTestBytes("/tmp/testbytes.crt", write_bytes)

	var readbuf = bytes.NewReader(write_bytes)
	var crt2 CRT

	if crt2, err = ReadCrt(readbuf); err != nil {
		t.Errorf("error parsing generated stream: %v", err)
		return
	}
	if len(crt.Voices) != len(crt2.Voices) {
		t.Errorf("diff number of voices %d and %d", len(crt.Voices), len(crt2.Voices))
	}
	for i := range crt.Voices {
		if !diffVCE(*crt.Voices[i], *crt2.Voices[i]) {
			t.Errorf("read/write data doesnt match for voice %d", i+1)
			fmt.Printf("read: %s\n", VceToJson(*crt.Voices[i]))
			fmt.Printf("wrote: %s\n", VceToJson(*crt2.Voices[i]))
			return
		}
	}

}

func TestCreateCRT(t *testing.T) {
	var err error
	var list []string
	var voicelistPath = "testfiles/INTERNAL.CRT.voice-list.txt"
	list, err = loadVoiceList(voicelistPath)
	if err != nil {
		t.Logf("error parsing %s: %v - skipping contents check", voicelistPath, err)
		return
	}
	var vcePaths []string
	for _, name := range list {
		vcePaths = append(vcePaths, "testfiles/"+name+".VCE")
	}
	_ = WriteCrtFileFromVCEPaths("testfiles/gen/INTERNAL.CRT", vcePaths)
	_, _ = testParseCRT(t, "testfiles/gen/INTERNAL.CRT")
}

func testParseCRT(t *testing.T, path string) (crt CRT, err error) {
	log.Println("test ", path)

	crt, err = ReadCrtFile(path)
	if err != nil {
		t.Errorf("error parsing %s: %v", path, err)
		return
	}

	voicelistPath := path + ".voice-list.txt"
	var list []string
	list, err = loadVoiceList(voicelistPath)
	if err != nil {
		//t.Logf("error parsing %s: %v - skipping contents check", voicelistPath, err)
		// don't report the error back to the caller - need to continue with the Write Test
		err = nil
		return
	}

	compareVoiceList(t, path, crt, list)
	return
}

func TestAllCRT(t *testing.T) {
	fileList := []string{}
	_ = filepath.Walk(*testfilepath,
		func(path string, f os.FileInfo, err error) error {
			if filepath.Ext(path) == ".CRT" {
				fileList = append(fileList, path)
			}
			return nil
		})
	for _, path := range fileList {
		crt, err := testParseCRT(t, path)

		if err != nil {
			t.Errorf("cannot test write CRT since parser failed - %s %s\n", path, err)
		} else {
			testWriteCRT(t, crt)
		}
	}
}
