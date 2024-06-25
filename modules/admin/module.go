package admin

import "github.com/cosys-io/cosys/common"

var Module = &common.Module{
	Routes: SchemaRoutes,
	Controllers: map[string]common.Controller{
		"admin": SchemaController,
	},
	Middlewares: nil,
	Policies:    nil,
	Models:      nil,
	Services:    nil,

	OnRegister: OnRegister,
	OnDestroy:  nil,
}
