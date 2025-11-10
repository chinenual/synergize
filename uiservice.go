package main

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/chinenual/synergize/data"
	"github.com/chinenual/synergize/io"
	"github.com/chinenual/synergize/logger"
	"github.com/chinenual/synergize/osc"
	"github.com/chinenual/synergize/seq"
	"github.com/chinenual/synergize/synio"
	"github.com/chinenual/synergize/zeroconf"
	"github.com/pkg/errors"
)

// UIService: methods called by the frontend/javascript

type UIService struct {
}

func NewUIService() *UIService {
	return new(UIService)
}

type ConnectionStatusResponse struct {
	SynergyName        string
	ControlSurfaceName string
}

func (s *UIService) GetVersion() (version string, newVersionAvailable bool, err error) {
	version = AppVersion
	newVersionAvailable = CheckForNewVersion(false, io.SynergyConnectionType(), osc.ControlSurfaceConfigured())
	err = nil
	return
}

func (s *UIService) CheckVersion(synergyWasDisconnected bool, controlSurfaceWasDisconnected bool) (err error) {
	if synergyWasDisconnected || (controlSurfaceWasDisconnected && prefsUserPreferences.UseOsc) {
		CheckForNewVersion(true, io.SynergyConnectionType(), osc.ControlSurfaceConfigured())
	}
	err = nil
	return
}

func (s *UIService) ShowAbout() (err error) {
	logger.Infof("Show About (from messages)\n")
	wailsAboutWindow.Show()
	return
}
func (s *UIService) ShowPreferences() (err error) {
	logger.Infof("Show Preferences (from messages)\n")
	wailsPrefsWindow.Show()
	return
}

func (s *UIService) CancelPreferences() (err error) {
	wailsPrefsWindow.Hide()
	return
}

func (s *UIService) GetPreferences() (Os string, preferences Preferences, err error) {
	Os = runtime.GOOS
	preferences = prefsUserPreferences
	return
}

func refreshNavPane(path string) {
	wailsApp.Event.Emit("explore", path)
	return
}

func (s *UIService) SavePreferences(preferences Preferences) (err error) {
	logger.Info("INFO: SavePreferences called with ", preferences)
	oldPath := prefsUserPreferences.LibraryPath
	prefsUserPreferences = preferences

	if err = prefsSavePreferences(); err != nil {
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

	wailsPrefsWindow.Hide()

	return
}

type ConnectSynergyResponseType struct {
	AlreadyConnected bool
	Status           ConnectionStatusResponse
}

func (s *UIService) ConnectSynergy(zeroconfChoice *zeroconf.Service) (result ConnectSynergyResponseType, err error) {
	var alreadyConfgured = io.SynergyConfigured()
	if zeroconfChoice != nil {
		logger.Infof("ZEROCONF: config Synergy selected by user: %#v\n", *zeroconfChoice)
		if err = ConnectToSynergy(zeroconfChoice); err != nil {
			return
		}
	} else if !alreadyConfgured {
		err = errors.New("invalid argument to ConnectSynergy")
		return
	}
	result = ConnectSynergyResponseType{
		AlreadyConnected: alreadyConfgured,
		Status: ConnectionStatusResponse{
			SynergyName:        io.SynergyName(),
			ControlSurfaceName: osc.ControlSurfaceName(),
		}}
	return
}

// 	case "connectSynergy":
// 		var args struct {
// 			ZeroconfChoice *zeroconf.Service
// 		}
// 		if len(m.Payload) > 0 {
// 			// Unmarshal payload
// 			if err = json.Unmarshal(m.Payload, &args); err != nil {
// 				payload = err.Error()
// 				return
// 			}
// 		}

// 		var alreadyConfgured = io.SynergyConfigured()
// 		if args.ZeroconfChoice != nil {
// 			logger.Infof("ZEROCONF: config Synergy selected by user: %#v\n", *args.ZeroconfChoice)
// 			if err = ConnectToSynergy(args.ZeroconfChoice); err != nil {
// 				payload = err.Error()
// 				return
// 			}
// 		} else if !alreadyConfgured {
// 			err = errors.New("invalid argument to ConnectSynergy")
// 			payload = err.Error()
// 			return
// 		}
// 		type responseType struct {
// 			AlreadyConnected bool
// 			Status           connectionStatusResponse
// 		}
// 		response := responseType{
// 			AlreadyConnected: alreadyConfgured,
// 			Status: connectionStatusResponse{
// 				SynergyName:        io.SynergyName(),
// 				ControlSurfaceName: osc.ControlSurfaceName(),
// 			}}
// 		payload = response

func (s *UIService) CrtEditAddVoice(crt data.CRT, vcePath string, slot int) (result data.CRT, err error) {
	var vce data.VCE
	if vce, err = data.ReadVceFile(vcePath); err != nil {
		return
	} else {
		logger.Infof("Add vce %s to CRT at slot %d\n", vcePath, slot)
		if len(crt.Voices) < slot {
			// grow the slice
			newVoices := make([]*data.VCE, slot)
			copy(newVoices, crt.Voices)
			crt.Voices = newVoices
		}
		crt.Voices[slot-1] = &vce
		result = crt

		return
	}
}

func (s *UIService) CrtEditLoadCRT(crt data.CRT) (err error) {
	if err = synio.LoadCRT(crt); err != nil {
		return
	}
	return
}

func (s *UIService) CrtEditSaveCRT(path string, crt data.CRT) (err error) {
	if err = data.WriteCrtFileFromVCEArray(path, crt.Voices); err != nil {
		return
	}
	return
}

func (s *UIService) DisableVRAM() (err error) {
	if err = synio.DisableVRAM(); err != nil {
		return
	}
	return
}
func (s *UIService) DisconnectControlSurface() (status ConnectionStatusResponse, err error) {
	if err = DisconnectControlSurface(); err != nil {
		return
	}
	status.SynergyName = io.SynergyName()
	status.ControlSurfaceName = osc.ControlSurfaceName()
	return
}
func (s *UIService) DisconnectSynergy() (status ConnectionStatusResponse, err error) {
	if err = DisconnectSynergy(); err != nil {
		return
	}
	status.SynergyName = io.SynergyName()
	status.ControlSurfaceName = osc.ControlSurfaceName()
	return
}

func (s *UIService) Dx2synCancel() (err error) {
	if err = dx2SynProcessCancel(); err != nil {
		return
	}
	return
}

func (s *UIService) Dx2synStart(path string) (err error) {
	if err = dx2synProcessStart(path); err != nil {
		return
	}
	return
}

func (s *UIService) Syn2midi(path string, tempo float64, raw bool, maxClockSeconds uint32, trackButtons [4]seq.TrackPlayMode) (err error) {
	if err = seq.ConvertSYNToMIDI(path, seq.TrackPerVoice, tempo, raw, maxClockSeconds*1000, trackButtons); err != nil {
		return
	}
	return
}

func (s *UIService) GetSynSequencerState(path string) (trackButtons [4]seq.TrackPlayMode, err error) {
	if trackButtons, err = seq.GetSYNSequencerState(path); err != nil {
		return
	} else {
		logger.Infof("GetSYNSequencerState: %v\n", trackButtons)
	}
	return
}

// 	case "getCWD":
// 		payload, _ = os.Getwd()
// 		logger.Infof("CWD: %s\n", payload)

func (s *UIService) GetConnectionStatus() (status ConnectionStatusResponse, err error) {
	status.SynergyName = io.SynergyName()
	status.ControlSurfaceName = osc.ControlSurfaceName()
	return
}

func (s *UIService) GetPatchTypeNames() (result []string, err error) {
	result = data.PatchTypeNames
	return
}

// 	case "isHTTPDebug":
// 		payload = prefsUserPreferences.HTTPDebug

// 	case "loadCRT":
// 		var path string
// 		if len(m.Payload) > 0 {
// 			// Unmarshal payload
// 			if err = json.Unmarshal(m.Payload, &path); err != nil {
// 				payload = err.Error()
// 				return
// 			}
// 		}
// 		if err = diagLoadCRT(path); err != nil {
// 			payload = err.Error()
// 			return
// 		} else {
// 			payload = "ok"
// 		}

// case "loadSYN":
//
//	var path string
//	if len(m.Payload) > 0 {
//		// Unmarshal payload
//		if err = json.Unmarshal(m.Payload, &path); err != nil {
//			payload = err.Error()
//			return
//		}
//	}
//	if err = diagLoadSYN(path); err != nil {
//		payload = err.Error()
//		return
//	} else {
//		payload = "ok"
//	}

func (s *UIService) loadVceVoicingMode(path string) (vce data.VCE, err error) {
	if vce, err = data.ReadVceFile(path); err != nil {
		return
	}
	if err = synio.LoadVceVoicingMode(vce); err != nil {
		return
	}
	return
}

func (s *UIService) ReadCRT(path string) (crt data.CRT, err error) {
	if crt, err = data.ReadCrtFile(path); err != nil {
		return
	}
	return
}

func (s *UIService) ReadVCE(path string) (vce data.VCE, err error) {
	if vce, err = data.ReadVceFile(path); err != nil {
		return
	}
	return
}

func (s *UIService) runCOMTST() (status string, err error) {
	if err = synio.DiagCOMTST(); err != nil {
		return
	}
	status = "Success!"
	return
}

func (s *UIService) SaveSYN(path string) (err error) {
	if err = diagSaveSYN(path); err != nil {
		return
	}
	return
}

func (s *UIService) SaveVCE(path string) (err error) {
	if err = diagSaveVCE(path); err != nil {
		return
	}
	return
}

func (s *UIService) SetEnvEle(funcname string, osc int, index int, value int) (err error) {
	switch funcname {
	case "setEnvFreqLowVal":
		err = synio.SetEnvFreqLowVal(osc, index, byte(value))
		break
	case "setEnvFreqUpVal":
		err = synio.SetEnvFreqUpVal(osc, index, byte(value))
		break
	case "setEnvAmpLowVal":
		err = synio.SetEnvAmpLowVal(osc, index, byte(value))
		break
	case "setEnvAmpUpVal":
		err = synio.SetEnvAmpUpVal(osc, index, byte(value))
		break
	case "setEnvFreqLowTime":
		err = synio.SetEnvFreqLowTime(osc, index, byte(value))
		break
	case "setEnvFreqUpTime":
		err = synio.SetEnvFreqUpTime(osc, index, byte(value))
		break
	case "setEnvAmpLowTime":
		err = synio.SetEnvAmpLowTime(osc, index, byte(value))
		break
	case "setEnvAmpUpTime":
		err = synio.SetEnvAmpUpTime(osc, index, byte(value))
		break
	default:
		err = errors.New("invalid funcname: " + funcname)
	}
	if err != nil {
		return
	}
	return
}
func (s *UIService) SendToCSurface(field string, value int) (err error) {
	if err = osc.OscSendToCSurface(field, value); err != nil {
		logger.Errorf("Error sending to csurface: %v\n", err)
		return
	}
	return
}

func (s *UIService) SetEnvelopes(osc int, envelopes data.Envelope) (err error) {
	if err = synio.SetEnvelopes(osc, envelopes); err != nil {
		return
	}
	return
}

func (s *UIService) SetFilterArray(uiFilterIndex int, values []int) (err error) {
	if err = synio.SetFilterArray(uiFilterIndex, values); err != nil {
		return
	}
	return
}

func (s *UIService) SetFilterEle(uiFilterIndex int, index int, value int) (err error) {
	if err = synio.SetFilterEle(uiFilterIndex, index, value); err != nil {
		return
	}
	return
}
func (s *UIService) SetLoopPoint(osc int, env string, envType int, sustainPt int, loopPt int) (err error) {
	if err = synio.SetEnvLoopPoint(osc, env, envType, sustainPt, loopPt); err != nil {
		return
	}
	return
}

// 	case "setNumOscillators":
// 		var args struct {
// 			NumOsc    int
// 			PatchType int
// 		}
// 		if len(m.Payload) > 0 {
// 			// Unmarshal payload
// 			if err = json.Unmarshal(m.Payload, &args); err != nil {
// 				payload = err.Error()
// 				return
// 			}
// 		}
// 		var resultPayload struct {
// 			EnvelopeTemplate data.Envelope
// 			PatchBytes       [16]byte
// 		}
// 		if resultPayload.PatchBytes, err = synio.SetNumOscillators(args.NumOsc, args.PatchType); err != nil {
// 			payload = err.Error()
// 			return
// 		}
// 		resultPayload.EnvelopeTemplate = data.DefaultEnvelope
// 		payload = resultPayload

func (s *UIService) SetOscEnvLengths(osc int, freqLength int, ampLength int) (err error) {
	if err = synio.SetOscEnvLengths(osc, freqLength, ampLength); err != nil {
		return
	}
	return
}

// 	case "setOscFILTER":
// 		var args struct {
// 			Args []int
// 		}
// 		if len(m.Payload) > 0 {
// 			// Unmarshal payload
// 			if err = json.Unmarshal(m.Payload, &args); err != nil {
// 				payload = err.Error()
// 				return
// 			}
// 		}
// 		if err = synio.SetOscFILTER(args.Args[0], args.Args[1]); err != nil {
// 			payload = err.Error()
// 			return
// 		}
// 		payload = "ok"

// 	case "setOscKEYPROP":
// 		var args struct {
// 			Args []int
// 		}
// 		if len(m.Payload) > 0 {
// 			// Unmarshal payload
// 			if err = json.Unmarshal(m.Payload, &args); err != nil {
// 				payload = err.Error()
// 				return
// 			}
// 		}
// 		var val = false
// 		if args.Args[1] == 1 {
// 			val = true
// 		}
// 		if err = synio.SetOscKEYPROP(args.Args[0], val); err != nil {
// 			payload = err.Error()
// 			return
// 		}
// 		payload = "ok"

func (s *UIService) SetOscSolo(mute []bool, solo []bool) (oscStatus [16]bool, err error) {
	if oscStatus, err = synio.SetOscSolo(mute, solo); err != nil {
		return
	}
	return
}

// 	case "setOscWAVE":
// 		var args struct {
// 			Args []int
// 		}
// 		if len(m.Payload) > 0 {
// 			// Unmarshal payload
// 			if err = json.Unmarshal(m.Payload, &args); err != nil {
// 				payload = err.Error()
// 				return
// 			}
// 		}
// 		var val = false
// 		if args.Args[1] == 1 {
// 			val = true
// 		}
// 		if err = synio.SetOscWAVE(args.Args[0], val); err != nil {
// 			payload = err.Error()
// 			return
// 		}
// 		payload = "ok"

func (s *UIService) SetVoiceOscDataByte(osc int, value int) (err error) {
	if err = synio.SetVoiceOscDataByte(osc, "OPTCH_reloadGenerators", byte(value)); err != nil {
		return
	}
	return
}

func (s *UIService) SetPatchType(index int) (patchBytes [16]byte, err error) {
	if patchBytes, err = synio.SetPatchType(index); err != nil {
		return
	}
	return
}

// 	case "setVNAME":
// 		var args struct {
// 			Param string
// 			Args  string // HACK: just a string - the JS code shares some logic with the other voice bytes that use setVoiceByte
// 		}
// 		if len(m.Payload) > 0 {
// 			// Unmarshal payload
// 			if err = json.Unmarshal(m.Payload, &args); err != nil {
// 				payload = err.Error()
// 				return
// 			}
// 		}
// 		if err = synio.SetVNAME(args.Args); err != nil {
// 			payload = err.Error()
// 			return
// 		}
// 		payload = "ok"

// 	case "setVoiceByte":
// 		var args struct {
// 			Param string
// 			Args  []int
// 		}
// 		if len(m.Payload) > 0 {
// 			// Unmarshal payload
// 			if err = json.Unmarshal(m.Payload, &args); err != nil {
// 				payload = err.Error()
// 				return
// 			}
// 		}
// 		if len(args.Args) == 2 {
// 			if err = synio.SetVoiceOscDataByte(args.Args[0], args.Param, byte(args.Args[1])); err != nil {
// 				payload = err.Error()
// 				return
// 			}
// 		} else {
// 			if err = synio.SetVoiceHeadDataByte(args.Param, byte(args.Args[0])); err != nil {
// 				payload = err.Error()
// 				return
// 			}
// 		}
// 		payload = "ok"

func (s *UIService) SetVoiceKPROPEle(index int, value int) (err error) {
	if err = synio.SetVoiceKPROPEle(index, value); err != nil {
		return
	}
	return
}

func (s *UIService) SetVoiceVEQEle(index int, value int) (err error) {
	if err = synio.SetVoiceVEQEle(index, value); err != nil {
		return
	}
	return
}

// 	case "getSynergy":
// 		var response [2]struct {
// 			HasDevice         bool
// 			AlreadyConfigured bool
// 			Name              string
// 			Choices           *[]zeroconf.Service
// 		}
// 		if response[0].HasDevice, response[0].AlreadyConfigured, response[0].Name, response[0].Choices, err = GetSynergyConfig(); err != nil {
// 			logger.Infof("ZEROCONF: GetSynergyConfig failed: %v\n", err)
// 			payload = err.Error()
// 		} else {
// 			logger.Infof("ZEROCONF: GetSynergyConfig success: %#v\n", response)
// 			payload = response
// 		}

// 	case "getSynergyAndControlSurface":
// 		var response [2]struct {
// 			HasDevice         bool
// 			AlreadyConfigured bool
// 			Name              string
// 			Choices           *[]zeroconf.Service
// 		}
// 		if response[0].HasDevice, response[0].AlreadyConfigured, response[0].Name, response[0].Choices, err = GetSynergyConfig(); err != nil {
// 			logger.Infof("ZEROCONF: GetSynergyConfig failed: %v\n", err)
// 			payload = err.Error()
// 		} else if response[1].HasDevice, response[1].AlreadyConfigured, response[1].Name, response[1].Choices, err = GetControlSurfaceConfig(); err != nil {
// 			logger.Infof("ZEROCONF: GetControlSurfaceConfig failed: %v\n", err)
// 			payload = err.Error()
// 		} else {
// 			logger.Infof("ZEROCONF: GetControlSurfaceConfig success: %#v\n", response)
// 			payload = response
// 		}

// 	case "getControlSurface":
// 		var response [2]struct {
// 			HasDevice         bool
// 			AlreadyConfigured bool
// 			Name              string
// 			Choices           *[]zeroconf.Service
// 		}
// 		if response[0].HasDevice, response[0].AlreadyConfigured, response[0].Name, response[0].Choices, err = GetControlSurfaceConfig(); err != nil {
// 			logger.Infof("ZEROCONF: GetControlSurfaceConfig failed: %v\n", err)
// 			payload = err.Error()
// 		} else {
// 			logger.Infof("ZEROCONF: GetControlSurfaceConfig success: %#v\n", response)
// 			payload = response
// 		}

// 	case "rescanZeroconf":
// 		// NOP - basically just waiting a bit for the listener to find new stuff

// 		// HACK: the javascript modal gets confused if we return too fast (attempting to open a new modal before the
// 		// previous incarnation has finished transitioning causes the events to be ignored):
// 		//    https://getbootstrap.com/docs/4.0/components/modal/).
// 		// So if we returned too fast, add a bit of artificial delay...
// 		time.Sleep(time.Second * 3)
// 		payload = "ok"

func (s *UIService) ToggleVoicingMode(
	mode bool,
	disconnect bool,
	useVce *data.VCE,
	zeroconfSynergy *zeroconf.Service,
	zeroconfCs *zeroconf.Service) (loaded_vce *data.VCE,
	csEnabled bool,
	csName string,
	synergyName string,
	err error) {

	// 	case "toggleVoicingMode":
	// 		var args struct {
	// 			Mode            bool
	// 			Disconnect      bool
	// 			Vce             *data.VCE
	// 			ZeroconfSynergy *zeroconf.Service
	// 			ZeroconfCs      *zeroconf.Service
	// 		}

	if mode {
		if zeroconfSynergy != nil {
			logger.Infof("ZEROCONF: config Synergy selected by user: %#v\n", *zeroconfSynergy)
			if err = ConnectSynergy(*zeroconfSynergy); err != nil {
				return
			}
		}
		if zeroconfCs != nil {
			logger.Infof("ZEROCONF: config Control Surface selected by user: %#v\n", *zeroconfCs)
			osc.SetControlSurface((*zeroconfCs).InstanceName, (*zeroconfCs).HostName, (*zeroconfCs).Port)
		}
		csEnabled = osc.ControlSurfaceConfigured()
		csName = osc.ControlSurfaceName()

		if csEnabled {
			if err = osc.Init(prefsUserPreferences.OscPort, *verboseOscIn, *verboseOscOut, io.SynergyName()); err != nil {
				return
			}
		}
		var vce data.VCE
		if vce, err = synio.EnableVoicingMode(useVce); err != nil {
			return
		}
		loaded_vce = &vce
		synergyName = io.SynergyName()

	} else {
		if err = osc.Quit(); err != nil {
			return
		}
		if err = synio.DisableVoicingMode(); err != nil {
			return
		}
		if disconnect {
			if err = DisconnectSynergy(); err != nil {
				return
			}
		}
	}
	return
}

func (s *UIService) Explore(path string) (exploration Exploration, err error) {
	// 	case "explore":
	// 		// Unmarshal payload
	// 		var path string
	// 		if len(m.Payload) > 0 {
	// 			// Unmarshal payload
	// 			if err = json.Unmarshal(m.Payload, &path); err != nil {
	// 				payload = err.Error()
	// 				return
	// 			}
	// 		}

	// 		// Explore
	// 		if payload, err = explore(path); err != nil {
	// 			payload = err.Error()
	// 			return
	// 		}
	exploration, err = explore(path)
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
