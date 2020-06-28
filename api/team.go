package api

import (
	"net/http"
	"strings"
)

// CreateTeam interacts with the API to create a new team at the given remote server with the specified token
func CreateTeam(remote, token string) (*http.Response, error) {
	resp, err := http.Post(remote+"/api/team/create", "application/json", strings.NewReader("{\"token\": \""+token+"\"}"))
	return resp, err
}
