package models

import (
	"fmt"
	"io"
	"strconv"
)

type RoleEnum string

const (
	RoleEnumAdmin          RoleEnum = "ADMIN"
	RoleEnumRead           RoleEnum = "READ"
	RoleEnumModify         RoleEnum = "MODIFY"
	RoleEnumGenerateAPIKey RoleEnum = "GENERATE_API_KEY"
)

func (e RoleEnum) Implies(other RoleEnum) bool {
	// admin has all roles
	if e == RoleEnumAdmin {
		return true
	}

	// until we add a NONE value, all values imply read
	if e.IsValid() && other == RoleEnumRead {
		return true
	}

	// all others only imply themselves
	return e == other
}

func (e RoleEnum) IsValid() bool {
	switch e {
	case RoleEnumRead, RoleEnumModify, RoleEnumAdmin, RoleEnumGenerateAPIKey:
		return true
	}
	return false
}

func (e RoleEnum) String() string {
	return string(e)
}

func (e *RoleEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RoleEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RoleEnum", str)
	}
	return nil
}

func (e RoleEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Roles []RoleEnum

func (r Roles) HasRole(role RoleEnum) bool {
	for _, r := range r {
		if r.Implies(role) {
			return true
		}
	}
	return false
}

func (r Roles) Strings() []string {
	strs := make([]string, len(r))
	for i, role := range r {
		strs[i] = role.String()
	}
	return strs
}
