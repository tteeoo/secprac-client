package api

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

// GetReport will download the team's scoring report and save it
func GetReport(remote, id, token string) error {

	// Send GET request
	client := &http.Client{}
	req, err := http.NewRequest("GET", remote+"/api/report", nil)
	req.Header.Set("token", token)
	req.Header.Set("name", id)
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	// Read response data
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("server responded with bad status code when getting report: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
	}

	// Write to file
	err = ioutil.WriteFile("/usr/local/share/secprac/report.html", body, 0644)
	if err != nil {
		return errors.New("error writing report to file:" + err.Error())
	}

	return nil
}
