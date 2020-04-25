package storage

import (
	"time"

	"github.com/pkg/errors"

	"github.com/SteMak/house-tyan/util"

	"github.com/jackc/pgx/log/logrusadapter"

	"github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/out"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var (
	db  *sqlx.DB
	log *logrus.Logger

	psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

var (
	Awards  awards
	Rewards rewards
	Users   users
)

func Init() {
	out.Infoln("\nInit storage...")

	connCfg, err := pgx.ParseURI(config.Storage.Connection)
	if err != nil {
		out.Fatal(err)
	}

	log, err = util.Logger(config.Storage.Log)
	if err != nil {
		if !errors.Is(err, util.ErrNoLogger) {
			out.Fatal(err)
		}
	}

	if log != nil {
		connCfg.Logger = logrusadapter.NewLogger(log)
	}
	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connCfg,
		MaxConnections: 20,
		AcquireTimeout: 30 * time.Second,
	})
	if err != nil {
		out.Fatal(err)
	}

	native := stdlib.OpenDBFromPool(pool)
	db = sqlx.NewDb(native, config.Storage.Driver)

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
