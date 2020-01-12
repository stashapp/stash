package urlbuilders

import "strconv"

type DvdURLBuilder struct {
	BaseURL  string
	DvdID string
}

func NewDvdURLBuilder(baseURL string, dvdID int) DvdURLBuilder {
	return DvdURLBuilder{
		BaseURL:  baseURL,
		DvdID: strconv.Itoa(dvdID),
	}
}

func (b DvdURLBuilder) GetDvdFrontImageURL() string {
	return b.BaseURL + "/dvd/" + b.DvdID + "/frontimage"
}

func (b DvdURLBuilder) GetDvdBackImageURL() string {
	return b.BaseURL + "/dvd/" + b.DvdID + "/backimage"
}

