package clubs

import (
	"time"

	"github.com/SteMak/house-tyan/libs/dgutils"
	"github.com/SteMak/house-tyan/modules"
	"github.com/SteMak/house-tyan/out"
	"github.com/SteMak/house-tyan/storage"
	"github.com/SteMak/house-tyan/util"
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
						_module.middlewareClubOwner,
					},
					Function: _module.onClubDelete,
				},
				"info": &dgutils.Command{
					Raw: true,
					Handlers: []func(*dgutils.MessageContext){
						_module.middlewareChannel,
						_module.middlewareClubInfo,
					},
					Function: _module.onClubInfo,
				},
				"desc": &dgutils.Command{
					Raw: true,
					Handlers: []func(*dgutils.MessageContext){
						_module.middlewareChannel,
						_module.middlewareArgsPresented,
						_module.middlewareClubOwner,
					},
					Function: _module.onClubEditDescription,
				},
				"kick": &dgutils.Command{
					Raw: true,
					Handlers: []func(*dgutils.MessageContext){
						_module.middlewareChannel,
						_module.middlewareClubOwner,
						_module.middlewareClubKick,
					},
					Function: _module.onClubKick,
				},
				"leave": &dgutils.Command{
					Raw: true,
					Handlers: []func(*dgutils.MessageContext){
						_module.middlewareChannel,
						_module.middlewareNotClubOwner,
					},
					Function: _module.onClubLeave,
				},
				"apply": &dgutils.Command{
					Raw: true,
					Handlers: []func(*dgutils.MessageContext){
						_module.middlewareChannel,
						_module.middlewareClubApply,
					},
					Function: _module.onClubApply,
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
		go log.Error(err)
		return
	}

	expiredAt := util.Midnight(time.Now().UTC().Add(bot.config.NotVerifiedLifetime))

	club := storage.Club{
		OwnerID:   ctx.Message.Author.ID,
		Title:     ctx.Param("name").(string),
		Symbol:    ctx.Param("symbol").(string),
		ExpiredAt: &expiredAt,
	}

	err = storage.Clubs.Create(tx, &club)
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "Создание клуба полетело", "Попробуйте снова позже.")
		go log.Error(err)
		tx.Rollback()
		return
	}

	data := map[string]interface{}{
		"Prefix":     bot.cmds.Prefix,
		"Owner":      ctx.Message.Author,
		"Club":       club,
		"MinMembers": bot.config.MinimumMembers,
		"Price":      int64(bot.config.Price),
	}

	m := modules.Send(ctx.Message.ChannelID, "clubs/created.xml", data, nil)
	if m == nil {
		go modules.SendFail(ctx.Message.ChannelID, "Создание клуба полетело", "Попробуйте снова позже.")
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "База крашнулась на закрытии", "Попробуйте снова позже.")
		go log.Error(err)
		return
	}
}

func (bot *module) onClubDelete(ctx *dgutils.MessageContext) {
	tx, err := storage.Tx()
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "База крашнулась на открытии", "Попробуйте снова позже.")
		go log.Error(err)
		return
	}

	club := ctx.Param("club").(*storage.Club)
	if err = club.Delete(tx); err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "Удаление клуба полетело", "Попробуйте снова позже.")
		go log.Error(err)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "База крашнулась на закрытии", "Попробуйте снова позже.")
		go log.Error(err)
		return
	}

	modules.SendGood(ctx.Message.ChannelID, "Клуб удалён", "Стирание прошло успешно")
}

func (bot *module) onClubKick(ctx *dgutils.MessageContext) {
	tx, err := storage.Tx()
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "База крашнулась на открытии", "Попробуйте снова позже.")
		go log.Error(err)
		return
	}

	club := ctx.Param("club").(*storage.Club)
	userID := ctx.Param("userID").(string)

	err = club.DeleteMember(tx, userID)
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "Ошибка удаления учасника", "Попробуйте снова позже.")
		go log.Error(err)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "База крашнулась на закрытии", "Попробуйте снова позже.")
		go log.Error(err)
		return
	}

	modules.SendGood(ctx.Message.ChannelID, "Участник исключён", "Стирание прошло успешно")
}

func (bot *module) onClubLeave(ctx *dgutils.MessageContext) {
	tx, err := storage.Tx()
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "База крашнулась на открытии", "Попробуйте снова позже.")
		go log.Error(err)
		return
	}

	club := ctx.Param("club").(*storage.Club)
	userID := ctx.Message.Author.ID

	err = club.DeleteMember(tx, userID)
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "Ошибка удаления учасника", "Попробуйте снова позже.")
		go log.Error(err)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "База крашнулась на закрытии", "Попробуйте снова позже.")
		go log.Error(err)
		return
	}

	modules.SendGood(ctx.Message.ChannelID, "Участник исключён", "Стирание прошло успешно")
}

func (bot *module) onClubInfo(ctx *dgutils.MessageContext) {
	club := ctx.Param("club").(*storage.Club)
	members, err := club.GetMembers()
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "Не удалось получить список участников", "Попробуйте снова позже.")
		go log.Error(err)
		return
	}

	m := modules.Send(ctx.Message.ChannelID, "clubs/info.xml", map[string]interface{}{
		"Club":    club,
		"Members": members,
	}, nil)
	if m == nil {
		go modules.SendFail(ctx.Message.ChannelID, "Не удалось вывести информацию о клубе", "Попробуйте снова позже.")
		return
	}
}

func (bot *module) onClubEditDescription(ctx *dgutils.MessageContext) {
	tx, err := storage.Tx()
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "База крашнулась на открытии", "Попробуйте снова позже.")
		go log.Error(err)
		return
	}

	club := ctx.Param("club").(*storage.Club)
	club.EditDescription(tx, ctx.Args[0])

	err = tx.Commit()
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "База крашнулась на закрытии", "Попробуйте снова позже.")
		go log.Error(err)
		return
	}

	m := modules.Send(ctx.Message.ChannelID, "clubs/description_edit.xml", club, nil)
	if m == nil {
		go modules.SendFail(ctx.Message.ChannelID, "Не удалось вывести информацию об изменении описания", "Попробуйте снова позже.")
		return
	}
}

func (bot *module) onClubApply(ctx *dgutils.MessageContext) {
	tx, err := storage.Tx()
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "База крашнулась на открытии", "Попробуйте снова позже.")
		go log.Error(err)
		return
	}

	club := ctx.Param("club").(*storage.Club)

	dm, err := ctx.Session.UserChannelCreate(club.OwnerID)
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "Не удалось создать DM с главой клуба", "Попробуйте снова позже.")
		go log.Error(err)
		return
	}

	m := modules.Send(dm.ID, "clubs/club_accept.xml", ctx.Message.Author, nil)
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "Не удалось отправить запрос", "Попробуйте снова позже.")
		go log.Error(err)
		return
	}

	err = ctx.Session.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	if err != nil {
		go out.Err(false, err)
		go log.Error(err)
		return
	}

	err = ctx.Session.MessageReactionAdd(m.ChannelID, m.ID, "❎")
	if err != nil {
		go out.Err(false, err)
		go log.Error(err)
		return
	}

	err = club.ClubApply(tx, ctx.Message.Author.ID, m.ID)
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "База крашнулась на комите", "Попробуйте снова позже.")
		go log.Error(err)
		return
	}

	err = tx.Commit()
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "База крашнулась на закрытии", "Попробуйте снова позже.")
		go log.Error(err)
		return
	}

	modules.Send(ctx.Message.ChannelID, "clubs/club_apply.xml", club, nil)
}
