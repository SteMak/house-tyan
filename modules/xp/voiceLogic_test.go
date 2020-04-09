package xp

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

type addXpUsersDataType struct {
	afkID     string
	config    *config
	states    []*discordgo.VoiceState
	getMember func(guildID, userID string) (*discordgo.Member, error)

	Expected map[string]int
}

var addXpUsersData = []addXpUsersDataType{
	{
		afkID: "",
		config: &config{
			VoiceFarm:  voiceFarm{},
			RoleHermit: "",
		},
		states:    []*discordgo.VoiceState{},
		getMember: func(guildID, userID string) (*discordgo.Member, error) { return nil, nil },

		Expected: map[string]int{},
	},
	{
		afkID: "afkID",
		config: &config{
			VoiceFarm: voiceFarm{
				MaxRoomBoost: 1,
				XpForVoice:   15,
			},
			RoleHermit: "roleHermit",
		},
		states: []*discordgo.VoiceState{
			&discordgo.VoiceState{
				ChannelID: "firstChannel",
				GuildID:   "guildID",
				Deaf:      false,
				SelfDeaf:  false,
				Mute:      false,
				SelfMute:  false,
				UserID:    "firstUser",
			},
			&discordgo.VoiceState{
				ChannelID: "firstChannel",
				GuildID:   "guildID",
				Deaf:      false,
				SelfDeaf:  false,
				Mute:      false,
				SelfMute:  false,
				UserID:    "secondUser",
			},
		},
		getMember: func(guildID, userID string) (*discordgo.Member, error) {
			return &discordgo.Member{
				User: &discordgo.User{
					Bot: false,
					ID:  userID,
				},
			}, nil
		},

		Expected: map[string]int{
			"firstUser":  15,
			"secondUser": 15,
		},
	},
	{
		afkID: "afkID",
		config: &config{
			VoiceFarm: voiceFarm{
				MaxRoomBoost: 3,
				XpForVoice:   15,
			},
			RoleHermit: "roleHermit",
		},
		states: []*discordgo.VoiceState{
			&discordgo.VoiceState{
				ChannelID: "firstChannel",
				GuildID:   "guildID",
				Deaf:      false,
				SelfDeaf:  false,
				Mute:      false,
				SelfMute:  false,
				UserID:    "firstUser",
			},
			&discordgo.VoiceState{
				ChannelID: "firstChannel",
				GuildID:   "guildID",
				Deaf:      false,
				SelfDeaf:  false,
				Mute:      false,
				SelfMute:  false,
				UserID:    "secondUser",
			},
		},
		getMember: func(guildID, userID string) (*discordgo.Member, error) {
			return &discordgo.Member{
				User: &discordgo.User{
					Bot: userID == "firstUser",
					ID:  userID,
				},
			}, nil
		},

		Expected: map[string]int{},
	},
	{
		afkID: "afkID",
		config: &config{
			VoiceFarm: voiceFarm{
				MaxRoomBoost: 3,
				XpForVoice:   15,
			},
			RoleHermit: "roleHermit",
		},
		states: []*discordgo.VoiceState{
			&discordgo.VoiceState{
				ChannelID: "afkID",
				GuildID:   "guildID",
				Deaf:      false,
				SelfDeaf:  false,
				Mute:      false,
				SelfMute:  false,
				UserID:    "firstUser",
			},
		},
		getMember: func(guildID, userID string) (*discordgo.Member, error) {
			return &discordgo.Member{
				User: &discordgo.User{
					Bot: false,
					ID:  userID,
				},
			}, nil
		},

		Expected: map[string]int{},
	},
	{
		afkID: "afkID",
		config: &config{
			VoiceFarm: voiceFarm{
				MaxRoomBoost: 3,
				XpForVoice:   15,
			},
			RoleHermit: "roleHermit",
		},
		states: []*discordgo.VoiceState{
			&discordgo.VoiceState{
				ChannelID: "firstChannel",
				GuildID:   "guildID",
				Deaf:      false,
				SelfDeaf:  false,
				Mute:      false,
				SelfMute:  false,
				UserID:    "firstUser",
			},
			&discordgo.VoiceState{
				ChannelID: "firstChannel",
				GuildID:   "guildID",
				Deaf:      false,
				SelfDeaf:  false,
				Mute:      false,
				SelfMute:  false,
				UserID:    "secondUser",
			},
		},
		getMember: func(guildID, userID string) (*discordgo.Member, error) {
			var roles []string
			if userID == "secondUser" {
				roles = []string{
					"onerole",
					"roleHermit",
					"another",
				}
			}
			return &discordgo.Member{
				User: &discordgo.User{
					Bot: false,
					ID:  userID,
				},
				Roles: roles,
			}, nil
		},

		Expected: map[string]int{
			"firstUser": 30,
		},
	},
	{
		afkID: "afkID",
		config: &config{
			VoiceFarm: voiceFarm{
				MaxRoomBoost: 3,
				XpForVoice:   15,
			},
			RoleHermit: "roleHermit",
		},
		states: []*discordgo.VoiceState{
			&discordgo.VoiceState{
				ChannelID: "firstChannel",
				GuildID:   "guildID",
				Deaf:      true,
				SelfDeaf:  false,
				Mute:      false,
				SelfMute:  false,
				UserID:    "firstUser",
			},
			&discordgo.VoiceState{
				ChannelID: "firstChannel",
				GuildID:   "guildID",
				Deaf:      false,
				SelfDeaf:  false,
				Mute:      false,
				SelfMute:  true,
				UserID:    "secondUser",
			},
		},
		getMember: func(guildID, userID string) (*discordgo.Member, error) {
			return &discordgo.Member{
				User: &discordgo.User{
					Bot: false,
					ID:  userID,
				},
			}, nil
		},

		Expected: map[string]int{},
	},
	{
		afkID: "afkID",
		config: &config{
			VoiceFarm: voiceFarm{
				MaxRoomBoost: 3,
				XpForVoice:   15,
			},
			RoleHermit: "roleHermit",
		},
		states: []*discordgo.VoiceState{
			&discordgo.VoiceState{
				ChannelID: "firstChannel",
				GuildID:   "guildID",
				Deaf:      false,
				SelfDeaf:  false,
				Mute:      false,
				SelfMute:  false,
				UserID:    "firstUser",
			},
			&discordgo.VoiceState{
				ChannelID: "firstChannel",
				GuildID:   "guildID",
				Deaf:      false,
				SelfDeaf:  false,
				Mute:      false,
				SelfMute:  false,
				UserID:    "secondUser",
			},
			&discordgo.VoiceState{
				ChannelID: "secondChannel",
				GuildID:   "guildID",
				Deaf:      false,
				SelfDeaf:  false,
				Mute:      false,
				SelfMute:  false,
				UserID:    "thirdUser",
			},
			&discordgo.VoiceState{
				ChannelID: "secondChannel",
				GuildID:   "guildID",
				Deaf:      false,
				SelfDeaf:  false,
				Mute:      false,
				SelfMute:  false,
				UserID:    "fourthUser",
			},
			&discordgo.VoiceState{
				ChannelID: "secondChannel",
				GuildID:   "guildID",
				Deaf:      false,
				SelfDeaf:  false,
				Mute:      false,
				SelfMute:  false,
				UserID:    "fifthUser",
			},
		},
		getMember: func(guildID, userID string) (*discordgo.Member, error) {
			return &discordgo.Member{
				User: &discordgo.User{
					Bot: false,
					ID:  userID,
				},
			}, nil
		},

		Expected: map[string]int{
			"firstUser":  30,
			"secondUser": 30,
			"thirdUser":  45,
			"fourthUser": 45,
			"fifthUser":  45,
		},
	},
}

func TestAddXpUsers(t *testing.T) {
	for _, data := range addXpUsersData {
		actual := make(map[string]int)

		addXpUsers(
			data.afkID,
			data.config,
			data.states,
			data.getMember,
			func(userID string, xp int) {
				actual[userID] = xp
			},
		)

		if len(actual) != len(data.Expected) {
			t.Error(
				"Expected length:", len(data.Expected),
				"Actual length:", len(actual),
			)
		}

		for i, value := range actual {
			if data.Expected[i] != value {
				t.Error(
					"User:", i,
					"Expected:", data.Expected[i],
					"Actual:", value,
				)
			}
		}
	}
}
