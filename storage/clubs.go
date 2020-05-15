package storage

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/jmoiron/sqlx"
)

type Club struct {
	ID          uint       `db:"id"`
	InsertedAt  *time.Time `db:"inserted_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	OwnerID     string     `db:"owner_id"`
	RoleID      *string    `db:"role_id"`
	ChannelID   *string    `db:"channel_id"`
	Title       string     `db:"title"`
	Description *string    `db:"description"`
	Symbol      string     `db:"symbol"`
	IconURL     *string    `db:"icon_url"`
	XP          uint64     `db:"xp"`
	Verified    bool       `db:"verified"`
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

func (c *Club) DeleteMember(tx *sqlx.Tx, memberID string) error {
	return exec(tx, psql.Delete("club_members").
		Where(squirrel.Eq{"user_id": memberID}),
	)
}

func (c *Club) DeleteMembers(tx *sqlx.Tx) error {
	return exec(tx, psql.Delete("club_members").
		Where(squirrel.Eq{"club_id": c.ID}),
	)
}

func (c *Club) HasMember(memberID string) (result bool, err error) {
	err = db.Get(&result, `SELECT EXISTS(SELECT 1 FROM club_members WHERE user_id = $1)`, memberID)
	return
}

func (c *Club) Delete(tx *sqlx.Tx) error {
	return exec(tx, psql.Delete("clubs").
		Where(squirrel.Eq{"id": c.ID}),
	)
}

type ClubMember struct {
	ClubID     uint       `db:"club_id"`
	UserID     string     `db:"user_id"`
	InsertedAt *time.Time `db:"inserted_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
	XP         uint64     `db:"xp"`
}

type clubs struct{}

func (c *clubs) Create(tx *sqlx.Tx, club *Club) error {
	err := psql.Insert("clubs").
		Columns("owner_id", "title", "symbol").
		Values(club.OwnerID, club.Title, club.Symbol).
		Suffix("RETURNING id").
		RunWith(tx).
		QueryRow().Scan(&club.ID)

	if err != nil {
		return err
	}

	return club.AddMember(tx, club.OwnerID)
}

func (c *clubs) DeleteByOwner(tx *sqlx.Tx, ownerID string) error {
	return exec(tx, psql.Delete("clubs").
		Where(squirrel.Eq{"owner_id": ownerID}),
	)
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

func (c *clubs) GetExpired(expiredAfter time.Duration) (clubs []Club, err error) {
	err = db.Select(&clubs, `
		SELECT * FROM clubs
		WHERE NOT verified
			AND localtimestamp >= date_trunc('day', inserted_at) + $1
	`, expiredAfter)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return
}
