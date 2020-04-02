package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"io/ioutil"
)

type AmpEnvelopeHead struct {
	OPTCH  uint8
	OHARM  int8
	FDETUN int8
	FENVL  uint8
}

type AmpEnvelope struct {
	Head AmpEnvelopeHead
//	Tables []AmpEnvelopeTable
}

/*
type VCE struct {
	Head VCEHead
	AmpEnvelopes []AmpEnvelope
}
*/

type VCE struct {
	VOITAB uint8
	OSCPTR [16]uint16
	VTRANS uint8
	VTCENT uint8
	VTSENS uint8
	UNUSED uint8
	VEQ    [24]uint8
	VNAME  [8]uint8
	VACENT uint8
	VASENS uint8
	VIBRAT uint8
	VIBDEL uint8
	KPROP  [24]uint8
	APVIB  uint8
	FILTER [16]int8
}

func Name(vce VCE) (name string) {
	name=""
	for i := 0; i < 8; i++ {
		if vce.VNAME[i] == ' ' {
			break;
		}
		name = name + string(vce.VNAME[i]);
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

	err = binary.Read(buf, binary.LittleEndian, &vce)
	if err != nil {
		log.Println("binary.Read failed:", err)
		return
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
