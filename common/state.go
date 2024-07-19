package common

import (
	"os"
	"os/signal"
	"syscall"
)

type Environment string

const (
	Dev  Environment = "development"
	Test Environment = "test"
	Prod Environment = "production"
	Cmd  Environment = "command"
)

type State string

const (
	Registration State = "registration"
	Bootstrap    State = "bootstrap"
	Execution    State = "execution"
	Cleanup      State = "cleanup"
)

func shutdownChannel() <-chan os.Signal {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	return shutdownChan
}
