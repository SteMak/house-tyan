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
		data["Stack"] = fmt.Sprintf("%+v", st.StackTrace())
	}

	m, err := messages.Get("main/error.xml", data)
	if err != nil {
		out.Err(false, err)
		return
	}

	_, err = session.ChannelMessageSendComplex(*config.Bot.ErrorsChannel, &m.MessageSend)
	if err != nil {
		out.Err(false, err)
	}
}

func SendLog(tplName string, data interface{}) {
	Send(*config.Bot.LogChannel, tplName, data, nil)
}
