package clubs

import (
	"io/ioutil"

	"github.com/SteMak/house-tyan/libs/dgutils"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

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

	bot.cmds.Start(session)

	log.Trace("Started.")
}

func (bot *module) Stop() {
	bot.cmds.Stop()

	for _, stopHandler := range bot.stopHandlers {
		stopHandler()
	}
	bot.running = false
	log.Trace("Stoped.")
}
