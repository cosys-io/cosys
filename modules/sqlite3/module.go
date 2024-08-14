package sqlite3

import (
	"github.com/cosys-io/cosys/common"
	"github.com/cosys-io/cosys/modules/sqlite3/internal"
	_ "github.com/mattn/go-sqlite3"
)

var (
	database         *internal.Database
	BootstrapHookKey string
)

func init() {
	_ = common.RegisterModule(func(cosys *common.Cosys) error {
		var err error

		database = internal.NewDatabase(cosys)
		if err = cosys.UseDatabase(database); err != nil {
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
	if err := database.Open("data.db"); err != nil {
		return err
	}

	return database.LoadSchema()
}
