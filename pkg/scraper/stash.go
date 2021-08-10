package scraper

import (
	"context"
	"database/sql"
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

func (s *stashScraper) scrapePerformerByFragment(scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error) {
	client := s.getStashClient()

	var q struct {
		FindPerformer *scrapedPerformerStash `graphql:"findPerformer(id: $f)"`
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

type scrapedStudioStash struct {
	Name string  `graphql:"name" json:"name"`
	URL  *string `graphql:"url" json:"url"`
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

func (s *stashScraper) scrapeSceneByScene(scene *models.Scene) (*models.ScrapedScene, error) {
	// query by MD5
	var q struct {
		FindScene *scrapedSceneStash `graphql:"findSceneByHash(input: $c)"`
	}

	type SceneHashInput struct {
		Checksum *string `graphql:"checksum" json:"checksum"`
		Oshash   *string `graphql:"oshash" json:"oshash"`
	}

	input := SceneHashInput{
		Checksum: &scene.Checksum.String,
		Oshash:   &scene.OSHash.String,
	}

	vars := map[string]interface{}{
		"c": &input,
	}

	client := s.getStashClient()
	if err := client.Query(context.Background(), &q, vars); err != nil {
		return nil, err
	}

	// need to copy back to a scraped scene
	ret := models.ScrapedScene{}
	if err := copier.Copy(&ret, q.FindScene); err != nil {
		return nil, err
	}

	// get the performer image directly
	var err error
	ret.Image, err = getStashSceneImage(s.config.StashServer.URL, q.FindScene.ID, s.globalConfig)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (s *stashScraper) scrapeSceneByFragment(scene models.ScrapedSceneInput) (*models.ScrapedScene, error) {
	return nil, errors.New("scrapeSceneByFragment not supported for stash scraper")
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

func (s *stashScraper) scrapeGalleryByGallery(gallery *models.Gallery) (*models.ScrapedGallery, error) {
	var q struct {
		FindGallery *scrapedGalleryStash `graphql:"findGalleryByHash(input: $c)"`
	}

	type GalleryHashInput struct {
		Checksum *string `graphql:"checksum" json:"checksum"`
	}

	input := GalleryHashInput{
		Checksum: &gallery.Checksum,
	}

	vars := map[string]interface{}{
		"c": &input,
	}

	client := s.getStashClient()
	if err := client.Query(context.Background(), &q, vars); err != nil {
		return nil, err
	}

	// need to copy back to a scraped scene
	ret := models.ScrapedGallery{}
	if err := copier.Copy(&ret, q.FindGallery); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (s *stashScraper) scrapeGalleryByFragment(scene models.ScrapedGalleryInput) (*models.ScrapedGallery, error) {
	return nil, errors.New("scrapeGalleryByFragment not supported for stash scraper")
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

func getScene(sceneID int, txnManager models.TransactionManager) (*models.Scene, error) {
	var ret *models.Scene
	if err := txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		var err error
		ret, err = r.Scene().Find(sceneID)
		return err
	}); err != nil {
		return nil, err
	}
	return ret, nil
}

func sceneToUpdateInput(scene *models.Scene) models.SceneUpdateInput {
	toStringPtr := func(s sql.NullString) *string {
		if s.Valid {
			return &s.String
		}

		return nil
	}

	dateToStringPtr := func(s models.SQLiteDate) *string {
		if s.Valid {
			return &s.String
		}

		return nil
	}

	return models.SceneUpdateInput{
		ID:      strconv.Itoa(scene.ID),
		Title:   toStringPtr(scene.Title),
		Details: toStringPtr(scene.Details),
		URL:     toStringPtr(scene.URL),
		Date:    dateToStringPtr(scene.Date),
	}
}

func getGallery(galleryID int, txnManager models.TransactionManager) (*models.Gallery, error) {
	var ret *models.Gallery
	if err := txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		var err error
		ret, err = r.Gallery().Find(galleryID)
		return err
	}); err != nil {
		return nil, err
	}
	return ret, nil
}

func galleryToUpdateInput(gallery *models.Gallery) models.GalleryUpdateInput {
	toStringPtr := func(s sql.NullString) *string {
		if s.Valid {
			return &s.String
		}

		return nil
	}

	dateToStringPtr := func(s models.SQLiteDate) *string {
		if s.Valid {
			return &s.String
		}

		return nil
	}

	return models.GalleryUpdateInput{
		ID:      strconv.Itoa(gallery.ID),
		Title:   toStringPtr(gallery.Title),
		Details: toStringPtr(gallery.Details),
		URL:     toStringPtr(gallery.URL),
		Date:    dateToStringPtr(gallery.Date),
	}
}
