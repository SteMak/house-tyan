package clubs

import (
	"github.com/sirupsen/logrus"

	"github.com/SteMak/house-tyan/libs/dgutils"

	"github.com/bwmarrin/discordgo"
)

var log *logrus.Logger

type module struct {
	session *discordgo.Session
	config  config

	running bool

	cmds         *dgutils.Discord
	stopHandlers []func()
}
