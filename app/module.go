package app

import (
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type Module struct {
	CommandGroup

	Prefix string
	Name   string
}

func (m *Module) buildCommands(wg *sync.WaitGroup, s *discordgo.Session, mess *discordgo.Message) *Context {
	defer wg.Done()

	if !strings.HasPrefix(mess.Content, m.alias) {
		return nil
	}

	ctx := &Context{
		module:   m,
		handlers: m.middlewares,
		params:   make(map[string]interface{}),
		Message:  mess,
		Session:  s,
	}
	ctx.args.Args = strings.TrimPrefix(mess.Content, m.alias)

	m.CommandGroup.buildCommands(ctx)

	return ctx
}
