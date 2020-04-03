package xp

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (bot *module) handlerXpMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if isntRightChannel(m.ChannelID, bot.config.Channels.XpFarming) {
		return
	}
	if len(m.Content) < bot.config.XpConfig.XpMesLen {
		return
	}

	fmt.Println(m.Author.ID, bot.config.XpConfig.XpMessage)
}
