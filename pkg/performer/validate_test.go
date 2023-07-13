package performer

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestValidateDeathDate(t *testing.T) {
	assert := assert.New(t)

	date1 := "2001-01-01"
	date2 := "2002-01-01"
	date3 := "2003-01-01"
	date4 := "2004-01-01"
	empty := ""

	md2, _ := models.ParseDate(date2)
	md3, _ := models.ParseDate(date3)

	emptyPerformer := models.Performer{}
	invalidPerformer := models.Performer{
		Birthdate: &md3,
		DeathDate: &md2,
	}
	validPerformer := models.Performer{
		Birthdate: &md2,
		DeathDate: &md3,
	}

	// nil values should always return nil
	assert.Nil(ValidateDeathDate(nil, nil, &date1))
	assert.Nil(ValidateDeathDate(nil, &date2, nil))
	assert.Nil(ValidateDeathDate(&emptyPerformer, nil, &date1))
	assert.Nil(ValidateDeathDate(&emptyPerformer, &date2, nil))

	// empty strings should always return nil
	assert.Nil(ValidateDeathDate(nil, &empty, &date1))
	assert.Nil(ValidateDeathDate(nil, &date2, &empty))
	assert.Nil(ValidateDeathDate(&emptyPerformer, &empty, &date1))
	assert.Nil(ValidateDeathDate(&emptyPerformer, &date2, &empty))
	assert.Nil(ValidateDeathDate(&validPerformer, &empty, &date1))
	assert.Nil(ValidateDeathDate(&validPerformer, &date2, &empty))

	// nil inputs should return nil even if performer is invalid
	assert.Nil(ValidateDeathDate(&invalidPerformer, nil, nil))

	// invalid input values should return error
	assert.NotNil(ValidateDeathDate(nil, &date2, &date1))
	assert.NotNil(ValidateDeathDate(&validPerformer, &date2, &date1))

	// valid input values should return nil
	assert.Nil(ValidateDeathDate(nil, &date1, &date2))

	// use performer values if performer set and values available
	assert.NotNil(ValidateDeathDate(&validPerformer, nil, &date1))
	assert.NotNil(ValidateDeathDate(&validPerformer, &date4, nil))
	assert.Nil(ValidateDeathDate(&validPerformer, nil, &date4))
	assert.Nil(ValidateDeathDate(&validPerformer, &date1, nil))
}
