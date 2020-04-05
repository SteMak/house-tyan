package awards

import (
	"github.com/SteMak/house-tyan/cache"
	"github.com/SteMak/house-tyan/libs/dgutils"
	"github.com/SteMak/house-tyan/modules"
	"github.com/SteMak/house-tyan/out"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

var (
	commands = map[string]interface{}{
		"запрос": &dgutils.Command{
			Raw: true,
			Handlers: []func(*dgutils.MessageContext){
				_module.middlewareChannel,
				_module.middlewareRole,
				_module.middlewareUsage,
			},
			Function: _module.onCreateBlank,
		},
		"выдать": &dgutils.Command{
			Raw: true,
			Handlers: []func(*dgutils.MessageContext){
				_module.middlewareChannel,
				_module.middlewareRole,
				_module.middlewareBlank,
			},
			Function: _module.onGiveOn,
		},
		// "сумма": &dgutils.Command{
		// 	Raw: true,
		// 	Handlers: []func(*dgutils.MessageContext){
		// 		_module.middlewareChannel,
		// 		_module.middlewareRole,
		// 	},
		// 	Function: _module.onRewadrUsers,
		// },
		// "отправить": &dgutils.Command{
		// 	Raw: true,
		// 	Handlers: []func(*dgutils.MessageContext){
		// 		_module.middlewareChannel,
		// 		_module.middlewareRole,
		// 	},
		// 	Function: _module.onRewadrUsers,
		// },
		// "отменить": &dgutils.Command{
		// 	Raw: true,
		// 	Handlers: []func(*dgutils.MessageContext){
		// 		_module.middlewareChannel,
		// 		_module.middlewareRole,
		// 	},
		// 	Function: _module.onRewadrUsers,
		// },
	}
)

func (bot *module) onGiveOn(ctx *dgutils.MessageContext) {
	blank := ctx.Param("blank").(*cache.Blank)

	if !blank.Actions.CanSetUsers {
		err := ctx.Session.ChannelMessageDelete(ctx.Message.ChannelID, ctx.Message.ID)
		if err != nil {
			out.Err(true, errors.WithStack(err))
		}
		return
	}

	users := ctx.Message.Mentions
	if len(users) == 0 {
		modules.Edit(blank.Message.ID, ctx.Message.ChannelID, "awards/blank.howmuch.xml",
			map[string]interface{}{
				"Err":    "Необходимо указать пользователей",
				"Author": ctx.Message.Author,
				"Blank":  blank,
			}, nil)
		return
	}

	var last cache.Reward
	last.Users = make(map[string]discordgo.User)
	for _, user := range users {
		last.Users[user.ID] = *user
	}
	blank.Rewards = append(blank.Rewards, last)
	blank.Actions = cache.BlankActions{
		SetReason: true,
		SetAmount: true,
		Discard:   true,
	}

	err := cache.Blanks.Set(blank)
	if err != nil {
		out.Err(true, errors.WithStack(err))
		modules.Send(ctx.Message.ChannelID, "main/common_error.xml",
			map[string]interface{}{
				"Title":   "Ошибка",
				"Message": "Не удалось отредактировать заявку",
			}, nil)
		return
	}

	modules.Edit(blank.Message.ID, ctx.Message.ChannelID, "awards/blank.howmuch.xml",
		map[string]interface{}{
			"Author": ctx.Message.Author,
			"Blank":  blank,
		}, nil)
}

func (bot *module) onCreateBlank(ctx *dgutils.MessageContext) {
	reason := ctx.Args[0]

	blank, err := cache.Blanks.Create(ctx.Message.Author.ID, reason, ctx.Message.Author, ctx.Message)
	if err != nil {
		out.Err(true, errors.WithStack(err))
		return
	}

	modules.Send(ctx.Message.ChannelID, "awards/blank.created.xml",
		map[string]interface{}{
			"Author":    ctx.Message.Author,
			"ExpiresAt": blank.ExpiresAt,
			"Reason":    reason,
		}, nil)

	modules.Send(ctx.Message.ChannelID, "awards/blank.reason.xml",
		map[string]interface{}{
			"Author":    ctx.Message.Author,
			"ExpiresAt": blank.ExpiresAt,
		}, nil)
	modules.Send(ctx.Message.ChannelID, "awards/blank.howmuch.xml",
		map[string]interface{}{
			"Author":    ctx.Message.Author,
			"ExpiresAt": blank.ExpiresAt,
			"Reason":    reason,
		}, nil)
	modules.Send(ctx.Message.ChannelID, "awards/blank.anythingelse.xml",
		map[string]interface{}{
			"Author":    ctx.Message.Author,
			"ExpiresAt": blank.ExpiresAt,
			"Reason":    reason,
		}, nil)
	modules.Send(ctx.Message.ChannelID, "awards/black.sended.xml",
		map[string]interface{}{}, nil)
	modules.Send(ctx.Message.ChannelID, "awards/blank.discarded.xml",
		map[string]interface{}{}, nil)
}
