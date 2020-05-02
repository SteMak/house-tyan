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

	config.Load("./../cli/bot/config/dev/config.yaml")
	Init()

	os.Exit(t.Run())
}
