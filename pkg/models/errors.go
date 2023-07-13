package models

import "errors"

var (
	// ErrNotFound signifies entities which are not found
	ErrNotFound = errors.New("not found")

	// ErrConversion signifies conversion errors
	ErrConversion = errors.New("conversion error")

	ErrScraperSource = errors.New("invalid ScraperSource")
)
