package logger

import (
	"io/ioutil"
	"log"
	"os"
)

var Logger *log.Logger

func init() {
	// Initialize logger to discard output by default
	Logger = log.New(ioutil.Discard, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// SetLogFile sets the output file for the logger.
func SetLogFile(filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	Logger.SetOutput(file)
}
