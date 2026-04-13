// TODO(audio): update this file

package api

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/plugin/hook"
	"github.com/stashapp/stash/pkg/audio"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/utils"
)

// used to refetch audio after hooks run
func (r *mutationResolver) getAudio(ctx context.Context, id int) (ret *models.Audio, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Audio.Find(ctx, id)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) AudioCreate(ctx context.Context, input models.AudioCreateInput) (ret *models.Audio, err error) {
	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	fileIDs, err := translator.fileIDSliceFromStringSlice(input.FileIds)
	if err != nil {
		return nil, fmt.Errorf("converting file ids: %w", err)
	}

	// Populate a new audio from the input
	newAudio := models.NewAudio()

	newAudio.Title = translator.string(input.Title)
	newAudio.Code = translator.string(input.Code)
	newAudio.Details = translator.string(input.Details)
	newAudio.Director = translator.string(input.Director)
	newAudio.Rating = input.Rating100
	newAudio.Organized = translator.bool(input.Organized)
	newAudio.StashIDs = models.NewRelatedStashIDs(models.StashIDInputs(input.StashIds).ToStashIDs())

	newAudio.Date, err = translator.datePtr(input.Date)
	if err != nil {
		return nil, fmt.Errorf("converting date: %w", err)
	}
	newAudio.StudioID, err = translator.intPtrFromString(input.StudioID)
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}

	if input.Urls != nil {
		newAudio.URLs = models.NewRelatedStrings(stringslice.TrimSpace(input.Urls))
	} else if input.URL != nil {
		newAudio.URLs = models.NewRelatedStrings([]string{strings.TrimSpace(*input.URL)})
	}

	newAudio.PerformerIDs, err = translator.relatedIds(input.PerformerIds)
	if err != nil {
		return nil, fmt.Errorf("converting performer ids: %w", err)
	}
	newAudio.TagIDs, err = translator.relatedIds(input.TagIds)
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}
	newAudio.GalleryIDs, err = translator.relatedIds(input.GalleryIds)
	if err != nil {
		return nil, fmt.Errorf("converting gallery ids: %w", err)
	}

	// prefer groups over movies
	if len(input.Groups) > 0 {
		newAudio.Groups, err = translator.relatedGroups(input.Groups)
		if err != nil {
			return nil, fmt.Errorf("converting groups: %w", err)
		}
	} else if len(input.Movies) > 0 {
		newAudio.Groups, err = translator.relatedGroupsFromMovies(input.Movies)
		if err != nil {
			return nil, fmt.Errorf("converting movies: %w", err)
		}
	}

	var coverImageData []byte
	if input.CoverImage != nil {
		var err error
		coverImageData, err = utils.ProcessImageInput(ctx, *input.CoverImage)
		if err != nil {
			return nil, fmt.Errorf("processing cover image: %w", err)
		}
	}

	customFields := convertMapJSONNumbers(input.CustomFields)

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.Resolver.audioService.Create(ctx, models.CreateAudioInput{
			Audio:        &newAudio,
			FileIDs:      fileIDs,
			CoverImage:   coverImageData,
			CustomFields: customFields,
		})
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) AudioUpdate(ctx context.Context, input models.AudioUpdateInput) (ret *models.Audio, err error) {
	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Start the transaction and save the audio
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.audioUpdate(ctx, input, translator)
		return err
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, ret.ID, hook.AudioUpdatePost, input, translator.getFields())
	return r.getAudio(ctx, ret.ID)
}

func (r *mutationResolver) AudiosUpdate(ctx context.Context, input []*models.AudioUpdateInput) (ret []*models.Audio, err error) {
	inputMaps := getUpdateInputMaps(ctx)

	// Start the transaction and save the audios
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		for i, audio := range input {
			translator := changesetTranslator{
				inputMap: inputMaps[i],
			}

			thisAudio, err := r.audioUpdate(ctx, *audio, translator)
			if err != nil {
				return err
			}

			ret = append(ret, thisAudio)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// execute post hooks outside of txn
	var newRet []*models.Audio
	for i, audio := range ret {
		translator := changesetTranslator{
			inputMap: inputMaps[i],
		}

		r.hookExecutor.ExecutePostHooks(ctx, audio.ID, hook.AudioUpdatePost, input, translator.getFields())

		audio, err = r.getAudio(ctx, audio.ID)
		if err != nil {
			return nil, err
		}

		newRet = append(newRet, audio)
	}

	return newRet, nil
}

func audioPartialFromInput(input models.AudioUpdateInput, translator changesetTranslator) (*models.AudioPartial, error) {
	updatedAudio := models.NewAudioPartial()

	updatedAudio.Title = translator.optionalString(input.Title, "title")
	updatedAudio.Code = translator.optionalString(input.Code, "code")
	updatedAudio.Details = translator.optionalString(input.Details, "details")
	updatedAudio.Director = translator.optionalString(input.Director, "director")
	updatedAudio.Rating = translator.optionalInt(input.Rating100, "rating100")

	if input.OCounter != nil {
		logger.Warnf("o_counter is deprecated and no longer supported, use audioIncrementO/audioDecrementO instead")
	}

	if input.PlayCount != nil {
		logger.Warnf("play_count is deprecated and no longer supported, use audioIncrementPlayCount/audioDecrementPlayCount instead")
	}

	updatedAudio.PlayDuration = translator.optionalFloat64(input.PlayDuration, "play_duration")
	updatedAudio.Organized = translator.optionalBool(input.Organized, "organized")
	updatedAudio.StashIDs = translator.updateStashIDs(input.StashIds, "stash_ids")

	var err error

	updatedAudio.Date, err = translator.optionalDate(input.Date, "date")
	if err != nil {
		return nil, fmt.Errorf("converting date: %w", err)
	}
	updatedAudio.StudioID, err = translator.optionalIntFromString(input.StudioID, "studio_id")
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}

	updatedAudio.URLs = translator.optionalURLs(input.Urls, input.URL)

	updatedAudio.PrimaryFileID, err = translator.fileIDPtrFromString(input.PrimaryFileID)
	if err != nil {
		return nil, fmt.Errorf("converting primary file id: %w", err)
	}

	updatedAudio.PerformerIDs, err = translator.updateIds(input.PerformerIds, "performer_ids")
	if err != nil {
		return nil, fmt.Errorf("converting performer ids: %w", err)
	}
	updatedAudio.TagIDs, err = translator.updateIds(input.TagIds, "tag_ids")
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}
	updatedAudio.GalleryIDs, err = translator.updateIds(input.GalleryIds, "gallery_ids")
	if err != nil {
		return nil, fmt.Errorf("converting gallery ids: %w", err)
	}

	if translator.hasField("groups") {
		updatedAudio.GroupIDs, err = translator.updateGroupIDs(input.Groups, "groups")
		if err != nil {
			return nil, fmt.Errorf("converting groups: %w", err)
		}
	} else if translator.hasField("movies") {
		updatedAudio.GroupIDs, err = translator.updateGroupIDsFromMovies(input.Movies, "movies")
		if err != nil {
			return nil, fmt.Errorf("converting movies: %w", err)
		}
	}

	return &updatedAudio, nil
}

func (r *mutationResolver) audioUpdate(ctx context.Context, input models.AudioUpdateInput, translator changesetTranslator) (*models.Audio, error) {
	audioID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, fmt.Errorf("converting id: %w", err)
	}

	qb := r.repository.Audio

	originalAudio, err := qb.Find(ctx, audioID)
	if err != nil {
		return nil, err
	}

	if originalAudio == nil {
		return nil, fmt.Errorf("audio with id %d not found", audioID)
	}

	// Populate audio from the input
	updatedAudio, err := audioPartialFromInput(input, translator)
	if err != nil {
		return nil, err
	}

	// ensure that title is set where audio has no file
	if updatedAudio.Title.Set && updatedAudio.Title.Value == "" {
		if err := originalAudio.LoadFiles(ctx, r.repository.Audio); err != nil {
			return nil, err
		}

		if len(originalAudio.Files.List()) == 0 {
			return nil, errors.New("title must be set if audio has no files")
		}
	}

	if updatedAudio.PrimaryFileID != nil {
		newPrimaryFileID := *updatedAudio.PrimaryFileID

		// if file hash has changed, we should migrate generated files
		// after commit
		if err := originalAudio.LoadFiles(ctx, r.repository.Audio); err != nil {
			return nil, err
		}

		// ensure that new primary file is associated with audio
		var f *models.VideoFile
		for _, ff := range originalAudio.Files.List() {
			if ff.ID == newPrimaryFileID {
				f = ff
			}
		}

		if f == nil {
			return nil, fmt.Errorf("file with id %d not associated with audio", newPrimaryFileID)
		}
	}

	var coverImageData []byte
	coverImageIncluded := translator.hasField("cover_image")
	if input.CoverImage != nil {
		var err error
		coverImageData, err = utils.ProcessImageInput(ctx, *input.CoverImage)
		if err != nil {
			return nil, fmt.Errorf("processing cover image: %w", err)
		}
	}

	var customFields *models.CustomFieldsInput
	if input.CustomFields != nil {
		cfCopy := *input.CustomFields
		customFields = &cfCopy
		// convert json.Numbers to int/float
		customFields.Full = convertMapJSONNumbers(customFields.Full)
		customFields.Partial = convertMapJSONNumbers(customFields.Partial)
	}

	audio, err := qb.UpdatePartial(ctx, audioID, *updatedAudio)
	if err != nil {
		return nil, err
	}

	if coverImageIncluded {
		if err := r.audioUpdateCoverImage(ctx, audio, coverImageData); err != nil {
			return nil, err
		}
	}

	if customFields != nil {
		if err := qb.SetCustomFields(ctx, audio.ID, *customFields); err != nil {
			return nil, err
		}
	}

	return audio, nil
}

func (r *mutationResolver) audioUpdateCoverImage(ctx context.Context, s *models.Audio, coverImageData []byte) error {
	qb := r.repository.Audio

	// update cover table - empty data will clear the cover
	if err := qb.UpdateCover(ctx, s.ID, coverImageData); err != nil {
		return err
	}

	return nil
}

func (r *mutationResolver) BulkAudioUpdate(ctx context.Context, input BulkAudioUpdateInput) ([]*models.Audio, error) {
	audioIDs, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return nil, fmt.Errorf("converting ids: %w", err)
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate audio from the input
	updatedAudio := models.NewAudioPartial()

	updatedAudio.Title = translator.optionalString(input.Title, "title")
	updatedAudio.Code = translator.optionalString(input.Code, "code")
	updatedAudio.Details = translator.optionalString(input.Details, "details")
	updatedAudio.Director = translator.optionalString(input.Director, "director")
	updatedAudio.Rating = translator.optionalInt(input.Rating100, "rating100")
	updatedAudio.Organized = translator.optionalBool(input.Organized, "organized")

	updatedAudio.Date, err = translator.optionalDate(input.Date, "date")
	if err != nil {
		return nil, fmt.Errorf("converting date: %w", err)
	}
	updatedAudio.StudioID, err = translator.optionalIntFromString(input.StudioID, "studio_id")
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}

	updatedAudio.URLs = translator.optionalURLsBulk(input.Urls, input.URL)

	updatedAudio.PerformerIDs, err = translator.updateIdsBulk(input.PerformerIds, "performer_ids")
	if err != nil {
		return nil, fmt.Errorf("converting performer ids: %w", err)
	}
	updatedAudio.TagIDs, err = translator.updateIdsBulk(input.TagIds, "tag_ids")
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}
	updatedAudio.GalleryIDs, err = translator.updateIdsBulk(input.GalleryIds, "gallery_ids")
	if err != nil {
		return nil, fmt.Errorf("converting gallery ids: %w", err)
	}

	if translator.hasField("group_ids") {
		updatedAudio.GroupIDs, err = translator.updateGroupIDsBulk(input.GroupIds, "group_ids")
		if err != nil {
			return nil, fmt.Errorf("converting group ids: %w", err)
		}
	} else if translator.hasField("movie_ids") {
		updatedAudio.GroupIDs, err = translator.updateGroupIDsBulk(input.MovieIds, "movie_ids")
		if err != nil {
			return nil, fmt.Errorf("converting movie ids: %w", err)
		}
	}

	var customFields *models.CustomFieldsInput
	if input.CustomFields != nil {
		cf := handleUpdateCustomFields(*input.CustomFields)
		customFields = &cf
	}

	ret := []*models.Audio{}

	// Start the transaction and save the audios
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Audio

		for _, audioID := range audioIDs {
			audio, err := qb.UpdatePartial(ctx, audioID, updatedAudio)
			if err != nil {
				return err
			}

			if customFields != nil {
				if err := qb.SetCustomFields(ctx, audio.ID, *customFields); err != nil {
					return err
				}
			}

			ret = append(ret, audio)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// execute post hooks outside of txn
	var newRet []*models.Audio
	for _, audio := range ret {
		r.hookExecutor.ExecutePostHooks(ctx, audio.ID, hook.AudioUpdatePost, input, translator.getFields())

		audio, err = r.getAudio(ctx, audio.ID)
		if err != nil {
			return nil, err
		}

		newRet = append(newRet, audio)
	}

	return newRet, nil
}

func (r *mutationResolver) AudioDestroy(ctx context.Context, input models.AudioDestroyInput) (bool, error) {
	audioID, err := strconv.Atoi(input.ID)
	if err != nil {
		return false, fmt.Errorf("converting id: %w", err)
	}

	fileNamingAlgo := manager.GetInstance().Config.GetVideoFileNamingAlgorithm()
	trashPath := manager.GetInstance().Config.GetDeleteTrashPath()

	var s *models.Audio
	fileDeleter := &audio.FileDeleter{
		Deleter:        file.NewDeleterWithTrash(trashPath),
		FileNamingAlgo: fileNamingAlgo,
		Paths:          manager.GetInstance().Paths,
	}

	deleteGenerated := utils.IsTrue(input.DeleteGenerated)
	deleteFile := utils.IsTrue(input.DeleteFile)
	destroyFileEntry := utils.IsTrue(input.DestroyFileEntry)

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Audio
		var err error
		s, err = qb.Find(ctx, audioID)
		if err != nil {
			return err
		}

		if s == nil {
			return fmt.Errorf("audio with id %d not found", audioID)
		}

		// kill any running encoders
		manager.KillRunningStreams(s, fileNamingAlgo)

		return r.audioService.Destroy(ctx, s, fileDeleter, deleteGenerated, deleteFile, destroyFileEntry)
	}); err != nil {
		fileDeleter.Rollback()
		return false, err
	}

	// perform the post-commit actions
	fileDeleter.Commit()

	// call post hook after performing the other actions
	r.hookExecutor.ExecutePostHooks(ctx, s.ID, hook.AudioDestroyPost, plugin.AudioDestroyInput{
		AudioDestroyInput: input,
		Checksum:          s.Checksum,
		OSHash:            s.OSHash,
		Path:              s.Path,
	}, nil)

	return true, nil
}

func (r *mutationResolver) AudiosDestroy(ctx context.Context, input models.AudiosDestroyInput) (bool, error) {
	audioIDs, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return false, fmt.Errorf("converting ids: %w", err)
	}

	var audios []*models.Audio
	fileNamingAlgo := manager.GetInstance().Config.GetVideoFileNamingAlgorithm()
	trashPath := manager.GetInstance().Config.GetDeleteTrashPath()

	fileDeleter := &audio.FileDeleter{
		Deleter:        file.NewDeleterWithTrash(trashPath),
		FileNamingAlgo: fileNamingAlgo,
		Paths:          manager.GetInstance().Paths,
	}

	deleteGenerated := utils.IsTrue(input.DeleteGenerated)
	deleteFile := utils.IsTrue(input.DeleteFile)
	destroyFileEntry := utils.IsTrue(input.DestroyFileEntry)

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Audio

		for _, id := range audioIDs {
			audio, err := qb.Find(ctx, id)
			if err != nil {
				return err
			}
			if audio == nil {
				return fmt.Errorf("audio with id %d not found", id)
			}

			audios = append(audios, audio)

			// kill any running encoders
			manager.KillRunningStreams(audio, fileNamingAlgo)

			if err := r.audioService.Destroy(ctx, audio, fileDeleter, deleteGenerated, deleteFile, destroyFileEntry); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		fileDeleter.Rollback()
		return false, err
	}

	// perform the post-commit actions
	fileDeleter.Commit()

	for _, audio := range audios {
		// call post hook after performing the other actions
		r.hookExecutor.ExecutePostHooks(ctx, audio.ID, hook.AudioDestroyPost, plugin.AudiosDestroyInput{
			AudiosDestroyInput: input,
			Checksum:           audio.Checksum,
			OSHash:             audio.OSHash,
			Path:               audio.Path,
		}, nil)
	}

	return true, nil
}

func (r *mutationResolver) AudioAssignFile(ctx context.Context, input AssignAudioFileInput) (bool, error) {
	audioID, err := strconv.Atoi(input.AudioID)
	if err != nil {
		return false, fmt.Errorf("converting audio id: %w", err)
	}

	fileID, err := strconv.Atoi(input.FileID)
	if err != nil {
		return false, fmt.Errorf("converting file id: %w", err)
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		return r.Resolver.audioService.AssignFile(ctx, audioID, models.FileID(fileID))
	}); err != nil {
		return false, fmt.Errorf("assigning file to audio: %w", err)
	}

	return true, nil
}

func (r *mutationResolver) AudioMerge(ctx context.Context, input AudioMergeInput) (*models.Audio, error) {
	srcIDs, err := stringslice.StringSliceToIntSlice(input.Source)
	if err != nil {
		return nil, fmt.Errorf("converting source ids: %w", err)
	}

	destID, err := strconv.Atoi(input.Destination)
	if err != nil {
		return nil, fmt.Errorf("converting destination id: %w", err)
	}

	var values *models.AudioPartial
	var coverImageData []byte
	var customFields *models.CustomFieldsInput

	if input.Values != nil {
		translator := changesetTranslator{
			inputMap: getNamedUpdateInputMap(ctx, "input.values"),
		}

		values, err = audioPartialFromInput(*input.Values, translator)
		if err != nil {
			return nil, err
		}

		if input.Values.CoverImage != nil {
			var err error
			coverImageData, err = utils.ProcessImageInput(ctx, *input.Values.CoverImage)
			if err != nil {
				return nil, fmt.Errorf("processing cover image: %w", err)
			}
		}

		if input.Values.CustomFields != nil {
			cf := handleUpdateCustomFields(*input.Values.CustomFields)
			customFields = &cf
		}
	} else {
		v := models.NewAudioPartial()
		values = &v
	}

	mgr := manager.GetInstance()
	trashPath := mgr.Config.GetDeleteTrashPath()
	fileDeleter := &audio.FileDeleter{
		Deleter:        file.NewDeleterWithTrash(trashPath),
		FileNamingAlgo: mgr.Config.GetVideoFileNamingAlgorithm(),
		Paths:          mgr.Paths,
	}

	var ret *models.Audio
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		if err := r.Resolver.audioService.Merge(ctx, srcIDs, destID, fileDeleter, audio.MergeOptions{
			AudioPartial:       *values,
			IncludePlayHistory: utils.IsTrue(input.PlayHistory),
			IncludeOHistory:    utils.IsTrue(input.OHistory),
		}); err != nil {
			return err
		}

		ret, err = r.Resolver.repository.Audio.Find(ctx, destID)
		if err != nil {
			return err
		}
		if ret == nil {
			return fmt.Errorf("audio with id %d not found", destID)
		}

		// only update cover image if one was provided
		if len(coverImageData) > 0 {
			if err := r.audioUpdateCoverImage(ctx, ret, coverImageData); err != nil {
				return err
			}
		}

		if customFields != nil {
			if err := r.Resolver.repository.Audio.SetCustomFields(ctx, ret.ID, *customFields); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) getAudioMarker(ctx context.Context, id int) (ret *models.AudioMarker, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.AudioMarker.Find(ctx, id)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) AudioMarkerCreate(ctx context.Context, input AudioMarkerCreateInput) (*models.AudioMarker, error) {
	audioID, err := strconv.Atoi(input.AudioID)
	if err != nil {
		return nil, fmt.Errorf("converting audio id: %w", err)
	}

	primaryTagID, err := strconv.Atoi(input.PrimaryTagID)
	if err != nil {
		return nil, fmt.Errorf("converting primary tag id: %w", err)
	}

	// Populate a new audio marker from the input
	newMarker := models.NewAudioMarker()

	newMarker.Title = strings.TrimSpace(input.Title)
	newMarker.Seconds = input.Seconds
	newMarker.PrimaryTagID = primaryTagID
	newMarker.AudioID = audioID

	if input.EndSeconds != nil {
		if err := validateAudioMarkerEndSeconds(newMarker.Seconds, *input.EndSeconds); err != nil {
			return nil, err
		}
		newMarker.EndSeconds = input.EndSeconds
	}

	tagIDs, err := stringslice.StringSliceToIntSlice(input.TagIds)
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.AudioMarker

		err := qb.Create(ctx, &newMarker)
		if err != nil {
			return err
		}

		// Save the marker tags
		// If this tag is the primary tag, then let's not add it.
		tagIDs = sliceutil.Exclude(tagIDs, []int{newMarker.PrimaryTagID})
		return qb.UpdateTags(ctx, newMarker.ID, tagIDs)
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, newMarker.ID, hook.AudioMarkerCreatePost, input, nil)
	return r.getAudioMarker(ctx, newMarker.ID)
}

func validateAudioMarkerEndSeconds(seconds, endSeconds float64) error {
	if endSeconds < seconds {
		return fmt.Errorf("end_seconds (%f) must be greater than or equal to seconds (%f)", endSeconds, seconds)
	}
	return nil
}

func float64OrZero(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func (r *mutationResolver) AudioMarkerUpdate(ctx context.Context, input AudioMarkerUpdateInput) (*models.AudioMarker, error) {
	markerID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, fmt.Errorf("converting id: %w", err)
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate audio marker from the input
	updatedMarker := models.NewAudioMarkerPartial()

	updatedMarker.Title = translator.optionalString(input.Title, "title")
	updatedMarker.Seconds = translator.optionalFloat64(input.Seconds, "seconds")
	updatedMarker.EndSeconds = translator.optionalFloat64(input.EndSeconds, "end_seconds")
	updatedMarker.AudioID, err = translator.optionalIntFromString(input.AudioID, "audio_id")
	if err != nil {
		return nil, fmt.Errorf("converting audio id: %w", err)
	}
	updatedMarker.PrimaryTagID, err = translator.optionalIntFromString(input.PrimaryTagID, "primary_tag_id")
	if err != nil {
		return nil, fmt.Errorf("converting primary tag id: %w", err)
	}

	var tagIDs []int
	tagIdsIncluded := translator.hasField("tag_ids")
	if input.TagIds != nil {
		tagIDs, err = stringslice.StringSliceToIntSlice(input.TagIds)
		if err != nil {
			return nil, fmt.Errorf("converting tag ids: %w", err)
		}
	}

	mgr := manager.GetInstance()
	trashPath := mgr.Config.GetDeleteTrashPath()

	fileDeleter := &audio.FileDeleter{
		Deleter:        file.NewDeleterWithTrash(trashPath),
		FileNamingAlgo: mgr.Config.GetVideoFileNamingAlgorithm(),
		Paths:          mgr.Paths,
	}

	// Start the transaction and save the audio marker
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.AudioMarker
		sqb := r.repository.Audio

		// check to see if timestamp was changed
		existingMarker, err := qb.Find(ctx, markerID)
		if err != nil {
			return err
		}
		if existingMarker == nil {
			return fmt.Errorf("audio marker with id %d not found", markerID)
		}

		// Validate end_seconds
		shouldValidateEndSeconds := (updatedMarker.Seconds.Set || updatedMarker.EndSeconds.Set) && !updatedMarker.EndSeconds.Null
		if shouldValidateEndSeconds {
			seconds := existingMarker.Seconds
			if updatedMarker.Seconds.Set {
				seconds = updatedMarker.Seconds.Value
			}

			endSeconds := existingMarker.EndSeconds
			if updatedMarker.EndSeconds.Set {
				endSeconds = &updatedMarker.EndSeconds.Value
			}

			if endSeconds != nil {
				if err := validateAudioMarkerEndSeconds(seconds, *endSeconds); err != nil {
					return err
				}
			}
		}

		newMarker, err := qb.UpdatePartial(ctx, markerID, updatedMarker)
		if err != nil {
			return err
		}

		existingAudio, err := sqb.Find(ctx, existingMarker.AudioID)
		if err != nil {
			return err
		}
		if existingAudio == nil {
			return fmt.Errorf("audio with id %d not found", existingMarker.AudioID)
		}

		// remove the marker preview if the audio changed or if the timestamp was changed
		if existingMarker.AudioID != newMarker.AudioID || existingMarker.Seconds != newMarker.Seconds || float64OrZero(existingMarker.EndSeconds) != float64OrZero(newMarker.EndSeconds) {
			seconds := int(existingMarker.Seconds)
			if err := fileDeleter.MarkMarkerFiles(existingAudio, seconds); err != nil {
				return err
			}
		}

		if tagIdsIncluded {
			// Save the marker tags
			// If this tag is the primary tag, then let's not add it.
			tagIDs = sliceutil.Exclude(tagIDs, []int{newMarker.PrimaryTagID})
			if err := qb.UpdateTags(ctx, markerID, tagIDs); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		fileDeleter.Rollback()
		return nil, err
	}

	// perform the post-commit actions
	fileDeleter.Commit()

	r.hookExecutor.ExecutePostHooks(ctx, markerID, hook.AudioMarkerUpdatePost, input, translator.getFields())
	return r.getAudioMarker(ctx, markerID)
}

func (r *mutationResolver) BulkAudioMarkerUpdate(ctx context.Context, input BulkAudioMarkerUpdateInput) ([]*models.AudioMarker, error) {
	ids, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return nil, fmt.Errorf("converting ids: %w", err)
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate performer from the input
	partial := models.NewAudioMarkerPartial()

	partial.Title = translator.optionalString(input.Title, "title")

	partial.PrimaryTagID, err = translator.optionalIntFromString(input.PrimaryTagID, "primary_tag_id")
	if err != nil {
		return nil, fmt.Errorf("converting primary tag id: %w", err)
	}

	partial.TagIDs, err = translator.updateIdsBulk(input.TagIds, "tag_ids")
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}

	ret := []*models.AudioMarker{}

	// Start the transaction and save the performers
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.AudioMarker

		for _, id := range ids {
			l := partial

			if err := adjustMarkerPartialForTagExclusion(ctx, r.repository.AudioMarker, id, &l); err != nil {
				return err
			}

			updated, err := qb.UpdatePartial(ctx, id, l)
			if err != nil {
				return err
			}

			ret = append(ret, updated)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// execute post hooks outside of txn
	var newRet []*models.AudioMarker
	for _, m := range ret {
		r.hookExecutor.ExecutePostHooks(ctx, m.ID, hook.AudioMarkerUpdatePost, input, translator.getFields())

		m, err = r.getAudioMarker(ctx, m.ID)
		if err != nil {
			return nil, err
		}

		newRet = append(newRet, m)
	}

	return newRet, nil
}

// adjustMarkerPartialForTagExclusion adjusts the AudioMarkerPartial to exclude the primary tag from tag updates.
func adjustMarkerPartialForTagExclusion(ctx context.Context, r models.AudioMarkerReader, id int, partial *models.AudioMarkerPartial) error {
	if partial.TagIDs == nil && !partial.PrimaryTagID.Set {
		return nil
	}

	// exclude primary tag from tag updates
	var primaryTagID int
	if partial.PrimaryTagID.Set {
		primaryTagID = partial.PrimaryTagID.Value
	} else {
		existing, err := r.Find(ctx, id)
		if err != nil {
			return fmt.Errorf("finding existing primary tag id: %w", err)
		}

		primaryTagID = existing.PrimaryTagID
	}

	existingTagIDs, err := r.GetTagIDs(ctx, id)
	if err != nil {
		return fmt.Errorf("getting existing tag ids: %w", err)
	}

	tagIDAttr := partial.TagIDs

	if tagIDAttr == nil {
		tagIDAttr = &models.UpdateIDs{
			IDs:  existingTagIDs,
			Mode: models.RelationshipUpdateModeSet,
		}
	}

	newTagIDs := tagIDAttr.Apply(existingTagIDs)
	// Remove primary tag from newTagIDs if present
	newTagIDs = sliceutil.Exclude(newTagIDs, []int{primaryTagID})

	if len(existingTagIDs) != len(newTagIDs) {
		partial.TagIDs = &models.UpdateIDs{
			IDs:  newTagIDs,
			Mode: models.RelationshipUpdateModeSet,
		}
	} else {
		// no change to tags required
		partial.TagIDs = nil
	}

	return nil
}

func (r *mutationResolver) AudioMarkerDestroy(ctx context.Context, id string) (bool, error) {
	return r.AudioMarkersDestroy(ctx, []string{id})
}

func (r *mutationResolver) AudioMarkersDestroy(ctx context.Context, markerIDs []string) (bool, error) {
	ids, err := stringslice.StringSliceToIntSlice(markerIDs)
	if err != nil {
		return false, fmt.Errorf("converting ids: %w", err)
	}

	var markers []*models.AudioMarker
	fileNamingAlgo := manager.GetInstance().Config.GetVideoFileNamingAlgorithm()
	trashPath := manager.GetInstance().Config.GetDeleteTrashPath()

	fileDeleter := &audio.FileDeleter{
		Deleter:        file.NewDeleterWithTrash(trashPath),
		FileNamingAlgo: fileNamingAlgo,
		Paths:          manager.GetInstance().Paths,
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.AudioMarker
		sqb := r.repository.Audio

		for _, markerID := range ids {
			marker, err := qb.Find(ctx, markerID)

			if err != nil {
				return err
			}

			if marker == nil {
				return fmt.Errorf("audio marker with id %d not found", markerID)
			}

			s, err := sqb.Find(ctx, marker.AudioID)

			if err != nil {
				return err
			}

			if s == nil {
				return fmt.Errorf("audio with id %d not found", marker.AudioID)
			}

			markers = append(markers, marker)

			if err := audio.DestroyMarker(ctx, s, marker, qb, fileDeleter); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		fileDeleter.Rollback()
		return false, err
	}

	fileDeleter.Commit()

	for _, marker := range markers {
		r.hookExecutor.ExecutePostHooks(ctx, marker.ID, hook.AudioMarkerDestroyPost, markerIDs, nil)
	}

	return true, nil
}

func (r *mutationResolver) AudioSaveActivity(ctx context.Context, id string, resumeTime *float64, playDuration *float64) (ret bool, err error) {
	audioID, err := strconv.Atoi(id)
	if err != nil {
		return false, fmt.Errorf("converting id: %w", err)
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Audio

		ret, err = qb.SaveActivity(ctx, audioID, resumeTime, playDuration)
		return err
	}); err != nil {
		return false, err
	}

	return ret, nil
}

func (r *mutationResolver) AudioResetActivity(ctx context.Context, id string, resetResume *bool, resetDuration *bool) (ret bool, err error) {
	audioID, err := strconv.Atoi(id)
	if err != nil {
		return false, fmt.Errorf("converting id: %w", err)
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Audio

		ret, err = qb.ResetActivity(ctx, audioID, utils.IsTrue(resetResume), utils.IsTrue(resetDuration))
		return err
	}); err != nil {
		return false, err
	}

	return ret, nil
}

// deprecated
func (r *mutationResolver) AudioIncrementPlayCount(ctx context.Context, id string) (ret int, err error) {
	audioID, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("converting id: %w", err)
	}

	var updatedTimes []time.Time

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Audio

		updatedTimes, err = qb.AddViews(ctx, audioID, nil)
		return err
	}); err != nil {
		return 0, err
	}

	return len(updatedTimes), nil
}

func (r *mutationResolver) AudioAddPlay(ctx context.Context, id string, t []*time.Time) (*HistoryMutationResult, error) {
	audioID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("converting id: %w", err)
	}

	var times []time.Time

	// convert time to local time, so that sorting is consistent
	for _, tt := range t {
		times = append(times, tt.Local())
	}

	var updatedTimes []time.Time

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Audio

		updatedTimes, err = qb.AddViews(ctx, audioID, times)
		return err
	}); err != nil {
		return nil, err
	}

	return &HistoryMutationResult{
		Count:   len(updatedTimes),
		History: sliceutil.ValuesToPtrs(updatedTimes),
	}, nil
}

func (r *mutationResolver) AudioDeletePlay(ctx context.Context, id string, t []*time.Time) (*HistoryMutationResult, error) {
	audioID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var times []time.Time

	for _, tt := range t {
		times = append(times, *tt)
	}

	var updatedTimes []time.Time

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Audio

		updatedTimes, err = qb.DeleteViews(ctx, audioID, times)
		return err
	}); err != nil {
		return nil, err
	}

	return &HistoryMutationResult{
		Count:   len(updatedTimes),
		History: sliceutil.ValuesToPtrs(updatedTimes),
	}, nil
}

func (r *mutationResolver) AudioResetPlayCount(ctx context.Context, id string) (ret int, err error) {
	audioID, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Audio

		ret, err = qb.DeleteAllViews(ctx, audioID)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

// deprecated
func (r *mutationResolver) AudioIncrementO(ctx context.Context, id string) (ret int, err error) {
	audioID, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("converting id: %w", err)
	}

	var updatedTimes []time.Time

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Audio

		updatedTimes, err = qb.AddO(ctx, audioID, nil)
		return err
	}); err != nil {
		return 0, err
	}

	return len(updatedTimes), nil
}

// deprecated
func (r *mutationResolver) AudioDecrementO(ctx context.Context, id string) (ret int, err error) {
	audioID, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("converting id: %w", err)
	}

	var updatedTimes []time.Time

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Audio

		updatedTimes, err = qb.DeleteO(ctx, audioID, nil)
		return err
	}); err != nil {
		return 0, err
	}

	return len(updatedTimes), nil
}

func (r *mutationResolver) AudioResetO(ctx context.Context, id string) (ret int, err error) {
	audioID, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("converting id: %w", err)
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Audio

		ret, err = qb.ResetO(ctx, audioID)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *mutationResolver) AudioAddO(ctx context.Context, id string, t []*time.Time) (*HistoryMutationResult, error) {
	audioID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("converting id: %w", err)
	}

	var times []time.Time

	// convert time to local time, so that sorting is consistent
	for _, tt := range t {
		times = append(times, tt.Local())
	}

	var updatedTimes []time.Time

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Audio

		updatedTimes, err = qb.AddO(ctx, audioID, times)
		return err
	}); err != nil {
		return nil, err
	}

	return &HistoryMutationResult{
		Count:   len(updatedTimes),
		History: sliceutil.ValuesToPtrs(updatedTimes),
	}, nil
}

func (r *mutationResolver) AudioDeleteO(ctx context.Context, id string, t []*time.Time) (*HistoryMutationResult, error) {
	audioID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("converting id: %w", err)
	}

	var times []time.Time

	for _, tt := range t {
		times = append(times, *tt)
	}

	var updatedTimes []time.Time

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Audio

		updatedTimes, err = qb.DeleteO(ctx, audioID, times)
		return err
	}); err != nil {
		return nil, err
	}

	return &HistoryMutationResult{
		Count:   len(updatedTimes),
		History: sliceutil.ValuesToPtrs(updatedTimes),
	}, nil
}

func (r *mutationResolver) AudioGenerateScreenshot(ctx context.Context, id string, at *float64) (string, error) {
	if at != nil {
		manager.GetInstance().GenerateScreenshot(ctx, id, *at)
	} else {
		manager.GetInstance().GenerateDefaultScreenshot(ctx, id)
	}

	return "todo", nil
}
