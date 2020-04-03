package xp

import (
	"fmt"

	"github.com/SteMak/house-tyan/util"
	"github.com/bwmarrin/discordgo"
)

func (bot *module) handlerXpMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !util.EqualAny(m.ChannelID, bot.config.MessageFarm.Channels) {
		return
	}
	if len(m.Content) < bot.config.MessageFarm.MessageLength {
		return
	}

	member, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		return
	}
	if util.EqualAny(bot.config.RoleHermit, member.Roles) {
		return 
	}

	fmt.Println(m.Author.ID, bot.config.MessageFarm.XpForMessage)
}
