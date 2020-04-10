package storage

import (
	"database/sql"

	"github.com/SteMak/house-tyan/cache"
)

type AwardStatus uint8

const (
	AwardStatusUnknow = iota
	AwardStatusAccept
	AwardStatusDiscard
)

type Award struct {
	Base
	AuthorID string `db:"author_id"`
	Reason   string `db:"reason"`
	Status   uint8  `db:"status"`
}

type Reward struct {
	AwardID string `db:"award_id"`
	UserID  string `db:"user_id"`
	Amount  uint64 `db:"amount"`
	Paid    bool   `db:"paid"`
}

type awards struct{}

func (awards) Create(tx *sql.Tx, id string, blank *cache.Blank) error {
	_, err := tx.Exec(`
		INSERT INTO awards (id, author_id, reason)
		VALUES ($1, $2, $3)
	`, id, blank.Author.ID, blank.Reason)

	if err != nil {
		return err
	}

	for _, reward := range blank.Rewards {
		for _, user := range reward.Users {
			_, err := tx.Exec(`
					INSERT INTO rewards (award_id, user_id, amount)
					VALUES ($1, $2, $3)
				`,
				id,
				user.ID,
				reward.Amount,
			)

			if err != nil {
				return err
			}
		}
	}

	return nil
}
