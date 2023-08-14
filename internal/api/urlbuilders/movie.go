package urlbuilders

import (
	"github.com/stashapp/stash/pkg/models"
	"strconv"
)

type MovieURLBuilder struct {
	BaseURL   string
	MovieID   string
	UpdatedAt string
}

func NewMovieURLBuilder(baseURL string, movie *models.Movie) MovieURLBuilder {
	return MovieURLBuilder{
		BaseURL:   baseURL,
		MovieID:   strconv.Itoa(movie.ID),
		UpdatedAt: strconv.FormatInt(movie.UpdatedAt.Unix(), 10),
	}
}

func (b MovieURLBuilder) GetMovieFrontImageURL(hasImage bool) string {
	url := b.BaseURL + "/movie/" + b.MovieID + "/frontimage?t=" + b.UpdatedAt
	if !hasImage {
		url += "&default=true"
	}
	return url
}

func (b MovieURLBuilder) GetMovieBackImageURL() string {
	return b.BaseURL + "/movie/" + b.MovieID + "/backimage?t=" + b.UpdatedAt
}
