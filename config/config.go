package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/SteMak/house-tyan/out"
)

var cfg config

var (
	Session = &cfg.Session
	Bot     = &cfg.Bot
	Storage = &cfg.Storage
	Modules = &cfg.Modules
	Cache   = &cfg.Cache
)

type config struct {
	Session session           `json:"session,omitempty"`
	Bot     bot               `json:"bot,omitempty"`
	Storage storage           `json:"storage,omitempty"`
	Cache   cache             `json:"cache,omitempty"`
	Modules map[string]module `json:"modules,omitempty"`
}

func Load(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		out.Fatal(err)
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		out.Fatal(err)
	}

	Session.Token = os.Getenv("TOKEN")

	out.Infoln("Config loaded:", path)
}
