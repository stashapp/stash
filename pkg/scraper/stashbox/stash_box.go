package stashbox

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Yamashou/gqlgenc/client"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper/stashbox/graphql"
)

type Client struct {
	client *graphql.Client
}

// NewClient returns a new instance of a stash-box client.
func NewClient(box models.StashBox) *Client {
	authHeader := func(req *http.Request) {
		req.Header.Set("ApiKey", box.APIKey)
	}

	client := &graphql.Client{
		Client: client.NewClient(http.DefaultClient, box.Endpoint, authHeader),
	}

	return &Client{
		client: client,
	}
}

func (c Client) QueryStashBoxScene(queryStr string) ([]*models.ScrapedScene, error) {
	scenes, err := c.client.SearchScene(context.TODO(), queryStr)
	if err != nil {
		return nil, err
	}

	sceneFragments := scenes.SearchScene

	var ret []*models.ScrapedScene
	for _, s := range sceneFragments {
		ss, err := sceneFragmentToScrapedScene(s)
		if err != nil {
			return nil, err
		}
		ret = append(ret, ss)
	}

	return ret, nil
}

func (c Client) FindStashBoxSceneByFingerprint(sceneID int) ([]*models.ScrapedScene, error) {
	// find the scene hash
	qb := models.NewSceneQueryBuilder()
	scene, err := qb.Find(sceneID)
	if err != nil {
		return nil, err
	}

	// try MD5 first, otherwise try oshash
	var ret []*models.ScrapedScene
	if scene.Checksum.Valid {
		ret, err = c.findStashBoxSceneByFingerprint(scene.Checksum.String, graphql.FingerprintAlgorithmMd5)
		if err != nil {
			return nil, err
		}
	}

	if ret == nil && scene.OSHash.Valid {
		ret, err = c.findStashBoxSceneByFingerprint(scene.OSHash.String, graphql.FingerprintAlgorithmOshash)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func (c Client) findStashBoxSceneByFingerprint(hash string, algoType graphql.FingerprintAlgorithm) ([]*models.ScrapedScene, error) {
	scenes, err := c.client.FindSceneByFingerprint(context.TODO(), graphql.FingerprintQueryInput{
		Hash:      hash,
		Algorithm: algoType,
	})

	if err != nil {
		return nil, err
	}

	sceneFragments := scenes.FindSceneByFingerprint

	var ret []*models.ScrapedScene
	for _, s := range sceneFragments {
		ss, err := sceneFragmentToScrapedScene(s)
		if err != nil {
			return nil, err
		}
		ret = append(ret, ss)
	}

	return ret, nil
}

func findURL(urls []*graphql.URLFragment, urlType string) *string {
	for _, u := range urls {
		if u.Type == urlType {
			ret := u.URL
			return &ret
		}
	}

	return nil
}

func enumToStringPtr(e fmt.Stringer) *string {
	if e != nil {
		ret := e.String()
		return &ret
	}

	return nil
}

func formatMeasurements(m graphql.MeasurementsFragment) *string {
	if m.BandSize != nil && m.CupSize != nil && m.Hip != nil && m.Waist != nil {
		ret := fmt.Sprintf("%d%s-%d-%d", *m.BandSize, *m.CupSize, *m.Waist, *m.Hip)
		return &ret
	}

	return nil
}

func formatCareerLength(start, end *int) *string {
	if start == nil && end == nil {
		return nil
	}

	var ret string
	if end == nil {
		ret = fmt.Sprintf("%d -", *start)
	} else {
		ret = fmt.Sprintf("%d - %d", *start, *end)
	}

	return &ret
}

func formatBodyModifications(m []*graphql.BodyModificationFragment) *string {
	if len(m) == 0 {
		return nil
	}

	var retSlice []string
	for _, f := range m {
		if f.Description == nil {
			retSlice = append(retSlice, f.Location)
		} else {
			retSlice = append(retSlice, fmt.Sprintf("%s, %s", f.Location, *f.Description))
		}
	}

	ret := strings.Join(retSlice, "; ")
	return &ret
}

func performerFragmentToScrapedScenePerformer(p graphql.PerformerFragment) *models.ScrapedScenePerformer {
	sp := &models.ScrapedScenePerformer{
		Name:         p.Name,
		Country:      p.Country,
		Measurements: formatMeasurements(p.Measurements),
		CareerLength: formatCareerLength(p.CareerStartYear, p.CareerEndYear),
		Tattoos:      formatBodyModifications(p.Tattoos),
		Piercings:    formatBodyModifications(p.Piercings),
		Twitter:      findURL(p.Urls, "TWITTER"),
		// TODO - Image
	}

	if p.Height != nil {
		hs := strconv.Itoa(*p.Height)
		sp.Height = &hs
	}

	if p.Birthdate != nil {
		b := p.Birthdate.Date
		sp.Birthdate = &b
	}

	if p.Gender != nil {
		sp.Gender = enumToStringPtr(p.Gender)
	}

	if p.Ethnicity != nil {
		sp.Ethnicity = enumToStringPtr(p.Ethnicity)
	}

	if p.EyeColor != nil {
		sp.EyeColor = enumToStringPtr(p.EyeColor)
	}

	if p.BreastType != nil {
		sp.FakeTits = enumToStringPtr(p.BreastType)
	}

	return sp
}

func sceneFragmentToScrapedScene(s *graphql.SceneFragment) (*models.ScrapedScene, error) {
	ss := &models.ScrapedScene{
		Title:   s.Title,
		Date:    s.Date,
		Details: s.Details,
		URL:     findURL(s.Urls, "STUDIO"),
		// Image
		// stash_id
	}

	if s.Studio != nil {
		ss.Studio = &models.ScrapedSceneStudio{
			Name: s.Studio.Name,
			URL:  findURL(s.Studio.Urls, "HOME"),
		}

		err := models.MatchScrapedSceneStudio(ss.Studio)
		if err != nil {
			return nil, err
		}
	}

	for _, p := range s.Performers {
		sp := performerFragmentToScrapedScenePerformer(p.Performer)

		err := models.MatchScrapedScenePerformer(sp)
		if err != nil {
			return nil, err
		}

		ss.Performers = append(ss.Performers, sp)
	}

	for _, t := range s.Tags {
		st := &models.ScrapedSceneTag{
			Name: t.Name,
		}

		err := models.MatchScrapedSceneTag(st)
		if err != nil {
			return nil, err
		}

		ss.Tags = append(ss.Tags, st)
	}

	return ss, nil
}
