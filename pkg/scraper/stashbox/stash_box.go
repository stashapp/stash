package stashbox

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Yamashou/gqlgenc/client"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stashapp/stash/pkg/scraper/stashbox/graphql"
	"github.com/stashapp/stash/pkg/utils"
)

// Timeout to get the image. Includes transfer time. May want to make this
// configurable at some point.
const imageGetTimeout = time.Second * 30

// Client represents the client interface to a stash-box server instance.
type Client struct {
	client     *graphql.Client
	txnManager models.TransactionManager
}

// NewClient returns a new instance of a stash-box client.
func NewClient(box models.StashBox, txnManager models.TransactionManager) *Client {
	authHeader := func(req *http.Request) {
		req.Header.Set("ApiKey", box.APIKey)
	}

	client := &graphql.Client{
		Client: client.NewClient(http.DefaultClient, box.Endpoint, authHeader),
	}

	return &Client{
		client:     client,
		txnManager: txnManager,
	}
}

// QueryStashBoxScene queries stash-box for scenes using a query string.
func (c Client) QueryStashBoxScene(queryStr string) ([]*models.ScrapedScene, error) {
	scenes, err := c.client.SearchScene(context.TODO(), queryStr)
	if err != nil {
		return nil, err
	}

	sceneFragments := scenes.SearchScene

	var ret []*models.ScrapedScene
	for _, s := range sceneFragments {
		ss, err := sceneFragmentToScrapedScene(c.txnManager, s)
		if err != nil {
			return nil, err
		}
		ret = append(ret, ss)
	}

	return ret, nil
}

// FindStashBoxScenesByFingerprints queries stash-box for scenes using every
// scene's MD5 checksum and/or oshash.
func (c Client) FindStashBoxScenesByFingerprints(sceneIDs []string) ([]*models.ScrapedScene, error) {
	ids, err := utils.StringSliceToIntSlice(sceneIDs)
	if err != nil {
		return nil, err
	}

	var fingerprints []string

	if err := c.txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		qb := r.Scene()

		for _, sceneID := range ids {
			scene, err := qb.Find(sceneID)
			if err != nil {
				return err
			}

			if scene == nil {
				return fmt.Errorf("scene with id %d not found", sceneID)
			}

			if scene.Checksum.Valid {
				fingerprints = append(fingerprints, scene.Checksum.String)
			}

			if scene.OSHash.Valid {
				fingerprints = append(fingerprints, scene.OSHash.String)
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return c.findStashBoxScenesByFingerprints(fingerprints)
}

func (c Client) findStashBoxScenesByFingerprints(fingerprints []string) ([]*models.ScrapedScene, error) {
	var ret []*models.ScrapedScene
	for i := 0; i < len(fingerprints); i += 100 {
		end := i + 100
		if end > len(fingerprints) {
			end = len(fingerprints)
		}
		scenes, err := c.client.FindScenesByFingerprints(context.TODO(), fingerprints[i:end])

		if err != nil {
			return nil, err
		}

		sceneFragments := scenes.FindScenesByFingerprints

		for _, s := range sceneFragments {
			ss, err := sceneFragmentToScrapedScene(c.txnManager, s)
			if err != nil {
				return nil, err
			}
			ret = append(ret, ss)
		}
	}

	return ret, nil
}

func (c Client) SubmitStashBoxFingerprints(sceneIDs []string, endpoint string) (bool, error) {
	ids, err := utils.StringSliceToIntSlice(sceneIDs)
	if err != nil {
		return false, err
	}

	var fingerprints []graphql.FingerprintSubmission

	if err := c.txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		qb := r.Scene()

		for _, sceneID := range ids {
			scene, err := qb.Find(sceneID)
			if err != nil {
				return err
			}

			if scene == nil {
				continue
			}

			stashIDs, err := qb.GetStashIDs(sceneID)
			if err != nil {
				return err
			}

			sceneStashID := ""
			for _, stashID := range stashIDs {
				if stashID.Endpoint == endpoint {
					sceneStashID = stashID.StashID
				}
			}

			if sceneStashID != "" {
				if scene.Checksum.Valid && scene.Duration.Valid {
					fingerprint := graphql.FingerprintInput{
						Hash:      scene.Checksum.String,
						Algorithm: graphql.FingerprintAlgorithmMd5,
						Duration:  int(scene.Duration.Float64),
					}
					fingerprints = append(fingerprints, graphql.FingerprintSubmission{
						SceneID:     sceneStashID,
						Fingerprint: &fingerprint,
					})
				}

				if scene.OSHash.Valid && scene.Duration.Valid {
					fingerprint := graphql.FingerprintInput{
						Hash:      scene.OSHash.String,
						Algorithm: graphql.FingerprintAlgorithmOshash,
						Duration:  int(scene.Duration.Float64),
					}
					fingerprints = append(fingerprints, graphql.FingerprintSubmission{
						SceneID:     sceneStashID,
						Fingerprint: &fingerprint,
					})
				}
			}
		}

		return nil
	}); err != nil {
		return false, err
	}

	return c.submitStashBoxFingerprints(fingerprints)
}

func (c Client) submitStashBoxFingerprints(fingerprints []graphql.FingerprintSubmission) (bool, error) {
	for _, fingerprint := range fingerprints {
		_, err := c.client.SubmitFingerprint(context.TODO(), fingerprint)
		if err != nil {
			return false, err
		}
	}

	return true, nil
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

func fetchImage(url string) (*string, error) {
	client := &http.Client{
		Timeout: imageGetTimeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// determine the image type and set the base64 type
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(body)
	}

	img := "data:" + contentType + ";base64," + utils.GetBase64StringFromData(body)
	return &img, nil
}

func performerFragmentToScrapedScenePerformer(p graphql.PerformerFragment) *models.ScrapedScenePerformer {
	id := p.ID
	images := []string{}
	for _, image := range p.Images {
		images = append(images, image.URL)
	}
	sp := &models.ScrapedScenePerformer{
		Name:         p.Name,
		Country:      p.Country,
		Measurements: formatMeasurements(p.Measurements),
		CareerLength: formatCareerLength(p.CareerStartYear, p.CareerEndYear),
		Tattoos:      formatBodyModifications(p.Tattoos),
		Piercings:    formatBodyModifications(p.Piercings),
		Twitter:      findURL(p.Urls, "TWITTER"),
		RemoteSiteID: &id,
		Images:       images,
		// TODO - tags not currently supported
		// TODO - Image - should be returned as a set of URLs. Will need a
		// graphql schema change to accommodate this. Leave off for now.
	}

	if p.Height != nil && *p.Height > 0 {
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

func getFirstImage(images []*graphql.ImageFragment) *string {
	ret, err := fetchImage(images[0].URL)
	if err != nil {
		logger.Warnf("Error fetching image %s: %s", images[0].URL, err.Error())
	}

	return ret
}

func getFingerprints(scene *graphql.SceneFragment) []*models.StashBoxFingerprint {
	fingerprints := []*models.StashBoxFingerprint{}
	for _, fp := range scene.Fingerprints {
		fingerprint := models.StashBoxFingerprint{
			Algorithm: fp.Algorithm.String(),
			Hash:      fp.Hash,
			Duration:  fp.Duration,
		}
		fingerprints = append(fingerprints, &fingerprint)
	}
	return fingerprints
}

func sceneFragmentToScrapedScene(txnManager models.TransactionManager, s *graphql.SceneFragment) (*models.ScrapedScene, error) {
	stashID := s.ID
	ss := &models.ScrapedScene{
		Title:        s.Title,
		Date:         s.Date,
		Details:      s.Details,
		URL:          findURL(s.Urls, "STUDIO"),
		Duration:     s.Duration,
		RemoteSiteID: &stashID,
		Fingerprints: getFingerprints(s),
		// Image
		// stash_id
	}

	if len(s.Images) > 0 {
		// TODO - #454 code sorts images by aspect ratio according to a wanted
		// orientation. I'm just grabbing the first for now
		ss.Image = getFirstImage(s.Images)
	}

	if err := txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		pqb := r.Performer()
		tqb := r.Tag()

		if s.Studio != nil {
			studioID := s.Studio.ID
			ss.Studio = &models.ScrapedSceneStudio{
				Name:         s.Studio.Name,
				URL:          findURL(s.Studio.Urls, "HOME"),
				RemoteSiteID: &studioID,
			}

			err := scraper.MatchScrapedSceneStudio(r.Studio(), ss.Studio)
			if err != nil {
				return err
			}
		}

		for _, p := range s.Performers {
			sp := performerFragmentToScrapedScenePerformer(p.Performer)

			err := scraper.MatchScrapedScenePerformer(pqb, sp)
			if err != nil {
				return err
			}

			ss.Performers = append(ss.Performers, sp)
		}

		for _, t := range s.Tags {
			st := &models.ScrapedSceneTag{
				Name: t.Name,
			}

			err := scraper.MatchScrapedSceneTag(tqb, st)
			if err != nil {
				return err
			}

			ss.Tags = append(ss.Tags, st)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ss, nil
}
