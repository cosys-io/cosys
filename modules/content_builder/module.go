package content_builder

import "github.com/cosys-io/cosys/common"

var Module = &common.Module{
	Routes: []*common.Route{
		common.NewRoute("GET", `/admin/schema`, "admin.schema"),
		common.NewRoute("POST", `/admin/schema/([a-zA-Z]+)`, "admin.build"),
	},
	Controllers: map[string]common.Controller{
		"admin": Controller,
	},
	Middlewares: nil,
	Policies:    nil,
	Models:      nil,
	Services:    nil,

	OnRegister: nil,
	OnDestroy:  nil,
}
