package models

type StashIDInput struct {
	Endpoint string `json:"endpoint"`
	StashID  string `json:"stash_id"`
}

func StashIDsFromInput(i []*StashIDInput) []StashID {
	var ret []StashID
	for _, stashID := range i {
		newJoin := StashID{
			StashID:  stashID.StashID,
			Endpoint: stashID.Endpoint,
		}
		ret = append(ret, newJoin)
	}

	return ret
}
