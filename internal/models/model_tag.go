package models

type Tag struct {
	ID        int             `db:"id" json:"id"`
	Name      string          `db:"name" json:"name"` // TODO make schema not null
	CreatedAt SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}
