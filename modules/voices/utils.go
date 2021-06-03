package voices

import (
	"strings"

	"github.com/SteMak/house-tyan/out"

	//	"github.com/SteMak/house-tyan/cache"
	"github.com/SteMak/house-tyan/libs/dgutils"
	"github.com/bwmarrin/discordgo"
)

func getObjID(str string) (string, discordgo.PermissionOverwriteType) {
	var objID string
	var objType discordgo.PermissionOverwriteType
	if strings.HasPrefix(str, "<@") {
		objID = strings.TrimPrefix(str, "<@")
	}

	if strings.HasPrefix(objID, "&") {
		objID = strings.TrimPrefix(objID, "&")
		objType = discordgo.PermissionOverwriteTypeRole
	} else {
		objType = discordgo.PermissionOverwriteTypeMember
	}

	objID = strings.TrimPrefix(objID, "!")

	objID = strings.TrimSuffix(objID, ">")

	return objID, objType
}

func setPermissions(ctx *dgutils.MessageContext, permID int64, isAllow bool) {
	var channel *discordgo.Channel
	var args = ctx.Args

	for _, state := range voiceStatesCache {
		if state.UserID == ctx.Message.Author.ID {
			channel, err = ctx.Session.Channel(state.ChannelID)
			if err != nil {
				out.Err(true, err)
				return
			}
		}
	}

	isEveryone := false

	if len(ctx.Args) < 1 {
		args = append(args, "<@&"+ctx.Message.GuildID+">")
		isEveryone = true
	}

	permID1 := permID
	permID2 := int64(0)

	perms := channel.PermissionOverwrites

	for _, arg := range args {
		objID, objType := getObjID(arg)

		for _, perm := range perms {
			perm1 := perm.Deny
			perm2 := perm.Allow
			if isAllow {
				perm1, perm2 = perm2, perm1
			}
			if perm.ID == objID {
				if permID1 == permID1&perm2 {
					permID2 = perm2 - permID1
				} else {
					permID2 = perm2 - permID2
				}
				permID1 = perm1 + permID1
				break
			}
		}
		if isAllow {
			permID1, permID2 = permID2, permID1
		}
		perm := discordgo.PermissionOverwrite{
			ID:    objID,
			Type:  objType,
			Deny:  permID1,
			Allow: permID2,
		}

		perms = append(perms, &perm)

	}

	data := discordgo.ChannelEdit{
		Position:             channel.Position,
		PermissionOverwrites: perms,
	}

	_, err = ctx.Session.ChannelEditComplex(channel.ID, &data)
	if err != nil {
		out.Err(true, err)
		return
	}

	content := "Успешно [value] для: "
	value := ""
	if isEveryone {
		args = []string{"всех"}
	}

	if permID == permConnect && isAllow {
		value = "открыто"
	} else if permID == permConnect && !isAllow {
		value = "закрыто"
	} else if permID == permView && isAllow {
		value = "показано"
	} else if permID == permView && !isAllow {
		value = "скрыто"
	}

	content = strings.Replace(content, "[value]", value, 1) + strings.Join(args, ", ")

	ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, content)

}

func getInfo(s *discordgo.Session, channelID string) discordgo.MessageEmbed {
	ch, err := s.Channel(channelID)
	if err != nil {
		out.Err(true, err)
	}

	perms := ch.PermissionOverwrites
	permsConnect := []string{}
	permsNoConnect := []string{}
	permsView := []string{}
	permsNoView := []string{}
	ownerID := ""

	for _, perm := range perms {
		if permManage == permManage&perm.Allow {
			ownerID = perm.ID
			continue
		}
		if permConnect == permConnect&perm.Allow {
			permsConnect = append(permsConnect, mention(perm.ID, perm.Type))
		}
		if permView == permView&perm.Allow {
			permsView = append(permsView, mention(perm.ID, perm.Type))
		}
		if permConnect == permConnect&perm.Deny {
			permsNoConnect = append(permsNoConnect, mention(perm.ID, perm.Type))
		}
		if permView == permView&perm.Deny {
			permsNoView = append(permsNoView, mention(perm.ID, perm.Type))
		}
	}

	embedFields := []*discordgo.MessageEmbedField{}

	ownerInfoField := discordgo.MessageEmbedField{
		Name:  "Владелец",
		Value: "• " + mention(ownerID, 1) + " \n • Канал [ <#" + channelID + "> ]",
	}

	embedFields = append(embedFields, &ownerInfoField)

	if len(permsConnect) != 0 {
		embed := discordgo.MessageEmbedField{
			Name:   "Доступно",
			Value:  strings.Join(permsConnect, ", "),
			Inline: true,
		}
		embedFields = append(embedFields, &embed)
	}
	if len(permsNoConnect) != 0 {
		embed := discordgo.MessageEmbedField{
			Name:   "Закрыто",
			Value:  strings.Join(permsNoConnect, ", "),
			Inline: true,
		}
		embedFields = append(embedFields, &embed)
	}
	if len(permsView) != 0 {
		embed := discordgo.MessageEmbedField{
			Name:   "Видно",
			Value:  strings.Join(permsView, ", "),
			Inline: true,
		}
		embedFields = append(embedFields, &embed)
	}
	if len(permsNoView) != 0 {
		embed := discordgo.MessageEmbedField{
			Name:   "Скрыто",
			Value:  strings.Join(permsNoView, ", "),
			Inline: true,
		}
		embedFields = append(embedFields, &embed)
	}

	embedInfo := discordgo.MessageEmbed{
		Title:  "Информация о приватном канале",
		Fields: embedFields,
	}

	return embedInfo

}

func mention(id string, objType discordgo.PermissionOverwriteType) string {
	if objType == 0 {
		return "<@&" + id + ">"
	} else if objType == 1 {
		return "<@" + id + ">"
	}

	return ""
}
