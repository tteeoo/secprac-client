package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/blueberry-jam/secprac-client/util"
)

// Script represents a vulnerability-checking script provided by the server
type Script struct {
	Name     string `json:"name"`
	Points   int    `json:"points"`
	Shell    string `json:"shell"`
	Script   string `json:"script"`
	URL      string `json:"url"`
	Fixed    bool   `json:"fixed"`
	Setup    string `json:"setup"`
	SetupURL string `json:"setup_url"`
}

// GetScripts fetches the vulnerability-checking scripts from the specified remote server\
func GetScripts(remote, token string) ([]Script, error) {

	// Send GET request
	client := &http.Client{}
	url := remote + "/api/vuln/vulns.json"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("token", token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Read response data
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("server responded with bad status code: " + strconv.Itoa(resp.StatusCode) + ", url: " + url + ", body: " + string(body))
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

// DownloadScripts downloads the scripts from the given information and populates the Script.Script struct field of all the scripts.
// Does not download setup scripts if reco is true
func DownloadScripts(remote, token string, scripts []Script, reco bool) ([]Script, error) {
	var c = make(chan error, len(scripts))
	for index := range scripts {
		i := index
		script := &scripts[i]
		go func() {
			url := remote + "/api/scripts/" + script.URL
			util.Logger.Println("downloading script (" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(scripts)) + ")")
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

			// Read response data
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				c <- err
			}
			if resp.StatusCode != 200 {
				c <- errors.New("server responded with bad status code: " + strconv.Itoa(resp.StatusCode) + ", url: " + url + ", body: " + string(body))
			}
			script.Script = string(body)

			// Get setup script (if it exists)
			if script.SetupURL != "" && !reco {
				url := remote + "/api/scripts/setup/" + script.SetupURL
				util.Logger.Println("downloading setup script (" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(scripts)) + ")")
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

				// Read response data
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					c <- err
				}
				if resp.StatusCode != 200 {
					c <- errors.New("server responded with bad status code: " + strconv.Itoa(resp.StatusCode) + ", body: " + string(body))
				}
				script.Setup = string(body)
			}
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

func vuln(remote, endpoint, token string, script Script) (int, error) {

	// Send POST request
	resp, err := http.Post(remote+endpoint, "application/json", strings.NewReader(`{"token": "`+token+`", "name": "`+script.Name+`"}`))
	if err != nil {
		return 0, err
	}

	// Read response data
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != 200 {
		return 0, errors.New("server responded with bad status code: " + strconv.Itoa(resp.StatusCode) + ", body: " + string(body))
	}

	// Parse JSON into map
	var pointJSON map[string]int
	err = json.Unmarshal(body, &pointJSON)
	if err != nil {
		return 0, err
	}
	if _, ok := pointJSON["awarded"]; !ok {
		return 0, errors.New("server response JSON does not contain key \"id\"")
	}

	return pointJSON["awarded"], nil
}

// VulnDone will tell the server that the client is done with a vulnerability, returns the points gained (or lost)
func VulnDone(remote, token string, script Script) (int, error) {
	return vuln(remote, "/api/vuln/done", token, script)
}

// VulnUndo will tell the server to undo a fixed vulnerability, returns the points lost (or gained)
func VulnUndo(remote, token string, script Script) (int, error) {
	return vuln(remote, "/api/vuln/undo", token, script)
}
