package util

import (
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

// Logger points to a log.Logger which writes to stdout and a unique file
var Logger *log.Logger

// LogFileName is the name of the log file created
var LogFileName string

func init() {

	// Get logging directory
	logDir := os.Getenv("SECPRAC_LOG_DIR")
	if len(logDir) < 1 {
		logDir = "/var/log/secprac-client"
	}
	LogFileName := logDir + "/secprac-client-" + strconv.Itoa(int(time.Now().Unix())) + ".log"

	// Log to terminal and a file
	logFile, err := os.Create(LogFileName)
	if err != nil {
		log.Fatal(err)
	}

	logFile.Sync()

	Logger = log.New(io.MultiWriter(logFile, os.Stdout), "", log.Ldate|log.Ltime)
}
