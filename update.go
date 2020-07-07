package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"
)

var newVersionAvailable bool = false
var checkedForVersion bool = false

var versionUrl string = "https://api.chinenual.com/api/v1/version?latest&app=Synergize&v="+Version+"&OS=" + runtime.GOOS + "&ARCH=" + runtime.GOARCH

type GithubReleaseApiResponse []struct {
	Redirect string `json:"redirect"`
	TagName string `json:"tag_name"`
}

func getLatest(url string) (version string,err error) {
	var response *http.Response
	if response, err = http.Get(url); err != nil {
		log.Printf("Check for new version API call failed with %v\n", err)
		return
	} else {
		var bytes []byte
		if bytes, err = ioutil.ReadAll(response.Body); err != nil {
			log.Printf("could not read response stream %v\n", err)
			return
		} else {
			var res GithubReleaseApiResponse
			if err = json.Unmarshal(bytes, &res); err != nil {
				log.Printf("could not decode response %v\n", err)
				return
			} else if res[0].Redirect != "" {
				version,err = getLatest(res[0].Redirect)
				return
			} else {
				version = res[0].TagName;
				return
			}
		}
	}
	return
}

func CheckForNewVersion(forceRecheck bool, connected bool) (newVersion bool) {
	if forceRecheck || (!checkedForVersion) {
		checkedForVersion = true
		var latestVersion string
		var err error
		var url = versionUrl
		if forceRecheck {
			url = versionUrl + "&connected=" + strconv.FormatBool(connected)
		}
		if latestVersion,err = getLatest(url); err != nil {
			log.Printf("Error checking for new version: %v",err)
			return
		}			
		log.Printf("Latest version is %s\n", latestVersion)
		if "v"+Version != latestVersion {
			newVersionAvailable = true
		}
	}
	newVersion = newVersionAvailable;
	return
}
