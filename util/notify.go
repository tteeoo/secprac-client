package util

import (
	"os/exec"
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
func PointNotif(points int, name, user string) {
	if points > 0 {
		Notify(user, "gained points", "gained "+strconv.Itoa(points)+" point(s) for "+name, IconPlus, false)
	} else if points < 0 {
		Notify(user, "lost points", "lost "+strconv.Itoa(0-points)+" point(s) for "+name, IconMinus, false)
	} else {
		Notify(user, "vuln fixed", "fixed vulnerability: "+name, IconInfo, false)
	}
}

// Notify sends a desktop notification
func Notify(user, title, text, icon string, urgent bool) {
	var cmd *exec.Cmd
	if urgent {
		cmd = exec.Command("su", "-c", "DISPLAY=:* notify-send -u critical -a secprac -i \""+icon+"\" \""+title+"\" \""+text+"\"", user)
	} else {
		cmd = exec.Command("su", "-c", "notify-send -a secprac -i \""+icon+"\" \""+title+"\" \""+text+"\"", user)
	}

	cmd.Env = append(cmd.Env, "DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/1000/bus")
	cmd.Env = append(cmd.Env, "DISPLAY=:*")
	err := cmd.Run()
	if err != nil {
		Logger.Println("error sending notification:", err)
	}
}
