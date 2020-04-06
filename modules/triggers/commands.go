package triggers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/dgraph-io/badger"

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
					Raw:      true,
					Function: _module.onTriggerAdd,
				},
				"удалить": &dgutils.Command{
					Function: _module.onTriggerDelete,
				},
			},
		},
	}
)

func (bot *module) onTriggerAdd(ctx *dgutils.MessageContext) {
	if ctx.Args == nil || !strings.Contains(ctx.Args[0], " ") {
		modules.Send(ctx.Message.ChannelID, "triggers/usage.xml", nil, nil)
		return
	}

	args := strings.SplitN(ctx.Args[0], " ", 2)
	args[0] = strings.ToLower(args[0])

	trigger, err := cache.Triggers.Get(args[0])
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			trigger = &cache.Trigger{
				Name:    args[0],
				Answers: []string{args[1]},
			}
		} else {
			out.Err(true, err)
			modules.Send(ctx.Message.ChannelID, "common_error.xml", map[string]interface{}{
				"Title":   "Ошибка",
				"Message": "Что то пошло не так :0",
			}, nil)
			return
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
		modules.Send(ctx.Message.ChannelID, "common_error.xml", map[string]interface{}{
			"Title":   "Ошибка",
			"Message": "Что то пошло не так :0",
		}, nil)
		return
	}

	modules.Send(ctx.Message.ChannelID, "triggers/trigger.added.xml", trigger, nil)
}

func (bot *module) onTriggerDelete(ctx *dgutils.MessageContext) {
	if ctx.Args == nil || len(ctx.Args) > 2 {
		modules.Send(ctx.Message.ChannelID, "triggers/usage.xml", nil, nil)
		return
	}

	name := strings.ToLower(ctx.Args[0])

	if len(ctx.Args) == 2 {
		i, err := strconv.Atoi(ctx.Args[1])
		if err != nil {
			modules.Send(ctx.Message.ChannelID, "common_error.xml", map[string]interface{}{
				"Title":   "Ошибка",
				"Message": "Индекс должен быть числом",
			}, nil)
			return
		}

		trigger, err := cache.Triggers.Get(name)
		if err != nil {
			modules.Send(ctx.Message.ChannelID, "triggers/trigger.not.found.xml", name, nil)
			return
		}

		if i > len(trigger.Answers) || i < 0 {
			modules.Send(ctx.Message.ChannelID, "common_error.xml", map[string]interface{}{
				"Title":   "Ошибка",
				"Message": "Триггер под индексом " + strconv.Itoa(i) + " отсутствует.",
			}, nil)
			return
		}

		answer := trigger.Answers[i-1]
		trigger.Answers = append(trigger.Answers[:i-1], trigger.Answers[i:]...)

		err = cache.Triggers.Set(trigger)
		if err != nil {
			out.Err(true, err)
			modules.Send(ctx.Message.ChannelID, "common_error.xml", map[string]interface{}{
				"Title":   "Ошибка",
				"Message": "Не удалось удалить триггер",
			}, nil)
			return
		}

		modules.Send(ctx.Message.ChannelID, "triggers/trigger.deleted.xml", map[string]string{
			"Name":   name,
			"Answer": answer,
		}, nil)

		return
	}

	err := cache.Triggers.Delete(name)
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			modules.Send(ctx.Message.ChannelID, "triggers/trigger.not.found.xml", name, nil)
			return
		}
		out.Err(true, err)
		modules.Send(ctx.Message.ChannelID, "common_error.xml", map[string]interface{}{
			"Title":   "Ошибка",
			"Message": "Что то пошло не так :0",
		}, nil)
	}

	modules.Send(ctx.Message.ChannelID, "triggers/trigger.deleted.xml", map[string]string{
		"Name":   name,
		"Answer": "",
	}, nil)
}
