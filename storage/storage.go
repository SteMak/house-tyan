package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/out"
	"github.com/jmoiron/sqlx"

	// postgres driver
	_ "github.com/jackc/pgx/v4/stdlib"
)

var (
	db *sqlx.DB

	psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

var (
	Awards  awards
	Rewards rewards
	Users   users
)

func Init() {
	out.Infoln("\nInit storage...")

	var err error
	db, err = sqlx.Connect(config.Storage.Driver, config.Storage.Connection)
	if err != nil {
		out.Infoln("[FAIL]")
		out.Fatal(err)
	}

	db.SetMaxIdleConns(config.Storage.MaxIdleConnection)
	db.SetMaxOpenConns(config.Storage.MaxOpenConnection)

	out.Infoln("Done.")
}

func Tx() (*sqlx.Tx, error) {
	return db.Beginx()
}

func exec(tx *sqlx.Tx, stmt squirrel.Sqlizer) error {
	query, args, err := stmt.ToSql()
	if err != nil {
		return err
	}

	if tx == nil {
		_, err = db.Exec(query, args...)
	} else {
		_, err = tx.Exec(query, args...)
	}
	return err
}
