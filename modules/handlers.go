package modules

import (
	"github.com/SteMak/house-tyan/config"
	"github.com/bwmarrin/discordgo"
)

func onReady(s *discordgo.Session, e *discordgo.Ready) {
	if config.Bot.LogChannel == nil {
		return
	}

	data := map[string]interface{}{
		"Name": e.User,
	}

	Send(*config.Bot.LogChannel, "main/started.xml", data, nil)
}
