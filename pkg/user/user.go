package user

type RoleEnum string

type User struct {
	Username string
	Roles    []RoleEnum
}
