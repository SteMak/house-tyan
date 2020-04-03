package awards

import (
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/bwmarrin/discordgo"
)

type module struct {
	session *discordgo.Session
	config  config

	running bool

	unb *unbelievaBoat

	stopHandlers []func()
}

func (module) ID() string {
	return "awards"
}

func (bot module) IsRunning() bool {
	return bot.running
}

func (bot *module) Init(prefix, configPath string) error {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return errors.WithStack(err)
	}

	err = yaml.Unmarshal(data, &bot.config)
	if err != nil {
		return errors.WithStack(err)
	}

	bot.loadEnv()

	bot.unb = &unbelievaBoat{
		token:  bot.config.Bank.Token,
		client: &http.Client{},
	}

	return nil
}

func (bot *module) Start(session *discordgo.Session) {
	bot.session = session
	bot.running = true

	bot.stopHandlers = []func(){
		bot.session.AddHandler(bot.handlerUp),
		bot.session.AddHandler(bot.handlerRequest),
	}
}

func (bot *module) Stop() {
	for _, stopHandler := range bot.stopHandlers {
		stopHandler()
	}
	bot.running = false
}
