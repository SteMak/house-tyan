package clubs

import (
	"github.com/bwmarrin/discordgo"

	conf "github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/libs/dgutils"
	"github.com/SteMak/house-tyan/modules"
	"github.com/SteMak/house-tyan/storage"
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

	tx, err := storage.Tx()
	if err != nil {
		modules.SendFail(ctx.Message.ChannelID, "База крашнулась на открытии", "Попробуйте снова позже.")
		return
	}

	clubName := ctx.Args[0]
	clubChannel, err := ctx.Session.GuildChannelCreateComplex(conf.Bot.GuildID, discordgo.GuildChannelCreateData{
		NSFW:     true,
		Name:     clubName,
		ParentID: "",
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			&discordgo.PermissionOverwrite{
				
			},
		},
		Topic: "",
		Type:  0,
	})
	if err != nil {
		modules.SendFail(ctx.Message.ChannelID, "Канал создать не удалось", "Попробуйте снова позже.")
		return
	}
	clubRole, err := ctx.Session.GuildRoleCreate(conf.Bot.GuildID)
	if err != nil {
		modules.SendFail(ctx.Message.ChannelID, "Роль создать не удалось", "Попробуйте снова позже.")
		return
	}
	clubRole, err = ctx.Session.GuildRoleEdit(conf.Bot.GuildID, clubRole.ID, "[] "+clubName+"_club", 0, false, 0, false)
	if err != nil {
		modules.SendFail(ctx.Message.ChannelID, "Отредачить роль не удалось", "Попробуйте снова позже.")
		return
	}

	err = storage.Clubs.Create(tx, &storage.Club{
		OwnerID:   ctx.Message.Author.ID,
		ChannelID: clubChannel.ID,
		RoleID:    clubRole.ID,
		Title:     clubName,
	})
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		modules.SendFail(ctx.Message.ChannelID, "База крашнулась на закрытии", "Обратитесь к админам.")
		return
	}
}
