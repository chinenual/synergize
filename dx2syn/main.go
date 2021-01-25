package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/chinenual/synergize/data"
)

var allFlag = flag.Bool("all", false, "extract all patches")
var nameFlag = flag.String("name", "", "Name of the patch to extract")
var indexFlag = flag.Int("index", -1, "Index of the patch to extract")
var sysexFlag = flag.String("sysex", "", "Pathname of the sysex file to parse")
var verboseFlag = flag.Bool("verbose", false, "Verbose debugging")
var statsFlag = flag.Bool("stats", false, "print statistics about the sysex - dont generate vce")
var stevealgoFlag = flag.Bool("stevealgo", false, "create 32 different vce's - one for each algo")
var makecrtFlag = flag.String("makecrt", "", "Pathname of a directory containing VCEs")

func usage(msg string) {
	log.Printf("ERROR: %s\n", msg)
	log.Printf("Usage: \n")
	log.Printf("\tdx7-to-synergy ( -all | -index <n> | -name <name> ) <sysex-pathname>\n")
	os.Exit(1)
}

func main() {
	flag.Parse()

	if *makecrtFlag != "" {
		if err := makeCrt(*makecrtFlag); err != nil {
			log.Printf("ERROR: could not create CRT from %s: %v", *makecrtFlag, err)
		}
		return
	}
	if *stevealgoFlag {
		for a := 0; a < 32; a++ {
			var vce data.VCE
			var err error
			if vce, err = helperBlankVce(); err != nil {
				return
			}
			// DX7 always uses 6 oscillators
			vce.Head.VOITAB = 5
			if err = helperSetAlgorithmPatchType(&vce, byte(a), 0); err != nil {
				log.Printf("ERROR: could not set algo %d: %v", a, err)
			} else {
				vcePathname := fmt.Sprintf("algo%d.VCE", a)
				if err = data.WriteVceFile(vcePathname, vce, true); err != nil {
					log.Printf("ERROR: could not write VCEfile %s: %v", vcePathname, err)
				}
			}
		}
		return
	}

	if *sysexFlag == "" {
		usage("must specify a SYSEX file pathname")
	}
	if *allFlag && (*indexFlag >= 0 || *nameFlag != "") {
		usage("must specify only one of -all, -index or -name option")
	}
	if *indexFlag >= 0 && *nameFlag != "" {
		usage("must specify only one of -all, -index or -name option")
	}

	var err error
	var sysex Dx7Sysex
	var selectedVoices []Dx7Voice
	if sysex, err = ReadDx7Sysex(*sysexFlag); err != nil {
		log.Printf("ERROR: could not parse sysex file %s: %v", *sysexFlag, err)
		os.Exit(1)
	}

	if *allFlag {
		selectedVoices = sysex.Voices
	} else if *indexFlag >= 0 {
		selectedVoices = sysex.Voices[*indexFlag : *indexFlag+1]
	} else if *nameFlag != "" {
		// search for the name
		for i, v := range sysex.Voices {
			if v.VoiceName == *nameFlag {
				selectedVoices = sysex.Voices[i : i+1]
				break
			}
		}
		if selectedVoices == nil {
			log.Printf("ERROR: no such voice name '%s'. Valid names:\n", *nameFlag)
			for _, v := range sysex.Voices {
				log.Printf("'%s'\n", v.VoiceName)
			}
			os.Exit(1)
		}
	} else {
		usage("Must specify at least one of -all, -index or -name option")
	}

	//	fmt.Printf("Selected: %v\n", selectedVoices)
	if *statsFlag {
		log.Printf("sysex: %s\n", *sysexFlag)
	}
	for _, v := range selectedVoices {
		if *statsFlag {
			//log.Printf("feedback: %s %d\n", v.VoiceName, v.Feedback)
			log.Printf("COARSE:%s,  %d,  %d,  %d,  %d,  %d,  %d  \n", v.VoiceName, v.Osc[0].OscFreqCoarse, v.Osc[1].OscFreqCoarse, v.Osc[2].OscFreqCoarse, v.Osc[3].OscFreqCoarse, v.Osc[4].OscFreqCoarse, v.Osc[5].OscFreqCoarse)
			log.Printf("FINE:  %s,  %d,  %d,  %d,  %d,  %d, %d  \n", v.VoiceName, v.Osc[0].OscFreqFine, v.Osc[1].OscFreqFine, v.Osc[2].OscFreqFine, v.Osc[3].OscFreqFine, v.Osc[4].OscFreqFine, v.Osc[5].OscFreqFine)
			log.Printf(" \n")

			//log.Printf("OSCMode: %s %t %t %t %t %t %t \n", v.VoiceName, v.Osc[0].OscMode, v.Osc[1].OscMode, v.Osc[2].OscMode, v.Osc[3].OscMode, v.Osc[4].OscMode, v.Osc[5].OscMode)
			/*
				for i := 0; i < 6; i++ {
					log.Printf("KeyScale Break-LeftD-Rightd-LC-RC OP%d %s %d %d %d %d %d \n", i, v.VoiceName, v.Osc[i].KeyLevelScalingBreakPoint,
						v.Osc[i].KeyLevelScalingLeftDepth,
						v.Osc[i].KeyLevelScalingRightDepth,
						v.Osc[i].KeyLevelScalingLeftCurve,
						v.Osc[i].KeyLevelScalingRightCurve)
				}
			*/
		} else {
			hasError := false
			var vce data.VCE
			if *verboseFlag {
				log.Printf("Translating %s %s...\n", v.VoiceName, Dx7VoiceToJSON(v))
			} else {
				log.Printf("Translating %s...\n", v.VoiceName)
			}
			if vce, err = TranslateDx7ToVce(v); err != nil {
				log.Printf("ERROR: could not translate Dx7 voice %s: %v", v.VoiceName, err)
			} else {
				if *verboseFlag {
					log.Printf("Result VCE: %s %s\n", v.VoiceName, helperVCEToJSON(vce))
				}
				const IGNORE_VALIDATION = true
				if err = data.VceValidate(vce); (err != nil) && (!IGNORE_VALIDATION) {
					log.Printf("ERROR: validation error on translate Dx7 voice %s: %v\n", v.VoiceName, err)
					hasError = true
				} else {
					vcePathname := v.VoiceName + ".VCE"
					if err = data.WriteVceFile(vcePathname, vce, false); err != nil {
						log.Printf("ERROR: could not write VCEfile %s: %v\n", vcePathname, err)
						hasError = true
					}
				}
			}
			if hasError {
				log.Printf("ERROR: Error during conversion\n")
			}
		}
	}
	os.Exit(0)
}
