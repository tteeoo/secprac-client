package util

import (
	"os/exec"
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

// Notify sends a desktop notification
func Notify(title, text, icon string, urgent bool) {
	var cmd *exec.Cmd
	if urgent {
		cmd = exec.Command("notify-send", "-a", "secprac", "-i", icon, title, text, "-u", "critical")
	} else {
		cmd = exec.Command("notify-send", "-a", "secprac", "-i", icon, title, text)
	}

	err := cmd.Run()
	if err != nil {
		Logger.Println("error sending notification:", err)
	}
}
