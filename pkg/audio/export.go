// TODO(audio): update this file

package audio

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/utils"
)

type ExportGetter interface {
	models.ViewDateReader
	models.ODateReader
	models.CustomFieldsReader
	GetCover(ctx context.Context, audioID int) ([]byte, error)
}

type TagFinder interface {
	models.TagGetter
	FindByAudioID(ctx context.Context, audioID int) ([]*models.Tag, error)
	FindByAudioMarkerID(ctx context.Context, audioMarkerID int) ([]*models.Tag, error)
}

// ToBasicJSON converts a audio object into its JSON object equivalent. It
// does not convert the relationships to other objects, with the exception
// of cover image.
func ToBasicJSON(ctx context.Context, reader ExportGetter, audio *models.Audio) (*jsonschema.Audio, error) {
	newAudioJSON := jsonschema.Audio{
		Title:     audio.Title,
		Code:      audio.Code,
		URLs:      audio.URLs.List(),
		Details:   audio.Details,
		Director:  audio.Director,
		CreatedAt: json.JSONTime{Time: audio.CreatedAt},
		UpdatedAt: json.JSONTime{Time: audio.UpdatedAt},
	}

	if audio.Date != nil {
		newAudioJSON.Date = audio.Date.String()
	}

	if audio.Rating != nil {
		newAudioJSON.Rating = *audio.Rating
	}

	newAudioJSON.Organized = audio.Organized

	for _, f := range audio.Files.List() {
		newAudioJSON.Files = append(newAudioJSON.Files, f.Base().Path)
	}

	cover, err := reader.GetCover(ctx, audio.ID)
	if err != nil {
		logger.Errorf("Error getting audio cover: %v", err)
	}

	if len(cover) > 0 {
		newAudioJSON.Cover = utils.GetBase64StringFromData(cover)
	}

	var ret []models.StashID
	for _, stashID := range audio.StashIDs.List() {
		newJoin := models.StashID{
			StashID:  stashID.StashID,
			Endpoint: stashID.Endpoint,
		}
		ret = append(ret, newJoin)
	}

	newAudioJSON.StashIDs = ret

	dates, err := reader.GetViewDates(ctx, audio.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting view dates: %v", err)
	}

	for _, date := range dates {
		newAudioJSON.PlayHistory = append(newAudioJSON.PlayHistory, json.JSONTime{Time: date})
	}

	odates, err := reader.GetODates(ctx, audio.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting o dates: %v", err)
	}

	for _, date := range odates {
		newAudioJSON.OHistory = append(newAudioJSON.OHistory, json.JSONTime{Time: date})
	}

	newAudioJSON.CustomFields, err = reader.GetCustomFields(ctx, audio.ID)
	if err != nil {
		return nil, fmt.Errorf("getting audio custom fields: %v", err)
	}

	return &newAudioJSON, nil
}

// GetStudioName returns the name of the provided audio's studio. It returns an
// empty string if there is no studio assigned to the audio.
func GetStudioName(ctx context.Context, reader models.StudioGetter, audio *models.Audio) (string, error) {
	if audio.StudioID != nil {
		studio, err := reader.Find(ctx, *audio.StudioID)
		if err != nil {
			return "", err
		}

		if studio != nil {
			return studio.Name, nil
		}
	}

	return "", nil
}

// GetTagNames returns a slice of tag names corresponding to the provided
// audio's tags.
func GetTagNames(ctx context.Context, reader TagFinder, audio *models.Audio) ([]string, error) {
	tags, err := reader.FindByAudioID(ctx, audio.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting audio tags: %v", err)
	}

	return getTagNames(tags), nil
}

func getTagNames(tags []*models.Tag) []string {
	var results []string
	for _, tag := range tags {
		if tag.Name != "" {
			results = append(results, tag.Name)
		}
	}

	return results
}

// GetDependentTagIDs returns a slice of unique tag IDs that this audio references.
func GetDependentTagIDs(ctx context.Context, tags TagFinder, markerReader models.AudioMarkerFinder, audio *models.Audio) ([]int, error) {
	var ret []int

	t, err := tags.FindByAudioID(ctx, audio.ID)
	if err != nil {
		return nil, err
	}

	for _, tt := range t {
		ret = sliceutil.AppendUnique(ret, tt.ID)
	}

	sm, err := markerReader.FindByAudioID(ctx, audio.ID)
	if err != nil {
		return nil, err
	}

	for _, smm := range sm {
		ret = sliceutil.AppendUnique(ret, smm.PrimaryTagID)
		smmt, err := tags.FindByAudioMarkerID(ctx, smm.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid tags for audio marker: %v", err)
		}

		for _, smmtt := range smmt {
			ret = sliceutil.AppendUnique(ret, smmtt.ID)
		}
	}

	return ret, nil
}

// GetAudioGroupsJSON returns a slice of AudioGroup JSON representation objects
// corresponding to the provided audio's audio group relationships.
func GetAudioGroupsJSON(ctx context.Context, groupReader models.GroupGetter, audio *models.Audio) ([]jsonschema.AudioGroup, error) {
	audioGroups := audio.Groups.List()

	var results []jsonschema.AudioGroup
	for _, audioGroup := range audioGroups {
		group, err := groupReader.Find(ctx, audioGroup.GroupID)
		if err != nil {
			return nil, fmt.Errorf("error getting group: %v", err)
		}

		if group != nil {
			audioGroupJSON := jsonschema.AudioGroup{
				GroupName: group.Name,
			}
			if audioGroup.AudioIndex != nil {
				audioGroupJSON.AudioIndex = *audioGroup.AudioIndex
			}
			results = append(results, audioGroupJSON)
		}
	}

	return results, nil
}

// GetDependentGroupIDs returns a slice of group IDs that this audio references.
func GetDependentGroupIDs(ctx context.Context, audio *models.Audio) ([]int, error) {
	var ret []int

	m := audio.Groups.List()
	for _, mm := range m {
		ret = append(ret, mm.GroupID)
	}

	return ret, nil
}

// GetAudioMarkersJSON returns a slice of AudioMarker JSON representation
// objects corresponding to the provided audio's markers.
func GetAudioMarkersJSON(ctx context.Context, markerReader models.AudioMarkerFinder, tagReader TagFinder, audio *models.Audio) ([]jsonschema.AudioMarker, error) {
	audioMarkers, err := markerReader.FindByAudioID(ctx, audio.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting audio markers: %v", err)
	}

	var results []jsonschema.AudioMarker

	for _, audioMarker := range audioMarkers {
		primaryTag, err := tagReader.Find(ctx, audioMarker.PrimaryTagID)
		if err != nil {
			return nil, fmt.Errorf("invalid primary tag for audio marker: %v", err)
		}

		audioMarkerTags, err := tagReader.FindByAudioMarkerID(ctx, audioMarker.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid tags for audio marker: %v", err)
		}

		audioMarkerJSON := jsonschema.AudioMarker{
			Title:      audioMarker.Title,
			Seconds:    getDecimalString(audioMarker.Seconds),
			PrimaryTag: primaryTag.Name,
			Tags:       getTagNames(audioMarkerTags),
			CreatedAt:  json.JSONTime{Time: audioMarker.CreatedAt},
			UpdatedAt:  json.JSONTime{Time: audioMarker.UpdatedAt},
		}

		if audioMarker.EndSeconds != nil {
			audioMarkerJSON.EndSeconds = getDecimalString(*audioMarker.EndSeconds)
		}

		results = append(results, audioMarkerJSON)
	}

	return results, nil
}

func getDecimalString(num float64) string {
	if num == 0 {
		return ""
	}

	precision := getPrecision(num)
	if precision == 0 {
		precision = 1
	}
	return fmt.Sprintf("%."+strconv.Itoa(precision)+"f", num)
}

func getPrecision(num float64) int {
	if num == 0 {
		return 0
	}

	e := 1.0
	p := 0
	for (math.Round(num*e) / e) != num {
		e *= 10
		p++
	}
	return p
}
