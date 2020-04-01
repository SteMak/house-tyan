package logic

import (
	"fmt"

	"github.com/SteMak/ani-helper/workerTools/config"
	"github.com/SteMak/ani-helper/workerTools/database"
	"github.com/bwmarrin/discordgo"
)

var (
	chMonitorWriters []simplifiedUser
)

type simplifiedUser struct {
	id     string
	strify string
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.GuildID == config.GdHouseID {
		if m.Content == "R U TYT?" && m.ChannelID == config.ChForLogsID {
			s.ChannelMessageSend(m.ChannelID, "I TYT KUSHAU!")
			return
		}

		if m.ChannelID == config.ChForBumpSiupID {
			if len(chMonitorWriters) >= 30 {
				chMonitorWriters = chMonitorWriters[1:]
			}

			chMonitorWriters = append(chMonitorWriters, simplifiedUser{
				id:     m.Author.ID,
				strify: m.Author.String(),
			})
		}

		if m.ChannelID == config.ChForBumpSiupID && len(m.Embeds) > 0 {
			detectBumpSiup(s, m)
			return
		}

		isRequest, request := checkRequest(m.Content)
		if isRequest {
			member, err := s.GuildMember(config.GdHouseID, m.Author.ID)
			if err != nil {
				return
			}
			if len(member.Roles) > 0 && hasRole(member, config.RoRequestMakerID) {
				detectRequest(s, m.ChannelID, "-запрос"+request)
				return
			}
		}
	}
}

func reactionHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.ChannelID == config.ChForRequestID && r.UserID == config.UsConfirmatorID {
		item, err := database.Records.Record(r.MessageID)
		if err != nil {
			return
		}

		fmt.Println("FOUND "+item.EmbedID+" reation added", r.Emoji.Name)

		emojiOnRequest(s, r, item)

		fmt.Println("GUILD " + item.EmbedID + " request processed successfuly")
	}
}
