package urlbuilders

import "strconv"

type performerURLBuilder struct {
	BaseURL string
	PerformerID string
}

func NewPerformerURLBuilder(baseURL string, performerID int) performerURLBuilder {
	return performerURLBuilder{
		BaseURL: baseURL,
		PerformerID: strconv.Itoa(performerID),
	}
}

func (b performerURLBuilder) GetPerformerImageUrl() string {
	return b.BaseURL + "/performer/" + b.PerformerID + "/image"
}
