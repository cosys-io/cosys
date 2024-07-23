package sqlite3

import (
	"database/sql"
	"fmt"
	"github.com/cosys-io/cosys/common"
	"github.com/cosys-io/cosys/modules/sqlite3/internal"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db               *sql.DB
	BootstrapHookKey string
)

func init() {
	_ = common.RegisterModule(func(cosys *common.Cosys) error {
		var err error

		if err = cosys.UseDatabase(internal.NewDatabase(db, cosys)); err != nil {
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
	db, err = sql.Open("sqlite3", "data.db")
	if err != nil {
		return err
	}

	database, err := cosys.Database()
	if err != nil {
		return err
	}

	sqlite3, ok := database.(internal.Database)
	if !ok {
		return fmt.Errorf("database is not a sqlite3 database")
	}

	return sqlite3.LoadSchema()
}
