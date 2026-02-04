package models

type User struct {
	Username string
	Roles    Roles
	ApiKey   string
}

type UserInput struct {
	Username string
	Roles    Roles
	Password string
}
