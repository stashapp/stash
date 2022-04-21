package graphql

import "github.com/99designs/gqlgen/graphql"

// Override for generated struct due to mistaken omitempty
// https://github.com/Yamashou/gqlgenc/issues/77
type SceneDraftInput struct {
	Title        *string             `json:"title,omitempty"`
	Details      *string             `json:"details,omitempty"`
	URL          *string             `json:"url,omitempty"`
	Date         *string             `json:"date,omitempty"`
	Studio       *DraftEntityInput   `json:"studio,omitempty"`
	Performers   []*DraftEntityInput `json:"performers"`
	Tags         []*DraftEntityInput `json:"tags,omitempty"`
	Image        *graphql.Upload     `json:"image,omitempty"`
	Fingerprints []*FingerprintInput `json:"fingerprints"`
}
