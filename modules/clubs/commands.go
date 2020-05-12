package clubs

import (
	"github.com/SteMak/house-tyan/libs/dgutils"
	"github.com/SteMak/house-tyan/modules"
	"github.com/SteMak/house-tyan/out"
	"github.com/SteMak/house-tyan/storage"
)

var (
	commands = map[string]interface{}{
		"club": &dgutils.Group{
			Commands: map[string]interface{}{
				"create": &dgutils.Command{
					Raw: true,
					Handlers: []func(*dgutils.MessageContext){
						_module.middlewareChannel,
						_module.middlewareClubCreate,
					},
					Function: _module.onClubCreate,
				},
				"delete": &dgutils.Command{
					Raw: true,
					Handlers: []func(*dgutils.MessageContext){
						_module.middlewareChannel,
						_module.middlewareClubDelete,
					},
					Function: _module.onClubDelete,
				},
				"kick": &dgutils.Command{
					Raw: true,
					Handlers: []func(*dgutils.MessageContext){
						_module.middlewareChannel,
						_module.middlewareClubKick,
					},
					Function: _module.onClubKick,
				},
			},
		},
	}
)

func (bot *module) onClubCreate(ctx *dgutils.MessageContext) {
	tx, err := storage.Tx()
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "База крашнулась на открытии", "Попробуйте снова позже.")
		return
	}

	err = storage.Clubs.Create(tx, &storage.Club{
		OwnerID: ctx.Message.Author.ID,
		Title:   ctx.Param("name").(string),
		Symbol:  ctx.Param("symbol").(string),
	})
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "Создание клуба полетело", "Попробуйте снова позже.")
		return
	}

	err = tx.Commit()
	if err != nil {
		go out.Err(true, err)
		modules.SendFail(ctx.Message.ChannelID, "База крашнулась на закрытии", "Попробуйте снова позже.")
		return
	}

	modules.SendGood(ctx.Message.ChannelID, "Клуб создан", "Создание прошло успешно")
}

func (bot *module) onClubDelete(ctx *dgutils.MessageContext) {
	tx, err := storage.Tx()
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "База крашнулась на открытии", "Попробуйте снова позже.")
		return
	}

	err = storage.Clubs.DeleteByOwner(tx, ctx.Message.Author.ID)
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "Удаление клуба полетело", "Попробуйте снова позже.")
		return
	}

	err = tx.Commit()
	if err != nil {
		go out.Err(true, err)
		modules.SendFail(ctx.Message.ChannelID, "База крашнулась на закрытии", "Попробуйте снова позже.")
		return
	}

	modules.SendGood(ctx.Message.ChannelID, "Клуб удалён", "Стирание прошло успешно")
}

func (bot *module) onClubKick(ctx *dgutils.MessageContext) {
	tx, err := storage.Tx()
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "База крашнулась на открытии", "Попробуйте снова позже.")
		return
	}

	club := ctx.Param("club").(*storage.Club)
	userID := ctx.Param("userID").(string)

	club.DeleteMember(tx, userID)

	err = tx.Commit()
	if err != nil {
		go out.Err(true, err)
		modules.SendFail(ctx.Message.ChannelID, "База крашнулась на закрытии", "Попробуйте снова позже.")
		return
	}

	modules.SendGood(ctx.Message.ChannelID, "Участник исключён", "Стирание прошло успешно")
}
