package common

import (
	"net/http"
)

// Routes

type Route struct {
	Method      string
	Path        string
	Action      string
	Middlewares []string
	Policies    []string
}

func NewRoute(method string, path string, action string, options ...RouteOption) *Route {
	newRoute := &Route{
		method,
		path,
		action,
		[]string{},
		[]string{},
	}

	for _, option := range options {
		option(newRoute)
	}

	return newRoute
}

type RouteOption func(*Route)

func UseMiddlewares(middlewares ...string) RouteOption {
	return func(route *Route) {
		route.Middlewares = append(route.Middlewares, middlewares...)
	}
}

func UsePolicies(policies ...string) RouteOption {
	return func(route *Route) {
		route.Policies = append(route.Policies, policies...)
	}
}

// Controllers

type Controller map[string]Action

type Action func(Cosys) http.HandlerFunc

// Middlewares

type Middleware func(Cosys) func(http.HandlerFunc) http.HandlerFunc

// Policies

type Policy func(Cosys, *http.Request) bool
