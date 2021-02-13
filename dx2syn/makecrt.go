package main

import (
	"errors"
	"fmt"
	"github.com/chinenual/synergize/data"
	"github.com/orcaman/writerseeker"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func checkWriteCrtFromVCEArray(vces []*data.VCE) (err error) {
	var writebuf = writerseeker.WriterSeeker{}

	if err = data.WriteCrt(&writebuf, vces); err != nil {
		return
	}
	return
}

func makeCrt(dirPath string) (err error) {

	crtCount := 1
	crtPath := filepath.Join(dirPath, filepath.Base(dirPath)+".CRT")
	var vces []*data.VCE

	var fileInfos []os.FileInfo
	if fileInfos, err = ioutil.ReadDir(dirPath); err != nil {
		return err
	}
	for _, f := range fileInfos {
		path := filepath.Join(dirPath, f.Name())
		if strings.EqualFold(filepath.Ext(path), ".vce") {
			var vce data.VCE
			if vce, err = data.ReadVceFile(path); err != nil {
				return err
			} else {
				// will it fit?
				newVces := append(vces, &vce)
				if err = checkWriteCrtFromVCEArray(newVces); err != nil {
					// doesn't fit.  Write what we have to the current crtPath, then start a new
					// list of voices for a new crtPath
					if len(vces) < 1 {
						return errors.New("Must have at least one VCE file")
					}
					// write the last set of vce's that fit:
					if err = data.WriteCrtFileFromVCEArray(crtPath, vces); err != nil {
						fmt.Printf("diag %v \n", vces)
						return
					}
					// initialize the new list of vces to just this next voice
					newVces = nil
					newVces = append(newVces, &vce)
					// start a new file
					crtCount += 1
					crtPath = filepath.Join(dirPath, filepath.Base(dirPath)+"-"+strconv.Itoa(crtCount)+".CRT")
					log.Printf("%s: Add %s ...\n", crtPath, path)
				} else {
					log.Printf("%s: Add %s ...\n", crtPath, path)
				}
				vces = newVces
			}
		}
	}
	return
}
