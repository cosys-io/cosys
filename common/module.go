package common

import (
	"sync"
)

var (
	mdMutex    sync.RWMutex
	mdRegister = newPermRegister[Module]()
)

// RegisterModule registers a module to the cosys app.
func RegisterModule(moduleName string, module Module) error {
	return mdRegister.Register(moduleName, module)
}

// Module is a hook that is called during the registration stage.
type Module func(*Cosys) error
