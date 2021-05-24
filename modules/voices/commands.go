package voices

import (
	
	"github.com/SteMak/house-tyan/libs/dgutils"
//	"github.com/SteMak/house-tyan/modules"
//	"github.com/SteMak/house-tyan/out"
//	"github.com/SteMak/house-tyan/storage"
//	"github.com/SteMak/house-tyan/util"
)

var (
	commands = map[string] interface{}{
		"voice": &dgutils.Group{
			Commands: map[string]interface{}{

				"unlock": &dgutils.Command{
					Raw: false,
					Description: "Открыть доступ к каналу пользователю или роли",
					Handlers: []func(*dgutils.MessageContext){
						_module.middlewareChannel,
					},
					Function: _module.onVoiceUnlock,
				},
				
				"lock": &dgutils.Command{
					Raw:false,
					Description: "Закрыть доступ к каналу пользователю или роли",
					Handlers: []func(*dgutils.MessageContext){
						_module.middlewareChannel,
					},
					Function: _module.onVoiceLock,
				},

				"show": &dgutils.Command{
					Raw: false,
					Description: "Показать канал для пользователя или роли",
					Handlers: []func(*dgutils.MessageContext){
						_module.middlewareChannel,
					},
					Function: _module.onVoiceShow,
				},
				
				"hide": &dgutils.Command{
					Raw: false,
					Description: "Скрыть канал от пользователя или роли",
					Handlers: []func(*dgutils.MessageContext){
						_module.middlewareChannel,
					},
					Function: _module.middlewareChannel,
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
	
}
