package xp

import (
	"github.com/SteMak/house-tyan/out"
	"github.com/SteMak/house-tyan/util"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

func addXpUsers(
	afk string,
	config *config,
	voiceStates []*discordgo.VoiceState,
	getMember func(guildID, userID string) (*discordgo.Member, error),
	addXpToUser func(userID string, xp int),
) {

	channals := generateChannals(voiceStates, getMember)

	for channelID, users := range channals {
		if skipUser(channelID, afk, len(users)) != nil {
			continue
		}

		roomBoost := countRoomBoost(len(users), config.VoiceFarm)

		for _, member := range users {
			if util.EqualAny(config.RoleHermit, member.Roles) {
				continue
			}

			addXpToUser(member.User.ID, config.VoiceFarm.XpForVoice*roomBoost)
		}
	}
}

func generateChannals(
	voiceStates []*discordgo.VoiceState,
	getMember func(guildID, userID string) (*discordgo.Member, error),
) map[string][]*discordgo.Member {

	channels := make(map[string][]*discordgo.Member)
	for _, state := range voiceStates {
		member, err := getMember(state.GuildID, state.UserID)
		if err != nil || member.User.Bot {
			out.Err(false, err)
			continue
		}
		if state.SelfDeaf || state.Deaf {
			continue
		}
		if state.SelfMute || state.Mute {
			continue
		}

		channels[state.ChannelID] = append(channels[state.ChannelID], member)
	}

	return channels
}

func skipUser(channelID, afk string, countMembers int) error {
	if channelID == afk {
		return errors.New("he is in afk")
	}

	if countMembers < 2 {
		return errors.New("there are a few people")
	}

	return nil
}

func countRoomBoost(countMembers int, voiceConf voiceFarm) int {
	roomBoost := 1
	if countMembers > voiceConf.MaxRoomBoost {
		roomBoost = roomBoost * voiceConf.MaxRoomBoost
	} else {
		roomBoost = roomBoost * countMembers
	}

	return roomBoost
}
