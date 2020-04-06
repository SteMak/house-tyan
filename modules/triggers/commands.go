package triggers

import (
	"strings"

	"github.com/SteMak/house-tyan/cache"
	"github.com/SteMak/house-tyan/libs/dgutils"
	"github.com/SteMak/house-tyan/modules"
	"github.com/SteMak/house-tyan/out"
)

var (
	commands = map[string]interface{}{
		"триггер": &dgutils.Group{
			Commands: map[string]interface{}{
				"добавить": &dgutils.Command{
					Raw:         true,
					Description: "Добавить триггер",
					Function:    _module.onTriggerAdd,
				},
			},
		},
	}
)

func (bot *module) onTriggerAdd(ctx *dgutils.MessageContext) {
	args := strings.SplitN(ctx.Args[0], " ", 2)

	var trigger *cache.Trigger

	trigger, err := cache.Triggers.Get(args[0])
	if err != nil {
		trigger = &cache.Trigger{
			Name:    args[0],
			Answers: []string{args[1]},
		}
	} else {
		trigger = &cache.Trigger{
			Name:    args[0],
			Answers: append(trigger.Answers, args[1]),
		}
	}

	err = cache.Triggers.Set(trigger)
	if err != nil {
		out.Err(true, err)
		return
	}

	modules.Send(ctx.Message.ChannelID, "triggers/trigger.added.xml", trigger, nil)
}
