package scene

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

// LockedFieldsError is returned when an update would change a locked field.
type LockedFieldsError struct {
	Fields []models.SceneLockableField
}

func (e *LockedFieldsError) Error() string {
	names := make([]string, len(e.Fields))
	for i, f := range e.Fields {
		names[i] = f.String()
	}
	return fmt.Sprintf("cannot change locked fields: %s", strings.Join(names, ", "))
}

// LockViolations returns the fields that partial would change while they are
// locked both before and after the update. s must have its LockedFields and
// URLs loaded.
func LockViolations(s *models.Scene, partial models.ScenePartial) []models.SceneLockableField {
	before := s.LockedFields.List()
	after := partial.LockedFields.Apply(before)

	var ret []models.SceneLockableField
	for _, f := range models.AllSceneLockableField {
		if slices.Contains(before, f.String()) && slices.Contains(after, f.String()) && fieldChanged(s, partial, f) {
			ret = append(ret, f)
		}
	}

	return ret
}

// CheckLockedFields returns a LockedFieldsError if partial would change a locked field.
func CheckLockedFields(s *models.Scene, partial models.ScenePartial) error {
	if v := LockViolations(s, partial); len(v) > 0 {
		return &LockedFieldsError{Fields: v}
	}
	return nil
}

// RemoveLockedFields unsets the fields in partial that would change a locked
// field, returning the fields that were removed. Used by writers that should
// skip locked fields rather than fail.
func RemoveLockedFields(s *models.Scene, partial *models.ScenePartial) []models.SceneLockableField {
	v := LockViolations(s, *partial)
	for _, f := range v {
		switch f {
		case models.SceneLockableFieldTitle:
			partial.Title = models.OptionalString{}
		case models.SceneLockableFieldCode:
			partial.Code = models.OptionalString{}
		case models.SceneLockableFieldDetails:
			partial.Details = models.OptionalString{}
		case models.SceneLockableFieldDirector:
			partial.Director = models.OptionalString{}
		case models.SceneLockableFieldDate:
			partial.Date = models.OptionalDate{}
		case models.SceneLockableFieldProductionDate:
			partial.ProductionDate = models.OptionalDate{}
		case models.SceneLockableFieldStudio:
			partial.StudioID = models.OptionalInt{}
		case models.SceneLockableFieldUrls:
			partial.URLs = nil
		}
	}

	return v
}

type LockStateLoader interface {
	models.LockedFieldLoader
	models.URLLoader
}

// LoadLockState loads the relationships needed by LockViolations.
func LoadLockState(ctx context.Context, s *models.Scene, l LockStateLoader) error {
	if err := s.LoadLockedFields(ctx, l); err != nil {
		return fmt.Errorf("loading locked fields: %w", err)
	}
	if err := s.LoadURLs(ctx, l); err != nil {
		return fmt.Errorf("loading urls: %w", err)
	}
	return nil
}

// IsFieldLocked returns true if f is locked on s. s must have its LockedFields loaded.
func IsFieldLocked(s *models.Scene, f models.SceneLockableField) bool {
	return slices.Contains(s.LockedFields.List(), f.String())
}

func fieldChanged(s *models.Scene, partial models.ScenePartial, f models.SceneLockableField) bool {
	switch f {
	case models.SceneLockableFieldTitle:
		return stringChanged(s.Title, partial.Title)
	case models.SceneLockableFieldCode:
		return stringChanged(s.Code, partial.Code)
	case models.SceneLockableFieldDetails:
		return stringChanged(s.Details, partial.Details)
	case models.SceneLockableFieldDirector:
		return stringChanged(s.Director, partial.Director)
	case models.SceneLockableFieldDate:
		return dateChanged(s.Date, partial.Date)
	case models.SceneLockableFieldProductionDate:
		return dateChanged(s.ProductionDate, partial.ProductionDate)
	case models.SceneLockableFieldStudio:
		return partial.StudioID.Set && !intPtrEqual(s.StudioID, partial.StudioID.Ptr())
	case models.SceneLockableFieldUrls:
		return partial.URLs != nil && !slices.Equal(s.URLs.List(), partial.URLs.Apply(s.URLs.List()))
	}

	return false
}

func stringChanged(existing string, v models.OptionalString) bool {
	// a null value clears the field
	return v.Set && v.Value != existing
}

func dateChanged(existing *models.Date, v models.OptionalDate) bool {
	if !v.Set {
		return false
	}

	n := v.Ptr()
	if existing == nil || n == nil {
		return existing != n
	}

	return existing.String() != n.String()
}

func intPtrEqual(a, b *int) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
