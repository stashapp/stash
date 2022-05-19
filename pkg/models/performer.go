package models

import (
	"context"
	"fmt"
	"io"
	"strconv"
)

type GenderEnum string

const (
	GenderEnumMale              GenderEnum = "MALE"
	GenderEnumFemale            GenderEnum = "FEMALE"
	GenderEnumTransgenderMale   GenderEnum = "TRANSGENDER_MALE"
	GenderEnumTransgenderFemale GenderEnum = "TRANSGENDER_FEMALE"
	GenderEnumIntersex          GenderEnum = "INTERSEX"
	GenderEnumNonBinary         GenderEnum = "NON_BINARY"
)

var AllGenderEnum = []GenderEnum{
	GenderEnumMale,
	GenderEnumFemale,
	GenderEnumTransgenderMale,
	GenderEnumTransgenderFemale,
	GenderEnumIntersex,
	GenderEnumNonBinary,
}

func (e GenderEnum) IsValid() bool {
	switch e {
	case GenderEnumMale, GenderEnumFemale, GenderEnumTransgenderMale, GenderEnumTransgenderFemale, GenderEnumIntersex, GenderEnumNonBinary:
		return true
	}
	return false
}

func (e GenderEnum) String() string {
	return string(e)
}

func (e *GenderEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = GenderEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid GenderEnum", str)
	}
	return nil
}

func (e GenderEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type GenderCriterionInput struct {
	Value    *GenderEnum       `json:"value"`
	Modifier CriterionModifier `json:"modifier"`
}

type PerformerFilterType struct {
	And     *PerformerFilterType  `json:"AND"`
	Or      *PerformerFilterType  `json:"OR"`
	Not     *PerformerFilterType  `json:"NOT"`
	Name    *StringCriterionInput `json:"name"`
	Details *StringCriterionInput `json:"details"`
	// Filter by favorite
	FilterFavorites *bool `json:"filter_favorites"`
	// Filter by birth year
	BirthYear *IntCriterionInput `json:"birth_year"`
	// Filter by age
	Age *IntCriterionInput `json:"age"`
	// Filter by ethnicity
	Ethnicity *StringCriterionInput `json:"ethnicity"`
	// Filter by country
	Country *StringCriterionInput `json:"country"`
	// Filter by eye color
	EyeColor *StringCriterionInput `json:"eye_color"`
	// Filter by height
	Height *StringCriterionInput `json:"height"`
	// Filter by measurements
	Measurements *StringCriterionInput `json:"measurements"`
	// Filter by fake tits value
	FakeTits *StringCriterionInput `json:"fake_tits"`
	// Filter by career length
	CareerLength *StringCriterionInput `json:"career_length"`
	// Filter by tattoos
	Tattoos *StringCriterionInput `json:"tattoos"`
	// Filter by piercings
	Piercings *StringCriterionInput `json:"piercings"`
	// Filter by aliases
	Aliases *StringCriterionInput `json:"aliases"`
	// Filter by gender
	Gender *GenderCriterionInput `json:"gender"`
	// Filter to only include performers missing this property
	IsMissing *string `json:"is_missing"`
	// Filter to only include performers with these tags
	Tags *HierarchicalMultiCriterionInput `json:"tags"`
	// Filter by tag count
	TagCount *IntCriterionInput `json:"tag_count"`
	// Filter by scene count
	SceneCount *IntCriterionInput `json:"scene_count"`
	// Filter by image count
	ImageCount *IntCriterionInput `json:"image_count"`
	// Filter by gallery count
	GalleryCount *IntCriterionInput `json:"gallery_count"`
	// Filter by StashID
	StashID *StringCriterionInput `json:"stash_id"`
	// Filter by rating
	Rating *IntCriterionInput `json:"rating"`
	// Filter by url
	URL *StringCriterionInput `json:"url"`
	// Filter by hair color
	HairColor *StringCriterionInput `json:"hair_color"`
	// Filter by weight
	Weight *IntCriterionInput `json:"weight"`
	// Filter by death year
	DeathYear *IntCriterionInput `json:"death_year"`
	// Filter by studios where performer appears in scene/image/gallery
	Studios *HierarchicalMultiCriterionInput `json:"studios"`
	// Filter by autotag ignore value
	IgnoreAutoTag *bool `json:"ignore_auto_tag"`
}

type PerformerReader interface {
	Find(ctx context.Context, id int) (*Performer, error)
	FindMany(ctx context.Context, ids []int) ([]*Performer, error)
	FindBySceneID(ctx context.Context, sceneID int) ([]*Performer, error)
	FindNamesBySceneID(ctx context.Context, sceneID int) ([]*Performer, error)
	FindByImageID(ctx context.Context, imageID int) ([]*Performer, error)
	FindByGalleryID(ctx context.Context, galleryID int) ([]*Performer, error)
	FindByNames(ctx context.Context, names []string, nocase bool) ([]*Performer, error)
	FindByStashID(ctx context.Context, stashID StashID) ([]*Performer, error)
	FindByStashIDStatus(ctx context.Context, hasStashID bool, stashboxEndpoint string) ([]*Performer, error)
	CountByTagID(ctx context.Context, tagID int) (int, error)
	Count(ctx context.Context) (int, error)
	All(ctx context.Context) ([]*Performer, error)
	// TODO - this interface is temporary until the filter schema can fully
	// support the query needed
	QueryForAutoTag(ctx context.Context, words []string) ([]*Performer, error)
	Query(ctx context.Context, performerFilter *PerformerFilterType, findFilter *FindFilterType) ([]*Performer, int, error)
	GetImage(ctx context.Context, performerID int) ([]byte, error)
	GetStashIDs(ctx context.Context, performerID int) ([]*StashID, error)
	GetTagIDs(ctx context.Context, performerID int) ([]int, error)
}

type PerformerWriter interface {
	Create(ctx context.Context, newPerformer Performer) (*Performer, error)
	Update(ctx context.Context, updatedPerformer PerformerPartial) (*Performer, error)
	UpdateFull(ctx context.Context, updatedPerformer Performer) (*Performer, error)
	Destroy(ctx context.Context, id int) error
	UpdateImage(ctx context.Context, performerID int, image []byte) error
	DestroyImage(ctx context.Context, performerID int) error
	UpdateStashIDs(ctx context.Context, performerID int, stashIDs []StashID) error
	UpdateTags(ctx context.Context, performerID int, tagIDs []int) error
}

type PerformerReaderWriter interface {
	PerformerReader
	PerformerWriter
}
