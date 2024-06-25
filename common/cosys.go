package common

import (
	"fmt"
	"log"
)

type Cosys struct {
	Configs  Configs
	Modules  map[string]*Module
	Services map[string]Service
	Models   map[string]Model
}

func (c Cosys) Database() Database {
	database, ok := dbRegister["sqlite3"]
	if !ok {
		log.Fatal("database not found: " + "sqlite3")
	}

	return database(&c)
}

func (c Cosys) Logger() Logger {
	logger, ok := loggerRegister["default"]
	if !ok {
		log.Fatal("logger not found: " + "default")
	}

	return logger
}

func (c Cosys) Server() Server {
	server, ok := svRegister["default"]
	if !ok {
		log.Fatal("server not found: " + "default")
	}

	return server(&c)
}

func NewCosys(cfg Configs) *Cosys {
	return &Cosys{
		Configs:  cfg,
		Modules:  make(map[string]*Module),
		Services: make(map[string]Service),
		Models:   make(map[string]Model),
	}
}

func (c Cosys) Register(modules map[string]*Module) (*Cosys, error) {
	newCosys := c
	var err error

	if modules == nil {
		modules = make(map[string]*Module)
	}

	newCosys.Modules = make(map[string]*Module)

	for _, moduleName := range newCosys.Configs.Module.Modules {
		module, ok := modules[moduleName]
		if !ok {
			return nil, fmt.Errorf("module not found: %s", moduleName)
		}

		newCosys.Modules[moduleName] = module

		for modelUid, model := range module.Models {
			if _, dup := newCosys.Models[modelUid]; dup {
				return nil, fmt.Errorf("model already exists: %s", modelUid)
			}
			newCosys.Models[modelUid] = model
		}

		for serviceUid, service := range module.Services {
			if _, dup := newCosys.Services[serviceUid]; dup {
				return nil, fmt.Errorf("service already exists: %s", serviceUid)
			}
			newCosys.Services[serviceUid] = service
		}
	}

	for _, module := range newCosys.Modules {
		if module.OnRegister == nil {
			continue
		}
		newCosys, err = module.OnRegister(newCosys)
		if err != nil {
			return nil, err
		}
	}

	return &newCosys, nil
}

func (c Cosys) Start() error {
	if err := c.Server().Start(); err != nil {
		return err
	}

	return nil
}

func (c Cosys) Destroy() (*Cosys, error) {
	newCosys := c
	var err error

	for _, module := range newCosys.Modules {
		if module.OnDestroy == nil {
			continue
		}
		newCosys, err = module.OnDestroy(newCosys)
		if err != nil {
			return nil, err
		}
	}

	return &newCosys, nil
}
