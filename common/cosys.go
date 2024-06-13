package common

import "fmt"

type Cosys struct {
	Configs  *Configs
	Modules  map[string]*Module
	Services map[string]Service
	Models   map[string]Model
	Db       func(*Cosys) Database
	Ms       func(*Cosys) *ModuleService
	Sv       func(*Cosys) *Server
}

func (c Cosys) Database() Database {
	return c.Db(&c)
}

func (c Cosys) ModuleService() *ModuleService {
	return c.Ms(&c)
}

func (c Cosys) Server() *Server {
	return c.Sv(&c)
}

func NewCosys(configs *Configs) *Cosys {
	cosys := &Cosys{
		Configs:  configs,
		Modules:  map[string]*Module{},
		Services: map[string]Service{},
		Models:   map[string]Model{},
		Db:       nil,
		Ms:       func(cosys *Cosys) *ModuleService { return &ModuleService{Cosys: cosys} },
		Sv:       func(cosys *Cosys) *Server { return &Server{Port: "3000", Cosys: cosys} },
	}

	return cosys
}

func (c Cosys) Register(modules map[string]*Module) (*Cosys, error) {
	cosys := c
	var err error

	cosys.Modules = modules

	for _, module := range modules {
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
