package voices

import (
	"github.com/SteMak/house-tyan/out"

	"github.com/bwmarrin/discordgo"
)

var (
	err error
)

func (bot *module) voiceHandler(s *discordgo.Session, vs *discordgo.VoicrStateUpdate) {
	privateVoices := map[string]map[string]string{
		"692382001545871501": {
			"coreChannelID": "844141092180852736",
			"coreParentID": "843056759287447572",
			},
		}

	privateVoice := privateVoices[vs.GuildID]

	if vs.ChannelID == privateVoice["coreChannelID"] {
		voiceDefaultData := discordgo.ChannelEdit {
			Position: 99,
			ParentID: privateVoice["coreParentID"],
		}

		member, err := s.GuildMember(vs.GuildID, vs.UserID)
		if err != nil {
			out.Err(true, err)
			return
		}

		permInitData := discordgo.PermissonOverwrite{
			ID: vs.UserID,
			Type: 1,
			Deny: 0,
			Allow: 16,
		}

		permissionOverwrites := []*discordgo.PermissonOverwrite{
			&permInitData,
		}

		data := discordgo.GuildChannelCreateData{
			Name: member.User.Username,
			Type: 2,
			Position: 99,
			PermissionOverwrites: permissionOverwrites,
			ParentID: privateVoice["coreParentID"],
		}

		channel, err := s.GuildChannelCreateComplex(vs.GuildID, data)
		if err != nil {
			out.Err(true, err)
			return
		}

		err = s.GuildMemberMove(vs.GuildID, vs.UserID, &channel.ID)
		if err != nil {
			out.Err(true, err)
			return
		}
	}
	if vs.BeforeUpdate != nil && vs.BeforeUpdate.ChannelID != privateVoice["coreChannelID"] {
		channelID := vs.BeforeUpdate.ChannelID
		channel, err := s.Channel(channelID)
		if err != nil {
			out.Err(true, err)
			return
		}

		if channel.ParentID != privateVoice["coreParentID"] {
			return
		}

		guild, err := s.Guild(vs.GuildID)
		if err != nil {
			out.Err(true, err)
			return
		}

		for _, voiceState := range guild.VoiceStates {
			if channelID == voiceState.ChannelID {
				return
			}
		}

		s.ChannelDelete(channelID)

	}
}
