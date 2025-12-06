package api

import "github.com/stashapp/stash/pkg/models"

// FacetCount represents a count for a specific entity in faceted search
type FacetCount struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Count int    `json:"count"`
}

// ResolutionFacetCount represents a count for a specific resolution
type ResolutionFacetCount struct {
	Resolution models.ResolutionEnum `json:"resolution"`
	Count      int                   `json:"count"`
}

// OrientationFacetCount represents a count for a specific orientation
type OrientationFacetCount struct {
	Orientation models.OrientationEnum `json:"orientation"`
	Count       int                    `json:"count"`
}

// GenderFacetCount represents a count for a specific gender
type GenderFacetCount struct {
	Gender models.GenderEnum `json:"gender"`
	Count  int               `json:"count"`
}

// BooleanFacetCount represents a count for a boolean value
type BooleanFacetCount struct {
	Value bool `json:"value"`
	Count int  `json:"count"`
}

// RatingFacetCount represents a count for a rating value
type RatingFacetCount struct {
	Rating int `json:"rating"`
	Count  int `json:"count"`
}

// CircumcisedFacetCount represents a count for a circumcised value
type CircumcisedFacetCount struct {
	Value models.CircumisedEnum `json:"value"`
	Count int                   `json:"count"`
}

// CaptionFacetCount represents a count for a caption language
type CaptionFacetCount struct {
	Language string `json:"language"`
	Count    int    `json:"count"`
}

// SceneFacetsResult contains all facet counts for scene filtering
type SceneFacetsResult struct {
	Tags          []*FacetCount            `json:"tags"`
	Performers    []*FacetCount            `json:"performers"`
	Studios       []*FacetCount            `json:"studios"`
	Groups        []*FacetCount            `json:"groups"`
	PerformerTags []*FacetCount            `json:"performer_tags"`
	Resolutions   []*ResolutionFacetCount  `json:"resolutions"`
	Orientations  []*OrientationFacetCount `json:"orientations"`
	Organized     []*BooleanFacetCount     `json:"organized"`
	Interactive   []*BooleanFacetCount     `json:"interactive"`
	Ratings       []*RatingFacetCount      `json:"ratings"`
	Captions      []*CaptionFacetCount     `json:"captions"`
}

// PerformerFacetsResult contains all facet counts for performer filtering
type PerformerFacetsResult struct {
	Tags        []*FacetCount            `json:"tags"`
	Studios     []*FacetCount            `json:"studios"`
	Genders     []*GenderFacetCount      `json:"genders"`
	Countries   []*FacetCount            `json:"countries"`
	Circumcised []*CircumcisedFacetCount `json:"circumcised"`
	Favorite    []*BooleanFacetCount     `json:"favorite"`
	Ratings     []*RatingFacetCount      `json:"ratings"`
}

// GalleryFacetsResult contains all facet counts for gallery filtering
type GalleryFacetsResult struct {
	Tags       []*FacetCount        `json:"tags"`
	Performers []*FacetCount        `json:"performers"`
	Studios    []*FacetCount        `json:"studios"`
	Organized  []*BooleanFacetCount `json:"organized"`
	Ratings    []*RatingFacetCount  `json:"ratings"`
}

// GroupFacetsResult contains all facet counts for group filtering
type GroupFacetsResult struct {
	Tags       []*FacetCount `json:"tags"`
	Performers []*FacetCount `json:"performers"`
	Studios    []*FacetCount `json:"studios"`
}

// StudioFacetsResult contains all facet counts for studio filtering
type StudioFacetsResult struct {
	Tags     []*FacetCount        `json:"tags"`
	Parents  []*FacetCount        `json:"parents"`
	Favorite []*BooleanFacetCount `json:"favorite"`
}

// TagFacetsResult contains all facet counts for tag filtering
type TagFacetsResult struct {
	Parents  []*FacetCount        `json:"parents"`
	Children []*FacetCount        `json:"children"`
	Favorite []*BooleanFacetCount `json:"favorite"`
}
