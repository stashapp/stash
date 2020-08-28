package performer

import (
	"database/sql"
	"errors"

	"github.com/stashapp/stash/pkg/manager/jsonschema"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/models/modelstest"
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

const (
	performerID = 1
	noImageID   = 2
	errImageID  = 3
)

const (
	performerName = "testPerformer"
	url           = "url"
	aliases       = "aliases"
	careerLength  = "careerLength"
	country       = "country"
	ethnicity     = "ethnicity"
	eyeColor      = "eyeColor"
	fakeTits      = "fakeTits"
	gender        = "gender"
	height        = "height"
	instagram     = "instagram"
	measurements  = "measurements"
	piercings     = "piercings"
	tattoos       = "tattoos"
	twitter       = "twitter"
)

var imageBytes = []byte("imageBytes")

const image = "aW1hZ2VCeXRlcw=="

var birthDate = models.SQLiteDate{
	String: "2001-01-01",
	Valid:  true,
}
var createTime time.Time = time.Date(2001, 01, 01, 0, 0, 0, 0, time.UTC)
var updateTime time.Time = time.Date(2002, 01, 01, 0, 0, 0, 0, time.UTC)

func createFullPerformer(id int) models.Performer {
	return models.Performer{
		ID:           id,
		Name:         modelstest.NullString(performerName),
		URL:          modelstest.NullString(url),
		Aliases:      modelstest.NullString(aliases),
		Birthdate:    birthDate,
		CareerLength: modelstest.NullString(careerLength),
		Country:      modelstest.NullString(country),
		Ethnicity:    modelstest.NullString(ethnicity),
		EyeColor:     modelstest.NullString(eyeColor),
		FakeTits:     modelstest.NullString(fakeTits),
		Favorite: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
		Gender:       modelstest.NullString(gender),
		Height:       modelstest.NullString(height),
		Instagram:    modelstest.NullString(instagram),
		Measurements: modelstest.NullString(measurements),
		Piercings:    modelstest.NullString(piercings),
		Tattoos:      modelstest.NullString(tattoos),
		Twitter:      modelstest.NullString(twitter),
		CreatedAt: models.SQLiteTimestamp{
			Timestamp: createTime,
		},
		UpdatedAt: models.SQLiteTimestamp{
			Timestamp: updateTime,
		},
	}
}

func createEmptyPerformer(id int) models.Performer {
	return models.Performer{
		ID: id,
		CreatedAt: models.SQLiteTimestamp{
			Timestamp: createTime,
		},
		UpdatedAt: models.SQLiteTimestamp{
			Timestamp: updateTime,
		},
	}
}

func createFullJSONPerformer(image string) *jsonschema.Performer {
	return &jsonschema.Performer{
		Name:         performerName,
		URL:          url,
		Aliases:      aliases,
		Birthdate:    birthDate.String,
		CareerLength: careerLength,
		Country:      country,
		Ethnicity:    ethnicity,
		EyeColor:     eyeColor,
		FakeTits:     fakeTits,
		Favorite:     true,
		Gender:       gender,
		Height:       height,
		Instagram:    instagram,
		Measurements: measurements,
		Piercings:    piercings,
		Tattoos:      tattoos,
		Twitter:      twitter,
		CreatedAt: models.JSONTime{
			Time: createTime,
		},
		UpdatedAt: models.JSONTime{
			Time: updateTime,
		},
		Image: image,
	}
}

func createEmptyJSONPerformer() *jsonschema.Performer {
	return &jsonschema.Performer{
		CreatedAt: models.JSONTime{
			Time: createTime,
		},
		UpdatedAt: models.JSONTime{
			Time: updateTime,
		},
	}
}

type testScenario struct {
	input    models.Performer
	expected *jsonschema.Performer
	err      bool
}

var scenarios []testScenario

func initTestTable() {
	scenarios = []testScenario{
		testScenario{
			createFullPerformer(performerID),
			createFullJSONPerformer(image),
			false,
		},
		testScenario{
			createEmptyPerformer(noImageID),
			createEmptyJSONPerformer(),
			false,
		},
		testScenario{
			createFullPerformer(errImageID),
			nil,
			true,
		},
	}
}

func TestToJSON(t *testing.T) {
	initTestTable()

	mockPerformerReader := &mocks.PerformerReaderWriter{}

	imageErr := errors.New("error getting image")

	mockPerformerReader.On("GetPerformerImage", performerID).Return(imageBytes, nil).Once()
	mockPerformerReader.On("GetPerformerImage", noImageID).Return(nil, nil).Once()
	mockPerformerReader.On("GetPerformerImage", errImageID).Return(nil, imageErr).Once()

	for i, s := range scenarios {
		tag := s.input
		json, err := ToJSON(mockPerformerReader, &tag)

		if !s.err && err != nil {
			t.Errorf("[%d] unexpected error: %s", i, err.Error())
		} else if s.err && err == nil {
			t.Errorf("[%d] expected error not returned", i)
		} else {
			assert.Equal(t, s.expected, json, "[%d]", i)
		}
	}

	mockPerformerReader.AssertExpectations(t)
}
