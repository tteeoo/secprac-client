package util

import (
	"io/ioutil"
	"os"
	"os/exec"
	ou "os/user"
	"strconv"
)

const (
	// IconPlus is the path to the plus icon
	IconPlus = "/usr/local/share/secprac/secprac-plus.png"

	// IconMinus is the path to the minus icon
	IconMinus = "/usr/local/share/secprac/secprac-minus.png"

	// IconInfo is the path to the info icon
	IconInfo = "/usr/local/share/secprac/secprac-info.png"

	// IconLogo is the path to the logo icon
	IconLogo = "/usr/local/share/secprac/secprac-logo.png"
)

// PointNotif takes a value of points and script name, then sends an appropriate notification
func PointNotif(points int, name string, user *ou.User) {
	if points > 0 {
		Notify(user, "gained points", "gained "+strconv.Itoa(points)+" point(s) for "+name, IconPlus, false)
	} else if points < 0 {
		Notify(user, "lost points", "lost "+strconv.Itoa(0-points)+" point(s) for "+name, IconMinus, false)
	} else {
		Notify(user, "vuln fixed", "fixed vulnerability: "+name, IconInfo, false)
	}
}

// Notify sends a desktop notification
func Notify(user *ou.User, title, text, icon string, urgent bool) {
	var cmd *exec.Cmd

	if urgent {
		cmd = exec.Command("su", "-c", "notify-send -u critical -a secprac -i \""+icon+"\" \""+title+"\" \""+text+"\"", user.Username)
	} else {
		cmd = exec.Command("su", "-c", "notify-send -a secprac -i \""+icon+"\" \""+title+"\" \""+text+"\"", user.Username)
	}

	// Get dbus address (method varies from distro to distro)
	if _, err := os.Stat("/run/user/"+user.Uid+"/dbus-session"); os.IsNotExist(err) {
		if _, err := os.Stat("/run/user/"+user.Uid+"/bus"); os.IsNotExist(err) {
			Logger.Println("error getting dbus session address for notification:", err)
			return
		}
		cmd.Env = append(cmd.Env, "DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/"+user.Uid+"/bus")
	} else {
		b, err := ioutil.ReadFile("/run/user/"+user.Uid+"/dbus-session")
		cmd.Env = append(cmd.Env, string(b))
		if err != nil {
			Logger.Println("error getting dbus session address for notification:", err)
			return
		}
	}
	cmd.Env = append(cmd.Env, "DISPLAY=:*")

	err := cmd.Run()
	if err != nil {
		Logger.Println("error sending notification:", err)
	}
}
