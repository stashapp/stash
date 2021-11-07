// package documents represents indexed documents
package documents

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/stashapp/stash/pkg/models"
)

type Scene struct {
	Title   string `json:"title"`
	Details string `json:"details"`

	Date *string `json:"date"`
}

func NewScene(in models.Scene) Scene {
	details := ""
	if in.Details.Valid {
		details = in.Details.String
	}

	var date *string
	if in.Date.Valid {
		date = &in.Date.String
	}

	return Scene{
		Title:   in.GetTitle(),
		Details: details,
		Date:    date,
	}
}

func (s Scene) Classifier() string {
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
