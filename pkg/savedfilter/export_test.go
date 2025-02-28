package savedfilter

import (
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"

	"testing"
)

const (
	savedFilterID = 1
	noImageID     = 2
	errImageID    = 3
	errAliasID    = 4
	withParentsID = 5
	errParentsID  = 6
)

const (
	filterName = "testFilter"
	mode       = models.FilterModeGalleries
)

var (
	findFilter   = models.FindFilterType{}
	objectFilter = make(map[string]interface{})
	uiOptions    = make(map[string]interface{})
)

func createSavedFilter(id int) models.SavedFilter {
	return models.SavedFilter{
		ID:           id,
		Name:         filterName,
		Mode:         mode,
		FindFilter:   &findFilter,
		ObjectFilter: objectFilter,
		UIOptions:    uiOptions,
	}
}

func createJSONSavedFilter() *jsonschema.SavedFilter {
	return &jsonschema.SavedFilter{
		Name:         filterName,
		Mode:         mode,
		FindFilter:   &findFilter,
		ObjectFilter: objectFilter,
		UIOptions:    uiOptions,
	}
}

type testScenario struct {
	savedFilter models.SavedFilter
	expected    *jsonschema.SavedFilter
	err         bool
}

var scenarios []testScenario

func initTestTable() {
	scenarios = []testScenario{
		{
			createSavedFilter(savedFilterID),
			createJSONSavedFilter(),
			false,
		},
	}
}

func TestToJSON(t *testing.T) {
	initTestTable()

	db := mocks.NewDatabase()

	for i, s := range scenarios {
		savedFilter := s.savedFilter
		json, err := ToJSON(testCtx, &savedFilter)

		switch {
		case !s.err && err != nil:
			t.Errorf("[%d] unexpected error: %s", i, err.Error())
		case s.err && err == nil:
			t.Errorf("[%d] expected error not returned", i)
		default:
			assert.Equal(t, s.expected, json, "[%d]", i)
		}
	}

	db.AssertExpectations(t)
}
