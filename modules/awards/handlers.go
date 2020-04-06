package awards

import (
	"sort"

	"github.com/SteMak/house-tyan/messages"
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

	err = bot.unb.addToBalance(user.ID, 0, bot.config.AwardAmount, "for S.up")
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

	tpl, err := messages.Random(`^awards/answers/*\.xml$`)
	if err != nil {
		out.Err(true, errors.WithStack(err))
		return
	}

	go modules.Send(m.ChannelID, tpl, map[string]interface{}{
		"Amount":  bot.config.AwardAmount,
		"Mention": m.Author.Mention(),
		"Reason":  "S.up",
	}, nil)
}

func (bot *module) handlerRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	ok, content := bot.matchRequest(s, m.Message)
	if !ok {
		return
	}

	if m.Content == "-запрос" {
		modules.Send(m.ChannelID, "awards/usage.xml", nil, nil)
		return
	}

	item, err := parseRequest(s, content)
	if err != nil {
		out.Err(true, err)
		return
	}

	sort.SliceStable(item.Rewards, func(i, j int) bool {
		return item.Rewards[i].Amount > item.Rewards[j].Amount
	})

	blank := modules.Send(bot.config.Channels.Confirm, "awards/blank.xml", map[string]interface{}{
		"Reason":  item.Reason,
		"Rewards": item.Rewards,
	}, func(msg *messages.Message) error {
		if m.Author == nil {
			return errors.New("Message author is nil")
		}
		msg.Embed.Author = &discordgo.MessageEmbedAuthor{
			Name:    m.Author.String(),
			IconURL: m.Author.AvatarURL("16"),
		}
		return nil
	})

	if blank == nil {
		return
	}

	item.ID = blank.ID
	item.Author = *m.Author

	err = cache.Awards.Set(item)
	if err != nil {
		out.Err(true, errors.WithStack(err))
		return
	}
}

func (bot *module) handleCommand(s *discordgo.Session, m *discordgo.MessageCreate) {

}
