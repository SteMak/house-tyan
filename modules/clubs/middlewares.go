package clubs

import (
	"github.com/SteMak/house-tyan/libs/dgutils"
	"github.com/SteMak/house-tyan/storage"
)

func (bot *module) middlewareChannel(ctx *dgutils.MessageContext) {
	ctx.Next()
}

func (bot *module) middlewareClub(ctx *dgutils.MessageContext) {
	club, err := storage.Clubs.GetClubByUser(ctx.Message.Author.ID)
	if err != nil {
		return
	}
	if club == nil {
		ctx.SetParam("club", nil)
		ctx.Next()
	}

	ctx.SetParam("club", club)
	ctx.Next()
}
