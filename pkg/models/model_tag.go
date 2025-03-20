package models

import (
	"context"
	"strings"
	"time"
)

type Tag struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	SortName      string    `json:"sort_name"`
	Favorite      bool      `json:"favorite"`
	Description   string    `json:"description"`
	IgnoreAutoTag bool      `json:"ignore_auto_tag"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Aliases   RelatedStrings `json:"aliases"`
	ParentIDs RelatedIDs     `json:"parent_ids"`
	ChildIDs  RelatedIDs     `json:"tag_ids"`
}

func NewTag() Tag {
	currentTime := time.Now()
	return Tag{
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
}

func (s *Tag) LoadAliases(ctx context.Context, l AliasLoader) error {
	return s.Aliases.load(func() ([]string, error) {
		return l.GetAliases(ctx, s.ID)
	})
}

func (s *Tag) LoadParentIDs(ctx context.Context, l TagRelationLoader) error {
	return s.ParentIDs.load(func() ([]int, error) {
		return l.GetParentIDs(ctx, s.ID)
	})
}

func (s *Tag) LoadChildIDs(ctx context.Context, l TagRelationLoader) error {
	return s.ChildIDs.load(func() ([]int, error) {
		return l.GetChildIDs(ctx, s.ID)
	})
}

func TagSortFunction(a, b *Tag) int {
	sortA := a.SortName
	if sortA == "" {
		sortA = strings.ToLower(a.Name)
	}

	sortB := b.SortName
	if sortB == "" {
		sortB = strings.ToLower(b.Name)
	}

	return strings.Compare(sortA, sortB)
}

type TagPartial struct {
	Name          OptionalString
	SortName      OptionalString
	Description   OptionalString
	Favorite      OptionalBool
	IgnoreAutoTag OptionalBool
	CreatedAt     OptionalTime
	UpdatedAt     OptionalTime

	Aliases   *UpdateStrings
	ParentIDs *UpdateIDs
	ChildIDs  *UpdateIDs
}

func NewTagPartial() TagPartial {
	currentTime := time.Now()
	return TagPartial{
		UpdatedAt: NewOptionalTime(currentTime),
	}
}

type TagPath struct {
	Tag
	Path string `json:"path"`
}
