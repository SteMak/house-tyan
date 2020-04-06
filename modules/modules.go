package modules

import (
	"fmt"

	"github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/out"
	"github.com/bwmarrin/discordgo"
)

type Module interface {
	Init(prefix, configPath string) error

	ID() string

	Start(session *discordgo.Session)
	Stop()

	IsRunning() bool
}

var (
	modules = make(map[string]Module)
)

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

		out.Infoln("Config file:", m.Config)
		out.Infoln("Prefix:", m.Prefix)

		if err := module.Init(m.Prefix, m.Config); err != nil {
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
}

func Stop() {
	session.Close()
}
