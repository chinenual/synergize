package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/chinenual/synergize/logger"
)

type Preferences struct {
	SerialPort         string
	SerialBaud         uint
	LibraryPath        string
	UseOsc             bool
	OscAutoConfig      bool
	OscPort            uint
	OscCSurfaceAddress string
	OscCSurfacePort    uint
	VstAutoConfig      bool

	// hidden from user
	HTTPDebug      bool
	VstServiceType string
}

var preferencesPathname = getWorkingDirectory() + "/preferences.json"

var prefsUserPreferences = Preferences{
	UseOsc:          false,
	SerialBaud:      9600,
	OscAutoConfig:   false,
	OscPort:         8000,
	OscCSurfacePort: 9000,
	VstAutoConfig:   false,

	VstServiceType: "_synergia._tcp",
}

func prefsLoadPreferences() (err error) {
	var b []byte
	if b, err = ioutil.ReadFile(preferencesPathname); err != nil {
		logger.Errorf("Error loading preferences.  Using defaults %#v: %v", prefsUserPreferences, err)
		return
	}
	if err = json.Unmarshal(b, &prefsUserPreferences); err != nil {
		logger.Errorf("Error parsing preferences.  Using defaults %#v: %v", prefsUserPreferences, err)
		return
	}
	logger.Infof("Loaded preferences %#v from file %s\n", prefsUserPreferences, preferencesPathname)
	return
}

func prefsSavePreferences() (err error) {
	var b []byte
	if b, err = json.MarshalIndent(prefsUserPreferences, "", " "); err != nil {
		logger.Error("Error saving preferences", err)
	}
	logger.Infof("Save preferences %#v to file %s\n", prefsUserPreferences, preferencesPathname)
	if err = ioutil.WriteFile(preferencesPathname, b, 0644); err != nil {
		logger.Error("Error saving preferences", err)
	}
	return
}

func prefsSynergyName() string {
	hostname, _ := os.Hostname()
	return fmt.Sprintf("%s [%s]", hostname, prefsUserPreferences.SerialPort)
}
