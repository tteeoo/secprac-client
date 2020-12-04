package util

import (
	"io/ioutil"
	"os"
	"os/exec"
	ou "os/user"
	"strconv"
	"strings"
)

var (
	// IconPlus is the path to the plus icon
	IconPlus = "/usr/local/share/secprac/secprac-plus.png"

	// IconMinus is the path to the minus icon
	IconMinus = "/usr/local/share/secprac/secprac-minus.png"

	// IconInfo is the path to the info icon
	IconInfo = "/usr/local/share/secprac/secprac-info.png"
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
	run := "/run/user/" + user.Uid
	if d := os.Getenv("DBUS_SESSION_BUS_ADDRESS"); d != "" {
		cmd.Env = append(cmd.Env, "DBUS_SESSION_BUS_ADDRESS="+d)
	} else if _, err := os.Stat(run + "/bus"); !os.IsNotExist(err) {
		cmd.Env = append(cmd.Env, "DBUS_SESSION_BUS_ADDRESS=unix:path="+run+"/bus")
	} else if b, err := ioutil.ReadFile(run + "/dbus-session"); err == nil {
		if b[len(b)-1] == byte('\n') {
			cmd.Env = append(cmd.Env, string(b[:len(b)-1]))
		} else {
			cmd.Env = append(cmd.Env, string(b))
		}
	} else {
		launch := exec.Command("su", "-c", "dbus-launch --sh-syntax", user.Username)
		b, err := launch.Output()
		if err != nil {
			Logger.Println("error getting dbus session address for notification")
			return
		}
		cmd.Env = append(cmd.Env, strings.Replace(strings.Split(string(b), ";")[0], "'", "", 2))
	}
	cmd.Env = append(cmd.Env, "DISPLAY=:*")

	err := cmd.Run()
	if err != nil {
		Logger.Println("error sending notification:", err)
	}
}
