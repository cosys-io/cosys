package cms

import (
	"github.com/cosys-io/cosys/common"
	"log"
)

func init() {
	if err := common.RegisterCommand(rootCmd); err != nil {
		log.Fatal(err)
	}

	if err := common.RegisterModule("admin", adminModule); err != nil {
		log.Fatal(err)
	}
}

var adminModule = &common.Module{
	Routes: schemaRoutes,
	Controllers: map[string]common.Controller{
		"admin": schemaController,
	},
	Middlewares: nil,
	Policies:    nil,
	Models:      nil,
	Services:    nil,

	OnRegister: onRegister,
	OnDestroy:  nil,
}
