package clubs

import (
	conf "github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/modules"
	"github.com/bwmarrin/discordgo"
)

func CreateClubRole(title string, color int, mentionable bool) *discordgo.Role {
	return modules.CreateRole(title, color, false, 0, mentionable)
}

func CreateClubChannel(title string, roleID string, managers []string) *discordgo.Channel {
	perms := []*discordgo.PermissionOverwrite{
		{
			Type:  discordgo.PermissionOverwriteTypeRole,
			ID:    roleID,
			Allow: discordgo.PermissionViewChannel,
		},
		{
			Type: discordgo.PermissionOverwriteTypeRole,
			ID:   conf.Bot.GuildID,
			Deny: discordgo.PermissionViewChannel,
		},
	}
	for i := 0; i < len(managers); i++ {
		perms = append(perms, &discordgo.PermissionOverwrite{
			Type:  discordgo.PermissionOverwriteTypeMember,
			ID:    managers[i],
			Allow: discordgo.PermissionManageMessages | discordgo.PermissionMentionEveryone,
		})
	}

	return modules.CreateChannel(title, discordgo.ChannelTypeGuildText, "", conf.Bot.Channels.Category, perms)
}
