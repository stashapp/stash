// Package stashbox provides a client interface to a stash-box server instance.
package stashbox

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/Yamashou/gqlgenc/clientv2"
	"github.com/Yamashou/gqlgenc/graphqljson"
	"github.com/gofrs/uuid/v5"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stashapp/stash/pkg/scraper/stashbox/graphql"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
)

type SceneReader interface {
	models.SceneGetter
	models.StashIDLoader
	models.VideoFileLoader
}

type PerformerReader interface {
	models.PerformerGetter
	match.PerformerFinder
	models.AliasLoader
	models.StashIDLoader
	models.URLLoader
	FindBySceneID(ctx context.Context, sceneID int) ([]*models.Performer, error)
	GetImage(ctx context.Context, performerID int) ([]byte, error)
}

type StudioReader interface {
	models.StudioGetter
	match.StudioFinder
	models.StashIDLoader
}

type TagFinder interface {
	models.TagQueryer
	FindBySceneID(ctx context.Context, sceneID int) ([]*models.Tag, error)
}

type Repository struct {
	TxnManager models.TxnManager

	Scene     SceneReader
	Performer PerformerReader
	Tag       TagFinder
	Studio    StudioReader
}

func NewRepository(repo models.Repository) Repository {
	return Repository{
		TxnManager: repo.TxnManager,
		Scene:      repo.Scene,
		Performer:  repo.Performer,
		Tag:        repo.Tag,
		Studio:     repo.Studio,
	}
}

func (r *Repository) WithReadTxn(ctx context.Context, fn txn.TxnFunc) error {
	return txn.WithReadTxn(ctx, r.TxnManager, fn)
}

// Client represents the client interface to a stash-box server instance.
type Client struct {
	client     *graphql.Client
	repository Repository
	box        models.StashBox
}

// NewClient returns a new instance of a stash-box client.
func NewClient(box models.StashBox, repo Repository) *Client {
	authHeader := func(ctx context.Context, req *http.Request, gqlInfo *clientv2.GQLRequestInfo, res interface{}, next clientv2.RequestInterceptorFunc) error {
		req.Header.Set("ApiKey", box.APIKey)
		return next(ctx, req, gqlInfo, res)
	}

	client := &graphql.Client{
		Client: clientv2.NewClient(http.DefaultClient, box.Endpoint, nil, authHeader),
	}

	return &Client{
		client:     client,
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

	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		qb := r.Scene

		for _, sceneID := range ids {
			scene, err := qb.Find(ctx, sceneID)
			if err != nil {
				return err
			}

			if scene == nil {
				return fmt.Errorf("scene with id %d not found", sceneID)
			}

			if err := scene.LoadFiles(ctx, r.Scene); err != nil {
				return err
			}

			var sceneFPs []*graphql.FingerprintQueryInput

			for _, f := range scene.Files.List() {
				checksum := f.Fingerprints.GetString(models.FingerprintTypeMD5)
				if checksum != "" {
					sceneFPs = append(sceneFPs, &graphql.FingerprintQueryInput{
						Hash:      checksum,
						Algorithm: graphql.FingerprintAlgorithmMd5,
					})
				}

				oshash := f.Fingerprints.GetString(models.FingerprintTypeOshash)
				if oshash != "" {
					sceneFPs = append(sceneFPs, &graphql.FingerprintQueryInput{
						Hash:      oshash,
						Algorithm: graphql.FingerprintAlgorithmOshash,
					})
				}

				phash := f.Fingerprints.GetInt64(models.FingerprintTypePhash)
				if phash != 0 {
					phashStr := utils.PhashToString(phash)
					sceneFPs = append(sceneFPs, &graphql.FingerprintQueryInput{
						Hash:      phashStr,
						Algorithm: graphql.FingerprintAlgorithmPhash,
					})
				}
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
	var results [][]*scraper.ScrapedScene

	// filter out nils
	var validScenes [][]*graphql.FingerprintQueryInput
	for _, s := range scenes {
		if len(s) > 0 {
			validScenes = append(validScenes, s)
		}
	}

	for i := 0; i < len(validScenes); i += 40 {
		end := i + 40
		if end > len(validScenes) {
			end = len(validScenes)
		}
		scenes, err := c.client.FindScenesBySceneFingerprints(ctx, validScenes[i:end])

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
			results = append(results, sceneResults)
		}
	}

	// repopulate the results to be the same order as the input
	ret := make([][]*scraper.ScrapedScene, len(scenes))
	upTo := 0

	for i, v := range scenes {
		if len(v) > 0 {
			ret[i] = results[upTo]
			upTo++
		}
	}

	return ret, nil
}

func (c Client) SubmitStashBoxFingerprints(ctx context.Context, sceneIDs []string) (bool, error) {
	ids, err := stringslice.StringSliceToIntSlice(sceneIDs)
	if err != nil {
		return false, err
	}

	endpoint := c.box.Endpoint

	var fingerprints []graphql.FingerprintSubmission

	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		qb := r.Scene

		for _, sceneID := range ids {
			scene, err := qb.Find(ctx, sceneID)
			if err != nil {
				return err
			}

			if scene == nil {
				continue
			}

			if err := scene.LoadStashIDs(ctx, qb); err != nil {
				return err
			}

			if err := scene.LoadFiles(ctx, qb); err != nil {
				return err
			}

			stashIDs := scene.StashIDs.List()
			sceneStashID := ""
			for _, stashID := range stashIDs {
				if stashID.Endpoint == endpoint {
					sceneStashID = stashID.StashID
				}
			}

			if sceneStashID != "" {
				for _, f := range scene.Files.List() {
					duration := f.Duration

					if duration != 0 {
						if checksum := f.Fingerprints.GetString(models.FingerprintTypeMD5); checksum != "" {
							fingerprint := graphql.FingerprintInput{
								Hash:      checksum,
								Algorithm: graphql.FingerprintAlgorithmMd5,
								Duration:  int(duration),
							}
							fingerprints = append(fingerprints, graphql.FingerprintSubmission{
								SceneID:     sceneStashID,
								Fingerprint: &fingerprint,
							})
						}

						if oshash := f.Fingerprints.GetString(models.FingerprintTypeOshash); oshash != "" {
							fingerprint := graphql.FingerprintInput{
								Hash:      oshash,
								Algorithm: graphql.FingerprintAlgorithmOshash,
								Duration:  int(duration),
							}
							fingerprints = append(fingerprints, graphql.FingerprintSubmission{
								SceneID:     sceneStashID,
								Fingerprint: &fingerprint,
							})
						}

						if phash := f.Fingerprints.GetInt64(models.FingerprintTypePhash); phash != 0 {
							fingerprint := graphql.FingerprintInput{
								Hash:      utils.PhashToString(phash),
								Algorithm: graphql.FingerprintAlgorithmPhash,
								Duration:  int(duration),
							}
							fingerprints = append(fingerprints, graphql.FingerprintSubmission{
								SceneID:     sceneStashID,
								Fingerprint: &fingerprint,
							})
						}
					}
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
		performer := performerFragmentToScrapedPerformer(*fragment)
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
	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		qb := r.Performer

		for _, performerID := range ids {
			performer, err := qb.Find(ctx, performerID)
			if err != nil {
				return err
			}

			if performer == nil {
				return fmt.Errorf("performer with id %d not found", performerID)
			}

			if performer.Name != "" {
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

	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		qb := r.Performer

		for _, performerID := range ids {
			performer, err := qb.Find(ctx, performerID)
			if err != nil {
				return err
			}

			if performer == nil {
				return fmt.Errorf("performer with id %d not found", performerID)
			}

			if performer.Name != "" {
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
		if performer.Name != "" {
			performerResults, err := c.queryStashBoxPerformer(ctx, performer.Name)
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
	case graphql.GenderEnumNonBinary:
		res = models.GenderEnumNonBinary
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

func performerFragmentToScrapedPerformer(p graphql.PerformerFragment) *models.ScrapedPerformer {
	images := []string{}
	for _, image := range p.Images {
		images = append(images, image.URL)
	}

	sp := &models.ScrapedPerformer{
		Name:           &p.Name,
		Disambiguation: p.Disambiguation,
		Country:        p.Country,
		Measurements:   formatMeasurements(*p.Measurements),
		CareerLength:   formatCareerLength(p.CareerStartYear, p.CareerEndYear),
		Tattoos:        formatBodyModifications(p.Tattoos),
		Piercings:      formatBodyModifications(p.Piercings),
		Twitter:        findURL(p.Urls, "TWITTER"),
		RemoteSiteID:   &p.ID,
		Images:         images,
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
		// #4437 - stash-box may return aliases that are equal to the performer name
		// filter these out
		p.Aliases = sliceutil.Filter(p.Aliases, func(s string) bool {
			return !strings.EqualFold(s, p.Name)
		})

		// #4596 - stash-box may return duplicate aliases. Filter these out
		p.Aliases = stringslice.UniqueFold(p.Aliases)

		alias := strings.Join(p.Aliases, ", ")
		sp.Aliases = &alias
	}

	for _, u := range p.Urls {
		sp.URLs = append(sp.URLs, u.URL)
	}

	return sp
}

func studioFragmentToScrapedStudio(s graphql.StudioFragment) *models.ScrapedStudio {
	images := []string{}
	for _, image := range s.Images {
		images = append(images, image.URL)
	}

	st := &models.ScrapedStudio{
		Name:         s.Name,
		URL:          findURL(s.Urls, "HOME"),
		Images:       images,
		RemoteSiteID: &s.ID,
	}

	if len(st.Images) > 0 {
		st.Image = &st.Images[0]
	}

	return st
}

func getFirstImage(ctx context.Context, client *http.Client, images []*graphql.ImageFragment) *string {
	ret, err := fetchImage(ctx, client, images[0].URL)
	if err != nil && !errors.Is(err, context.Canceled) {
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
		Code:         s.Code,
		Date:         s.Date,
		Details:      s.Details,
		Director:     s.Director,
		URL:          findURL(s.Urls, "STUDIO"),
		Duration:     s.Duration,
		RemoteSiteID: &stashID,
		Fingerprints: getFingerprints(s),
		// Image
		// stash_id
	}

	for _, u := range s.Urls {
		ss.URLs = append(ss.URLs, u.URL)
	}

	if len(ss.URLs) > 0 {
		ss.URL = &ss.URLs[0]
	}

	if len(s.Images) > 0 {
		// TODO - #454 code sorts images by aspect ratio according to a wanted
		// orientation. I'm just grabbing the first for now
		ss.Image = getFirstImage(ctx, c.getHTTPClient(), s.Images)
	}

	if ss.URL == nil && len(s.Urls) > 0 {
		// The scene in Stash-box may not have a Studio URL but it does have another URL.
		// For example it has a www.manyvids.com URL, which is auto set as type ManyVids.
		// This should be re-visited once Stashapp can support more than one URL.
		ss.URL = &s.Urls[0].URL
	}

	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		pqb := r.Performer
		tqb := r.Tag

		if s.Studio != nil {
			ss.Studio = studioFragmentToScrapedStudio(*s.Studio)

			err := match.ScrapedStudio(ctx, r.Studio, ss.Studio, &c.box.Endpoint)
			if err != nil {
				return err
			}

			var parentStudio *graphql.FindStudio
			if s.Studio.Parent != nil {
				parentStudio, err = c.client.FindStudio(ctx, &s.Studio.Parent.ID, nil)
				if err != nil {
					return err
				}

				if parentStudio.FindStudio != nil {
					ss.Studio.Parent = studioFragmentToScrapedStudio(*parentStudio.FindStudio)

					err = match.ScrapedStudio(ctx, r.Studio, ss.Studio.Parent, &c.box.Endpoint)
					if err != nil {
						return err
					}
				}
			}
		}

		for _, p := range s.Performers {
			sp := performerFragmentToScrapedPerformer(*p.Performer)

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

	if performer.FindPerformer == nil {
		return nil, nil
	}

	ret := performerFragmentToScrapedPerformer(*performer.FindPerformer)

	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		err := match.ScrapedPerformer(ctx, r.Performer, ret, &c.box.Endpoint)
		return err
	}); err != nil {
		return nil, err
	}

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
			ret = performerFragmentToScrapedPerformer(*performer)
		}
	}

	if ret == nil {
		return nil, nil
	}

	r := c.repository
	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		err := match.ScrapedPerformer(ctx, r.Performer, ret, &c.box.Endpoint)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (c Client) FindStashBoxStudio(ctx context.Context, query string) (*models.ScrapedStudio, error) {
	var studio *graphql.FindStudio

	_, err := uuid.FromString(query)
	if err == nil {
		// Confirmed the user passed in a Stash ID
		studio, err = c.client.FindStudio(ctx, &query, nil)
	} else {
		// Otherwise assume they're searching on a name
		studio, err = c.client.FindStudio(ctx, nil, &query)
	}

	if err != nil {
		return nil, err
	}

	var ret *models.ScrapedStudio
	if studio.FindStudio != nil {
		r := c.repository
		if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
			ret = studioFragmentToScrapedStudio(*studio.FindStudio)

			err = match.ScrapedStudio(ctx, r.Studio, ret, &c.box.Endpoint)
			if err != nil {
				return err
			}

			if studio.FindStudio.Parent != nil {
				parentStudio, err := c.client.FindStudio(ctx, &studio.FindStudio.Parent.ID, nil)
				if err != nil {
					return err
				}

				if parentStudio.FindStudio != nil {
					ret.Parent = studioFragmentToScrapedStudio(*parentStudio.FindStudio)

					err = match.ScrapedStudio(ctx, r.Studio, ret.Parent, &c.box.Endpoint)
					if err != nil {
						return err
					}
				}
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func (c Client) GetUser(ctx context.Context) (*graphql.Me, error) {
	return c.client.Me(ctx)
}

func appendFingerprintUnique(v []*graphql.FingerprintInput, toAdd *graphql.FingerprintInput) []*graphql.FingerprintInput {
	for _, vv := range v {
		if vv.Algorithm == toAdd.Algorithm && vv.Hash == toAdd.Hash {
			return v
		}
	}

	return append(v, toAdd)
}

func (c Client) SubmitSceneDraft(ctx context.Context, scene *models.Scene, cover []byte) (*string, error) {
	draft := graphql.SceneDraftInput{}
	var image io.Reader
	r := c.repository
	pqb := r.Performer
	sqb := r.Studio
	endpoint := c.box.Endpoint

	if scene.Title != "" {
		draft.Title = &scene.Title
	}
	if scene.Code != "" {
		draft.Code = &scene.Code
	}
	if scene.Details != "" {
		draft.Details = &scene.Details
	}
	if scene.Director != "" {
		draft.Director = &scene.Director
	}
	// TODO - draft does not accept multiple URLs. Use single URL for now.
	if len(scene.URLs.List()) > 0 {
		url := strings.TrimSpace(scene.URLs.List()[0])
		draft.URL = &url
	}
	if scene.Date != nil {
		v := scene.Date.String()
		draft.Date = &v
	}

	if scene.StudioID != nil {
		studio, err := sqb.Find(ctx, *scene.StudioID)
		if err != nil {
			return nil, err
		}
		if studio == nil {
			return nil, fmt.Errorf("studio with id %d not found", *scene.StudioID)
		}

		studioDraft := graphql.DraftEntityInput{
			Name: studio.Name,
		}

		stashIDs, err := sqb.GetStashIDs(ctx, studio.ID)
		if err != nil {
			return nil, err
		}
		for _, stashID := range stashIDs {
			c := stashID
			if stashID.Endpoint == endpoint {
				studioDraft.ID = &c.StashID
				break
			}
		}
		draft.Studio = &studioDraft
	}

	fingerprints := []*graphql.FingerprintInput{}

	// submit all file fingerprints
	if err := scene.LoadFiles(ctx, r.Scene); err != nil {
		return nil, err
	}

	for _, f := range scene.Files.List() {
		duration := f.Duration

		if duration != 0 {
			if oshash := f.Fingerprints.GetString(models.FingerprintTypeOshash); oshash != "" {
				fingerprint := graphql.FingerprintInput{
					Hash:      oshash,
					Algorithm: graphql.FingerprintAlgorithmOshash,
					Duration:  int(duration),
				}
				fingerprints = appendFingerprintUnique(fingerprints, &fingerprint)
			}

			if checksum := f.Fingerprints.GetString(models.FingerprintTypeMD5); checksum != "" {
				fingerprint := graphql.FingerprintInput{
					Hash:      checksum,
					Algorithm: graphql.FingerprintAlgorithmMd5,
					Duration:  int(duration),
				}
				fingerprints = appendFingerprintUnique(fingerprints, &fingerprint)
			}

			if phash := f.Fingerprints.GetInt64(models.FingerprintTypePhash); phash != 0 {
				fingerprint := graphql.FingerprintInput{
					Hash:      utils.PhashToString(phash),
					Algorithm: graphql.FingerprintAlgorithmPhash,
					Duration:  int(duration),
				}
				fingerprints = appendFingerprintUnique(fingerprints, &fingerprint)
			}
		}
	}
	draft.Fingerprints = fingerprints

	scenePerformers, err := pqb.FindBySceneID(ctx, scene.ID)
	if err != nil {
		return nil, err
	}

	performers := []*graphql.DraftEntityInput{}
	for _, p := range scenePerformers {
		performerDraft := graphql.DraftEntityInput{
			Name: p.Name,
		}

		stashIDs, err := pqb.GetStashIDs(ctx, p.ID)
		if err != nil {
			return nil, err
		}

		for _, stashID := range stashIDs {
			c := stashID
			if stashID.Endpoint == endpoint {
				performerDraft.ID = &c.StashID
				break
			}
		}

		performers = append(performers, &performerDraft)
	}
	draft.Performers = performers

	var tags []*graphql.DraftEntityInput
	sceneTags, err := r.Tag.FindBySceneID(ctx, scene.ID)
	if err != nil {
		return nil, err
	}
	for _, tag := range sceneTags {
		tags = append(tags, &graphql.DraftEntityInput{Name: tag.Name})
	}
	draft.Tags = tags

	if len(cover) > 0 {
		image = bytes.NewReader(cover)
	}

	if err := scene.LoadStashIDs(ctx, r.Scene); err != nil {
		return nil, err
	}

	stashIDs := scene.StashIDs.List()
	var stashID *string
	for _, v := range stashIDs {
		if v.Endpoint == endpoint {
			vv := v.StashID
			stashID = &vv
			break
		}
	}
	draft.ID = stashID

	var id *string
	var ret graphql.SubmitSceneDraft
	err = c.submitDraft(ctx, graphql.SubmitSceneDraftDocument, draft, image, &ret)
	id = ret.SubmitSceneDraft.ID

	return id, err

	// ret, err := c.client.SubmitSceneDraft(ctx, draft, uploadImage(image))
	// if err != nil {
	// 	return nil, err
	// }

	// id := ret.SubmitSceneDraft.ID
	// return id, nil
}

func (c Client) SubmitPerformerDraft(ctx context.Context, performer *models.Performer) (*string, error) {
	draft := graphql.PerformerDraftInput{}
	var image io.Reader
	pqb := c.repository.Performer
	endpoint := c.box.Endpoint

	if err := performer.LoadAliases(ctx, pqb); err != nil {
		return nil, err
	}

	if err := performer.LoadURLs(ctx, pqb); err != nil {
		return nil, err
	}

	if err := performer.LoadStashIDs(ctx, pqb); err != nil {
		return nil, err
	}

	img, _ := pqb.GetImage(ctx, performer.ID)
	if img != nil {
		image = bytes.NewReader(img)
	}

	if performer.Name != "" {
		draft.Name = performer.Name
	}
	if performer.Disambiguation != "" {
		draft.Disambiguation = &performer.Disambiguation
	}
	if performer.Birthdate != nil {
		d := performer.Birthdate.String()
		draft.Birthdate = &d
	}
	if performer.Country != "" {
		draft.Country = &performer.Country
	}
	if performer.Ethnicity != "" {
		draft.Ethnicity = &performer.Ethnicity
	}
	if performer.EyeColor != "" {
		draft.EyeColor = &performer.EyeColor
	}
	if performer.FakeTits != "" {
		draft.BreastType = &performer.FakeTits
	}
	if performer.Gender != nil && performer.Gender.IsValid() {
		v := performer.Gender.String()
		draft.Gender = &v
	}
	if performer.HairColor != "" {
		draft.HairColor = &performer.HairColor
	}
	if performer.Height != nil {
		v := strconv.Itoa(*performer.Height)
		draft.Height = &v
	}
	if performer.Measurements != "" {
		draft.Measurements = &performer.Measurements
	}
	if performer.Piercings != "" {
		draft.Piercings = &performer.Piercings
	}
	if performer.Tattoos != "" {
		draft.Tattoos = &performer.Tattoos
	}
	if len(performer.Aliases.List()) > 0 {
		aliases := strings.Join(performer.Aliases.List(), ",")
		draft.Aliases = &aliases
	}
	if performer.CareerLength != "" {
		var career = strings.Split(performer.CareerLength, "-")
		if i, err := strconv.Atoi(strings.TrimSpace(career[0])); err == nil {
			draft.CareerStartYear = &i
		}
		if len(career) == 2 {
			if y, err := strconv.Atoi(strings.TrimSpace(career[1])); err == nil {
				draft.CareerEndYear = &y
			}
		}
	}

	if len(performer.URLs.List()) > 0 {
		draft.Urls = performer.URLs.List()
	}

	stashIDs, err := pqb.GetStashIDs(ctx, performer.ID)
	if err != nil {
		return nil, err
	}
	var stashID *string
	for _, v := range stashIDs {
		c := v
		if v.Endpoint == endpoint {
			stashID = &c.StashID
			break
		}
	}
	draft.ID = stashID

	var id *string
	var ret graphql.SubmitPerformerDraft
	err = c.submitDraft(ctx, graphql.SubmitPerformerDraftDocument, draft, image, &ret)
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

	r := &clientv2.Request{
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

	responseBytes, err := io.ReadAll(resp.Body)
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

	if len(respGQL.Errors) > 0 {
		// try to parse standard graphql error
		errors := &clientv2.GqlErrorList{}
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
