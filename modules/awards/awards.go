package awards

import (
	"io/ioutil"

	"github.com/SteMak/house-tyan/libs"
	"github.com/SteMak/house-tyan/libs/dgutils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/bwmarrin/discordgo"
)

var (
	log *logrus.Logger
)

type module struct {
	session *discordgo.Session
	config  config

	running bool

	cmds *dgutils.Discord
	unb  *libs.UnbelievaBoatAPI

	stopHandlers []func()
}

func (module) ID() string {
	return "awards"
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

func (bot module) IsRunning() bool {
	return bot.running
}

func (bot *module) Init(prefix string) error {
	if log == nil {
		return errors.New("Logger is required")
	}

	bot.loadEnv()

	bot.cmds = &dgutils.Discord{
		Prefix:   prefix,
		Commands: commands,
	}

	bot.unb = libs.NewUnbelievaBoatAPI(bot.config.Bank.Token)

	return nil
}

func (bot *module) Start(session *discordgo.Session) {
	bot.session = session
	bot.running = true

	bot.stopHandlers = []func(){
		bot.session.AddHandler(bot.handlerBlankProcess),
	}

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
