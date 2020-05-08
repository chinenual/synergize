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
	
	"github.com/pkg/errors"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
)

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "connectToSynergy":
		err = connectToSynergy()
		if err != nil {
			payload = err.Error()
		} else {
			payload = FirmwareVersion
		}
		
	case "disableVRAM":
		err = connectToSynergyIfNotConnected();
		if err != nil {
			payload = err.Error()
			return
		} 
		err = synioDisableVRAM()
		if err != nil {
			payload = err.Error()
			return
		} else {
			payload = "ok"
		}
		
	case "getVersion":
		payload = AppVersion

	case "showAbout":
		about_w.Show()
		
	case "showPreferences":
		log.Printf("Show Preferences (from messages)\n");
		prefs_w.Show()
		
	case "getPreferences":
		payload = struct {
			Os string
			Preferences Preferences
		} {runtime.GOOS, prefsUserPreferences}
		
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
		err = diagLoadSYN(path)
		if err != nil {
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
		err = diagSaveSYN(path)
		if err != nil {
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
		err = diagLoadCRT(path)
		if err != nil {
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
		crt,err = data.CrtReadFile(path);
		if err != nil {
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
		vce,err = data.VceReadFile(path);
		if err != nil {
			payload = err.Error()
			return
		} else {
			payload = vce
		}
		
	case "getFirmwareVersion":
		payload = FirmwareVersion

	case "getCWD":
		payload,_ = os.Getwd()
		log.Printf("CWD: %s\n",payload)
		
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
			err = synioInit(prefsUserPreferences.SerialPort, prefsUserPreferences.SerialBaud, true, *serialVerboseFlag)
			if err != nil {
				err = errors.Wrapf(err, "Cannot connect to synergy on port %s at %d baud\n", prefsUserPreferences.SerialPort,prefsUserPreferences.SerialBaud)
				payload = err.Error()
				return
			}
			FirmwareVersion = "Connected"
		}
		err = synioDiagCOMTST()
		if err != nil {
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
	Dirs       []Dir              `json:"dirs"`
	SYNFiles   []Dir              `json:"SYNfiles"`
	CRTFiles   []Dir              `json:"CRTfiles"`
	VCEFiles   []Dir              `json:"VCEfiles"`
	Path       string             `json:"path"`
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
			case ".vce" :
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
