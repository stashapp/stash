package models

import (
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
	Find(id int) (*Performer, error)
	FindMany(ids []int) ([]*Performer, error)
	FindBySceneID(sceneID int) ([]*Performer, error)
	FindNamesBySceneID(sceneID int) ([]*Performer, error)
	FindByImageID(imageID int) ([]*Performer, error)
	FindByGalleryID(galleryID int) ([]*Performer, error)
	FindByNames(names []string, nocase bool) ([]*Performer, error)
	FindByStashID(stashID StashID) ([]*Performer, error)
	FindByStashIDStatus(hasStashID bool, stashboxEndpoint string) ([]*Performer, error)
	CountByTagID(tagID int) (int, error)
	Count() (int, error)
	All() ([]*Performer, error)
	// TODO - this interface is temporary until the filter schema can fully
	// support the query needed
	QueryForAutoTag(words []string) ([]*Performer, error)
	Query(performerFilter *PerformerFilterType, findFilter *FindFilterType) ([]*Performer, int, error)
	GetImage(performerID int) ([]byte, error)
	GetStashIDs(performerID int) ([]*StashID, error)
	GetTagIDs(performerID int) ([]int, error)
}

type PerformerWriter interface {
	Create(newPerformer Performer) (*Performer, error)
	Update(updatedPerformer PerformerPartial) (*Performer, error)
	UpdateFull(updatedPerformer Performer) (*Performer, error)
	Destroy(id int) error
	UpdateImage(performerID int, image []byte) error
	DestroyImage(performerID int) error
	UpdateStashIDs(performerID int, stashIDs []StashID) error
	UpdateTags(performerID int, tagIDs []int) error
}

type PerformerReaderWriter interface {
	PerformerReader
	PerformerWriter
}
