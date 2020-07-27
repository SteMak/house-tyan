package app

import (
	"math"

	"github.com/bwmarrin/discordgo"
)

const abortIndex int8 = math.MaxInt8 / 2

type Context struct {
	args
	module   *Module
	params   map[string]interface{}
	handlers HandlerChain
	index    int8

	Session *discordgo.Session
	Message *discordgo.Message
	Error   error
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		if c.index == abortIndex {
			break
		}

		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Context) Abort() {
	c.index = abortIndex
}

func (c *Context) AboutWithError(err error) {
	c.index = abortIndex
	c.Error = err
}
