package awards

import (
	"encoding/json"
	"io/ioutil"

	"github.com/bwmarrin/discordgo"
)

type module struct {
	session *discordgo.Session
	config  config

	running bool
}

func (module) ID() string {
	return "awards"
}

func (bot *module) Init(prefix, configPath string) error {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &bot.config)
	if err != nil {
		return err
	}

	bot.app = &router.App{
		Prefix:      prefix,
		Description: bot.config.Description,
	}

	bot.initCommands()

	return nil
}

func (bot *module) Start(session *discordgo.Session) {
	bot.session = session
	bot.running = true

	bot.session.AddHandler(bot.onMessageCreate)
	bot.session.AddHandler(bot.onReactionAdd)
}

func (bot *module) Stop() {
	bot.running = false
}
