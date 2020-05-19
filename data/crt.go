package data

import (
	"encoding/binary"
	"bytes"
	"io"
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
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
	Voices []VCE
}


func ReadCrtFile(filename string) (crt CRT, err error) {
	// A CRT file is a long header containing filter info, followed by a list of CCE fragments (each voice missing the filter params since they are concatenated elsewhere in the file).
	
	var b []byte

	if b,err = ioutil.ReadFile(filename); err != nil {
		return 
	}
	buf := bytes.NewReader(b)

	if err = binary.Read(buf, binary.LittleEndian, &crt.Head); err != nil {
		log.Println("binary.Read failed:", err)
		return
	}

//	log.Println(crt.Head)

	// voice Offsets are from the VOIDTAB field
 	var voitabOffset uint16 = 50
	// filter Offsets are from the FILTAB field (after the last AFILTER entry)
 	var filtabOffset uint16 = 160
	
	for i,offset := range(crt.Head.VOIPTR) {
		if offset != 0 {
//		log.Printf("seek to %d\n",voitabOffset + offset)
		
			if _,err = buf.Seek(int64(voitabOffset+offset), io.SeekStart); err != nil {
				err = errors.Wrapf(err,"failed to seek to voice #%d start", i)
				return
			}
			var vce VCE
			if vce,err = ReadVce(buf, true); err != nil {
				err = errors.Wrapf(err,"failed to read voice #%d start", i)
				return
			}
			
			if VceAFilterCount(vce) > 0 {
				offset = uint16(crt.Head.AFILTR[i]-1) * 32
				// log.Printf("v:%d A FILTER index %d so offset %d from %x\n", i,crt.Head.AFILTR[i],offset,filtabOffset);
				if _,err = buf.Seek(int64(filtabOffset+offset), io.SeekStart); err != nil {
					err = errors.Wrapf(err,"failed to seek to voice #%d filter-b start", i)
					return
				}
				vceReadAFilters(buf, &vce)
			}
			if VceBFilterCount(vce) > 0 {
				offset = uint16(crt.Head.BFILTR[i]-1) * 32
				// log.Printf("v:%d B FILTER index %d so offset %d from %x\n", i,crt.Head.BFILTR[i],offset,filtabOffset);
				if _,err = buf.Seek(int64(filtabOffset+offset), io.SeekStart); err != nil {
					err = errors.Wrapf(err,"failed to seek to voice #%d filter-b start", i)
					return
				}
				vceReadBFilters(buf, &vce)
			}
			
			crt.Voices = append(crt.Voices, vce)
		}
	}
	
	return
}
