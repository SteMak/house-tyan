package storage

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/jmoiron/sqlx"
)

type Club struct {
	ID              uint       `db:"id"`
	InsertedAt      *time.Time `db:"inserted_at"`
	UpdatedAt       *time.Time `db:"updated_at"`
	OwnerID         string     `db:"owner_id"`
	RoleID          *string    `db:"role_id"`
	RoleColor       string     `db:"role_color"`
	RoleMentionable bool       `db:"role_mentionable"`
	ChannelID       *string    `db:"channel_id"`
	CardMessageID   *string    `db:"card_message_id"`
	Title           string     `db:"title"`
	Description     *string    `db:"description"`
	Symbol          string     `db:"symbol"`
	IconURL         *string    `db:"icon_url"`
	XP              uint64     `db:"xp"`
	ExpiredAt       *time.Time `db:"expired_at"`
	Verified        bool       `db:"verified"`
}

func (c *Club) randomize() {
	desc := gofakeit.Paragraph(1, 1, 10, "")
	chid := gofakeit.Numerify("test##############")
	rlid := gofakeit.Numerify("test##############")

	c.OwnerID = gofakeit.Numerify("test##############")
	c.ChannelID = &chid
	c.RoleID = &rlid
	c.Title = gofakeit.Word()
	c.Symbol = gofakeit.Emoji()
	c.Description = &desc
}

func (c *Club) AddMember(tx *sqlx.Tx, memberID string) error {
	return exec(tx, psql.Insert("club_members").
		Values(c.ID, memberID).
		Suffix("ON CONFLICT DO NOTHING"),
	)
}

func (c *Club) MakeMemberManager(tx *sqlx.Tx, memberID string) error {
	return exec(tx, psql.Update("club_members").
		Where(squirrel.Eq{"club_members.user_id": memberID}).
		Set("manager", true),
	)
}

func (c *Club) MakeMemberUser(tx *sqlx.Tx, memberID string) error {
	return exec(tx, psql.Update("club_members").
		Where(squirrel.Eq{"club_members.user_id": memberID}).
		Set("manager", false),
	)
}

func (c *Club) DeleteMember(tx *sqlx.Tx, memberID string) error {
	return exec(tx, psql.Delete("club_members").
		Where(squirrel.Eq{"user_id": memberID}),
	)
}

func (c *Club) HasMember(memberID string) (result bool, err error) {
	err = db.Get(&result, `
		SELECT EXISTS(SELECT 1 FROM club_members WHERE club_id = $1 AND user_id = $2)
	`, c.ID, memberID)
	return
}

func (c *Club) HasInvite(memberID string) (result bool, err error) {
	err = db.Get(&result, `
		SELECT EXISTS(SELECT 1 FROM club_invites WHERE club_id = $1 AND user_id = $2)
	`, c.ID, memberID)
	return
}

func (c *clubs) GetInvite(messageID string) (*ClubInvite, error) {
	invite := new(ClubInvite)
	err := db.Get(invite, `SELECT * FROM club_invites WHERE message_id = $1`, messageID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return invite, err
}

func (c *clubs) DeleteUserInvites(userID string) error {
	_, err := db.Exec(`DELETE FROM club_invites WHERE user_id = $1`, userID)
	return err
}

func (c *Club) GetMembers() (*[]ClubMember, error) {
	query, args, err := psql.Select("cm.*").
		From("club_members cm").
		Where(squirrel.Eq{"cm.club_id": c.ID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	members := new([]ClubMember)
	err = db.Select(members, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return members, err
}

func (c *Club) Delete(tx *sqlx.Tx) error {
	if err := exec(tx, psql.Delete("club_members").Where(squirrel.Eq{"club_id": c.ID})); err != nil {
		return err
	}
	if err := exec(tx, psql.Delete("club_invites").Where(squirrel.Eq{"club_id": c.ID})); err != nil {
		return err
	}

	return exec(tx, psql.Delete("clubs").
		Where(squirrel.Eq{"id": c.ID}),
	)
}

func (c *Club) EditRoleID(tx *sqlx.Tx, roleID string) error {
	c.RoleID = &roleID
	return exec(tx, psql.Update("clubs").
		Where(squirrel.Eq{"id": c.ID}).
		Set("role_id", roleID),
	)
}

func (c *Club) EditChannelID(tx *sqlx.Tx, channelID string) error {
	c.ChannelID = &channelID
	return exec(tx, psql.Update("clubs").
		Where(squirrel.Eq{"id": c.ID}).
		Set("channel_id", channelID),
	)
}

func (c *Club) EditDescription(tx *sqlx.Tx, desc string) error {
	c.Description = &desc
	return exec(tx, psql.Update("clubs").
		Where(squirrel.Eq{"id": c.ID}).
		Set("description", desc),
	)
}

func (c *Club) ClubApply(tx *sqlx.Tx, userID, messageID string) error {
	return exec(tx, psql.Insert("club_invites").
		Columns("club_id", "user_id", "is_request", "message_id").
		Values(c.ID, userID, true, messageID).
		Suffix("ON CONFLICT DO NOTHING"),
	)
}

type ClubMember struct {
	ClubID     uint       `db:"club_id"`
	UserID     string     `db:"user_id"`
	Manager    bool       `db:"manager"`
	InsertedAt *time.Time `db:"inserted_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
	XP         uint64     `db:"xp"`
}

type ClubInvite struct {
	ClubID     uint       `db:"club_id"`
	UserID     string     `db:"user_id"`
	IsRequest  bool       `db:"is_request"`
	MessageID  string     `db:"message_id"`
	InsertedAt *time.Time `db:"inserted_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}

type clubs struct{}

func (c *clubs) Create(tx *sqlx.Tx, club *Club) error {
	err := psql.Insert("clubs").
		Columns("owner_id", "title", "symbol", "expired_at").
		Values(club.OwnerID, club.Title, club.Symbol, club.ExpiredAt).
		Suffix("RETURNING id").
		RunWith(tx).
		QueryRow().Scan(&club.ID)
	if err != nil {
		return err
	}

	err = club.AddMember(tx, club.OwnerID)
	if err != nil {
		return err
	}

	return club.MakeMemberManager(tx, club.OwnerID)
}

func (c *clubs) DeleteByOwner(tx *sqlx.Tx, ownerID string) error {
	return exec(tx, psql.Delete("clubs").
		Where(squirrel.Eq{"owner_id": ownerID}),
	)
}

func (c *clubs) GetClubByID(id uint) (*Club, error) {
	query, args, err := psql.Select("c.*").
		From("clubs c").
		Where(squirrel.Eq{"c.id": id}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, err
	}

	club := new(Club)
	err = db.Get(club, query, args...)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return club, err
}

func (c *clubs) GetClubByUser(userID string) (*Club, error) {
	query, args, err := psql.Select("c.*").
		From("club_members cm").
		Join("clubs c ON c.id=cm.club_id").
		Where(squirrel.Eq{"cm.user_id": userID}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, err
	}

	club := new(Club)
	err = db.Get(club, query, args...)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return club, err
}

func (c *clubs) GetClubByTitle(title string) (*Club, error) {
	query, args, err := psql.Select("c.*").
		From("clubs c").
		Where(squirrel.Eq{"c.title": title}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, err
	}

	club := new(Club)
	err = db.Get(club, query, args...)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return club, err
}

func (c *clubs) GetClubByPrefix(title string) (*Club, error) {
	query, args, err := psql.Select("c.*").
		From("clubs c").
		Where(squirrel.Eq{"c.symbol": title}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, err
	}

	club := new(Club)
	err = db.Get(club, query, args...)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return club, err
}

func (c *clubs) GetClub(arg string) (*Club, error) {
	var (
		club *Club
		err  error
	)

	userID := arg
	userID = strings.TrimPrefix(userID, "<@")
	userID = strings.TrimPrefix(userID, "!")
	userID = strings.TrimSuffix(userID, ">")

	if club, err = Clubs.GetClubByTitle(arg); err != nil || club != nil {
		if err != nil {
			return nil, err
		}
	} else if club, err = Clubs.GetClubByPrefix(arg); err != nil || club != nil {
		if err != nil {
			return nil, err
		}
	} else if club, err = Clubs.GetClubByUser(userID); err != nil || club != nil {
		if err != nil {
			return nil, err
		}
	}

	return club, nil
}

func (c *clubs) GetExpired() (clubs []Club, err error) {
	err = db.Select(&clubs, `
		SELECT * FROM clubs
		WHERE NOT verified
			AND localtimestamp >= expired_at
	`)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return
}
