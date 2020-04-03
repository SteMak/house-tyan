package config

import (
	"io/ioutil"
	"os"

	"github.com/SteMak/house-tyan/out"
	"gopkg.in/yaml.v2"
)

var cfg config

var (
	Session = &cfg.Session
	Bot     = &cfg.Bot
	Storage = &cfg.Storage
	Modules = &cfg.Modules
)

type config struct {
	Session session           `yaml:"session,omitempty"`
	Bot     bot               `yaml:"bot,omitempty"`
	Storage storage           `yaml:"storage,omitempty"`
	Modules map[string]module `yaml:"modules,omitempty"`
}

func Load(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		out.Fatal(err)
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		out.Fatal(err)
	}

	Session.Token = os.Getenv("TOKEN")

	out.Infoln("Config loaded:", path)
}
