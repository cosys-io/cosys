package main

import (
	"database/sql"
	"log"

	"github.com/cosys-io/cosys/internal/cosys"
)

func main() {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}

	_ = &cosys.Cosys{
		DB: db,
	}
}
