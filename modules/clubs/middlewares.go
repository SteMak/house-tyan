package clubs

import (
	"strings"

	"github.com/SteMak/house-tyan/libs/dgutils"
	"github.com/SteMak/house-tyan/modules"
	"github.com/SteMak/house-tyan/storage"
)

func (bot *module) middlewareChannel(ctx *dgutils.MessageContext) {
	ctx.Next()
}

func (bot *module) middlewareCreateClub(ctx *dgutils.MessageContext) {
	club, err := storage.Clubs.GetClubByUser(ctx.Message.Author.ID)
	if err != nil {
		return
	}

	if club != nil {
		modules.SendFail(ctx.Message.ChannelID, "Вы уже состоите в клубе", "Покинте текущий клуб, чтобы создать новый.")
		return
	}

		if ctx.Args == nil || ctx.Args[0] == "" {
		modules.SendFail(ctx.Message.ChannelID, "Имя клуба и символ не обнаружены", "Имя клуба не должно быть пустым.")
		return
	}
	if !strings.Contains(ctx.Args[0], " ") {
		modules.SendFail(ctx.Message.ChannelID, "Нет пробела", "Между символом и именем клуба должен быть пробел.")
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

	ctx.SetParam("symbol", clubSymbol)
	ctx.SetParam("name", clubName)

	ctx.Next()
}
