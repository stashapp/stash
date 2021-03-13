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
		UpdatedAt: strconv.FormatInt(movie.UpdatedAt.Timestamp.Unix(), 10),
	}
}

func (b MovieURLBuilder) GetMovieFrontImageURL() string {
	return b.BaseURL + "/movie/" + b.MovieID + "/frontimage?" + b.UpdatedAt
}

func (b MovieURLBuilder) GetMovieBackImageURL() string {
	return b.BaseURL + "/movie/" + b.MovieID + "/backimage?" + b.UpdatedAt
}
