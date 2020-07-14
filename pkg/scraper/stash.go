package scraper

import (
	"context"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/shurcooL/graphql"

	"github.com/stashapp/stash/pkg/models"
)

func getStashClient(c scraperTypeConfig) *graphql.Client {
	url := c.scraperConfig.StashServer.URL
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

func scrapePerformerNamesStash(c scraperTypeConfig, name string) ([]*models.ScrapedPerformer, error) {
	client := getStashClient(c)

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

func scrapePerformerFragmentStash(c scraperTypeConfig, scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error) {
	client := getStashClient(c)

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
	ret.Image, err = getStashPerformerImage(c.scraperConfig.StashServer.URL, performerID)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func scrapeSceneFragmentStash(c scraperTypeConfig, scene models.SceneUpdateInput) (*models.ScrapedScene, error) {
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

	// TODO - change findScene to accept multiple inputs
	var q struct {
		FindScene *models.ScrapedSceneStash `graphql:"findScene(checksum: $c)"`
	}

	checksum := graphql.String(storedScene.Checksum.String)
	vars := map[string]interface{}{
		"c": &checksum,
	}

	client := getStashClient(c)
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
	ret.Image, err = getStashSceneImage(c.scraperConfig.StashServer.URL, q.FindScene.ID)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}
