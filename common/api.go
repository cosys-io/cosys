package common

import (
	"net/http"
)

type Route struct {
	Method      string
	Path        string
	Action      string
	Middlewares []string
	Policies    []string
}

func NewRoute(method string, path string, action string, options ...RouteOption) *Route {
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

	return &newRoute
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

type Controller map[string]Action

type Action func(Cosys) http.HandlerFunc

type Middleware func(Cosys) func(http.HandlerFunc) http.HandlerFunc

type Policy func(Cosys, *http.Request) bool
