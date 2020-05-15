package storage

import (
	"database/sql"
	"time"

	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/SteMak/house-tyan/cache"
	"github.com/jmoiron/sqlx"
)

type AwardStatus uint8

const (
	AwardStatusUnknow = AwardStatus(iota)
	AwardStatusAccept
	AwardStatusReject
)

type Award struct {
	ID         uint64      `db:"id"`
	InsertedAt *time.Time  `db:"inserted_at"`
	UpdatedAt  *time.Time  `db:"updated_at"`
	AuthorID   string      `db:"author_id"`
	BlankID    *string     `db:"blank_mid"`
	Reason     string      `db:"reason"`
	Status     AwardStatus `db:"status"`
}

func (award Award) Reawrds() ([]Reward, error) {
	var rewards []Reward

	query, args, err := psql.Select("*").From("rewards").
		Where(squirrel.Eq{
			"award_id": award.ID,
		}).ToSql()

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	err = db.Select(&rewards, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return rewards, err
}

type awards struct{}

func (awards) Create(tx *sqlx.Tx, blank *cache.Blank) (uint64, error) {
	var id uint64
	err := psql.Insert("awards").
		Columns("author_id", "reason").
		Values(blank.Author.ID, blank.Reason).
		Suffix("RETURNING id").
		RunWith(tx).
		QueryRow().Scan(&id)

	if err != nil {
		return 0, err
	}

	for _, r := range blank.Rewards {
		for _, u := range r.Users {
			_, err = psql.Insert("rewards").
				Columns("award_id", "user_id", "amount").
				Values(id, u.ID, r.Amount).
				RunWith(tx).Exec()

			if err != nil {
				return 0, err
			}
		}
	}
	return id, nil
}

func (awards) SetStatus(tx *sqlx.Tx, awardID uint64, status AwardStatus) error {
	return exec(tx, psql.Update("awards").
		Where(squirrel.Eq{"id": awardID}).
		Set("status", status),
	)
}

func (awards) Accept(tx *sqlx.Tx, blankID string) error {
	return exec(tx, psql.Update("awards").
		Where(squirrel.Eq{"blank_mid": blankID}).
		Set("status", AwardStatusAccept),
	)
}

func (awards) Reject(tx *sqlx.Tx, blankID string) error {
	return exec(tx, psql.Update("awards").
		Where(squirrel.Eq{"blank_mid": blankID}).
		Set("status", AwardStatusReject),
	)
}

func (awards) SetBlankID(tx *sqlx.Tx, awardID uint64, blankID string) error {
	return exec(tx, psql.Update("awards").
		Where(squirrel.Eq{"id": awardID}).
		Set("blank_mid", blankID),
	)
}

func (awards) SetPaid(tx *sqlx.Tx, awardID uint64, userID string) error {
	return exec(tx,
		psql.Update("rewards").
			Where(squirrel.Eq{
				"award_id": awardID,
				"user_id":  userID,
			}).
			Set("paid", true),
	)
}

func (awards) Get(awardID uint64) (*Award, error) {
	query := `SELECT * FROM awards WHERE id = $1 LIMIT 1`
	award := new(Award)
	err := db.Get(award, query, awardID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return award, nil
}

func (awards) GetByBlankID(blankID string) (*Award, error) {
	query := `SELECT * FROM awards WHERE blank_mid = $1 LIMIT 1`
	award := new(Award)
	err := db.Get(award, query, blankID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return award, nil
}
