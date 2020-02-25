package models

import (
	"database/sql"
)

type Movie struct {
	ID             int             `db:"id" json:"id"`
	Front_Image    []byte          `db:"front_image" json:"front_image"`
	Back_Image     []byte          `db:"back_image" json:"back_image"`
	Checksum       string          `db:"checksum" json:"checksum"`
	Name           sql.NullString  `db:"name" json:"name"`
	Aliases        sql.NullString  `db:"aliases" json:"aliases"`
	Duration 	   sql.NullString  `db:"duration" json:"duration"`
	Date     	   SQLiteDate      `db:"date" json:"date"`
	Rating   	   sql.NullString  `db:"rating" json:"rating"`
	Director       sql.NullString  `db:"director" json:"director"`
	Synopsis       sql.NullString  `db:"synopsis" json:"synopsis"`
	URL            sql.NullString  `db:"url" json:"url"`
	CreatedAt      SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt      SQLiteTimestamp `db:"updated_at" json:"updated_at"`
	Scene_index    string          `db:"scene_index" json:"scene_index"`
}

var DefaultMovieImage string = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAYAAABw4pVUAAAABmJLR0QA/wD/AP+gvaeTAAAACXBIWXMAAA3XAAAN1wFCKJt4AAAAB3RJTUUH4wgVBQsJl1CMZAAAASJJREFUeNrt3N0JwyAYhlEj3cj9R3Cm5rbkqtAP+qrnGaCYHPwJpLlaa++mmLpbAERAgAgIEAEBIiBABERAgAgIEAEBIiBABERAgAgIEAHZuVflj40x4i94zhk9vqsVvEq6AsQqMP1EjORx20OACAgQRRx7T+zzcFBxcjNDfoB4ntQqTm5Awo7MlqywZxcgYQ+RlqywJ3ozJAQCSBiEJSsQA0gYBpDAgAARECACAkRAgAgIEAERECACAmSjUv6eAOSB8m8YIGGzBUjYbAESBgMkbBkDEjZbgITBAClcxiqQvEoatreYIWEBASIgJ4Gkf11ntXH3nS9uxfGWfJ5J9hAgAgJEQAQEiIAAERAgAgJEQAQEiIAAERAgAgJEQAQEiL7qBuc6RKLHxr0CAAAAAElFTkSuQmCC"
