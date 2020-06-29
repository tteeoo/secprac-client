package util

import (
	"strconv"
	"time"
)

// GetTimestamp returns a timestamp of the current time
func GetTimestamp() string {
	var timestamp string
	t := time.Now().UTC()
	y, m, d := t.Date()
	timestamp += strconv.Itoa(y) + "/" + strconv.Itoa(int(m)) + "/" + strconv.Itoa(d) + " " + strconv.Itoa(t.Hour()) + ":" + strconv.Itoa(t.Minute())
	return timestamp
}
