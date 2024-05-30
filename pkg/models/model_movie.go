package models

import (
	"context"
	"time"
)

type Movie struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Aliases  string `json:"aliases"`
	Duration *int   `json:"duration"`
	Date     *Date  `json:"date"`
	// Rating expressed in 1-100 scale
	Rating    *int      `json:"rating"`
	StudioID  *int      `json:"studio_id"`
	Director  string    `json:"director"`
	Synopsis  string    `json:"synopsis"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	URLs RelatedStrings `json:"urls"`
}

func NewMovie() Movie {
	currentTime := time.Now()
	return Movie{
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
}

func (g *Movie) LoadURLs(ctx context.Context, l URLLoader) error {
	return g.URLs.load(func() ([]string, error) {
		return l.GetURLs(ctx, g.ID)
	})
}

type MoviePartial struct {
	Name     OptionalString
	Aliases  OptionalString
	Duration OptionalInt
	Date     OptionalDate
	// Rating expressed in 1-100 scale
	Rating    OptionalInt
	StudioID  OptionalInt
	Director  OptionalString
	Synopsis  OptionalString
	URLs      *UpdateStrings
	CreatedAt OptionalTime
	UpdatedAt OptionalTime
}

func NewMoviePartial() MoviePartial {
	currentTime := time.Now()
	return MoviePartial{
		UpdatedAt: NewOptionalTime(currentTime),
	}
}
