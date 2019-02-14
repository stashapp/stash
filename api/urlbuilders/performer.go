package urlbuilders

import "strconv"

type PerformerURLBuilder struct {
	BaseURL     string
	PerformerID string
}

func NewPerformerURLBuilder(baseURL string, performerID int) PerformerURLBuilder {
	return PerformerURLBuilder{
		BaseURL:     baseURL,
		PerformerID: strconv.Itoa(performerID),
	}
}

func (b PerformerURLBuilder) GetPerformerImageURL() string {
	return b.BaseURL + "/performer/" + b.PerformerID + "/image"
}
