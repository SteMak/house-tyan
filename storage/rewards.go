package storage

type Reward struct {
	AwardID string `db:"award_id"`
	UserID  string `db:"user_id"`
	Amount  uint64 `db:"amount"`
	Paid    bool   `db:"paid"`
}

type rewards struct{}
