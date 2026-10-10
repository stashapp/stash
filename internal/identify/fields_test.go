package identify

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func Test_setPartialDate_invalidDate(t *testing.T) {
	var target models.OptionalDate
	invalidDate := "invalid"

	setPartialDate(&target, &invalidDate, nil, &FieldOptions{
		Strategy: FieldStrategyOverwrite,
	}, "date")

	assert.Equal(t, models.OptionalDate{}, target)
}
