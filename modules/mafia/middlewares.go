package mafia

import (
	"github.com/SteMak/house-tyan/libs/dgutils"
	"github.com/SteMak/house-tyan/out"
)

func (bot *module) middlewareIsOwner(ctx *dgutils.MessageContext) {
	if ctx.Message.Author.ID == bot.config.OwnerID {
		ctx.Next()
	}
}

func (bot *module) middlewareMafia(ctx *dgutils.MessageContext) {
	if bot.game != nil {
		ctx.Next()
		return
	}

	joinMsg, err = bot.session.ChannelMessageSend(ctx.Message.ChannelID, "Мафия не создана.")
	if err != nil {
		out.Err(true, err)
	}
}

func (bot *module) middlewareMafiaIsNotStarted(ctx *dgutils.MessageContext) {
	if joinMsg != nil {
		ctx.Next()
		return
	}

	joinMsg, err = bot.session.ChannelMessageSend(ctx.Message.ChannelID, "Мафия уже началась.")
	if err != nil {
		out.Err(true, err)
	}
}

func (bot *module) middlewareMafiaIsStarted(ctx *dgutils.MessageContext) {
	if joinMsg == nil {
		ctx.Next()
		return
	}

	joinMsg, err = bot.session.ChannelMessageSend(ctx.Message.ChannelID, "Мафия ещё не началась.")
	if err != nil {
		out.Err(true, err)
	}
}

func (bot *module) middlewareMafiaVoteIsNotStated(ctx *dgutils.MessageContext) {
	if voteMsg == nil {
		ctx.Next()
		return
	}

	joinMsg, err = bot.session.ChannelMessageSend(ctx.Message.ChannelID, "Идёт голосование.")
	if err != nil {
		out.Err(true, err)
	}
}

func (bot *module) middlewareMafiaVoteIsStated(ctx *dgutils.MessageContext) {
	if voteMsg != nil {
		ctx.Next()
		return
	}

	joinMsg, err = bot.session.ChannelMessageSend(ctx.Message.ChannelID, "Голосование не идёт.")
	if err != nil {
		out.Err(true, err)
	}
}
