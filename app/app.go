package app

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

type (
	HandlerFunc  func(*Context)
	HandlerChain []HandlerFunc
)

func (chain *HandlerChain) append(handlers ...HandlerFunc) {
	*chain = append(*chain, handlers...)
}

var (
	app *Application = &Application{}
)

func init() {
	app.modules = make(map[string]*Module)
}

type Application struct {
	middleware

	modules  map[string]*Module
	session  *discordgo.Session
	commands []*Command
}

func (a *Application) onHandle(s *discordgo.Session, m *discordgo.MessageCreate) {
	var ctxs []*Context
	wg := sync.WaitGroup{}
	for _, module := range a.modules {
		wg.Add(1)
		ctx := module.buildCommands(&wg, s, m.Message)
		if ctx != nil {
			ctxs = append(ctxs, ctx)
		}
	}
	wg.Wait()

	for _, ctx := range ctxs {
		ctx.index = -1
		ctx.Next()
	}
}

func Init(s *discordgo.Session) {
	s.AddHandler(app.onHandle)
}

func NewModule(module, prefix string) *Module {
	m := &Module{
		CommandGroup: CommandGroup{
			alias: prefix,
		},
		Prefix: prefix,
		Name:   module,
	}
	app.modules[module] = m
	return m
}

func Use(handlers ...HandlerFunc) {
	app.Use(handlers...)
}
