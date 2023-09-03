package heresphere

import (
	"github.com/stashapp/stash/pkg/models"
)

type StringSingletonInput struct {
	Value string `json:"value"`
}

// Technically the "modifier" in IntCriterionInput is redundant when i do this, meh
type IntCriterionStored struct {
	Modifier models.CriterionModifier `json:"modifier"`
	Value    models.IntCriterionInput `json:"value"`
}

func (m IntCriterionStored) ToOriginal() *models.IntCriterionInput {
	obj := m.Value
	obj.Modifier = m.Modifier
	return &obj
}

type DateCriterionStored struct {
	Modifier models.CriterionModifier  `json:"modifier"`
	Value    models.DateCriterionInput `json:"value"`
}

func (m DateCriterionStored) ToOriginal() *models.DateCriterionInput {
	obj := m.Value
	obj.Modifier = m.Modifier
	return &obj
}

type TimeCriterionStored struct {
	Modifier models.CriterionModifier       `json:"modifier"`
	Value    models.TimestampCriterionInput `json:"value"`
}

func (m TimeCriterionStored) ToOriginal() *models.TimestampCriterionInput {
	obj := m.Value
	obj.Modifier = m.Modifier
	return &obj
}

type HierarchicalMultiCriterionInputStoredEntry struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}
type HierarchicalMultiCriterionInputStoredAux struct {
	Items    *[]HierarchicalMultiCriterionInputStoredEntry `json:"items"`
	Depth    *int                                          `json:"depth"`
	Excludes *[]HierarchicalMultiCriterionInputStoredEntry `json:"excludes"`
}
type HierarchicalMultiCriterionInputStored struct {
	Modifier models.CriterionModifier                 `json:"modifier"`
	Value    HierarchicalMultiCriterionInputStoredAux `json:"value"`
}

func (m HierarchicalMultiCriterionInputStored) ToHierarchicalCriterion() *models.HierarchicalMultiCriterionInput {
	data := &models.HierarchicalMultiCriterionInput{
		Value:    []string{},
		Modifier: m.Modifier,
		Depth:    m.Value.Depth,
		Excludes: []string{},
	}

	if m.Value.Items != nil {
		for _, entry := range *m.Value.Items {
			data.Value = append(data.Value, entry.Id)
		}
	}
	if m.Value.Excludes != nil {
		for _, entry := range *m.Value.Excludes {
			data.Excludes = append(data.Excludes, entry.Id)
		}
	}

	return data
}
func (m HierarchicalMultiCriterionInputStored) ToMultiCriterion() *models.MultiCriterionInput {
	filled := m.ToHierarchicalCriterion()
	return &models.MultiCriterionInput{
		Value:    filled.Value,
		Modifier: m.Modifier,
		Excludes: filled.Excludes,
	}
}

type SceneFilterTypeStored struct {
	And      *SceneFilterTypeStored       `json:"AND"`
	Or       *SceneFilterTypeStored       `json:"OR"`
	Not      *SceneFilterTypeStored       `json:"NOT"`
	ID       *IntCriterionStored          `json:"id"`
	Title    *models.StringCriterionInput `json:"title"`
	Code     *models.StringCriterionInput `json:"code"`
	Details  *models.StringCriterionInput `json:"details"`
	Director *models.StringCriterionInput `json:"director"`
	// Filter by file oshash
	Oshash *models.StringCriterionInput `json:"oshash"`
	// Filter by file checksum
	Checksum *models.StringCriterionInput `json:"checksum"`
	// Filter by file phash
	Phash *models.StringCriterionInput `json:"phash"`
	// Filter by phash distance
	PhashDistance *models.PhashDistanceCriterionInput `json:"phash_distance"`
	// Filter by path
	Path *models.StringCriterionInput `json:"path"`
	// Filter by file count
	FileCount *IntCriterionStored `json:"file_count"`
	// Filter by rating expressed as 1-5
	Rating *IntCriterionStored `json:"rating"`
	// Filter by rating expressed as 1-100
	Rating100 *IntCriterionStored `json:"rating100"`
	// Filter by organized
	Organized *StringSingletonInput `json:"organized"`
	// Filter by o-counter
	OCounter *IntCriterionStored `json:"o_counter"`
	// Filter Scenes that have an exact phash match available
	Duplicated *models.PHashDuplicationCriterionInput `json:"duplicated"`
	// Filter by resolution
	Resolution *models.ResolutionCriterionInput `json:"resolution"`
	// Filter by video codec
	VideoCodec *models.StringCriterionInput `json:"video_codec"`
	// Filter by audio codec
	AudioCodec *models.StringCriterionInput `json:"audio_codec"`
	// Filter by duration (in seconds)
	Duration *IntCriterionStored `json:"duration"`
	// Filter to only include scenes which have markers. `true` or `false`
	HasMarkers *models.StringCriterionInput `json:"has_markers"`
	// Filter to only include scenes missing this property
	IsMissing *models.StringCriterionInput `json:"is_missing"`
	// Filter to only include scenes with this studio
	Studios *HierarchicalMultiCriterionInputStored `json:"studios"`
	// Filter to only include scenes with this movie
	Movies *HierarchicalMultiCriterionInputStored `json:"movies"`
	// Filter to only include scenes with these tags
	Tags *HierarchicalMultiCriterionInputStored `json:"tags"`
	// Filter by tag count
	TagCount *IntCriterionStored `json:"tag_count"`
	// Filter to only include scenes with performers with these tags
	PerformerTags *HierarchicalMultiCriterionInputStored `json:"performer_tags"`
	// Filter scenes that have performers that have been favorited
	PerformerFavorite *StringSingletonInput `json:"performer_favorite"`
	// Filter scenes by performer age at time of scene
	PerformerAge *IntCriterionStored `json:"performer_age"`
	// Filter to only include scenes with these performers
	Performers *HierarchicalMultiCriterionInputStored `json:"performers"`
	// Filter by performer count
	PerformerCount *IntCriterionStored `json:"performer_count"`
	// Filter by StashID
	StashID *models.StringCriterionInput `json:"stash_id"`
	// Filter by StashID Endpoint
	StashIDEndpoint *models.StashIDCriterionInput `json:"stash_id_endpoint"`
	// Filter by url
	URL *models.StringCriterionInput `json:"url"`
	// Filter by interactive
	Interactive *StringSingletonInput `json:"interactive"`
	// Filter by InteractiveSpeed
	InteractiveSpeed *IntCriterionStored `json:"interactive_speed"`
	// Filter by captions
	Captions *models.StringCriterionInput `json:"captions"`
	// Filter by resume time
	ResumeTime *IntCriterionStored `json:"resume_time"`
	// Filter by play count
	PlayCount *IntCriterionStored `json:"play_count"`
	// Filter by play duration (in seconds)
	PlayDuration *IntCriterionStored `json:"play_duration"`
	// Filter by date
	Date *DateCriterionStored `json:"date"`
	// Filter by created at
	CreatedAt *TimeCriterionStored `json:"created_at"`
	// Filter by updated at
	UpdatedAt *TimeCriterionStored `json:"updated_at"`
}

func (fsf SceneFilterTypeStored) ToOriginal() *models.SceneFilterType {
	model := models.SceneFilterType{
		Title:           fsf.Title,
		Code:            fsf.Code,
		Details:         fsf.Details,
		Director:        fsf.Director,
		Oshash:          fsf.Oshash,
		Checksum:        fsf.Checksum,
		Phash:           fsf.Phash,
		Path:            fsf.Path,
		Duplicated:      fsf.Duplicated,
		Resolution:      fsf.Resolution,
		VideoCodec:      fsf.VideoCodec,
		AudioCodec:      fsf.AudioCodec,
		StashID:         fsf.StashID,
		StashIDEndpoint: fsf.StashIDEndpoint,
		URL:             fsf.URL,
		Captions:        fsf.Captions,
	}

	if fsf.And != nil {
		model.And = fsf.And.ToOriginal()
	}
	if fsf.Or != nil {
		model.Or = fsf.Or.ToOriginal()
	}
	if fsf.Not != nil {
		model.Not = fsf.Not.ToOriginal()
	}
	if fsf.ID != nil {
		model.ID = fsf.ID.ToOriginal()
	}
	if fsf.FileCount != nil {
		model.FileCount = fsf.FileCount.ToOriginal()
	}
	if fsf.Rating100 != nil {
		model.Rating = fsf.Rating100.ToOriginal()
	}
	if fsf.Organized != nil {
		b := fsf.Organized.Value == "true"
		model.Organized = &b
	}
	if fsf.OCounter != nil {
		model.OCounter = fsf.OCounter.ToOriginal()
	}
	if fsf.Duration != nil {
		model.Duration = fsf.Duration.ToOriginal()
	}
	if fsf.HasMarkers != nil {
		model.HasMarkers = &fsf.HasMarkers.Value
	}
	if fsf.IsMissing != nil {
		model.IsMissing = &fsf.IsMissing.Value
	}
	if fsf.Studios != nil {
		model.Studios = fsf.Studios.ToHierarchicalCriterion()
	}
	if fsf.Movies != nil {
		model.Movies = fsf.Movies.ToMultiCriterion()
	}
	if fsf.Tags != nil {
		model.Tags = fsf.Tags.ToHierarchicalCriterion()
	}
	if fsf.TagCount != nil {
		model.TagCount = fsf.TagCount.ToOriginal()
	}
	if fsf.PerformerTags != nil {
		model.PerformerTags = fsf.PerformerTags.ToHierarchicalCriterion()
	}
	if fsf.PerformerFavorite != nil {
		b := fsf.PerformerFavorite.Value == "true"
		model.PerformerFavorite = &b
	}
	if fsf.PerformerAge != nil {
		model.PerformerAge = fsf.PerformerAge.ToOriginal()
	}
	if fsf.Performers != nil {
		model.Performers = fsf.Performers.ToMultiCriterion()
	}
	if fsf.PerformerCount != nil {
		model.PerformerCount = fsf.PerformerCount.ToOriginal()
	}
	if fsf.Interactive != nil {
		b := fsf.Interactive.Value == "true"
		model.Interactive = &b
	}
	if fsf.InteractiveSpeed != nil {
		model.InteractiveSpeed = fsf.InteractiveSpeed.ToOriginal()
	}
	if fsf.ResumeTime != nil {
		model.ResumeTime = fsf.ResumeTime.ToOriginal()
	}
	if fsf.PlayCount != nil {
		model.PlayCount = fsf.PlayCount.ToOriginal()
	}
	if fsf.PlayDuration != nil {
		model.PlayDuration = fsf.PlayDuration.ToOriginal()
	}
	if fsf.Date != nil {
		model.Date = fsf.Date.ToOriginal()
	}
	if fsf.CreatedAt != nil {
		model.CreatedAt = fsf.CreatedAt.ToOriginal()
	}
	if fsf.UpdatedAt != nil {
		model.UpdatedAt = fsf.UpdatedAt.ToOriginal()
	}

	return &model
}
