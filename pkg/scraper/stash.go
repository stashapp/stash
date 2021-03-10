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
	txnManager   models.TransactionManager
}

func newStashScraper(scraper scraperTypeConfig, txnManager models.TransactionManager, config config, globalConfig GlobalConfig) *stashScraper {
	return &stashScraper{
		scraper:      scraper,
		config:       config,
		globalConfig: globalConfig,
		txnManager:   txnManager,
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

	if q.FindPerformer != nil {
		// the ids of the tags must be nilled
		for _, t := range q.FindPerformer.Tags {
			t.ID = nil
		}
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
	id, err := strconv.Atoi(scene.ID)
	if err != nil {
		return nil, err
	}

	var storedScene *models.Scene
	if err := s.txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		var err error
		storedScene, err = r.Scene().Find(id)
		return err
	}); err != nil {
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

func (s *stashScraper) scrapeGalleryByFragment(scene models.GalleryUpdateInput) (*models.ScrapedGallery, error) {
	id, err := strconv.Atoi(scene.ID)
	if err != nil {
		return nil, err
	}

	// query by MD5
	// assumes that the gallery exists in the database
	var storedGallery *models.Gallery
	if err := s.txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		qb := r.Gallery()

		var err error
		storedGallery, err = qb.Find(id)
		return err
	}); err != nil {
		return nil, err
	}

	var q struct {
		FindGallery *models.ScrapedGalleryStash `graphql:"findGalleryByHash(input: $c)"`
	}

	type GalleryHashInput struct {
		Checksum *string `graphql:"checksum" json:"checksum"`
	}

	input := GalleryHashInput{
		Checksum: &storedGallery.Checksum,
	}

	vars := map[string]interface{}{
		"c": &input,
	}

	client := s.getStashClient()
	err = client.Query(context.Background(), &q, vars)
	if err != nil {
		return nil, err
	}

	if q.FindGallery != nil {
		// the ids of the studio, performers and tags must be nilled
		if q.FindGallery.Studio != nil {
			q.FindGallery.Studio.ID = nil
		}

		for _, p := range q.FindGallery.Performers {
			p.ID = nil
		}

		for _, t := range q.FindGallery.Tags {
			t.ID = nil
		}
	}

	// need to copy back to a scraped scene
	ret := models.ScrapedGallery{}
	err = copier.Copy(&ret, q.FindGallery)
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

func (s *stashScraper) scrapeGalleryByURL(url string) (*models.ScrapedGallery, error) {
	return nil, errors.New("scrapeGalleryByURL not supported for stash scraper")
}

func (s *stashScraper) scrapeMovieByURL(url string) (*models.ScrapedMovie, error) {
	return nil, errors.New("scrapeMovieByURL not supported for stash scraper")
}

func sceneFromUpdateFragment(scene models.SceneUpdateInput, txnManager models.TransactionManager) (*models.Scene, error) {
	id, err := strconv.Atoi(scene.ID)
	if err != nil {
		return nil, err
	}

	// TODO - should we modify it with the input?
	var ret *models.Scene
	if err := txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		var err error
		ret, err = r.Scene().Find(id)
		return err
	}); err != nil {
		return nil, err
	}
	return ret, nil
}

func galleryFromUpdateFragment(gallery models.GalleryUpdateInput, txnManager models.TransactionManager) (ret *models.Gallery, err error) {
	id, err := strconv.Atoi(gallery.ID)
	if err != nil {
		return nil, err
	}

	if err := txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		ret, err = r.Gallery().Find(id)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
