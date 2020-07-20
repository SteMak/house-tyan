package modules

import (
	"fmt"
	"path/filepath"

	"github.com/SteMak/house-tyan/middleware"

	"github.com/SteMak/house-tyan/app"

	"github.com/SteMak/house-tyan/util"
	"github.com/robfig/cron/v3"

	"github.com/sirupsen/logrus"

	"github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/out"
	"github.com/bwmarrin/discordgo"
)

var (
	modules = make(map[string]Module)
)

//
var (
	Cron *cron.Cron
)

type Module interface {
	Init(prefix string) error

	LoadConfig(string) error
	SetLogger(*logrus.Logger)

	ID() string

	Start(session *discordgo.Session)
	Stop()

	IsRunning() bool
}

func Register(name string, module Module) {
	modules[name] = module
}

func Get(name string) Module {
	if module, ok := modules[name]; ok {
		return module
	}
	return nil
}

func Attach(module Module) {
	module.Start(session)
}

func loadModules() {
	for id, m := range *config.Modules {
		out.Infof("\nLoading %s...\n", id)
		module, exists := modules[id]
		if !exists {
			out.Err(false, fmt.Errorf("Module %s not found", id))
			continue
		}

		out.Infoln("Prefix:", m.Prefix)

		if m.Log != nil {
			log, err := util.Logger(m.Log)
			if err != nil {
				out.Err(false, err)
				continue
			}

			module.SetLogger(log)
		}

		if m.Config != nil {
			configPath := *m.Config
			if !filepath.IsAbs(configPath) {
				configPath = filepath.Join(config.Path, *m.Config)
			}

			out.Infoln("Config file:", configPath)

			if err := module.LoadConfig(configPath); err != nil {
				out.Err(false, err)
				continue
			}
		}
		if err := module.Init(m.Prefix); err != nil {
			out.Err(false, err)
			continue
		}

		if m.Enabled {
			Attach(module)
			out.Infoln("[ENABLED]")
		} else {
			out.Infoln("[DISABLED]")
		}
	}
}

func Run() {
	out.ErrorHandler = SendError

	authentificate()

	logger, err := util.Logger(config.Bot.Log)
	if err == nil {
		app.Use(middleware.Log(logger))
	}
	app.Init(session)

}

func Stop() {
	session.Close()
}
