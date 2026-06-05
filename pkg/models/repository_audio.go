package models

import (
	"context"
)

// AudioGetter provides methods to get audios by ID.
type AudioGetter interface {
	// TODO - rename this to Find and remove existing method
	FindMany(ctx context.Context, ids []int) ([]*Audio, error)
	Find(ctx context.Context, id int) (*Audio, error)
	// FindByIDs works the same way as FindMany, but it ignores any audios not found
	// Audios are not guaranteed to be in the same order as the input
	FindByIDs(ctx context.Context, ids []int) ([]*Audio, error)
}

// AudioFinder provides methods to find audios.
type AudioFinder interface {
	AudioGetter
	IDsFromFileIDsLoader
	FindByFingerprints(ctx context.Context, fp []Fingerprint) ([]*Audio, error)
	FindByChecksum(ctx context.Context, checksum string) ([]*Audio, error)
	FindByPath(ctx context.Context, path string) ([]*Audio, error)
	FindByFileID(ctx context.Context, fileID FileID) ([]*Audio, error)
	FindByPrimaryFileID(ctx context.Context, fileID FileID) ([]*Audio, error)
	FindByPerformerID(ctx context.Context, performerID int) ([]*Audio, error)
	FindByGroupID(ctx context.Context, groupID int) ([]*Audio, error)
}

// AudioQueryer provides methods to query audios.
type AudioQueryer interface {
	Query(ctx context.Context, options AudioQueryOptions) (*AudioQueryResult, error)
	QueryCount(ctx context.Context, audioFilter *AudioFilterType, findFilter *FindFilterType) (int, error)
}

// AudioCounter provides methods to count audios.
type AudioCounter interface {
	Count(ctx context.Context) (int, error)
	CountByPerformerID(ctx context.Context, performerID int) (int, error)
	CountByFileID(ctx context.Context, fileID FileID) (int, error)
	OCountByPerformerID(ctx context.Context, performerID int) (int, error)
	OCountByGroupID(ctx context.Context, groupID int) (int, error)
	OCountByStudioID(ctx context.Context, studioID int) (int, error)
}

// AudioCreator provides methods to create audios.
type AudioCreator interface {
	Create(ctx context.Context, newAudio *Audio, fileIDs []FileID) error
}

// AudioUpdater provides methods to update audios.
type AudioUpdater interface {
	Update(ctx context.Context, updatedAudio *Audio) error
	UpdatePartial(ctx context.Context, id int, updatedAudio AudioPartial) (*Audio, error)
}

// AudioDestroyer provides methods to destroy audios.
type AudioDestroyer interface {
	Destroy(ctx context.Context, id int) error
}

type AudioCreatorUpdater interface {
	AudioCreator
	AudioUpdater
}

// AudioReader provides all methods to read audios.
type AudioReader interface {
	AudioFinder
	AudioQueryer
	AudioCounter

	URLLoader
	ViewDateReader
	ODateReader
	FileIDLoader
	PerformerIDLoader
	TagIDLoader
	AudioGroupLoader
	AudioFileLoader
	CustomFieldsReader

	All(ctx context.Context) ([]*Audio, error)
	Size(ctx context.Context) (float64, error)
	Duration(ctx context.Context) (float64, error)
	PlayDuration(ctx context.Context) (float64, error)
}

// AudioWriter provides all methods to modify audios.
type AudioWriter interface {
	AudioCreator
	AudioUpdater
	AudioDestroyer

	AddFileID(ctx context.Context, id int, fileID FileID) error
	AssignFiles(ctx context.Context, audioID int, fileID []FileID) error

	OHistoryWriter
	ViewHistoryWriter
	SaveActivity(ctx context.Context, audioID int, resumeTime *float64, playDuration *float64) (bool, error)
	ResetActivity(ctx context.Context, audioID int, resetResume bool, resetDuration bool) (bool, error)
	CustomFieldsWriter
}

// AudioReaderWriter provides all audio methods.
type AudioReaderWriter interface {
	AudioReader
	AudioWriter
}
