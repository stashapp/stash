package models

type StashID struct {
	StashID  string `db:"stash_id" json:"stash_id"`
	Endpoint string `db:"endpoint" json:"endpoint"`
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
