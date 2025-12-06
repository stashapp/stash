package models

// FacetCount represents a count for a specific entity in faceted search
type FacetCount struct {
	ID    string
	Label string
	Count int
}

// ResolutionFacetCount represents a count for a specific resolution
type ResolutionFacetCount struct {
	Resolution ResolutionEnum
	Count      int
}

// OrientationFacetCount represents a count for a specific orientation
type OrientationFacetCount struct {
	Orientation OrientationEnum
	Count       int
}

// GenderFacetCount represents a count for a specific gender
type GenderFacetCount struct {
	Gender GenderEnum
	Count  int
}

// BooleanFacetCount represents a count for a boolean value
type BooleanFacetCount struct {
	Value bool
	Count int
}

// RatingFacetCount represents a count for a specific rating value (1-5 stars, stored as 20-100)
type RatingFacetCount struct {
	Rating int // Rating value (100 = 5 stars, 80 = 4 stars, etc.)
	Count  int
}

// CircumcisedFacetCount represents a count for a circumcised value
type CircumcisedFacetCount struct {
	Value CircumisedEnum
	Count int
}

// CaptionFacetCount represents a count for a caption language
type CaptionFacetCount struct {
	Language string
	Count    int
}

// SceneFacetOptions controls which expensive facets to include
type SceneFacetOptions struct {
	// Include performer_tags facet (expensive - 3 joins)
	IncludePerformerTags bool
	// Include captions facet (expensive - file joins)
	IncludeCaptions bool
}

// SceneFacets contains all facet counts for scene filtering
type SceneFacets struct {
	Tags          []FacetCount
	Performers    []FacetCount
	Studios       []FacetCount
	Groups        []FacetCount
	PerformerTags []FacetCount
	Resolutions   []ResolutionFacetCount
	Orientations  []OrientationFacetCount
	Organized     []BooleanFacetCount
	Interactive   []BooleanFacetCount
	Ratings       []RatingFacetCount
	Captions      []CaptionFacetCount
}

// PerformerFacets contains all facet counts for performer filtering
type PerformerFacets struct {
	Tags        []FacetCount
	Studios     []FacetCount
	Genders     []GenderFacetCount
	Countries   []FacetCount
	Circumcised []CircumcisedFacetCount
	Favorite    []BooleanFacetCount
	Ratings     []RatingFacetCount
}

// GalleryFacets contains all facet counts for gallery filtering
type GalleryFacets struct {
	Tags       []FacetCount
	Performers []FacetCount
	Studios    []FacetCount
	Organized  []BooleanFacetCount
	Ratings    []RatingFacetCount
}

// GroupFacets contains all facet counts for group filtering
type GroupFacets struct {
	Tags       []FacetCount
	Performers []FacetCount
	Studios    []FacetCount
}

// StudioFacets contains all facet counts for studio filtering
type StudioFacets struct {
	Tags     []FacetCount
	Parents  []FacetCount
	Favorite []BooleanFacetCount
}

// TagFacets contains all facet counts for tag filtering
type TagFacets struct {
	Parents  []FacetCount
	Children []FacetCount
	Favorite []BooleanFacetCount
}
