package storage

import (
	"database/sql"

	"github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/out"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"

	// postgres driver
	_ "github.com/lib/pq"
)

var db *sql.DB

func Init() {
	out.Info("\nstorage connection         ")

	var err error
	db, err = sql.Open(config.Storage.Driver, config.Storage.Connection)
	if err != nil {
		out.Infoln("[FAIL]")
		out.Fatal(err)
	}

	if config.Storage.Driver == "sqlite3" {
		db.SetMaxOpenConns(1)
	}
	out.Infoln("[OK]")
}
