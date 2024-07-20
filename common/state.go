package common

import (
	"os"
	"os/signal"
	"syscall"
)

// Environment specifies which environment the app is running in.
type Environment string

const (
	Dev  Environment = "development"
	Test Environment = "test"
	Prod Environment = "production"
	Cmd  Environment = "command"
)

// State specifies which stage the cosys app is in.
type State string

const (
	Registration State = "registration"
	Bootstrap    State = "bootstrap"
	Execution    State = "execution"
	Cleanup      State = "cleanup"
)

// shutdownChannel returns a read-only channel that is sent to
// when the cosys app is interrupted or terminated.
func shutdownChannel() <-chan os.Signal {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	return shutdownChan
}
