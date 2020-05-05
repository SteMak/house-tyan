package clubs

import (
	"github.com/SteMak/house-tyan/libs/dgutils"
	"github.com/SteMak/house-tyan/modules"
	"github.com/SteMak/house-tyan/out"
	"github.com/SteMak/house-tyan/storage"
	"strings"
)

var (
	commands = map[string]interface{}{
		"createclub": &dgutils.Command{
			Raw: true,
			Handlers: []func(*dgutils.MessageContext){
				_module.middlewareChannel,
				_module.middlewareClub,
			},
			Function: _module.onCreateClub,
		},
	}
)

func (bot *module) onCreateClub(ctx *dgutils.MessageContext) {
	if ctx.Param("club") != nil {
		modules.SendFail(ctx.Message.ChannelID, "Вы уже состоите в клубе", "Покинте текущий клуб, чтобы создать новый.")
		return
	}

	if len(ctx.Args) == 0 || ctx.Args[0] == "" {
		modules.SendFail(ctx.Message.ChannelID, "Имя клуба не обнаружено", "Имя клуба не должно быть пустым.")
		return
	}
	if !strings.Contains(ctx.Args[0], " ") {
		modules.SendFail(ctx.Message.ChannelID, "Не найден символ клуба", "Через пробел нужно указать символ клуба и его название.")
		return
	}

	clubSymbolName := strings.SplitN(ctx.Args[0], " ", 2)
	clubSymbol := clubSymbolName[0]
	clubName := clubSymbolName[1]
	if len(clubSymbol) == 0 {
		modules.SendFail(ctx.Message.ChannelID, "Символ не найден", "Символ не должен быть пустым.")
		return
	}
	if len(clubName) == 0 {
		modules.SendFail(ctx.Message.ChannelID, "Имя клуба не обнаружено", "Имя клуба не должно быть пустым.")
		return
	}

	tx, err := storage.Tx()
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "База крашнулась на открытии", "Попробуйте снова позже.")
		return
	}

	err = storage.Clubs.Create(tx, &storage.Club{
		OwnerID: ctx.Message.Author.ID,
		Title:   clubName,
		Symbol:  clubSymbol,
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

	modules.SendFail(ctx.Message.ChannelID, "Всё ок!", "Получилось!")
}
