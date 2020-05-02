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
}

var preferencesPathname string = "preferences.json"

var prefsUserPreferences Preferences

func prefsLoadPreferences() (err error) {
	var b []byte
	b,err = ioutil.ReadFile(preferencesPathname)
	if err != nil {
		log.Println("Error loading preferences", err)
		return
	}
	err = json.Unmarshal(b, &prefsUserPreferences)
	if err != nil {
		log.Println("Error parsing preferences", err)
		return
	}
	log.Printf("Loaded preferences %#v from file %s\n", prefsUserPreferences, preferencesPathname)
	return
}

func prefsSavePreferences() (err error) {
	var b []byte
	b,err = json.MarshalIndent(prefsUserPreferences, "", " ")
	if err != nil {
		log.Println("Error saving preferences", err)
	}
	log.Printf("Save preferences %#v to file %s\n", prefsUserPreferences, preferencesPathname)
	err = ioutil.WriteFile(preferencesPathname, b, 0644)
	if err != nil {
		log.Println("Error saving preferences", err)
	}
	return
}
