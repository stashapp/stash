// package documents represents indexed documents
package documents

import (
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/stashapp/stash/pkg/models"
)

type Scene struct {
	Title   string `json:"title"`
	Details string `json:"details"`

	Date *string `json:"date"`
	Year *int    `json:"year"` // Computed from Date
}

func NewScene(in models.Scene) Scene {
	details := ""
	if in.Details.Valid {
		details = in.Details.String
	}

	var date *string
	var year *int
	if in.Date.Valid {
		date = &in.Date.String
		layout := "2006-01-02"
		t, err := time.Parse(layout, in.Date.String)
		if err != nil {
			year = nil
		}

		y := t.Year()
		year = &y
	}

	return Scene{
		Title:   in.GetTitle(),
		Details: details,
		Date:    date,
		Year:    year,
	}
}

func (s Scene) Type() string {
	return "scene"
}

func BuildSceneIndexMapping() (mapping.IndexMapping, error) {
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	dateMapping := bleve.NewDateTimeFieldMapping()

	sceneMapping := bleve.NewDocumentMapping()

	sceneMapping.AddFieldMappingsAt("title", englishTextFieldMapping)
	sceneMapping.AddFieldMappingsAt("details", englishTextFieldMapping)
	sceneMapping.AddFieldMappingsAt("date", dateMapping)

	indexMapping := bleve.NewIndexMapping()

	indexMapping.AddDocumentMapping("scene", sceneMapping)

	indexMapping.TypeField = "Type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil
}
