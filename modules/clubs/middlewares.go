package clubs

import (
	"strings"

	conf "github.com/SteMak/house-tyan/config"

	"github.com/SteMak/house-tyan/libs/dgutils"
	"github.com/SteMak/house-tyan/modules"
	"github.com/SteMak/house-tyan/storage"
)

func (bot *module) middlewareChannel(ctx *dgutils.MessageContext) {
	if conf.Bot.Channels.Console != ctx.Message.ChannelID {
		return
	}
	ctx.Next()
}

func (bot *module) middlewareClubCreate(ctx *dgutils.MessageContext) {
	for _, banedCombinamions := range []string{"<@", "<#", "@here", "@everyone"} {
		if !strings.Contains(ctx.Message.Content, banedCombinamions) {
			continue
		}

		modules.SendFail(ctx.Message.ChannelID, "Не правильные данные", "Нельзя использовать `"+banedCombinamions+"`.")
		return
	}

	for _, banWord := range bot.config.BadWords {
		if !strings.Contains(ctx.Message.Content, banWord) {
			continue
		}

		modules.SendFail(ctx.Message.ChannelID, "Не правильные данные", "Нельзя использовать слово ||"+banWord+"||.")
		return
	}

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

	ctx.SetParam("club", club)

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

	if userID == ctx.Message.Author.ID {
		modules.SendFail(ctx.Message.ChannelID, "Нельзя кикнуть себя", "Найдите кого нибудь другого.")
		return
	}

	_, err = bot.session.User(userID)
	if err != nil {
		modules.SendFail(ctx.Message.ChannelID, "Вааааа! Кто это?!", "Что это за существо???")
		return
	}

	userClub, err := storage.Clubs.GetClubByUser(userID)
	if err != nil || userClub == nil || userClub.ID != club.ID {
		modules.SendFail(ctx.Message.ChannelID, "Этот пользователь неуязвим", "Данный человек не состоит в вашем клубе.")
		return
	}

	ctx.SetParam("userID", userID)
	ctx.SetParam("club", club)

	ctx.Next()
}

func (bot *module) middlewareClubInfo(ctx *dgutils.MessageContext) {
	var (
		club *storage.Club
		err  error
	)

	if len(ctx.Args) == 0 {
		club, err = storage.Clubs.GetClubByUser(ctx.Message.Author.ID)
		if err != nil {
			return
		}
		if club == nil {
			go modules.SendFail(ctx.Message.ChannelID, "Вы не в клубе", "Попробуйте когда будете в клубе")
			return
		}
	} else {
		userID := ctx.Args[0]
		userID = strings.TrimPrefix(userID, "<@")
		userID = strings.TrimPrefix(userID, "!")
		userID = strings.TrimSuffix(userID, ">")

		if club, err = storage.Clubs.GetClubByTitle(ctx.Args[0]); err != nil || club != nil {
			if err != nil {
				return
			}
		} else if club, err = storage.Clubs.GetClubByUser(userID); err != nil || club != nil {
			if err != nil {
				return
			}
		}
	}

	if club == nil {
		go modules.SendFail(ctx.Message.ChannelID, "Клуб не найден", "Укажите члена клуба или его название")
		return
	}

	ctx.SetParam("club", club)

	ctx.Next()
}

func (bot *module) middlewareClubEditInfo(ctx *dgutils.MessageContext) {
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

	ctx.SetParam("club", club)

	ctx.Next()
}
