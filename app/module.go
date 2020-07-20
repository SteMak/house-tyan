package app

import "sync"

type Module struct {
	CommandGroup
}

func (m *Module) buildCommands(wg *sync.WaitGroup, args string, ctxs *[]Context, chain HandlerChain) {
	m.CommandGroup.buildCommands(wg, args, ctxs, chain)
}
