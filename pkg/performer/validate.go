package performer

import (
	"errors"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func ValidateDeathDate(performer *models.Performer, birthdate *string, deathDate *string) error {
	// don't validate existing values
	if birthdate == nil && deathDate == nil {
		return nil
	}

	if performer != nil {
		if birthdate == nil && performer.Birthdate.Valid {
			birthdate = &performer.Birthdate.String
		}
		if deathDate == nil && performer.DeathDate.Valid {
			deathDate = &performer.DeathDate.String
		}
	}

	if birthdate == nil || deathDate == nil || *birthdate == "" || *deathDate == "" {
		return nil
	}

	f, _ := utils.ParseDateStringAsTime(*birthdate)
	t, _ := utils.ParseDateStringAsTime(*deathDate)

	if f.After(t) {
		return errors.New("the date of death should be higher than the date of birth")
	}

	return nil
}
