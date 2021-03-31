package dx2syn

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/chinenual/synergize/data"
)

const VOICENAME_EXT = ".txt"

var _voiceNotes string = ""

func VoiceNotesStart(sysexPath string, dxVoiceName string) (err error) {
	_voiceNotes = ""
	_voiceNotes += fmt.Sprintf("Sysx: \t\t\"%s\"\nDX Voice: \t\"%s\"\n", sysexPath, dxVoiceName)
	return
}

func VoiceNotesSynVNAME(vce data.VCE) (err error) {
	_voiceNotes += fmt.Sprintf("Synergy Voice: \t\"%s\"\n", data.VceName(vce.Head))
	return
}

func VoiceNotesAlgorithm(dxAlg byte) (err error) {
	_voiceNotes += fmt.Sprintf("\nDX Algorithm: \t%d\n", dxAlg+1)
	return
}

func VoiceNotesFeedback(feedback byte) (err error) {
	if feedback > 0 {
		_voiceNotes += fmt.Sprintf("\nNon-zero DX Feedback: %d.  Simulated via Triangle waveshape\n", feedback)
	}
	return
}

func VoiceNotesClose(vcePath string) (err error) {
	if _voiceNotes != "" {
		vceExt := path.Ext(vcePath)
		base := (vcePath)[0 : len(vcePath)-len(vceExt)]
		pathname := base + VOICENAME_EXT

		err = ioutil.WriteFile(pathname, []byte(_voiceNotes), 0777)
		if err != nil {
			return
		}
	}
	_voiceNotes = ""
	return
}
