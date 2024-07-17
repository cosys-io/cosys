package common

import (
	"fmt"
	"sync"
)

var (
	mdMutex    sync.RWMutex
	mdRegister = make(map[string]Module)
)

func RegisterModule(moduleName string, module Module) error {
	if module == nil {
		return fmt.Errorf("module is nil: %s", moduleName)
	}

	mdMutex.Lock()
	defer mdMutex.Unlock()

	if _, dup := mdRegister[moduleName]; dup {
		return fmt.Errorf("duplicate module:" + moduleName)
	}

	mdRegister[moduleName] = module

	return nil
}

type Module func(*Cosys) error
