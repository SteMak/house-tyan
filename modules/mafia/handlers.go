package mafia

import (
	"strconv"
	"sync"

	"github.com/SteMak/house-tyan/out"
	"github.com/bwmarrin/discordgo"
)

var (
	joinMsg *discordgo.Message
	infoMsg *discordgo.Message

	votePlayers  sync.Map
	votedPlayers sync.Map
	voteMsg      *discordgo.Message
)

func (bot *module) mafiaJoinHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if joinMsg == nil || r.MessageID != joinMsg.ID || r.Emoji.Name != "✅" || r.UserID == s.State.User.ID {
		return
	}

	bot.game.Add(r.UserID)
	updateInfoEmbed(bot)
}

func (bot *module) mafiaLeaveHandler(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	if joinMsg == nil || r.MessageID != joinMsg.ID || r.Emoji.Name != "✅" || r.UserID == s.State.User.ID {
		return
	}

	bot.game.Pop(r.UserID)
	updateInfoEmbed(bot)
}

func (bot *module) mafiaVoteHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if voteMsg == nil || m.Message.ChannelID != voteMsg.ChannelID {
		return
	}

	if _, player := bot.game.GetPlayer(m.Author.ID); player != nil {
		if _, ok := votedPlayers.Load(player); ok {
			return
		}

		index, err := strconv.Atoi(m.Content)
		if err != nil {
			return
		}

		if _, target := bot.game.GetPlayerByIndex(index); target != nil {
			if _, ok := votePlayers.Load(target); ok {
				votedPlayers.Store(player, true)
				if votes, ok := votePlayers.Load(target); ok {
					votePlayers.Store(target, votes.(int)+1)
				}
				updateVoteEmbed(bot)
			}
			return
		}

		_, err = s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" такого игрока нет")
		if err != nil {
			out.Err(true, err)
		}
	}
}
