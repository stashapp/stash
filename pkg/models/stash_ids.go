package models

type StashID struct {
	StashID  string `db:"stash_id" json:"stash_id"`
	Endpoint string `db:"endpoint" json:"endpoint"`
}

type UpdateStashIDs struct {
	StashIDs []StashID              `json:"stash_ids"`
	Mode     RelationshipUpdateMode `json:"mode"`
}
