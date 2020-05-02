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
	ChannelID   string     `db:"channel_id"`
	RoleID      string     `db:"role_id"`
	Title       string     `db:"title"`
	Description *string    `db:"description"`
	Symbol      string     `db:"symbol"`
	IconURL     *string    `db:"icon_url"`
	XP          uint64     `db:"xp"`
}

func (c *Club) randomize() {
	desc := gofakeit.Paragraph(1, 1, 10, "")

	c.OwnerID = gofakeit.Numerify("test##############")
	c.ChannelID = gofakeit.Numerify("test##############")
	c.RoleID = gofakeit.Numerify("test##############")
	c.Title = gofakeit.Word()
	c.Symbol = gofakeit.Emoji()
	c.Description = &desc
}

func (c *Club) AddMember(tx *sqlx.Tx, userID string) error {
	return exec(tx, psql.Insert("club_members").
		Values(c.ID, userID).
		Suffix("ON CONFLICT DO NOTHING"),
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
		Columns("owner_id", "role_id", "title", "description", "symbol").
		Values(club.OwnerID, club.RoleID, club.Title, club.Description, club.Symbol).
		Suffix("RETURNING id").
		RunWith(tx).
		QueryRow().Scan(&club.ID)

	if err != nil {
		return err
	}

	return club.AddMember(tx, club.OwnerID)
}

func (clubs) UserClub(userID string) (*Club, error) {
	query, args, err := psql.Select("*").
		From("clubs").
		Where(squirrel.Eq{"owner_id": userID}).
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

	if err != nil {
		return nil, err
	}
	return club, nil
}
