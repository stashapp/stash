package scraper

import (
	"context"
	"errors"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/shurcooL/graphql"

	"github.com/stashapp/stash/pkg/models"
)

type stashScraper struct {
	scraper      scraperTypeConfig
	config       config
	globalConfig GlobalConfig
}

func newStashScraper(scraper scraperTypeConfig, config config, globalConfig GlobalConfig) *stashScraper {
	return &stashScraper{
		scraper:      scraper,
		config:       config,
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

func (s *stashScraper) scrapePerformersByName(name string) ([]*models.ScrapedPerformer, error) {
	client := s.getStashClient()

	var q struct {
		FindPerformers stashFindPerformerNamesResultType `graphql:"findPerformers(filter: $f)"`
	}

	page := 1
	perPage := 10

	vars := map[string]interface{}{
		"f": models.FindFilterType{
			Q:       &name,
			Page:    &page,
			PerPage: &perPage,
		},
	}

	err := client.Query(context.Background(), &q, vars)
	if err != nil {
		return nil, err
	}

	var ret []*models.ScrapedPerformer
	for _, p := range q.FindPerformers.Performers {
		ret = append(ret, p.toPerformer())
	}

	return ret, nil
}

func (s *stashScraper) scrapePerformerByFragment(scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error) {
	client := s.getStashClient()

	var q struct {
		FindPerformer *models.ScrapedPerformerStash `graphql:"findPerformer(id: $f)"`
	}

	performerID := *scrapedPerformer.URL

	// get the id from the URL field
	vars := map[string]interface{}{
		"f": performerID,
	}

	err := client.Query(context.Background(), &q, vars)
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
	ret.Image, err = getStashPerformerImage(s.config.StashServer.URL, performerID, s.globalConfig)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (s *stashScraper) scrapeSceneByFragment(scene models.SceneUpdateInput) (*models.ScrapedScene, error) {
	// query by MD5
	// assumes that the scene exists in the database
	qb := models.NewSceneQueryBuilder()
	id, err := strconv.Atoi(scene.ID)
	if err != nil {
		return nil, err
	}

	storedScene, err := qb.Find(id)

	if err != nil {
		return nil, err
	}

	var q struct {
		FindScene *models.ScrapedSceneStash `graphql:"findSceneByHash(input: $c)"`
	}

	type SceneHashInput struct {
		Checksum *string `graphql:"checksum" json:"checksum"`
		Oshash   *string `graphql:"oshash" json:"oshash"`
	}

	input := SceneHashInput{
		Checksum: &storedScene.Checksum.String,
		Oshash:   &storedScene.OSHash.String,
	}

	vars := map[string]interface{}{
		"c": &input,
	}

	client := s.getStashClient()
	err = client.Query(context.Background(), &q, vars)
	if err != nil {
		return nil, err
	}

	if q.FindScene != nil {
		// the ids of the studio, performers and tags must be nilled
		if q.FindScene.Studio != nil {
			q.FindScene.Studio.ID = nil
		}

		for _, p := range q.FindScene.Performers {
			p.ID = nil
		}

		for _, t := range q.FindScene.Tags {
			t.ID = nil
		}
	}

	// need to copy back to a scraped scene
	ret := models.ScrapedScene{}
	err = copier.Copy(&ret, q.FindScene)
	if err != nil {
		return nil, err
	}

	// get the performer image directly
	ret.Image, err = getStashSceneImage(s.config.StashServer.URL, q.FindScene.ID, s.globalConfig)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (s *stashScraper) scrapePerformerByURL(url string) (*models.ScrapedPerformer, error) {
	return nil, errors.New("scrapePerformerByURL not supported for stash scraper")
}

func (s *stashScraper) scrapeSceneByURL(url string) (*models.ScrapedScene, error) {
	return nil, errors.New("scrapeSceneByURL not supported for stash scraper")
}

func (s *stashScraper) scrapeMovieByURL(url string) (*models.ScrapedMovie, error) {
	return nil, errors.New("scrapeMovieByURL not supported for stash scraper")
}

func sceneFromUpdateFragment(scene models.SceneUpdateInput) (*models.Scene, error) {
	qb := models.NewSceneQueryBuilder()
	id, err := strconv.Atoi(scene.ID)
	if err != nil {
		return nil, err
	}

	// TODO - should we modify it with the input?
	return qb.Find(id)
}
