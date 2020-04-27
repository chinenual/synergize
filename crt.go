package main

import (
	"encoding/binary"
	"bytes"
	"io"
	"io/ioutil"
	"log"
)

type CRTHead struct {
	VRAM [51]byte  		// "free storage"
	VOITAB byte   		// always zero
	VCHK [5]byte		// check bytes - each should be 170
	TSTVEC uint16		// test vector
	IPDVEC uint16		// ret. from input w/ data
	IPNVEC uint16		// ret. from input w/out data
	VERSON byte     	// version
	VOIPTR [24]uint16	// 24 pointers to voice data (offsets from VOITAB)
	BFILTR [24]byte		// b-filter # start for voices
	AFILTR [24]int8		// a-filter # start for voices
}

type CRT struct {
	Head CRTHead
	Voices []VCEHead
}


func crtReadFile(filename string) (crt CRT, err error) {
	// A CRT file is a long header containing filter info, followed by a list of CCE fragments (each voice missing the filter params since they are concatenated elsewhere in the file).
	
	var b []byte

	b,err = ioutil.ReadFile(filename)
	if err != nil {
		return 
	}
	buf := bytes.NewReader(b)

	err = binary.Read(buf, binary.LittleEndian, &crt.Head)
	if err != nil {
		log.Println("binary.Read failed:", err)
		return
	}

	//	log.Println(crt.Head)

	// Offsets are from the VOIDTAB field
 	var startOffset uint16 = 50
	for i,offset := range(crt.Head.VOIPTR) {
//		log.Printf("seek to %d\n",startOffset + offset)
		
		_,err = buf.Seek(int64(startOffset+offset), io.SeekStart)
		if err != nil {
			log.Printf("failed to seek to voice #%d start: %v\n", i, err)
			return
		}
		var vceHead VCEHead
		err = vceReadHead(buf, &vceHead)
//		log.Printf("  HEAD: %v:  NAME: %s\n",vceHead, vceName(vceHead))
		crt.Voices = append(crt.Voices, vceHead)
	}
	
	return
}
