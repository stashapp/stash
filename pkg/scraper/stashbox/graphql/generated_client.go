// Code generated by github.com/Yamashou/gqlgenc, DO NOT EDIT.

package graphql

import (
	"context"
	"net/http"

	"github.com/Yamashou/gqlgenc/client"
)

type StashBoxGraphQLClient interface {
	FindSceneByFingerprint(ctx context.Context, fingerprint FingerprintQueryInput, httpRequestOptions ...client.HTTPRequestOption) (*FindSceneByFingerprint, error)
	FindScenesByFullFingerprints(ctx context.Context, fingerprints []*FingerprintQueryInput, httpRequestOptions ...client.HTTPRequestOption) (*FindScenesByFullFingerprints, error)
	FindScenesBySceneFingerprints(ctx context.Context, fingerprints [][]*FingerprintQueryInput, httpRequestOptions ...client.HTTPRequestOption) (*FindScenesBySceneFingerprints, error)
	SearchScene(ctx context.Context, term string, httpRequestOptions ...client.HTTPRequestOption) (*SearchScene, error)
	SearchPerformer(ctx context.Context, term string, httpRequestOptions ...client.HTTPRequestOption) (*SearchPerformer, error)
	FindPerformerByID(ctx context.Context, id string, httpRequestOptions ...client.HTTPRequestOption) (*FindPerformerByID, error)
	FindSceneByID(ctx context.Context, id string, httpRequestOptions ...client.HTTPRequestOption) (*FindSceneByID, error)
	SubmitFingerprint(ctx context.Context, input FingerprintSubmission, httpRequestOptions ...client.HTTPRequestOption) (*SubmitFingerprint, error)
	Me(ctx context.Context, httpRequestOptions ...client.HTTPRequestOption) (*Me, error)
	SubmitSceneDraft(ctx context.Context, input SceneDraftInput, httpRequestOptions ...client.HTTPRequestOption) (*SubmitSceneDraft, error)
	SubmitPerformerDraft(ctx context.Context, input PerformerDraftInput, httpRequestOptions ...client.HTTPRequestOption) (*SubmitPerformerDraft, error)
}

type Client struct {
	Client *client.Client
}

func NewClient(cli *http.Client, baseURL string, options ...client.HTTPRequestOption) StashBoxGraphQLClient {
	return &Client{Client: client.NewClient(cli, baseURL, options...)}
}

type Query struct {
	FindPerformer                 *Performer                   "json:\"findPerformer\" graphql:\"findPerformer\""
	QueryPerformers               QueryPerformersResultType    "json:\"queryPerformers\" graphql:\"queryPerformers\""
	FindStudio                    *Studio                      "json:\"findStudio\" graphql:\"findStudio\""
	QueryStudios                  QueryStudiosResultType       "json:\"queryStudios\" graphql:\"queryStudios\""
	FindTag                       *Tag                         "json:\"findTag\" graphql:\"findTag\""
	QueryTags                     QueryTagsResultType          "json:\"queryTags\" graphql:\"queryTags\""
	FindTagCategory               *TagCategory                 "json:\"findTagCategory\" graphql:\"findTagCategory\""
	QueryTagCategories            QueryTagCategoriesResultType "json:\"queryTagCategories\" graphql:\"queryTagCategories\""
	FindScene                     *Scene                       "json:\"findScene\" graphql:\"findScene\""
	FindSceneByFingerprint        []*Scene                     "json:\"findSceneByFingerprint\" graphql:\"findSceneByFingerprint\""
	FindScenesByFingerprints      []*Scene                     "json:\"findScenesByFingerprints\" graphql:\"findScenesByFingerprints\""
	FindScenesByFullFingerprints  []*Scene                     "json:\"findScenesByFullFingerprints\" graphql:\"findScenesByFullFingerprints\""
	FindScenesBySceneFingerprints [][]*Scene                   "json:\"findScenesBySceneFingerprints\" graphql:\"findScenesBySceneFingerprints\""
	QueryScenes                   QueryScenesResultType        "json:\"queryScenes\" graphql:\"queryScenes\""
	FindSite                      *Site                        "json:\"findSite\" graphql:\"findSite\""
	QuerySites                    QuerySitesResultType         "json:\"querySites\" graphql:\"querySites\""
	FindEdit                      *Edit                        "json:\"findEdit\" graphql:\"findEdit\""
	QueryEdits                    QueryEditsResultType         "json:\"queryEdits\" graphql:\"queryEdits\""
	FindUser                      *User                        "json:\"findUser\" graphql:\"findUser\""
	QueryUsers                    QueryUsersResultType         "json:\"queryUsers\" graphql:\"queryUsers\""
	Me                            *User                        "json:\"me\" graphql:\"me\""
	SearchPerformer               []*Performer                 "json:\"searchPerformer\" graphql:\"searchPerformer\""
	SearchScene                   []*Scene                     "json:\"searchScene\" graphql:\"searchScene\""
	FindDraft                     *Draft                       "json:\"findDraft\" graphql:\"findDraft\""
	FindDrafts                    []*Draft                     "json:\"findDrafts\" graphql:\"findDrafts\""
	Version                       Version                      "json:\"version\" graphql:\"version\""
	GetConfig                     StashBoxConfig               "json:\"getConfig\" graphql:\"getConfig\""
}
type Mutation struct {
	SceneCreate          *Scene                "json:\"sceneCreate\" graphql:\"sceneCreate\""
	SceneUpdate          *Scene                "json:\"sceneUpdate\" graphql:\"sceneUpdate\""
	SceneDestroy         bool                  "json:\"sceneDestroy\" graphql:\"sceneDestroy\""
	PerformerCreate      *Performer            "json:\"performerCreate\" graphql:\"performerCreate\""
	PerformerUpdate      *Performer            "json:\"performerUpdate\" graphql:\"performerUpdate\""
	PerformerDestroy     bool                  "json:\"performerDestroy\" graphql:\"performerDestroy\""
	StudioCreate         *Studio               "json:\"studioCreate\" graphql:\"studioCreate\""
	StudioUpdate         *Studio               "json:\"studioUpdate\" graphql:\"studioUpdate\""
	StudioDestroy        bool                  "json:\"studioDestroy\" graphql:\"studioDestroy\""
	TagCreate            *Tag                  "json:\"tagCreate\" graphql:\"tagCreate\""
	TagUpdate            *Tag                  "json:\"tagUpdate\" graphql:\"tagUpdate\""
	TagDestroy           bool                  "json:\"tagDestroy\" graphql:\"tagDestroy\""
	UserCreate           *User                 "json:\"userCreate\" graphql:\"userCreate\""
	UserUpdate           *User                 "json:\"userUpdate\" graphql:\"userUpdate\""
	UserDestroy          bool                  "json:\"userDestroy\" graphql:\"userDestroy\""
	ImageCreate          *Image                "json:\"imageCreate\" graphql:\"imageCreate\""
	ImageDestroy         bool                  "json:\"imageDestroy\" graphql:\"imageDestroy\""
	NewUser              *string               "json:\"newUser\" graphql:\"newUser\""
	ActivateNewUser      *User                 "json:\"activateNewUser\" graphql:\"activateNewUser\""
	GenerateInviteCode   *string               "json:\"generateInviteCode\" graphql:\"generateInviteCode\""
	RescindInviteCode    bool                  "json:\"rescindInviteCode\" graphql:\"rescindInviteCode\""
	GrantInvite          int                   "json:\"grantInvite\" graphql:\"grantInvite\""
	RevokeInvite         int                   "json:\"revokeInvite\" graphql:\"revokeInvite\""
	TagCategoryCreate    *TagCategory          "json:\"tagCategoryCreate\" graphql:\"tagCategoryCreate\""
	TagCategoryUpdate    *TagCategory          "json:\"tagCategoryUpdate\" graphql:\"tagCategoryUpdate\""
	TagCategoryDestroy   bool                  "json:\"tagCategoryDestroy\" graphql:\"tagCategoryDestroy\""
	SiteCreate           *Site                 "json:\"siteCreate\" graphql:\"siteCreate\""
	SiteUpdate           *Site                 "json:\"siteUpdate\" graphql:\"siteUpdate\""
	SiteDestroy          bool                  "json:\"siteDestroy\" graphql:\"siteDestroy\""
	RegenerateAPIKey     string                "json:\"regenerateAPIKey\" graphql:\"regenerateAPIKey\""
	ResetPassword        bool                  "json:\"resetPassword\" graphql:\"resetPassword\""
	ChangePassword       bool                  "json:\"changePassword\" graphql:\"changePassword\""
	SceneEdit            Edit                  "json:\"sceneEdit\" graphql:\"sceneEdit\""
	PerformerEdit        Edit                  "json:\"performerEdit\" graphql:\"performerEdit\""
	StudioEdit           Edit                  "json:\"studioEdit\" graphql:\"studioEdit\""
	TagEdit              Edit                  "json:\"tagEdit\" graphql:\"tagEdit\""
	EditVote             Edit                  "json:\"editVote\" graphql:\"editVote\""
	EditComment          Edit                  "json:\"editComment\" graphql:\"editComment\""
	ApplyEdit            Edit                  "json:\"applyEdit\" graphql:\"applyEdit\""
	CancelEdit           Edit                  "json:\"cancelEdit\" graphql:\"cancelEdit\""
	SubmitFingerprint    bool                  "json:\"submitFingerprint\" graphql:\"submitFingerprint\""
	SubmitSceneDraft     DraftSubmissionStatus "json:\"submitSceneDraft\" graphql:\"submitSceneDraft\""
	SubmitPerformerDraft DraftSubmissionStatus "json:\"submitPerformerDraft\" graphql:\"submitPerformerDraft\""
	DestroyDraft         bool                  "json:\"destroyDraft\" graphql:\"destroyDraft\""
	FavoritePerformer    bool                  "json:\"favoritePerformer\" graphql:\"favoritePerformer\""
	FavoriteStudio       bool                  "json:\"favoriteStudio\" graphql:\"favoriteStudio\""
}
type URLFragment struct {
	URL  string "json:\"url\" graphql:\"url\""
	Type string "json:\"type\" graphql:\"type\""
}
type ImageFragment struct {
	ID     string "json:\"id\" graphql:\"id\""
	URL    string "json:\"url\" graphql:\"url\""
	Width  int    "json:\"width\" graphql:\"width\""
	Height int    "json:\"height\" graphql:\"height\""
}
type StudioFragment struct {
	Name   string           "json:\"name\" graphql:\"name\""
	ID     string           "json:\"id\" graphql:\"id\""
	Urls   []*URLFragment   "json:\"urls\" graphql:\"urls\""
	Images []*ImageFragment "json:\"images\" graphql:\"images\""
}
type TagFragment struct {
	Name string "json:\"name\" graphql:\"name\""
	ID   string "json:\"id\" graphql:\"id\""
}
type FuzzyDateFragment struct {
	Date     string           "json:\"date\" graphql:\"date\""
	Accuracy DateAccuracyEnum "json:\"accuracy\" graphql:\"accuracy\""
}
type MeasurementsFragment struct {
	BandSize *int    "json:\"band_size\" graphql:\"band_size\""
	CupSize  *string "json:\"cup_size\" graphql:\"cup_size\""
	Waist    *int    "json:\"waist\" graphql:\"waist\""
	Hip      *int    "json:\"hip\" graphql:\"hip\""
}
type BodyModificationFragment struct {
	Location    string  "json:\"location\" graphql:\"location\""
	Description *string "json:\"description\" graphql:\"description\""
}
type PerformerFragment struct {
	ID              string                      "json:\"id\" graphql:\"id\""
	Name            string                      "json:\"name\" graphql:\"name\""
	Disambiguation  *string                     "json:\"disambiguation\" graphql:\"disambiguation\""
	Aliases         []string                    "json:\"aliases\" graphql:\"aliases\""
	Gender          *GenderEnum                 "json:\"gender\" graphql:\"gender\""
	MergedIds       []string                    "json:\"merged_ids\" graphql:\"merged_ids\""
	Urls            []*URLFragment              "json:\"urls\" graphql:\"urls\""
	Images          []*ImageFragment            "json:\"images\" graphql:\"images\""
	Birthdate       *FuzzyDateFragment          "json:\"birthdate\" graphql:\"birthdate\""
	Ethnicity       *EthnicityEnum              "json:\"ethnicity\" graphql:\"ethnicity\""
	Country         *string                     "json:\"country\" graphql:\"country\""
	EyeColor        *EyeColorEnum               "json:\"eye_color\" graphql:\"eye_color\""
	HairColor       *HairColorEnum              "json:\"hair_color\" graphql:\"hair_color\""
	Height          *int                        "json:\"height\" graphql:\"height\""
	Measurements    MeasurementsFragment        "json:\"measurements\" graphql:\"measurements\""
	BreastType      *BreastTypeEnum             "json:\"breast_type\" graphql:\"breast_type\""
	CareerStartYear *int                        "json:\"career_start_year\" graphql:\"career_start_year\""
	CareerEndYear   *int                        "json:\"career_end_year\" graphql:\"career_end_year\""
	Tattoos         []*BodyModificationFragment "json:\"tattoos\" graphql:\"tattoos\""
	Piercings       []*BodyModificationFragment "json:\"piercings\" graphql:\"piercings\""
}
type PerformerAppearanceFragment struct {
	As        *string           "json:\"as\" graphql:\"as\""
	Performer PerformerFragment "json:\"performer\" graphql:\"performer\""
}
type FingerprintFragment struct {
	Algorithm FingerprintAlgorithm "json:\"algorithm\" graphql:\"algorithm\""
	Hash      string               "json:\"hash\" graphql:\"hash\""
	Duration  int                  "json:\"duration\" graphql:\"duration\""
}
type SceneFragment struct {
	ID           string                         "json:\"id\" graphql:\"id\""
	Title        *string                        "json:\"title\" graphql:\"title\""
	Details      *string                        "json:\"details\" graphql:\"details\""
	Duration     *int                           "json:\"duration\" graphql:\"duration\""
	Date         *string                        "json:\"date\" graphql:\"date\""
	Urls         []*URLFragment                 "json:\"urls\" graphql:\"urls\""
	Images       []*ImageFragment               "json:\"images\" graphql:\"images\""
	Studio       *StudioFragment                "json:\"studio\" graphql:\"studio\""
	Tags         []*TagFragment                 "json:\"tags\" graphql:\"tags\""
	Performers   []*PerformerAppearanceFragment "json:\"performers\" graphql:\"performers\""
	Fingerprints []*FingerprintFragment         "json:\"fingerprints\" graphql:\"fingerprints\""
}
type FindSceneByFingerprint struct {
	FindSceneByFingerprint []*SceneFragment "json:\"findSceneByFingerprint\" graphql:\"findSceneByFingerprint\""
}
type FindScenesByFullFingerprints struct {
	FindScenesByFullFingerprints []*SceneFragment "json:\"findScenesByFullFingerprints\" graphql:\"findScenesByFullFingerprints\""
}
type FindScenesBySceneFingerprints struct {
	FindScenesBySceneFingerprints [][]*SceneFragment "json:\"findScenesBySceneFingerprints\" graphql:\"findScenesBySceneFingerprints\""
}
type SearchScene struct {
	SearchScene []*SceneFragment "json:\"searchScene\" graphql:\"searchScene\""
}
type SearchPerformer struct {
	SearchPerformer []*PerformerFragment "json:\"searchPerformer\" graphql:\"searchPerformer\""
}
type FindPerformerByID struct {
	FindPerformer *PerformerFragment "json:\"findPerformer\" graphql:\"findPerformer\""
}
type FindSceneByID struct {
	FindScene *SceneFragment "json:\"findScene\" graphql:\"findScene\""
}
type SubmitFingerprint struct {
	SubmitFingerprint bool "json:\"submitFingerprint\" graphql:\"submitFingerprint\""
}
type Me struct {
	Me *struct {
		Name string "json:\"name\" graphql:\"name\""
	} "json:\"me\" graphql:\"me\""
}
type SubmitSceneDraft struct {
	SubmitSceneDraft struct {
		ID *string "json:\"id\" graphql:\"id\""
	} "json:\"submitSceneDraft\" graphql:\"submitSceneDraft\""
}
type SubmitPerformerDraft struct {
	SubmitPerformerDraft struct {
		ID *string "json:\"id\" graphql:\"id\""
	} "json:\"submitPerformerDraft\" graphql:\"submitPerformerDraft\""
}

const FindSceneByFingerprintDocument = `query FindSceneByFingerprint ($fingerprint: FingerprintQueryInput!) {
	findSceneByFingerprint(fingerprint: $fingerprint) {
		... SceneFragment
	}
}
fragment SceneFragment on Scene {
	id
	title
	details
	duration
	date
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
	studio {
		... StudioFragment
	}
	tags {
		... TagFragment
	}
	performers {
		... PerformerAppearanceFragment
	}
	fingerprints {
		... FingerprintFragment
	}
}
fragment ImageFragment on Image {
	id
	url
	width
	height
}
fragment PerformerAppearanceFragment on PerformerAppearance {
	as
	performer {
		... PerformerFragment
	}
}
fragment PerformerFragment on Performer {
	id
	name
	disambiguation
	aliases
	gender
	merged_ids
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
	birthdate {
		... FuzzyDateFragment
	}
	ethnicity
	country
	eye_color
	hair_color
	height
	measurements {
		... MeasurementsFragment
	}
	breast_type
	career_start_year
	career_end_year
	tattoos {
		... BodyModificationFragment
	}
	piercings {
		... BodyModificationFragment
	}
}
fragment FuzzyDateFragment on FuzzyDate {
	date
	accuracy
}
fragment URLFragment on URL {
	url
	type
}
fragment StudioFragment on Studio {
	name
	id
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
}
fragment TagFragment on Tag {
	name
	id
}
fragment MeasurementsFragment on Measurements {
	band_size
	cup_size
	waist
	hip
}
fragment BodyModificationFragment on BodyModification {
	location
	description
}
fragment FingerprintFragment on Fingerprint {
	algorithm
	hash
	duration
}
`

func (c *Client) FindSceneByFingerprint(ctx context.Context, fingerprint FingerprintQueryInput, httpRequestOptions ...client.HTTPRequestOption) (*FindSceneByFingerprint, error) {
	vars := map[string]interface{}{
		"fingerprint": fingerprint,
	}

	var res FindSceneByFingerprint
	if err := c.Client.Post(ctx, "FindSceneByFingerprint", FindSceneByFingerprintDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const FindScenesByFullFingerprintsDocument = `query FindScenesByFullFingerprints ($fingerprints: [FingerprintQueryInput!]!) {
	findScenesByFullFingerprints(fingerprints: $fingerprints) {
		... SceneFragment
	}
}
fragment ImageFragment on Image {
	id
	url
	width
	height
}
fragment PerformerAppearanceFragment on PerformerAppearance {
	as
	performer {
		... PerformerFragment
	}
}
fragment PerformerFragment on Performer {
	id
	name
	disambiguation
	aliases
	gender
	merged_ids
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
	birthdate {
		... FuzzyDateFragment
	}
	ethnicity
	country
	eye_color
	hair_color
	height
	measurements {
		... MeasurementsFragment
	}
	breast_type
	career_start_year
	career_end_year
	tattoos {
		... BodyModificationFragment
	}
	piercings {
		... BodyModificationFragment
	}
}
fragment MeasurementsFragment on Measurements {
	band_size
	cup_size
	waist
	hip
}
fragment SceneFragment on Scene {
	id
	title
	details
	duration
	date
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
	studio {
		... StudioFragment
	}
	tags {
		... TagFragment
	}
	performers {
		... PerformerAppearanceFragment
	}
	fingerprints {
		... FingerprintFragment
	}
}
fragment StudioFragment on Studio {
	name
	id
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
}
fragment TagFragment on Tag {
	name
	id
}
fragment FuzzyDateFragment on FuzzyDate {
	date
	accuracy
}
fragment BodyModificationFragment on BodyModification {
	location
	description
}
fragment FingerprintFragment on Fingerprint {
	algorithm
	hash
	duration
}
fragment URLFragment on URL {
	url
	type
}
`

func (c *Client) FindScenesByFullFingerprints(ctx context.Context, fingerprints []*FingerprintQueryInput, httpRequestOptions ...client.HTTPRequestOption) (*FindScenesByFullFingerprints, error) {
	vars := map[string]interface{}{
		"fingerprints": fingerprints,
	}

	var res FindScenesByFullFingerprints
	if err := c.Client.Post(ctx, "FindScenesByFullFingerprints", FindScenesByFullFingerprintsDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const FindScenesBySceneFingerprintsDocument = `query FindScenesBySceneFingerprints ($fingerprints: [[FingerprintQueryInput!]!]!) {
	findScenesBySceneFingerprints(fingerprints: $fingerprints) {
		... SceneFragment
	}
}
fragment PerformerFragment on Performer {
	id
	name
	disambiguation
	aliases
	gender
	merged_ids
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
	birthdate {
		... FuzzyDateFragment
	}
	ethnicity
	country
	eye_color
	hair_color
	height
	measurements {
		... MeasurementsFragment
	}
	breast_type
	career_start_year
	career_end_year
	tattoos {
		... BodyModificationFragment
	}
	piercings {
		... BodyModificationFragment
	}
}
fragment FuzzyDateFragment on FuzzyDate {
	date
	accuracy
}
fragment StudioFragment on Studio {
	name
	id
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
}
fragment TagFragment on Tag {
	name
	id
}
fragment ImageFragment on Image {
	id
	url
	width
	height
}
fragment PerformerAppearanceFragment on PerformerAppearance {
	as
	performer {
		... PerformerFragment
	}
}
fragment MeasurementsFragment on Measurements {
	band_size
	cup_size
	waist
	hip
}
fragment BodyModificationFragment on BodyModification {
	location
	description
}
fragment FingerprintFragment on Fingerprint {
	algorithm
	hash
	duration
}
fragment SceneFragment on Scene {
	id
	title
	details
	duration
	date
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
	studio {
		... StudioFragment
	}
	tags {
		... TagFragment
	}
	performers {
		... PerformerAppearanceFragment
	}
	fingerprints {
		... FingerprintFragment
	}
}
fragment URLFragment on URL {
	url
	type
}
`

func (c *Client) FindScenesBySceneFingerprints(ctx context.Context, fingerprints [][]*FingerprintQueryInput, httpRequestOptions ...client.HTTPRequestOption) (*FindScenesBySceneFingerprints, error) {
	vars := map[string]interface{}{
		"fingerprints": fingerprints,
	}

	var res FindScenesBySceneFingerprints
	if err := c.Client.Post(ctx, "FindScenesBySceneFingerprints", FindScenesBySceneFingerprintsDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const SearchSceneDocument = `query SearchScene ($term: String!) {
	searchScene(term: $term) {
		... SceneFragment
	}
}
fragment ImageFragment on Image {
	id
	url
	width
	height
}
fragment BodyModificationFragment on BodyModification {
	location
	description
}
fragment PerformerFragment on Performer {
	id
	name
	disambiguation
	aliases
	gender
	merged_ids
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
	birthdate {
		... FuzzyDateFragment
	}
	ethnicity
	country
	eye_color
	hair_color
	height
	measurements {
		... MeasurementsFragment
	}
	breast_type
	career_start_year
	career_end_year
	tattoos {
		... BodyModificationFragment
	}
	piercings {
		... BodyModificationFragment
	}
}
fragment FuzzyDateFragment on FuzzyDate {
	date
	accuracy
}
fragment MeasurementsFragment on Measurements {
	band_size
	cup_size
	waist
	hip
}
fragment SceneFragment on Scene {
	id
	title
	details
	duration
	date
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
	studio {
		... StudioFragment
	}
	tags {
		... TagFragment
	}
	performers {
		... PerformerAppearanceFragment
	}
	fingerprints {
		... FingerprintFragment
	}
}
fragment URLFragment on URL {
	url
	type
}
fragment StudioFragment on Studio {
	name
	id
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
}
fragment TagFragment on Tag {
	name
	id
}
fragment PerformerAppearanceFragment on PerformerAppearance {
	as
	performer {
		... PerformerFragment
	}
}
fragment FingerprintFragment on Fingerprint {
	algorithm
	hash
	duration
}
`

func (c *Client) SearchScene(ctx context.Context, term string, httpRequestOptions ...client.HTTPRequestOption) (*SearchScene, error) {
	vars := map[string]interface{}{
		"term": term,
	}

	var res SearchScene
	if err := c.Client.Post(ctx, "SearchScene", SearchSceneDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const SearchPerformerDocument = `query SearchPerformer ($term: String!) {
	searchPerformer(term: $term) {
		... PerformerFragment
	}
}
fragment FuzzyDateFragment on FuzzyDate {
	date
	accuracy
}
fragment MeasurementsFragment on Measurements {
	band_size
	cup_size
	waist
	hip
}
fragment BodyModificationFragment on BodyModification {
	location
	description
}
fragment PerformerFragment on Performer {
	id
	name
	disambiguation
	aliases
	gender
	merged_ids
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
	birthdate {
		... FuzzyDateFragment
	}
	ethnicity
	country
	eye_color
	hair_color
	height
	measurements {
		... MeasurementsFragment
	}
	breast_type
	career_start_year
	career_end_year
	tattoos {
		... BodyModificationFragment
	}
	piercings {
		... BodyModificationFragment
	}
}
fragment URLFragment on URL {
	url
	type
}
fragment ImageFragment on Image {
	id
	url
	width
	height
}
`

func (c *Client) SearchPerformer(ctx context.Context, term string, httpRequestOptions ...client.HTTPRequestOption) (*SearchPerformer, error) {
	vars := map[string]interface{}{
		"term": term,
	}

	var res SearchPerformer
	if err := c.Client.Post(ctx, "SearchPerformer", SearchPerformerDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const FindPerformerByIDDocument = `query FindPerformerByID ($id: ID!) {
	findPerformer(id: $id) {
		... PerformerFragment
	}
}
fragment PerformerFragment on Performer {
	id
	name
	disambiguation
	aliases
	gender
	merged_ids
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
	birthdate {
		... FuzzyDateFragment
	}
	ethnicity
	country
	eye_color
	hair_color
	height
	measurements {
		... MeasurementsFragment
	}
	breast_type
	career_start_year
	career_end_year
	tattoos {
		... BodyModificationFragment
	}
	piercings {
		... BodyModificationFragment
	}
}
fragment URLFragment on URL {
	url
	type
}
fragment ImageFragment on Image {
	id
	url
	width
	height
}
fragment FuzzyDateFragment on FuzzyDate {
	date
	accuracy
}
fragment MeasurementsFragment on Measurements {
	band_size
	cup_size
	waist
	hip
}
fragment BodyModificationFragment on BodyModification {
	location
	description
}
`

func (c *Client) FindPerformerByID(ctx context.Context, id string, httpRequestOptions ...client.HTTPRequestOption) (*FindPerformerByID, error) {
	vars := map[string]interface{}{
		"id": id,
	}

	var res FindPerformerByID
	if err := c.Client.Post(ctx, "FindPerformerByID", FindPerformerByIDDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const FindSceneByIDDocument = `query FindSceneByID ($id: ID!) {
	findScene(id: $id) {
		... SceneFragment
	}
}
fragment URLFragment on URL {
	url
	type
}
fragment ImageFragment on Image {
	id
	url
	width
	height
}
fragment StudioFragment on Studio {
	name
	id
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
}
fragment TagFragment on Tag {
	name
	id
}
fragment PerformerAppearanceFragment on PerformerAppearance {
	as
	performer {
		... PerformerFragment
	}
}
fragment PerformerFragment on Performer {
	id
	name
	disambiguation
	aliases
	gender
	merged_ids
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
	birthdate {
		... FuzzyDateFragment
	}
	ethnicity
	country
	eye_color
	hair_color
	height
	measurements {
		... MeasurementsFragment
	}
	breast_type
	career_start_year
	career_end_year
	tattoos {
		... BodyModificationFragment
	}
	piercings {
		... BodyModificationFragment
	}
}
fragment FuzzyDateFragment on FuzzyDate {
	date
	accuracy
}
fragment SceneFragment on Scene {
	id
	title
	details
	duration
	date
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
	studio {
		... StudioFragment
	}
	tags {
		... TagFragment
	}
	performers {
		... PerformerAppearanceFragment
	}
	fingerprints {
		... FingerprintFragment
	}
}
fragment MeasurementsFragment on Measurements {
	band_size
	cup_size
	waist
	hip
}
fragment BodyModificationFragment on BodyModification {
	location
	description
}
fragment FingerprintFragment on Fingerprint {
	algorithm
	hash
	duration
}
`

func (c *Client) FindSceneByID(ctx context.Context, id string, httpRequestOptions ...client.HTTPRequestOption) (*FindSceneByID, error) {
	vars := map[string]interface{}{
		"id": id,
	}

	var res FindSceneByID
	if err := c.Client.Post(ctx, "FindSceneByID", FindSceneByIDDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const SubmitFingerprintDocument = `mutation SubmitFingerprint ($input: FingerprintSubmission!) {
	submitFingerprint(input: $input)
}
`

func (c *Client) SubmitFingerprint(ctx context.Context, input FingerprintSubmission, httpRequestOptions ...client.HTTPRequestOption) (*SubmitFingerprint, error) {
	vars := map[string]interface{}{
		"input": input,
	}

	var res SubmitFingerprint
	if err := c.Client.Post(ctx, "SubmitFingerprint", SubmitFingerprintDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const MeDocument = `query Me {
	me {
		name
	}
}
`

func (c *Client) Me(ctx context.Context, httpRequestOptions ...client.HTTPRequestOption) (*Me, error) {
	vars := map[string]interface{}{}

	var res Me
	if err := c.Client.Post(ctx, "Me", MeDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const SubmitSceneDraftDocument = `mutation SubmitSceneDraft ($input: SceneDraftInput!) {
	submitSceneDraft(input: $input) {
		id
	}
}
`

func (c *Client) SubmitSceneDraft(ctx context.Context, input SceneDraftInput, httpRequestOptions ...client.HTTPRequestOption) (*SubmitSceneDraft, error) {
	vars := map[string]interface{}{
		"input": input,
	}

	var res SubmitSceneDraft
	if err := c.Client.Post(ctx, "SubmitSceneDraft", SubmitSceneDraftDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const SubmitPerformerDraftDocument = `mutation SubmitPerformerDraft ($input: PerformerDraftInput!) {
	submitPerformerDraft(input: $input) {
		id
	}
}
`

func (c *Client) SubmitPerformerDraft(ctx context.Context, input PerformerDraftInput, httpRequestOptions ...client.HTTPRequestOption) (*SubmitPerformerDraft, error) {
	vars := map[string]interface{}{
		"input": input,
	}

	var res SubmitPerformerDraft
	if err := c.Client.Post(ctx, "SubmitPerformerDraft", SubmitPerformerDraftDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}
