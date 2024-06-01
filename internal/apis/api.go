package apis

import "fmt"

type API struct {
	Routes      []*Route
	controllers map[string]*Controller
	middlewares map[string]Middleware
	policies    map[string]Policy
}

func (a *API) Controller(uid string) (*Controller, error) {
	controller := a.controllers[uid]
	if controller == nil {
		return nil, fmt.Errorf("invalid uid: %s", uid)
	}

	return controller, nil
}

func (a *API) Middleware(uid string) (Middleware, error) {
	middleware := a.middlewares[uid]
	if middleware == nil {
		return nil, fmt.Errorf("invalid uid: %s", uid)
	}

	return middleware, nil
}

func (a *API) Policy(uid string) (Policy, error) {
	policy := a.policies[uid]
	if policy == nil {
		return nil, fmt.Errorf("invalid uid: %s", uid)
	}

	return policy, nil
}

func NewAPI(routes []*Route, controllers map[string]*Controller, middlewares map[string]Middleware, policies map[string]Policy) *API {
	return &API{
		routes,
		controllers,
		middlewares,
		policies,
	}
}
