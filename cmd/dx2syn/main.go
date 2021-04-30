package main

import (
	"flag"
	"os"

	"github.com/chinenual/synergize/logger"

	"github.com/chinenual/synergize/dx2syn"
)

var allFlag = flag.Bool("all", false, "extract all patches")
var nameFlag = flag.String("name", "", "Name of the patch to extract")
var indexFlag = flag.Int("index", -1, "Index of the patch to extract")
var sysexFlag = flag.String("sysex", "", "Pathname of the sysex file to parse")
var verboseFlag = flag.Bool("verbose", false, "Verbose debugging")
var statsFlag = flag.Bool("stats", false, "print statistics about the sysex - dont generate vce")

//var stevealgoFlag = flag.Bool("stevealgo", false, "create 32 different vce's - one for each algo")
var makecrtFlag = flag.String("makecrt", "", "Pathname of a directory containing VCEs")
var algoFlag = flag.Int("algo", -1, "DX Algorithm Number")
var loglevelFlag = flag.String("loglevel", "INFO", "Set log level (DEBUG,INFO,WARN or ERROR)")

func usage(msg string) {
	logger.Errorf("%s\n", msg)
	logger.Printf("Usage: \n")
	logger.Printf("\tdx2syn ( -all | -index <n> | -name <name> ) <sysex-pathname>\n")
	os.Exit(1)
}

func main() {
	flag.Parse()
	logger.InitViaString("", *loglevelFlag)

	if *makecrtFlag != "" {
		if err := dx2syn.MakeCrt(*makecrtFlag, *verboseFlag); err != nil {
			logger.Errorf("Could not create CRT from %s: %v", *makecrtFlag, err)
		}
		return
	}
	/***********
	if *stevealgoFlag {
		for a := 0; a < 32; a++ {
			var vce data.VCE
			var err error
			if vce, err = data.BlankVce(); err != nil {
				return
			}
			// DX7 always uses 6 oscillators
			vce.Head.VOITAB = 5
			if err = dx2syn.SetAlgorithmPatchType(&vce, byte(a), 0); err != nil {
				logger.Errorf("Could not set algo %d: %v", a, err)
			} else {
				vcePathname := fmt.Sprintf("algo%d.VCE", a)
				if err = data.WriteVceFile(vcePathname, vce, true); err != nil {
					logger.Errorf("Could not write VCEfile %s: %v", vcePathname, err)
				}
			}
		}
		return
	}
	*********/

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
	var sysex dx2syn.Dx7Sysex
	var selectedVoices []dx2syn.Dx7Voice
	if sysex, err = dx2syn.ReadDx7Sysex(*sysexFlag); err != nil {
		logger.Errorf("Could not parse sysex file %s: %v", *sysexFlag, err)
		os.Exit(1)
	}

	if *allFlag {
		selectedVoices = sysex.Voices
	} else if *indexFlag >= 0 {
		selectedVoices = sysex.Voices[*indexFlag : *indexFlag+1]
	} else if *algoFlag != -1 {
		// search for the algo
		for i, v := range sysex.Voices {
			if int(v.Algorithm) == (*algoFlag - 1) {
				selectedVoices = sysex.Voices[i : i+1]
				break
			}
		}
		if selectedVoices == nil {
			logger.Errorf("No such Algorithm '%d'. Valid Algorithms:\n", *algoFlag)
			for _, v := range sysex.Voices {
				logger.Printf("'%d'\n", v.Algorithm+1)
			}
			os.Exit(1)
		}

	} else if *nameFlag != "" {
		// search for the name
		for i, v := range sysex.Voices {
			if v.VoiceName == *nameFlag {
				selectedVoices = sysex.Voices[i : i+1]
				break
			}
		}
		if selectedVoices == nil {
			logger.Errorf("No such voice name '%s'. Valid names:\n", *nameFlag)
			for _, v := range sysex.Voices {
				logger.Printf("'%s'\n", v.VoiceName)
			}
			os.Exit(1)
		}
	} else {
		usage("Must specify at least one of -all, -index or -name option")
	}

	logger.Debugf("Selected: %v\n", selectedVoices)
	if *statsFlag {
		logger.Infof("sysex: %s\n", *sysexFlag)
	}

	nameMap := make(map[string]bool)

	for _, v := range selectedVoices {
		if *statsFlag {
			//logger.Infof("feedback: %s %d\n", v.VoiceName, v.Feedback
			logger.Infof("Trans:ALG:FB: %d,  %d,  %d  \n", v.Transpose, v.Algorithm, v.Feedback)
			logger.Infof("COARSE:  %d,  %d,  %d,  %d,  %d,  %d  \n", v.Osc[0].OscFreqCoarse, v.Osc[1].OscFreqCoarse, v.Osc[2].OscFreqCoarse, v.Osc[3].OscFreqCoarse, v.Osc[4].OscFreqCoarse, v.Osc[5].OscFreqCoarse)
			logger.Infof("FINE:    %d,  %d,  %d,  %d,  %d, %d  \n", v.Osc[0].OscFreqFine, v.Osc[1].OscFreqFine, v.Osc[2].OscFreqFine, v.Osc[3].OscFreqFine, v.Osc[4].OscFreqFine, v.Osc[5].OscFreqFine)
			logger.Infof("Volume:  %d,  %d,  %d,  %d,  %d, %d  \n", v.Osc[0].OperatorOutputLevel, v.Osc[1].OperatorOutputLevel, v.Osc[2].OperatorOutputLevel, v.Osc[3].OperatorOutputLevel, v.Osc[4].OperatorOutputLevel, v.Osc[5].OperatorOutputLevel)
			logger.Infof(" \n")

			//logger.Infof("OSCMode: %s %t %t %t %t %t %t \n", v.VoiceName, v.Osc[0].OscMode, v.Osc[1].OscMode, v.Osc[2].OscMode, v.Osc[3].OscMode, v.Osc[4].OscMode, v.Osc[5].OscMode)
			/*
				for i := 0; i < 6; i++ {
					logger.Infof("KeyScale Break-LeftD-Rightd-LC-RC OP%d %s %d %d %d %d %d \n", i, v.VoiceName, v.Osc[i].KeyLevelScalingBreakPoint,
						v.Osc[i].KeyLevelScalingLeftDepth,
						v.Osc[i].KeyLevelScalingRightDepth,
						v.Osc[i].KeyLevelScalingLeftCurve,
						v.Osc[i].KeyLevelScalingRightCurve)
				}
			*/
		} else {
			hasError := false
			if *verboseFlag {
				logger.Infof("Translating '%s' %s...\n", v.VoiceName, dx2syn.Dx7VoiceToJSON(v))
			} else {
				logger.Debugf("Translating '%s'...\n", v.VoiceName)
			}
			if _, err = dx2syn.TranslateDx7ToVceFile(*sysexFlag, *verboseFlag, &nameMap, v); err != nil {
				logger.Errorf("Could not translate Dx7 voice %s: %v", v.VoiceName, err)
				hasError = true
			}
			if hasError {
				logger.Errorf("Error during conversion\n")
			}
		}
	}
	os.Exit(0)
}
