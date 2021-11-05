// package documents represents indexed documents
package documents

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/stashapp/stash/pkg/models"
)

type Scene struct {
	Title   string
	Details string
}

func NewScene(in models.Scene) Scene {
	details := ""
	if in.Details.Valid {
		details = in.Details.String
	}

	return Scene{
		Title:   in.GetTitle(),
		Details: details,
	}
}

func (s Scene) Classifier() string {
	return "Scene"
}

func BuildSceneIndexMapping() (mapping.IndexMapping, error) {
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	sceneMapping := bleve.NewDocumentMapping()

	sceneMapping.AddFieldMappingsAt("Title", englishTextFieldMapping)
	sceneMapping.AddFieldMappingsAt("Details", englishTextFieldMapping)

	indexMapping := bleve.NewIndexMapping()

	indexMapping.AddDocumentMapping("Scene", sceneMapping)

	indexMapping.TypeField = "Type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil
}
