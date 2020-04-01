package modules

import (
	"runtime/debug"
	"time"

	tyan "github.com/SteMak/house-tyan"

	"github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/messages"
	"github.com/SteMak/house-tyan/out"
)

func Send(channelID string, tplName string, data interface{}) (err error) {
	m, err := messages.Get(tplName, data)
	if err != nil {
		out.Err(true, err)
		return
	}

	_, err = session.ChannelMessageSendComplex(channelID, m)
	if err != nil {
		out.Err(true, err)
		return
	}
	return
}

func SendError(msg string) {
	data := map[string]interface{}{
		"Timestamp": time.Now().UTC().Format(time.StampNano),
		"Version":   tyan.Vesion,
		"Message":   msg,
		"Stack":     string(debug.Stack()),
	}

	m, err := messages.Get("main/error.xml", data)
	if err != nil {
		out.Err(false, err)
		return
	}

	_, err = session.ChannelMessageSendComplex(*config.Bot.ErrorsChannel, m)
	if err != nil {
		out.Err(false, err)
	}
}
