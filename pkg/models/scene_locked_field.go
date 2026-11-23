package models

import (
	"fmt"
	"io"
	"slices"
	"strconv"
)

// SceneLockableField is a scene field that the user can lock against changes.
type SceneLockableField string

const (
	SceneLockableFieldTitle          SceneLockableField = "TITLE"
	SceneLockableFieldCode           SceneLockableField = "CODE"
	SceneLockableFieldDetails        SceneLockableField = "DETAILS"
	SceneLockableFieldDirector       SceneLockableField = "DIRECTOR"
	SceneLockableFieldDate           SceneLockableField = "DATE"
	SceneLockableFieldProductionDate SceneLockableField = "PRODUCTION_DATE"
	SceneLockableFieldStudio         SceneLockableField = "STUDIO"
	SceneLockableFieldUrls           SceneLockableField = "URLS"
	SceneLockableFieldPerformers     SceneLockableField = "PERFORMERS"
	SceneLockableFieldTags           SceneLockableField = "TAGS"
	SceneLockableFieldGroups         SceneLockableField = "GROUPS"
	SceneLockableFieldGalleries      SceneLockableField = "GALLERIES"
)

var AllSceneLockableField = []SceneLockableField{
	SceneLockableFieldTitle,
	SceneLockableFieldCode,
	SceneLockableFieldDetails,
	SceneLockableFieldDirector,
	SceneLockableFieldDate,
	SceneLockableFieldProductionDate,
	SceneLockableFieldStudio,
	SceneLockableFieldUrls,
	SceneLockableFieldPerformers,
	SceneLockableFieldTags,
	SceneLockableFieldGroups,
	SceneLockableFieldGalleries,
}

func (e SceneLockableField) IsValid() bool {
	switch e {
	case SceneLockableFieldTitle, SceneLockableFieldCode, SceneLockableFieldDetails, SceneLockableFieldDirector,
		SceneLockableFieldDate, SceneLockableFieldProductionDate, SceneLockableFieldStudio, SceneLockableFieldUrls,
		SceneLockableFieldPerformers, SceneLockableFieldTags, SceneLockableFieldGroups, SceneLockableFieldGalleries:
		return true
	}
	return false
}

func (e SceneLockableField) String() string {
	return string(e)
}

func (e *SceneLockableField) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SceneLockableField(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SceneLockableField", str)
	}
	return nil
}

func (e SceneLockableField) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// IsLocked returns true if f is locked on s. s must have its LockedFields loaded.
func (s *Scene) IsLocked(f SceneLockableField) bool {
	return slices.Contains(s.LockedFields.List(), f.String())
}
