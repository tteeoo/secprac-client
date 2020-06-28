package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// Team represents a participating team authenticated on the server
type Team struct {
	Token string `json:"token"`
	ID    string `json:"id"`
}

// NewTeam interacts with the API to create a new team at the given remote server with the specified token
func NewTeam(remote, token string) (*Team, error) {

	// Send POST request
	resp, err := http.Post(remote+"/api/team/create", "application/json", strings.NewReader("{\"token\": \""+token+"\"}"))
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
