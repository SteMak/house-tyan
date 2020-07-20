package app

import (
	"github.com/bwmarrin/discordgo"
)

type Context struct {
	args     string
	module   *Module
	params   map[string]interface{}
	handlers HandlerChain
	index    uint8
	Session  *discordgo.Session
	Message  *discordgo.Message
}
