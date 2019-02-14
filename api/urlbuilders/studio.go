package urlbuilders

import "strconv"

type StudioURLBuilder struct {
	BaseURL  string
	StudioID string
}

func NewStudioURLBuilder(baseURL string, studioID int) StudioURLBuilder {
	return StudioURLBuilder{
		BaseURL:  baseURL,
		StudioID: strconv.Itoa(studioID),
	}
}

func (b StudioURLBuilder) GetStudioImageURL() string {
	return b.BaseURL + "/studio/" + b.StudioID + "/image"
}
