package common

import (
	"fmt"
	"sync"
)

var (
	loggerMutex    sync.RWMutex
	loggerRegister = make(map[string]Logger)
)

func RegisterLogger(loggerName string, logger Logger) error {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	if logger == nil {
		return fmt.Errorf("logger is nil: %s", loggerName)
	}

	if _, dup := loggerRegister[loggerName]; dup {
		return fmt.Errorf("duplicate logger: %s", loggerName)
	}

	loggerRegister[loggerName] = logger
	return nil
}

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
