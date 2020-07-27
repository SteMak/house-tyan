package app

import (
	"fmt"
	"testing"

	"github.com/SteMak/house-tyan/dstype"

	"github.com/stretchr/testify/assert"

	"github.com/bwmarrin/discordgo"
)

func mess(s string) *discordgo.MessageCreate {
	return &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Content: s,
		},
	}
}

func TestCommand(t *testing.T) {
	mod := NewModule("testing", ".")
	mod.On("ping").Handle(func(ctx *Context) {
		fmt.Println("pong")
	})
	mod.Enable()

	app.onHandle(nil, mess(".testing ping    argv0 argv1 argv2"))
}

func TestGroup(t *testing.T) {
	mod := NewModule("testing", ".")
	mod.Enable()

	var result string

	clubgroup := mod.Group("club")

	clubgroup.Use(func(ctx *Context) {
		fmt.Println("middleware")
	})

	clubgroup.On().Handle(func(ctx *Context) {
		result = "club"
		fmt.Println(result)
	})

	clubgroup.On("create").Handle(func(ctx *Context) {
		result = "club create"
		fmt.Println(result)
	})

	clubgroup.On("lb").Handle(func(ctx *Context) {
		result = "club lb"
		fmt.Println(result)
	})

	app.onHandle(nil, mess(".club"))
	assert.Equal(t, "club", result)

	app.onHandle(nil, mess(".club lb"))
	assert.Equal(t, "club lb", result)

	app.onHandle(nil, mess(".club create"))
	assert.Equal(t, "club create", result)
}

func TestArgs(t *testing.T) {
	mod := NewModule("testing", ".")
	mod.On("ping").Handle(func(ctx *Context) {
		var (
			arg1 dstype.Grapheme
			arg2 int
			arg3 string
		)

		err := ctx.Scan(
			&arg1,
			&arg2,
			&arg3,
		)
		assert.NoError(t, err)
	})
	mod.Enable()

	app.onHandle(nil, mess(".ping 👍🏼 2 string arg"))
}
