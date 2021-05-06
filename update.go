package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"

	"github.com/chinenual/synergize/logger"
)

var newVersionAvailable bool = false
var checkedForVersion bool = false

var versionUrl string = "https://api.chinenual.com/api/v1/version?latest&app=Synergize&v=" + Version + "&OS=" + runtime.GOOS + "&ARCH=" + runtime.GOARCH + "&OSVER=" // version filled in in the function call - not yet initialized at variable init time

type GithubReleaseApiResponse []struct {
	Redirect string `json:"redirect"`
	TagName  string `json:"tag_name"`
}

func getLatest(url string) (version string, err error) {
	var response *http.Response
	if response, err = http.Get(url); err != nil {
		logger.Infof("Check for new version API call failed with %v\n", err)
		return
	} else {
		var bytes []byte
		if bytes, err = ioutil.ReadAll(response.Body); err != nil {
			logger.Infof("could not read response stream %v\n", err)
			return
		} else {
			var res GithubReleaseApiResponse
			if err = json.Unmarshal(bytes, &res); err != nil {
				logger.Infof("could not decode response %v\n", err)
				return
			} else if res[0].Redirect != "" {
				version, err = getLatest(res[0].Redirect)
			} else {
				version = res[0].TagName
			}
		}
	}
	return
}

func CheckForNewVersion(forceRecheck bool, synergyType string, hasCs bool, other string) (newVersion bool) {
	if forceRecheck || (!checkedForVersion) {
		checkedForVersion = true
		var latestVersion string
		var err error
		var url = versionUrl + OsVersion
		if forceRecheck {
			url = versionUrl + "&synergy=" + synergyType + "&cs=" + strconv.FormatBool(hasCs)
			if other != "" {
				url += "&other=" + other
			}
		}
		logger.Debugf("url: %s\n", url)
		if latestVersion, err = getLatest(url); err != nil {
			logger.Errorf("Error checking for new version: %v", err)
			return
		}
		if !forceRecheck {
			logger.Infof("Latest version is %s\n", latestVersion)
		}
		if "v"+Version != latestVersion {
			newVersionAvailable = true
		}
	}
	newVersion = newVersionAvailable
	return
}
