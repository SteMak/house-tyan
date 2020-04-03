package xp

import (
	"io/ioutil"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v2"
)

type module struct {
	session *discordgo.Session
	config  config

	running bool

	stopHandlers []func()
}

func (module) ID() string {
	return "xp"
}

func (bot module) IsRunning() bool {
	return bot.running
}

func (bot *module) Init(prefix, configPath string) error {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &bot.config)
	if err != nil {
		return err
	}

	return nil
}

func (bot *module) Start(session *discordgo.Session) {
	bot.session = session
	bot.running = true

	// TODO runMagicGorooting()

	bot.stopHandlers = []func(){
		bot.session.AddHandler(bot.handlerXpMessage),
	}
}

func (bot *module) Stop() {
	for _, stopHandler := range bot.stopHandlers {
		stopHandler()
	}
	bot.running = false
}
