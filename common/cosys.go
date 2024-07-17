package common

import (
	"github.com/spf13/cobra"
)

type Cosys struct {
	server   *singleRegister[Server]
	database *singleRegister[Database]
	logger   *singleRegister[Logger]

	routes      *multiRegister[*Route]
	controllers *multiRegister[Controller]
	middlewares *multiRegister[Middleware]
	policies    *multiRegister[Policy]

	commands *multiRegister[*cobra.Command]
	models   *multiRegister[Model]
	services *multiRegister[Service]
}

func New() *Cosys {
	return &Cosys{
		server:   newSingleRegister[Server](itemName("server")),
		database: newSingleRegister[Database](itemName("database")),
		logger:   newSingleRegister[Logger](itemName("logger")),

		routes:      newMultiRegister[*Route](itemName("routes"), allowUpdate, allowRemove),
		controllers: newMultiRegister[Controller](itemName("controller"), allowUpdate, allowRemove),
		middlewares: newMultiRegister[Middleware](itemName("middleware"), allowUpdate, allowRemove),
		policies:    newMultiRegister[Policy](itemName("policies"), allowUpdate, allowRemove),

		commands: newMultiRegister[*cobra.Command](itemName("command")),
		models:   newMultiRegister[Model](itemName("model"), allowUpdate, allowRemove),
		services: newMultiRegister[Service](itemName("model")),
	}
}

func (c Cosys) Server() Server {
	server, _ := c.server.Get()

	return server
}

func (c Cosys) Database() Database {
	database, _ := c.database.Get()

	return database
}

func (c Cosys) Logger() Logger {
	logger, _ := c.logger.Get()

	return logger
}

func (c Cosys) UseServer(serverFunc func(*Cosys) (Server, error)) error {
	server, err := serverFunc(&c)
	if err != nil {
		return err
	}

	if err = c.server.Register(server); err != nil {
		return err
	}

	return nil
}

func (c Cosys) UseDatabase(databaseFunc func(*Cosys) (Database, error)) error {
	database, err := databaseFunc(&c)
	if err != nil {
		return err
	}

	if err = c.database.Register(database); err != nil {
		return err
	}

	return nil
}

func (c Cosys) UseLogger(loggerFunc func(*Cosys) (Logger, error)) error {
	logger, err := loggerFunc(&c)
	if err != nil {
		return err
	}

	if err = c.logger.Register(logger); err != nil {
		return err
	}

	return nil
}

func (c Cosys) AddRoutes(routes map[string]*Route) (*Cosys, error) {
	for uid, route := range routes {
		if err := c.routes.Register(uid, route); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

func (c Cosys) AddControllers(controllers map[string]Controller) (*Cosys, error) {
	for uid, controller := range controllers {
		if err := c.controllers.Register(uid, controller); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

func (c Cosys) AddMiddlewares(middlewares map[string]Middleware) (*Cosys, error) {
	for uid, middleware := range middlewares {
		if err := c.middlewares.Register(uid, middleware); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

func (c Cosys) AddPolicies(policies map[string]Policy) (*Cosys, error) {
	for uid, policy := range policies {
		if err := c.policies.Register(uid, policy); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

func (c Cosys) AddCommands(commands ...*cobra.Command) (*Cosys, error) {
	for _, command := range commands {
		if err := c.commands.Register(command.Name(), command); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

func (c Cosys) AddModels(models map[string]Model) (*Cosys, error) {
	for uid, model := range models {
		if err := c.models.Register(uid, model); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

func (c Cosys) AddServices(services map[string]Service) (*Cosys, error) {
	for uid, service := range services {
		if err := c.services.Register(uid, service); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

func (c Cosys) Start() error {
	cosys, err := c.register()
	if err != nil {
		return err
	}

	return cosys.Command.Execute()
}

func (c Cosys) register() (*Cosys, error) {
	for _, module := range mdRegister {
		if err := module(&c); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

func (c Cosys) Bootstrap() (*Cosys, error) {
	return &c, nil
}

func (c Cosys) Destroy() (*Cosys, error) {
	return &c, nil
}
