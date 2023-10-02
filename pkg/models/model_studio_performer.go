package models

type StudioPerformer struct {
	ID        int       `json:"id"`
	StudioID  int       `json:"studio_id"`
	Performer Performer `json:"performer"`
	Depth     *int
}
