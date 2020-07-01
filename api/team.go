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

// Team represents a participating team authenticated on the server
type Team struct {
	Token   string   `json:"token"`
	ID      string   `json:"id"`
	Scripts []Script `json:"scripts"`
	Points  int      `json:"points"`
}

// NewTeam interacts with the API to create a new team at the given remote server with the specified token
func NewTeam(remote, token string) (*Team, error) {

	// Send POST request
	resp, err := http.Post(remote+"/api/team/create", "application/json", strings.NewReader("{\"token\": \""+token+"\",\"time\": \""+util.GetTimestamp()+"\"}"))
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
		return nil, errors.New("server responded with bad status code: " + strconv.Itoa(resp.StatusCode) + ", body: " + string(body))
	}

	// Parse JSON into map
	var idJSON map[string]string
	err = json.Unmarshal(body, &idJSON)
	if err != nil {
		return nil, err
	}
	if _, ok := idJSON["id"]; !ok {
		return nil, errors.New("server response JSON does not contain key \"id\"")
	}

	return &Team{
		Token: token,
		ID:    idJSON["id"],
	}, nil
}

// TeamDone will tell the server the client thinks it is done, and the server will respond with 200 if it agrees
func TeamDone(remote, token string) error {

	// Send POST request
	resp, err := http.Post(remote+"/api/team/done", "application/json", strings.NewReader("{\"token\": \""+token+"\"}"))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("server responded with bad status code: " + strconv.Itoa(resp.StatusCode))
	}

	return nil
}
