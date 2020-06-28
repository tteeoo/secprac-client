package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/blueberry-jam/secprac-client/util"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	// "github.com/blueberry-jam/secprac-client/script"
)

func main() {

	// Get remote server from command line args
	if len(os.Args) < 2 {
		util.Logger.Fatalln("no server provided; run again like this: `secprac-client <server ip/url>`")
	}
	remote := os.Args[1]

	// Generate random Token
	b := make([]byte, 18)
	_, err := rand.Read(b)
	if err != nil {
		util.Logger.Fatalln(err)
	}
	token := base64.URLEncoding.EncodeToString(b)
	util.Logger.Println("generated team token:", token)

	// Authenticate with the server
	// Send post req
	util.Logger.Println("attempting to authenticate with server (" + remote + ")")
	resp, err := http.Post(remote+"/api/team/create", "application/json", strings.NewReader("{\"Token\": \""+token+"\"}"))
	if err != nil {
		util.Notify("error", "failed to authenticate with the server, check the log at: "+util.LogFileName, "", true)
		util.Logger.Fatalln("error authenticating with server:", err)
	}
	if resp.StatusCode != 200 {
		util.Notify("error", "failed to authenticate with the server, check the log at: "+util.LogFileName, "", true)
		util.Logger.Fatalln("error server responded with HTTP code", resp.Status)
	}

	// Read response data
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Notify("error", "failed to authenticate with the server, check the log at: "+util.LogFileName, "", true)
		util.Logger.Fatalln("error parsing server response body JSON:", err)
	}

	// Parse into JSON
	var idJSON map[string]string
	err = json.Unmarshal(body, &idJSON)
	if _, ok := idJSON["ID"]; !ok {
		util.Notify("error", "failed to authenticate with the server, check the log at: "+util.LogFileName, "", true)
		util.Logger.Fatalln("server response body JSON did not contain key 'ID'")
	}

	teamID := idJSON["ID"]
	util.Notify("authenticated", "successfully authenticated with server, your team ID is: "+teamID, "", false)
	util.Logger.Println("connected to", remote, "with ID", teamID)
}
