package api

import (
	"context"
	"errors"
	"sort"
	"strconv"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/scraper"
)

var (
	// ErrNotImplemented is an error which means the given functionality isn't implemented by the API.
	ErrNotImplemented = errors.New("not implemented")

	// ErrNotSupported is returned whenever there's a test, which can be used to guard against the error,
	// but the given parameters aren't supported by the system.
	ErrNotSupported = errors.New("not supported")

	// ErrInput signifies errors where the input isn't valid for some reason. And no more specific error exists.
	ErrInput = errors.New("input error")
)

type hookExecutor interface {
	ExecutePostHooks(ctx context.Context, id int, hookType plugin.HookTriggerEnum, input interface{}, inputFields []string)
}

type Resolver struct {
	txnManager   models.TransactionManager
	hookExecutor hookExecutor
}

func (r *Resolver) scraperCache() *scraper.Cache {
	return manager.GetInstance().ScraperCache
}

func (r *Resolver) Gallery() GalleryResolver {
	return &galleryResolver{r}
}
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Performer() PerformerResolver {
	return &performerResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Scene() SceneResolver {
	return &sceneResolver{r}
}
func (r *Resolver) Image() ImageResolver {
	return &imageResolver{r}
}
func (r *Resolver) SceneMarker() SceneMarkerResolver {
	return &sceneMarkerResolver{r}
}
func (r *Resolver) Studio() StudioResolver {
	return &studioResolver{r}
}
func (r *Resolver) Movie() MovieResolver {
	return &movieResolver{r}
}
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}
func (r *Resolver) Tag() TagResolver {
	return &tagResolver{r}
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

func (r *Resolver) withTxn(ctx context.Context, fn func(r models.Repository) error) error {
	return r.txnManager.WithTxn(ctx, fn)
}

func (r *Resolver) withReadTxn(ctx context.Context, fn func(r models.ReaderRepository) error) error {
	return r.txnManager.WithReadTxn(ctx, fn)
}

func (r *queryResolver) MarkerWall(ctx context.Context, q *string) (ret []*models.SceneMarker, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		ret, err = repo.SceneMarker().Wall(q)
		return err
	}); err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *queryResolver) SceneWall(ctx context.Context, q *string) (ret []*models.Scene, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		ret, err = repo.Scene().Wall(q)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) MarkerStrings(ctx context.Context, q *string, sort *string) (ret []*models.MarkerStringsResultType, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		ret, err = repo.SceneMarker().GetMarkerStrings(q, sort)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) Stats(ctx context.Context) (*StatsResultType, error) {
	var ret StatsResultType
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		scenesQB := repo.Scene()
		imageQB := repo.Image()
		galleryQB := repo.Gallery()
		studiosQB := repo.Studio()
		performersQB := repo.Performer()
		moviesQB := repo.Movie()
		tagsQB := repo.Tag()
		scenesCount, _ := scenesQB.Count()
		scenesSize, _ := scenesQB.Size()
		scenesDuration, _ := scenesQB.Duration()
		imageCount, _ := imageQB.Count()
		imageSize, _ := imageQB.Size()
		galleryCount, _ := galleryQB.Count()
		performersCount, _ := performersQB.Count()
		studiosCount, _ := studiosQB.Count()
		moviesCount, _ := moviesQB.Count()
		tagsCount, _ := tagsQB.Count()

		ret = StatsResultType{
			SceneCount:     scenesCount,
			ScenesSize:     scenesSize,
			ScenesDuration: scenesDuration,
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

func (r *queryResolver) Version(ctx context.Context) (*Version, error) {
	version, hash, buildtime := GetVersion()

	return &Version{
		Version:   &version,
		Hash:      hash,
		BuildTime: buildtime,
	}, nil
}

// Latestversion returns the latest git shorthash commit.
func (r *queryResolver) Latestversion(ctx context.Context) (*ShortVersion, error) {
	ver, url, err := GetLatestVersion(ctx, true)
	if err == nil {
		logger.Infof("Retrieved latest hash: %s", ver)
	} else {
		logger.Errorf("Error while retrieving latest hash: %s", err)
	}

	return &ShortVersion{
		Shorthash: ver,
		URL:       url,
	}, err
}

// Get scene marker tags which show up under the video.
func (r *queryResolver) SceneMarkerTags(ctx context.Context, scene_id string) ([]*SceneMarkerTag, error) {
	sceneID, err := strconv.Atoi(scene_id)
	if err != nil {
		return nil, err
	}

	var keys []int
	tags := make(map[int]*SceneMarkerTag)

	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		sceneMarkers, err := repo.SceneMarker().FindBySceneID(sceneID)
		if err != nil {
			return err
		}

		tqb := repo.Tag()
		for _, sceneMarker := range sceneMarkers {
			markerPrimaryTag, err := tqb.Find(sceneMarker.PrimaryTagID)
			if err != nil {
				return err
			}
			_, hasKey := tags[markerPrimaryTag.ID]
			if !hasKey {
				sceneMarkerTag := &SceneMarkerTag{Tag: markerPrimaryTag}
				tags[markerPrimaryTag.ID] = sceneMarkerTag
				keys = append(keys, markerPrimaryTag.ID)
			}
			tags[markerPrimaryTag.ID].SceneMarkers = append(tags[markerPrimaryTag.ID].SceneMarkers, sceneMarker)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// Sort so that primary tags that show up earlier in the video are first.
	sort.Slice(keys, func(i, j int) bool {
		a := tags[keys[i]]
		b := tags[keys[j]]
		return a.SceneMarkers[0].Seconds < b.SceneMarkers[0].Seconds
	})

	var result []*SceneMarkerTag
	for _, key := range keys {
		result = append(result, tags[key])
	}

	return result, nil
}
