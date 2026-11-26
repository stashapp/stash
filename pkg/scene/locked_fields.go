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
		case models.SceneLockableFieldPerformers:
			partial.PerformerIDs = nil
		case models.SceneLockableFieldTags:
			partial.TagIDs = nil
		case models.SceneLockableFieldGroups:
			partial.GroupIDs = nil
		case models.SceneLockableFieldGalleries:
			partial.GalleryIDs = nil
		}
	}

	return v
}

type LockStateLoader interface {
	models.LockedFieldLoader
	models.URLLoader
	models.PerformerIDLoader
	models.TagIDLoader
	models.SceneGroupLoader
	models.GalleryIDLoader
}

// LoadLockState loads the relationships needed by LockViolations.
// Relationship fields are only loaded when locked.
func LoadLockState(ctx context.Context, s *models.Scene, l LockStateLoader) error {
	if err := s.LoadLockedFields(ctx, l); err != nil {
		return fmt.Errorf("loading locked fields: %w", err)
	}
	if err := s.LoadURLs(ctx, l); err != nil {
		return fmt.Errorf("loading urls: %w", err)
	}
	if s.IsLocked(models.SceneLockableFieldPerformers) {
		if err := s.LoadPerformerIDs(ctx, l); err != nil {
			return fmt.Errorf("loading performers: %w", err)
		}
	}
	if s.IsLocked(models.SceneLockableFieldTags) {
		if err := s.LoadTagIDs(ctx, l); err != nil {
			return fmt.Errorf("loading tags: %w", err)
		}
	}
	if s.IsLocked(models.SceneLockableFieldGroups) {
		if err := s.LoadGroups(ctx, l); err != nil {
			return fmt.Errorf("loading groups: %w", err)
		}
	}
	if s.IsLocked(models.SceneLockableFieldGalleries) {
		if err := s.LoadGalleryIDs(ctx, l); err != nil {
			return fmt.Errorf("loading galleries: %w", err)
		}
	}
	return nil
}

type LockedFieldFinder interface {
	models.SceneGetter
	models.LockedFieldLoader
}

// CheckGalleriesLocked returns a LockedFieldsError naming the first scene in
// sceneIDs that has its galleries locked. Used when a gallery-side edit would
// change the scenes' gallery links.
func CheckGalleriesLocked(ctx context.Context, r LockedFieldFinder, sceneIDs []int) error {
	if len(sceneIDs) == 0 {
		return nil
	}

	scenes, err := r.FindByIDs(ctx, sceneIDs)
	if err != nil {
		return err
	}

	for _, s := range scenes {
		if err := s.LoadLockedFields(ctx, r); err != nil {
			return fmt.Errorf("loading locked fields: %w", err)
		}
		if s.IsLocked(models.SceneLockableFieldGalleries) {
			return fmt.Errorf("scene %q: %w", s.DisplayName(), &LockedFieldsError{Fields: []models.SceneLockableField{models.SceneLockableFieldGalleries}})
		}
	}

	return nil
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
	case models.SceneLockableFieldPerformers:
		return len(partial.PerformerIDs.ImpactedIDs(s.PerformerIDs.List())) > 0
	case models.SceneLockableFieldTags:
		return len(partial.TagIDs.ImpactedIDs(s.TagIDs.List())) > 0
	case models.SceneLockableFieldGalleries:
		return len(partial.GalleryIDs.ImpactedIDs(s.GalleryIDs.List())) > 0
	case models.SceneLockableFieldGroups:
		return groupsChanged(s.Groups.List(), partial.GroupIDs)
	}

	return false
}

// groupsChanged mirrors the add/remove/set semantics of the scenes_groups store.
func groupsChanged(existing []models.GroupsScenes, u *models.UpdateGroupIDs) bool {
	if u == nil {
		return false
	}

	has := func(list []models.GroupsScenes, v models.GroupsScenes, matchIndex bool) bool {
		return slices.ContainsFunc(list, func(e models.GroupsScenes) bool {
			return e.GroupID == v.GroupID && (!matchIndex || intPtrEqual(e.SceneIndex, v.SceneIndex))
		})
	}

	switch u.Mode {
	case models.RelationshipUpdateModeAdd:
		// existing groups are not modified by add
		for _, v := range u.Groups {
			if !has(existing, v, false) {
				return true
			}
		}
	case models.RelationshipUpdateModeRemove:
		for _, v := range u.Groups {
			if has(existing, v, false) {
				return true
			}
		}
	case models.RelationshipUpdateModeSet:
		if len(existing) != len(u.Groups) {
			return true
		}
		for _, v := range u.Groups {
			if !has(existing, v, true) {
				return true
			}
		}
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
