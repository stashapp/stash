package models

import "time"

type PerformerImage struct {
	ID         int       `json:"id"`
	PerformerID int      `json:"performer_id"`
	ImageBlob  string    `json:"image_blob"`
	SortOrder  int       `json:"sort_order"`
	CreatedAt  time.Time `json:"created_at"`
}
