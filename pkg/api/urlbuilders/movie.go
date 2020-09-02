package urlbuilders

import "strconv"

type MovieURLBuilder struct {
	BaseURL string
	MovieID string
}

func NewMovieURLBuilder(baseURL string, movieID int) MovieURLBuilder {
	return MovieURLBuilder{
		BaseURL: baseURL,
		MovieID: strconv.Itoa(movieID),
	}
}

func (b MovieURLBuilder) GetMovieFrontImageURL() string {
	return b.BaseURL + "/movie/" + b.MovieID + "/frontimage"
}

func (b MovieURLBuilder) GetMovieBackImageURL() string {
	return b.BaseURL + "/movie/" + b.MovieID + "/backimage"
}
