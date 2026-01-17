package models

type User struct {
	Username string
	Roles    Roles
}

type UserInput struct {
	Username string
	Roles    Roles
	Password string
}
