package clubs

import (
	"time"

	"github.com/SteMak/house-tyan/app"
	"github.com/SteMak/house-tyan/modules"
	"github.com/SteMak/house-tyan/out"
	"github.com/SteMak/house-tyan/storage"
	"github.com/SteMak/house-tyan/util"
)

func (bot *module) onClubCreate(ctx *app.Context) {
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
		"Prefix":     bot.cmd.Prefix,
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

func (bot *module) onClubDelete(ctx *app.Context) {
	tx, err := storage.Tx()
	if err != nil {
		go out.Err(true, err)
		go modules.SendFail(ctx.Message.ChannelID, "База крашнулась на открытии", "Попробуйте снова позже.")
		go log.Error(err)
		return
	}

	err = storage.Clubs.DeleteByOwner(tx, ctx.Message.Author.ID)
	if err != nil {
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

func (bot *module) onClubKick(ctx *app.Context) {
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
