package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/chinenual/synergize/data"
	"github.com/chinenual/synergize/osc"
	"github.com/chinenual/synergize/synio"
	"github.com/chinenual/synergize/zeroconf"

	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"github.com/pkg/errors"
)

var chooseZeroconfServiceChan chan int

func chooseZeroconfService(prompt string, choices []zeroconf.Service) (choice *zeroconf.Service, err error) {
	log.Println("m")
	choice = nil
	var stringified []string
	log.Println("l")
	for _, s := range choices {
		log.Println("k")
		stringified = append(stringified, fmt.Sprintf("%s (%s:%d)", s.InstanceName(), s.Address(), s.Port()))
	}
	log.Println("j")
	var msg = struct {
		Prompt  string
		Choices []string
	}{prompt, stringified}
	log.Println("i")

	var wg sync.WaitGroup
	wg.Add(1)
	// we want to treat this as a synchronous call - need to play games to wait on the result:
	if err = bootstrap.SendMessage(w, "chooseZeroconfService", msg,
		func(m *bootstrap.MessageIn) {
			defer wg.Done()
			log.Println("h")
			// Unmarshal payload
			var strval string
			log.Println("g")
			if err = json.Unmarshal(m.Payload, &strval); err != nil {
				log.Println("f")
				log.Printf("chooseZeroconfService failed to decode json response : %v\n", err)
				return
			}
			log.Println("")
			if strval != "ok" {
				log.Println("e")
				err = errors.Errorf("chooseZeroconfService failed to return ok: : %s", strval)
				return
			} else {
				log.Println("d")
				log.Printf(".... waiting for callback\n")
				intval := <-chooseZeroconfServiceChan
				log.Println("c")
				log.Printf(".... callback returned %d\n", intval)

				if intval >= 0 && intval < len(choices) {
					choice = &choices[intval]
				}
			}
			log.Println("b")
		}); err != nil {
		wg.Done()
	}
	wg.Wait()
	log.Println("a")
	return
}

func getZeroconfAddress(serviceType string, choices *[]zeroconf.Service) (addr string, port uint, err error) {
	if len(*choices) == 0 {
		// we will try one more time to browse
		log.Printf("ZEROCONF: No %s services found on previous browse - trying again\n", serviceType)
		zeroconf.Browse()
	}
	var choice *zeroconf.Service
	if len(*choices) == 0 {
		err = errors.Errorf("Cannot find %s via Bonjour/zeroconf", serviceType)
		return
	} else if len(*choices) == 1 {
		choice = &(*choices)[0]
	} else {
		// user needs to choose
		if choice, err = chooseZeroconfService("Choose "+serviceType, *choices); err != nil {
			return
		}
		if choice == nil {
			log.Printf("ZEROCONF: user chose 'none of the above'\n")
		}
	}
	addr = ""
	port = 0
	if choice != nil {
		addr = choice.Address()
		port = uint(choice.Port())

		log.Printf("ZEROCONF: auto configuring %s: %s:%d [%s: %s]\n",
			serviceType,
			addr, port, choice.HostName(), choice.InstanceName())
	}
	return
}

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	log.Printf("Handle message: %s %s\n", m.Name, string(m.Payload))
	switch m.Name {
	case "cancelPreferences":
		prefs_w.Hide()

	case "chooseZeroconfServiceCallback":
		var selected int
		// Unmarshal payload
		if err = json.Unmarshal(m.Payload, &selected); err != nil {
			payload = err.Error()
			return
		}
		log.Printf("...got callback (%d) send to channel\n", selected)
		chooseZeroconfServiceChan <- selected

	case "connectToSynergy":
		if err = connectToSynergy(); err != nil {
			payload = err.Error()
		} else {
			payload = FirmwareVersion
		}

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
			log.Printf("Add vce %s to CRT at slot %d\n", args.VcePath, args.Slot)
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
		if err = connectToSynergyIfNotConnected(); err != nil {
			payload = err.Error()
			return
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
		if err = connectToSynergyIfNotConnected(); err != nil {
			payload = err.Error()
			return
		}
		if err = synio.DisableVRAM(); err != nil {
			payload = err.Error()
			return
		} else {
			payload = "ok"
		}

	case "getCWD":
		payload, _ = os.Getwd()
		log.Printf("CWD: %s\n", payload)

	case "getFirmwareVersion":
		payload = FirmwareVersion

	case "getPreferences":
		payload = struct {
			Os          string
			Preferences Preferences
		}{runtime.GOOS, prefsUserPreferences}

	case "getVersion":
		payload = struct {
			Version             string
			NewVersionAvailable bool
		}{AppVersion, CheckForNewVersion(false, false)}

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
		if diagLoadCRT(path); err != nil {
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
		if FirmwareVersion == "" {
			// not yet connected to the Synergy.
			// A conundrum: if user has already put the synergy into
			// test mode, we can't query the firmware.  If we just
			// initialize the serial connection and dont update
			// firmware version, the UI will continue to show
			// "not connected".
			//
			// Run the serial init without querying the firmware version
			if err = synio.Init(prefsUserPreferences.SerialPort, prefsUserPreferences.SerialBaud, true, *serialVerboseFlag, *mockSynio); err != nil {
				err = errors.Wrapf(err, "Cannot connect to synergy on port %s at %d baud\n", prefsUserPreferences.SerialPort, prefsUserPreferences.SerialBaud)
				payload = err.Error()
				return
			}
			FirmwareVersion = "Connected"
		}
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
				log.Printf("ERROR: could not start zeroconf: %v\n", err)
			}
		} else {
			zeroconf.CloseServer()
		}

		prefs_w.Hide()
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
		osc.OscSendToCSurface(args.Field, args.Value)
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
		about_w.Show()

	case "showPreferences":
		log.Printf("Show Preferences (from messages)\n")
		prefs_w.Show()

	case "toggleVoicingMode":
		if err = connectToSynergyIfNotConnected(); err != nil {
			payload = err.Error()
			return
		}
		var mode bool
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &mode); err != nil {
				payload = err.Error()
				return
			}
		}
		if mode {
			var vce data.VCE
			payload = nil
			var csEnabled = false
			if prefsUserPreferences.UseOsc {
				port := prefsUserPreferences.OscPort
				csurfaceAddress := prefsUserPreferences.OscCSurfaceAddress
				csurfacePort := prefsUserPreferences.OscCSurfacePort
				log.Println("FFF")
				if prefsUserPreferences.OscAutoConfig {
					log.Println("EEE")
					if csurfaceAddress, csurfacePort, err = getZeroconfAddress("Control Surface", &zeroconf.OscServices); err != nil {
						log.Println("DDD")
						log.Println(err)
						payload = err.Error()
						return
					}
					log.Println("CCC")
				}
				log.Println("BBBB")
				if csurfaceAddress != "" {
					log.Println("AAAA")
					if err = osc.OscInit(port,
						csurfaceAddress,
						csurfacePort,
						*verboseOscIn, *verboseOscOut,
						prefsSynergyName()); err != nil {

						log.Println(err)
						payload = err.Error()
					} else {
						csEnabled = true
					}
				}
			}
			log.Println("4444")
			if payload == nil {
				log.Println("3333")
				if vce, err = synio.EnableVoicingMode(); err != nil {
					payload = err.Error()
					return
				}
				log.Println("2222")
				// NOTE: need to pass reference in order to get the custom JSON marshalling to notice the VNAME
				resultPayload := struct {
					Vce       *data.VCE
					CsEnabled bool
				}{
					Vce:       &vce,
					CsEnabled: csEnabled,
				}
				log.Println("1111")
				payload = resultPayload
				log.Println("0000")
			}

		} else {
			if err = osc.OscQuit(); err != nil {
				payload = err.Error()
			}
			if err = synio.DisableVoicingMode(); err != nil {
				payload = err.Error()
				return
			}
			payload = "ok"
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
		log.Printf("ERROR: %v %v\n", payload, err)
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
