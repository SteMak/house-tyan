package voices

import (
	"github.com/SteMak/house-tyan/cache"
	"github.com/SteMak/house-tyan/libs/dgutils"
	"github.com/SteMak/house-tyan/out"
	//	"github.com/SteMak/house-tyan/out"
	//	"github.com/SteMak/house-tyan/storage"
	//	"github.com/SteMak/house-tyan/util"
)

var (
	commands = map[string]interface{}{
		"voice": &dgutils.Group{
			Commands: map[string]interface{}{

				"unlock": &dgutils.Command{
					Raw:         false,
					Description: "Открыть доступ",
					Handlers: []func(*dgutils.MessageContext){
						_module.middlewareChannelManage,
					},
					Function: _module.onVoiceUnlock,
				},

				"lock": &dgutils.Command{
					Raw:         false,
					Description: "Закрыть доступ",
					Handlers: []func(*dgutils.MessageContext){
						_module.middlewareChannelManage,
					},
					Function: _module.onVoiceLock,
				},

				"show": &dgutils.Command{
					Raw:         false,
					Description: "Показать канал",
					Handlers: []func(*dgutils.MessageContext){
						_module.middlewareChannelManage,
					},
					Function: _module.onVoiceShow,
				},

				"hide": &dgutils.Command{
					Raw:         false,
					Description: "Скрыть канал",
					Handlers: []func(*dgutils.MessageContext){
						_module.middlewareChannelManage,
					},
					Function: _module.onVoiceHide,
				},

				"info": &dgutils.Command{
					Raw:         false,
					Description: "Информация о приватном канале",
					Handlers: []func(*dgutils.MessageContext){
						_module.middlewareChannelInfo,
					},
					Function: _module.onVoiceInfo,
				},
				"save": &dgutils.Command{
					Raw:         false,
					Description: "Сохранение данных канала",
					Handlers: []func(*dgutils.MessageContext){
						_module.middlewareChannelManage,
					},
					Function: _module.onVoiceSave,
				},
			},
		},
	}
)

func (bot *module) onVoiceUnlock(ctx *dgutils.MessageContext) {
	setPermissions(ctx, permConnect, true)
}

func (bot *module) onVoiceLock(ctx *dgutils.MessageContext) {
	setPermissions(ctx, permConnect, false)
}

func (bot *module) onVoiceShow(ctx *dgutils.MessageContext) {
	setPermissions(ctx, permView, true)
}

func (bot *module) onVoiceHide(ctx *dgutils.MessageContext) {
	setPermissions(ctx, permView, false)
}

func (bot *module) onVoiceInfo(ctx *dgutils.MessageContext) {
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

	messageEmbed := getInfo(ctx.Session, channelID)
	ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, &messageEmbed)

}

func (bot *module) onVoiceSave(ctx *dgutils.MessageContext) {
	var voice *cache.Voice
	var channelID string

	for _, state := range voiceStatesCache {
		if state.UserID == ctx.Message.Author.ID {
			channelID = state.ChannelID
		}
	}

	channel, err := ctx.Session.Channel(channelID)
	if err != nil {
		out.Err(true, err)
		return
	}

	voice = &cache.Voice{
		ID:          ctx.Message.Author.ID,
		Name:        channel.Name,
		Permissions: channel.PermissionOverwrites,
	}

	err = cache.Voices.Set(voice)
	if err != nil {
		out.Err(true, err)
		return
	}
}
