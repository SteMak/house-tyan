package storage

import (
	"database/sql"

	"github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/out"

	// postgres driver
	_ "github.com/lib/pq"
)

var db *sql.DB

func Init() {
	out.Infoln("\nInit storage...")

	var err error
	db, err = sql.Open(config.Storage.Driver, config.Storage.Connection)
	if err != nil {
		out.Infoln("[FAIL]")
		out.Fatal(err)
	}
	out.Infoln("Done.")
}

func Tx() (*sql.Tx, error) {
	tx, err := db.Begin()
	return tx, err
}
