package performer

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

type NameExistsError struct {
	Name           string
	Disambiguation string
}

func (e *NameExistsError) Error() string {
	if e.Disambiguation != "" {
		return fmt.Sprintf("performer with name '%s' and disambiguation '%s' already exists", e.Name, e.Disambiguation)
	}
	return fmt.Sprintf("performer with name '%s' already exists", e.Name)
}

// EnsureNameUnique returns an error if the performer name and disambiguation provided
// is used by another performer
func EnsureNameUnique(ctx context.Context, id int, name string, disambig string, qb models.PerformerReaderWriter) error {
	performerFilter := models.PerformerFilterType{
		Name: &models.StringCriterionInput{
			Value:    name,
			Modifier: models.CriterionModifierEquals,
		},
	}

	if disambig != "" {
		performerFilter.Disambiguation = &models.StringCriterionInput{
			Value:    disambig,
			Modifier: models.CriterionModifierEquals,
		}
	}

	pp := 1
	findFilter := models.FindFilterType{
		PerPage: &pp,
	}

	existing, _, err := qb.Query(ctx, &performerFilter, &findFilter)
	if err != nil {
		return err
	}

	if len(existing) > 0 && existing[0].ID != id {
		return &NameExistsError{
			Name:           name,
			Disambiguation: disambig,
		}
	}

	return nil
}
