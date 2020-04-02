package awards

import (
	"regexp"
	"strings"

	"github.com/SteMak/house-tyan/util"

	"github.com/pkg/errors"

	conf "github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/out"
	"github.com/bwmarrin/discordgo"
)

func (bot *module) matchUp(s *discordgo.Session, m *discordgo.Message) bool {
	if m.ChannelID != bot.config.Channels.Bump {
		return false
	}
	if m.Author == nil {
		return false
	}
	if m.Author.ID != bot.config.Bots.Uper {
		return false
	}

	if len(m.Embeds) != 1 {
		return false
	}

	embed := m.Embeds[0]

	if embed.Title != "Сервер Up" {
		return false
	}
	if embed.Footer == nil {
		return false
	}
	return true
}

func (bot *module) matchBump(s *discordgo.Session, m *discordgo.Message) bool {
	if m.ChannelID != bot.config.Channels.Bump {
		return false
	}
	if m.Author == nil {
		return false
	}
	if m.Author.ID != bot.config.Bots.Bumper {
		return false
	}

	if len(m.Embeds) != 1 {
		return false
	}

	embed := m.Embeds[0]

	return regexp.
		MustCompile(`Server bumped by <@\d+>`).
		MatchString(embed.Description)
}

func (bot *module) matchRequest(s *discordgo.Session, m *discordgo.Message) (bool, string) {
	if m.GuildID != conf.Bot.GuildID {
		return false, ""
	}

	content := m.Content

	var request string
	isRequest := func() bool {
		if strings.HasPrefix(content, "-запрос") {
			request = strings.TrimPrefix(content, "-запрос")
			return true
		}
		if strings.HasPrefix(content, "- запрос") {
			request = strings.TrimPrefix(content, "- запрос")
			return true
		}
		if strings.HasPrefix(content, "-  запрос") {
			request = strings.TrimPrefix(content, "-  запрос")
			return true
		}
		return false
	}

	if !isRequest() {
		return false, ""
	}

	member, err := s.GuildMember(conf.Bot.GuildID, m.Author.ID)
	if err != nil {
		out.Err(true, errors.WithStack(err))
		return false, ""
	}

	if !util.HasRole(member, bot.config.Roles.Requester) {
		return false, ""
	}

	return true, request
}
