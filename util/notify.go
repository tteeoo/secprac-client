package util

import (
	"os/exec"
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
