package models

import (
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
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	URL       OptionalString
	CreatedAt OptionalTime
	UpdatedAt OptionalTime
}

var DefaultMovieImage = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAYAAABw4pVUAAAABmJLR0QA/wD/AP+gvaeTAAAACXBIWXMAAA3XAAAN1wFCKJt4AAAAB3RJTUUH4wgVBQsJl1CMZAAAASJJREFUeNrt3N0JwyAYhlEj3cj9R3Cm5rbkqtAP+qrnGaCYHPwJpLlaa++mmLpbAERAgAgIEAEBIiBABERAgAgIEAEBIiBABERAgAgIEAHZuVflj40x4i94zhk9vqsVvEq6AsQqMP1EjORx20OACAgQRRx7T+zzcFBxcjNDfoB4ntQqTm5Awo7MlqywZxcgYQ+RlqywJ3ozJAQCSBiEJSsQA0gYBpDAgAARECACAkRAgAgIEAERECACAmSjUv6eAOSB8m8YIGGzBUjYbAESBgMkbBkDEjZbgITBAClcxiqQvEoatreYIWEBASIgJ4Gkf11ntXH3nS9uxfGWfJ5J9hAgAgJEQAQEiIAAERAgAgJEQAQEiIAAERAgAgJEQAQEiL7qBuc6RKLHxr0CAAAAAElFTkSuQmCC"

func NewMovie(name string) *Movie {
	currentTime := time.Now()
	return &Movie{
		Name:      name,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
}

func NewMoviePartial() MoviePartial {
	updatedTime := time.Now()
	return MoviePartial{
		UpdatedAt: NewOptionalTime(updatedTime),
	}
}

type Movies []*Movie

func (m *Movies) Append(o interface{}) {
	*m = append(*m, o.(*Movie))
}

func (m *Movies) New() interface{} {
	return &Movie{}
}
