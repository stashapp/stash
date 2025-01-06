package models

import (
	"context"
	"time"
)

type Studio struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	ParentID  *int      `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Rating expressed in 1-100 scale
	Rating        *int   `json:"rating"`
	Favorite      bool   `json:"favorite"`
	Details       string `json:"details"`
	IgnoreAutoTag bool   `json:"ignore_auto_tag"`

	Aliases  RelatedStrings  `json:"aliases"`
	TagIDs   RelatedIDs      `json:"tag_ids"`
	StashIDs RelatedStashIDs `json:"stash_ids"`
}

func NewStudio() Studio {
	currentTime := time.Now()
	return Studio{
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
}

// StudioPartial represents part of a Studio object. It is used to update the database entry.
type StudioPartial struct {
	ID       int
	Name     OptionalString
	URL      OptionalString
	ParentID OptionalInt
	// Rating expressed in 1-100 scale
	Rating        OptionalInt
	Favorite      OptionalBool
	Details       OptionalString
	CreatedAt     OptionalTime
	UpdatedAt     OptionalTime
	IgnoreAutoTag OptionalBool

	Aliases  *UpdateStrings
	TagIDs   *UpdateIDs
	StashIDs *UpdateStashIDs
}

func NewStudioPartial() StudioPartial {
	currentTime := time.Now()
	return StudioPartial{
		UpdatedAt: NewOptionalTime(currentTime),
	}
}

func (s *Studio) LoadAliases(ctx context.Context, l AliasLoader) error {
	return s.Aliases.load(func() ([]string, error) {
		return l.GetAliases(ctx, s.ID)
	})
}

func (s *Studio) LoadTagIDs(ctx context.Context, l TagIDLoader) error {
	return s.TagIDs.load(func() ([]int, error) {
		return l.GetTagIDs(ctx, s.ID)
	})
}

func (s *Studio) LoadStashIDs(ctx context.Context, l StashIDLoader) error {
	return s.StashIDs.load(func() ([]StashID, error) {
		return l.GetStashIDs(ctx, s.ID)
	})
}

func (s *Studio) LoadRelationships(ctx context.Context, l PerformerReader) error {
	if err := s.LoadAliases(ctx, l); err != nil {
		return err
	}

	if err := s.LoadTagIDs(ctx, l); err != nil {
		return err
	}

	if err := s.LoadStashIDs(ctx, l); err != nil {
		return err
	}

	return nil
}
