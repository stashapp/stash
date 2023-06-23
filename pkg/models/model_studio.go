package models

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Studio struct {
	ID        int       `json:"id"`
	Checksum  string    `json:"checksum"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	ParentID  *int      `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Rating expressed in 1-100 scale
	Rating        *int   `json:"rating"`
	Details       string `json:"details"`
	IgnoreAutoTag bool   `json:"ignore_auto_tag"`

	ImageBytes []byte

	Aliases  RelatedStrings  `json:"aliases"`
	StashIDs RelatedStashIDs `json:"stash_ids"`
}

func (s *Studio) LoadAliases(ctx context.Context, l AliasLoader) error {
	return s.Aliases.load(func() ([]string, error) {
		return l.GetAliases(ctx, s.ID)
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

	if err := s.LoadStashIDs(ctx, l); err != nil {
		return err
	}

	return nil
}

// StudioPartial represents part of a Studio object. It is used to update the database entry.
type StudioPartial struct {
	ID       int
	Checksum OptionalString
	Name     OptionalString
	URL      OptionalString
	ParentID OptionalInt
	// Rating expressed in 1-100 scale
	Rating        OptionalInt
	Details       OptionalString
	CreatedAt     OptionalTime
	UpdatedAt     OptionalTime
	IgnoreAutoTag OptionalBool

	// True if the image should be updated with ImageBytes
	ImageIncluded bool
	// Either contains the image, or is empty if the image should be removed
	ImageBytes []byte

	Aliases  *UpdateStrings
	StashIDs *UpdateStashIDs
}

type Studios []*Studio

func (s *Studios) Append(o interface{}) {
	*s = append(*s, o.(*Studio))
}

func (s *Studios) New() interface{} {
	return &Studio{}
}

// Checks to make sure that:
// 1. The studio exists locally
// 2. If the studio has a parent, it is not itself
// 3. If the studio has a parent, it exists locally and the parent does not have the studio as its parent
func (s *StudioPartial) ValidateModifyStudio(ctx context.Context, qb StudioReader) error {
	existing, err := qb.Find(ctx, s.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("studio with id %d not found", s.ID)
	}

	currentParentID := s.ParentID.Ptr()

	if currentParentID != nil {
		if *currentParentID == s.ID {
			return errors.New("studio cannot be an ancestor of itself")
		}

		// ensure there is no cyclic dependency
		parentStudio, err := qb.Find(ctx, *currentParentID)
		if err != nil || parentStudio == nil {
			return fmt.Errorf("error finding parent studio: %v", err)
		} else if parentStudio.ParentID == &s.ID {
			return errors.New("studio is already parent studio of the new parent studio")
		}
	}

	return nil
}
