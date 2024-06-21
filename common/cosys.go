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
	database, ok := dbMap["sqlite3"]
	if !ok {
		log.Fatal("database not found: " + "sqlite3")
	}

	return database(&c)
}

func (c Cosys) ModuleService() ModuleService {
	moduleService, ok := msMap["default"]
	if !ok {
		log.Fatal("module service not found: " + "default")
	}

	return moduleService(&c)
}

func (c Cosys) Server() Server {
	server, ok := svMap["default"]
	if !ok {
		log.Fatal("server not found: " + "default")
	}

	return server(&c)
}

func NewCosys(configs Configs) *Cosys {
	cosys := &Cosys{
		Configs:  configs,
		Modules:  map[string]*Module{},
		Services: map[string]Service{},
		Models:   map[string]Model{},
	}

	return cosys
}

func (c Cosys) Register(modules map[string]*Module) (*Cosys, error) {
	cosys := c
	var err error

	cosys.Modules = map[string]*Module{}

	for _, moduleName := range c.Configs.Module.Modules {
		module, ok := modules[moduleName]
		if !ok {
			return nil, fmt.Errorf("module %s not found", moduleName)
		}

		cosys.Modules[moduleName] = module

		for name, model := range module.Models {
			if _, ok := cosys.Models[name]; ok {
				return nil, fmt.Errorf("model already exists: %s", name)
			}
			cosys.Models[name] = model
		}

		for name, service := range module.Services {
			if _, ok := cosys.Services[name]; ok {
				return nil, fmt.Errorf("service already exists: %s", name)
			}
			cosys.Services[name] = service
		}
	}

	for _, module := range cosys.Modules {
		if module.OnRegister == nil {
			continue
		}
		cosys, err = module.OnRegister(cosys)
		if err != nil {
			return nil, err
		}
	}

	return &cosys, nil
}

func (c Cosys) Start() error {
	if err := c.Server().Start(); err != nil {
		return err
	}

	return nil
}

func (c Cosys) Destroy() (*Cosys, error) {
	cosys := c
	var err error

	for _, module := range c.Modules {
		if module.OnDestroy == nil {
			continue
		}
		cosys, err = module.OnDestroy(cosys)
		if err != nil {
			return nil, err
		}
	}

	return &cosys, nil
}
