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
		out.Debugln(m.Author.ID, int(bot.config.MessageFarm.XpForEmpty))
		return
	}

	var (
		content       = m.Content
		foundings     []string
		countOfRunes  = 0
		countOfCommon = 0
	)

	link := regexp.MustCompile(`<@!?\d+>`)
	foundings = link.FindAllString(content, -1)
	content = link.ReplaceAllLiteralString(content, "")
	countOfRunes += len(foundings)

	emodji := regexp.MustCompile(`<a?:\w+:\d+>`)
	foundings = emodji.FindAllString(content, -1)
	content = emodji.ReplaceAllLiteralString(content, "")
	countOfRunes += len(foundings)

	runes := []rune(content)
	for _, r := range runes {
		if len(string(r)) != 1 {
			countOfRunes++
		} else {
			countOfCommon++
		}
	}

	charXP := bot.config.MessageFarm.XpForChar * float32(countOfCommon)
	runeXP := bot.config.MessageFarm.XpForRune * float32(countOfRunes)
	out.Debugln(m.Author.ID, int(charXP)+int(runeXP))
}
