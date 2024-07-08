package logging

import (
	"log"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Logging function called from main.go
func init() {
	Logging()
	log.Println("Logging initialised")
}
func Logging() {

	rollingLog()
	log.SetPrefix("SRV: ")                       // All messages will be prefixed by OWS:
	log.SetFlags(log.LstdFlags | log.Lshortfile) // Time, date,

}

func rollingLog() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "app.log",
		MaxSize:    1, // megabytes
		MaxBackups: 30,
		MaxAge:     28,    //days
		Compress:   false, // disabled by default
	})
}
