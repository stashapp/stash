package api

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/audio"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/plugin/hook"
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
	newAudio.Rating = input.Rating100
	newAudio.Organized = translator.bool(input.Organized)

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

	if len(input.Groups) > 0 {
		newAudio.Groups, err = translator.relatedGroupsAudio(input.Groups)
		if err != nil {
			return nil, fmt.Errorf("converting groups: %w", err)
		}
	}

	customFields := convertMapJSONNumbers(input.CustomFields)

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.Resolver.audioService.Create(ctx, models.CreateAudioInput{
			Audio:        &newAudio,
			FileIDs:      fileIDs,
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
	updatedAudio.Rating = translator.optionalInt(input.Rating100, "rating100")

	if input.OCounter != nil {
		logger.Warnf("o_counter is deprecated and no longer supported, use audioIncrementO/audioDecrementO instead")
	}

	if input.PlayCount != nil {
		logger.Warnf("play_count is deprecated and no longer supported, use audioIncrementPlayCount/audioDecrementPlayCount instead")
	}

	updatedAudio.PlayDuration = translator.optionalFloat64(input.PlayDuration, "play_duration")
	updatedAudio.Organized = translator.optionalBool(input.Organized, "organized")

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

	if translator.hasField("groups") {
		updatedAudio.GroupIDs, err = translator.updateGroupIDsAudio(input.Groups, "groups")
		if err != nil {
			return nil, fmt.Errorf("converting groups: %w", err)
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
		var f *models.AudioFile
		for _, ff := range originalAudio.Files.List() {
			if ff.ID == newPrimaryFileID {
				f = ff
			}
		}

		if f == nil {
			return nil, fmt.Errorf("file with id %d not associated with audio", newPrimaryFileID)
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

	if customFields != nil {
		if err := qb.SetCustomFields(ctx, audio.ID, *customFields); err != nil {
			return nil, err
		}
	}

	return audio, nil
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

	if translator.hasField("group_ids") {
		updatedAudio.GroupIDs, err = translator.updateGroupIDsBulkAudio(input.GroupIds, "group_ids")
		if err != nil {
			return nil, fmt.Errorf("converting group ids: %w", err)
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

	fileNamingAlgo := manager.GetInstance().Config.GetAudioFileNamingAlgorithm()
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
		manager.KillRunningStreamsAudio(s, fileNamingAlgo)

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
	fileNamingAlgo := manager.GetInstance().Config.GetAudioFileNamingAlgorithm()
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
			manager.KillRunningStreamsAudio(audio, fileNamingAlgo)

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
	var customFields *models.CustomFieldsInput

	if input.Values != nil {
		translator := changesetTranslator{
			inputMap: getNamedUpdateInputMap(ctx, "input.values"),
		}

		values, err = audioPartialFromInput(*input.Values, translator)
		if err != nil {
			return nil, err
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
		FileNamingAlgo: mgr.Config.GetAudioFileNamingAlgorithm(),
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
