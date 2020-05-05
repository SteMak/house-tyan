package xp

import (
	"io/ioutil"

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

	stopHandlers []func()

	voiceXpWorker *voiceXpWorker

	xpHandlers []modules.HandlerXP
}

func (module) ID() string {
	return "xp"
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

}

func (bot *module) Init(prefix string) error {
	bot.voiceXpWorker = newVoiceXpWorker(&bot.config)

	return nil
}

func (bot *module) Start(session *discordgo.Session) {
	bot.session = session
	bot.running = true

	bot.stopHandlers = []func(){
		bot.session.AddHandler(bot.handlerXpMessage),
	}
	bot.voiceXpWorker.start(session)
}

func (bot *module) Stop() {
	for _, stopHandler := range bot.stopHandlers {
		stopHandler()
	}
	bot.voiceXpWorker.stop()

	bot.running = false
}
