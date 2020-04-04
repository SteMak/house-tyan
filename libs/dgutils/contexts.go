package dgutils

import (
	"github.com/bwmarrin/discordgo"
)

type MessageContext struct {
	Session *discordgo.Session
	Message *discordgo.Message
	Args    []string

	handlers []func(*MessageContext)
	index    int
}

func (mc *MessageContext) Next() {
	mc.index++
	if len(mc.handlers) >= mc.index {
		mc.handlers[mc.index-1](mc)
	}
}

type ReactionContext struct {
	Session  *discordgo.Session
	Reaction *discordgo.MessageReaction
}
