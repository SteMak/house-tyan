package dgutils

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type Menu struct {
	Message  *discordgo.MessageSend
	Items    map[string]func(*ReactionContext)
	Duration time.Duration
}

func (discord *Discord) menuHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}

	if menu := discord.menus[r.MessageID]; menu != nil {
		if f, ok := menu.Items[r.Emoji.Name]; ok {
			delete(discord.menus, r.MessageID)
			if err := s.MessageReactionsRemoveAll(r.ChannelID, r.MessageID); err != nil {
				discord.ErrorHandler(err)
			}
			f(&ReactionContext{
				Session:  discord.Session,
				Reaction: r.MessageReaction,
			})
		}
	}
}

func (discord *Discord) SendMenu(channelID string, m *Menu) (*discordgo.Message, error) {
	msg, err := discord.Session.ChannelMessageSendComplex(channelID, m.Message)
	if err != nil {
		return nil, err
	}

	for id := range m.Items {
		err = discord.Session.MessageReactionAdd(channelID, msg.ID, id)
		if err != nil {
			return nil, err
		}
	}

	discord.menus[msg.ID] = m

	go func() {
		time.Sleep(m.Duration)
		if _, ok := discord.menus[msg.ID]; ok {
			delete(discord.menus, msg.ID)
			if err := discord.Session.MessageReactionsRemoveAll(channelID, msg.ID); err != nil {
				discord.ErrorHandler(err)
			}
		}
	}()

	return msg, nil
}
