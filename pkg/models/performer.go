package models

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
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

type CircumcisedEnum string

const (
	CircumcisedEnumCut   CircumcisedEnum = "CUT"
	CircumcisedEnumUncut CircumcisedEnum = "UNCUT"
)

var AllCircumcisionEnum = []CircumcisedEnum{
	CircumcisedEnumCut,
	CircumcisedEnumUncut,
}

func (e CircumcisedEnum) IsValid() bool {
	switch e {
	case CircumcisedEnumCut, CircumcisedEnumUncut:
		return true
	}
	return false
}

func (e CircumcisedEnum) String() string {
	return string(e)
}

func (e *CircumcisedEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = CircumcisedEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid CircumcisedEnum", str)
	}
	return nil
}

func (e CircumcisedEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type CircumcisionCriterionInput struct {
	Value    []CircumcisedEnum `json:"value"`
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
	// Filter by band size
	BandSize *IntCriterionInput `json:"band_size"`
	// Filter by cup size
	CupSize *StringCriterionInput `json:"cup_size"`
	// Filter by waist size
	WaistSize *IntCriterionInput `json:"waist_size"`
	// Filter by hip size
	HipSize *IntCriterionInput `json:"hip_size"`
	// Filter by measurements - deprecated: use band_size/cup_size/waist_size/hip_size
	Measurements *StringCriterionInput `json:"measurements"`
	// Filter by fake tits value
	FakeTits *StringCriterionInput `json:"fake_tits"`
	// Filter by penis length value
	PenisLength *FloatCriterionInput `json:"penis_length"`
	// Filter by circumcision
	Circumcised *CircumcisionCriterionInput `json:"circumcised"`
	// Filter by career length
	CareerLength *StringCriterionInput `json:"career_length"` // deprecated
	// Filter by career start year
	CareerStart *DateCriterionInput `json:"career_start"`
	// Filter by career end year
	CareerEnd *DateCriterionInput `json:"career_end"`
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
	// Filter by scene marker count (via scene)
	MarkerCount *IntCriterionInput `json:"marker_count"`
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
	// Filter by StashIDs Endpoint
	StashIDsEndpoint *StashIDsCriterionInput `json:"stash_ids_endpoint"`
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
	// Filter by groups where performer appears in scene
	Groups *HierarchicalMultiCriterionInput `json:"groups"`
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
	// Filter by related scene markers (via scene) that meet this criteria
	MarkersFilter *SceneMarkerFilterType `json:"markers_filter"`
	// Filter by created at
	CreatedAt *TimestampCriterionInput `json:"created_at"`
	// Filter by updated at
	UpdatedAt *TimestampCriterionInput `json:"updated_at"`

	// Filter by custom fields
	CustomFields []CustomFieldCriterionInput `json:"custom_fields"`
}

type PerformerCreateInput struct {
	Name           string           `json:"name"`
	Disambiguation *string          `json:"disambiguation"`
	URL            *string          `json:"url"` // deprecated
	Urls           []string         `json:"urls"`
	Gender         *GenderEnum      `json:"gender"`
	Birthdate      *string          `json:"birthdate"`
	Ethnicity      *string          `json:"ethnicity"`
	Country        *string          `json:"country"`
	EyeColor       *string          `json:"eye_color"`
	Height         *string          `json:"height"`
	HeightCm       *int             `json:"height_cm"`
	BandSize       *int             `json:"band_size"`
	CupSize        *string          `json:"cup_size"`
	WaistSize      *int             `json:"waist_size"`
	HipSize        *int             `json:"hip_size"`
	Measurements   *string          `json:"measurements"` // deprecated: use band_size/cup_size/waist_size/hip_size
	FakeTits       *string          `json:"fake_tits"`
	PenisLength    *float64         `json:"penis_length"`
	Circumcised    *CircumcisedEnum `json:"circumcised"`
	CareerLength   *string          `json:"career_length"`
	CareerStart    *string          `json:"career_start"`
	CareerEnd      *string          `json:"career_end"`
	Tattoos        *string          `json:"tattoos"`
	Piercings      *string          `json:"piercings"`
	Aliases        *string          `json:"aliases"`
	AliasList      []string         `json:"alias_list"`
	Twitter        *string          `json:"twitter"`   // deprecated
	Instagram      *string          `json:"instagram"` // deprecated
	Favorite       *bool            `json:"favorite"`
	TagIds         []string         `json:"tag_ids"`
	// This should be a URL or a base64 encoded data URL
	Image         *string        `json:"image"`
	StashIds      []StashIDInput `json:"stash_ids"`
	Rating100     *int           `json:"rating100"`
	Details       *string        `json:"details"`
	DeathDate     *string        `json:"death_date"`
	HairColor     *string        `json:"hair_color"`
	Weight        *int           `json:"weight"`
	IgnoreAutoTag *bool          `json:"ignore_auto_tag"`

	CustomFields map[string]interface{} `json:"custom_fields"`
}

type PerformerUpdateInput struct {
	ID             string           `json:"id"`
	Name           *string          `json:"name"`
	Disambiguation *string          `json:"disambiguation"`
	URL            *string          `json:"url"` // deprecated
	Urls           []string         `json:"urls"`
	Gender         *GenderEnum      `json:"gender"`
	Birthdate      *string          `json:"birthdate"`
	Ethnicity      *string          `json:"ethnicity"`
	Country        *string          `json:"country"`
	EyeColor       *string          `json:"eye_color"`
	Height         *string          `json:"height"`
	HeightCm       *int             `json:"height_cm"`
	BandSize       *int             `json:"band_size"`
	CupSize        *string          `json:"cup_size"`
	WaistSize      *int             `json:"waist_size"`
	HipSize        *int             `json:"hip_size"`
	Measurements   *string          `json:"measurements"` // deprecated: use band_size/cup_size/waist_size/hip_size
	FakeTits       *string          `json:"fake_tits"`
	PenisLength    *float64         `json:"penis_length"`
	Circumcised    *CircumcisedEnum `json:"circumcised"`
	CareerLength   *string          `json:"career_length"`
	CareerStart    *string          `json:"career_start"`
	CareerEnd      *string          `json:"career_end"`
	Tattoos        *string          `json:"tattoos"`
	Piercings      *string          `json:"piercings"`
	Aliases        *string          `json:"aliases"`
	AliasList      []string         `json:"alias_list"`
	Twitter        *string          `json:"twitter"`   // deprecated
	Instagram      *string          `json:"instagram"` // deprecated
	Favorite       *bool            `json:"favorite"`
	TagIds         []string         `json:"tag_ids"`
	// This should be a URL or a base64 encoded data URL
	Image         *string        `json:"image"`
	StashIds      []StashIDInput `json:"stash_ids"`
	Rating100     *int           `json:"rating100"`
	Details       *string        `json:"details"`
	DeathDate     *string        `json:"death_date"`
	HairColor     *string        `json:"hair_color"`
	Weight        *int           `json:"weight"`
	IgnoreAutoTag *bool          `json:"ignore_auto_tag"`

	CustomFields CustomFieldsInput `json:"custom_fields"`
}

// measurementsRe matches strings like "34DD-24-34" or "34B-24-34"
// Group 1: band (digits), Group 2: cup (letters), Group 3: waist (digits), Group 4: hip (digits)
var measurementsRe = regexp.MustCompile(`^(\d+)([A-Za-z]+)-(\d+)-(\d+)$`)

// ParseMeasurementsString parses a measurements string of the form "34DD-24-34"
// into its constituent parts. Returns an error if the string doesn't match the
// expected format.
func ParseMeasurementsString(s string) (bandSize *int, cupSize *string, waistSize *int, hipSize *int, err error) {
	s = strings.TrimSpace(s)
	m := measurementsRe.FindStringSubmatch(s)
	if m == nil {
		return nil, nil, nil, nil, fmt.Errorf("measurements %q does not match expected format (e.g. 34DD-24-34)", s)
	}

	band, _ := strconv.Atoi(m[1])
	cup := strings.ToUpper(m[2])
	waist, _ := strconv.Atoi(m[3])
	hip, _ := strconv.Atoi(m[4])

	return &band, &cup, &waist, &hip, nil
}

// FormatMeasurements formats the individual measurement fields back into the
// legacy "34DD-24-34" string format. Returns an empty string if all fields are nil/empty.
func FormatMeasurements(bandSize *int, cupSize string, waistSize *int, hipSize *int) string {
	if bandSize == nil && cupSize == "" && waistSize == nil && hipSize == nil {
		return ""
	}

	band := ""
	if bandSize != nil {
		band = strconv.Itoa(*bandSize)
	}
	waist := ""
	if waistSize != nil {
		waist = strconv.Itoa(*waistSize)
	}
	hip := ""
	if hipSize != nil {
		hip = strconv.Itoa(*hipSize)
	}

	return fmt.Sprintf("%s%s-%s-%s", band, cupSize, waist, hip)
}
