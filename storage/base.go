package storage

import (
	"time"
)

type Base struct {
	ID         string     `db:"id"`
	InsertedAt *time.Time `db:"inserted_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}
