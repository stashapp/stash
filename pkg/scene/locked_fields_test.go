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
			PerformerIDs: models.NewRelatedIDs([]int{1, 2}),
			TagIDs:       models.NewRelatedIDs([]int{1}),
			GalleryIDs:   models.NewRelatedIDs([]int{}),
			Groups:       models.NewRelatedGroups([]models.GroupsScenes{{GroupID: 1, SceneIndex: &studioID}}),
		}
	}

	ids := func(mode models.RelationshipUpdateMode, v ...int) *models.UpdateIDs {
		return &models.UpdateIDs{IDs: v, Mode: mode}
	}
	groups := func(mode models.RelationshipUpdateMode, v ...models.GroupsScenes) *models.UpdateGroupIDs {
		return &models.UpdateGroupIDs{Groups: v, Mode: mode}
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
		{"locked performers added", newScene("PERFORMERS"), models.ScenePartial{PerformerIDs: ids(models.RelationshipUpdateModeAdd, 3)}, []models.SceneLockableField{models.SceneLockableFieldPerformers}},
		{"locked performers add existing", newScene("PERFORMERS"), models.ScenePartial{PerformerIDs: ids(models.RelationshipUpdateModeAdd, 1)}, nil},
		{"locked performers same set", newScene("PERFORMERS"), models.ScenePartial{PerformerIDs: ids(models.RelationshipUpdateModeSet, 2, 1)}, nil},
		{"locked tags removed", newScene("TAGS"), models.ScenePartial{TagIDs: ids(models.RelationshipUpdateModeRemove, 1)}, []models.SceneLockableField{models.SceneLockableFieldTags}},
		{"locked tags remove missing", newScene("TAGS"), models.ScenePartial{TagIDs: ids(models.RelationshipUpdateModeRemove, 2)}, nil},
		{"locked galleries set", newScene("GALLERIES"), models.ScenePartial{GalleryIDs: ids(models.RelationshipUpdateModeSet, 1)}, []models.SceneLockableField{models.SceneLockableFieldGalleries}},
		{"locked groups added", newScene("GROUPS"), models.ScenePartial{GroupIDs: groups(models.RelationshipUpdateModeAdd, models.GroupsScenes{GroupID: 2})}, []models.SceneLockableField{models.SceneLockableFieldGroups}},
		{"locked groups add existing", newScene("GROUPS"), models.ScenePartial{GroupIDs: groups(models.RelationshipUpdateModeAdd, models.GroupsScenes{GroupID: 1})}, nil},
		{"locked groups same set", newScene("GROUPS"), models.ScenePartial{GroupIDs: groups(models.RelationshipUpdateModeSet, models.GroupsScenes{GroupID: 1, SceneIndex: &studioID})}, nil},
		{"locked groups index changed", newScene("GROUPS"), models.ScenePartial{GroupIDs: groups(models.RelationshipUpdateModeSet, models.GroupsScenes{GroupID: 1})}, []models.SceneLockableField{models.SceneLockableFieldGroups}},
		{"locked groups removed", newScene("GROUPS"), models.ScenePartial{GroupIDs: groups(models.RelationshipUpdateModeRemove, models.GroupsScenes{GroupID: 1})}, []models.SceneLockableField{models.SceneLockableFieldGroups}},
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

func TestRemoveLockedFields_relationships(t *testing.T) {
	s := &models.Scene{
		LockedFields: models.NewRelatedStrings([]string{"PERFORMERS", "TAGS", "GROUPS", "GALLERIES"}),
		PerformerIDs: models.NewRelatedIDs([]int{}),
		TagIDs:       models.NewRelatedIDs([]int{}),
		GalleryIDs:   models.NewRelatedIDs([]int{}),
		Groups:       models.NewRelatedGroups([]models.GroupsScenes{}),
	}
	add := &models.UpdateIDs{IDs: []int{1}, Mode: models.RelationshipUpdateModeAdd}
	p := models.ScenePartial{
		PerformerIDs: add,
		TagIDs:       add,
		GalleryIDs:   add,
		GroupIDs:     &models.UpdateGroupIDs{Groups: []models.GroupsScenes{{GroupID: 1}}, Mode: models.RelationshipUpdateModeAdd},
	}

	removed := RemoveLockedFields(s, &p)

	assert.Len(t, removed, 4)
	assert.Nil(t, p.PerformerIDs)
	assert.Nil(t, p.TagIDs)
	assert.Nil(t, p.GalleryIDs)
	assert.Nil(t, p.GroupIDs)
}
