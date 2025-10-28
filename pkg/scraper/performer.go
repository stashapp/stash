package scraper

type ScrapedPerformerInput struct {
	// Set if performer matched
	StoredID       *string  `json:"stored_id"`
	Name           *string  `json:"name"`
	Disambiguation *string  `json:"disambiguation"`
	Gender         *string  `json:"gender"`
	URLs           []string `json:"urls"`
	URL            *string  `json:"url"`       // deprecated
	Twitter        *string  `json:"twitter"`   // deprecated
	Instagram      *string  `json:"instagram"` // deprecated
	Birthdate      *string  `json:"birthdate"`
	Ethnicity      *string  `json:"ethnicity"`
	Country        *string  `json:"country"`
	EyeColor       *string  `json:"eye_color"`
	Height         *string  `json:"height"`
	Measurements   *string  `json:"measurements"`
	FakeTits       *string  `json:"fake_tits"`
	PenisLength    *string  `json:"penis_length"`
	Circumcised    *string  `json:"circumcised"`
	CareerLength   *string  `json:"career_length"`
	Tattoos        *string  `json:"tattoos"`
	Piercings      *string  `json:"piercings"`
	Aliases        *string  `json:"aliases"`
	Details        *string  `json:"details"`
	DeathDate      *string  `json:"death_date"`
	HairColor      *string  `json:"hair_color"`
	Weight         *string  `json:"weight"`
	RemoteSiteID   *string  `json:"remote_site_id"`
}
