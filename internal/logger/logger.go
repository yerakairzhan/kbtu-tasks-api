package logger

import (
	"log"
	"os"
)

var Log = log.New(os.Stdout, "[tasks-api] ", log.Ldate|log.Ltime|log.Lshortfile)

func Infof(format string, args ...any) {
	Log.Printf("INFO: "+format, args...)
}

func Errorf(format string, args ...any) {
	Log.Printf("ERROR: "+format, args...)
}

func Fatalf(format string, args ...any) {
	Log.Fatalf("FATAL: "+format, args...)
}
