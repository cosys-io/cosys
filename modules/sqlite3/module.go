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
	dbCtor := func(cosys *common.Cosys) common.Database {
		return Database{
			cosys: cosys,
			db:    db,
		}
	}

	if err := common.RegisterDatabase("sqlite3", dbCtor); err != nil {
		log.Fatal(err)
	}

	if err := common.RegisterModule("sqlite3", module); err != nil {
		log.Fatal(err)
	}
}

func OnRegister(cosys common.Cosys) (common.Cosys, error) {
	var err error
	db, err = sql.Open("sqlite3", "data.db")
	if err != nil {
		return common.Cosys{}, err
	}

	if err = loadSchema(cosys); err != nil {
		return common.Cosys{}, err
	}

	return cosys, nil
}

var module = &common.Module{
	OnRegister: OnRegister,
	OnDestroy:  nil,
}
