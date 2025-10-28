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
		return sliceutil.Exclude(u.IDs, existing)
	case RelationshipUpdateModeRemove:
		return sliceutil.Intersect(existing, u.IDs)
	case RelationshipUpdateModeSet:
		// get the difference between the two lists
		return sliceutil.NotIntersect(existing, u.IDs)
	}

	return nil
}

// Apply applies the update to a list of existing ids, returning the result.
func (u *UpdateIDs) Apply(existing []int) []int {
	if u == nil {
		return existing
	}

	return applyUpdate(u.IDs, u.Mode, existing)
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

// Apply applies the update to a list of existing strings, returning the result.
func (u *UpdateStrings) Apply(existing []string) []string {
	if u == nil {
		return existing
	}

	return applyUpdate(u.Values, u.Mode, existing)
}

// applyUpdate applies values to existing, using the update mode specified.
func applyUpdate[T comparable](values []T, mode RelationshipUpdateMode, existing []T) []T {
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

type UpdateGroupDescriptions struct {
	Groups []GroupIDDescription   `json:"groups"`
	Mode   RelationshipUpdateMode `json:"mode"`
}

// Apply applies the update to a list of existing ids, returning the result.
func (u *UpdateGroupDescriptions) Apply(existing []GroupIDDescription) []GroupIDDescription {
	if u == nil {
		return existing
	}

	switch u.Mode {
	case RelationshipUpdateModeAdd:
		return u.applyAdd(existing)
	case RelationshipUpdateModeRemove:
		return u.applyRemove(existing)
	case RelationshipUpdateModeSet:
		return u.Groups
	}

	return nil
}

func (u *UpdateGroupDescriptions) applyAdd(existing []GroupIDDescription) []GroupIDDescription {
	// overwrite any existing values with the same id
	ret := append([]GroupIDDescription{}, existing...)
	for _, v := range u.Groups {
		found := false
		for i, vv := range ret {
			if vv.GroupID == v.GroupID {
				ret[i] = v
				found = true
				break
			}
		}

		if !found {
			ret = append(ret, v)
		}
	}

	return ret
}

func (u *UpdateGroupDescriptions) applyRemove(existing []GroupIDDescription) []GroupIDDescription {
	// remove any existing values with the same id
	var ret []GroupIDDescription
	for _, v := range existing {
		found := false
		for _, vv := range u.Groups {
			if vv.GroupID == v.GroupID {
				found = true
				break
			}
		}

		// if not found in the remove list, keep it
		if !found {
			ret = append(ret, v)
		}
	}

	return ret
}
