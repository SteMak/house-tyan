package modules

import (
	tyan "github.com/SteMak/house-tyan"
	"github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/out"
	"github.com/SteMak/house-tyan/util"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var (
	session *discordgo.Session
	log     *logrus.Logger
)

func authentificate() {
	out.Infoln("\nAuthentifications...")
	s, err := discordgo.New(config.Session.Token)
	if err != nil {
		out.Fatal(err)
	}
	session = s

	log, err = util.Logger(config.Session.Log)
	if err != nil {
		out.Fatal(err)
	}

	switch log.Level {
	case logrus.ErrorLevel:
		session.LogLevel = discordgo.LogError
	case logrus.WarnLevel:
		session.LogLevel = discordgo.LogWarning
	case logrus.InfoLevel:
		session.LogLevel = discordgo.LogInformational
	case logrus.DebugLevel:
		session.LogLevel = discordgo.LogDebug
	}

	discordgo.Logger = func(msgL, caller int, format string, a ...interface{}) {
		switch msgL {
		case discordgo.LogError:
			log.Errorf(format, a...)
		case discordgo.LogWarning:
			log.Warnf(format, a...)
		case discordgo.LogInformational:
			log.Infof(format, a...)
		case discordgo.LogDebug:
			log.Debugf(format, a...)
		}
	}

	session.StateEnabled = true
	session.SyncEvents = false

	if tyan.GlobalCtx != nil {
		session.Debug = tyan.GlobalCtx.GlobalBool("debug")
	}

	session.AddHandler(onReady)

	if err := session.Open(); err != nil {
		out.Fatal(err)
	}
}
