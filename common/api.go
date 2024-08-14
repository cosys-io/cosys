package common

import (
	"fmt"
	"net/http"
	"strings"
)

// Route

// Route specifies an api endpoint.
type Route struct {
	Method      string
	Path        string
	Action      ActionFunc
	Middlewares []MiddlewareFunc
	Policies    []PolicyFunc
}

// String returns the route path.
func (r Route) String() string {
	return r.Method + " " + r.Path
}

// NewRoute returns a new route with configurations.
func NewRoute(method string, path string, action ActionFunc, options ...RouteOption) Route {
	newRoute := Route{
		Method:      method,
		Path:        path,
		Action:      action,
		Middlewares: []MiddlewareFunc{},
		Policies:    []PolicyFunc{},
	}

	for _, option := range options {
		option(&newRoute)
	}

	return newRoute
}

// GetAction returns the actionFunc based on the actionFunc
// added to the cosys instance under the given uid.
// The actionFunc with the is retrieved from the cosys instance
// when the returned actionFunc is called, not when GetAction is called.
// If the actionFunc added under the given uid is updated or removed,
// the resulting actionFunc will use the updated actionFunc or throw an error if removed.
func GetAction(uid string) ActionFunc {
	return func(cosys *Cosys) (http.HandlerFunc, error) {
		uids := strings.Split(uid, ".")
		if len(uids) != 2 {
			return nil, fmt.Errorf("invalid uid")
		}
		controllerUid := uids[0]
		actionUid := uids[1]

		controller, err := cosys.controllers.Get(controllerUid)
		if err != nil {
			return nil, err
		}

		action, err := controller.Action(actionUid)
		if err != nil {
			return nil, err
		}

		actionFunc, err := action.Action(cosys)
		if err != nil {
			return nil, err
		}

		return actionFunc, nil
	}
}

// RouteOption is a route configuration.
type RouteOption func(*Route)

// UseMiddlewares adds middlewareFuncs to the route.
func UseMiddlewares(middlewares ...MiddlewareFunc) RouteOption {
	return func(route *Route) {
		if route == nil {
			return
		}

		route.Middlewares = append(route.Middlewares, middlewares...)
	}
}

// GetMiddlewares adds middlewareFuncs to the route,
// based on the middlewareFuncs added to the cosys instance under the given uids.
// The middlewareFuncs are retrieved from the cosys instance
// when the middlewareFuncs added to the route are called, not when GetMiddlewares is called.
// If the middlewareFuncs added under the given uid are updated or removed,
// the resulting middlewareFuncs will use the updated middlewareFuncs or throw errors if removed.
func GetMiddlewares(uids ...string) RouteOption {
	return func(route *Route) {
		middlewares := make([]MiddlewareFunc, len(uids))
		for _, uid := range uids {
			middleware := func(cosys *Cosys) (func(http.HandlerFunc) http.HandlerFunc, error) {
				m, err := cosys.middlewares.Get(uid)
				if err != nil {
					return nil, err
				}

				middlewareFunc, err := m.Middleware(cosys)
				if err != nil {
					return nil, err
				}

				return middlewareFunc, nil
			}

			middlewares = append(middlewares, middleware)
		}

		route.Middlewares = append(route.Middlewares, middlewares...)
	}
}

// UsePolicies adds policyFuncs to the route.
func UsePolicies(policies ...PolicyFunc) RouteOption {
	return func(route *Route) {
		if route == nil {
			return
		}

		route.Policies = append(route.Policies, policies...)
	}
}

// GetPolicies adds policyFuncs to the route,
// based on the policyFuncs added to the cosys instance under the given uids.
// The policyFuncs are retrieved from the cosys instance
// when the policyFuncs added to the route are called, not when GetPolicies is called.
// If the policyFuncs added under the given uid are updated or removed,
// the resulting policyFuncs will use the updated policyFuncs or throw errors if removed.
func GetPolicies(uids ...string) RouteOption {
	return func(route *Route) {
		policies := make([]PolicyFunc, len(uids))
		for _, uid := range uids {
			policy := func(cosys *Cosys) (func(*http.Request) bool, error) {
				p, err := cosys.policies.Get(uid)
				if err != nil {
					return nil, err
				}

				policyFunc, err := p.Policy(cosys)
				if err != nil {
					return nil, err
				}

				return policyFunc, nil
			}

			policies = append(policies, policy)
		}

		route.Policies = append(route.Policies, policies...)
	}
}

// Controller

// Controller is a group of actions.
type Controller struct {
	uid        string
	controller map[string]Action
}

// String returns the controller uid.
func (c Controller) String() string {
	return c.uid
}

// Action returns the controller action by its uid.
func (c Controller) Action(uid string) (Action, error) {
	if c.controller == nil {
		return Action{}, fmt.Errorf("controller is nil")
	}

	action, ok := c.controller[uid]
	if !ok {
		return Action{}, fmt.Errorf("action not found: %s", uid)
	}

	return action, nil
}

// NewController returns a new controller with the given uid and actions,
// or throws an error if multiple actions have the same uid.
func NewController(uid string, actions map[string]ActionFunc) (Controller, error) {
	controller := map[string]Action{}
	for actionUid, actionFunc := range actions {
		action, err := NewAction(actionUid, actionFunc)
		if err != nil {
			return Controller{}, err
		}

		controller[actionUid] = action
	}

	return Controller{
		uid:        uid,
		controller: controller,
	}, nil
}

// Action

// ActionFunc takes in cosys instance and returns a handler.
type ActionFunc func(*Cosys) (http.HandlerFunc, error)

// Action is a wrapper around ActionFunc that allows them to be identifiable by uid.
type Action struct {
	uid        string
	actionFunc ActionFunc
}

// String returns the action uid.
func (a Action) String() string {
	return a.uid
}

// Action returns the handler.
func (a Action) Action(cosys *Cosys) (http.HandlerFunc, error) {
	if cosys == nil {
		return nil, fmt.Errorf("cosys is nil")
	}

	if a.actionFunc == nil {
		return nil, fmt.Errorf("actionFunc is nil")
	}

	return a.actionFunc(cosys)
}

// NewAction returns a new action with the given uid and actionFunc.
func NewAction(uid string, action ActionFunc) (Action, error) {
	if action == nil {
		return Action{}, fmt.Errorf("action is nil")
	}

	return Action{
		uid:        uid,
		actionFunc: action,
	}, nil
}

// Middleware

// MiddlewareFunc takes in a cosys instance and returns a middleware.
type MiddlewareFunc func(*Cosys) (func(http.HandlerFunc) http.HandlerFunc, error)

// Middleware is a wrapper around MiddlewareFunc that allows them to be identifiable by uid.
type Middleware struct {
	uid            string
	middlewareFunc MiddlewareFunc
}

// String returns the middleware uid.
func (m Middleware) String() string {
	return m.uid
}

// Middleware returns the middleware.
func (m Middleware) Middleware(cosys *Cosys) (func(http.HandlerFunc) http.HandlerFunc, error) {
	if cosys == nil {
		return nil, fmt.Errorf("middlewareFunc is nil")
	}

	if m.middlewareFunc == nil {
		return nil, fmt.Errorf("middleware is nil")
	}

	return m.middlewareFunc(cosys)
}

// NewMiddleware returns a new middleware with the given uid and middlewareFunc.
func NewMiddleware(uid string, middleware MiddlewareFunc) (Middleware, error) {
	if middleware == nil {
		return Middleware{}, fmt.Errorf("middleware is nil")
	}

	return Middleware{
		uid:            uid,
		middlewareFunc: middleware,
	}, nil
}

// Policy

// PolicyFunc takes in a cosys instance and returns a policy.
type PolicyFunc func(*Cosys) (func(*http.Request) bool, error)

// Policy is a wrapper around PolicyFunc that allows them to be identifiable by uid.
type Policy struct {
	uid        string
	policyFunc PolicyFunc
}

// String returns the policy uid.
func (p Policy) String() string {
	return p.uid
}

// Policy returns the policy.
func (p Policy) Policy(cosys *Cosys) (func(*http.Request) bool, error) {
	if cosys == nil {
		return nil, fmt.Errorf("policyFunc is nil")
	}

	if p.policyFunc == nil {
		return nil, fmt.Errorf("policyFunc is nil")
	}

	return p.policyFunc(cosys)
}

// NewPolicy returns a new policy with the given uid and policyFunc.
func NewPolicy(uid string, policy PolicyFunc) (Policy, error) {
	if policy == nil {
		return Policy{}, fmt.Errorf("policyFunc is nil")
	}

	return Policy{
		uid:        uid,
		policyFunc: policy,
	}, nil
}
