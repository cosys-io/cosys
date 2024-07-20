package common

import (
	"fmt"
)

// LogLevel is the alert level of a log message.
type LogLevel string

const (
	Debug LogLevel = "DEBUG"
	Info  LogLevel = "INFO"
	Warn  LogLevel = "WARN"
	Error LogLevel = "ERROR"
)

// Logger is a core service for logging.
type Logger interface {
	Log(stringer fmt.Stringer, logLevel LogLevel)
	Info(stringer fmt.Stringer)
	Warn(stringer fmt.Stringer)
	Debug(stringer fmt.Stringer)
	Error(stringer fmt.Stringer)
}
