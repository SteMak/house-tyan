package storage

import (
	"time"

	"github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/out"
	"github.com/SteMak/house-tyan/util"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/log/logrusadapter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	_ "github.com/jackc/pgx/pgtype"
)

var (
	pgxconn *pgx.Conn
	log     *logrus.Logger
)

// Tables
var (
	Awards  awards
	Rewards rewards
	Users   users
	Clubs   clubs
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

	pgxconn, err = pool.Acquire()
	if err != nil {
		out.Fatal(err)
	}
}

func Tx() (*pgx.Tx, error) {
	return pgxconn.Begin()
}
