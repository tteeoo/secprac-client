package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/blueberry-jam/secprac-client/util"
)

// Script represents a vulnerability-checking script provided by the server
type Script struct {
	Name   string `json:"name"`
	Points int    `json:"points"`
	Shell  string `json:"shell"`
	Script string `json:"script"`
	URL    string `json:"url"`
}

// GetScripts fetches the vulnerability-checking scripts from the specified remote server\
func GetScripts(remote, token string) ([]Script, error) {

	// Send GET request
	client := &http.Client{}
	req, err := http.NewRequest("GET", remote+"/api/vuln/vulns.json", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("token", token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("server responded with bad status code: " + strconv.Itoa(resp.StatusCode))
	}

	// Read response data
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse JSON into slice
	var scripts []Script
	var scriptMap map[string]Script
	err = json.Unmarshal(body, &scriptMap)
	if err != nil {
		return nil, err
	}

	for k := range scriptMap {
		scripts = append(scripts, scriptMap[k])
	}
	return scripts, nil
}

// DownloadScripts downloads the scripts from the given information and populates the Script.Script struct field of all the scripts
func DownloadScripts(remote, token string, scripts []Script) ([]Script, error) {
	var c = make(chan error, len(scripts))
	for i := range scripts {
		script := &scripts[i]
		go func() {
			url := remote + "/api/scripts/" + script.URL
			util.Logger.Println("downloading script:", url)
			client := &http.Client{}
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				c <- err
			}
			req.Header.Set("token", token)
			resp, err := client.Do(req)
			if err != nil {
				c <- err
			}
			if resp.StatusCode != 200 {
				c <- errors.New("server responded with bad status code: " + strconv.Itoa(resp.StatusCode))
			}

			// Read response data
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				c <- err
			}
			script.Script = string(body)
			c <- nil
		}()
	}
	for range scripts {
		if err := <-c; err != nil {
			return []Script{}, err
		}
	}
	return scripts, nil
}
