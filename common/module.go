package common

import (
	"sync"
)

var (
	mdMutex    sync.RWMutex
	mdRegister = newPermRegister[Module]()
)

// RegisterModule registers a module to the cosys app.
// Safe for concurrent use.
func RegisterModule(module Module) error {
	_, err := mdRegister.RegisterRandom(module)
	return err
}

// Module is a hook that is called during the registration stage.
type Module func(*Cosys) error
