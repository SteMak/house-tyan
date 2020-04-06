package triggers

import (
	"math/rand"
	"strings"

	"github.com/SteMak/house-tyan/out"

	"github.com/SteMak/house-tyan/cache"
	"github.com/bwmarrin/discordgo"
)

func (bot *module) triggerHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || m.Content == "" || strings.Contains(m.Content, " ") {
		return
	}

	// if m.Content == "test" {
	// 	err := cache.Triggers.Set(&cache.Trigger{
	// 		Name:    "trigger",
	// 		Answers: []string{"Bang, bang, bang, pull my Devil Trigger"},
	// 	})
	// 	if err != nil {
	// 		out.Debug(err)
	// 		return
	// 	}
	// }

	trigger, err := cache.Triggers.Get(strings.ToLower(m.Content))
	if err != nil {
		return
	}

	_, err = s.ChannelMessageSend(m.ChannelID, trigger.Answers[rand.Intn(len(trigger.Answers))])
	if err != nil {
		out.Err(true, err)
	}
}
