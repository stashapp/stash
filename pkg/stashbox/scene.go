package stashbox

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/stashbox/graphql"
	"github.com/stashapp/stash/pkg/utils"
)

// QueryScene queries stash-box for scenes using a query string.
func (c Client) QueryScene(ctx context.Context, queryStr string) ([]*models.ScrapedScene, error) {
	scenes, err := c.client.SearchScene(ctx, queryStr)
	if err != nil {
		return nil, err
	}

	sceneFragments := scenes.SearchScene

	var ret []*models.ScrapedScene
	var ignoredTags []string
	for _, s := range sceneFragments {
		ss, err := c.sceneFragmentToScrapedScene(ctx, s)
		if err != nil {
			return nil, err
		}

		var thisIgnoredTags []string
		ss.Tags, thisIgnoredTags = scraper.FilterTags(c.excludeTagRE, ss.Tags)
		ignoredTags = sliceutil.AppendUniques(ignoredTags, thisIgnoredTags)

		ret = append(ret, ss)
	}

	scraper.LogIgnoredTags(ignoredTags)

	return ret, nil
}

// FindStashBoxScenesByFingerprints queries stash-box for a scene using the
// scene's MD5/OSHASH checksum, or PHash.
func (c Client) FindSceneByFingerprints(ctx context.Context, fps models.Fingerprints) ([]*models.ScrapedScene, error) {
	res, err := c.FindScenesByFingerprints(ctx, []models.Fingerprints{fps})
	if len(res) > 0 {
		return res[0], err
	}
	return nil, err
}

// FindScenesByFingerprints queries stash-box for scenes using every
// scene's MD5/OSHASH checksum, or PHash, and returns results in the same order
// as the input slice.
func (c Client) FindScenesByFingerprints(ctx context.Context, fps []models.Fingerprints) ([][]*models.ScrapedScene, error) {
	var fingerprints [][]*graphql.FingerprintQueryInput

	for _, fp := range fps {
		fingerprints = append(fingerprints, convertFingerprints(fp))
	}

	return c.findScenesByFingerprints(ctx, fingerprints)
}

func convertFingerprints(fps models.Fingerprints) []*graphql.FingerprintQueryInput {
	var ret []*graphql.FingerprintQueryInput

	for _, f := range fps {
		var i = &graphql.FingerprintQueryInput{}
		switch f.Type {
		case models.FingerprintTypeMD5:
			i.Algorithm = graphql.FingerprintAlgorithmMd5
			i.Hash = f.String()
		case models.FingerprintTypeOshash:
			i.Algorithm = graphql.FingerprintAlgorithmOshash
			i.Hash = f.String()
		case models.FingerprintTypePhash:
			i.Algorithm = graphql.FingerprintAlgorithmPhash
			i.Hash = utils.PhashToString(f.Int64())
		default:
			continue
		}

		if !i.Algorithm.IsValid() {
			continue
		}

		ret = append(ret, i)
	}

	return ret
}

func (c Client) findScenesByFingerprints(ctx context.Context, scenes [][]*graphql.FingerprintQueryInput) ([][]*models.ScrapedScene, error) {
	var results [][]*models.ScrapedScene

	// filter out nils
	var validScenes [][]*graphql.FingerprintQueryInput
	for _, s := range scenes {
		if len(s) > 0 {
			validScenes = append(validScenes, s)
		}
	}

	var ignoredTags []string

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
			var sceneResults []*models.ScrapedScene
			for _, scene := range sceneFragments {
				ss, err := c.sceneFragmentToScrapedScene(ctx, scene)
				if err != nil {
					return nil, err
				}

				var thisIgnoredTags []string
				ss.Tags, thisIgnoredTags = scraper.FilterTags(c.excludeTagRE, ss.Tags)
				ignoredTags = sliceutil.AppendUniques(ignoredTags, thisIgnoredTags)

				sceneResults = append(sceneResults, ss)
			}
			results = append(results, sceneResults)
		}
	}

	scraper.LogIgnoredTags(ignoredTags)

	// repopulate the results to be the same order as the input
	ret := make([][]*models.ScrapedScene, len(scenes))
	upTo := 0

	for i, v := range scenes {
		if len(v) > 0 {
			ret[i] = results[upTo]
			upTo++
		}
	}

	return ret, nil
}

func (c Client) sceneFragmentToScrapedScene(ctx context.Context, s *graphql.SceneFragment) (*models.ScrapedScene, error) {
	stashID := s.ID

	ss := &models.ScrapedScene{
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
		ss.Image = getFirstImage(ctx, c.httpClient, s.Images)
	}

	ss.URLs = make([]string, len(s.Urls))
	for i, u := range s.Urls {
		ss.URLs[i] = u.URL
	}

	if s.Studio != nil {
		var err error
		ss.Studio, err = c.resolveStudio(ctx, s.Studio)
		if err != nil {
			return nil, err
		}
	}

	for _, p := range s.Performers {
		sp := performerFragmentToScrapedPerformer(*p.Performer)
		ss.Performers = append(ss.Performers, sp)
	}

	for _, t := range s.Tags {
		st := &models.ScrapedTag{
			Name:         t.Name,
			RemoteSiteID: &t.ID,
		}
		ss.Tags = append(ss.Tags, st)
	}

	return ss, nil
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

type SceneDraft struct {
	// Files, URLs, StashIDs must be loaded
	Scene *models.Scene
	// StashIDs must be loaded
	Performers []*models.Performer
	// StashIDs must be loaded
	Studio *models.Studio
	// StashIDs must be loaded
	Tags  []*models.Tag
	Cover []byte
}

func (c Client) SubmitSceneDraft(ctx context.Context, d SceneDraft) (*string, error) {
	draft := newSceneDraftInput(d, c.box.Endpoint)
	var image io.Reader

	if len(d.Cover) > 0 {
		image = bytes.NewReader(d.Cover)
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

func newSceneDraftInput(d SceneDraft, endpoint string) graphql.SceneDraftInput {
	scene := d.Scene

	draft := graphql.SceneDraftInput{}

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
	draft.Urls = scene.URLs.List()

	if scene.Date != nil {
		v := scene.Date.String()
		draft.Date = &v
	}

	if d.Studio != nil {
		studio := d.Studio

		studioDraft := graphql.DraftEntityInput{
			Name: studio.Name,
		}

		stashIDs := studio.StashIDs.List()
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

	for _, f := range scene.Files.List() {
		duration := f.Duration

		if duration != 0 {
			fingerprints = appendFingerprintsUnique(fingerprints, fileFingerprintsToInputGraphQL(f.Fingerprints, int(duration))...)
		}
	}
	draft.Fingerprints = fingerprints

	scenePerformers := d.Performers

	inputPerformers := []*graphql.DraftEntityInput{}
	for _, p := range scenePerformers {
		performerDraft := graphql.DraftEntityInput{
			Name: p.Name,
		}

		stashIDs := p.StashIDs.List()
		for _, stashID := range stashIDs {
			c := stashID
			if stashID.Endpoint == endpoint {
				performerDraft.ID = &c.StashID
				break
			}
		}

		inputPerformers = append(inputPerformers, &performerDraft)
	}
	draft.Performers = inputPerformers

	var tags []*graphql.DraftEntityInput
	sceneTags := d.Tags
	for _, tag := range sceneTags {
		tagDraft := graphql.DraftEntityInput{Name: tag.Name}

		stashIDs := tag.StashIDs.List()
		for _, stashID := range stashIDs {
			if stashID.Endpoint == endpoint {
				tagDraft.ID = &stashID.StashID
				break
			}
		}

		tags = append(tags, &tagDraft)
	}
	draft.Tags = tags

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

	return draft
}

func fileFingerprintsToInputGraphQL(fps models.Fingerprints, duration int) []*graphql.FingerprintInput {
	var ret []*graphql.FingerprintInput

	for _, f := range fps {
		var i = &graphql.FingerprintInput{
			Duration: duration,
		}
		switch f.Type {
		case models.FingerprintTypeMD5:
			i.Algorithm = graphql.FingerprintAlgorithmMd5
			i.Hash = f.String()
		case models.FingerprintTypeOshash:
			i.Algorithm = graphql.FingerprintAlgorithmOshash
			i.Hash = f.String()
		case models.FingerprintTypePhash:
			i.Algorithm = graphql.FingerprintAlgorithmPhash
			i.Hash = utils.PhashToString(f.Int64())
		default:
			continue
		}

		if !i.Algorithm.IsValid() {
			continue
		}

		ret = appendFingerprintUnique(ret, i)
	}

	return ret
}

func (c Client) SubmitFingerprints(ctx context.Context, scenes []*models.Scene) (bool, error) {
	endpoint := c.box.Endpoint

	var fingerprints []graphql.FingerprintSubmission

	for _, scene := range scenes {
		stashIDs := scene.StashIDs.List()
		sceneStashID := ""
		for _, stashID := range stashIDs {
			if stashID.Endpoint == endpoint {
				sceneStashID = stashID.StashID
			}
		}

		if sceneStashID == "" {
			continue
		}

		for _, f := range scene.Files.List() {
			duration := f.Duration

			if duration == 0 {
				continue
			}

			fps := fileFingerprintsToInputGraphQL(f.Fingerprints, int(duration))
			for _, fp := range fps {
				fingerprints = append(fingerprints, graphql.FingerprintSubmission{
					SceneID:     sceneStashID,
					Fingerprint: fp,
				})
			}
		}
	}

	return c.submitFingerprints(ctx, fingerprints)
}

func (c Client) submitFingerprints(ctx context.Context, fingerprints []graphql.FingerprintSubmission) (bool, error) {
	for _, fingerprint := range fingerprints {
		_, err := c.client.SubmitFingerprint(ctx, fingerprint)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func appendFingerprintUnique(v []*graphql.FingerprintInput, toAdd *graphql.FingerprintInput) []*graphql.FingerprintInput {
	for _, vv := range v {
		if vv.Algorithm == toAdd.Algorithm && vv.Hash == toAdd.Hash {
			return v
		}
	}

	return append(v, toAdd)
}

func appendFingerprintsUnique(v []*graphql.FingerprintInput, toAdd ...*graphql.FingerprintInput) []*graphql.FingerprintInput {
	for _, a := range toAdd {
		v = appendFingerprintUnique(v, a)
	}

	return v
}
