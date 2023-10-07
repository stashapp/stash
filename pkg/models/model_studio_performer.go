package models

type StudioPerformer struct {
	PerformerID int       `json:"performer_id"`
	StudioID    int       `json:"studio_id"`
	Performer   Performer `json:"performer"`
	Depth       *int
}
