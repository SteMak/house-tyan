package triggers

import (
	"github.com/SteMak/house-tyan/libs/dgutils"
	"github.com/SteMak/house-tyan/out"
	"github.com/bwmarrin/discordgo"
)

type module struct {
	session *discordgo.Session
	cmds    *dgutils.Discord

	running bool

	stopHandlers []func()
}

func (module) ID() string {
	return "triggers"
}

func (bot module) IsRunning() bool {
	return bot.running
}

func (bot *module) Init(prefix, configPath string) error {
	// data, err := ioutil.ReadFile(configPath)
	// if err != nil {
	// 	return err
	// }

	// err = yaml.Unmarshal(data, &bot.config)
	// if err != nil {
	// 	return err
	// }

	bot.cmds = &dgutils.Discord{
		Prefix:       prefix,
		Commands:     commands,
		ErrorHandler: func(err error) { out.Err(true, err) },
	}

	return nil
}

func (bot *module) Start(session *discordgo.Session) {
	bot.session = session
	bot.running = true

	bot.stopHandlers = []func(){
		bot.session.AddHandler(bot.triggerHandler),
	}

	bot.cmds.Start(session)
}

func (bot *module) Stop() {
	bot.cmds.Stop()

	for _, stopHandler := range bot.stopHandlers {
		stopHandler()
	}

	bot.running = false
}
