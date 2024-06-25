package internal

import (
	"fmt"
	"github.com/cosys-io/cosys/common"
	"log"
)

type Logger struct{}

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

func (l Logger) Info(stringer fmt.Stringer) {
	l.Log(stringer, common.Info)
}

func (l Logger) Debug(stringer fmt.Stringer) {
	l.Log(stringer, common.Debug)
}

func (l Logger) Warn(stringer fmt.Stringer) {
	l.Log(stringer, common.Warn)
}

func (l Logger) Error(stringer fmt.Stringer) {
	l.Log(stringer, common.Error)
}
