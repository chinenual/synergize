package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

var newVersionAvailable bool = false
var checkedForVersion bool = false

var url string = "https://api.github.com/repos/chinenual/synergize/releases/latest"

type GithubReleaseApiResponse struct {
	// we only care about one field:
	TagName string `json:"tag_name"`
}

func CheckForNewVersion() (newVersion bool) {
	if !checkedForVersion {
		checkedForVersion = true
		if response, err := http.Get(url); err != nil {
			log.Printf("Check for new version API call failed with %v\n", err)
			newVersionAvailable = false

		} else {
			if bytes, err := ioutil.ReadAll(response.Body); err != nil {
				log.Printf("could not read response stream %v\n", err)
			} else {
				var res GithubReleaseApiResponse
				if err := json.Unmarshal(bytes, &res); err != nil {
					log.Printf("could not decode response %v\n", err)
				} else {
					log.Printf("Latest version is %s\n", res.TagName)
					if "v"+Version != res.TagName {
						newVersionAvailable = true
					}

				}
			}
		}
	}
	return newVersionAvailable
}
