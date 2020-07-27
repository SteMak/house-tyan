package app

type middleware struct {
	middlewares HandlerChain
}

func (m *middleware) Use(handlers ...HandlerFunc) {
	m.middlewares.append(handlers...)
}

type CommandGroup struct {
	middleware

	alias    string
	enabled  bool
	root     *Command
	commands []*Command
	groups   []*CommandGroup
}

func (g *CommandGroup) On(aliases ...string) *Command {
	cmd := new(Command)
	cmd.owner = g
	cmd.aliases = aliases

	app.commands = append(app.commands, cmd)

	if len(aliases) == 0 {
		g.root = cmd
		return g.root
	}

	g.commands = append(g.commands, cmd)
	return cmd
}

func (g *CommandGroup) Group(alias string) *CommandGroup {
	var group CommandGroup
	group.Enable()
	group.alias = alias
	g.groups = append(g.groups, &group)
	return &group
}

func (g *CommandGroup) Enable() {
	g.enabled = true
}

func (g *CommandGroup) Disable() {
	g.enabled = false
}

func (g *CommandGroup) hasChild() bool {
	return len(g.groups) > 0
}

func matchAliases(args string, cmd *Command) (bool, string) {
	for _, alias := range cmd.aliases {
		match, a := matchCommand(alias, args)
		if match {
			return true, a
		}
	}
	return false, ""
}

func (g *CommandGroup) buildCommands(ctx *Context) {
	if !g.enabled {
		return
	}

	ctx.handlers.append(g.middlewares...)

	for _, group := range g.groups {
		match, args := matchCommand(group.alias, ctx.args.Args)
		if !match {
			continue
		}

		ctx.args.Args = args
		group.buildCommands(ctx)

		return
	}

	for _, cmd := range g.commands {
		if cmd.Handler == nil {
			continue
		}

		match, args := matchAliases(ctx.args.Args, cmd)
		if !match {
			continue
		}

		ctx.args.Args = args
		ctx.handlers.append(cmd.Handler)

		return
	}

	if g.root == nil {
		return
	}

	if g.root.Handler == nil {
		return
	}
	ctx.handlers.append(g.root.Handler)
}
