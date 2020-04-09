package awards

import (
	"strings"

	"github.com/SteMak/house-tyan/modules"
	"github.com/pkg/errors"

	"github.com/SteMak/house-tyan/cache"
	"github.com/SteMak/house-tyan/out"
	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/badger"
)

func (bot *module) handlerUp(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !bot.matchUp(s, m.Message) {
		return
	}

	out.Debugln("FOUND S.up")

	user, err := cache.Usernames.Get(m.Author.String())
	if errors.Is(err, badger.ErrKeyNotFound) {
		// TODO: обработать ситуацию с не найденным пользователем
		return
	}
	if err != nil {
		out.Err(true, errors.WithStack(err))
		return
	}

	err = bot.unb.AddToBalance(user.ID, int64(bot.config.AwardAmount), "for S.up")
	if err != nil {
		out.Err(true, errors.WithStack(err))
		modules.Send(m.ChannelID, "awards/bump_fail.xml", map[string]interface{}{
			"Mention": m.Author.Mention(),
		}, nil)
		return
	}

	go modules.SendLog("awards/awarded.xml", map[string]interface{}{
		"Amount":  bot.config.AwardAmount,
		"Mention": m.Author.Mention(),
		"Reason":  "S.up",
	})

	// tpl, err := messages.Random(`^awards/answers/*\.xml$`)
	// if err != nil {
	// 	out.Err(true, errors.WithStack(err))
	// 	return
	// }

	// go modules.Send(m.ChannelID, tpl, map[string]interface{}{
	// 	"Amount":  bot.config.AwardAmount,
	// 	"Mention": m.Author.Mention(),
	// 	"Reason":  "S.up",
	// }, nil)
}

func (bot *module) handlerBlankAccept(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.ChannelID != bot.config.Channels.Confirm {
		return
	}

	emoji := m.Emoji.Name
	if !strings.ContainsAny(emoji, "✅🇽") {
		return
	}

	award, exists, err := cache.Awards.Get(m.MessageID)
	if !exists {
		return
	}
	if err != nil {
		out.Err(true, errors.WithStack(err))
		return
	}

	switch emoji {
	case "✅":
		for _, reward := range award.Rewards {
			for _, user := range reward.Users {
				if err := bot.unb.AddToBalance(user.ID, int64(reward.Amount), award.Reason); err != nil {
					out.Err(true, errors.WithStack(err))
					return
				}
			}
		}

		cache.Blanks.Delete(award.Author.ID)
		cache.Awards.Delete(award.ID)

	case "🇽":
		cache.Blanks.Delete(award.Author.ID)
		cache.Awards.Delete(award.ID)
	}
}
