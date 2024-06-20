package server

import (
	"github.com/cosys-io/cosys/common"
	"log"
)

var (
	port string
)

func init() {
	svFunc := func(cosys *common.Cosys) common.Server {
		return Server{
			Port:  port,
			Cosys: cosys,
		}
	}

	if err := common.RegisterServer("default", svFunc); err != nil {
		log.Fatal(err)
	}
}

func OnRegister(cosys common.Cosys) (common.Cosys, error) {
	port = cosys.Configs.Server.Port
	return cosys, nil
}

var Module = &common.Module{
	Routes:      nil,
	Controllers: nil,
	Middlewares: nil,
	Policies:    nil,
	Models:      nil,
	Services:    nil,

	OnRegister: OnRegister,
	OnDestroy:  nil,
}
