package sqlite3

import (
	"github.com/cosys-io/cosys/common"
	"github.com/cosys-io/cosys/modules/sqlite3/internal"
	_ "github.com/mattn/go-sqlite3"
)

var (
	database         *internal.Database // database is the Database core service.
	BootstrapHookKey string             // BootstrapHookKey can be used to update or remove the bootstrap hook.
)

// init registers the module to register the Database core service and the bootstrap hook.
func init() {
	_ = common.RegisterModule(func(cosys *common.Cosys) error {
		var err error

		database = internal.NewDatabase(cosys)
		if err = cosys.UseDatabase(database); err != nil {
			return err
		}

		BootstrapHookKey, err = cosys.AddBootstrapHook(bootstrap)
		if err != nil {
			return err
		}

		return nil
	})
}

// bootstrap opens the connection to the SQLite3 database and
// loads the schema for all registered models.
func bootstrap(cosys *common.Cosys) error {
	if err := database.Open("data.db"); err != nil {
		return err
	}

	return database.LoadSchema()
}
