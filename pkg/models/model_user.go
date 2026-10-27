package models

import "time"

type User struct {
	ID         int
	Username   string
	Notes      string
	Roles      Roles
	ApiKeyHash string
	Locked     bool
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
