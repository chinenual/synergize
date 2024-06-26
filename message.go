package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/chinenual/synergize/seq"

	"github.com/chinenual/synergize/data"
	"github.com/chinenual/synergize/io"
	"github.com/chinenual/synergize/logger"
	"github.com/chinenual/synergize/osc"
	"github.com/chinenual/synergize/synio"
	"github.com/chinenual/synergize/zeroconf"

	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"github.com/pkg/errors"
)

type connectionStatusResponse struct {
	SynergyName        string
	ControlSurfaceName string
}

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	logger.Debugf("Handle message: %s %s\n", m.Name, string(m.Payload))

	defer func() {
		if err != nil {
			logger.Errorf("Error handling message: {Name: %s, Payload: %s}: %v", m.Name, string(m.Payload), err)
		}
	}()

	switch m.Name {
	case "cancelPreferences":
		_ = prefs_w.Hide()

	case "checkVersion":
		var args struct {
			SynergyWasDisconnected        bool
			ControlSurfaceWasDisconnected bool
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		if args.SynergyWasDisconnected || (args.ControlSurfaceWasDisconnected && prefsUserPreferences.UseOsc) {
			CheckForNewVersion(true, io.SynergyConnectionType(), osc.ControlSurfaceConfigured())
		}
		payload = "ok"

	case "connectSynergy":
		var args struct {
			ZeroconfChoice *zeroconf.Service
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}

		var alreadyConfgured = io.SynergyConfigured()
		if args.ZeroconfChoice != nil {
			logger.Infof("ZEROCONF: config Synergy selected by user: %#v\n", *args.ZeroconfChoice)
			if err = ConnectToSynergy(args.ZeroconfChoice); err != nil {
				payload = err.Error()
				return
			}
		} else if !alreadyConfgured {
			err = errors.New("invalid argument to ConnectSynergy")
			payload = err.Error()
			return
		}
		type responseType struct {
			AlreadyConnected bool
			Status           connectionStatusResponse
		}
		response := responseType{
			AlreadyConnected: alreadyConfgured,
			Status: connectionStatusResponse{
				SynergyName:        io.SynergyName(),
				ControlSurfaceName: osc.ControlSurfaceName(),
			}}
		payload = response

	case "crtEditAddVoice":
		var args struct {
			Crt     data.CRT
			VcePath string
			Slot    int
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		var vce data.VCE
		if vce, err = data.ReadVceFile(args.VcePath); err != nil {
			payload = err.Error()
			return
		} else {
			logger.Infof("Add vce %s to CRT at slot %d\n", args.VcePath, args.Slot)
			if len(args.Crt.Voices) < args.Slot {
 				// grow the slice
				newVoices := make([]*data.VCE, args.Slot)
				copy(newVoices, args.Crt.Voices)
				args.Crt.Voices = newVoices
			}
 			args.Crt.Voices[args.Slot-1] = &vce
			payload = args.Crt
		}

	case "crtEditLoadCRT":
		var args struct {
			Crt data.CRT
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		if err = synio.LoadCRT(args.Crt); err != nil {
			payload = err.Error()
			return
		}
		payload = "ok"

	case "crtEditSaveCRT":
		var args struct {
			Path string
			Crt  data.CRT
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		if err = data.WriteCrtFileFromVCEArray(args.Path, args.Crt.Voices); err != nil {
			payload = err.Error()
			return
		}
		payload = "ok"

	case "disableVRAM":
		if err = synio.DisableVRAM(); err != nil {
			payload = err.Error()
			return
		} else {
			payload = "ok"
		}

	case "disconnectControlSurface":
		if err = DisconnectControlSurface(); err != nil {
			payload = err.Error()
		} else {
			response := connectionStatusResponse{io.SynergyName(), osc.ControlSurfaceName()}
			payload = response
		}

	case "disconnectSynergy":
		if err = DisconnectSynergy(); err != nil {
			payload = err.Error()
		} else {
			response := connectionStatusResponse{io.SynergyName(), osc.ControlSurfaceName()}
			payload = response
		}

	case "dx2synCancel":
		if err = dx2SynProcessCancel(); err != nil {
			payload = err.Error()
		} else {
			payload = "Ok"
		}

	case "dx2synStart":
		var args struct {
			Path string
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		if err = dx2synProcessStart(args.Path); err != nil {
			payload = err.Error()
		} else {
			payload = "Ok"
		}

	case "syn2midi":
		var args struct {
			Path            string
			Tempo           float64
			Raw             bool
			MaxClockSeconds uint32
			TrackButtons    [4]seq.TrackPlayMode
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		logger.Infof("syn2midi: %v\n", args)

		if err = seq.ConvertSYNToMIDI(args.Path, seq.TrackPerVoice, args.Tempo, args.Raw, args.MaxClockSeconds*1000, args.TrackButtons); err != nil {
			payload = err.Error()
		} else {
			payload = "Ok"
		}
	case "getSynSequencerState":
		var args struct {
			Path string
		}
		var response struct {
			TrackButtons [4]seq.TrackPlayMode
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		if response.TrackButtons, err = seq.GetSYNSequencerState(args.Path); err != nil {
			payload = err.Error()
		} else {
			logger.Infof("GetSYNSequencerState: %v\n", response)
			payload = response
		}

	case "getCWD":
		payload, _ = os.Getwd()
		logger.Infof("CWD: %s\n", payload)

	case "getConnectionStatus":
		response := connectionStatusResponse{io.SynergyName(), osc.ControlSurfaceName()}
		payload = response

	case "getPatchTypeNames":
		response := data.PatchTypeNames
		payload = response

	case "getPreferences":
		payload = struct {
			Os          string
			Preferences Preferences
		}{runtime.GOOS, prefsUserPreferences}

	case "getVersion":
		payload = struct {
			Version             string
			NewVersionAvailable bool
		}{AppVersion, CheckForNewVersion(false, io.SynergyConnectionType(), osc.ControlSurfaceConfigured())}

	case "isHTTPDebug":
		payload = prefsUserPreferences.HTTPDebug

	case "loadCRT":
		var path string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &path); err != nil {
				payload = err.Error()
				return
			}
		}
		if err = diagLoadCRT(path); err != nil {
			payload = err.Error()
			return
		} else {
			payload = "ok"
		}

	case "loadSYN":
		var path string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &path); err != nil {
				payload = err.Error()
				return
			}
		}
		if err = diagLoadSYN(path); err != nil {
			payload = err.Error()
			return
		} else {
			payload = "ok"
		}

	case "loadVceVoicingMode":
		var vce data.VCE
		var path string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &path); err != nil {
				payload = err.Error()
				return
			}
		}
		if vce, err = data.ReadVceFile(path); err != nil {
			payload = err.Error()
			return
		}
		if err = synio.LoadVceVoicingMode(vce); err != nil {
			payload = err.Error()
			return
		} else {
			// NOTE: need to pass reference in order to get the custom JSON marshalling to notice the VNAME
			payload = &vce
		}

	case "readCRT":
		var crt data.CRT
		var path string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &path); err != nil {
				payload = err.Error()
				return
			}
		}
		if crt, err = data.ReadCrtFile(path); err != nil {
			payload = err.Error()
			return
		} else {
			payload = crt
		}

	case "readVCE":
		var vce data.VCE
		var path string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &path); err != nil {
				payload = err.Error()
				return
			}
		}
		if vce, err = data.ReadVceFile(path); err != nil {
			payload = err.Error()
			return
		} else {
			// NOTE: need to pass reference in order to get the custom JSON marshalling to notice the VNAME
			payload = &vce
		}

	case "runCOMTST":
		// nothing interesting in the payload - just start the test and return results
		if err = synio.DiagCOMTST(); err != nil {
			payload = err.Error()
			return
		} else {
			payload = "Success!"
		}
	case "savePreferences":
		oldPath := prefsUserPreferences.LibraryPath
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &prefsUserPreferences); err != nil {
				payload = err.Error()
				return
			}
		}
		if err = prefsSavePreferences(); err != nil {
			payload = err.Error()
			return
		}
		if oldPath != prefsUserPreferences.LibraryPath {
			refreshNavPane(prefsUserPreferences.LibraryPath)
		}

		if prefsUserPreferences.UseOsc {
			if err := zeroconf.StartServer(prefsUserPreferences.OscPort, prefsSynergyName()); err != nil {
				logger.Errorf("could not start zeroconf: %v\n", err)
			}
		} else {
			zeroconf.CloseServer()
		}
		if (!zeroconf.ListenerRunning()) &&
			(prefsUserPreferences.UseOsc && prefsUserPreferences.OscAutoConfig) {
			zeroconf.StartListener()
		}
		_ = prefs_w.Hide()
		payload = "ok"

	case "saveSYN":
		var path string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &path); err != nil {
				payload = err.Error()
				return
			}
		}
		if err = diagSaveSYN(path); err != nil {
			payload = err.Error()
			return
		} else {
			payload = "ok"
		}

	case "saveVCE":
		var path string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &path); err != nil {
				payload = err.Error()
				return
			}
		}
		if err = diagSaveVCE(path); err != nil {
			payload = err.Error()
			return
		} else {
			payload = "ok"
		}

	case "setEnvFreqLowVal",
		"setEnvFreqUpVal",
		"setEnvAmpLowVal",
		"setEnvAmpUpVal",
		"setEnvFreqLowTime",
		"setEnvFreqUpTime",
		"setEnvAmpLowTime",
		"setEnvAmpUpTime":
		var args struct {
			Osc   int
			Index int
			Value int
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		switch m.Name {
		case "setEnvFreqLowVal":
			err = synio.SetEnvFreqLowVal(args.Osc, args.Index, byte(args.Value))
		case "setEnvFreqUpVal":
			err = synio.SetEnvFreqUpVal(args.Osc, args.Index, byte(args.Value))
		case "setEnvAmpLowVal":
			err = synio.SetEnvAmpLowVal(args.Osc, args.Index, byte(args.Value))
		case "setEnvAmpUpVal":
			err = synio.SetEnvAmpUpVal(args.Osc, args.Index, byte(args.Value))
		case "setEnvFreqLowTime":
			err = synio.SetEnvFreqLowTime(args.Osc, args.Index, byte(args.Value))
		case "setEnvFreqUpTime":
			err = synio.SetEnvFreqUpTime(args.Osc, args.Index, byte(args.Value))
		case "setEnvAmpLowTime":
			err = synio.SetEnvAmpLowTime(args.Osc, args.Index, byte(args.Value))
		case "setEnvAmpUpTime":
			err = synio.SetEnvAmpUpTime(args.Osc, args.Index, byte(args.Value))
		}
		if err != nil {
			payload = err.Error()
			return
		}
		payload = "ok"

	case "sendToCSurface":
		var args struct {
			Field string
			Value int
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		if err = osc.OscSendToCSurface(args.Field, args.Value); err != nil {
			logger.Errorf("Error sending to csurface: %v\n", err)
			payload = err.Error()
			return
		}
		payload = "ok"

	case "setEnvelopes":
		var args struct {
			Osc       int // 1-based
			Envelopes data.Envelope
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		if err = synio.SetEnvelopes(args.Osc, args.Envelopes); err != nil {
			payload = err.Error()
			return
		}
		payload = "ok"

	case "setFilterArray":
		var args struct {
			UiFilterIndex int // 0=Af, 1..16=Bf
			Values        []int
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		if err = synio.SetFilterArray(args.UiFilterIndex, args.Values); err != nil {
			payload = err.Error()
			return
		}
		payload = "ok"

	case "setFilterEle":
		var args struct {
			UiFilterIndex int // 0=Af, 1..16=Bf
			Index         int
			Value         int
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		if err = synio.SetFilterEle(args.UiFilterIndex, args.Index, args.Value); err != nil {
			payload = err.Error()
			return
		}
		payload = "ok"

	case "setLoopPoint":
		var args struct {
			Osc       int
			Env       string
			EnvType   int
			SustainPt int
			LoopPt    int
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		if err = synio.SetEnvLoopPoint(args.Osc, args.Env, args.EnvType, args.SustainPt, args.LoopPt); err != nil {
			payload = err.Error()
			return
		}
		payload = "ok"

	case "setNumOscillators":
		var args struct {
			NumOsc    int
			PatchType int
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		var resultPayload struct {
			EnvelopeTemplate data.Envelope
			PatchBytes       [16]byte
		}
		if resultPayload.PatchBytes, err = synio.SetNumOscillators(args.NumOsc, args.PatchType); err != nil {
			payload = err.Error()
			return
		}
		resultPayload.EnvelopeTemplate = data.DefaultEnvelope
		payload = resultPayload

	case "setOscEnvLengths":
		var args struct {
			Osc        int
			FreqLength int
			AmpLength  int
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		if err = synio.SetOscEnvLengths(args.Osc, args.FreqLength, args.AmpLength); err != nil {
			payload = err.Error()
			return
		}
		payload = "ok"

	case "setOscFILTER":
		var args struct {
			Args []int
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		if err = synio.SetOscFILTER(args.Args[0], args.Args[1]); err != nil {
			payload = err.Error()
			return
		}
		payload = "ok"

	case "setOscKEYPROP":
		var args struct {
			Args []int
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		var val = false
		if args.Args[1] == 1 {
			val = true
		}
		if err = synio.SetOscKEYPROP(args.Args[0], val); err != nil {
			payload = err.Error()
			return
		}
		payload = "ok"

	case "setOscSolo":
		var args struct {
			Mute []bool
			Solo []bool
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		if payload, err = synio.SetOscSolo(args.Mute, args.Solo); err != nil {
			payload = err.Error()
			return
		}

	case "setOscWAVE":
		var args struct {
			Args []int
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		var val = false
		if args.Args[1] == 1 {
			val = true
		}
		if err = synio.SetOscWAVE(args.Args[0], val); err != nil {
			payload = err.Error()
			return
		}
		payload = "ok"

	case "setPatchByte":
		var args struct {
			Osc   int
			Value int
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		if err = synio.SetVoiceOscDataByte(args.Osc, "OPTCH_reloadGenerators", byte(args.Value)); err != nil {
			payload = err.Error()
			return
		}
		payload = "ok"

	case "setPatchType":
		var index int
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &index); err != nil {
				payload = err.Error()
				return
			}
		}
		if payload, err = synio.SetPatchType(index); err != nil {
			payload = err.Error()
			return
		}

	case "setVNAME":
		var args struct {
			Param string
			Args  string // HACK: just a string - the JS code shares some logic with the other voice bytes that use setVoiceByte
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		if err = synio.SetVNAME(args.Args); err != nil {
			payload = err.Error()
			return
		}
		payload = "ok"

	case "setVoiceByte":
		var args struct {
			Param string
			Args  []int
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		if len(args.Args) == 2 {
			if err = synio.SetVoiceOscDataByte(args.Args[0], args.Param, byte(args.Args[1])); err != nil {
				payload = err.Error()
				return
			}
		} else {
			if err = synio.SetVoiceHeadDataByte(args.Param, byte(args.Args[0])); err != nil {
				payload = err.Error()
				return
			}
		}
		payload = "ok"

	case "setVoiceKPROPEle":
		var args struct {
			Index int
			Value int
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		if err = synio.SetVoiceKPROPEle(args.Index, args.Value); err != nil {
			payload = err.Error()
			return
		}
		payload = "ok"

	case "setVoiceVEQEle":
		var args struct {
			Index int
			Value int
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		if err = synio.SetVoiceVEQEle(args.Index, args.Value); err != nil {
			payload = err.Error()
			return
		}
		payload = "ok"

	case "showAbout":
		_ = about_w.Show()

	case "showPreferences":
		logger.Infof("Show Preferences (from messages)\n")
		_ = prefs_w.Show()

	case "getSynergy":
		var response [2]struct {
			HasDevice         bool
			AlreadyConfigured bool
			Name              string
			Choices           *[]zeroconf.Service
		}
		if response[0].HasDevice, response[0].AlreadyConfigured, response[0].Name, response[0].Choices, err = GetSynergyConfig(); err != nil {
			logger.Infof("ZEROCONF: GetSynergyConfig failed: %v\n", err)
			payload = err.Error()
		} else {
			logger.Infof("ZEROCONF: GetSynergyConfig success: %#v\n", response)
			payload = response
		}

	case "getSynergyAndControlSurface":
		var response [2]struct {
			HasDevice         bool
			AlreadyConfigured bool
			Name              string
			Choices           *[]zeroconf.Service
		}
		if response[0].HasDevice, response[0].AlreadyConfigured, response[0].Name, response[0].Choices, err = GetSynergyConfig(); err != nil {
			logger.Infof("ZEROCONF: GetSynergyConfig failed: %v\n", err)
			payload = err.Error()
		} else if response[1].HasDevice, response[1].AlreadyConfigured, response[1].Name, response[1].Choices, err = GetControlSurfaceConfig(); err != nil {
			logger.Infof("ZEROCONF: GetControlSurfaceConfig failed: %v\n", err)
			payload = err.Error()
		} else {
			logger.Infof("ZEROCONF: GetControlSurfaceConfig success: %#v\n", response)
			payload = response
		}

	case "getControlSurface":
		var response [2]struct {
			HasDevice         bool
			AlreadyConfigured bool
			Name              string
			Choices           *[]zeroconf.Service
		}
		if response[0].HasDevice, response[0].AlreadyConfigured, response[0].Name, response[0].Choices, err = GetControlSurfaceConfig(); err != nil {
			logger.Infof("ZEROCONF: GetControlSurfaceConfig failed: %v\n", err)
			payload = err.Error()
		} else {
			logger.Infof("ZEROCONF: GetControlSurfaceConfig success: %#v\n", response)
			payload = response
		}

	case "rescanZeroconf":
		// NOP - basically just waiting a bit for the listener to find new stuff

		// HACK: the javascript modal gets confused if we return too fast (attempting to open a new modal before the
		// previous incarnation has finished transitioning causes the events to be ignored):
		//    https://getbootstrap.com/docs/4.0/components/modal/).
		// So if we returned too fast, add a bit of artificial delay...
		time.Sleep(time.Second * 3)
		payload = "ok"

	case "toggleVoicingMode":
		var args struct {
			Mode            bool
			Disconnect      bool
			Vce             *data.VCE
			ZeroconfSynergy *zeroconf.Service
			ZeroconfCs      *zeroconf.Service
		}
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &args); err != nil {
				payload = err.Error()
				return
			}
		}
		if args.Mode {
			if args.ZeroconfSynergy != nil {
				logger.Infof("ZEROCONF: config Synergy selected by user: %#v\n", *args.ZeroconfSynergy)
				if err = ConnectSynergy(*args.ZeroconfSynergy); err != nil {
					payload = err.Error()
					return
				}
			}
			if args.ZeroconfCs != nil {
				logger.Infof("ZEROCONF: config Control Surface selected by user: %#v\n", *args.ZeroconfCs)
				osc.SetControlSurface((*args.ZeroconfCs).InstanceName, (*args.ZeroconfCs).HostName, (*args.ZeroconfCs).Port)
			}
			payload = nil
			csEnabled := osc.ControlSurfaceConfigured()
			csName := osc.ControlSurfaceName()

			if csEnabled {
				if err = osc.Init(prefsUserPreferences.OscPort, *verboseOscIn, *verboseOscOut, io.SynergyName()); err != nil {
					payload = err.Error()
					return
				}
			}
			if payload == nil {
				var vce data.VCE
				if vce, err = synio.EnableVoicingMode(args.Vce); err != nil {
					payload = err.Error()
					return
				}
				// NOTE: need to pass reference in order to get the custom JSON marshalling to notice the VNAME
				resultPayload := struct {
					Vce         *data.VCE
					CsEnabled   bool
					CsName      string
					SynergyName string
				}{
					Vce:         &vce,
					CsEnabled:   csEnabled,
					CsName:      csName,
					SynergyName: io.SynergyName(),
				}
				payload = resultPayload
			}

		} else {
			if err = osc.Quit(); err != nil {
				payload = err.Error()
			}
			if err = synio.DisableVoicingMode(); err != nil {
				payload = err.Error()
				return
			}
			if args.Disconnect {
				if err = DisconnectSynergy(); err != nil {
					payload = err.Error()
					return
				}
			}
			response := connectionStatusResponse{io.SynergyName(), osc.ControlSurfaceName()}
			payload = response
		}

	case "explore":
		// Unmarshal payload
		var path string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &path); err != nil {
				payload = err.Error()
				return
			}
		}

		// Explore
		if payload, err = explore(path); err != nil {
			payload = err.Error()
			return
		}

	default:
		err = errors.New("Unhandled message " + m.Name)
		payload = err.Error()
		logger.Errorf("%v %v\n", payload, err)
	}
	return
}

// Exploration represents the results of an exploration
type Exploration struct {
	Dirs     []Dir  `json:"dirs"`
	SYNFiles []Dir  `json:"SYNfiles"`
	CRTFiles []Dir  `json:"CRTfiles"`
	VCEFiles []Dir  `json:"VCEfiles"`
	Path     string `json:"path"`
}

// PayloadDir represents a dir payload
type Dir struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// explore explores a path.
// If path is empty, it explores preference's library; if that's empty the user's home directory
func explore(path string) (e Exploration, err error) {
	// If no path is provided, use the preference library path
	if len(path) == 0 {
		path = prefsUserPreferences.LibraryPath
	}
	// if still no path, then use the user's home directory
	if len(path) == 0 {
		var u *user.User
		if u, err = user.Current(); err != nil {
			return
		}
		path = u.HomeDir
	}

	// Read dir
	var files []os.FileInfo
	if files, err = ioutil.ReadDir(path); err != nil {
		return
	}

	// Init exploration
	e = Exploration{
		Dirs:     []Dir{},
		SYNFiles: []Dir{},
		CRTFiles: []Dir{},
		VCEFiles: []Dir{},
		Path:     filepath.Base(path),
	}

	// Add previous dir
	if filepath.Dir(path) != path {
		e.Dirs = append(e.Dirs, Dir{
			Name: "..",
			Path: filepath.Dir(path),
		})
	}

	// Loop through files
	for _, f := range files {
		if f.IsDir() {
			e.Dirs = append(e.Dirs, Dir{
				Name: f.Name(),
				Path: filepath.Join(path, f.Name()),
			})
		} else {

			// Only collect files with Synergy related extensions
			switch strings.ToLower(filepath.Ext(f.Name())) {
			case ".syn":
				e.SYNFiles = append(e.SYNFiles, Dir{
					Name: strings.TrimSuffix(f.Name(), filepath.Ext(f.Name())),
					Path: filepath.Join(path, f.Name()),
				})
			case ".crt":
				e.CRTFiles = append(e.CRTFiles, Dir{
					Name: strings.TrimSuffix(f.Name(), filepath.Ext(f.Name())),
					Path: filepath.Join(path, f.Name()),
				})
			case ".vce":
				e.VCEFiles = append(e.VCEFiles, Dir{
					Name: strings.TrimSuffix(f.Name(), filepath.Ext(f.Name())),
					Path: filepath.Join(path, f.Name()),
				})
			default:
				// ignore
			}
		}
	}

	return
}
