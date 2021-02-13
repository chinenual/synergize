package main

import (
	"github.com/chinenual/synergize/data"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func makeCrt(dirPath string) (err error) {

	crtPath := filepath.Join(dirPath, filepath.Base(dirPath)+".CRT")
	var vces []*data.VCE

	var fileInfos []os.FileInfo
	if fileInfos, err = ioutil.ReadDir(dirPath); err != nil {
		return err
	}
	for _, f := range fileInfos {
		path := filepath.Join(dirPath, f.Name())
		if strings.EqualFold(filepath.Ext(path), ".vce") {
			if len(vces) < 24 {
				var vce data.VCE
				if vce, err = data.ReadVceFile(path); err != nil {
					return err
				} else {
					vces = append(vces, &vce)
					log.Printf("Add %s ...\n", path)
				}
			} else {
				log.Printf("WARNING: ignore extra VCEs in folder: %s\n", path)
			}
		}
	}

	if err = data.WriteCrtFileFromVCEArray(crtPath, vces); err != nil {
		return
	}
	return
}
