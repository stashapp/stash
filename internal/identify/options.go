package identify

import (
	"fmt"
	"io"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

type IdentifyMetadataTaskOptions struct {
	// An ordered list of sources to identify items with. Only the first source that finds a match is used.
	Sources []*IdentifySource `json:"sources"`
	// Options defined here override the configured defaults
	Options *IdentifyMetadataOptions `json:"options"`
}

type IdentifySource struct {
	Source *models.ScraperSource `json:"source"`
	// Options defined for a source override the defaults
	Options *IdentifyMetadataOptions `json:"options"`
}

type IdentifySourceInput struct {
	Source *models.ScraperSourceInput `json:"source"`
	// Options defined for a source override the defaults
	Options *IdentifyMetadataOptionsInput `json:"options"`
}

type IdentifyMetadataInput struct {
	// An ordered list of sources to identify items with. Only the first source that finds a match is used.
	Sources []*IdentifySourceInput `json:"sources"`
	// Options defined here override the configured defaults
	Options *IdentifyMetadataOptionsInput `json:"options"`
	// scene ids to identify
	SceneIDs []string `json:"sceneIDs"`
	// paths of scenes to identify - ignored if scene ids are set
	Paths []string `json:"paths"`
}

type IdentifyMetadataOptions struct {
	// any fields missing from here are defaulted to MERGE and createMissing false
	FieldOptions []*IdentifyFieldOptions `json:"fieldOptions"`
	// defaults to true if not provided
	SetCoverImage *bool `json:"setCoverImage"`
	SetOrganized  *bool `json:"setOrganized"`
	// defaults to true if not provided
	IncludeMalePerformers *bool `json:"includeMalePerformers"`
}

type IdentifyMetadataOptionsInput struct {
	// any fields missing from here are defaulted to MERGE and createMissing false
	FieldOptions []*IdentifyFieldOptionsInput `json:"fieldOptions"`
	// defaults to true if not provided
	SetCoverImage *bool `json:"setCoverImage"`
	SetOrganized  *bool `json:"setOrganized"`
	// defaults to true if not provided
	IncludeMalePerformers *bool `json:"includeMalePerformers"`
}

type IdentifyFieldOptions struct {
	Field    string                `json:"field"`
	Strategy IdentifyFieldStrategy `json:"strategy"`
	// creates missing objects if needed - only applicable for performers, tags and studios
	CreateMissing *bool `json:"createMissing"`
}

type IdentifyFieldOptionsInput struct {
	Field    string                `json:"field"`
	Strategy IdentifyFieldStrategy `json:"strategy"`
	// creates missing objects if needed - only applicable for performers, tags and studios
	CreateMissing *bool `json:"createMissing"`
}

type IdentifyFieldStrategy string

const (
	// Never sets the field value
	IdentifyFieldStrategyIgnore IdentifyFieldStrategy = "IGNORE"
	// For multi-value fields, merge with existing.
	// For single-value fields, ignore if already set
	IdentifyFieldStrategyMerge IdentifyFieldStrategy = "MERGE"
	// Always replaces the value if a value is found.
	//   For multi-value fields, any existing values are removed and replaced with the
	//   scraped values.
	IdentifyFieldStrategyOverwrite IdentifyFieldStrategy = "OVERWRITE"
)

var AllIdentifyFieldStrategy = []IdentifyFieldStrategy{
	IdentifyFieldStrategyIgnore,
	IdentifyFieldStrategyMerge,
	IdentifyFieldStrategyOverwrite,
}

func (e IdentifyFieldStrategy) IsValid() bool {
	switch e {
	case IdentifyFieldStrategyIgnore, IdentifyFieldStrategyMerge, IdentifyFieldStrategyOverwrite:
		return true
	}
	return false
}

func (e IdentifyFieldStrategy) String() string {
	return string(e)
}

func (e *IdentifyFieldStrategy) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = IdentifyFieldStrategy(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid IdentifyFieldStrategy", str)
	}
	return nil
}

func (e IdentifyFieldStrategy) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
