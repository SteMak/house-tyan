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

	xp := howMuchXp(m.Content, bot.config.MessageFarm)
	err = storage.Users.AddXP(tx, m.Author, int64(xp))
	if err != nil {
		out.Err(true, errors.WithStack(err))
		tx.Rollback()
		return
	}

	tx.Commit()
	bot.handleXp(m.Author.ID, uint64(xp))
}
