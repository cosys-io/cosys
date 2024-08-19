package internal

import (
	"fmt"
	"github.com/cosys-io/cosys/common"
	"log"
)

// Logger is an implementation of the Logger core service using the native log package.
type Logger struct{}

// Log logs a message at the given log level.
func (l Logger) Log(stringer fmt.Stringer, logLevel common.LogLevel) {
	switch logLevel {
	case common.Debug:
		log.Println("DEBUG: ", stringer.String())
	case common.Info:
		log.Println("INFO: ", stringer.String())
	case common.Warn:
		log.Println("WARN: ", stringer.String())
	case common.Error:
		log.Println("ERROR: ", stringer.String())
	default:
		log.Println(stringer.String())
	}
}

// Info logs a message at the info level.
func (l Logger) Info(stringer fmt.Stringer) {
	l.Log(stringer, common.Info)
}

// Debug logs a message at the debug level.
func (l Logger) Debug(stringer fmt.Stringer) {
	l.Log(stringer, common.Debug)
}

// Warn logs a message at the warn level.
func (l Logger) Warn(stringer fmt.Stringer) {
	l.Log(stringer, common.Warn)
}

// Error logs a message at the error level.
func (l Logger) Error(stringer fmt.Stringer) {
	l.Log(stringer, common.Error)
}
