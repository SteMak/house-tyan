package clubs

import (
	"github.com/SteMak/house-tyan/libs/dgutils"
	"github.com/SteMak/house-tyan/modules"
	"github.com/SteMak/house-tyan/out"
	"github.com/SteMak/house-tyan/storage"
)

var (
	commands = map[string]interface{}{
		"createclub": &dgutils.Command{
			Raw: true,
			Handlers: []func(*dgutils.MessageContext){
				_module.middlewareChannel,
				_module.middlewareCreateClub,
			},
			Function: _module.onCreateClub,
		},
	}
)

func (bot *module) onCreateClub(ctx *dgutils.MessageContext) {
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

	modules.SendGood(ctx.Message.ChannelID, "Всё ок!", "Получилось!")
}
