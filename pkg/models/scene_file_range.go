package models

import (
	"context"
)

// SceneFileRange is the time range within a file that a scene covers.
// Bounds are in seconds and are open-ended: a nil StartTime means 0 and a nil
// EndTime means the end of the file. Both nil therefore means the whole file,
// which is the pre-existing behavior for scenes that predate ranges.
type SceneFileRange struct {
	FileID    FileID
	StartTime *float64 // seconds from the start of the file; nil = 0
	EndTime   *float64 // seconds from the start of the file; nil = end of file
}

// Bounds returns the effective start and end for the given file duration,
// applying the open-ended defaults.
func (r SceneFileRange) Bounds(fileDuration float64) (start, end float64) {
	start = 0
	if r.StartTime != nil {
		start = *r.StartTime
	}
	end = fileDuration
	if r.EndTime != nil {
		end = *r.EndTime
	}
	return start, end
}

// Duration returns the effective duration of the ranged scene given the
// duration of its file.
func (r SceneFileRange) Duration(fileDuration float64) float64 {
	start, end := r.Bounds(fileDuration)
	return end - start
}

// SceneFileRangeLoader provides methods to load file ranges for scenes.
type SceneFileRangeLoader interface {
	GetFileRanges(ctx context.Context, relatedID int) ([]SceneFileRange, error)
}

// SceneFileSharingChecker provides methods to check how files are shared
// between scenes.
type SceneFileSharingChecker interface {
	// FileSharedWithUnrangedScene returns true if a scene other than
	// excludeSceneID is joined to the file without a range.
	FileSharedWithUnrangedScene(ctx context.Context, fileID FileID, excludeSceneID int) (bool, error)
}

// SceneFileRangeWriter provides methods to modify file ranges for scenes.
type SceneFileRangeWriter interface {
	// SetFileRange sets the range for the given file on the given scene.
	// Both bounds may be nil to clear the range.
	SetFileRange(ctx context.Context, sceneID int, fileID FileID, startTime, endTime *float64) error
	// AddFileWithRange joins the file to the scene with the given range
	// without removing existing joins for the file from other scenes.
	// The file is marked primary only if the scene has no primary file.
	AddFileWithRange(ctx context.Context, sceneID int, fileID FileID, startTime, endTime *float64) error
}
