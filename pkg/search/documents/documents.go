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
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`

	StashType string `json:"stash_type"`
}

func NewTag(in models.Tag) Tag {
	return Tag{
		ID:   in.ID,
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

	Performer []Performer `json:"performer,omitempty"`
	Tag       []string    `json:"tag,omitempty"`
	TagID     []int       `json:"tag_id,omitempty"`

	StashType string `json:"stash_type"`
}

func NewScene(in models.Scene, inPerformers []Performer, inTags []Tag) Scene {
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

	var tags []string
	var tagIDs []int

	for _, t := range inTags {
		tags = append(tags, t.Name)
		tagIDs = append(tagIDs, t.ID)
	}

	return Scene{
		Title:     in.GetTitle(),
		Details:   details,
		Date:      date,
		Year:      year,
		Performer: inPerformers,
		Tag:       tags,
		TagID:     tagIDs,

		StashType: "scene",
	}
}

func (s Scene) Type() string {
	return "scene"
}

func BuildSceneDocumentMapping() *mapping.DocumentMapping {
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	numericalFieldMapping := bleve.NewNumericFieldMapping()

	dateMapping := bleve.NewDateTimeFieldMapping()

	sceneMapping := bleve.NewDocumentMapping()

	sceneMapping.AddFieldMappingsAt("title", englishTextFieldMapping)
	sceneMapping.AddFieldMappingsAt("details", englishTextFieldMapping)
	sceneMapping.AddFieldMappingsAt("date", dateMapping)

	sceneMapping.AddSubDocumentMapping("tag", BuildTagDocumentMapping())

	// Tags are flattened into the structure
	sceneMapping.AddFieldMappingsAt("tag", englishTextFieldMapping)
	sceneMapping.AddFieldMappingsAt("tag_id", numericalFieldMapping)

	sceneMapping.AddSubDocumentMapping("performer", BuildPerformerDocumentMapping())

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
