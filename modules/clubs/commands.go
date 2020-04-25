package clubs

import (
	"time"

	conf "github.com/SteMak/house-tyan/config"

	"github.com/SteMak/house-tyan/cache"
	"github.com/SteMak/house-tyan/libs/dgutils"
	"github.com/SteMak/house-tyan/modules"
	"github.com/SteMak/house-tyan/out"
	"github.com/pkg/errors"
)

var (
	commands = map[string]interface{}{
		"запрос": &dgutils.Command{
			Raw: true,
			Handlers: []func(*dgutils.MessageContext){
				_module.middlewareChannel,
			},
			Function: _module.onCreateBlank,
		},
	}
)

func (bot *module) onCreateBlank(ctx *dgutils.MessageContext) {
	exists, err := cache.Blanks.Exists(ctx.Message.Author.ID)
	if err != nil {
		go out.Err(true, errors.WithStack(err))
		go modules.SendFail(ctx.Message.ChannelID, "Ошибка", "Не удалось создать заявку")
		return
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
		go modules.SendFail(ctx.Message.ChannelID, "Ошибка", "Не удалось создать заявку")
		return
	}

	blank.Message = *m

	if err := cache.Blanks.Create(blank); err != nil {
		go out.Err(true, errors.WithStack(err))
		go modules.SendFail(ctx.Message.ChannelID, "Ошибка", "Не удалось создать заявку")
		return
	}
}
