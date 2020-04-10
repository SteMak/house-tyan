package awards

import (
	"strconv"
	"time"

	"github.com/SteMak/house-tyan/storage"

	conf "github.com/SteMak/house-tyan/config"

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
				_module.middlewareDeleteMessage,
			},
			Function: _module.onCreateBlank,
		},
		"выдать": &dgutils.Command{
			Raw: true,
			Handlers: []func(*dgutils.MessageContext){
				_module.middlewareChannel,
				_module.middlewareRole,
				_module.middlewareBlank,
				_module.middlewareDeleteMessage,
			},
			Function: _module.onUsers,
		},
		"сумма": &dgutils.Command{
			Raw: true,
			Handlers: []func(*dgutils.MessageContext){
				_module.middlewareChannel,
				_module.middlewareRole,
				_module.middlewareBlank,
				_module.middlewareDeleteMessage,
			},
			Function: _module.onAmount,
		},
		"отправить": &dgutils.Command{
			Raw: true,
			Handlers: []func(*dgutils.MessageContext){
				_module.middlewareChannel,
				_module.middlewareRole,
				_module.middlewareBlank,
				_module.middlewareDeleteMessage,
			},
			Function: _module.onSend,
		},
		"отменить": &dgutils.Command{
			Raw: true,
			Handlers: []func(*dgutils.MessageContext){
				_module.middlewareChannel,
				_module.middlewareRole,
				_module.middlewareBlank,
				_module.middlewareDeleteMessage,
			},
			Function: _module.onDiscard,
		},
	}
)

func (bot *module) onSend(ctx *dgutils.MessageContext) {
	blank := ctx.Param("blank").(*cache.Blank)
	if !blank.Actions.Send {
		return
	}

	m := modules.Send(bot.config.Channels.Confirm, "awards/blank.xml", map[string]interface{}{
		"Reason":  blank.Reason,
		"Rewards": blank.Rewards,
	}, nil)
	if m == nil {
		modules.Send(ctx.Message.ChannelID, "common_error.xml", map[string]interface{}{
			"Title":   "Ошибка",
			"Message": "Не удалось отправить заявку",
		}, nil)
		return
	}

	cache.Blanks.Delete(blank.ID)

	tx, err := storage.Tx()
	if err != nil {
		out.Err(true, errors.WithStack(err))
		modules.Send(ctx.Message.ChannelID, "common_error.xml", map[string]interface{}{
			"Title":   "Ошибка",
			"Message": "Не удалось отправить заявку",
		}, nil)
		tx.Rollback()
		return
	}

	err = storage.Awards.Create(tx, m.ID, blank)
	if err != nil {
		out.Err(true, errors.WithStack(err))
		modules.Send(ctx.Message.ChannelID, "common_error.xml", map[string]interface{}{
			"Title":   "Ошибка",
			"Message": "Не удалось отправить заявку",
		}, nil)
		tx.Rollback()
		return
	}

	tx.Commit()

	modules.Edit(blank.Message.ID, ctx.Message.ChannelID, "awards/black.sended.xml", map[string]interface{}{
		"Blank": blank,
	}, nil)
}

func (bot *module) onDiscard(ctx *dgutils.MessageContext) {
	blank := ctx.Param("blank").(*cache.Blank)
	if !blank.Actions.Discard {
		return
	}

	cache.Blanks.Delete(blank.ID)
	go modules.Edit(blank.Message.ID, ctx.Message.ChannelID, "awards/blank.discarded.xml", nil, nil)
}

func (bot *module) onAmount(ctx *dgutils.MessageContext) {
	blank := ctx.Param("blank").(*cache.Blank)
	if !blank.Actions.SetAmount {
		return
	}

	amount, err := strconv.ParseUint(ctx.Args[0], 10, 64)
	if err != nil {
		modules.Edit(blank.Message.ID, ctx.Message.ChannelID, "awards/blank.amount.xml", map[string]interface{}{
			"Err":    "Укажите корректную сумму (целое, положительное число)",
			"Author": ctx.Message.Author,
			"Blank":  blank,
		}, nil)
		return
	}
	if amount == 0 {
		modules.Edit(blank.Message.ID, ctx.Message.ChannelID, "awards/blank.amount.xml", map[string]interface{}{
			"Err":    "Сумма не должна быть 0",
			"Author": ctx.Message.Author,
			"Blank":  blank,
		}, nil)
		return
	}

	blank.Rewards[len(blank.Rewards)-1].Amount = amount

	blank.Actions = cache.BlankActions{
		Send:     true,
		SetUsers: true,
		Discard:  true,
	}

	if err := cache.Blanks.Set(blank); err != nil {
		out.Err(true, errors.WithStack(err))
		modules.Send(ctx.Message.ChannelID, "common_error.xml", map[string]interface{}{
			"Title":   "Ошибка",
			"Message": "Не удалось отредактировать заявку",
		}, nil)
		return
	}

	modules.Edit(blank.Message.ID, ctx.Message.ChannelID, "awards/blank.users.xml", map[string]interface{}{
		"Author": ctx.Message.Author,
		"Blank":  blank,
	}, nil)
}

func (bot *module) onUsers(ctx *dgutils.MessageContext) {
	blank := ctx.Param("blank").(*cache.Blank)
	if !blank.Actions.SetUsers {
		return
	}

	users := ctx.Message.Mentions
	if len(users) == 0 {
		modules.Edit(blank.Message.ID, ctx.Message.ChannelID, "awards/blank.users.xml", map[string]interface{}{
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
		SetAmount: true,
		Discard:   true,
	}

	if err := cache.Blanks.Set(blank); err != nil {
		out.Err(true, errors.WithStack(err))
		modules.Send(ctx.Message.ChannelID, "common_error.xml", map[string]interface{}{
			"Title":   "Ошибка",
			"Message": "Не удалось отредактировать заявку",
		}, nil)
		return
	}

	modules.Edit(blank.Message.ID, ctx.Message.ChannelID, "awards/blank.amount.xml", map[string]interface{}{
		"Author": ctx.Message.Author,
		"Blank":  blank,
	}, nil)
}

func (bot *module) onCreateBlank(ctx *dgutils.MessageContext) {
	exists, err := cache.Blanks.Exists(ctx.Message.Author.ID)
	if err != nil {
		out.Err(true, errors.WithStack(err))
		modules.Send(ctx.Message.ChannelID, "common_error.xml", map[string]interface{}{
			"Title":   "Ошибка",
			"Message": "Не удалось отредактировать заявку",
		}, nil)
	}

	if exists {
		return
	}

	blank := &cache.Blank{
		ID:     ctx.Message.Author.ID,
		Reason: ctx.Args[0],
		Author: *ctx.Message.Author,
		Actions: cache.BlankActions{
			SetUsers: true,
			Discard:  true,
		},
		ExpiresAt: time.Now().UTC().Add(conf.Cache.TTL.Blank),
	}

	m := modules.Send(ctx.Message.ChannelID, "awards/blank.users.xml", map[string]interface{}{
		"Author": ctx.Message.Author,
		"Blank":  blank,
	}, nil)

	if m == nil {
		modules.Send(ctx.Message.ChannelID, "common_error.xml", map[string]interface{}{
			"Title":   "Ошибка",
			"Message": "Не удалось отредактировать заявку",
		}, nil)
	}

	blank.Message = *m

	if err := cache.Blanks.Create(blank); err != nil {
		out.Err(true, errors.WithStack(err))
		modules.Send(ctx.Message.ChannelID, "common_error.xml", map[string]interface{}{
			"Title":   "Ошибка",
			"Message": "Не удалось отредактировать заявку",
		}, nil)
		return
	}
}
