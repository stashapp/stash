package models

import (
	"time"

	"github.com/stashapp/stash/pkg/hash/md5"
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
}

type StudioPartial struct {
	Checksum  OptionalString
	Name      OptionalString
	URL       OptionalString
	ParentID  OptionalInt
	CreatedAt OptionalTime
	UpdatedAt OptionalTime
	// Rating expressed in 1-100 scale
	Rating        OptionalInt
	Details       OptionalString
	IgnoreAutoTag OptionalBool
}

func NewStudio(name string) *Studio {
	currentTime := time.Now()
	return &Studio{
		Checksum:  md5.FromString(name),
		Name:      name,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
}

func NewStudioPartial() StudioPartial {
	updatedTime := time.Now()
	return StudioPartial{
		UpdatedAt: NewOptionalTime(updatedTime),
	}
}

type Studios []*Studio

func (s *Studios) Append(o interface{}) {
	*s = append(*s, o.(*Studio))
}

func (s *Studios) New() interface{} {
	return &Studio{}
}
