package performer

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

var (
	ErrNameMissing = errors.New("performer name must not be blank")
)

type NotFoundError struct {
	id int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("performer with id %d not found", e.id)
}

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

type DuplicateAliasError struct {
	Alias string
}

func (e *DuplicateAliasError) Error() string {
	return fmt.Sprintf("performer contains duplicate alias '%s'", e.Alias)
}

type DeathDateError struct {
	Birthdate models.Date
	DeathDate models.Date
}

func (e *DeathDateError) Error() string {
	return fmt.Sprintf("death date %s should be after birthdate %s", e.DeathDate, e.Birthdate)
}

func ValidateCreate(ctx context.Context, performer models.Performer, qb models.PerformerReader) error {
	if err := ValidateName(ctx, performer.Name, performer.Disambiguation, qb); err != nil {
		return err
	}

	if err := ValidateAliases(performer.Name, performer.Aliases); err != nil {
		return err
	}

	if err := ValidateDeathDate(performer.Birthdate, performer.DeathDate); err != nil {
		return err
	}

	return nil
}

func ValidateUpdate(ctx context.Context, id int, partial models.PerformerPartial, qb models.PerformerReader) error {
	existing, err := qb.Find(ctx, id)
	if err != nil {
		return err
	}

	if existing == nil {
		return &NotFoundError{id}
	}

	if err := ValidateUpdateName(ctx, *existing, partial.Name, partial.Disambiguation, qb); err != nil {
		return err
	}

	if err := existing.LoadAliases(ctx, qb); err != nil {
		return err
	}
	if err := ValidateUpdateAliases(*existing, partial.Name, partial.Aliases); err != nil {
		return err
	}

	if err := ValidateUpdateDeathDate(*existing, partial.Birthdate, partial.DeathDate); err != nil {
		return err
	}

	return nil
}

func validateName(ctx context.Context, name string, disambig string, existingID *int, qb models.PerformerQueryer) error {
	performerFilter := models.PerformerFilterType{
		Name: &models.StringCriterionInput{
			Value:    name,
			Modifier: models.CriterionModifierEquals,
		},
	}

	modifier := models.CriterionModifierIsNull

	if disambig != "" {
		modifier = models.CriterionModifierEquals
	}

	performerFilter.Disambiguation = &models.StringCriterionInput{
		Value:    disambig,
		Modifier: modifier,
	}

	if existingID == nil {
		// creating: error if any existing performer matches

		pp := 1
		findFilter := models.FindFilterType{
			PerPage: &pp,
		}

		count, err := qb.QueryCount(ctx, &performerFilter, &findFilter)
		if err != nil {
			return err
		}

		if count > 0 {
			return &NameExistsError{
				Name:           name,
				Disambiguation: disambig,
			}
		}

		return nil
	} else {
		// updating: check for matches, but ignore self

		pp := 2
		findFilter := models.FindFilterType{
			PerPage: &pp,
		}

		conflicts, _, err := qb.Query(ctx, &performerFilter, &findFilter)
		if err != nil {
			return err
		}

		if len(conflicts) > 0 {
			// valid if the only conflict is the existing performer
			if len(conflicts) > 1 || conflicts[0].ID != *existingID {
				return &NameExistsError{
					Name:           name,
					Disambiguation: disambig,
				}
			}
		}

		return nil
	}
}

// ValidateName returns an error if the performer name and disambiguation provided is used by another performer.
func ValidateName(ctx context.Context, name string, disambig string, qb models.PerformerQueryer) error {
	if name == "" {
		return ErrNameMissing
	}

	return validateName(ctx, name, disambig, nil, qb)
}

// ValidateUpdateName performs the same check as ValidateName, but is used when modifying an existing performer.
func ValidateUpdateName(ctx context.Context, existing models.Performer, name models.OptionalString, disambig models.OptionalString, qb models.PerformerQueryer) error {
	// if neither name nor disambig is set, don't check anything
	if !name.Set && !disambig.Set {
		return nil
	}

	newName := existing.Name
	if name.Set {
		newName = name.Value
	}

	if newName == "" {
		return ErrNameMissing
	}

	newDisambig := existing.Disambiguation
	if disambig.Set {
		newDisambig = disambig.Value
	}

	return validateName(ctx, newName, newDisambig, &existing.ID, qb)
}

func ValidateAliases(name string, aliases models.RelatedStrings) error {
	if !aliases.Loaded() {
		return nil
	}

	m := make(map[string]bool)
	nameL := strings.ToLower(name)
	m[nameL] = true

	for _, alias := range aliases.List() {
		aliasL := strings.ToLower(alias)
		if m[aliasL] {
			return &DuplicateAliasError{alias}
		}
		m[aliasL] = true
	}

	return nil
}

func ValidateUpdateAliases(existing models.Performer, name models.OptionalString, aliases *models.UpdateStrings) error {
	// if neither name nor aliases is set, don't check anything
	if !name.Set && aliases == nil {
		return nil
	}

	newName := existing.Name
	if name.Set {
		newName = name.Value
	}

	newAliases := aliases.Apply(existing.Aliases.List())

	return ValidateAliases(newName, models.NewRelatedStrings(newAliases))
}

// ValidateDeathDate returns an error if the birthdate is after the death date.
func ValidateDeathDate(birthdate *models.Date, deathDate *models.Date) error {
	if birthdate == nil || deathDate == nil {
		return nil
	}

	if birthdate.After(*deathDate) {
		return &DeathDateError{Birthdate: *birthdate, DeathDate: *deathDate}
	}

	return nil
}

// ValidateUpdateDeathDate performs the same check as ValidateDeathDate, but is used when modifying an existing performer.
func ValidateUpdateDeathDate(existing models.Performer, birthdate models.OptionalDate, deathDate models.OptionalDate) error {
	// if neither birthdate nor deathDate is set, don't check anything
	if !birthdate.Set && !deathDate.Set {
		return nil
	}

	newBirthdate := existing.Birthdate
	if birthdate.Set {
		newBirthdate = birthdate.Ptr()
	}

	newDeathDate := existing.DeathDate
	if deathDate.Set {
		newDeathDate = deathDate.Ptr()
	}

	return ValidateDeathDate(newBirthdate, newDeathDate)
}
