package clubs

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"

	"github.com/SteMak/house-tyan/libs/dgutils"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/bwmarrin/discordgo"
)

var log *logrus.Logger

type module struct {
	session *discordgo.Session
	config  config

	running bool

	cmds         *dgutils.Discord
	stopHandlers []func()
}

func (module) ID() string {
	return "clubs"
}

func (bot module) IsRunning() bool {
	return bot.running
}

func (bot *module) LoadConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.WithStack(err)
	}

	err = yaml.Unmarshal(data, &bot.config)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (bot *module) SetLogger(logger *logrus.Logger) {
	log = logger
}

func (bot *module) Init(prefix string) error {
	bot.loadEnv()

	bot.cmds = &dgutils.Discord{
		Prefix:   prefix,
		Commands: commands,
	}

	return nil
}

func (bot *module) Start(session *discordgo.Session) {
	bot.session = session
	bot.running = true

	bot.stopHandlers = []func(){}

}

func (bot *module) Stop() {
	bot.cmds.Stop()

	for _, stopHandler := range bot.stopHandlers {
		stopHandler()
	}
	bot.running = false
}
