package voices

import (
	"github.com/SteMak/house-tyan/libs/dgutils"
	"github.com/SteMak/house-tyan/out"
)

func (bot *module) middlewareChannelManage(ctx *dgutils.MessageContext) {
	channelID := "" 
	for _, state := range voiceStatesCache {
		if state.UserID == ctx.Message.Author.ID {
			channelID = state.ChannelID
		}
	}

	if channelID == "" {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Вы не находитесь в голосовом канале!")
		return
	}

	channel, err := ctx.Session.Channel(channelID) 
	if err != nil {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Упс, мы не смогли получить данные канала!")
		out.Err(true, err)
		return
	}

	if channel.ParentID != privateVoices[ctx.Message.GuildID]["coreParentID"] {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Вы не находитесь в приватном голосовом канале!")
		out.Err(true, err)
		return
	}

	isOwner := false 
	for _, perm := range channel.PermissionOverwrites {
		if perm.ID == ctx.Message.Author.ID && permManage == permManage & perm.Allow {
			isOwner = true
		}
	}

	if !isOwner {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Вы не являетесь владельцем канала!")
		out.Err(true, err)
		return
	}

	ctx.Next()
}

func (bot *module) middlewareChannelInfo(ctx *dgutils.MessageContext) {
	channelID := ""

    if len(ctx.Args) != 0 {
        channelID = ctx.Args[0]
    } else {
        for _, state := range voiceStatesCache {
            if ctx.Message.Author.ID == state.UserID {
                channelID = state.ChannelID
                break
            }
        }
    }

    if channelID == "" {
    	ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Вы не находитесь в голосовом канале!")
    	out.Err(true, err)
    	return
    }

    channel, err := ctx.Session.Channel(channelID)
    if err != nil {
    	ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Не удалось получить данные о категории приватных каналов! Ой.")
    	out.Err(true, err)
    	return
    }

	if channel.ParentID != privateVoices[ctx.Message.GuildID]["coreParentID"] {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Информацию о голосовом канале не удалось получить, это не приватный канал")
		return
	}

	ctx.Next()
}
