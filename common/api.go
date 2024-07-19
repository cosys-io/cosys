package common

import (
	"net/http"
)

type Api struct {
	Routes      []*Route
	Controllers map[string]Controller
	Middlewares map[string]Middleware
	Policies    map[string]Policy
}

// Route

type Route struct {
	Method      string
	Path        string
	Action      string
	Middlewares []string
	Policies    []string
}

func (r Route) String() string {
	return r.Path
}

func NewRoute(method string, path string, action string, options ...RouteOption) Route {
	newRoute := Route{
		Method:      method,
		Path:        path,
		Action:      action,
		Middlewares: []string{},
		Policies:    []string{},
	}

	for _, option := range options {
		option(&newRoute)
	}

	return newRoute
}

type RouteOption func(*Route)

func UseMiddlewares(middlewares ...string) RouteOption {
	return func(route *Route) {
		if route == nil {
			return
		}

		route.Middlewares = append(route.Middlewares, middlewares...)
	}
}

func UsePolicies(policies ...string) RouteOption {
	return func(route *Route) {
		if route == nil {
			return
		}

		route.Policies = append(route.Policies, policies...)
	}
}

// Controller

type Controller struct {
	uid        string
	controller map[string]Action
}

func (c Controller) String() string {
	return c.uid
}

func NewController(uid string, controller map[string]Action) Controller {
	return Controller{
		uid:        uid,
		controller: controller,
	}
}

type Action struct {
	uid    string
	action func(Cosys) http.HandlerFunc
}

func (a Action) String() string {
	return a.uid
}

func NewAction(uid string, action func(Cosys) http.HandlerFunc) Action {
	return Action{
		uid:    uid,
		action: action,
	}
}

// Middleware

type Middleware struct {
	uid        string
	middleware func(Cosys) func(http.HandlerFunc) http.HandlerFunc
}

func (m Middleware) String() string {
	return m.uid
}

func NewMiddleware(uid string, middleware func(Cosys) func(http.HandlerFunc) http.HandlerFunc) Middleware {
	return Middleware{
		uid:        uid,
		middleware: middleware,
	}
}

// Policy

type Policy struct {
	uid    string
	policy func(Cosys, *http.Request) bool
}

func (p Policy) String() string {
	return p.uid
}

func NewPolicy(uid string, policy func(Cosys, *http.Request) bool) Policy {
	return Policy{
		uid:    uid,
		policy: policy,
	}
}
