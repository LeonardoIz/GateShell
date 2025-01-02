package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// InitLogging sets up logging to a file and the console
func InitLogging(logDir string) error {
	// Create the necessary directories if they don't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// Set up logging to file
	logFileName := filepath.Join(logDir, fmt.Sprintf("%s.txt", time.Now().Format("2006-01-02--15-04-05")))
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	// Set up logging to console and file
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.Println("Logging to file", logFileName)

	return nil
}
