package models

import (
	"fmt"
	"io"
	"strconv"

	"github.com/stashapp/stash/pkg/sliceutil/intslice"
)

type RelationshipUpdateMode string

const (
	RelationshipUpdateModeSet    RelationshipUpdateMode = "SET"
	RelationshipUpdateModeAdd    RelationshipUpdateMode = "ADD"
	RelationshipUpdateModeRemove RelationshipUpdateMode = "REMOVE"
)

var AllRelationshipUpdateMode = []RelationshipUpdateMode{
	RelationshipUpdateModeSet,
	RelationshipUpdateModeAdd,
	RelationshipUpdateModeRemove,
}

func (e RelationshipUpdateMode) IsValid() bool {
	switch e {
	case RelationshipUpdateModeSet, RelationshipUpdateModeAdd, RelationshipUpdateModeRemove:
		return true
	}
	return false
}

func (e RelationshipUpdateMode) String() string {
	return string(e)
}

func (e *RelationshipUpdateMode) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RelationshipUpdateMode(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RelationshipUpdateMode", str)
	}
	return nil
}

func (e RelationshipUpdateMode) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type UpdateIDs struct {
	IDs  []int                  `json:"ids"`
	Mode RelationshipUpdateMode `json:"mode"`
}

func (u *UpdateIDs) IDStrings() []string {
	if u == nil {
		return nil
	}

	return intslice.IntSliceToStringSlice(u.IDs)
}
