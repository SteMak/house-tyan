package xp

import (
	"github.com/SteMak/house-tyan/out"

	"github.com/bwmarrin/discordgo"
)

func (bot *module) handlerXpMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if xpMessageChecks(m.ChannelID, bot.config, m.GuildID, m.Author.ID, s.State.Member) != nil {
		return
	}

	out.Debugln(m.Author, howMuchXp(m.Content, bot.config.MessageFarm))
}
