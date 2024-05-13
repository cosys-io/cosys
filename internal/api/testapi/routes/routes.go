package routes

import "github.com/cosys-io/cosys/internal/apis"

var Routes = []*apis.Route{
	apis.NewRoute("GET", `/users/([0-9]+)`, "users.findOne", apis.UseMiddlewares("dummy"), apis.UsePolicies("dummy")),
	apis.NewRoute("POST", `/users`, "users.create"),
	apis.NewRoute("PUT", `/users/([0-9]+)`, "users.update"),
	apis.NewRoute("DELETE", `/users/([0-9]+)`, "users.delete"),
}
