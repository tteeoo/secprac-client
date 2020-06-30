package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"
	"os/exec"
	ou "os/user"
	"strconv"
	"time"

	"github.com/blueberry-jam/secprac-client/api"
	"github.com/blueberry-jam/secprac-client/util"
)

// Run as root!
func init() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-V") {
		println("secprac-client verison 0.1.4\nCopyright (C) Theo Henson\nMIT License\nOpen source at: https://github.com/blueberry-jam/secprac-client")
		os.Exit(0)
	}
	if os.Geteuid() != 0 {
		util.Logger.Fatalln("run the secprac client as root")
	}
}

// Entry point
func main() {

	// Handle command line arguments
	if len(os.Args) < 2 {
		util.Logger.Fatalln("no user provided, run again like this: `secprac-client <user> <server url>`")
	}
	username := os.Args[1]
	user, err := ou.Lookup(username)
	if err != nil {
		util.Logger.Fatalln("error fetch user with username", username+":", err)
	}
	if len(os.Args) < 3 {
		util.Notify(user, "error", "no remote server was provided as a command-line argument", util.IconMinus, true)
		util.Logger.Fatalln("no server provided, run again like this: `secprac-client <user> <server url>`")
	}
	remote := os.Args[2]

	// Generate random Token
	b := make([]byte, 18)
	_, err = rand.Read(b)
	if err != nil {
		util.Logger.Fatalln(err)
	}
	token := base64.URLEncoding.EncodeToString(b)
	util.Logger.Println("generated team token")

	// Authenticate with the server
	util.Logger.Println("attempting to authenticate with server (" + remote + ")")
	team, err := api.NewTeam(remote, token)
	if err != nil {
		util.Notify(user, "error", "failed to authenticate with the server, check the log at: "+util.LogFileName, util.IconMinus, true)
		util.Logger.Fatalln("error authenticating with server:", err)
	}

	util.Notify(user, "authenticated", "successfully authenticated with server, your team ID is "+team.ID, util.IconInfo, false)
	util.Logger.Println("authenticated with", remote, "given ID", team.ID)

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
	util.Notify(user, "downloaded scripts", "successfully downloaded "+strconv.Itoa(len(team.Scripts))+" scripts, start hacking!", util.IconInfo, false)

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
			break
		}
	}
}
