// package documents represents indexed documents
package documents

import (
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/stashapp/stash/pkg/models"
)

type Performer struct {
	Name string `json:"name"`

	StashType string `json:"stash_type"`
}

func NewPerformer(in models.Performer) Performer {
	name := ""
	if in.Name.Valid {
		name = in.Name.String
	}

	return Performer{
		StashType: "performer",
		Name:      name,
	}
}

func (p Performer) Type() string {
	return "performer"
}

func BuildPerformerDocumentMapping() *mapping.DocumentMapping {
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	performerMapping := bleve.NewDocumentMapping()

	performerMapping.AddFieldMappingsAt("title", englishTextFieldMapping)

	return performerMapping
}

type Scene struct {
	Title   string `json:"title"`
	Details string `json:"details"`

	Date *string `json:"date"`
	Year *int    `json:"year"` // Computed from Date

	Performer []*Performer `json:"performer"`

	StashType string `json:"stash_type"`
}

func NewScene(in models.Scene, performers []*Performer) Scene {
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
		Title:     in.GetTitle(),
		Details:   details,
		Date:      date,
		Year:      year,
		Performer: performers,

		StashType: "scene",
	}
}

func (s Scene) Type() string {
	return "scene"
}

func BuildSceneDocumentMapping() *mapping.DocumentMapping {
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	dateMapping := bleve.NewDateTimeFieldMapping()

	sceneMapping := bleve.NewDocumentMapping()

	sceneMapping.AddFieldMappingsAt("title", englishTextFieldMapping)
	sceneMapping.AddFieldMappingsAt("details", englishTextFieldMapping)
	sceneMapping.AddFieldMappingsAt("date", dateMapping)

	sceneMapping.AddSubDocumentMapping("performer", BuildPerformerDocumentMapping())
	return sceneMapping
}

func BuildIndexMapping() (mapping.IndexMapping, error) {
	sceneMapping := BuildSceneDocumentMapping()
	performerMapping := BuildPerformerDocumentMapping()
	indexMapping := bleve.NewIndexMapping()

	indexMapping.AddDocumentMapping("scene", sceneMapping)
	indexMapping.AddDocumentMapping("performer", performerMapping)

	indexMapping.TypeField = "Type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil
}
