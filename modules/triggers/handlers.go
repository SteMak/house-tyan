package triggers

import (
	"github.com/bwmarrin/discordgo"
)

func (bot *module) triggerHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || m.Content == "" {
		return
	}

}
