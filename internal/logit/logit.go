package logit

import (
	"io"
	"log"
	"os"
)

// Create loggers for File and File/Console output
var (
	Log  *log.Logger
	Raw  *log.Logger
	Loud *log.Logger
)

func init() {
	// Open Log file in temp dir
	logfile, err := os.CreateTemp("", "3gophers-listener-*.log")
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	// Open Raw file to log unchanged requests, this is going to come in handy when the indicator randomly change it's data structure
	rawfile, err := os.CreateTemp("", "3gophers-listener-raw-*.log")
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	// Set Log to logfile and console
	Log = log.New(logfile, "", log.LstdFlags)
	Loud = log.New(io.MultiWriter(os.Stdout, logfile), "", log.LstdFlags)
	Raw = log.New(rawfile, "", log.LstdFlags)

	Log.Println("Logit initialized: " + logfile.Name())
	Log.Println("Raw log initialized: " + rawfile.Name())
}
