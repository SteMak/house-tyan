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

	params map[string]interface{}
}

func newContext(s *discordgo.Session, m *discordgo.Message, args []string, handlers []func(*MessageContext)) *MessageContext {
	return &MessageContext{
		Session:  s,
		Message:  m,
		Args:     args,
		handlers: handlers,
		params:   make(map[string]interface{}),
	}
}

func (mc *MessageContext) SetParam(key string, value interface{}) {
	mc.params[key] = value
}

func (mc *MessageContext) Param(key string) interface{} {
	return mc.params[key]
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
