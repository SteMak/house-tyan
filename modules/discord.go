package modules

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/pkg/errors"

	tyan "github.com/SteMak/house-tyan"

	"github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/messages"
	"github.com/SteMak/house-tyan/out"
)

func Send(channelID string, tplName string, data interface{}, beforeSend func(*messages.Message) error) *discordgo.Message {
	m, err := messages.Get(tplName, data)
	if err != nil {
		out.Err(true, errors.WithStack(err))
		return nil
	}

	if beforeSend != nil {
		err = beforeSend(m)
		if err != nil {
			out.Err(true, err)
			return nil
		}
	}

	message, err := session.ChannelMessageSendComplex(channelID, &m.MessageSend)
	if err != nil {
		out.Err(true, errors.WithStack(err))
		return nil
	}

	for _, reaction := range m.Reactions {
		err = session.MessageReactionAdd(message.ChannelID, message.ID, reaction)
		if err != nil {
			out.Err(true, errors.WithStack(err))
			return nil
		}
	}

	return message
}

func Edit(messageID, channelID string, tplName string, data interface{}, beforeSend func(*messages.Message) error) *discordgo.Message {
	m, err := messages.Get(tplName, data)
	if err != nil {
		out.Err(true, errors.WithStack(err))
		return nil
	}

	if beforeSend != nil {
		err = beforeSend(m)
		if err != nil {
			out.Err(true, err)
			return nil
		}
	}

	edit := new(discordgo.MessageEdit)
	edit.ID = messageID
	edit.Channel = channelID
	edit.SetContent(m.Content).
		SetEmbed(m.Embed)

	message, err := session.ChannelMessageEditComplex(edit)
	if err != nil {
		out.Err(true, errors.WithStack(err))
		return nil
	}

	err = session.MessageReactionsRemoveAll(channelID, messageID)
	if err != nil {
		out.Err(true, errors.WithStack(err))
		return nil
	}

	for _, reaction := range m.Reactions {
		err = session.MessageReactionAdd(message.ChannelID, message.ID, reaction)
		if err != nil {
			out.Err(true, errors.WithStack(err))
			return nil
		}
	}

	return message
}

func SendError(err error) {
	data := map[string]interface{}{
		"Timestamp": time.Now().UTC().Format(time.StampNano),
		"Version":   tyan.Vesion,
		"Message":   err.Error(),
	}

	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	if st, ok := err.(stackTracer); ok {
		stack := st.StackTrace()
		if len(stack) > 5 {
			stack = stack[:5]
		}
		data["Stack"] = fmt.Sprintf("%+v", stack)
	}

	m, err := messages.Get("error.xml", data)
	if err != nil {
		out.Err(false, err)
		return
	}

	_, err = session.ChannelMessageSendComplex(config.Bot.Channels.Errors, &m.MessageSend)
	if err != nil {
		out.Err(false, err)
	}
}

func SendLog(tplName string, data interface{}) {
	Send(config.Bot.Channels.Logs, tplName, data, nil)
}

func SendFail(channelID string, title, msg string) {
	data := map[string]interface{}{
		"Message": msg,
	}

	if title != "" {
		data["Title"] = title
	}

	Send(channelID, "fail.xml", data, nil)
}

func SendGood(channelID string, title, msg string) {
	data := map[string]interface{}{
		"Message": msg,
	}

	if title != "" {
		data["Title"] = title
	}

	Send(channelID, "good.xml", data, nil)
}

func CreateRole(title string, color int, hoist bool, perm int64, mention bool) *discordgo.Role {
	role, err := session.GuildRoleCreate(config.Bot.GuildID)
	if err != nil {
		out.Err(true, errors.WithStack(err))
		return nil
	}

	role, err = session.GuildRoleEdit(config.Bot.GuildID, role.ID, title, color, hoist, perm, mention)
	if err != nil {
		out.Err(true, errors.WithStack(err))
		return nil
	}

	return role
}

func CreateChannel(title string, channelType discordgo.ChannelType, topic, parent string, permissions []*discordgo.PermissionOverwrite) *discordgo.Channel {
	channel, err := session.GuildChannelCreateComplex(config.Bot.GuildID, discordgo.GuildChannelCreateData{
		Name:                 title,
		Type:                 channelType,
		Topic:                topic,
		ParentID:             parent,
		PermissionOverwrites: permissions,
	})
	if err != nil {
		out.Err(true, errors.WithStack(err))
		return nil
	}

	return channel
}
