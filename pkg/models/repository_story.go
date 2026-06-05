package models

import "context"

type StoryGetter interface {
	Find(ctx context.Context, id int) (*Story, error)
	FindMany(ctx context.Context, ids []int) ([]*Story, error)
}

type StoryFinder interface {
	StoryGetter
	FindByPerformerID(ctx context.Context, performerID int) ([]*Story, error)
	FindByStudioID(ctx context.Context, studioID int) ([]*Story, error)
}

type StoryQueryer interface {
	Query(ctx context.Context, options StoryFilterType, findFilter *FindFilterType) ([]*Story, int, error)
}

type StoryCreator interface {
	Create(ctx context.Context, input *CreateStoryInput) error
}

type StoryUpdater interface {
	Update(ctx context.Context, updatedStory *UpdateStoryInput) error
	UpdatePartial(ctx context.Context, id int, updated StoryPartial) (*Story, error)
	UpdateFrontImage(ctx context.Context, storyID int, image []byte) error
	UpdateBackImage(ctx context.Context, storyID int, image []byte) error
}

type StoryDestroyer interface {
	Destroy(ctx context.Context, id int) error
}

type StoryReader interface {
	StoryFinder
	StoryQueryer
	All(ctx context.Context) ([]*Story, error)
	GetFrontImage(ctx context.Context, storyID int) ([]byte, error)
	GetBackImage(ctx context.Context, storyID int) ([]byte, error)
	HasFrontImage(ctx context.Context, storyID int) (bool, error)
	HasBackImage(ctx context.Context, storyID int) (bool, error)
	GetTagIDs(ctx context.Context, storyID int) ([]int, error)
	GetPerformerIDs(ctx context.Context, storyID int) ([]int, error)
	GetURLs(ctx context.Context, storyID int) ([]string, error)
}

type StoryWriter interface {
	StoryCreator
	StoryUpdater
	StoryDestroyer
}

type StoryReaderWriter interface {
	StoryReader
	StoryWriter
}
