package common

import (
	"fmt"
	"log"
)

type Cosys struct {
	Configs  Configs
	Api      *Api
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
		Configs: NewConfigs(),
		Api: &Api{
			Routes:      []*Route{},
			Controllers: map[string]Controller{},
			Middlewares: map[string]Middleware{},
			Policies:    map[string]Policy{},
		},
		Services: make(map[string]Service),
		Models:   make(map[string]Model),
	}
}

func (c Cosys) AddRoutes(routes ...*Route) (*Cosys, error) {
	newCosys := &c
	api := newCosys.Api

	api.Routes = append(api.Routes, routes...)

	return newCosys, nil
}

func (c Cosys) AddControllers(controllers map[string]Controller) (*Cosys, error) {
	newCosys := &c
	api := newCosys.Api

	for controllerUid, controller := range controllers {
		if _, dup := api.Controllers[controllerUid]; dup {
			return nil, fmt.Errorf("controller already exists: %s", controllerUid)
		}
		api.Controllers[controllerUid] = controller
	}

	return newCosys, nil
}

func (c Cosys) AddMiddlewares(middlewares map[string]Middleware) (*Cosys, error) {
	newCosys := &c
	api := newCosys.Api

	for middlewareUid, middleware := range middlewares {
		if _, dup := api.Middlewares[middlewareUid]; dup {
			return nil, fmt.Errorf("middleware already exists: %s", middlewareUid)
		}
		api.Middlewares[middlewareUid] = middleware
	}

	return newCosys, nil
}

func (c Cosys) AddPolicies(policies map[string]Policy) (*Cosys, error) {
	newCosys := &c
	api := newCosys.Api

	for policyUid, policy := range policies {
		if _, dup := api.Policies[policyUid]; dup {
			return nil, fmt.Errorf("policy already exists: %s", policyUid)
		}
		api.Policies[policyUid] = policy
	}

	return newCosys, nil
}

func (c Cosys) register() (*Cosys, error) {
	newCosys := c
	var err error

	for _, module := range mdRegister {
		api := newCosys.Api
		api.Routes = append(api.Routes, module.Routes...)

		for controllerUid, controller := range module.Controllers {
			if _, dup := api.Controllers[controllerUid]; dup {
				return nil, fmt.Errorf("controller already exists: %s", controllerUid)
			}
			api.Controllers[controllerUid] = controller
		}

		for middlewareUid, middleware := range module.Middlewares {
			if _, dup := api.Middlewares[middlewareUid]; dup {
				return nil, fmt.Errorf("middleware already exists: %s", middlewareUid)
			}
			api.Middlewares[middlewareUid] = middleware
		}

		for policyUid, policy := range module.Policies {
			if _, dup := api.Policies[policyUid]; dup {
				return nil, fmt.Errorf("policy already exists: %s", policyUid)
			}
			api.Policies[policyUid] = policy
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

	defer func(cosys *Cosys) {
		_, err := cosys.destroy()
		if err != nil {

		}
	}(cosys)

	if err = cosys.Server().Start(); err != nil {
		return err
	}

	return nil
}

func (c Cosys) destroy() (*Cosys, error) {
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
