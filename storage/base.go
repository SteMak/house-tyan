package storage

import "database/sql"

type Base struct {
	ID         string       `db:"id"`
	InsertedAt sql.NullTime `db:"inserted_at"`
	UpdatedAt  sql.NullTime `db:"updated_at"`
}
