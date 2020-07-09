package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	ou "os/user"
	"strconv"
	"strings"
	"time"

	"github.com/blueberry-jam/secprac-client/api"
	"github.com/blueberry-jam/secprac-client/util"
)

// Entry point
func main() {

	// Print version
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-V") {
		println("secprac-client verison 0.1.4\nCopyright (C) Theo Henson\nMIT License\nOpen source at: https://github.com/blueberry-jam/secprac-client")
		os.Exit(0)
	}

	// Ensure running as root
	if os.Geteuid() != 0 {
		println("run the secprac client as root")
		os.Exit(1)
	}

	// init util.Logger
	util.Init()

	// Get user
	if len(os.Args) < 2 {
		util.Logger.Fatalln("no user provided, run again like this: `secprac-client <user> <server url>`")
	}
	username := os.Args[1]
	user, err := ou.Lookup(username)
	if err != nil {
		util.Logger.Fatalln("error fetching user with username", username+":", err)
	}

	// Get remote
	if len(os.Args) < 3 {
		util.Notify(user, "error", "no remote server was provided as a command-line argument", util.IconMinus, true)
		util.Logger.Fatalln("no server provided, run again like this: `secprac-client <user> <server url>`")
	}
	remote := os.Args[2]

	var token string
	team := &api.Team{}
	reco := false

	// Attempt to recover team data
	if _, err := os.Stat("/usr/local/share/secprac/team"); err == nil {
		reco = true
		b, err := ioutil.ReadFile("/usr/local/share/secprac/team")
		if err != nil {
			util.Logger.Println("failed to read existing team data file:", err)
			reco = false
		} else {
			data := strings.Split(string(b), " ")
			if len(data) > 1 {
				team.Token = data[0]
				team.ID = data[1]
				_, err = api.GetScripts(remote, team.Token)
				if err != nil {
					reco = false
					util.Logger.Println("token from team data file likely invalid, error communicating with server:", err)
				}
				if reco {
					util.Notify(user, "recovered", "successfully recovered valid team data from a crash, your ID is: "+team.ID, util.IconInfo, true)
					util.Logger.Println("successfully recovered token and ID")
				}
			} else {
				util.Logger.Println("existing team data file does not contain token and id")
				reco = false
			}
		}
	}

	// Generate new random token and authenticate with server
	if !reco {
		b := make([]byte, 18)
		_, err = rand.Read(b)
		if err != nil {
			util.Logger.Fatalln(err)
		}
		token = base64.URLEncoding.EncodeToString(b)
		util.Logger.Println("generated new team token")

		// Authenticate with the server
		util.Logger.Println("attempting to authenticate with server (" + remote + ")")
		team, err = api.NewTeam(remote, token)
		if err != nil {
			util.Notify(user, "error", "failed to authenticate with the server, check the log at: "+util.LogFileName, util.IconMinus, true)
			util.Logger.Fatalln("error authenticating with server:", err)
		}
		util.Notify(user, "authenticated", "successfully authenticated with the server, your team ID is "+team.ID, util.IconInfo, true)
		util.Logger.Println("authenticated with", remote, "given ID", team.ID)

		// Attempt to write to data file
		err = ioutil.WriteFile("/usr/local/share/secprac/team", []byte(team.Token+" "+team.ID), 0600)
		if err != nil {
			util.Logger.Println("error writing to team data file:", err)
		} else {
			util.Logger.Println("successfully saved team data")
		}
	}

	// Get the vulnerability-checking scripts
	team.Scripts, err = api.GetScripts(remote, team.Token)
	if err != nil {
		util.Notify(user, "error", "failed to get the script information from the server, check the log at: "+util.LogFileName, util.IconMinus, true)
		util.Logger.Fatalln("error getting script information from the server:", err)
	}
	if len(team.Scripts) < 1 {
		util.Notify(user, "error", "the server did not provide any scripts... you win?", util.IconPlus, true)
		util.Logger.Fatalln("server provided no scripts")
	}

	// Download scripts
	team.Scripts, err = api.DownloadScripts(remote, team.Token, team.Scripts)
	if err != nil {
		util.Notify(user, "error", "failed to download scripts from the server, check the log at: "+util.LogFileName, util.IconMinus, true)
		util.Logger.Fatalln("error downloading a script from the server:", err)
	}
	util.Notify(user, "downloaded scripts", "successfully downloaded "+strconv.Itoa(len(team.Scripts))+" scripts", util.IconInfo, false)

	// Run setup scripts
	count := 0
	for _, script := range team.Scripts {
		if script.Setup != "" {
			count++
		}
	}
	if count > 0 {
		util.Notify(user, "running setup scripts", "successfully downloaded "+strconv.Itoa(count)+" setup scripts and running them", util.IconInfo, false)
		for i := range team.Scripts {
			script := &team.Scripts[i]

			// Pipe script into shell and run
			cmd := exec.Command(script.Shell)
			cmd.Env = append(cmd.Env, "SECPRAC_USER="+user.Username)
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
				util.Logger.Println("error running setup script:", err)
				continue
			}
			if string(out) != "SETUP\n" {
				util.Logger.Println("setup script for", script.Name, "failed to complete successfully")
				util.Notify(user, "error", "failed to run some setup scripts, check the log at: "+util.LogFileName, util.IconMinus, true)
			}
		}
	}
	util.Notify(user, "start", "you may now start!", util.IconInfo, false)

	// Main loop
	for {
		done := true
		for i := range team.Scripts {
			script := &team.Scripts[i]

			// Pipe script into shell and run
			cmd := exec.Command(script.Shell)
			cmd.Env = append(cmd.Env, "SECPRAC_USER="+user.Username)
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

			// Check if fixed, and if the vuln has been fixed or undone
			if script.Fixed {
				if string(out) != "FIXED\n" {

					// Undo vuln
					points, err := api.VulnUndo(remote, team.Token, *script)
					if err != nil {
						util.Logger.Println("error when submitting undone vuln:", err)
						continue
					}
					util.Logger.Println("script undone:", script.Name)
					script.Fixed = false
					team.Points += points
					util.PointNotif(points, "undoing "+script.Name, user)
				}
			} else {
				if string(out) == "FIXED\n" {

					// Done vuln
					points, err := api.VulnDone(remote, team.Token, *script)
					if err != nil {
						util.Logger.Println("error when submitting done vuln:", err)
						continue
					}
					util.Logger.Println("script fixed:", script.Name)
					script.Fixed = true
					team.Points += points
					util.PointNotif(points, script.Name, user)
				}
			}

			// Team is not done if script isn't fixed
			if !script.Fixed {
				done = false
			}

			// Sleep for performance reasons
			time.Sleep(time.Second / 5)
		}

		// Check if team is done
		if done {
			err := api.TeamDone(remote, team.Token)
			if err != nil {
				util.Logger.Println("error telling server team is done:", err)
				continue
			}
			util.Logger.Println("done!")
			util.Notify(user, "complete", "you've successfully secured the system!", util.IconInfo, false)
			err = os.Remove("/usr/local/share/secprac/team")
			if err != nil {
				util.Logger.Println("error cleaning up old team data file:", err)
			} else {
				util.Logger.Println("successfully cleaned up old team data file")
			}
			break
		}
	}
}
