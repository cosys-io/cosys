package common

type Module struct {
	Routes      []*Route
	Controllers map[string]*Controller
	Middlewares map[string]Middleware
	Policies    map[string]Policy
	Lifecycles  map[string]LifeCycle
}
