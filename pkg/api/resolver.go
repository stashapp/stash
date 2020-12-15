package api

import (
	"context"
	"sort"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
)

type Resolver struct {
	txnManager models.TransactionManager
}

func (r *Resolver) Gallery() models.GalleryResolver {
	return &galleryResolver{r}
}
func (r *Resolver) Mutation() models.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Performer() models.PerformerResolver {
	return &performerResolver{r}
}
func (r *Resolver) Query() models.QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Scene() models.SceneResolver {
	return &sceneResolver{r}
}
func (r *Resolver) Image() models.ImageResolver {
	return &imageResolver{r}
}
func (r *Resolver) SceneMarker() models.SceneMarkerResolver {
	return &sceneMarkerResolver{r}
}
func (r *Resolver) Studio() models.StudioResolver {
	return &studioResolver{r}
}
func (r *Resolver) Movie() models.MovieResolver {
	return &movieResolver{r}
}
func (r *Resolver) Subscription() models.SubscriptionResolver {
	return &subscriptionResolver{r}
}
func (r *Resolver) Tag() models.TagResolver {
	return &tagResolver{r}
}

func (r *Resolver) ScrapedSceneTag() models.ScrapedSceneTagResolver {
	return &scrapedSceneTagResolver{r}
}

func (r *Resolver) ScrapedSceneMovie() models.ScrapedSceneMovieResolver {
	return &scrapedSceneMovieResolver{r}
}

func (r *Resolver) ScrapedScenePerformer() models.ScrapedScenePerformerResolver {
	return &scrapedScenePerformerResolver{r}
}

func (r *Resolver) ScrapedSceneStudio() models.ScrapedSceneStudioResolver {
	return &scrapedSceneStudioResolver{r}
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

type galleryResolver struct{ *Resolver }
type performerResolver struct{ *Resolver }
type sceneResolver struct{ *Resolver }
type sceneMarkerResolver struct{ *Resolver }
type imageResolver struct{ *Resolver }
type studioResolver struct{ *Resolver }
type movieResolver struct{ *Resolver }
type tagResolver struct{ *Resolver }
type scrapedSceneTagResolver struct{ *Resolver }
type scrapedSceneMovieResolver struct{ *Resolver }
type scrapedScenePerformerResolver struct{ *Resolver }
type scrapedSceneStudioResolver struct{ *Resolver }

func (r *Resolver) withTxn(ctx context.Context, fn func(r models.Repository) error) error {
	return r.txnManager.WithTxn(ctx, fn)
}

func (r *Resolver) withReadTxn(ctx context.Context, fn func(r models.ReaderRepository) error) error {
	return r.txnManager.WithReadTxn(ctx, fn)
}

func (r *queryResolver) MarkerWall(ctx context.Context, q *string) ([]*models.SceneMarker, error) {
	qb := sqlite.NewSceneMarkerQueryBuilder()
	return qb.Wall(q)
}

func (r *queryResolver) SceneWall(ctx context.Context, q *string) ([]*models.Scene, error) {
	qb := sqlite.NewSceneQueryBuilder()
	return qb.Wall(q)
}

func (r *queryResolver) MarkerStrings(ctx context.Context, q *string, sort *string) ([]*models.MarkerStringsResultType, error) {
	qb := sqlite.NewSceneMarkerQueryBuilder()
	return qb.GetMarkerStrings(q, sort)
}

func (r *queryResolver) ValidGalleriesForScene(ctx context.Context, scene_id *string) ([]*models.Gallery, error) {
	if scene_id == nil {
		panic("nil scene id") // TODO make scene_id mandatory
	}
	sceneID, err := strconv.Atoi(*scene_id)
	if err != nil {
		return nil, err
	}

	var validGalleries []*models.Gallery
	var sceneGallery *models.Gallery
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		sqb := sqlite.NewSceneQueryBuilder()
		scene, err := sqb.Find(sceneID)
		if err != nil {
			return err
		}

		qb := repo.Gallery()
		validGalleries, err = qb.ValidGalleriesForScenePath(scene.Path)
		if err != nil {
			return err
		}
		sceneGallery, err = qb.FindBySceneID(sceneID)
		return err
	}); err != nil {
		return nil, err
	}

	if sceneGallery != nil {
		validGalleries = append(validGalleries, sceneGallery)
	}
	return validGalleries, nil
}

func (r *queryResolver) Stats(ctx context.Context) (*models.StatsResultType, error) {
	var ret models.StatsResultType
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		scenesQB := sqlite.NewSceneQueryBuilder()
		imageQB := repo.Image()
		galleryQB := repo.Gallery()
		studiosQB := sqlite.NewStudioQueryBuilder()
		performersQB := repo.Performer()
		moviesQB := repo.Movie()
		tagsQB := sqlite.NewTagQueryBuilder()
		scenesCount, _ := scenesQB.Count()
		scenesSize, _ := scenesQB.Size()
		imageCount, _ := imageQB.Count()
		imageSize, _ := imageQB.Size()
		galleryCount, _ := galleryQB.Count()
		performersCount, _ := performersQB.Count()
		studiosCount, _ := studiosQB.Count()
		moviesCount, _ := moviesQB.Count()
		tagsCount, _ := tagsQB.Count()

		ret = models.StatsResultType{
			SceneCount:     scenesCount,
			ScenesSize:     scenesSize,
			ImageCount:     imageCount,
			ImagesSize:     imageSize,
			GalleryCount:   galleryCount,
			PerformerCount: performersCount,
			StudioCount:    studiosCount,
			MovieCount:     moviesCount,
			TagCount:       tagsCount,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (r *queryResolver) Version(ctx context.Context) (*models.Version, error) {
	version, hash, buildtime := GetVersion()

	return &models.Version{
		Version:   &version,
		Hash:      hash,
		BuildTime: buildtime,
	}, nil
}

//Gets latest version (git shorthash commit for now)
func (r *queryResolver) Latestversion(ctx context.Context) (*models.ShortVersion, error) {
	ver, url, err := GetLatestVersion(true)
	if err == nil {
		logger.Infof("Retrieved latest hash: %s", ver)
	} else {
		logger.Errorf("Error while retrieving latest hash: %s", err)
	}

	return &models.ShortVersion{
		Shorthash: ver,
		URL:       url,
	}, err
}

// Get scene marker tags which show up under the video.
func (r *queryResolver) SceneMarkerTags(ctx context.Context, scene_id string) ([]*models.SceneMarkerTag, error) {
	sceneID, _ := strconv.Atoi(scene_id)
	sqb := sqlite.NewSceneMarkerQueryBuilder()
	sceneMarkers, err := sqb.FindBySceneID(sceneID, nil)
	if err != nil {
		return nil, err
	}

	tags := make(map[int]*models.SceneMarkerTag)
	var keys []int
	tqb := sqlite.NewTagQueryBuilder()
	for _, sceneMarker := range sceneMarkers {
		markerPrimaryTag, err := tqb.Find(sceneMarker.PrimaryTagID, nil)
		if err != nil {
			return nil, err
		}
		_, hasKey := tags[markerPrimaryTag.ID]
		var sceneMarkerTag *models.SceneMarkerTag
		if !hasKey {
			sceneMarkerTag = &models.SceneMarkerTag{Tag: markerPrimaryTag}
			tags[markerPrimaryTag.ID] = sceneMarkerTag
			keys = append(keys, markerPrimaryTag.ID)
		} else {
			sceneMarkerTag = tags[markerPrimaryTag.ID]
		}
		tags[markerPrimaryTag.ID].SceneMarkers = append(tags[markerPrimaryTag.ID].SceneMarkers, sceneMarker)
	}

	// Sort so that primary tags that show up earlier in the video are first.
	sort.Slice(keys, func(i, j int) bool {
		a := tags[keys[i]]
		b := tags[keys[j]]
		return a.SceneMarkers[0].Seconds < b.SceneMarkers[0].Seconds
	})

	var result []*models.SceneMarkerTag
	for _, key := range keys {
		result = append(result, tags[key])
	}

	return result, nil
}

// wasFieldIncluded returns true if the given field was included in the request.
// Slices are unmarshalled to empty slices even if the field was omitted. This
// method determines if it was omitted altogether.
func wasFieldIncluded(ctx context.Context, field string) bool {
	rctx := graphql.GetRequestContext(ctx)

	_, ret := rctx.Variables[field]
	return ret
}
