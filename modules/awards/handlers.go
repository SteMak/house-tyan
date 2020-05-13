package awards

import (
	"strings"

	"github.com/SteMak/house-tyan/libs"

	"github.com/SteMak/house-tyan/storage"

	"github.com/SteMak/house-tyan/modules"
	"github.com/pkg/errors"

	"github.com/SteMak/house-tyan/out"
	"github.com/bwmarrin/discordgo"
)

func (bot *module) handlerBlankProcess(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == s.State.User.ID {
		return
	}

	if m.ChannelID != bot.config.Channels.Confirm {
		return
	}

	emoji := m.Emoji.Name
	if !strings.ContainsAny(emoji, "✅🇽") {
		return
	}

	blankID := m.MessageID

	award, err := storage.Awards.GetByBlankID(blankID)
	if err != nil {
		go out.Err(true, errors.WithStack(err))
		go modules.Send(m.ChannelID, "awards/fail_storage.xml", nil, nil)
		return
	}

	if award == nil {
		return
	}

	if award.Status != storage.AwardStatusUnknow {
		return
	}

	switch emoji {
	case "✅":
		err := storage.Awards.Accept(nil, blankID)
		if err != nil {
			go out.Err(true, errors.WithStack(err))
			go modules.Send(m.ChannelID, "awards/fail_storage.xml", nil, nil)
			return
		}

		rewards, err := award.Reawrds()
		if err != nil {
			go out.Err(true, errors.WithStack(err))
			go modules.Send(m.ChannelID, "awards/fail_storage.xml", nil, nil)
			return
		}

		var failed []string
		for _, reward := range rewards {
			if err := libs.Unb.AddToBalance(reward.UserID, int64(reward.Amount), award.Reason); err != nil {
				go out.Err(true, errors.WithStack(err))
				failed = append(failed, "<@"+reward.UserID+">")
				continue
			}

			if err := storage.Awards.SetPaid(nil, award.ID, reward.UserID); err != nil {
				go out.Err(true, errors.WithStack(err))
				failed = append(failed, "<@"+reward.UserID+">")
				continue
			}
		}

		data := make(map[string]interface{})
		if len(failed) > 0 {
			data["Failed"] = strings.Join(failed, ", ")
		}

		modules.Send(m.ChannelID, "awards/accepted.xml", data, nil)
	case "🇽":
		err := storage.Awards.Reject(nil, blankID)
		if err != nil {
			go out.Err(true, errors.WithStack(err))
			go modules.Send(m.ChannelID, "awards/fail_storage.xml", nil, nil)
			return
		}
		modules.Send(m.ChannelID, "awards/rejected.xml", nil, nil)
	}
}
