package awards

import (
	"github.com/GrownNed/dgutils"
	conf "github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/modules"
	"github.com/SteMak/house-tyan/out"
	"github.com/SteMak/house-tyan/util"
	"github.com/pkg/errors"
)

func (bot *module) middlewareCreateBlank(ctx *dgutils.MessageContext) {
	if ctx.Message.GuildID != conf.Bot.GuildID {
		return
	}

	if ctx.Message.ChannelID != bot.config.Channels.Requests {
		return
	}

	member, err := bot.session.GuildMember(conf.Bot.GuildID, ctx.Message.Author.ID)
	if err != nil {
		out.Err(true, errors.WithStack(err))
		return
	}

	if !util.HasRole(member, bot.config.Roles.Requester) {
		return
	}

	if len(ctx.Args) == 0 {
		modules.Send(ctx.Message.ChannelID, "awards/usage.xml", nil, nil)
		return
	}

	ctx.Next()
}
