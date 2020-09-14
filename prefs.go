package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Preferences struct {
	SerialPort         string
	SerialBaud         uint
	LibraryPath        string
	HTTPDebug          bool
	UseOsc             bool
	OscAutoConfig      bool
	OscPort            uint
	OscCSurfaceAddress string
	OscCSurfacePort    uint
	VstAutoConfig      bool
	VstAddress         string
	VstPort            uint
}

var preferencesPathname = getWorkingDirectory() + "/preferences.json"

var prefsUserPreferences = Preferences{
	UseOsc:          false,
	SerialBaud:      9600,
	OscAutoConfig:   true,
	OscPort:         8000,
	OscCSurfacePort: 9000,
	VstAutoConfig:   true,
}

func prefsLoadPreferences() (err error) {
	var b []byte
	if b, err = ioutil.ReadFile(preferencesPathname); err != nil {
		log.Println("Error loading preferences", err)
		return
	}
	if err = json.Unmarshal(b, &prefsUserPreferences); err != nil {
		log.Println("Error parsing preferences", err)
		return
	}
	log.Printf("Loaded preferences %#v from file %s\n", prefsUserPreferences, preferencesPathname)
	return
}

func prefsSavePreferences() (err error) {
	var b []byte
	if b, err = json.MarshalIndent(prefsUserPreferences, "", " "); err != nil {
		log.Println("Error saving preferences", err)
	}
	log.Printf("Save preferences %#v to file %s\n", prefsUserPreferences, preferencesPathname)
	if err = ioutil.WriteFile(preferencesPathname, b, 0644); err != nil {
		log.Println("Error saving preferences", err)
	}
	return
}

func prefsSynergyName() string {
	hostname, _ := os.Hostname()
	return fmt.Sprintf("%s [%s]", hostname, prefsUserPreferences.SerialPort)
}
