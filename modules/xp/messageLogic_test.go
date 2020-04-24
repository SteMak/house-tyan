package xp

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

type howMuchXpDataType struct {
	Content string
	MesConf messageFarm
	NewXp   int
}

type xpMessageChecksDataType struct {
	ChannelID string
	Config    config
	GuildID   string
	AuthorID  string
	GetMember func(string, string) (*discordgo.Member, error)

	Expected error
}

var howMuchXpData = []howMuchXpDataType{
	howMuchXpDataType{
		Content: "",
		MesConf: messageFarm{
			XpForChar:  0.34,
			XpForRune:  1,
			XpForEmpty: 5,
		},
		NewXp: 5,
	},
	howMuchXpDataType{ 
		Content: "dghd<:AH_Whaat:579709315024158720>hdfh<a:AH_A_SayumiPeek:641710287002796032>f h<@306370108161392653>  <@!306370108161392653>f😪 f",
		MesConf: messageFarm{
			XpForChar:  1,
			XpForRune:  100,
			XpForEmpty: 10000,
		},
		NewXp: 515,
	},
	howMuchXpDataType{
		Content: "dghd<:AH_Whaat:5797094158720>hdfh<aa:AH_A_SayumiPeek:641710287002796032>fh<@3063701061392653>  <@!?306370108161392653>f😪 f",
		MesConf: messageFarm{
			XpForChar:  1,
			XpForRune:  100,
			XpForEmpty: 10000,
		},
		NewXp: 233,
	},
	howMuchXpDataType{
		Content: "*рандомный русский текст*",
		MesConf: messageFarm{
			XpForChar:  0.2,
			XpForRune:  3,
			XpForEmpty: 5,
		},
		NewXp: 5,
	},
	howMuchXpDataType{
		Content: "💏 👨‍❤️‍💋‍👨 👩‍❤️‍💋‍👩 💑 👨‍👨‍👦 👨‍👨‍👧 👨‍👨‍👧‍👦 👨‍👨‍👦‍👦",
		MesConf: messageFarm{
			XpForChar:  1,
			XpForRune:  100,
			XpForEmpty: 10000,
		},
		NewXp: 807,
	},
	howMuchXpDataType{
		Content: "🖇",
		MesConf: messageFarm{
			XpForChar:  0,
			XpForRune:  1,
			XpForEmpty: 0,
		},
		NewXp: 1,
	},
	howMuchXpDataType{
		Content: " !\"#$%&'()*+,-. /0123456789:;<= >?@[]^_`{|}",
		MesConf: messageFarm{
			XpForChar:  1,
			XpForRune:  10000000000,
			XpForEmpty: 5000000000000,
		},
		NewXp: 43,
	},
	howMuchXpDataType{
		Content: "ABCDEFGHIJKLMNO PQRSTUVWXYZabcd efghijklmnopqrs tuvwxyz",
		MesConf: messageFarm{
			XpForChar:  1,
			XpForRune:  10000000000,
			XpForEmpty: 5000000000000,
		},
		NewXp: 55,
	},
	howMuchXpDataType{
		Content: "ЙЦУКЕНГШЩЗХЪФЫВ АПРОЛДЖЭЯЧСМИТЬ БЮІЇЄҐйцукенгшщ зхъфывапролджэя чсмитьбюіїєґ",
		MesConf: messageFarm{
			XpForChar:  1,
			XpForRune:  10000000000,
			XpForEmpty: 5000000000000,
		},
		NewXp: 76,
	},
	howMuchXpDataType{
		Content: "google.com http://google.com https://google.com http://www.google.com",
		MesConf: messageFarm{
			XpForChar:  1,
			XpForRune:  100,
			XpForEmpty: 5000000000000,
		},
		NewXp: 403,
	},
	howMuchXpDataType{
		Content: "******************         *************************",
		MesConf: messageFarm{
			XpForChar:  1,
			XpForRune:  100,
			XpForEmpty: 5000000000000,
		},
		NewXp: 3,
	},
}

var xpMessageChecksData = []xpMessageChecksDataType{
	xpMessageChecksDataType{
		ChannelID: "",
		Config:    config{},
		GuildID:   "",
		AuthorID:  "",
		GetMember: func(guildID string, authorID string) (*discordgo.Member, error) {
			return &discordgo.Member{}, nil
		},

		Expected: errors.New("isn't a message fam channel"),
	},
	xpMessageChecksDataType{
		ChannelID: "thirdChannel",
		Config: config{
			MessageFarm: messageFarm{
				Channels: []string{
					"firstChannel",
					"seconChannel",
				},
			},
		},
		GuildID:  "guildID",
		AuthorID: "firstUser",
		GetMember: func(guildID string, authorID string) (*discordgo.Member, error) {
			return &discordgo.Member{}, nil
		},

		Expected: errors.New("isn't a message fam channel"),
	},
	xpMessageChecksDataType{
		ChannelID: "firstChannel",
		Config: config{
			MessageFarm: messageFarm{
				Channels: []string{
					"firstChannel",
					"seconChannel",
				},
			},
		},
		GuildID:  "guildID",
		AuthorID: "firstUser",
		GetMember: func(guildID string, authorID string) (*discordgo.Member, error) {
			return &discordgo.Member{
				User: &discordgo.User{
					Bot: true,
				},
			}, nil
		},

		Expected: errors.New("maybe bot"),
	},
	xpMessageChecksDataType{
		ChannelID: "firstChannel",
		Config: config{
			MessageFarm: messageFarm{
				Channels: []string{
					"firstChannel",
					"seconChannel",
				},
			},
			RoleHermit: "roleHermit",
		},
		GuildID:  "guildID",
		AuthorID: "firstUser",
		GetMember: func(guildID string, authorID string) (*discordgo.Member, error) {
			return &discordgo.Member{
				User: &discordgo.User{
					Bot: false,
				},
				Roles: []string{
					"onerole",
					"roleHermit",
					"another",
				},
			}, nil
		},

		Expected: errors.New("it's a hermit"),
	},
	xpMessageChecksDataType{
		ChannelID: "firstChannel",
		Config: config{
			MessageFarm: messageFarm{
				Channels: []string{
					"firstChannel",
					"seconChannel",
				},
			},
			RoleHermit: "roleHermit",
		},
		GuildID:  "guildID",
		AuthorID: "firstUser",
		GetMember: func(guildID string, authorID string) (*discordgo.Member, error) {
			return &discordgo.Member{
				User: &discordgo.User{
					Bot: false,
				},
				Roles: []string{
					"onerole",
					"another",
				},
			}, nil
		},

		Expected: nil,
	},
}

func TestHowMuchXp(t *testing.T) {
	for _, data := range howMuchXpData {
		replyXp := howMuchXp(data.Content, data.MesConf)
		if replyXp != data.NewXp {
			t.Error(
				"Content:", data.Content,
				"Expected xp:", data.NewXp,
				"Replied xp:", replyXp,
			)
		}
	}
}

func TestXpMessageChecks(t *testing.T) {
	for _, data := range xpMessageChecksData {
		actualErr := xpMessageChecks(
			data.ChannelID,
			data.Config,
			data.GuildID,
			data.AuthorID,
			data.GetMember,
		)
		if actualErr != data.Expected && actualErr.Error() != data.Expected.Error() {
			t.Error(
				"ChannelID:", data.ChannelID,
				"Expected error:", data.Expected,
				"Replied error:", actualErr,
			)
		}
	}
}
