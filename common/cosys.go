package common

import (
	"fmt"
	"log"
)

type Cosys struct {
	Configs  Configs
	Apis     map[string]Api
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

func NewCosys() *Cosys {
	return &Cosys{
		Configs:  NewConfigs(),
		Services: make(map[string]Service),
		Models:   make(map[string]Model),
	}
}

func (c Cosys) register() (*Cosys, error) {
	newCosys := c
	var err error

	for moduleName, module := range mdRegister {
		newCosys.Apis[moduleName] = Api{
			Routes:      module.Routes,
			Controllers: module.Controllers,
			Middlewares: module.Middlewares,
			Policies:    module.Policies,
		}

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

	for _, module := range mdRegister {
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
	cosys, err := c.register()
	if err != nil {
		return err
	}

	if err = cosys.Server().Start(); err != nil {
		return err
	}

	return nil
}

func (c Cosys) Destroy() (*Cosys, error) {
	newCosys := c
	var err error

	for _, module := range mdRegister {
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
