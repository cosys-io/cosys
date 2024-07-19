package common

type Cosys struct {
	server   *singleRegister[Server]
	database *singleRegister[Database]
	logger   *singleRegister[Logger]

	routes      *stringerRegister[Route]
	controllers *stringerRegister[Controller]
	middlewares *stringerRegister[Middleware]
	policies    *stringerRegister[Policy]

	commands *permRegister[Command]
	models   *permRegister[Model]
	services *permRegister[Service]

	bootstrapHooks *stringerRegister[BootstrapHook]
	cleanupHooks   *stringerRegister[CleanupHook]
}

func New() (*Cosys, error) {
	cosys := &Cosys{
		server:   newSingleRegister[Server](itemName("server")),
		database: newSingleRegister[Database](itemName("database")),
		logger:   newSingleRegister[Logger](itemName("logger")),

		routes:      newStringerRegister[Route](itemName("routes")),
		controllers: newStringerRegister[Controller](itemName("controller")),
		middlewares: newStringerRegister[Middleware](itemName("middleware")),
		policies:    newStringerRegister[Policy](itemName("policies")),

		commands: newPermRegister[Command](itemName("command")),
		models:   newPermRegister[Model](itemName("model")),
		services: newPermRegister[Service](itemName("service")),

		bootstrapHooks: newStringerRegister[BootstrapHook](itemName("bootstrap hook")),
		cleanupHooks:   newStringerRegister[CleanupHook](itemName("cleanup hook")),
	}

	if err := cosys.register(); err != nil {
		return nil, err
	}

	return cosys, nil
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

func (c *Cosys) AddRoutes(routes ...Route) error {
	return c.routes.RegisterStringers(routes...)
}

func (c *Cosys) UpdateRoute(uid string, route Route) error {
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

func (c *Cosys) AddCommand(uid string, command Command) error {
	return c.commands.Register(uid, command)
}

func (c *Cosys) AddCommands(commands map[string]Command) error {
	return c.commands.RegisterMany(commands)
}

func (c *Cosys) AddModel(uid string, model Model) error {
	return c.models.Register(uid, model)
}

func (c *Cosys) AddModels(models map[string]Model) error {
	return c.models.RegisterMany(models)
}

func (c *Cosys) AddService(uid string, service Service) error {
	return c.services.Register(uid, service)
}

func (c *Cosys) AddServices(services map[string]Service) error {
	return c.services.RegisterMany(services)
}

func (c *Cosys) AddBootstrapHooks(hooks ...BootstrapHook) error {
	return c.bootstrapHooks.RegisterStringers(hooks...)
}

func (c *Cosys) AddCleanupHooks(hooks ...CleanupHook) error {
	return c.cleanupHooks.RegisterStringers(hooks...)
}

func (c *Cosys) Start() error {
	command := rootCmd(c)

	for _, cmd := range c.commands.GetAll() {
		command.AddCommand(cmd(c))
	}

	return command.Execute()
}

func (c *Cosys) register() error {
	for _, module := range mdRegister {
		if err := module(c); err != nil {
			return err
		}
	}

	return nil
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
