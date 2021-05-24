package clubs

import (
	"errors"

	"github.com/SteMak/house-tyan/modules"
	"github.com/SteMak/house-tyan/out"
	"github.com/SteMak/house-tyan/storage"
	"github.com/bwmarrin/discordgo"
)

var ErrClubNoFound = errors.New("клуб не найден")

func (bot *module) acceptHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID || r.Emoji.Name != "✅" && r.Emoji.Name != "❎" {
		return
	}

	invite, err := storage.Clubs.GetInvite(r.MessageID)
	if err != nil {
		go out.Err(true, err)
		go log.Error(err)
		return
	}

	if invite == nil {
		return
	}

	club, err := storage.Clubs.GetClubByID(invite.ClubID)
	if err != nil {
		go out.Err(true, err)
		go log.Error(err)
		return
	}

	if club == nil {
		go out.Err(true, ErrClubNoFound)
		go log.Error(ErrClubNoFound)
		return
	}

	if r.Emoji.Name == "✅" {
		tx, err := storage.Tx()
		if err != nil {
			go out.Err(true, err)
			go modules.SendFail(r.ChannelID, "База крашнулась на открытии", "Попробуйте снова позже.")
			go log.Error(err)
			return
		}

		err = club.AddMember(tx, invite.UserID)
		if err != nil {
			go out.Err(true, err)
			go modules.SendFail(r.ChannelID, "База крашнулась на комите", "Попробуйте снова позже.")
			go log.Error(err)
			tx.Rollback()
			return
		}

		err = tx.Commit()
		if err != nil {
			go out.Err(true, err)
			go modules.SendFail(r.ChannelID, "База крашнулась на закрытии", "Попробуйте снова позже.")
			go log.Error(err)
			return
		}
	}

	err = storage.Clubs.DeleteUserInvites(invite.UserID)
	if err != nil {
		go out.Err(true, err)
		go log.Error(err)
		return
	}

	err = s.ChannelMessageDelete(r.ChannelID, r.MessageID)
	if err != nil {
		go out.Err(true, err)
		go log.Error(err)
	}
}
