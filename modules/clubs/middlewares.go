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

func (bot *module) middlewareClubCreate(ctx *dgutils.MessageContext) {
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

func (bot *module) middlewareClubDelete(ctx *dgutils.MessageContext) {
	club, err := storage.Clubs.GetClubByUser(ctx.Message.Author.ID)
	if err != nil {
		return
	}

	if club == nil {
		modules.SendFail(ctx.Message.ChannelID, "Вы не состоите в клубе", "Станьте главой клуба, тогда вы сможете удалить его.")
		return
	}

	if club.OwnerID != ctx.Message.Author.ID {
		modules.SendFail(ctx.Message.ChannelID, "Вы не владелец клуба", "Станьте главой клуба, тогда вы сможете удалить его.")
		return
	}

	ctx.Next()
}

func (bot *module) middlewareClubKick(ctx *dgutils.MessageContext) {
	club, err := storage.Clubs.GetClubByUser(ctx.Message.Author.ID)
	if err != nil {
		return
	}

	if club == nil {
		modules.SendFail(ctx.Message.ChannelID, "Вы не состоите в клубе", "Создайте клуб, тогда вы сможете удалить участников.")
		return
	}

	if club.OwnerID != ctx.Message.Author.ID {
		modules.SendFail(ctx.Message.ChannelID, "Вы не владелец клуба", "Станьте главой клуба, тогда вы сможете удалить участников.")
		return
	}

	userID := ctx.Args[0]
	userID = strings.TrimPrefix(userID, "<@")
	userID = strings.TrimPrefix(userID, "!")
	userID = strings.TrimSuffix(userID, ">")

	_, err = bot.session.User(userID)
	if err != nil {
		modules.SendFail(ctx.Message.ChannelID, "Вааааа! Кто это?!", "Что это за существо???")
		return
	}

	bot.session.User(userID)
	userClub, err := storage.Clubs.GetClubByUser(userID)
	if err != nil || userClub == nil || userClub.ID != club.ID {
		modules.SendFail(ctx.Message.ChannelID, "Этот пользователь неуязвим", "Данный человек не состоит в вашем клубе.")
		return
	}

	ctx.SetParam("userID", userID)
	ctx.SetParam("club", club)

	ctx.Next()
}
