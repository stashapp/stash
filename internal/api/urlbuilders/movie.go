package urlbuilders

import (
	"github.com/stashapp/stash/pkg/models"
	"strconv"
)

type MovieURLBuilder struct {
	BaseURL string
	MovieID string
}

func NewMovieURLBuilder(baseURL string, movie *models.Movie) MovieURLBuilder {
	return MovieURLBuilder{
		BaseURL: baseURL,
		MovieID: strconv.Itoa(movie.ID),
	}
}

func (b MovieURLBuilder) GetMovieFrontImageURL() string {
	return b.BaseURL + "/movie/" + b.MovieID + "/frontimage"
}

func (b MovieURLBuilder) GetMovieBackImageURL() string {
	return b.BaseURL + "/movie/" + b.MovieID + "/backimage"
}
