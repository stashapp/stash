package scene

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestLockViolations(t *testing.T) {
	title := models.SceneLockableFieldTitle
	studioID := 1
	date, _ := models.ParseDate("2020-01-02")

	newScene := func(locked ...string) *models.Scene {
		if locked == nil {
			locked = []string{}
		}
		return &models.Scene{
			Title:        "title",
			Date:         &date,
			StudioID:     &studioID,
			URLs:         models.NewRelatedStrings([]string{"a"}),
			LockedFields: models.NewRelatedStrings(locked),
		}
	}

	setLocks := func(p models.ScenePartial, locked ...string) models.ScenePartial {
		p.LockedFields = &models.UpdateStrings{Values: locked, Mode: models.RelationshipUpdateModeSet}
		return p
	}

	tests := []struct {
		name    string
		scene   *models.Scene
		partial models.ScenePartial
		want    []models.SceneLockableField
	}{
		{"unlocked change", newScene(), models.ScenePartial{Title: models.NewOptionalString("new")}, nil},
		{"locked change", newScene("TITLE"), models.ScenePartial{Title: models.NewOptionalString("new")}, []models.SceneLockableField{title}},
		{"locked clear", newScene("TITLE"), models.ScenePartial{Title: models.NewOptionalStringPtr(nil)}, []models.SceneLockableField{title}},
		{"locked unchanged value", newScene("TITLE"), models.ScenePartial{Title: models.NewOptionalString("title")}, nil},
		{"unlock and change", newScene("TITLE"), setLocks(models.ScenePartial{Title: models.NewOptionalString("new")}), nil},
		{"change and lock", newScene(), setLocks(models.ScenePartial{Title: models.NewOptionalString("new")}, "TITLE"), nil},
		{"locked other field changed", newScene("CODE"), models.ScenePartial{Title: models.NewOptionalString("new")}, nil},
		{"locked studio cleared", newScene("STUDIO"), models.ScenePartial{StudioID: models.NewOptionalIntPtr(nil)}, []models.SceneLockableField{models.SceneLockableFieldStudio}},
		{"locked date changed", newScene("DATE"), models.ScenePartial{Date: models.NewOptionalDatePtr(nil)}, []models.SceneLockableField{models.SceneLockableFieldDate}},
		{"locked urls added", newScene("URLS"), models.ScenePartial{URLs: &models.UpdateStrings{Values: []string{"b"}, Mode: models.RelationshipUpdateModeAdd}}, []models.SceneLockableField{models.SceneLockableFieldUrls}},
		{"locked urls same", newScene("URLS"), models.ScenePartial{URLs: &models.UpdateStrings{Values: []string{"a"}, Mode: models.RelationshipUpdateModeSet}}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, LockViolations(tt.scene, tt.partial))
		})
	}
}

func TestRemoveLockedFields(t *testing.T) {
	s := &models.Scene{
		Title:        "title",
		LockedFields: models.NewRelatedStrings([]string{"TITLE"}),
	}
	p := models.ScenePartial{
		Title:   models.NewOptionalString("new"),
		Details: models.NewOptionalString("details"),
	}

	removed := RemoveLockedFields(s, &p)

	assert.Equal(t, []models.SceneLockableField{models.SceneLockableFieldTitle}, removed)
	assert.False(t, p.Title.Set)
	assert.True(t, p.Details.Set)
}
