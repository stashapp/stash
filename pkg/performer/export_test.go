package performer

import (
	"database/sql"
	"errors"

	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/models/mocks"
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
	performerName  = "testPerformer"
	url            = "url"
	aliases        = "aliases"
	careerLength   = "careerLength"
	country        = "country"
	ethnicity      = "ethnicity"
	eyeColor       = "eyeColor"
	fakeTits       = "fakeTits"
	gender         = "gender"
	height         = "height"
	instagram      = "instagram"
	measurements   = "measurements"
	piercings      = "piercings"
	tattoos        = "tattoos"
	twitter        = "twitter"
	rating         = 5
	details        = "details"
	hairColor      = "hairColor"
	weight         = 60
	autoTagIgnored = true
)

var imageBytes = []byte("imageBytes")

var stashID = models.StashID{
	StashID:  "StashID",
	Endpoint: "Endpoint",
}
var stashIDs = []*models.StashID{
	&stashID,
}

const image = "aW1hZ2VCeXRlcw=="

var birthDate = models.SQLiteDate{
	String: "2001-01-01",
	Valid:  true,
}
var deathDate = models.SQLiteDate{
	String: "2021-02-02",
	Valid:  true,
}

var (
	createTime = time.Date(2001, 01, 01, 0, 0, 0, 0, time.Local)
	updateTime = time.Date(2002, 01, 01, 0, 0, 0, 0, time.Local)
)

func createFullPerformer(id int, name string) *models.Performer {
	return &models.Performer{
		ID:           id,
		Name:         models.NullString(name),
		Checksum:     md5.FromString(name),
		URL:          models.NullString(url),
		Aliases:      models.NullString(aliases),
		Birthdate:    birthDate,
		CareerLength: models.NullString(careerLength),
		Country:      models.NullString(country),
		Ethnicity:    models.NullString(ethnicity),
		EyeColor:     models.NullString(eyeColor),
		FakeTits:     models.NullString(fakeTits),
		Favorite: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
		Gender:       models.NullString(gender),
		Height:       models.NullString(height),
		Instagram:    models.NullString(instagram),
		Measurements: models.NullString(measurements),
		Piercings:    models.NullString(piercings),
		Tattoos:      models.NullString(tattoos),
		Twitter:      models.NullString(twitter),
		CreatedAt: models.SQLiteTimestamp{
			Timestamp: createTime,
		},
		UpdatedAt: models.SQLiteTimestamp{
			Timestamp: updateTime,
		},
		Rating:    models.NullInt64(rating),
		Details:   models.NullString(details),
		DeathDate: deathDate,
		HairColor: models.NullString(hairColor),
		Weight: sql.NullInt64{
			Int64: weight,
			Valid: true,
		},
		IgnoreAutoTag: autoTagIgnored,
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

func createFullJSONPerformer(name string, image string) *jsonschema.Performer {
	return &jsonschema.Performer{
		Name:         name,
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
		CreatedAt: json.JSONTime{
			Time: createTime,
		},
		UpdatedAt: json.JSONTime{
			Time: updateTime,
		},
		Rating:    rating,
		Image:     image,
		Details:   details,
		DeathDate: deathDate.String,
		HairColor: hairColor,
		Weight:    weight,
		StashIDs: []models.StashID{
			stashID,
		},
		IgnoreAutoTag: autoTagIgnored,
	}
}

func createEmptyJSONPerformer() *jsonschema.Performer {
	return &jsonschema.Performer{
		CreatedAt: json.JSONTime{
			Time: createTime,
		},
		UpdatedAt: json.JSONTime{
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
		{
			*createFullPerformer(performerID, performerName),
			createFullJSONPerformer(performerName, image),
			false,
		},
		{
			createEmptyPerformer(noImageID),
			createEmptyJSONPerformer(),
			false,
		},
		{
			*createFullPerformer(errImageID, performerName),
			nil,
			true,
		},
	}
}

func TestToJSON(t *testing.T) {
	initTestTable()

	mockPerformerReader := &mocks.PerformerReaderWriter{}

	imageErr := errors.New("error getting image")

	mockPerformerReader.On("GetImage", testCtx, performerID).Return(imageBytes, nil).Once()
	mockPerformerReader.On("GetImage", testCtx, noImageID).Return(nil, nil).Once()
	mockPerformerReader.On("GetImage", testCtx, errImageID).Return(nil, imageErr).Once()

	mockPerformerReader.On("GetStashIDs", testCtx, performerID).Return(stashIDs, nil).Once()
	mockPerformerReader.On("GetStashIDs", testCtx, noImageID).Return(nil, nil).Once()

	for i, s := range scenarios {
		tag := s.input
		json, err := ToJSON(testCtx, mockPerformerReader, &tag)

		switch {
		case !s.err && err != nil:
			t.Errorf("[%d] unexpected error: %s", i, err.Error())
		case s.err && err == nil:
			t.Errorf("[%d] expected error not returned", i)
		default:
			assert.Equal(t, s.expected, json, "[%d]", i)
		}
	}

	mockPerformerReader.AssertExpectations(t)
}
