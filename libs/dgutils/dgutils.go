package dgutils

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Discord struct {
	Session *discordgo.Session

	OnReady   func()
	OnMessage func(*discordgo.Session, *discordgo.Message)

	Prefix   string
	Commands map[string]interface{}
	Usage    func(map[string]interface{}) *discordgo.MessageSend

	menus map[string]*Menu

	stopHandlers []func()
	ErrorHandler func(error)
}

func New(token string) (*Discord, error) {
	session, err := discordgo.New(token)
	if err != nil {
		return nil, err
	}

	return &Discord{
		Session: session,

		Prefix: "!",
		Usage:  Usage,

		menus:        make(map[string]*Menu),
		ErrorHandler: func(err error) { fmt.Println("error", err) },
	}, nil
}

func (discord *Discord) Start(s *discordgo.Session) {
	discord.Session = s

	discord.stopHandlers = []func(){
		discord.Session.AddHandler(discord.messageCreate),
		discord.Session.AddHandler(discord.menuHandler),
	}
}

func (discord *Discord) Stop() {
	for _, stopHandler := range discord.stopHandlers {
		stopHandler()
	}
}

func (discord *Discord) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if discord.OnMessage != nil {
		discord.OnMessage(s, m.Message)
	}

	if m.Author.Bot || !strings.HasPrefix(m.Content, discord.Prefix) || m.Author.ID == s.State.User.ID {
		return
	}

	str := strings.TrimPrefix(m.Content, discord.Prefix)
	str = strings.Join(strings.Fields(str), " ")

	command, args := findCommand(str, discord.Commands)
	if command == nil {
		if discord.Usage != nil {
			_, err := s.ChannelMessageSendComplex(m.ChannelID, discord.Usage(discord.Commands))
			if err != nil {
				discord.ErrorHandler(err)
			}
		}
		return
	}

	ctx := newContext(discord.Session,
		m.Message,
		args,
		append(command.Handlers, command.Function),
	)
	ctx.Next()
}
