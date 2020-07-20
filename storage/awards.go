package storage

import (
	"database/sql"
	"time"

	"errors"

	"github.com/SteMak/house-tyan/cache"
	"github.com/jackc/pgx"
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
	rows, err := pgxconn.Query(`
		SELECT
			award_id,
			user_id,
			amount,
			paid
		FROM rewards
		WHERE award_id = $1
	`,
		award.ID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var rewards []Reward
	for rows.Next() {
		var reward Reward
		err = rows.Scan(
			&reward.AwardID,
			&reward.UserID,
			&reward.Amount,
			&reward.Paid,
		)
		if err != nil {
			return nil, err
		}
		rewards = append(rewards, reward)
	}
	return rewards, nil
}

type awards struct{}

func (awards) Create(tx *pgx.Tx, blank *cache.Blank) (uint64, error) {
	var id uint64
	err := tx.QueryRow(`
		INSERT INTO awards(author_id,reason)
		VALUES($1,$2)
		RETURNING id
	`,
		blank.Author.ID,
		blank.Reason,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	for _, r := range blank.Rewards {
		for _, u := range r.Users {
			_, err = tx.Exec(`
				INSERT INTO rewards(award_id,user_id,amount)
				VALUES ($1,$2,$3)
			`,
				id,
				u.ID,
				r.Amount,
			)

			if err != nil {
				return 0, err
			}
		}
	}
	return id, nil
}

func (awards) SetStatus(tx *pgx.Tx, awardID uint64, status AwardStatus) error {
	_, err := tx.Exec(`
		UPDATE awards SET status=$2 WHERE id = $1
	`,
		awardID,
		status,
	)
	return err
}

func (awards) Accept(tx *pgx.Tx, blankID string) error {
	_, err := tx.Exec(`
		UPDATE awards SET status=$2 WHERE blank_mid = $1
	`,
		blankID,
		AwardStatusAccept,
	)
	return err
}

func (awards) Reject(tx *pgx.Tx, blankID string) error {
	_, err := tx.Exec(`
		UPDATE awards SET status=$2 WHERE blank_mid = $1
	`,
		blankID,
		AwardStatusReject,
	)
	return err
}

func (awards) SetBlankID(tx *pgx.Tx, awardID uint64, blankID string) error {
	_, err := tx.Exec(`
		UPDATE awards SET blank_mid=$2 WHERE id = $1
	`,
		awardID,
		blankID,
	)
	return err
}

func (awards) SetPaid(tx *pgx.Tx, awardID uint64, userID string) error {
	_, err := tx.Exec(`
		UPDATE rewards SET paid=true 
		WHERE award_id = $1 
			AND user_id = $2
	`,
		awardID,
		userID,
	)
	return err
}

func (awards) Get(awardID uint64) (*Award, error) {
	award := new(Award)
	err := pgxconn.QueryRow(`
		SELECT
			id, 
			inserted_at, 
			updated_at, 
			author_id, 
			blank_mid, 
			reason, 
			status
		FROM awards
		WHERE id=$1
	`,
		awardID,
	).Scan(
		&award.ID,
		&award.InsertedAt,
		&award.UpdatedAt,
		&award.AuthorID,
		&award.BlankID,
		&award.Reason,
		&award.Status,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return award, nil
}

func (awards) GetByBlankID(blankID string) (*Award, error) {
	award := new(Award)
	err := pgxconn.QueryRow(`
		SELECT
			id, 
			inserted_at, 
			updated_at, 
			author_id, 
			blank_mid, 
			reason, 
			status
		FROM awards
		WHERE blank_mid=$1
	`,
		blankID,
	).Scan(
		&award.ID,
		&award.InsertedAt,
		&award.UpdatedAt,
		&award.AuthorID,
		&award.BlankID,
		&award.Reason,
		&award.Status,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return award, nil
}
