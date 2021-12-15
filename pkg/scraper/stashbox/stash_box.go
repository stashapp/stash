package stashbox

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Yamashou/gqlgenc/client"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper/stashbox/graphql"
	"github.com/stashapp/stash/pkg/utils"
)

// Client represents the client interface to a stash-box server instance.
type Client struct {
	client     *graphql.Client
	txnManager models.TransactionManager
	box        models.StashBox
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
		box:        box,
	}
}

func (c Client) getHTTPClient() *http.Client {
	return c.client.Client.Client
}

// QueryStashBoxScene queries stash-box for scenes using a query string.
func (c Client) QueryStashBoxScene(ctx context.Context, queryStr string) ([]*models.ScrapedScene, error) {
	scenes, err := c.client.SearchScene(ctx, queryStr)
	if err != nil {
		return nil, err
	}

	sceneFragments := scenes.SearchScene

	var ret []*models.ScrapedScene
	for _, s := range sceneFragments {
		ss, err := c.sceneFragmentToScrapedScene(ctx, s)
		if err != nil {
			return nil, err
		}
		ret = append(ret, ss)
	}

	return ret, nil
}

// FindStashBoxScenesByFingerprints queries stash-box for scenes using every
// scene's MD5/OSHASH checksum, or PHash, and returns results in the same order
// as the input slice.
func (c Client) FindStashBoxScenesByFingerprints(ctx context.Context, sceneIDs []string) ([][]*models.ScrapedScene, error) {
	ids, err := utils.StringSliceToIntSlice(sceneIDs)
	if err != nil {
		return nil, err
	}

	var fingerprints []*graphql.FingerprintQueryInput
	// map fingerprints to their scene index
	fpToScene := make(map[string][]int)

	if err := c.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		qb := r.Scene()

		for index, sceneID := range ids {
			scene, err := qb.Find(sceneID)
			if err != nil {
				return err
			}

			if scene == nil {
				return fmt.Errorf("scene with id %d not found", sceneID)
			}

			if scene.Checksum.Valid {
				fingerprints = append(fingerprints, &graphql.FingerprintQueryInput{
					Hash:      scene.Checksum.String,
					Algorithm: graphql.FingerprintAlgorithmMd5,
				})
				fpToScene[scene.Checksum.String] = append(fpToScene[scene.Checksum.String], index)
			}

			if scene.OSHash.Valid {
				fingerprints = append(fingerprints, &graphql.FingerprintQueryInput{
					Hash:      scene.OSHash.String,
					Algorithm: graphql.FingerprintAlgorithmOshash,
				})
				fpToScene[scene.OSHash.String] = append(fpToScene[scene.OSHash.String], index)
			}

			if scene.Phash.Valid {
				phashStr := utils.PhashToString(scene.Phash.Int64)
				fingerprints = append(fingerprints, &graphql.FingerprintQueryInput{
					Hash:      phashStr,
					Algorithm: graphql.FingerprintAlgorithmPhash,
				})
				fpToScene[phashStr] = append(fpToScene[phashStr], index)
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	allScenes, err := c.findStashBoxScenesByFingerprints(ctx, fingerprints)
	if err != nil {
		return nil, err
	}

	// set the matched scenes back in their original order
	ret := make([][]*models.ScrapedScene, len(sceneIDs))
	for _, s := range allScenes {
		var addedTo []int
		for _, fp := range s.Fingerprints {
			sceneIndexes := fpToScene[fp.Hash]
			for _, index := range sceneIndexes {
				if !utils.IntInclude(addedTo, index) {
					addedTo = append(addedTo, index)
					ret[index] = append(ret[index], s)
				}
			}
		}
	}

	return ret, nil
}

// FindStashBoxScenesByFingerprintsFlat queries stash-box for scenes using every
// scene's MD5/OSHASH checksum, or PHash, and returns results a flat slice.
func (c Client) FindStashBoxScenesByFingerprintsFlat(ctx context.Context, sceneIDs []string) ([]*models.ScrapedScene, error) {
	ids, err := utils.StringSliceToIntSlice(sceneIDs)
	if err != nil {
		return nil, err
	}

	var fingerprints []*graphql.FingerprintQueryInput

	if err := c.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
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
				fingerprints = append(fingerprints, &graphql.FingerprintQueryInput{
					Hash:      scene.Checksum.String,
					Algorithm: graphql.FingerprintAlgorithmMd5,
				})
			}

			if scene.OSHash.Valid {
				fingerprints = append(fingerprints, &graphql.FingerprintQueryInput{
					Hash:      scene.OSHash.String,
					Algorithm: graphql.FingerprintAlgorithmOshash,
				})
			}

			if scene.Phash.Valid {
				fingerprints = append(fingerprints, &graphql.FingerprintQueryInput{
					Hash:      utils.PhashToString(scene.Phash.Int64),
					Algorithm: graphql.FingerprintAlgorithmPhash,
				})
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return c.findStashBoxScenesByFingerprints(ctx, fingerprints)
}

func (c Client) findStashBoxScenesByFingerprints(ctx context.Context, fingerprints []*graphql.FingerprintQueryInput) ([]*models.ScrapedScene, error) {
	var ret []*models.ScrapedScene
	for i := 0; i < len(fingerprints); i += 100 {
		end := i + 100
		if end > len(fingerprints) {
			end = len(fingerprints)
		}
		scenes, err := c.client.FindScenesByFullFingerprints(ctx, fingerprints[i:end])

		if err != nil {
			return nil, err
		}

		sceneFragments := scenes.FindScenesByFullFingerprints

		for _, s := range sceneFragments {
			ss, err := c.sceneFragmentToScrapedScene(ctx, s)
			if err != nil {
				return nil, err
			}
			ret = append(ret, ss)
		}
	}

	return ret, nil
}

func (c Client) SubmitStashBoxFingerprints(ctx context.Context, sceneIDs []string, endpoint string) (bool, error) {
	ids, err := utils.StringSliceToIntSlice(sceneIDs)
	if err != nil {
		return false, err
	}

	var fingerprints []graphql.FingerprintSubmission

	if err := c.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
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

				if scene.Phash.Valid && scene.Duration.Valid {
					fingerprint := graphql.FingerprintInput{
						Hash:      utils.PhashToString(scene.Phash.Int64),
						Algorithm: graphql.FingerprintAlgorithmPhash,
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

	return c.submitStashBoxFingerprints(ctx, fingerprints)
}

func (c Client) submitStashBoxFingerprints(ctx context.Context, fingerprints []graphql.FingerprintSubmission) (bool, error) {
	for _, fingerprint := range fingerprints {
		_, err := c.client.SubmitFingerprint(ctx, fingerprint)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

// QueryStashBoxPerformer queries stash-box for performers using a query string.
func (c Client) QueryStashBoxPerformer(ctx context.Context, queryStr string) ([]*models.StashBoxPerformerQueryResult, error) {
	performers, err := c.queryStashBoxPerformer(ctx, queryStr)

	res := []*models.StashBoxPerformerQueryResult{
		{
			Query:   queryStr,
			Results: performers,
		},
	}

	// set the deprecated image field
	for _, p := range res[0].Results {
		if len(p.Images) > 0 {
			p.Image = &p.Images[0]
		}
	}

	return res, err
}

func (c Client) queryStashBoxPerformer(ctx context.Context, queryStr string) ([]*models.ScrapedPerformer, error) {
	performers, err := c.client.SearchPerformer(ctx, queryStr)
	if err != nil {
		return nil, err
	}

	performerFragments := performers.SearchPerformer

	var ret []*models.ScrapedPerformer
	for _, fragment := range performerFragments {
		performer := performerFragmentToScrapedScenePerformer(*fragment)
		ret = append(ret, performer)
	}

	return ret, nil
}

// FindStashBoxPerformersByNames queries stash-box for performers by name
func (c Client) FindStashBoxPerformersByNames(ctx context.Context, performerIDs []string) ([]*models.StashBoxPerformerQueryResult, error) {
	ids, err := utils.StringSliceToIntSlice(performerIDs)
	if err != nil {
		return nil, err
	}

	var performers []*models.Performer

	if err := c.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		qb := r.Performer()

		for _, performerID := range ids {
			performer, err := qb.Find(performerID)
			if err != nil {
				return err
			}

			if performer == nil {
				return fmt.Errorf("performer with id %d not found", performerID)
			}

			if performer.Name.Valid {
				performers = append(performers, performer)
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return c.findStashBoxPerformersByNames(ctx, performers)
}

func (c Client) FindStashBoxPerformersByPerformerNames(ctx context.Context, performerIDs []string) ([][]*models.ScrapedPerformer, error) {
	ids, err := utils.StringSliceToIntSlice(performerIDs)
	if err != nil {
		return nil, err
	}

	var performers []*models.Performer

	if err := c.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		qb := r.Performer()

		for _, performerID := range ids {
			performer, err := qb.Find(performerID)
			if err != nil {
				return err
			}

			if performer == nil {
				return fmt.Errorf("performer with id %d not found", performerID)
			}

			if performer.Name.Valid {
				performers = append(performers, performer)
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	results, err := c.findStashBoxPerformersByNames(ctx, performers)
	if err != nil {
		return nil, err
	}

	var ret [][]*models.ScrapedPerformer
	for _, r := range results {
		ret = append(ret, r.Results)
	}

	return ret, nil
}

func (c Client) findStashBoxPerformersByNames(ctx context.Context, performers []*models.Performer) ([]*models.StashBoxPerformerQueryResult, error) {
	var ret []*models.StashBoxPerformerQueryResult
	for _, performer := range performers {
		if performer.Name.Valid {
			performerResults, err := c.queryStashBoxPerformer(ctx, performer.Name.String)
			if err != nil {
				return nil, err
			}

			result := models.StashBoxPerformerQueryResult{
				Query:   strconv.Itoa(performer.ID),
				Results: performerResults,
			}

			ret = append(ret, &result)
		}
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

func enumToStringPtr(e fmt.Stringer, titleCase bool) *string {
	if e != nil {
		ret := e.String()
		if titleCase {
			ret = strings.Title(strings.ToLower(ret))
		}
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
	switch {
	case end == nil:
		ret = fmt.Sprintf("%d -", *start)
	case start == nil:
		ret = fmt.Sprintf("- %d", *end)
	default:
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

func fetchImage(ctx context.Context, client *http.Client, url string) (*string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
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

func performerFragmentToScrapedScenePerformer(p graphql.PerformerFragment) *models.ScrapedPerformer {
	id := p.ID
	images := []string{}
	for _, image := range p.Images {
		images = append(images, image.URL)
	}
	sp := &models.ScrapedPerformer{
		Name:         &p.Name,
		Country:      p.Country,
		Measurements: formatMeasurements(p.Measurements),
		CareerLength: formatCareerLength(p.CareerStartYear, p.CareerEndYear),
		Tattoos:      formatBodyModifications(p.Tattoos),
		Piercings:    formatBodyModifications(p.Piercings),
		Twitter:      findURL(p.Urls, "TWITTER"),
		RemoteSiteID: &id,
		Images:       images,
		// TODO - tags not currently supported
		// graphql schema change to accommodate this. Leave off for now.
	}

	if len(sp.Images) > 0 {
		sp.Image = &sp.Images[0]
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
		sp.Gender = enumToStringPtr(p.Gender, false)
	}

	if p.Ethnicity != nil {
		sp.Ethnicity = enumToStringPtr(p.Ethnicity, true)
	}

	if p.EyeColor != nil {
		sp.EyeColor = enumToStringPtr(p.EyeColor, true)
	}

	if p.HairColor != nil {
		sp.HairColor = enumToStringPtr(p.HairColor, true)
	}

	if p.BreastType != nil {
		sp.FakeTits = enumToStringPtr(p.BreastType, true)
	}

	if len(p.Aliases) > 0 {
		alias := strings.Join(p.Aliases, ", ")
		sp.Aliases = &alias
	}

	return sp
}

func getFirstImage(ctx context.Context, client *http.Client, images []*graphql.ImageFragment) *string {
	ret, err := fetchImage(ctx, client, images[0].URL)
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

func (c Client) sceneFragmentToScrapedScene(ctx context.Context, s *graphql.SceneFragment) (*models.ScrapedScene, error) {
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
		ss.Image = getFirstImage(ctx, c.getHTTPClient(), s.Images)
	}

	if err := c.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		pqb := r.Performer()
		tqb := r.Tag()

		if s.Studio != nil {
			studioID := s.Studio.ID
			ss.Studio = &models.ScrapedStudio{
				Name:         s.Studio.Name,
				URL:          findURL(s.Studio.Urls, "HOME"),
				RemoteSiteID: &studioID,
			}

			err := match.ScrapedStudio(r.Studio(), ss.Studio, &c.box.Endpoint)
			if err != nil {
				return err
			}
		}

		for _, p := range s.Performers {
			sp := performerFragmentToScrapedScenePerformer(p.Performer)

			err := match.ScrapedPerformer(pqb, sp, &c.box.Endpoint)
			if err != nil {
				return err
			}

			ss.Performers = append(ss.Performers, sp)
		}

		for _, t := range s.Tags {
			st := &models.ScrapedTag{
				Name: t.Name,
			}

			err := match.ScrapedTag(tqb, st)
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

func (c Client) FindStashBoxPerformerByID(ctx context.Context, id string) (*models.ScrapedPerformer, error) {
	performer, err := c.client.FindPerformerByID(ctx, id)
	if err != nil {
		return nil, err
	}

	ret := performerFragmentToScrapedScenePerformer(*performer.FindPerformer)
	return ret, nil
}

func (c Client) FindStashBoxPerformerByName(ctx context.Context, name string) (*models.ScrapedPerformer, error) {
	performers, err := c.client.SearchPerformer(ctx, name)
	if err != nil {
		return nil, err
	}

	var ret *models.ScrapedPerformer
	for _, performer := range performers.SearchPerformer {
		if strings.EqualFold(performer.Name, name) {
			ret = performerFragmentToScrapedScenePerformer(*performer)
		}
	}

	return ret, nil
}
