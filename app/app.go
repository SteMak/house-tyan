package app

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

type (
	HandlerFunc  func(*Context) error
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
	var ctxs []Context
	wg := sync.WaitGroup{}
	for _, module := range a.modules {
		wg.Add(1)
		module.buildCommands(&wg, m.Content, &ctxs, a.middlewares)
	}
	wg.Wait()

}

func Init(s *discordgo.Session) {
	s.AddHandler(app.onHandle)
}

func NewModule(module, prefix string) *Module {
	m := &Module{
		CommandGroup{
			alias: prefix,
		},
	}
	app.modules[module] = m
	return m
}
