package models

import (
	"fmt"
	"io"
	"strconv"

	"github.com/stashapp/stash/pkg/sliceutil"
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

// GetImpactedIDs returns the IDs that will be impacted by the update.
// If the update is to add IDs, then the impacted IDs are the IDs being added.
// If the update is to remove IDs, then the impacted IDs are the IDs being removed.
// If the update is to set IDs, then the impacted IDs are the IDs being removed and the IDs being added.
// Any IDs that are already present and are being added are not returned.
// Likewise, any IDs that are not present that are being removed are not returned.
func (u *UpdateIDs) ImpactedIDs(existing []int) []int {
	if u == nil {
		return nil
	}

	switch u.Mode {
	case RelationshipUpdateModeAdd:
		return intslice.IntExclude(u.IDs, existing)
	case RelationshipUpdateModeRemove:
		return intslice.IntIntercect(existing, u.IDs)
	case RelationshipUpdateModeSet:
		// get the difference between the two lists
		return intslice.IntNotIntersect(existing, u.IDs)
	}

	return nil
}

// GetEffectiveIDs returns the new IDs that will be effective after the update.
func (u *UpdateIDs) EffectiveIDs(existing []int) []int {
	if u == nil {
		return nil
	}

	return effectiveValues(u.IDs, u.Mode, existing)
}

type UpdateStrings struct {
	Values []string               `json:"values"`
	Mode   RelationshipUpdateMode `json:"mode"`
}

func (u *UpdateStrings) Strings() []string {
	if u == nil {
		return nil
	}

	return u.Values
}

// GetEffectiveIDs returns the new IDs that will be effective after the update.
func (u *UpdateStrings) EffectiveValues(existing []string) []string {
	if u == nil {
		return nil
	}

	return effectiveValues(u.Values, u.Mode, existing)
}

// effectiveValues returns the new values that will be effective after the update.
func effectiveValues[T comparable](values []T, mode RelationshipUpdateMode, existing []T) []T {
	switch mode {
	case RelationshipUpdateModeAdd:
		return sliceutil.AppendUniques(existing, values)
	case RelationshipUpdateModeRemove:
		return sliceutil.Exclude(existing, values)
	case RelationshipUpdateModeSet:
		return values
	}

	return nil
}
