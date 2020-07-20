package clubs

import (
	"io/ioutil"

	"github.com/SteMak/house-tyan/app"

	"github.com/SteMak/house-tyan/modules"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type module struct {
	session *discordgo.Session
	config  config

	running bool

	cmd          *app.Module
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
	bot.cmd = app.NewModule(_module.ID(), prefix)
	group := bot.cmd.Group("club")
	group.On("create").Handle(bot.onClubCreate)
	group.On("delete").Handle(bot.onClubDelete)
	group.On("kick").Handle(bot.onClubKick)
	return nil
}

func (bot *module) Start(session *discordgo.Session) {
	bot.session = session
	bot.running = true

	bot.stopHandlers = []func(){}

	bot.cmd.Enable()

	modules.Cron.AddFunc("@daily", bot.removeNotVerified)

	log.Trace("Started.")
}

func (bot *module) Stop() {
	bot.cmd.Disable()

	for _, stopHandler := range bot.stopHandlers {
		stopHandler()
	}
	bot.running = false
	log.Trace("Stoped.")
}
