package app

type Command struct {
	aliases []string

	owner   *CommandGroup
	Args    string
	Handler HandlerFunc
}

func (c *Command) Handle(handler HandlerFunc) *Command {
	c.Handler = handler
	return c
}
