package dgutils

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Description string

	Raw      bool
	Function func(*MessageContext)

	Handlers []func(*MessageContext)
}

func (c *Command) Use(middleware ...func(*MessageContext)) {
	c.Handlers = append(c.Handlers, middleware...)
}

type Group struct {
	Description string
	Commands    map[string]interface{}
}

func commandsToString(commands map[string]interface{}) string {
	var usage string
	for key, value := range commands {
		switch v := value.(type) {
		case *Command:
			if v.Description != "" {
				usage += fmt.Sprintf("`%s` - %s\n", key, v.Description)
			} else {
				usage += fmt.Sprintf("`%s`\n", key)
			}
		case *Group:
			if v.Description != "" {
				usage += fmt.Sprintf("\n`[%s]` - %s\n", key, v.Description)
			} else {
				usage += fmt.Sprintf("\n`[%s]`\n", key)
			}
			usage += commandsToString(v.Commands) + "\n"
		}
	}
	return usage
}

func Usage(commands map[string]interface{}) *discordgo.MessageSend {
	usageEmbed := discordgo.MessageEmbed{
		Title:       "Использование",
		Description: commandsToString(commands),
		Color:       4437377,
	}

	return &discordgo.MessageSend{
		Embed: &usageEmbed,
	}
}

func findCommand(args string, commands map[string]interface{}) (*Command, []string) {
	an := strings.SplitN(args, " ", 2)

	if item, ok := commands[an[0]]; ok {
		switch v := item.(type) {
		case *Command:
			if len(an) == 2 {
				if v.Raw {
					return v, []string{an[1]}
				}
				return v, strings.Fields(an[1])
			}
			return v, nil
		case *Group:
			if len(an) == 2 {
				return findCommand(an[1], v.Commands)
			}
			return nil, nil
		}
	}

	return nil, nil
}
