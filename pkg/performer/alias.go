package performer

import (
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

// NormalizeAliases trims whitespace, deduplicates aliases (case-insensitively),
// and ensures aliases are not empty and do not match the performer's name.
func NormalizeAliases(performerName string, aliases []models.PerformerAlias) []models.PerformerAlias {
	sanitized := []models.PerformerAlias{}
	seen := make(map[string]bool)
	nameLower := strings.ToLower(strings.TrimSpace(performerName))

	for _, a := range aliases {
		trimmed := strings.TrimSpace(a.Alias)
		lower := strings.ToLower(trimmed)
		if trimmed != "" && lower != nameLower && !seen[lower] {
			seen[lower] = true
			a.Alias = trimmed
			sanitized = append(sanitized, a)
		}
	}

	return sanitized
}

// GetEffectiveAliases calculates the final list of aliases based on the update mode
// and optionally preserves the IgnoreAutoTag state for existing aliases.
func GetEffectiveAliases(existing []models.PerformerAlias, update []models.PerformerAlias, mode models.RelationshipUpdateMode, preserveIgnore bool) []models.PerformerAlias {
	var effective []models.PerformerAlias

	switch mode {
	case models.RelationshipUpdateModeSet:
		if preserveIgnore {
			existingMap := make(map[string]bool)
			for _, e := range existing {
				existingMap[e.Alias] = e.IgnoreAutoTag
			}

			for _, u := range update {
				if ignore, ok := existingMap[u.Alias]; ok {
					u.IgnoreAutoTag = ignore
				}
				effective = append(effective, u)
			}
		} else {
			effective = update
		}
	case models.RelationshipUpdateModeAdd:
		effective = append(effective, existing...)
		effective = append(effective, update...)
	case models.RelationshipUpdateModeRemove:
		updateMap := make(map[string]bool)
		for _, u := range update {
			updateMap[u.Alias] = true
		}

		for _, e := range existing {
			if !updateMap[e.Alias] {
				effective = append(effective, e)
			}
		}
	}

	return effective
}
