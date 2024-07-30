package server

import (
	"github.com/cosys-io/cosys/common"
	"github.com/cosys-io/cosys/modules/server/internal"
)

var (
	BootstrapHookKey string
)

func init() {
	_ = common.RegisterModule(func(cosys *common.Cosys) error {
		var err error

		if err := cosys.UseServer(internal.NewServer("3000", cosys)); err != nil {
			return err
		}

		BootstrapHookKey, err = cosys.AddBootstrapHook(Bootstrap)
		if err != nil {
			return err
		}

		return nil
	})
}

func Bootstrap(cosys *common.Cosys) error {
	var err error

	server, err := cosys.Server()
	if err != nil {
		return err
	}

	if err = server.ResolveEndpoints(); err != nil {
		return err
	}

	return nil
}
