package clubs

import (
	"github.com/SteMak/house-tyan/libs/dgutils"
)

func (bot *module) middlewareClub(ctx *dgutils.MessageContext) {
	ctx.Next()
}

func (bot *module) middlewareChannel(ctx *dgutils.MessageContext) {
	ctx.Next()
}
