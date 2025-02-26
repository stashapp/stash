package stashbox

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stashapp/stash/pkg/scraper/stashbox/graphql"
	"github.com/stashapp/stash/pkg/utils"
)

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

func appendFingerprintUnique(v []*graphql.FingerprintInput, toAdd *graphql.FingerprintInput) []*graphql.FingerprintInput {
	for _, vv := range v {
		if vv.Algorithm == toAdd.Algorithm && vv.Hash == toAdd.Hash {
			return v
		}
	}

	return append(v, toAdd)
}
