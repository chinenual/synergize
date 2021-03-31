package dx2syn

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chinenual/synergize/data"
	"github.com/orcaman/writerseeker"
)

func checkWriteCrtFromVCEArray(vces []*data.VCE) (err error) {
	var writebuf = writerseeker.WriterSeeker{}

	if err = data.WriteCrt(&writebuf, vces); err != nil {
		return
	}
	return
}

func isSysexExt(path string) bool {
	ext := strings.ToUpper(filepath.Ext(path))
	return ext == ".SYSX" || ext == ".SYSEX" || ext == ".SYX"
}

func MakeCrt(path string, verbose bool) (err error) {
	var fileInfo os.FileInfo
	if fileInfo, err = os.Stat(path); err != nil {
		return
	}
	if fileInfo.IsDir() {
		if err = recurseMakeCrt(path, verbose); err != nil {
			return
		}
	} else {
		if err = makeCrtFromSysex(path, verbose); err != nil {
			return
		}
	}
	return
}

func recurseMakeCrt(dirPath string, verbose bool) (err error) {
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && dirPath != path {
			if err = recurseMakeCrt(path, verbose); err != nil {
				return err
			}
		} else if isSysexExt(path) {
			if err = makeCrtFromSysex(path, verbose); err != nil {
				// Don't fail the whole job just because one SYX is bad
				// return err
			}
		}
		return nil
	})
	return
}
func makeCrtFromSysex(sysexPath string, verbose bool) (err error) {
	// make the vce's
	var sysex Dx7Sysex
	log.Printf("CONVERT VOICES FROM %s\n", sysexPath)
	if sysex, err = ReadDx7Sysex(sysexPath); err != nil {
		log.Printf("ERROR: could not parse sysex file %s: %v", sysexPath, err)
		return
	}
	nameMap := make(map[string]bool)

	for _, v := range sysex.Voices {
		hasError := false
		if verbose {
			log.Printf("Translating '%s' %s...\n", v.VoiceName, Dx7VoiceToJSON(v))
		} else {
			log.Printf("Translating '%s'...\n", v.VoiceName)
		}
		if _, err = TranslateDx7ToVceFile(sysexPath, verbose, &nameMap, v); err != nil {
			log.Printf("ERROR: could not translate Dx7 voice %s: %v", v.VoiceName, err)
			hasError = true
		}
		if hasError {
			log.Printf("ERROR: Error during conversion\n")
		}
	}

	ext := filepath.Ext(sysexPath)
	dirPath := sysexPath[:len(sysexPath)-len(ext)]
	if err = makeCrtFromSysexVces(dirPath); err != nil {
		return
	}
	return
}

func makeCrtFromSysexVces(sysexPath string) (err error) {

	log.Printf("MAKE CRT FROM %s\n", sysexPath)

	crtCount := 1
	crtPath := filepath.Join(sysexPath, filepath.Base(sysexPath)+".CRT")
	var vces []*data.VCE

	var fileInfos []os.FileInfo
	if fileInfos, err = ioutil.ReadDir(sysexPath); err != nil {
		return err
	}
	for _, f := range fileInfos {
		path := filepath.Join(sysexPath, f.Name())
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
						return
					}
					// initialize the new list of vces to just this next voice
					newVces = nil
					newVces = append(newVces, &vce)
					// start a new file
					crtCount++
					crtPath = filepath.Join(sysexPath, filepath.Base(sysexPath)+"-"+strconv.Itoa(crtCount)+".CRT")
					log.Printf("%s: Add %s ...\n", crtPath, path)
				} else {
					log.Printf("%s: Add %s ...\n", crtPath, path)
				}
				vces = newVces
			}

		}

	}
	if len(vces) > 1 {
		// write the last set of vce's that fit:
		if err = data.WriteCrtFileFromVCEArray(crtPath, vces); err != nil {
			return
		}
	}
	log.Printf("SUCCESS! MADE CRT FROM %s\n\n\n", sysexPath)
	return
}
