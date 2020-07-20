package app

import (
	"fmt"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestCommand(t *testing.T) {
	testmod := NewModule("testing", ".")
	testmod.On("ping").Handle(func(ctx *Context) error {
		fmt.Println("pong")
		return nil
	})
	testmod.Enable()

	app.onHandle(nil, &discordgo.MessageCreate{Message: &discordgo.Message{Content: ".testing ping"}})
}
