package models

import (
	"slices"
	"time"
)

type StashID struct {
	StashID   string    `db:"stash_id" json:"stash_id"`
	Endpoint  string    `db:"endpoint" json:"endpoint"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (s StashID) ToStashIDInput() StashIDInput {
	t := s.UpdatedAt
	return StashIDInput{
		StashID:   s.StashID,
		Endpoint:  s.Endpoint,
		UpdatedAt: &t,
	}
}

type StashIDs []StashID

func (s StashIDs) ToStashIDInputs() StashIDInputs {
	if s == nil {
		return nil
	}

	ret := make(StashIDInputs, len(s))
	for i, v := range s {
		ret[i] = v.ToStashIDInput()
	}
	return ret
}

// HasSameStashIDs returns true if the two lists of StashIDs are the same, ignoring order and updated at time.
func (s StashIDs) HasSameStashIDs(other StashIDs) bool {
	if len(s) != len(other) {
		return false
	}

	for _, v := range s {
		if !slices.ContainsFunc(other, func(o StashID) bool {
			return o.StashID == v.StashID && o.Endpoint == v.Endpoint
		}) {
			return false
		}
	}

	return true
}

type StashIDInput struct {
	StashID   string     `db:"stash_id" json:"stash_id"`
	Endpoint  string     `db:"endpoint" json:"endpoint"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}

func (s StashIDInput) ToStashID() StashID {
	ret := StashID{
		StashID:  s.StashID,
		Endpoint: s.Endpoint,
	}
	if s.UpdatedAt != nil {
		ret.UpdatedAt = *s.UpdatedAt
	} else {
		// default to now if not provided
		ret.UpdatedAt = time.Now()
	}

	return ret
}

type StashIDInputs []StashIDInput

func (s StashIDInputs) ToStashIDs() StashIDs {
	if s == nil {
		return nil
	}

	// #2800 - deduplicate StashIDs based on endpoint and stash_id
	ret := make(StashIDs, 0, len(s))
	seen := make(map[string]map[string]bool)

	for _, v := range s {
		stashID := v.ToStashID()

		if seen[stashID.Endpoint] == nil {
			seen[stashID.Endpoint] = make(map[string]bool)
		}

		if !seen[stashID.Endpoint][stashID.StashID] {
			seen[stashID.Endpoint][stashID.StashID] = true
			ret = append(ret, stashID)
		}
	}

	return ret
}

type UpdateStashIDs struct {
	StashIDs []StashID              `json:"stash_ids"`
	Mode     RelationshipUpdateMode `json:"mode"`
}

// AddUnique adds the stash id to the list, only if the endpoint/stashid pair does not already exist in the list.
func (u *UpdateStashIDs) AddUnique(v StashID) {
	for _, vv := range u.StashIDs {
		if vv.StashID == v.StashID && vv.Endpoint == v.Endpoint {
			return
		}
	}

	u.StashIDs = append(u.StashIDs, v)
}

// Set sets or replaces the stash id for the endpoint in the provided value.
func (u *UpdateStashIDs) Set(v StashID) {
	for i, vv := range u.StashIDs {
		if vv.Endpoint == v.Endpoint {
			u.StashIDs[i] = v
			return
		}
	}

	u.StashIDs = append(u.StashIDs, v)
}

type StashIDCriterionInput struct {
	// If present, this value is treated as a predicate.
	// That is, it will filter based on stash_id with the matching endpoint
	Endpoint *string           `json:"endpoint"`
	StashID  *string           `json:"stash_id"`
	Modifier CriterionModifier `json:"modifier"`
}

type StashIDsCriterionInput struct {
	// If present, this value is treated as a predicate.
	// That is, it will filter based on stash_ids with the matching endpoint
	Endpoint *string           `json:"endpoint"`
	StashIDs []*string         `json:"stash_ids"`
	Modifier CriterionModifier `json:"modifier"`
}
