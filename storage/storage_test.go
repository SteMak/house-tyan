package storage

import (
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v5"

	"github.com/SteMak/house-tyan/config"
)

func TestMain(t *testing.M) {
	gofakeit.Seed(time.Now().UnixNano())

	config.Storage.Driver = "pgx"
	config.Storage.Connection = "postgresql://bot:password@localhost:54320/anihouse?sslmode=disable"

	Init()

	os.Exit(t.Run())
}
