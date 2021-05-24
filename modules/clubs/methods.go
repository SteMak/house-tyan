package clubs

import (
	"github.com/SteMak/house-tyan/storage"
	"github.com/SteMak/house-tyan/util"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type club struct {
	storage.Club
}

func (c *club) onClubVerified(tx *sqlx.Tx) error {
	c.CreateRole(tx)
	c.CreateChannel(tx)
	// c.PostCard()

	return nil
}

func (c *club) CreateRole(tx *sqlx.Tx) error {
	color, _ := util.Hex2RGB(util.Hex(c.RoleColor))
	role := CreateClubRole(c.Title, color.GetDecimalRGB(), c.RoleMentionable)
	if role == nil {
		return errors.New("club role wasn't created")
	}

	return c.EditRoleID(tx, role.ID)
}

func (c *club) CreateChannel(tx *sqlx.Tx) error {
	channel := CreateClubChannel(c.Title, *c.RoleID, []string{c.OwnerID})
	if channel == nil {
		return errors.New("club chat wasn't created")
	}

	return c.EditChannelID(tx, channel.ID)
}
