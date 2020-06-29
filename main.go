package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/blueberry-jam/secprac-client/api"
	"github.com/blueberry-jam/secprac-client/util"
)

// Entry point
func main() {

	// Get remote server from command line args
	if len(os.Args) < 2 {
		util.Notify("error", "no remote server was provided as a command-line argument", util.IconMinus, true)
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
	util.Logger.Println("generated team token")

	// Authenticate with the server
	util.Logger.Println("attempting to authenticate with server (" + remote + ")")
	team, err := api.NewTeam(remote, token)
	if err != nil {
		util.Notify("error", "failed to authenticate with the server, check the log at: "+util.LogFileName, util.IconMinus, true)
		util.Logger.Fatalln("error authenticating with server:", err)
	}

	util.Notify("authenticated", "successfully authenticated with server, your team ID is "+team.ID, util.IconInfo, false)
	util.Logger.Println("authenticated with", remote, "given ID", team.ID)

	// Get the vulnerability-checking scripts
	team.Scripts, err = api.GetScripts(remote, team.Token)
	if err != nil {
		util.Notify("error", "failed to get the script information from the server, check the log at: "+util.LogFileName, util.IconMinus, true)
		util.Logger.Fatalln("error getting script information from the server:", err)
	}
	if len(team.Scripts) < 1 {
		util.Notify("error", "the server did not provide any scripts... you win?", util.IconPlus, true)
		util.Logger.Fatalln("server provided no scripts")
	}

	// Download scripts
	team.Scripts, err = api.DownloadScripts(remote, team.Token, team.Scripts)
	if err != nil {
		util.Notify("error", "failed to download scripts from the server, check the log at: "+util.LogFileName, util.IconMinus, true)
		util.Logger.Fatalln("error downloading a script from the server:", err)
	}
	util.Notify("downloaded scripts", "successfully downloaded "+strconv.Itoa(len(team.Scripts))+" scripts, start hacking!", util.IconInfo, false)

	// Main loop
	for {
		for i := range team.Scripts {
			script := &team.Scripts[i]

			// Pipe script into shell and run
			cmd := exec.Command(script.Shell)
			stdin, err := cmd.StdinPipe()
			if err != nil {
				util.Logger.Println("error getting command stdin:", err)
				continue
			}
			go func() {
				defer stdin.Close()
				io.WriteString(stdin, script.Script)
			}()
			out, err := cmd.CombinedOutput()
			if err != nil {
				util.Logger.Println("error running script:", err)
				continue
			}

			// Check if fixed
			if script.Fixed {
				if string(out) != "FIXED\n" {
					util.Logger.Println("script undone:", script.Name)
					// TODO: send request
					script.Fixed = false
				}
			} else {
				if string(out) == "FIXED\n" {
					util.Logger.Println("script fixed:", script.Name)
					// TODO: send request
					script.Fixed = true
				}
			}
			time.Sleep(time.Second / 5)
		}
	}
}
