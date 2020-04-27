package main

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
	Filters [][32]byte;
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

func vceReadHead(buf io.Reader, head *VCEHead) (err error) {
	err = binary.Read(buf, binary.LittleEndian, head)
	return
}

func vceReadFile(filename string) (vce VCE, err error) {
	var b []byte

	b,err = ioutil.ReadFile(filename)
	if err != nil {
		return 
	}
	buf := bytes.NewReader(b)

	err = vceReadHead(buf, &vce.Head)

	if err != nil {
		log.Println("binary.Read failed:", err)
		return
	}

	vce.Envelopes = make ([]Envelope,vce.Head.VOITAB+1)
	for i := byte(0); i <= vce.Head.VOITAB; i++ {
		var e Envelope
	
		err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.OPTCH)
		if err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.OHARM)
		if err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.FDETUN)
		if err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.FENVL)
		if err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.ENVTYPE)
		if err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.NPOINTS)
		if err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.SUSTAINPT)
		if err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.LOOPPT)
		if err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		// 4 values per point:
		e.FreqEnvelope.Table = make([]byte, e.FreqEnvelope.NPOINTS*4)
		for k := byte(0); k < e.FreqEnvelope.NPOINTS*4; k++ {
			err = binary.Read(buf, binary.LittleEndian, &e.FreqEnvelope.Table[k])
			if err != nil {
				log.Println("binary.Read failed:", err)
				return
			}
		}
		err = binary.Read(buf, binary.LittleEndian, &e.AmpEnvelope.ENVTYPE)
		if err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		err = binary.Read(buf, binary.LittleEndian, &e.AmpEnvelope.NPOINTS)
		if err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		err = binary.Read(buf, binary.LittleEndian, &e.AmpEnvelope.SUSTAINPT)
		if err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		err = binary.Read(buf, binary.LittleEndian, &e.AmpEnvelope.LOOPPT)
		if err != nil {
			log.Println("binary.Read failed:", err)
			return
		}
		// 4 values per point:
		e.AmpEnvelope.Table = make([]byte, e.AmpEnvelope.NPOINTS*4)
		for k := byte(0); k < e.AmpEnvelope.NPOINTS*4; k++ {
			err = binary.Read(buf, binary.LittleEndian, &e.AmpEnvelope.Table[k])
			if err != nil {
				log.Println("binary.Read failed:", err)
				return
			}
		}
		vce.Envelopes[i] = e
	}
	var filterCount = 0
	for _,f := range vce.Head.FILTER {
		if f != 0 {
			filterCount++
		}
	}
	vce.Filters = make ([][32]byte,filterCount);

	for i := 0; i < filterCount; i++ {
		for j := 0; j < 32; j++ {
			err = binary.Read(buf, binary.LittleEndian, &vce.Filters[i][j])
			if err != nil {
				log.Println(vceToString(vce))
				log.Println("binary.Read failed:", i, " ", j, " ", err)
				return
			}			
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
