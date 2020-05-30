package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/chinenual/synergize/data"
	"github.com/chinenual/synergize/synio"

	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"github.com/pkg/errors"
)

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "connectToSynergy":
		if err = connectToSynergy(); err != nil {
			payload = err.Error()
		} else {
			payload = FirmwareVersion
		}

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
			if vce, err = synio.EnableVoicingMode(); err != nil {
				payload = err.Error()
				return
			}
			payload = vce

		} else {
			if err = synio.DisableVoicingMode(); err != nil {
				payload = err.Error()
				return
			}
			payload = "ok"
		}

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
		var oscStatus [16]bool
		if oscStatus, err = synio.SetOscSolo(args.Mute, args.Solo); err != nil {
			payload = err.Error()
			return
		}
		payload = oscStatus

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

	case "getVersion":
		payload = struct {
			Version             string
			NewVersionAvailable bool
		}{AppVersion, CheckForNewVersion()}

	case "showAbout":
		about_w.Show()

	case "showPreferences":
		log.Printf("Show Preferences (from messages)\n")
		prefs_w.Show()

	case "getPreferences":
		payload = struct {
			Os          string
			Preferences Preferences
		}{runtime.GOOS, prefsUserPreferences}

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
		prefs_w.Hide()

	case "cancelPreferences":
		prefs_w.Hide()

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
			payload = vce
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
			payload = vce
		}

	case "getFirmwareVersion":
		payload = FirmwareVersion

	case "getCWD":
		payload, _ = os.Getwd()
		log.Printf("CWD: %s\n", payload)

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
			if err = synio.Init(prefsUserPreferences.SerialPort, prefsUserPreferences.SerialBaud, true, *serialVerboseFlag); err != nil {
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
		payload = errors.New("Unhandled message " + m.Name)
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
