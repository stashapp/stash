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
	Value     GenderEnum        `json:"value"`
	ValueList []GenderEnum      `json:"value_list"`
	Modifier  CriterionModifier `json:"modifier"`
}

type CircumisedEnum string

const (
	CircumisedEnumCut   CircumisedEnum = "CUT"
	CircumisedEnumUncut CircumisedEnum = "UNCUT"
)

var AllCircumcisionEnum = []CircumisedEnum{
	CircumisedEnumCut,
	CircumisedEnumUncut,
}

func (e CircumisedEnum) IsValid() bool {
	switch e {
	case CircumisedEnumCut, CircumisedEnumUncut:
		return true
	}
	return false
}

func (e CircumisedEnum) String() string {
	return string(e)
}

func (e *CircumisedEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = CircumisedEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid CircumisedEnum", str)
	}
	return nil
}

func (e CircumisedEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type CircumcisionCriterionInput struct {
	Value    []CircumisedEnum  `json:"value"`
	Modifier CriterionModifier `json:"modifier"`
}

type PerformerFilterType struct {
	OperatorFilter[PerformerFilterType]
	Name           *StringCriterionInput `json:"name"`
	Disambiguation *StringCriterionInput `json:"disambiguation"`
	Details        *StringCriterionInput `json:"details"`
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
	// Filter by height - deprecated: use height_cm instead
	Height *StringCriterionInput `json:"height"`
	// Filter by height in centimeters
	HeightCm *IntCriterionInput `json:"height_cm"`
	// Filter by measurements
	Measurements *StringCriterionInput `json:"measurements"`
	// Filter by fake tits value
	FakeTits *StringCriterionInput `json:"fake_tits"`
	// Filter by penis length value
	PenisLength *FloatCriterionInput `json:"penis_length"`
	// Filter by circumcision
	Circumcised *CircumcisionCriterionInput `json:"circumcised"`
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
	// Filter by play count
	PlayCount *IntCriterionInput `json:"play_count"`
	// Filter by O count
	OCounter *IntCriterionInput `json:"o_counter"`
	// Filter by StashID
	StashID *StringCriterionInput `json:"stash_id"`
	// Filter by StashID Endpoint
	StashIDEndpoint *StashIDCriterionInput `json:"stash_id_endpoint"`
	// Filter by rating expressed as 1-100
	Rating100 *IntCriterionInput `json:"rating100"`
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
	// Filter by performers where performer appears with another performer in scene/image/gallery
	Performers *MultiCriterionInput `json:"performers"`
	// Filter by autotag ignore value
	IgnoreAutoTag *bool `json:"ignore_auto_tag"`
	// Filter by birthdate
	Birthdate *DateCriterionInput `json:"birth_date"`
	// Filter by death date
	DeathDate *DateCriterionInput `json:"death_date"`
	// Filter by related scenes that meet this criteria
	ScenesFilter *SceneFilterType `json:"scenes_filter"`
	// Filter by related images that meet this criteria
	ImagesFilter *ImageFilterType `json:"images_filter"`
	// Filter by related galleries that meet this criteria
	GalleriesFilter *GalleryFilterType `json:"galleries_filter"`
	// Filter by related tags that meet this criteria
	TagsFilter *TagFilterType `json:"tags_filter"`
	// Filter by created at
	CreatedAt *TimestampCriterionInput `json:"created_at"`
	// Filter by updated at
	UpdatedAt *TimestampCriterionInput `json:"updated_at"`
}

type PerformerCreateInput struct {
	Name           string          `json:"name"`
	Disambiguation *string         `json:"disambiguation"`
	URL            *string         `json:"url"`
	Gender         *GenderEnum     `json:"gender"`
	Birthdate      *string         `json:"birthdate"`
	Ethnicity      *string         `json:"ethnicity"`
	Country        *string         `json:"country"`
	EyeColor       *string         `json:"eye_color"`
	Height         *string         `json:"height"`
	HeightCm       *int            `json:"height_cm"`
	Measurements   *string         `json:"measurements"`
	FakeTits       *string         `json:"fake_tits"`
	PenisLength    *float64        `json:"penis_length"`
	Circumcised    *CircumisedEnum `json:"circumcised"`
	CareerLength   *string         `json:"career_length"`
	Tattoos        *string         `json:"tattoos"`
	Piercings      *string         `json:"piercings"`
	Aliases        *string         `json:"aliases"`
	AliasList      []string        `json:"alias_list"`
	Twitter        *string         `json:"twitter"`
	Instagram      *string         `json:"instagram"`
	Favorite       *bool           `json:"favorite"`
	TagIds         []string        `json:"tag_ids"`
	// This should be a URL or a base64 encoded data URL
	Image         *string   `json:"image"`
	StashIds      []StashID `json:"stash_ids"`
	Rating100     *int      `json:"rating100"`
	Details       *string   `json:"details"`
	DeathDate     *string   `json:"death_date"`
	HairColor     *string   `json:"hair_color"`
	Weight        *int      `json:"weight"`
	IgnoreAutoTag *bool     `json:"ignore_auto_tag"`
}

type PerformerUpdateInput struct {
	ID             string          `json:"id"`
	Name           *string         `json:"name"`
	Disambiguation *string         `json:"disambiguation"`
	URL            *string         `json:"url"`
	Gender         *GenderEnum     `json:"gender"`
	Birthdate      *string         `json:"birthdate"`
	Ethnicity      *string         `json:"ethnicity"`
	Country        *string         `json:"country"`
	EyeColor       *string         `json:"eye_color"`
	Height         *string         `json:"height"`
	HeightCm       *int            `json:"height_cm"`
	Measurements   *string         `json:"measurements"`
	FakeTits       *string         `json:"fake_tits"`
	PenisLength    *float64        `json:"penis_length"`
	Circumcised    *CircumisedEnum `json:"circumcised"`
	CareerLength   *string         `json:"career_length"`
	Tattoos        *string         `json:"tattoos"`
	Piercings      *string         `json:"piercings"`
	Aliases        *string         `json:"aliases"`
	AliasList      []string        `json:"alias_list"`
	Twitter        *string         `json:"twitter"`
	Instagram      *string         `json:"instagram"`
	Favorite       *bool           `json:"favorite"`
	TagIds         []string        `json:"tag_ids"`
	// This should be a URL or a base64 encoded data URL
	Image         *string   `json:"image"`
	StashIds      []StashID `json:"stash_ids"`
	Rating100     *int      `json:"rating100"`
	Details       *string   `json:"details"`
	DeathDate     *string   `json:"death_date"`
	HairColor     *string   `json:"hair_color"`
	Weight        *int      `json:"weight"`
	IgnoreAutoTag *bool     `json:"ignore_auto_tag"`
}
