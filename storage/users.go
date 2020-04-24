package storage

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
)

type User struct {
	Base

	Name          string `db:"name"`
	Discriminator string `db:"discriminator"`
	XP            uint64 `db:"xp"`
	Balance       uint64 `db:"balance"`
}

type users struct{}

func (users) Set(tx *sqlx.Tx, user *discordgo.User) error {
	_, err := tx.Exec(`
		INSERT INTO users (id, username, discriminator)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET 
			username = $2,
			discriminator = $3,
	`, user.ID, user.Username, user.Discriminator)
	return err
}

func (users) Delete(tx *sql.Tx, id string) error {
	_, err := tx.Exec(`DELETE FROM users WHERE id = $1`, id)
	return err
}

func (users) AddXP(tx *sqlx.Tx, user *discordgo.User, xp int64) error {
	_, err := tx.Exec(`
		INSERT INTO users AS u (id, username, discriminator, xp)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			username = $2,
			discriminator = $3,
			xp = u.xp + $4
	`, user.ID, user.Username, user.Discriminator, xp)
	return err
}
