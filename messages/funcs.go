package messages

import (
	"html/template"

	"github.com/bwmarrin/discordgo"
)

var (
	funcs = make(template.FuncMap)
)

func init() {
	funcs["mention"] = func(user discordgo.User) string {
		return user.Mention()
	}
}
