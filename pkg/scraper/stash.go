package scraper

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/shurcooL/graphql"

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

func (s *stashScraper) getStashClient() *graphql.Client {
	url := s.config.StashServer.URL
	return graphql.NewClient(url+"/graphql", nil)
}

type stashFindPerformerNamePerformer struct {
	ID   string `json:"id" graphql:"id"`
	Name string `json:"name" graphql:"name"`
}

func (p stashFindPerformerNamePerformer) toPerformer() *models.ScrapedPerformer {
	return &models.ScrapedPerformer{
		Name: &p.Name,
		// put id into the URL field
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
	URL          *string            `graphql:"url" json:"url"`
	Twitter      *string            `graphql:"twitter" json:"twitter"`
	Instagram    *string            `graphql:"instagram" json:"instagram"`
	Birthdate    *string            `graphql:"birthdate" json:"birthdate"`
	Ethnicity    *string            `graphql:"ethnicity" json:"ethnicity"`
	Country      *string            `graphql:"country" json:"country"`
	EyeColor     *string            `graphql:"eye_color" json:"eye_color"`
	Height       *string            `graphql:"height" json:"height"`
	Measurements *string            `graphql:"measurements" json:"measurements"`
	FakeTits     *string            `graphql:"fake_tits" json:"fake_tits"`
	PenisLength  *string            `graphql:"penis_length" json:"penis_length"`
	Circumcised  *string            `graphql:"circumcised" json:"circumcised"`
	CareerLength *string            `graphql:"career_length" json:"career_length"`
	Tattoos      *string            `graphql:"tattoos" json:"tattoos"`
	Piercings    *string            `graphql:"piercings" json:"piercings"`
	Aliases      *string            `graphql:"aliases" json:"aliases"`
	Tags         []*scrapedTagStash `graphql:"tags" json:"tags"`
	Details      *string            `graphql:"details" json:"details"`
	DeathDate    *string            `graphql:"death_date" json:"death_date"`
	HairColor    *string            `graphql:"hair_color" json:"hair_color"`
	Weight       *string            `graphql:"weight" json:"weight"`
}

func (s *stashScraper) scrapeByFragment(ctx context.Context, input Input) (ScrapedContent, error) {
	if input.Gallery != nil || input.Scene != nil {
		return nil, fmt.Errorf("%w: using stash scraper as a fragment scraper", ErrNotSupported)
	}

	if input.Performer == nil {
		return nil, fmt.Errorf("%w: the given performer is nil", ErrNotSupported)
	}

	scrapedPerformer := input.Performer

	client := s.getStashClient()

	var q struct {
		FindPerformer *scrapedPerformerStash `graphql:"findPerformer(id: $f)"`
	}

	performerID := *scrapedPerformer.URL

	// get the id from the URL field
	vars := map[string]interface{}{
		"f": performerID,
	}

	err := client.Query(ctx, &q, vars)
	if err != nil {
		return nil, err
	}

	// need to copy back to a scraped performer
	ret := models.ScrapedPerformer{}
	err = copier.Copy(&ret, q.FindPerformer)
	if err != nil {
		return nil, err
	}

	// get the performer image directly
	ret.Image, err = getStashPerformerImage(ctx, s.config.StashServer.URL, performerID, s.client, s.globalConfig)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

type scrapedStudioStash struct {
	Name string  `graphql:"name" json:"name"`
	URL  *string `graphql:"url" json:"url"`
}

type stashFindSceneNamesResultType struct {
	Count  int                  `graphql:"count"`
	Scenes []*scrapedSceneStash `graphql:"scenes"`
}

func (s *stashScraper) scrapedStashSceneToScrapedScene(ctx context.Context, scene *scrapedSceneStash) (*ScrapedScene, error) {
	ret := ScrapedScene{}
	err := copier.Copy(&ret, scene)
	if err != nil {
		return nil, err
	}

	// get the performer image directly
	ret.Image, err = getStashSceneImage(ctx, s.config.StashServer.URL, scene.ID, s.client, s.globalConfig)
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
			return nil, err
		}

		for _, scene := range q.FindScenes.Scenes {
			converted, err := s.scrapedStashSceneToScrapedScene(ctx, scene)
			if err != nil {
				return nil, err
			}
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

type scrapedSceneStash struct {
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

func (s *stashScraper) scrapeSceneByScene(ctx context.Context, scene *models.Scene) (*ScrapedScene, error) {
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
		"c": &input,
	}

	client := s.getStashClient()
	if err := client.Query(ctx, &q, vars); err != nil {
		return nil, err
	}

	// need to copy back to a scraped scene
	ret, err := s.scrapedStashSceneToScrapedScene(ctx, q.FindScene)
	if err != nil {
		return nil, err
	}

	// get the performer image directly
	ret.Image, err = getStashSceneImage(ctx, s.config.StashServer.URL, q.FindScene.ID, s.client, s.globalConfig)
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

func (s *stashScraper) scrapeGalleryByGallery(ctx context.Context, gallery *models.Gallery) (*ScrapedGallery, error) {
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
	ret := ScrapedGallery{}
	if err := copier.Copy(&ret, q.FindGallery); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (s *stashScraper) scrapeByURL(_ context.Context, _ string, _ ScrapeContentType) (ScrapedContent, error) {
	return nil, ErrNotSupported
}
