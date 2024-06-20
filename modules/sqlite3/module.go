package sqlite3

import (
	"database/sql"
	"github.com/cosys-io/cosys/common"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

func init() {
	dbFunc := func(cosys *common.Cosys) common.Database {
		return Database{
			cosys: cosys,
			db:    db,
		}
	}

	if err := common.RegisterDatabase("sqlite3", dbFunc); err != nil {
		log.Fatal(err)
	}
}

func OnRegister(cosys common.Cosys) (common.Cosys, error) {
	var err error
	db, err = sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := loadSchema(cosys); err != nil {
		return cosys, err
	}

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
