package models

type PackageSpecInput struct {
	ID        string `json:"id"`
	SourceURL string `json:"sourceURL"`
}

type PackageSource struct {
	Name *string `json:"name"`
	URL  string  `json:"url"`
}
