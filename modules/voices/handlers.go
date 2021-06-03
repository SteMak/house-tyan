package voices

import (
	"github.com/SteMak/house-tyan/cache"
	"github.com/SteMak/house-tyan/out"

	"github.com/bwmarrin/discordgo"
)

var (
	voiceStatesCache []*discordgo.VoiceState
	privateVoices    map[string]map[string]string = map[string]map[string]string{
		"692382001545871501": {
			"coreChannelID": "844141092180852736",
			"coreParentID":  "843056759287447572",
		},
	}
	err error
)

const (
	permConnect int64 = 1048576
	permView    int64 = 1024
	permManage  int64 = 16
)

func (bot *module) voiceHandler(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {

	if len(voiceStatesCache) == 0 && vs.VoiceState.ChannelID != "" {
		voiceStatesCache = append(voiceStatesCache, vs.VoiceState)

	} else if vs.BeforeUpdate == nil {
		voiceStatesCache = append(voiceStatesCache, vs.VoiceState)

	} else if vs.VoiceState.ChannelID != "" && vs.BeforeUpdate != nil {
		for i, state := range voiceStatesCache {
			if state.UserID == vs.UserID {
				voiceStatesCache[i] = vs.VoiceState
			}
		}

	} else if vs.BeforeUpdate != nil {
		for i, state := range voiceStatesCache {
			if state.UserID == vs.BeforeUpdate.UserID {
				voiceStatesCache = append(voiceStatesCache[:i], voiceStatesCache[i+1:]...)
			}
		}
	}

	privateVoice := privateVoices[vs.GuildID]
	coreCh, err := s.Channel(privateVoice["coreChannelID"])
	if err != nil {
		out.Err(true, err)
		return
	}

	cg, err := s.Channel(coreCh.ParentID)
	if err != nil {
		out.Err(true, err)
		return
	}

	if vs.ChannelID == privateVoice["coreChannelID"] {
		var permissionOverwrites []*discordgo.PermissionOverwrite
		var channelName string

		voice, err := cache.Voices.Get(vs.UserID)
		if err != nil {
			channelName = voice.Name
			permissionOverwrites = append(permissionOverwrites, voice.Permissions...)

		} else {
			member, err := s.GuildMember(vs.GuildID, vs.UserID)
			if err != nil {
				out.Err(true, err)
				return
			}

			channelName = member.User.Username
			permissionOverwrite := discordgo.PermissionOverwrite{
				ID:    vs.UserID,
				Type:  1,
				Deny:  0,
				Allow: permManage,
			}

			cgPerms := cg.PermissionOverwrites

			permissionOverwrites = []*discordgo.PermissionOverwrite{
				&permissionOverwrite,
			}

			permissionOverwrites = append(permissionOverwrites, cgPerms...)
		}

		data := discordgo.GuildChannelCreateData{
			Name:                 channelName,
			Type:                 2,
			Position:             99,
			ParentID:             privateVoice["coreParentID"],
			PermissionOverwrites: permissionOverwrites,
		}

		ch, err := s.GuildChannelCreateComplex(vs.GuildID, data)
		if err != nil {
			out.Err(true, err)
			return
		}

		err = s.GuildMemberMove(vs.GuildID, vs.UserID, &ch.ID)
		if err != nil {
			out.Err(true, err)
			return
		}
	}

	if vs.BeforeUpdate != nil && vs.BeforeUpdate.ChannelID != privateVoice["coreChannelID"] {
		channelID := vs.BeforeUpdate.ChannelID

		for _, voiceState := range voiceStatesCache {
			if channelID == voiceState.ChannelID {
				return
			}
		}

		ch, err := s.Channel(channelID)
		if err != nil {
			out.Err(true, err)
			return
		}

		if ch.ParentID == privateVoice["coreParentID"] {
			s.ChannelDelete(channelID)
		}
	}
}

func (bot *module) readyHandler(s *discordgo.Session, r *discordgo.Ready) {
	guilds := r.Guilds
	for _, guild := range guilds {
		if guild.ID == "" {
			voiceStatesCache = guild.VoiceStates
		}
	}
}
