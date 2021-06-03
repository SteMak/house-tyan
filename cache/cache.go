package cache

import (
	"log"
	"os"

	"github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/out"
	"github.com/dgraph-io/badger"
)

var cache *badger.DB

var (
	Awards    awards
	Blanks    blanks
	Usernames usernames
	Triggers  triggers
	Voices    voices
)

func Init() {
	var err error

	if err := os.MkdirAll(config.Cache.Path, 0775); err != nil {
		log.Fatal("error init cache:", err)
	}

	cache, err = badger.Open(badger.DefaultOptions(config.Cache.Path))

	if err != nil {
		out.Fatal("error init cache:", err)
	}
}

func Close() {
	cache.Close()
}
