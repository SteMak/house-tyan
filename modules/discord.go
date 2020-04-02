package modules

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	tyan "github.com/SteMak/house-tyan"

	"github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/messages"
	"github.com/SteMak/house-tyan/out"
)

func Send(channelID string, tplName string, data interface{}) (err error) {
	m, err := messages.Get(tplName, data)
	if err != nil {
		out.Err(true, errors.WithStack(err))
		return
	}

	message, err := session.ChannelMessageSendComplex(channelID, &m.MessageSend)
	if err != nil {
		out.Err(true, errors.WithStack(err))
		return
	}

	for _, reaction := range m.Reactions {
		err = session.MessageReactionAdd(message.ChannelID, message.ID, reaction)
		if err != nil {
			out.Err(true, errors.WithStack(err))
			return
		}
	}

	return
}

func SendError(err error) {
	data := map[string]interface{}{
		"Timestamp": time.Now().UTC().Format(time.StampNano),
		"Version":   tyan.Vesion,
		"Error":     fmt.Sprintf("%+v", err),
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

func SendLog(tplName string, data interface{}) (err error) {
	return Send(*config.Bot.LogChannel, tplName, data)
}
