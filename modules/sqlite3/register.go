package sqlite3

import (
	"database/sql"
	"github.com/cosys-io/cosys/common"

	_ "github.com/mattn/go-sqlite3"
)

func Register(cosys common.Cosys) (common.Cosys, error) {
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		return cosys, err
	}

	cosys.Db = func(cosys *common.Cosys) common.Database {
		return Database{
			cosys: cosys,
			db:    db,
		}
	}

	return cosys, nil
}
