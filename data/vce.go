package data

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"io"
	"io/ioutil"
)

type FreqEnvelopeTable struct {
	OPTCH     byte
	OHARM     int8
	FDETUN    int8
	FENVL     byte
	ENVTYPE   byte
	NPOINTS   byte
	SUSTAINPT byte
	LOOPPT    byte
	Table     ArrayOfByte // force proper JSON encoding
}

type AmpEnvelopeTable struct {
	ENVTYPE   byte
	NPOINTS   byte
	SUSTAINPT byte
	LOOPPT    byte
	Table     ArrayOfByte // force proper JSON encoding
}

type Envelope struct {
	FreqEnvelope FreqEnvelopeTable
	AmpEnvelope  AmpEnvelopeTable
}

type VCE struct {
	Head VCEHead
	Envelopes []Envelope
	Filters [][32]int8;
}

type VCEHead struct {
	VOITAB byte
	OSCPTR [16]uint16
	VTRANS int8
	VTCENT byte
	VTSENS byte
	UNUSED byte
	VEQ    [24]int8
	VNAME  SpaceEncodedString // force string encoding for the name
	VACENT byte
	VASENS byte
	VIBRAT byte
	VIBDEL byte
	VIBDEP byte
	KPROP  [24]byte
	APVIB  byte
	FILTER [16]int8
}

func vceName(vceHead VCEHead) (name string) {
	name=""
	for i := 0; i < 8; i++ {
		if vceHead.VNAME[i] == ' ' {
			break;
		}
		name = name + string(vceHead.VNAME[i]);
	}
	return
}

func VceAFilterCount(vce VCE) (count int) {
	for i := byte(0); i < vce.Head.VOITAB; i++ {
		v := vce.Head.FILTER[i]
		if v < 0 {
			return 1
		}
	}
	return 0
}

func VceBFilterCount(vce VCE) (count int) {
	count = 0
	for i := byte(0); i < vce.Head.VOITAB; i++ {
		v := vce.Head.FILTER[i]
		if v > 0 {
			count = count + 1
		}
	}
	return
}

func ReadVceFile(filename string) (vce VCE, err error) {
	var b []byte

	if b,err = ioutil.ReadFile(filename); err != nil {
		return 
	}
	buf := bytes.NewReader(b)
	if vce,err = vceRead(buf, false); err != nil {
		return
	}
	return
}

func vceReadAFilters(buf io.Reader, vce *VCE) (err error) {
	// a voice has at most 1 A filter. These are at the head of the filter array
	// so we can unconditionally put it in slot 0 if there is one
	for _,f := range vce.Head.FILTER {
		if f < 0 {
			for j := 0; j < 32; j++ {
				if err = binary.Read(buf, binary.LittleEndian, &vce.Filters[0][j]); err != nil {
					log.Println(vceToString(*vce))
					log.Println("binary.Read failed:", 0, " ", j, " ", err)
					return
				}			
			}			
			return
		}
	}
	return
}

func vceReadBFilters(buf io.Reader, vce *VCE) (err error) {
	var filterCount = 0
	var hasAFilter = false;
	for _,f := range vce.Head.FILTER {
		if f != 0 {
			filterCount = filterCount + 1
		}
		if f < 0 {
			hasAFilter = true;
		}
	}

	var offset = 0
	if hasAFilter {
		offset = 1
	}
	for _,f := range vce.Head.FILTER {
		if f > 0 {
			// filters are one-based 
			var index = int(f) - 1 + offset
			for j := 0; j < 32; j++ {
				if err = binary.Read(buf, binary.LittleEndian, &vce.Filters[index][j]); err != nil {
					log.Println(vceToString(*vce))
					log.Println("binary.Read failed:", index, " ", j, " ", err)
					return
				}			
			}
		}
	}
	return
}
	
func vceRead(buf io.Reader, skipFilters bool) (vce VCE, err error) {
	if err = binary.Read(buf, binary.LittleEndian, &vce.Head); err != nil {
		log.Println("binary.Read failed:", err)
		return
	}

	vce.Envelopes = make ([]Envelope,vce.Head.VOITAB+1)
	for i := byte(0); i <= vce.Head.VOITAB; i++ {
		var e Envelope
	
		if err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.OPTCH); err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		if err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.OHARM); err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		if err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.FDETUN); err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		if err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.FENVL); err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		if err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.ENVTYPE); err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		if err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.NPOINTS); err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		if err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.SUSTAINPT); err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		if err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.LOOPPT); err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		// 4 values per point:
		e.FreqEnvelope.Table = make([]byte, e.FreqEnvelope.NPOINTS*4)
		for k := byte(0); k < e.FreqEnvelope.NPOINTS*4; k++ {
			if err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.Table[k]); err != nil {
				log.Println("binary.Read failed:", err)
				return
			}
		}
		if err = binary.Read(buf, binary.LittleEndian, &e.AmpEnvelope.ENVTYPE); err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		if err = binary.Read(buf, binary.LittleEndian, &e.AmpEnvelope.NPOINTS); err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		if err = binary.Read(buf, binary.LittleEndian, &e.AmpEnvelope.SUSTAINPT); err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		if err = binary.Read(buf, binary.LittleEndian, &e.AmpEnvelope.LOOPPT); err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		// 4 values per point:
		e.AmpEnvelope.Table = make([]byte, e.AmpEnvelope.NPOINTS*4)
		for k := byte(0); k < e.AmpEnvelope.NPOINTS*4; k++ {
			if err = binary.Read(buf, binary.LittleEndian, &e.AmpEnvelope.Table[k]); err != nil {
				log.Println("binary.Read failed:", err)
				return
			}
		}
		vce.Envelopes[i] = e
	}

	var filterCount = 0
	var hasAFilter = false
	for _,f := range vce.Head.FILTER {
		if f > 0 {
			filterCount++
		} else if f < 0 {
			hasAFilter = true;
		}
	}
	if hasAFilter {
		filterCount = filterCount+1;
	}
	vce.Filters = make ([][32]int8,filterCount);

	if ! skipFilters {
		if err = vceReadAFilters(buf, &vce); err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		if err = vceReadBFilters(buf, &vce); err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
	}
	return
}

func vceToString(vce VCE) (result string) {
       b,_ := json.MarshalIndent(vce, "", " ")
       result = string(b)

       return
}


func vceToJson(vce VCE) (result string) {
	b,_ := json.Marshal(vce)
	result = string(b)
	
	return
}
