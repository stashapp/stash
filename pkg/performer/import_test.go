package performer

import (
	"errors"

	"github.com/stashapp/stash/pkg/manager/jsonschema"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"

	"testing"
)

const invalidImage = "aW1hZ2VCeXRlcw&&"

const (
	existingPerformerName       = "existingPerformerName"
	errExistingPerformerName    = "errExistingPerformerName"
	errCreatePerformerName      = "errCreatePerformerName"
	errUpdatePerformerName      = "errUpdatePerformerName"
	errUpdateImagePerformerName = "errUpdateImagePerformerName"
)

type importTestScenario struct {
	input              *jsonschema.Performer
	duplicateBehaviour models.ImportDuplicateEnum
	created            *models.Performer
	updated            *models.Performer
	err                bool
}

var importScenarios = []importTestScenario{
	// create scenario
	importTestScenario{
		createFullJSONPerformer(performerName, image),
		models.ImportDuplicateEnumIgnore,
		createFullPerformer(performerID, performerName),
		nil,
		false,
	},
	// existing scenarios
	importTestScenario{
		createFullJSONPerformer(existingPerformerName, image),
		models.ImportDuplicateEnumIgnore,
		nil,
		nil,
		false,
	},
	importTestScenario{
		createFullJSONPerformer(existingPerformerName, image),
		models.ImportDuplicateEnumFail,
		nil,
		nil,
		true,
	},
	importTestScenario{
		createFullJSONPerformer(existingPerformerName, image),
		models.ImportDuplicateEnumOverwrite,
		nil,
		createFullPerformer(performerID, existingPerformerName),
		false,
	},
	// find by names error failure
	importTestScenario{
		createFullJSONPerformer(errExistingPerformerName, image),
		models.ImportDuplicateEnumIgnore,
		nil,
		nil,
		true,
	},
	// invalid image test case
	importTestScenario{
		createFullJSONPerformer(performerName, invalidImage),
		models.ImportDuplicateEnumIgnore,
		createFullPerformer(performerID, performerName),
		nil,
		true,
	},
	importTestScenario{
		createFullJSONPerformer(errCreatePerformerName, image),
		models.ImportDuplicateEnumIgnore,
		createFullPerformer(0, errCreatePerformerName),
		nil,
		true,
	},
	importTestScenario{
		createFullJSONPerformer(errUpdatePerformerName, image),
		models.ImportDuplicateEnumOverwrite,
		nil,
		createFullPerformer(performerID, errUpdatePerformerName),
		true,
	},
	importTestScenario{
		createFullJSONPerformer(errUpdateImagePerformerName, image),
		models.ImportDuplicateEnumIgnore,
		createFullPerformer(performerID, errUpdateImagePerformerName),
		nil,
		true,
	},
}

func TestImport(t *testing.T) {
	mockPerformer := &models.Performer{
		ID: performerID,
	}

	imageErr := errors.New("error getting image")

	for i, s := range importScenarios {
		mockPerformerReader := &mocks.PerformerReaderWriter{}

		mockPerformerReader.On("FindByNames", []string{performerName}, false).Return(nil, nil)
		mockPerformerReader.On("FindByNames", []string{existingPerformerName}, false).Return([]*models.Performer{
			mockPerformer,
		}, nil)
		mockPerformerReader.On("FindByNames", []string{errExistingPerformerName}, false).Return(nil, errors.New("FindByNames error"))
		mockPerformerReader.On("FindByNames", []string{errCreatePerformerName}, false).Return(nil, nil)
		mockPerformerReader.On("FindByNames", []string{errUpdatePerformerName}, false).Return([]*models.Performer{
			mockPerformer,
		}, nil)
		mockPerformerReader.On("FindByNames", []string{errUpdateImagePerformerName}, false).Return(nil, nil)

		importer := Importer{
			ReaderWriter:       mockPerformerReader,
			DuplicateBehaviour: s.duplicateBehaviour,
		}

		input := s.input
		if input.Name == errCreatePerformerName {
			mockPerformerReader.On("Create", *s.created).Return(nil, errors.New("Create error"))
		} else if input.Name == errUpdatePerformerName {
			mockPerformerReader.On("Update", *s.updated).Return(nil, errors.New("Update error"))
		} else {
			if s.created != nil {
				created := *s.created
				created.ID = 0
				mockPerformerReader.On("Create", created).Return(s.created, nil).Once()
			} else if s.updated != nil {
				mockPerformerReader.On("Update", *s.updated).Return(nil, nil).Once()
			}

			if input.Name == errUpdateImagePerformerName {
				mockPerformerReader.On("UpdatePerformerImage", performerID, imageBytes).Return(imageErr)
			} else {
				mockPerformerReader.On("UpdatePerformerImage", performerID, imageBytes).Return(nil)
			}
		}

		err := importer.Import(input)

		if !s.err && err != nil {
			t.Errorf("[%d] unexpected error: %s", i, err.Error())
		} else if s.err && err == nil {
			t.Errorf("[%d] expected error not returned", i)
		}
	}
}
