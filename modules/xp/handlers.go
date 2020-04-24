package xp

import (
	"github.com/SteMak/house-tyan/out"
	"github.com/SteMak/house-tyan/storage"
	"github.com/pkg/errors"

	"github.com/bwmarrin/discordgo"
)

func (bot *module) handlerXpMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if xpMessageChecks(m.ChannelID, bot.config, m.GuildID, m.Author.ID, s.State.Member) != nil {
		return
	}

	tx, err := storage.Tx()
	if err != nil {
		out.Err(true, errors.WithStack(err))
		return
	}

	err = storage.Users.AddXP(tx, m.Author, int64(howMuchXp(m.Content, bot.config.MessageFarm)))
	if err != nil {
		out.Err(true, errors.WithStack(err))
		tx.Rollback()
		return
	}

	tx.Commit()
}
