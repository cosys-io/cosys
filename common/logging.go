package common

import (
	"fmt"
)

type LogLevel string

const (
	Debug LogLevel = "DEBUG"
	Info  LogLevel = "INFO"
	Warn  LogLevel = "WARN"
	Error LogLevel = "ERROR"
)

type Logger interface {
	Log(stringer fmt.Stringer, logLevel LogLevel)
	Info(stringer fmt.Stringer)
	Warn(stringer fmt.Stringer)
	Debug(stringer fmt.Stringer)
	Error(stringer fmt.Stringer)
}
