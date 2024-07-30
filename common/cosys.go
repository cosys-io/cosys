package common

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// Cosys is the cosys app.
type Cosys struct {
	environment Environment
	state       State
	shutdown    <-chan os.Signal

	server   *singleRegister[Server]
	database *singleRegister[Database]
	logger   *singleRegister[Logger]

	routes      *stringerRegister[Route]
	controllers *stringerRegister[Controller]
	middlewares *stringerRegister[Middleware]
	policies    *stringerRegister[Policy]

	commands *stringerRegister[Command]
	models   *permRegister[Model]
	services *permRegister[Service]

	bootstrapHooks *multiRegister[BootstrapHook]
	cleanupHooks   *multiRegister[CleanupHook]
}

// New returns a new cosys instance, with modules registered.
func New() (*Cosys, error) {
	cosys := &Cosys{
		environment: "",
		state:       Registration,

		server:   newSingleRegister[Server](itemName("server")),
		database: newSingleRegister[Database](itemName("database")),
		logger:   newSingleRegister[Logger](itemName("logger")),

		routes:      newStringerRegister[Route](itemName("routes")),
		controllers: newStringerRegister[Controller](itemName("controller")),
		middlewares: newStringerRegister[Middleware](itemName("middleware")),
		policies:    newStringerRegister[Policy](itemName("policies")),

		commands: newStringerRegister[Command](itemName("command")),
		models:   newPermRegister[Model](itemName("model")),
		services: newPermRegister[Service](itemName("service")),

		bootstrapHooks: newMultiRegister[BootstrapHook](itemName("bootstrap hook")),
		cleanupHooks:   newMultiRegister[CleanupHook](itemName("cleanup hook")),
	}

	if err := cosys.AddCommands(serveCmd, devCmd, testCmd); err != nil {
		return nil, err
	}

	if err := cosys.register(); err != nil {
		return nil, err
	}

	return cosys, nil
}

// Environment returns the environment the cosys app is running in.
func (c *Cosys) Environment() Environment {
	return c.environment
}

// SetEnvironment specifies the environment the cosys app is running in.
func (c *Cosys) SetEnvironment(env Environment) {
	c.environment = env
}

// State returns the current state of the cosys app.
func (c *Cosys) State() State {
	return c.state
}

// ShutdownChannel returns a read-only channel that is sent to
// when the cosys app is interrupted or terminated.
func (c *Cosys) ShutdownChannel() <-chan os.Signal {
	return shutdownChannel()
}

// Server returns the server core service.
// Cannot be used during registration.
// Safe for concurrent use.
func (c *Cosys) Server() (Server, error) {
	if c.state == Registration {
		return nil, fmt.Errorf("server cannot be used during registration")
	}

	return c.server.Get()
}

// Database returns the database core service.
// Cannot be used during registration.
// Safe for concurrent use.
func (c *Cosys) Database() (Database, error) {
	if c.state == Registration {
		return nil, fmt.Errorf("database cannot be used during registration")
	}

	return c.database.Get()
}

// Logger returns the logger core service.
// Cannot be used during registration.
// Safe for concurrent use.
func (c *Cosys) Logger() (Logger, error) {
	if c.state == Registration {
		return nil, fmt.Errorf("logger cannot be used during registration")
	}

	return c.logger.Get()
}

// UseServer registers the server core service.
// Can only be used during registration.
// Safe for concurrent use.
func (c *Cosys) UseServer(server Server) error {
	if c.state != Registration {
		return fmt.Errorf("server must be registered during registration")
	}

	return c.server.Register(server)
}

// UseDatabase registers the database core service.
// Can only be used during registration.
// Safe for concurrent use.
func (c *Cosys) UseDatabase(database Database) error {
	if c.state != Registration {
		return fmt.Errorf("database must be registered during registration")
	}

	return c.database.Register(database)
}

// UseLogger registers the logger core service.
// Can only be used during registration.
// Safe for concurrent use.
func (c *Cosys) UseLogger(logger Logger) error {
	if c.state != Registration {
		return fmt.Errorf("logger must be registered during registration")
	}

	return c.logger.Register(logger)
}

func (c *Cosys) Routes() []Route {
	return c.routes.GetSlice()
}

// AddRoutes adds routes to the cosys app.
// Throws error if multiple routes have the same path.
// Safe for concurrent use.
func (c *Cosys) AddRoutes(routes ...Route) error {
	return c.routes.RegisterStringers(routes...)
}

// UpdateRoute updates a route specified by its path.
// Throws error if route with path does not exist.
// Safe for concurrent use.
func (c *Cosys) UpdateRoute(path string, route Route) error {
	return c.routes.Update(path, route)
}

// RemoveRoute removes a route specified by its path.
// Throws error if route with path does not exist.
// Safe for concurrent use.
func (c *Cosys) RemoveRoute(path string) error {
	return c.routes.Remove(path)
}

// AddControllers adds controllers to the cosys app.
// Throws error if multiple controllers have the same uid.
// Safe for concurrent use.
func (c *Cosys) AddControllers(controllers ...Controller) error {
	return c.controllers.RegisterStringers(controllers...)
}

// UpdateController updates a controller specified by its uid.
// Throws error if controller with uid does not exist.
// Safe for concurrent use.
func (c *Cosys) UpdateController(uid string, controller Controller) error {
	return c.controllers.Update(uid, controller)
}

// RemoveController removes a controller specified by its uid.
// Throws error if controller with uid does not exist.
// Safe for concurrent use.
func (c *Cosys) RemoveController(uid string) error {
	return c.controllers.Remove(uid)
}

// AddMiddlewares adds middlewares to the cosys app.
// Throws error if multiple middlewares have the same uid.
// Safe for concurrent use.
func (c *Cosys) AddMiddlewares(middlewares ...Middleware) error {
	return c.middlewares.RegisterStringers(middlewares...)
}

// UpdateMiddleware updates a middleware specified by its uid.
// Throws error if middleware with uid does not exist.
// Safe for concurrent use.
func (c *Cosys) UpdateMiddleware(uid string, middleware Middleware) error {
	return c.middlewares.Update(uid, middleware)
}

// RemoveMiddleware removes a middleware specified by its uid.
// Throws error if middleware with uid does not exist.
// Safe for concurrent use.
func (c *Cosys) RemoveMiddleware(uid string) error {
	return c.middlewares.Remove(uid)
}

// AddPolicies adds policies to the cosys app.
// Throws error if multiple policies have the same uid.
// Safe for concurrent use.
func (c *Cosys) AddPolicies(policies ...Policy) error {
	return c.policies.RegisterStringers(policies...)
}

// UpdatePolicy updates a policy specified by its uid.
// Throws error if policy with uid does not exist.
// Safe for concurrent use.
func (c *Cosys) UpdatePolicy(uid string, policy Policy) error {
	return c.policies.Update(uid, policy)
}

// RemovePolicy removes a policy specified by its uid.
// Throws error if policy with uid does not exist.
// Safe for concurrent use.
func (c *Cosys) RemovePolicy(uid string) error {
	return c.policies.Remove(uid)
}

// AddCommands adds commands to the cosys app.
// Throws error if multiple commands have the same name.
// Can only be used during registration.
// Safe for concurrent use.
func (c *Cosys) AddCommands(commands ...Command) error {
	if c.state != Registration {
		return fmt.Errorf("commands must be registered during registration")
	}

	return c.commands.RegisterStringers(commands...)
}

// AddModel adds a model to the cosys app.
// Throws error if multiple models have the same uid.
// Safe for concurrent use.
func (c *Cosys) AddModel(uid string, model Model) error {
	return c.models.Register(uid, model)
}

// AddModels adds models to the cosys app.
// Throws error if multiple models have the same uid.
// Safe for concurrent use.
func (c *Cosys) AddModels(models map[string]Model) error {
	return c.models.RegisterMany(models)
}

// Model returns a model by uid.
// Throws error if model with uid does not exist.
// Safe for concurrent use.
func (c *Cosys) Model(uid string) (Model, error) {
	return c.models.Get(uid)
}

// Models return a map of all models.
// Safe for concurrent use.
func (c *Cosys) Models() map[string]Model {
	return c.models.GetAll()
}

// AddService adds a service to the cosys app.
// Throws error if multiple routes have the same uid.
// Safe for concurrent use.
func (c *Cosys) AddService(uid string, service Service) error {
	return c.services.Register(uid, service)
}

// AddServices adds services to the cosys app.
// Throws error if multiple services have the same uid.
// Safe for concurrent use.
func (c *Cosys) AddServices(services map[string]Service) error {
	return c.services.RegisterMany(services)
}

// AddBootstrapHook adds a bootstrap hooks to the cosys app,
// and returns a uid that can be used to update or remove the hook.
// Safe for concurrent use.
func (c *Cosys) AddBootstrapHook(hook BootstrapHook) (string, error) {
	return c.bootstrapHooks.RegisterRandom(hook)
}

// UpdateBootstrapHook updates a bootstrap hook specified by its uid.
// Throws an error if hook with uid does not exist.
// Safe for concurrent use.
func (c *Cosys) UpdateBootstrapHook(uid string, hook BootstrapHook) error {
	return c.bootstrapHooks.Update(uid, hook)
}

// RemoveBootstrapHook removes a bootstrap hook specified by its uid.
// Throws an error if hook with uid does not exist.
// Safe for concurrent use.
func (c *Cosys) RemoveBootstrapHook(uid string) error {
	return c.bootstrapHooks.Remove(uid)
}

// AddCleanupHook adds a cleanup hook to the cosys app,
// and returns a uid that can be used to update or remove the hook.
// Safe for concurrent use.
func (c *Cosys) AddCleanupHook(hook CleanupHook) (string, error) {
	return c.cleanupHooks.RegisterRandom(hook)
}

// UpdateCleanupHook updates a cleanup hook specified by its uid.
// Throws an error if hook with uid does not exist.
// Safe for concurrent use.
func (c *Cosys) UpdateCleanupHook(uid string, hook BootstrapHook) error {
	return c.bootstrapHooks.Update(uid, hook)
}

// RemoveCleanupHook removes a cleanup hook specified by its uid.
// Throws an error if hook with uid does not exist.
// Safe for concurrent use.
func (c *Cosys) RemoveCleanupHook(uid string) error {
	return c.bootstrapHooks.Remove(uid)
}

// Start executes the cosys app command.
func (c *Cosys) Start() error {
	c.state = Execution

	command := &cobra.Command{}
	for _, cmd := range c.commands.GetAll() {
		command.AddCommand(cmd(c))
	}

	return command.Execute()
}

// register calls all registered module functions on the cosys instance.
func (c *Cosys) register() error {
	for _, module := range mdRegister.GetAll() {
		if err := module(c); err != nil {
			return err
		}
	}

	return nil
}

// startServer bootstraps the cosys app, starts the server,
// and cleanups the cosys app after shutdown.
func (c *Cosys) startServer() <-chan error {
	errCh := make(chan error, 1)

	if err := c.Bootstrap(); err != nil {
		errCh <- err
	}

	go func() {
		server, err := c.Server()
		if err != nil {
			errCh <- err
		}

		if err = server.Start(); err != nil {
			errCh <- err
		}
	}()

	go func() {
		<-c.ShutdownChannel()

		if err := c.Cleanup(); err != nil {
			errCh <- err
		}

		close(errCh)
	}()

	return errCh
}

// Bootstrap calls all bootstrap hooks added to the cosys instance.
func (c *Cosys) Bootstrap() error {
	c.state = Bootstrap

	for _, hook := range c.bootstrapHooks.GetAll() {
		if err := hook(c); err != nil {
			return err
		}
	}

	c.state = Execution

	return nil
}

// Cleanup calls all cleanup hooks added to the cosys instance.
func (c *Cosys) Cleanup() error {
	c.state = Cleanup

	for _, hook := range c.cleanupHooks.GetAll() {
		if err := hook(c); err != nil {
			return err
		}
	}

	c.state = Execution

	return nil
}
