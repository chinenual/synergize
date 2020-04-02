package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"io/ioutil"
)

type FreqEnvelopeTable struct {
	OPTCH     uint8
	OHARM     int8
	FDETUN    int8
	FENVL     uint8
	ENVTYPE   uint8
	NPOINTS   uint8
	SUSTAINPT uint8
	LOOPPT    uint8
	Table     []uint8
}

type AmpEnvelopeTable struct {
	ENVTYPE   uint8
	NPOINTS   uint8
	SUSTAINPT uint8
	LOOPPT    uint8
	Table     []uint8
}

type Envelope struct {
	FreqEnvelope FreqEnvelopeTable
	AmpEnvelope  AmpEnvelopeTable
}

type VCE struct {
	Head VCEHead
	Envelopes []Envelope
	Filters [][32]uint8;
}

type VCEHead struct {
	VOITAB uint8
	OSCPTR [16]uint16
	VTRANS uint8
	VTCENT uint8
	VTSENS uint8
	UNUSED uint8
	VEQ    [24]int8
	VNAME  [8]uint8
	VACENT uint8
	VASENS uint8
	VIBRAT uint8
	VIBDEL uint8
	VIBDEP uint8
	KPROP  [24]uint8
	APVIB  uint8
	FILTER [16]int8
}

func Name(vce VCE) (name string) {
	name=""
	for i := 0; i < 8; i++ {
		if vce.Head.VNAME[i] == ' ' {
			break;
		}
		name = name + string(vce.Head.VNAME[i]);
	}
	return
}

func ReadVCEFile(filename string) (vce VCE, err error) {
	var b []byte

	b,err = ioutil.ReadFile(filename)
	if err != nil {
		return 
	}
	buf := bytes.NewReader(b)

	err = binary.Read(buf, binary.LittleEndian, &vce.Head)
	if err != nil {
		log.Println("binary.Read failed:", err)
		return
	}

	vce.Envelopes = make ([]Envelope,vce.Head.VOITAB+1)
	for i := uint8(0); i <= vce.Head.VOITAB; i++ {
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
		e.FreqEnvelope.Table = make([]uint8, e.FreqEnvelope.NPOINTS*4)
		for k := uint8(0); k < e.FreqEnvelope.NPOINTS*4; k++ {
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
		e.AmpEnvelope.Table = make([]uint8, e.AmpEnvelope.NPOINTS*4)
		for k := uint8(0); k < e.AmpEnvelope.NPOINTS*4; k++ {
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
	vce.Filters = make ([][32]uint8,filterCount);

	for i := 0; i < filterCount; i++ {
		for j := 0; j < 32; j++ {
			err = binary.Read(buf, binary.LittleEndian, &vce.Filters[i][j])
			if err != nil {
				log.Println("binary.Read failed:", err)
				return
			}			
		}
	}
	return
}

func PrintVCEFile(filename string) (err error) {
	var vce VCE
	vce, err = ReadVCEFile(filename);
	if err != nil {
		return
	}
	log.Printf("%+v\n", vce)
	return
}
