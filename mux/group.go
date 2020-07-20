package mux

type Group struct {
	middlewares HandlerChain
	
}

func (g *Group) Use(handlers ...HandlerFunc) {
	g.middlewares = append(g.middlewares, handlers...)
}
