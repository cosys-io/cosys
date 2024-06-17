package sqlite3

import (
	"database/sql"
	"github.com/cosys-io/cosys/common"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Fatal(err)
	}

	dbFunc := func(cosys *common.Cosys) common.Database {
		return Database{
			cosys: cosys,
			db:    db,
		}
	}

	if err = common.RegisterDatabase("sqlite3", dbFunc); err != nil {
		log.Fatal(err)
	}
}

var Module = &common.Module{
	Routes:      nil,
	Controllers: nil,
	Middlewares: nil,
	Policies:    nil,
	Models:      nil,
	Services:    nil,

	OnRegister: nil,
	OnDestroy:  nil,
}
