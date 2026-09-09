package scene

import (
	"context"
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

var (
	ErrFileRangeNotFound  = errors.New("file not found on scene")
	ErrFileNotVideo       = errors.New("file is not a video file")
	ErrRangeSharedNoRange = errors.New("cannot clear file range: file is shared with another scene that has no range")
	errNoPrimaryFile      = errors.New("scene has no primary file")
)

// FileRangeInput is the service-level input for setting a scene file range.
// A nil FileID refers to the scene's primary file.
type FileRangeInput struct {
	FileID    *models.FileID
	StartTime *float64
	EndTime   *float64
}

func (i FileRangeInput) ResolveFileID(primaryFileID *models.FileID) (models.FileID, error) {
	if i.FileID != nil {
		return *i.FileID, nil
	}

	if primaryFileID == nil {
		return 0, errNoPrimaryFile
	}

	return *primaryFileID, nil
}

// validateFileRange validates the bounds of a single range against the
// duration of the given video file.
func validateFileRange(f *models.VideoFile, startTime, endTime *float64) error {
	if startTime != nil && *startTime < 0 {
		return fmt.Errorf("start time %v must not be negative", *startTime)
	}

	if endTime != nil && *endTime <= 0 {
		return fmt.Errorf("end time %v must be positive", *endTime)
	}

	if startTime != nil && endTime != nil && *startTime >= *endTime {
		return fmt.Errorf("start time %v must be before end time %v", *startTime, *endTime)
	}

	// tolerate a zero duration for files that have not been probed yet;
	// only reject bounds that are definitively beyond the file duration
	if f.Duration > 0 {
		if startTime != nil && *startTime >= f.Duration {
			return fmt.Errorf("start time %v must be before end of file (%v)", *startTime, f.Duration)
		}

		if endTime != nil && *endTime > f.Duration {
			return fmt.Errorf("end time %v must not exceed duration of file (%v)", *endTime, f.Duration)
		}
	}

	return nil
}

// validateFileRanges validates the input ranges against the given files.
// The scene's file IDs must be provided; a nil FileID in the input refers to
// the first (primary) file. Each file may be given a range at most once.
func validateFileRanges(fileIDs []models.FileID, primaryFileID *models.FileID, files []*models.VideoFile, input []FileRangeInput) (map[models.FileID]models.SceneFileRange, error) {
	filesByID := make(map[models.FileID]*models.VideoFile, len(files))
	for _, f := range files {
		filesByID[f.ID] = f
	}

	primary := primaryFileID
	if primary == nil && len(fileIDs) > 0 {
		p := fileIDs[0]
		primary = &p
	}

	ret := make(map[models.FileID]models.SceneFileRange)

	for _, i := range input {
		fileID, err := i.ResolveFileID(primary)
		if err != nil {
			return nil, err
		}

		if _, ok := ret[fileID]; ok {
			return nil, fmt.Errorf("duplicate range for file %d", fileID)
		}

		f, ok := filesByID[fileID]
		if !ok {
			return nil, fmt.Errorf("file %d: %w", fileID, ErrFileRangeNotFound)
		}

		if err := validateFileRange(f, i.StartTime, i.EndTime); err != nil {
			return nil, fmt.Errorf("file %s: %w", f.Path, err)
		}

		ret[fileID] = models.SceneFileRange{
			FileID:    fileID,
			StartTime: i.StartTime,
			EndTime:   i.EndTime,
		}
	}

	return ret, nil
}

// validateCreateFileRanges validates ranges provided to sceneCreate.
// All ranged files must be in the given file IDs.
func validateCreateFileRanges(fileIDs []models.FileID, files []*models.VideoFile, ranges []models.SceneFileRange) (map[models.FileID]models.SceneFileRange, error) {
	if len(ranges) == 0 {
		return nil, nil
	}

	if len(fileIDs) == 0 {
		return nil, errors.New("file ranges require file_ids to be set")
	}

	filesByID := make(map[models.FileID]*models.VideoFile, len(files))
	for _, f := range files {
		if f == nil {
			return nil, errors.New("file not found")
		}
		filesByID[f.ID] = f
	}

	ret := make(map[models.FileID]models.SceneFileRange)

	for _, r := range ranges {
		if _, ok := ret[r.FileID]; ok {
			return nil, fmt.Errorf("duplicate range for file %d", r.FileID)
		}

		f, ok := filesByID[r.FileID]
		if !ok {
			return nil, fmt.Errorf("file %d: %w", r.FileID, ErrFileRangeNotFound)
		}

		if err := validateFileRange(f, r.StartTime, r.EndTime); err != nil {
			return nil, fmt.Errorf("file %s: %w", f.Path, err)
		}

		ret[r.FileID] = r
	}

	return ret, nil
}

// assignFilesWithRanges assigns files to a newly created scene.
// Files with a range are joined without removing existing joins for the file
// from other scenes, so that multiple scenes can share one file.
// Files without a range use the existing exclusive assignment semantics.
func (s *Service) assignFilesWithRanges(ctx context.Context, sceneID int, fileIDs []models.FileID, ranges []models.SceneFileRange) error {
	if len(fileIDs) == 0 {
		return errors.New("file ranges require file_ids to be set")
	}

	files, err := s.File.Find(ctx, fileIDs...)
	if err != nil {
		return err
	}

	if len(files) != len(fileIDs) {
		return errors.New("file not found")
	}

	videoFiles := make([]*models.VideoFile, len(files))
	for i, f := range files {
		vf, ok := f.(*models.VideoFile)
		if !ok {
			return fmt.Errorf("%s: %w", f.Base().Path, ErrFileNotVideo)
		}
		videoFiles[i] = vf
	}

	rangesByFile, err := validateCreateFileRanges(fileIDs, videoFiles, ranges)
	if err != nil {
		return err
	}

	for _, vf := range videoFiles {
		if r, ok := rangesByFile[vf.ID]; ok {
			if err := s.Repository.AddFileWithRange(ctx, sceneID, vf.ID, r.StartTime, r.EndTime); err != nil {
				return fmt.Errorf("adding file %d with range to scene: %w", vf.ID, err)
			}
			continue
		}

		if err := s.AssignFile(ctx, sceneID, vf.ID); err != nil {
			return fmt.Errorf("assigning file %d to new scene: %w", vf.ID, err)
		}
	}

	return nil
}

// UpdateFileRanges replaces the file ranges for the scene.
// Ranges may only be set for files that are already joined to the scene.
// A range is cleared by omitting the file from the input; clearing is rejected
// if another scene without a range is joined to the same file.
func (s *Service) UpdateFileRanges(ctx context.Context, sceneID int, input []FileRangeInput) error {
	s2, err := s.Repository.Find(ctx, sceneID)
	if err != nil {
		return err
	}

	if err := s2.LoadFiles(ctx, s.Repository); err != nil {
		return err
	}

	if err := s2.LoadFileRanges(ctx, s.Repository); err != nil {
		return err
	}

	files := s2.Files.List()
	fileIDs := make([]models.FileID, len(files))
	for i, f := range files {
		fileIDs[i] = f.ID
	}

	newRanges, err := validateFileRanges(fileIDs, s2.PrimaryFileID, files, input)
	if err != nil {
		return err
	}

	oldRanges := make(map[models.FileID]models.SceneFileRange)
	for _, r := range s2.FileRanges.List() {
		oldRanges[r.FileID] = r
	}

	// apply new/changed ranges
	for fileID, r := range newRanges {
		old, found := oldRanges[fileID]
		if found && floatPtrEq(old.StartTime, r.StartTime) && floatPtrEq(old.EndTime, r.EndTime) {
			continue
		}

		if err := s.Repository.SetFileRange(ctx, sceneID, fileID, r.StartTime, r.EndTime); err != nil {
			return fmt.Errorf("setting file range: %w", err)
		}
	}

	// clear removed ranges
	for fileID := range oldRanges {
		if _, ok := newRanges[fileID]; ok {
			continue
		}

		// only allow clearing when no other scene holds the file without a range
		shared, err := s.Repository.FileSharedWithUnrangedScene(ctx, fileID, sceneID)
		if err != nil {
			return err
		}

		if shared {
			return fmt.Errorf("file %d: %w", fileID, ErrRangeSharedNoRange)
		}

		if err := s.Repository.SetFileRange(ctx, sceneID, fileID, nil, nil); err != nil {
			return fmt.Errorf("clearing file range: %w", err)
		}
	}

	return nil
}

func floatPtrEq(a, b *float64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
