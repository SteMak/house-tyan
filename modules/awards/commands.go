package awards

import (
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
				_module.middlewareCreateBlank,
			},
			Description: "Создает запрос на выдачу денег",
			Function:    _module.onCreateBlank,
		},
	}
)

func (bot *module) onCreateBlank(ctx *dgutils.MessageContext) {
	reason := ctx.Args[0]

	blank, err := cache.Blanks.Create(ctx.Message.Author.ID, reason, ctx.Message.Author)
	if err != nil {
		out.Err(true, errors.WithStack(err))
		return
	}

	modules.Send(ctx.Message.ChannelID, "awards/blank.reason.xml",
		map[string]interface{}{
			"Author":    ctx.Message.Author,
			"ExpiresAt": blank.ExpiresAt,
		}, nil)
	modules.Send(ctx.Message.ChannelID, "awards/blank.created.xml",
		map[string]interface{}{
			"Author":    ctx.Message.Author,
			"ExpiresAt": blank.ExpiresAt,
			"Reason":    reason,
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
