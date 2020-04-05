package xp

import (
	"regexp"

	"github.com/SteMak/house-tyan/out"

	"github.com/SteMak/house-tyan/util"
	"github.com/bwmarrin/discordgo"
)

func (bot *module) handlerXpMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !util.EqualAny(m.ChannelID, bot.config.MessageFarm.Channels) {
		return
	}

	member, err := s.State.Member(m.GuildID, m.Author.ID)
	if err != nil || member.User.Bot {
		return
	}
	if util.EqualAny(bot.config.RoleHermit, member.Roles) {
		return
	}

	if len(m.Content) == 0 {
		out.Debugln(m.Author.ID, bot.config.MessageFarm.XpForMessage)
		return
	}

	var (
		content         = m.Content
		foundings       [][]byte
		lengthOfSpecial = 0
		lengthOfCommon  = 0
	)

	link := regexp.MustCompile(`<@!?\d+>`)
	foundings = link.FindAll([]byte(content), -1)
	content = string(link.ReplaceAll([]byte(content), []byte("")))
	lengthOfSpecial += len(foundings)

	emodji := regexp.MustCompile(`<a?:\w+:\d+>`)
	foundings = emodji.FindAll([]byte(content), -1)
	content = string(emodji.ReplaceAll([]byte(content), []byte("")))
	lengthOfSpecial += len(foundings)

	runes := []rune(content)
	for _, r := range runes {
		if len(string(r)) != 1 {
			lengthOfSpecial++
		} else {
			lengthOfCommon++
		}
	}

	out.Debugln(m.Author.ID, int(bot.config.MessageFarm.XpForChar*float32(lengthOfCommon))+lengthOfSpecial)
}
