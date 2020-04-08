package triggers

import (
	conf "github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/libs/dgutils"
	"github.com/SteMak/house-tyan/modules"
	"github.com/SteMak/house-tyan/out"
)

func (bot *module) middlewareAdmin(ctx *dgutils.MessageContext) {
	if ctx.Message.GuildID != conf.Bot.GuildID {
		return
	}

	for _, roleID := range ctx.Message.Member.Roles {
		role, err := bot.session.State.Role(ctx.Message.GuildID, roleID)
		if err != nil {
			out.Err(true, err)
			return
		}

		if (role.Permissions & 0x8) == 0x8 {
			ctx.Next()
		}
	}
}

func (bot *module) middlewareUsage(ctx *dgutils.MessageContext) {
	if len(ctx.Args) == 0 {
		modules.Send(ctx.Message.ChannelID, "triggers/usage.xml", nil, nil)
		return
	}
	ctx.Next()
}
