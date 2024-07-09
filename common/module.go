package common

import (
	"fmt"
	"sync"
)

var (
	mdMutex    sync.RWMutex
	mdRegister = make(map[string]*Module)
)

func RegisterModule(moduleName string, module *Module) error {
	mdMutex.Lock()
	defer mdMutex.Unlock()

	if module == nil {
		return fmt.Errorf("module is nil: %s", moduleName)
	}

	if _, dup := mdRegister[moduleName]; dup {
		return fmt.Errorf("duplicate module:" + moduleName)
	}

	mdRegister[moduleName] = module
	return nil
}

type Module struct {
	Routes      []*Route
	Controllers map[string]Controller
	Middlewares map[string]Middleware
	Policies    map[string]Policy

	Models   map[string]Model
	Services map[string]Service

	OnRegister func(Cosys) (Cosys, error)
	OnDestroy  func(Cosys) (Cosys, error)
}
