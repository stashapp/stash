package stashbox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Yamashou/gqlgenc/client"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/Yamashou/gqlgenc/graphqljson"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stashapp/stash/pkg/scraper/stashbox/graphql"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/studio"
	"github.com/stashapp/stash/pkg/tag"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
)

type SceneReader interface {
	Find(ctx context.Context, id int) (*models.Scene, error)
}

type PerformerReader interface {
	match.PerformerFinder
	Find(ctx context.Context, id int) (*models.Performer, error)
	FindBySceneID(ctx context.Context, sceneID int) ([]*models.Performer, error)
	GetStashIDs(ctx context.Context, performerID int) ([]*models.StashID, error)
	GetImage(ctx context.Context, performerID int) ([]byte, error)
}

type StudioReader interface {
	match.StudioFinder
	studio.Finder
	GetStashIDs(ctx context.Context, studioID int) ([]*models.StashID, error)
}
type TagFinder interface {
	tag.Queryer
	FindBySceneID(ctx context.Context, sceneID int) ([]*models.Tag, error)
}

type Repository struct {
	Scene     SceneReader
	Performer PerformerReader
	Tag       TagFinder
	Studio    StudioReader
}

// Client represents the client interface to a stash-box server instance.
type Client struct {
	client     *graphql.Client
	txnManager txn.Manager
	repository Repository
	box        models.StashBox
}

// NewClient returns a new instance of a stash-box client.
func NewClient(box models.StashBox, txnManager txn.Manager, repo Repository) *Client {
	authHeader := func(req *http.Request) {
		req.Header.Set("ApiKey", box.APIKey)
	}

	client := &graphql.Client{
		Client: client.NewClient(http.DefaultClient, box.Endpoint, authHeader),
	}

	return &Client{
		client:     client,
		txnManager: txnManager,
		repository: repo,
		box:        box,
	}
}

func (c Client) getHTTPClient() *http.Client {
	return c.client.Client.Client
}

// QueryStashBoxScene queries stash-box for scenes using a query string.
func (c Client) QueryStashBoxScene(ctx context.Context, queryStr string) ([]*scraper.ScrapedScene, error) {
	scenes, err := c.client.SearchScene(ctx, queryStr)
	if err != nil {
		return nil, err
	}

	sceneFragments := scenes.SearchScene

	var ret []*scraper.ScrapedScene
	for _, s := range sceneFragments {
		ss, err := c.sceneFragmentToScrapedScene(ctx, s)
		if err != nil {
			return nil, err
		}
		ret = append(ret, ss)
	}

	return ret, nil
}

// FindStashBoxScenesByFingerprints queries stash-box for a scene using the
// scene's MD5/OSHASH checksum, or PHash.
func (c Client) FindStashBoxSceneByFingerprints(ctx context.Context, sceneID int) ([]*scraper.ScrapedScene, error) {
	res, err := c.FindStashBoxScenesByFingerprints(ctx, []int{sceneID})
	if len(res) > 0 {
		return res[0], err
	}
	return nil, err
}

// FindStashBoxScenesByFingerprints queries stash-box for scenes using every
// scene's MD5/OSHASH checksum, or PHash, and returns results in the same order
// as the input slice.
func (c Client) FindStashBoxScenesByFingerprints(ctx context.Context, ids []int) ([][]*scraper.ScrapedScene, error) {
	var fingerprints [][]*graphql.FingerprintQueryInput

	if err := txn.WithTxn(ctx, c.txnManager, func(ctx context.Context) error {
		qb := c.repository.Scene

		for _, sceneID := range ids {
			scene, err := qb.Find(ctx, sceneID)
			if err != nil {
				return err
			}

			if scene == nil {
				return fmt.Errorf("scene with id %d not found", sceneID)
			}

			var sceneFPs []*graphql.FingerprintQueryInput

			if scene.Checksum != nil {
				sceneFPs = append(sceneFPs, &graphql.FingerprintQueryInput{
					Hash:      *scene.Checksum,
					Algorithm: graphql.FingerprintAlgorithmMd5,
				})
			}

			if scene.OSHash != nil {
				sceneFPs = append(sceneFPs, &graphql.FingerprintQueryInput{
					Hash:      *scene.OSHash,
					Algorithm: graphql.FingerprintAlgorithmOshash,
				})
			}

			if scene.Phash != nil {
				phashStr := utils.PhashToString(*scene.Phash)
				sceneFPs = append(sceneFPs, &graphql.FingerprintQueryInput{
					Hash:      phashStr,
					Algorithm: graphql.FingerprintAlgorithmPhash,
				})
			}

			fingerprints = append(fingerprints, sceneFPs)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return c.findStashBoxScenesByFingerprints(ctx, fingerprints)
}

func (c Client) findStashBoxScenesByFingerprints(ctx context.Context, scenes [][]*graphql.FingerprintQueryInput) ([][]*scraper.ScrapedScene, error) {
	var ret [][]*scraper.ScrapedScene
	for i := 0; i < len(scenes); i += 40 {
		end := i + 40
		if end > len(scenes) {
			end = len(scenes)
		}
		scenes, err := c.client.FindScenesBySceneFingerprints(ctx, scenes[i:end])

		if err != nil {
			return nil, err
		}

		for _, sceneFragments := range scenes.FindScenesBySceneFingerprints {
			var sceneResults []*scraper.ScrapedScene
			for _, scene := range sceneFragments {
				ss, err := c.sceneFragmentToScrapedScene(ctx, scene)
				if err != nil {
					return nil, err
				}
				sceneResults = append(sceneResults, ss)
			}
			ret = append(ret, sceneResults)
		}
	}

	return ret, nil
}

func (c Client) SubmitStashBoxFingerprints(ctx context.Context, sceneIDs []string, endpoint string) (bool, error) {
	ids, err := stringslice.StringSliceToIntSlice(sceneIDs)
	if err != nil {
		return false, err
	}

	var fingerprints []graphql.FingerprintSubmission

	if err := txn.WithTxn(ctx, c.txnManager, func(ctx context.Context) error {
		qb := c.repository.Scene

		for _, sceneID := range ids {
			scene, err := qb.Find(ctx, sceneID)
			if err != nil {
				return err
			}

			if scene == nil {
				continue
			}

			stashIDs := scene.StashIDs
			sceneStashID := ""
			for _, stashID := range stashIDs {
				if stashID.Endpoint == endpoint {
					sceneStashID = stashID.StashID
				}
			}

			if sceneStashID != "" {
				if scene.Checksum != nil && scene.Duration != nil {
					fingerprint := graphql.FingerprintInput{
						Hash:      *scene.Checksum,
						Algorithm: graphql.FingerprintAlgorithmMd5,
						Duration:  int(*scene.Duration),
					}
					fingerprints = append(fingerprints, graphql.FingerprintSubmission{
						SceneID:     sceneStashID,
						Fingerprint: &fingerprint,
					})
				}

				if scene.OSHash != nil && scene.Duration != nil {
					fingerprint := graphql.FingerprintInput{
						Hash:      *scene.OSHash,
						Algorithm: graphql.FingerprintAlgorithmOshash,
						Duration:  int(*scene.Duration),
					}
					fingerprints = append(fingerprints, graphql.FingerprintSubmission{
						SceneID:     sceneStashID,
						Fingerprint: &fingerprint,
					})
				}

				if scene.Phash != nil && scene.Duration != nil {
					fingerprint := graphql.FingerprintInput{
						Hash:      utils.PhashToString(*scene.Phash),
						Algorithm: graphql.FingerprintAlgorithmPhash,
						Duration:  int(*scene.Duration),
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
func (c Client) QueryStashBoxPerformer(ctx context.Context, queryStr string) ([]*StashBoxPerformerQueryResult, error) {
	performers, err := c.queryStashBoxPerformer(ctx, queryStr)

	res := []*StashBoxPerformerQueryResult{
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
func (c Client) FindStashBoxPerformersByNames(ctx context.Context, performerIDs []string) ([]*StashBoxPerformerQueryResult, error) {
	ids, err := stringslice.StringSliceToIntSlice(performerIDs)
	if err != nil {
		return nil, err
	}

	var performers []*models.Performer

	if err := txn.WithTxn(ctx, c.txnManager, func(ctx context.Context) error {
		qb := c.repository.Performer

		for _, performerID := range ids {
			performer, err := qb.Find(ctx, performerID)
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
	ids, err := stringslice.StringSliceToIntSlice(performerIDs)
	if err != nil {
		return nil, err
	}

	var performers []*models.Performer

	if err := txn.WithTxn(ctx, c.txnManager, func(ctx context.Context) error {
		qb := c.repository.Performer

		for _, performerID := range ids {
			performer, err := qb.Find(ctx, performerID)
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

func (c Client) findStashBoxPerformersByNames(ctx context.Context, performers []*models.Performer) ([]*StashBoxPerformerQueryResult, error) {
	var ret []*StashBoxPerformerQueryResult
	for _, performer := range performers {
		if performer.Name.Valid {
			performerResults, err := c.queryStashBoxPerformer(ctx, performer.Name.String)
			if err != nil {
				return nil, err
			}

			result := StashBoxPerformerQueryResult{
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
		ret := strings.ReplaceAll(e.String(), "_", " ")
		if titleCase {
			c := cases.Title(language.Und)
			ret = c.String(strings.ToLower(ret))
		}
		return &ret
	}

	return nil
}

func translateGender(gender *graphql.GenderEnum) *string {
	var res models.GenderEnum
	switch *gender {
	case graphql.GenderEnumMale:
		res = models.GenderEnumMale
	case graphql.GenderEnumFemale:
		res = models.GenderEnumFemale
	case graphql.GenderEnumIntersex:
		res = models.GenderEnumIntersex
	case graphql.GenderEnumTransgenderFemale:
		res = models.GenderEnumTransgenderFemale
	case graphql.GenderEnumTransgenderMale:
		res = models.GenderEnumTransgenderMale
	}

	if res != "" {
		strVal := res.String()
		return &strVal
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
		sp.Gender = translateGender(p.Gender)
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

func (c Client) sceneFragmentToScrapedScene(ctx context.Context, s *graphql.SceneFragment) (*scraper.ScrapedScene, error) {
	stashID := s.ID
	ss := &scraper.ScrapedScene{
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

	if err := txn.WithTxn(ctx, c.txnManager, func(ctx context.Context) error {
		pqb := c.repository.Performer
		tqb := c.repository.Tag

		if s.Studio != nil {
			studioID := s.Studio.ID
			ss.Studio = &models.ScrapedStudio{
				Name:         s.Studio.Name,
				URL:          findURL(s.Studio.Urls, "HOME"),
				RemoteSiteID: &studioID,
			}

			err := match.ScrapedStudio(ctx, c.repository.Studio, ss.Studio, &c.box.Endpoint)
			if err != nil {
				return err
			}
		}

		for _, p := range s.Performers {
			sp := performerFragmentToScrapedScenePerformer(p.Performer)

			err := match.ScrapedPerformer(ctx, pqb, sp, &c.box.Endpoint)
			if err != nil {
				return err
			}

			ss.Performers = append(ss.Performers, sp)
		}

		for _, t := range s.Tags {
			st := &models.ScrapedTag{
				Name: t.Name,
			}

			err := match.ScrapedTag(ctx, tqb, st)
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

func (c Client) GetUser(ctx context.Context) (*graphql.Me, error) {
	return c.client.Me(ctx)
}

func (c Client) SubmitSceneDraft(ctx context.Context, sceneID int, endpoint string, imagePath string) (*string, error) {
	draft := graphql.SceneDraftInput{}
	var image *os.File
	if err := txn.WithTxn(ctx, c.txnManager, func(ctx context.Context) error {
		r := c.repository
		qb := r.Scene
		pqb := r.Performer
		sqb := r.Studio

		scene, err := qb.Find(ctx, sceneID)
		if err != nil {
			return err
		}

		if scene.Title != "" {
			draft.Title = &scene.Title
		}
		if scene.Details != "" {
			draft.Details = &scene.Details
		}
		if scene.URL != "" && len(strings.TrimSpace(scene.URL)) > 0 {
			url := strings.TrimSpace(scene.URL)
			draft.URL = &url
		}
		if scene.Date != nil {
			v := scene.Date.String()
			draft.Date = &v
		}

		if scene.StudioID != nil {
			studio, err := sqb.Find(ctx, int(*scene.StudioID))
			if err != nil {
				return err
			}
			studioDraft := graphql.DraftEntityInput{
				Name: studio.Name.String,
			}

			stashIDs, err := sqb.GetStashIDs(ctx, studio.ID)
			if err != nil {
				return err
			}
			for _, stashID := range stashIDs {
				if stashID.Endpoint == endpoint {
					studioDraft.ID = &stashID.StashID
					break
				}
			}
			draft.Studio = &studioDraft
		}

		fingerprints := []*graphql.FingerprintInput{}
		if scene.OSHash != nil && scene.Duration != nil {
			fingerprint := graphql.FingerprintInput{
				Hash:      *scene.OSHash,
				Algorithm: graphql.FingerprintAlgorithmOshash,
				Duration:  int(*scene.Duration),
			}
			fingerprints = append(fingerprints, &fingerprint)
		}

		if scene.Checksum != nil && scene.Duration != nil {
			fingerprint := graphql.FingerprintInput{
				Hash:      *scene.Checksum,
				Algorithm: graphql.FingerprintAlgorithmMd5,
				Duration:  int(*scene.Duration),
			}
			fingerprints = append(fingerprints, &fingerprint)
		}

		if scene.Phash != nil && scene.Duration != nil {
			fingerprint := graphql.FingerprintInput{
				Hash:      utils.PhashToString(*scene.Phash),
				Algorithm: graphql.FingerprintAlgorithmPhash,
				Duration:  int(*scene.Duration),
			}
			fingerprints = append(fingerprints, &fingerprint)
		}
		draft.Fingerprints = fingerprints

		scenePerformers, err := pqb.FindBySceneID(ctx, sceneID)
		if err != nil {
			return err
		}

		performers := []*graphql.DraftEntityInput{}
		for _, p := range scenePerformers {
			performerDraft := graphql.DraftEntityInput{
				Name: p.Name.String,
			}

			stashIDs, err := pqb.GetStashIDs(ctx, p.ID)
			if err != nil {
				return err
			}

			for _, stashID := range stashIDs {
				if stashID.Endpoint == endpoint {
					performerDraft.ID = &stashID.StashID
					break
				}
			}

			performers = append(performers, &performerDraft)
		}
		draft.Performers = performers

		var tags []*graphql.DraftEntityInput
		sceneTags, err := r.Tag.FindBySceneID(ctx, scene.ID)
		if err != nil {
			return err
		}
		for _, tag := range sceneTags {
			tags = append(tags, &graphql.DraftEntityInput{Name: tag.Name})
		}
		draft.Tags = tags

		exists, _ := fsutil.FileExists(imagePath)
		if exists {
			file, err := os.Open(imagePath)
			if err == nil {
				image = file
			}
		}

		stashIDs := scene.StashIDs
		var stashID *string
		for _, v := range stashIDs {
			if v.Endpoint == endpoint {
				stashID = &v.StashID
				break
			}
		}
		draft.ID = stashID

		return nil
	}); err != nil {
		return nil, err
	}

	var id *string
	var ret graphql.SubmitSceneDraft
	err := c.submitDraft(ctx, graphql.SubmitSceneDraftDocument, draft, image, &ret)
	id = ret.SubmitSceneDraft.ID

	return id, err

	// ret, err := c.client.SubmitSceneDraft(ctx, draft, uploadImage(image))
	// if err != nil {
	// 	return nil, err
	// }

	// id := ret.SubmitSceneDraft.ID
	// return id, nil
}

func (c Client) SubmitPerformerDraft(ctx context.Context, performer *models.Performer, endpoint string) (*string, error) {
	draft := graphql.PerformerDraftInput{}
	var image io.Reader
	if err := txn.WithTxn(ctx, c.txnManager, func(ctx context.Context) error {
		pqb := c.repository.Performer
		img, _ := pqb.GetImage(ctx, performer.ID)
		if img != nil {
			image = bytes.NewReader(img)
		}

		if performer.Name.Valid {
			draft.Name = performer.Name.String
		}
		if performer.Birthdate.Valid {
			draft.Birthdate = &performer.Birthdate.String
		}
		if performer.Country.Valid {
			draft.Country = &performer.Country.String
		}
		if performer.Ethnicity.Valid {
			draft.Ethnicity = &performer.Ethnicity.String
		}
		if performer.EyeColor.Valid {
			draft.EyeColor = &performer.EyeColor.String
		}
		if performer.FakeTits.Valid {
			draft.BreastType = &performer.FakeTits.String
		}
		if performer.Gender.Valid {
			draft.Gender = &performer.Gender.String
		}
		if performer.HairColor.Valid {
			draft.HairColor = &performer.HairColor.String
		}
		if performer.Height.Valid {
			draft.Height = &performer.Height.String
		}
		if performer.Measurements.Valid {
			draft.Measurements = &performer.Measurements.String
		}
		if performer.Piercings.Valid {
			draft.Piercings = &performer.Piercings.String
		}
		if performer.Tattoos.Valid {
			draft.Tattoos = &performer.Tattoos.String
		}
		if performer.Aliases.Valid {
			draft.Aliases = &performer.Aliases.String
		}

		var urls []string
		if len(strings.TrimSpace(performer.Twitter.String)) > 0 {
			urls = append(urls, "https://twitter.com/"+strings.TrimSpace(performer.Twitter.String))
		}
		if len(strings.TrimSpace(performer.Instagram.String)) > 0 {
			urls = append(urls, "https://instagram.com/"+strings.TrimSpace(performer.Instagram.String))
		}
		if len(strings.TrimSpace(performer.URL.String)) > 0 {
			urls = append(urls, strings.TrimSpace(performer.URL.String))
		}
		if len(urls) > 0 {
			draft.Urls = urls
		}

		stashIDs, err := pqb.GetStashIDs(ctx, performer.ID)
		if err != nil {
			return err
		}
		var stashID *string
		for _, v := range stashIDs {
			if v.Endpoint == endpoint {
				stashID = &v.StashID
				break
			}
		}
		draft.ID = stashID

		return nil
	}); err != nil {
		return nil, err
	}

	var id *string
	var ret graphql.SubmitPerformerDraft
	err := c.submitDraft(ctx, graphql.SubmitPerformerDraftDocument, draft, image, &ret)
	id = ret.SubmitPerformerDraft.ID

	return id, err

	// ret, err := c.client.SubmitPerformerDraft(ctx, draft, uploadImage(image))
	// if err != nil {
	// 	return nil, err
	// }

	// id := ret.SubmitPerformerDraft.ID
	// return id, nil
}

// we can't currently use this due to https://github.com/Yamashou/gqlgenc/issues/109
// func uploadImage(image io.Reader) client.HTTPRequestOption {
// 	return func(req *http.Request) {
// 		if image == nil {
// 			// return without changing anything
// 			return
// 		}

// 		// we can't handle errors in here, so if one happens, just return
// 		// without changing anything.

// 		// repackage the request to include the image
// 		bodyBytes, err := ioutil.ReadAll(req.Body)
// 		if err != nil {
// 			return
// 		}

// 		newBody := &bytes.Buffer{}
// 		writer := multipart.NewWriter(newBody)
// 		_ = writer.WriteField("operations", string(bodyBytes))

// 		if err := writer.WriteField("map", "{ \"0\": [\"variables.input.image\"] }"); err != nil {
// 			return
// 		}
// 		part, _ := writer.CreateFormFile("0", "draft")
// 		if _, err := io.Copy(part, image); err != nil {
// 			return
// 		}

// 		writer.Close()

// 		// now set the request body to this new body
// 		req.Body = io.NopCloser(newBody)
// 		req.ContentLength = int64(newBody.Len())
// 		req.Header.Set("Content-Type", writer.FormDataContentType())
// 	}
// }

func (c *Client) submitDraft(ctx context.Context, query string, input interface{}, image io.Reader, ret interface{}) error {
	vars := map[string]interface{}{
		"input": input,
	}

	r := &client.Request{
		Query:         query,
		Variables:     vars,
		OperationName: "",
	}

	requestBody, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err := writer.WriteField("operations", string(requestBody)); err != nil {
		return err
	}

	if image != nil {
		if err := writer.WriteField("map", "{ \"0\": [\"variables.input.image\"] }"); err != nil {
			return err
		}
		part, _ := writer.CreateFormFile("0", "draft")
		if _, err := io.Copy(part, image); err != nil {
			return err
		}
	} else if err := writer.WriteField("map", "{}"); err != nil {
		return err
	}

	writer.Close()

	req, _ := http.NewRequestWithContext(ctx, "POST", c.box.Endpoint, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Set("ApiKey", c.box.APIKey)

	httpClient := c.client.Client.Client
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type response struct {
		Data   json.RawMessage `json:"data"`
		Errors json.RawMessage `json:"errors"`
	}

	var respGQL response

	if err := json.Unmarshal(responseBytes, &respGQL); err != nil {
		return fmt.Errorf("failed to decode data %s: %w", string(responseBytes), err)
	}

	if respGQL.Errors != nil && len(respGQL.Errors) > 0 {
		// try to parse standard graphql error
		errors := &client.GqlErrorList{}
		if e := json.Unmarshal(responseBytes, errors); e != nil {
			return fmt.Errorf("failed to parse graphql errors. Response content %s - %w ", string(responseBytes), e)
		}

		return errors
	}

	if err := graphqljson.UnmarshalData(respGQL.Data, ret); err != nil {
		return err
	}

	return err
}
