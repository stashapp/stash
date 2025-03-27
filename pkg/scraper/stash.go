package scraper

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	graphql "github.com/hasura/go-graphql-client"
	"github.com/jinzhu/copier"

	"github.com/stashapp/stash/pkg/models"
)

type stashScraper struct {
	scraper      scraperTypeConfig
	config       config
	globalConfig GlobalConfig
	client       *http.Client
}

func newStashScraper(scraper scraperTypeConfig, client *http.Client, config config, globalConfig GlobalConfig) *stashScraper {
	return &stashScraper{
		scraper:      scraper,
		config:       config,
		client:       client,
		globalConfig: globalConfig,
	}
}

func setApiKeyHeader(apiKey string) func(req *http.Request) {
	return func(req *http.Request) {
		req.Header.Set("ApiKey", apiKey)
	}
}

func (s *stashScraper) getStashClient() *graphql.Client {
	url := s.config.StashServer.URL + "/graphql"
	ret := graphql.NewClient(url, s.client)

	if s.config.StashServer.ApiKey != "" {
		ret = ret.WithRequestModifier(setApiKeyHeader(s.config.StashServer.ApiKey))
	}

	return ret
}

type stashFindPerformerNamePerformer struct {
	ID   string `json:"id" graphql:"id"`
	Name string `json:"name" graphql:"name"`
}

func (p stashFindPerformerNamePerformer) toPerformer() *models.ScrapedPerformer {
	return &models.ScrapedPerformer{
		Name: &p.Name,
		// HACK - put id into the URL field
		URL: &p.ID,
	}
}

type stashFindPerformerNamesResultType struct {
	Count      int                                `graphql:"count"`
	Performers []*stashFindPerformerNamePerformer `graphql:"performers"`
}

// need a separate for scraped stash performers - does not include remote_site_id or image
type scrapedTagStash struct {
	Name string `graphql:"name" json:"name"`
}

type scrapedPerformerStash struct {
	Name         *string            `graphql:"name" json:"name"`
	Gender       *string            `graphql:"gender" json:"gender"`
	URLs         []string           `graphql:"urls" json:"urls"`
	Birthdate    *string            `graphql:"birthdate" json:"birthdate"`
	Ethnicity    *string            `graphql:"ethnicity" json:"ethnicity"`
	Country      *string            `graphql:"country" json:"country"`
	EyeColor     *string            `graphql:"eye_color" json:"eye_color"`
	Height       *int               `graphql:"height_cm" json:"height_cm"`
	Measurements *string            `graphql:"measurements" json:"measurements"`
	FakeTits     *string            `graphql:"fake_tits" json:"fake_tits"`
	PenisLength  *string            `graphql:"penis_length" json:"penis_length"`
	Circumcised  *string            `graphql:"circumcised" json:"circumcised"`
	CareerLength *string            `graphql:"career_length" json:"career_length"`
	Tattoos      *string            `graphql:"tattoos" json:"tattoos"`
	Piercings    *string            `graphql:"piercings" json:"piercings"`
	Aliases      []string           `graphql:"alias_list" json:"alias_list"`
	Tags         []*scrapedTagStash `graphql:"tags" json:"tags"`
	Details      *string            `graphql:"details" json:"details"`
	DeathDate    *string            `graphql:"death_date" json:"death_date"`
	HairColor    *string            `graphql:"hair_color" json:"hair_color"`
	Weight       *int               `graphql:"weight" json:"weight"`
}

func (s *stashScraper) imageGetter() imageGetter {
	ret := imageGetter{
		client:       s.client,
		globalConfig: s.globalConfig,
	}

	if s.config.StashServer.ApiKey != "" {
		ret.requestModifier = setApiKeyHeader(s.config.StashServer.ApiKey)
	}

	return ret
}

func (s *stashScraper) scrapeByFragment(ctx context.Context, input Input) (ScrapedContent, error) {
	if input.Performer != nil {
		return s.scrapeByPerformerFragment(ctx, *input.Performer)
	}

	if input.Scene != nil {
		return s.scrapeBySceneFragment(ctx, *input.Scene)
	}

	return nil, fmt.Errorf("%w: using stash scraper as a fragment scraper", ErrNotSupported)
}

func (s *stashScraper) scrapeByPerformerFragment(ctx context.Context, scrapedPerformer ScrapedPerformerInput) (ScrapedContent, error) {
	client := s.getStashClient()

	var q struct {
		FindPerformer *scrapedPerformerStash `graphql:"findPerformer(id: $f)"`
	}

	performerID := *scrapedPerformer.URL

	// get the id from the URL field
	vars := map[string]interface{}{
		"f": graphql.ID(performerID),
	}

	err := client.Query(ctx, &q, vars)
	if err != nil {
		return nil, convertGraphqlError(err)
	}

	// need to copy back to a scraped performer
	ret := models.ScrapedPerformer{}
	err = copier.Copy(&ret, q.FindPerformer)
	if err != nil {
		return nil, err
	}

	// convert alias list to aliases
	aliasStr := strings.Join(q.FindPerformer.Aliases, ", ")
	ret.Aliases = &aliasStr

	// convert numeric to string
	if q.FindPerformer.Height != nil {
		heightStr := strconv.Itoa(*q.FindPerformer.Height)
		ret.Height = &heightStr
	}
	if q.FindPerformer.Weight != nil {
		weightStr := strconv.Itoa(*q.FindPerformer.Weight)
		ret.Weight = &weightStr
	}

	// get the performer image directly
	ig := s.imageGetter()
	img, err := getStashPerformerImage(ctx, s.config.StashServer.URL, performerID, ig)
	if err != nil {
		return nil, err
	}
	ret.Images = []string{*img}
	ret.Image = img

	return &ret, nil
}

func (s *stashScraper) scrapeBySceneFragment(ctx context.Context, scrapedScene models.ScrapedSceneInput) (ScrapedContent, error) {
	client := s.getStashClient()

	var q struct {
		FindScene *scrapedSceneStash `graphql:"findScene(id: $f)"`
	}

	sceneID := scrapedScene.URLs[0]

	// get the id from the URL field
	vars := map[string]interface{}{
		"f": graphql.ID(sceneID),
	}

	err := client.Query(ctx, &q, vars)
	if err != nil {
		return nil, convertGraphqlError(err)
	}

	if q.FindScene == nil {
		return nil, nil
	}

	// need to copy back to a scraped scene
	ret, err := s.scrapedStashSceneToScrapedScene(ctx, q.FindScene)
	if err != nil {
		return nil, err
	}

	// get the scene image directly
	ig := s.imageGetter()
	ret.Image, err = getStashSceneImage(ctx, s.config.StashServer.URL, q.FindScene.ID, ig)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

type scrapedStudioStash struct {
	Name string  `graphql:"name" json:"name"`
	URL  *string `graphql:"url" json:"url"`
}

type stashFindSceneNamesResultType struct {
	Count  int                  `graphql:"count"`
	Scenes []*scrapedSceneStash `graphql:"scenes"`
}

func (s *stashScraper) scrapedStashSceneToScrapedScene(ctx context.Context, scene *scrapedSceneStash) (*models.ScrapedScene, error) {
	ret := models.ScrapedScene{}
	err := copier.Copy(&ret, scene)
	if err != nil {
		return nil, err
	}

	// convert first in files to file
	if len(scene.Files) > 0 {
		f := scene.Files[0].SceneFileType()
		ret.File = &f
	}

	// get the scene image directly
	ig := s.imageGetter()
	ret.Image, err = getStashSceneImage(ctx, s.config.StashServer.URL, scene.ID, ig)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (s *stashScraper) scrapeByName(ctx context.Context, name string, ty ScrapeContentType) ([]ScrapedContent, error) {
	client := s.getStashClient()

	page := 1
	perPage := 10

	vars := map[string]interface{}{
		"f": models.FindFilterType{
			Q:       &name,
			Page:    &page,
			PerPage: &perPage,
		},
	}

	var ret []ScrapedContent
	switch ty {
	case ScrapeContentTypeScene:
		var q struct {
			FindScenes stashFindSceneNamesResultType `graphql:"findScenes(filter: $f)"`
		}

		err := client.Query(ctx, &q, vars)
		if err != nil {
			return nil, convertGraphqlError(err)
		}

		for _, scene := range q.FindScenes.Scenes {
			converted, err := s.scrapedStashSceneToScrapedScene(ctx, scene)
			if err != nil {
				return nil, err
			}

			// HACK - put id into the URL field
			// put id into the URL field
			converted.URLs = []string{scene.ID}
			ret = append(ret, converted)
		}

		return ret, nil
	case ScrapeContentTypePerformer:
		var q struct {
			FindPerformers stashFindPerformerNamesResultType `graphql:"findPerformers(filter: $f)"`
		}

		err := client.Query(ctx, &q, vars)
		if err != nil {
			return nil, err
		}

		for _, p := range q.FindPerformers.Performers {
			ret = append(ret, p.toPerformer())
		}

		return ret, nil
	}

	return nil, ErrNotSupported
}

type stashVideoFile struct {
	Size       int64   `graphql:"size" json:"size"`
	Duration   float64 `graphql:"duration" json:"duration"`
	VideoCodec string  `graphql:"video_codec" json:"video_codec"`
	AudioCodec string  `graphql:"audio_codec" json:"audio_codec"`
	Width      int     `graphql:"width" json:"width"`
	Height     int     `graphql:"height" json:"height"`
	Framerate  float64 `graphql:"frame_rate" json:"frame_rate"`
	Bitrate    int     `graphql:"bit_rate" json:"bit_rate"`
}

func (f stashVideoFile) SceneFileType() models.SceneFileType {
	ret := models.SceneFileType{
		Duration:   &f.Duration,
		VideoCodec: &f.VideoCodec,
		AudioCodec: &f.AudioCodec,
		Width:      &f.Width,
		Height:     &f.Height,
		Framerate:  &f.Framerate,
		Bitrate:    &f.Bitrate,
	}

	size := strconv.FormatInt(f.Size, 10)
	ret.Size = &size

	return ret
}

type scrapedSceneStash struct {
	ID         string                   `graphql:"id" json:"id"`
	Title      *string                  `graphql:"title" json:"title"`
	Details    *string                  `graphql:"details" json:"details"`
	URLs       []string                 `graphql:"urls" json:"urls"`
	Date       *string                  `graphql:"date" json:"date"`
	Files      []stashVideoFile         `graphql:"files" json:"files"`
	Studio     *scrapedStudioStash      `graphql:"studio" json:"studio"`
	Tags       []*scrapedTagStash       `graphql:"tags" json:"tags"`
	Performers []*scrapedPerformerStash `graphql:"performers" json:"performers"`
}

func (s *stashScraper) scrapeSceneByScene(ctx context.Context, scene *models.Scene) (*models.ScrapedScene, error) {
	// query by MD5
	var q struct {
		FindScene *scrapedSceneStash `graphql:"findSceneByHash(input: $c)"`
	}

	type SceneHashInput struct {
		Checksum *string `graphql:"checksum" json:"checksum"`
		Oshash   *string `graphql:"oshash" json:"oshash"`
	}

	checksum := scene.Checksum
	oshash := scene.OSHash

	input := SceneHashInput{
		Checksum: &checksum,
		Oshash:   &oshash,
	}

	vars := map[string]interface{}{
		"c": input,
	}

	client := s.getStashClient()
	if err := client.Query(ctx, &q, vars); err != nil {
		return nil, convertGraphqlError(err)
	}

	if q.FindScene == nil {
		return nil, nil
	}

	// need to copy back to a scraped scene
	ret, err := s.scrapedStashSceneToScrapedScene(ctx, q.FindScene)
	if err != nil {
		return nil, err
	}

	// get the scene image directly
	ig := s.imageGetter()
	ret.Image, err = getStashSceneImage(ctx, s.config.StashServer.URL, q.FindScene.ID, ig)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

type scrapedGalleryStash struct {
	ID         string                   `graphql:"id" json:"id"`
	Title      *string                  `graphql:"title" json:"title"`
	Details    *string                  `graphql:"details" json:"details"`
	URL        *string                  `graphql:"url" json:"url"`
	Date       *string                  `graphql:"date" json:"date"`
	File       *models.SceneFileType    `graphql:"file" json:"file"`
	Studio     *scrapedStudioStash      `graphql:"studio" json:"studio"`
	Tags       []*scrapedTagStash       `graphql:"tags" json:"tags"`
	Performers []*scrapedPerformerStash `graphql:"performers" json:"performers"`
}

func (s *stashScraper) scrapeGalleryByGallery(ctx context.Context, gallery *models.Gallery) (*models.ScrapedGallery, error) {
	var q struct {
		FindGallery *scrapedGalleryStash `graphql:"findGalleryByHash(input: $c)"`
	}

	type GalleryHashInput struct {
		Checksum *string `graphql:"checksum" json:"checksum"`
	}

	checksum := gallery.PrimaryChecksum()
	input := GalleryHashInput{
		Checksum: &checksum,
	}

	vars := map[string]interface{}{
		"c": &input,
	}

	client := s.getStashClient()
	if err := client.Query(ctx, &q, vars); err != nil {
		return nil, err
	}

	// need to copy back to a scraped scene
	ret := models.ScrapedGallery{}
	if err := copier.Copy(&ret, q.FindGallery); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (s *stashScraper) scrapeImageByImage(ctx context.Context, image *models.Image) (*models.ScrapedImage, error) {
	return nil, ErrNotSupported
}

func (s *stashScraper) scrapeByURL(_ context.Context, _ string, _ ScrapeContentType) (ScrapedContent, error) {
	return nil, ErrNotSupported
}

func imageToUpdateInput(gallery *models.Image) models.ImageUpdateInput {
	dateToStringPtr := func(s *models.Date) *string {
		if s != nil {
			v := s.String()
			return &v
		}

		return nil
	}

	// fallback to file basename if title is empty
	title := gallery.GetTitle()
	urls := gallery.URLs.List()

	return models.ImageUpdateInput{
		ID:      strconv.Itoa(gallery.ID),
		Title:   &title,
		Details: &gallery.Details,
		Urls:    urls,
		Date:    dateToStringPtr(gallery.Date),
	}
}
