package api

import (
	"context"
	"github.com/stashapp/stash/internal/models"
	"github.com/stashapp/stash/internal/scraper"
	"sort"
	"strconv"
)

type Resolver struct{}

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
func (r *Resolver) SceneMarker() models.SceneMarkerResolver {
	return &sceneMarkerResolver{r}
}
func (r *Resolver) Studio() models.StudioResolver {
	return &studioResolver{r}
}
func (r *Resolver) Subscription() models.SubscriptionResolver {
	return &subscriptionResolver{r}
}
func (r *Resolver) Tag() models.TagResolver {
	return &tagResolver{r}
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

type galleryResolver struct{ *Resolver }
type performerResolver struct{ *Resolver }
type sceneResolver struct{ *Resolver }
type sceneMarkerResolver struct{ *Resolver }
type studioResolver struct{ *Resolver }
type tagResolver struct{ *Resolver }

func (r *queryResolver) MarkerWall(ctx context.Context, q *string) ([]models.SceneMarker, error) {
	qb := models.NewSceneMarkerQueryBuilder()
	return qb.Wall(q)
}

func (r *queryResolver) SceneWall(ctx context.Context, q *string) ([]models.Scene, error) {
	qb := models.NewSceneQueryBuilder()
	return qb.Wall(q)
}

func (r *queryResolver) MarkerStrings(ctx context.Context, q *string, sort *string) ([]*models.MarkerStringsResultType, error) {
	qb := models.NewSceneMarkerQueryBuilder()
	return qb.GetMarkerStrings(q, sort)
}

func (r *queryResolver) ValidGalleriesForScene(ctx context.Context, scene_id *string) ([]models.Gallery, error) {
	if scene_id == nil {
		panic("nil scene id") // TODO make scene_id mandatory
	}
	sceneID, _ := strconv.Atoi(*scene_id)
	sqb := models.NewSceneQueryBuilder()
	scene, err := sqb.Find(sceneID)
	if err != nil {
		return nil, err
	}

	qb := models.NewGalleryQueryBuilder()
	validGalleries, err := qb.ValidGalleriesForScenePath(scene.Path)
	sceneGallery, _ := qb.FindBySceneID(sceneID, nil)
	if sceneGallery != nil {
		validGalleries = append(validGalleries, *sceneGallery)
	}
	return validGalleries, nil
}

func (r *queryResolver) Stats(ctx context.Context) (models.StatsResultType, error) {
	//scenesCount, _ := runCountQuery(buildCountQuery(selectAll("scenes")), nil)
	//galleryCount, _ := runCountQuery(buildCountQuery(selectAll("galleries")), nil)
	//performersCount, _ := runCountQuery(buildCountQuery(selectAll("performers")), nil)
	//studiosCount, _ := runCountQuery(buildCountQuery(selectAll("studios")), nil)
	//tagsCount, _ := runCountQuery(buildCountQuery(selectAll("tags")), nil)
	//return StatsResultType{
	//	SceneCount: scenesCount,
	//	GalleryCount: galleryCount,
	//	PerformerCount: performersCount,
	//	StudioCount: studiosCount,
	//	TagCount: tagsCount,
	//}, nil
	return models.StatsResultType{}, nil // TODO
}

// Get scene marker tags which show up under the video.
func (r *queryResolver) SceneMarkerTags(ctx context.Context, scene_id string) ([]models.SceneMarkerTag, error) {
	sceneID, _ := strconv.Atoi(scene_id)
	sqb := models.NewSceneMarkerQueryBuilder()
	sceneMarkers, err := sqb.FindBySceneID(sceneID, nil)
	if err != nil {
		return nil, err
	}

	tags := make(map[int]*models.SceneMarkerTag)
	var keys []int
	tqb := models.NewTagQueryBuilder()
	for _, sceneMarker := range sceneMarkers {
		if !sceneMarker.PrimaryTagID.Valid {
			panic("missing primary tag id")
		}
		markerPrimaryTag, err := tqb.Find(int(sceneMarker.PrimaryTagID.Int64), nil)
		if err != nil {
			return nil, err
		}
		_, hasKey := tags[markerPrimaryTag.ID]
		var sceneMarkerTag *models.SceneMarkerTag
		if !hasKey {
			sceneMarkerTag = &models.SceneMarkerTag{ Tag: *markerPrimaryTag }
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

	var result []models.SceneMarkerTag
	for _, key := range keys {
		result = append(result, *tags[key])
	}

	return result, nil
}

func (r *queryResolver) ScrapeFreeones(ctx context.Context, performer_name string) (*models.ScrapedPerformer, error) {
	return scraper.GetPerformer(performer_name)
}

func (r *queryResolver) ScrapeFreeonesPerformerList(ctx context.Context, query string) ([]string, error) {
	return scraper.GetPerformerNames(query)
}
