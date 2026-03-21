package models

import "time"

type User struct {
	ID         int
	Username   string
	Roles      Roles
	ApiKeyHash string
	Locked     bool
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type UserInput struct {
	Username string
	Roles    Roles
	Password string
}
