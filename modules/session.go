package modules

import (
	tyan "github.com/SteMak/house-tyan"
	"github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/out"
	"github.com/bwmarrin/discordgo"
)

var (
	session *discordgo.Session
)

func authentificate() {
	out.Infoln("\nAuthentifications...")
	s, err := discordgo.New(config.Session.Token)
	if err != nil {
		out.Fatal(err)
	}
	session = s

	session.StateEnabled = true
	session.SyncEvents = false
	session.Debug = tyan.GlobalCtx.GlobalBool("debug")

	session.AddHandler(onReady)

	if err := session.Open(); err != nil {
		out.Fatal(err)
	}
}
