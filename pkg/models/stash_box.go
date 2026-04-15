package models

type StashBoxFingerprint struct {
	Algorithm string `json:"algorithm"`
	Hash      string `json:"hash"`
	Duration  int    `json:"duration"`
}

type StashBox struct {
	Endpoint             string `json:"endpoint"`
	APIKey               string `json:"api_key"`
	Name                 string `json:"name"`
	MaxRequestsPerMinute int    `json:"max_requests_per_minute" koanf:"max_requests_per_minute"`
}
