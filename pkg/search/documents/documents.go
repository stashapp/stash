// package documents represents indexed documents
package documents

import (
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/stashapp/stash/pkg/models"
)

type Tag struct {
	Name string `json:"name,omitempty"`

	StashType string `json:"stash_type"`
}

func NewTag(in models.Tag) Tag {
	return Tag{
		Name: in.Name,

		StashType: "tag",
	}
}

func (t Tag) Type() string {
	return "tag"
}

func BuildTagDocumentMapping() *mapping.DocumentMapping {
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	performerMapping := bleve.NewDocumentMapping()

	performerMapping.AddFieldMappingsAt("name", englishTextFieldMapping)

	return performerMapping
}

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

	performerMapping.AddFieldMappingsAt("name", englishTextFieldMapping)

	return performerMapping
}

type Scene struct {
	Title   string `json:"title,omitempty"`
	Details string `json:"details,omitempty"`

	Date *string `json:"date,omitempty"`
	Year *int    `json:"year,omitempty"` // Computed from Date

	Performer []*Performer `json:"performer,omitempty"`
	Tag       []*Tag       `json:"tag,omitempty"`

	StashType string `json:"stash_type"`
}

func NewScene(in models.Scene, performers []*Performer, tags []*Tag) Scene {
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
		Tag:       tags,

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
	sceneMapping.AddSubDocumentMapping("tag", BuildTagDocumentMapping())

	return sceneMapping
}

func BuildIndexMapping() (mapping.IndexMapping, error) {
	sceneMapping := BuildSceneDocumentMapping()
	performerMapping := BuildPerformerDocumentMapping()
	tagMapping := BuildTagDocumentMapping()

	indexMapping := bleve.NewIndexMapping()

	indexMapping.AddDocumentMapping("scene", sceneMapping)
	indexMapping.AddDocumentMapping("performer", performerMapping)
	indexMapping.AddDocumentMapping("tag", tagMapping)
	indexMapping.TypeField = "Type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil
}
