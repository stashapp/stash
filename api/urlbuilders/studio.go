package urlbuilders

import "strconv"

type studioURLBuilder struct {
	BaseURL string
	StudioID string
}

func NewStudioURLBuilder(baseURL string, studioID int) studioURLBuilder {
	return studioURLBuilder{
		BaseURL: baseURL,
		StudioID: strconv.Itoa(studioID),
	}
}

func (b studioURLBuilder) GetStudioImageUrl() string {
	return b.BaseURL + "/studio/" + b.StudioID + "/image"
}
