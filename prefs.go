package main

import (
	"encoding/json"
	"log"
	"io/ioutil"
)
	
type Preferences struct {
	SerialPort string
	SerialBaud uint
	LibraryPath string
	HTTPDebug bool
}

var preferencesPathname string = getWorkingDirectory() + "/preferences.json"

var prefsUserPreferences Preferences

func prefsLoadPreferences() (err error) {
	var b []byte
	if b,err = ioutil.ReadFile(preferencesPathname); err != nil {
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
	if	b,err = json.MarshalIndent(prefsUserPreferences, "", " "); err != nil {
		log.Println("Error saving preferences", err)
	}
	log.Printf("Save preferences %#v to file %s\n", prefsUserPreferences, preferencesPathname)
	if err = ioutil.WriteFile(preferencesPathname, b, 0644); err != nil {
		log.Println("Error saving preferences", err)
	}
	return
}

