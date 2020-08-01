package urlbuilders

import "strconv"

type TagURLBuilder struct {
	BaseURL string
	TagID   string
}

func NewTagURLBuilder(baseURL string, tagID int) TagURLBuilder {
	return TagURLBuilder{
		BaseURL: baseURL,
		TagID:   strconv.Itoa(tagID),
	}
}

func (b TagURLBuilder) GetTagImageURL() string {
	return b.BaseURL + "/tag/" + b.TagID + "/image"
}
