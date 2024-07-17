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

	bootstrapHooks *multiRegister[*BootstrapHook]
	cleanupHooks   *multiRegister[*CleanupHook]
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
		services: newMultiRegister[Service](itemName("service")),

		bootstrapHooks: newMultiRegister[*BootstrapHook](itemName("bootstrap hook"), allowUpdate, allowRemove),
		cleanupHooks:   newMultiRegister[*CleanupHook](itemName("cleanup hook"), allowUpdate, allowRemove),
	}
}

func (c *Cosys) Server() Server {
	server, _ := c.server.Get()

	return server
}

func (c *Cosys) Database() Database {
	database, _ := c.database.Get()

	return database
}

func (c *Cosys) Logger() Logger {
	logger, _ := c.logger.Get()

	return logger
}

func (c *Cosys) UseServer(server Server) error {
	return c.server.Register(server)
}

func (c *Cosys) UseDatabase(database Database) error {
	return c.database.Register(database)
}

func (c *Cosys) UseLogger(logger Logger) error {
	return c.logger.Register(logger)
}

func (c *Cosys) AddRoutes(routes ...*Route) error {
	return c.routes.RegisterStringers(routes...)
}

func (c *Cosys) UpdateRoute(uid string, route *Route) error {
	return c.routes.Update(uid, route)
}

func (c *Cosys) RemoveRoute(uid string) error {
	return c.routes.Remove(uid)
}

func (c *Cosys) AddControllers(controllers ...Controller) error {
	return c.controllers.RegisterStringers(controllers...)
}

func (c *Cosys) UpdateController(uid string, controller Controller) error {
	return c.controllers.Update(uid, controller)
}

func (c *Cosys) RemoveController(uid string) error {
	return c.controllers.Remove(uid)
}

func (c *Cosys) AddMiddlewares(middlewares ...Middleware) error {
	return c.middlewares.RegisterStringers(middlewares...)
}

func (c *Cosys) UpdateMiddleware(uid string, middleware Middleware) error {
	return c.middlewares.Update(uid, middleware)
}

func (c *Cosys) RemoveMiddleware(uid string) error {
	return c.middlewares.Remove(uid)
}

func (c *Cosys) AddPolicies(policies ...Policy) error {
	return c.policies.RegisterStringers(policies...)
}

func (c *Cosys) UpdatePolicy(uid string, policy Policy) error {
	return c.policies.Update(uid, policy)
}

func (c *Cosys) RemovePolicy(uid string) error {
	return c.policies.Remove(uid)
}

func (c *Cosys) AddCommands(commands ...*cobra.Command) error {
	return c.commands.RegisterStringers(commands...)
}

func (c *Cosys) AddModels(models ...Model) error {
	return c.models.RegisterStringers(models...)
}

func (c *Cosys) UpdateModel(uid string, model Model) error {
	return c.models.Update(uid, model)
}

func (c *Cosys) RemoveModel(uid string) error {
	return c.models.Remove(uid)
}

func (c *Cosys) AddServices(services ...Service) error {
	return c.services.RegisterStringers(services...)
}

func (c *Cosys) AddBootstrapHooks(hooks ...*BootstrapHook) error {
	return c.bootstrapHooks.RegisterStringers(hooks...)
}

func (c *Cosys) AddCleanupHooks(hooks ...*CleanupHook) error {
	return c.cleanupHooks.RegisterStringers(hooks...)
}

func (c *Cosys) Start() error {
	cosys, err := c.register()
	if err != nil {
		return err
	}

	return cosys.Command.Execute()
}

func (c *Cosys) register() (*Cosys, error) {
	for _, module := range mdRegister {
		if err := module(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Cosys) Bootstrap() error {
	for _, hook := range c.bootstrapHooks.GetAll() {
		if err := hook.Call(c); err != nil {
			return err
		}
	}

	return nil
}

func (c *Cosys) Destroy() error {
	for _, hook := range c.cleanupHooks.GetAll() {
		if err := hook.Call(c); err != nil {
			return err
		}
	}

	return nil
}
