package performer

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeAliases(t *testing.T) {
	name := "Performer Name"
	aliases := []models.PerformerAlias{
		{Alias: " Alias 1 "},
		{Alias: "alias 1"},        // Duplicate (case-insensitive)
		{Alias: "Performer Name"}, // Same as name
		{Alias: "performer name"}, // Same as name (case-insensitive)
		{Alias: ""},               // Empty
		{Alias: "Alias 2"},
	}

	normalized := NormalizeAliases(name, aliases)

	assert.Len(t, normalized, 2)
	assert.Equal(t, "Alias 1", normalized[0].Alias)
	assert.Equal(t, "Alias 2", normalized[1].Alias)
}

func TestGetEffectiveAliases(t *testing.T) {
	existing := []models.PerformerAlias{
		{Alias: "A", IgnoreAutoTag: false},
		{Alias: "B", IgnoreAutoTag: true},
	}
	update := []models.PerformerAlias{
		{Alias: "B", IgnoreAutoTag: false}, // Note: IgnoreAutoTag is false here
		{Alias: "C", IgnoreAutoTag: true},
	}

	t.Run("Mode SET with preservation", func(t *testing.T) {
		effective := GetEffectiveAliases(existing, update, models.RelationshipUpdateModeSet, true)
		assert.Len(t, effective, 2)
		assert.Equal(t, "B", effective[0].Alias)
		assert.True(t, effective[0].IgnoreAutoTag, "Should have preserved IgnoreAutoTag: true from existing")
		assert.Equal(t, "C", effective[1].Alias)
	})

	t.Run("Mode SET without preservation", func(t *testing.T) {
		effective := GetEffectiveAliases(existing, update, models.RelationshipUpdateModeSet, false)
		assert.Len(t, effective, 2)
		assert.Equal(t, "B", effective[0].Alias)
		assert.False(t, effective[0].IgnoreAutoTag, "Should NOT have preserved IgnoreAutoTag")
	})

	t.Run("Mode ADD", func(t *testing.T) {
		effective := GetEffectiveAliases(existing, update, models.RelationshipUpdateModeAdd, false)
		assert.Len(t, effective, 4)
		assert.Equal(t, "A", effective[0].Alias)
		assert.Equal(t, "B", effective[1].Alias)
		assert.Equal(t, "B", effective[2].Alias)
		assert.Equal(t, "C", effective[3].Alias)
	})

	t.Run("Mode REMOVE", func(t *testing.T) {
		effective := GetEffectiveAliases(existing, update, models.RelationshipUpdateModeRemove, false)
		assert.Len(t, effective, 1)
		assert.Equal(t, "A", effective[0].Alias)
	})
}
