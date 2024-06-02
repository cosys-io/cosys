package apis

import (
	"regexp"
)

type Route struct {
	Method      string
	Regex       *regexp.Regexp
	Action      string
	Middlewares []string
	Policies    []string
}

func NewRoute(method string, route string, action string, options ...RouteOption) *Route {
	regexp.MustCompile(`(.+)\.(.+)`)

	newRoute := &Route{
		method,
		regexp.MustCompile(`^` + route + `$`),
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
