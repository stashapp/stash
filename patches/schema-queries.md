# Schema Query Additions

After merging from upstream, add these queries to `graphql/schema/schema.graphql`.

## Location

Add after the `findTags` query (around line 148) and before `markerWall`:

## Queries to Add

```graphql
  # Facets - aggregated counts for filtering
  """
  Get facet counts for scenes based on current filter.
  Returns counts for each filter dimension, allowing the UI to show
  how many results each filter option would return.
  """
  sceneFacets(
    "Scene filter to apply when calculating facets"
    scene_filter: SceneFilterType
    "Maximum facets to return per category (default: 100)"
    limit: Int
    "Include performer_tags facet (expensive, default: false)"
    include_performer_tags: Boolean
    "Include captions facet (expensive, default: false)"
    include_captions: Boolean
  ): SceneFacetsResult!

  """
  Get facet counts for performers based on current filter.
  """
  performerFacets(
    performer_filter: PerformerFilterType
    limit: Int
  ): PerformerFacetsResult!

  """
  Get facet counts for galleries based on current filter.
  """
  galleryFacets(
    gallery_filter: GalleryFilterType
    limit: Int
  ): GalleryFacetsResult!

  """
  Get facet counts for groups based on current filter.
  """
  groupFacets(
    group_filter: GroupFilterType
    limit: Int
  ): GroupFacetsResult!

  """
  Get facet counts for studios based on current filter.
  """
  studioFacets(
    studio_filter: StudioFilterType
    limit: Int
  ): StudioFacetsResult!

  """
  Get facet counts for tags based on current filter.
  """
  tagFacets(
    tag_filter: TagFilterType
    limit: Int
  ): TagFacetsResult!
```

## Required Type Definitions

The result types are defined in `graphql/schema/types/facets.graphql` (a new file that won't conflict).

---

## Recommendations Queries

Add after `findScenes` query (around line 43) and before `findDuplicateScenes`:

```graphql
  "Get scene recommendations based on user interest"
  sceneRecommendations(limit: Int): SceneRecommendationsResultType!

  "Get scene recommendations based on a specific scene"
  sceneRecommendationsForScene(
    scene_id: ID!
    limit: Int
  ): SceneRecommendationsResultType!
```

Add after `findPerformers` query (around line 100) and before `findStudios`:

```graphql
  "Get performer recommendations based on user interest"
  performerRecommendations(limit: Int): PerformerRecommendationsResultType!

  "Get performer recommendations based on a specific performer"
  performerRecommendationsForPerformer(
    performer_id: ID!
    limit: Int
  ): PerformerRecommendationsResultType!
```

### Recommendations Result Types

These are defined in resolver files (no separate schema file):
- `internal/api/resolver_scene_recommendations_result_type.go`
- `internal/api/resolver_performer_recommendations_result_type.go`

### Recommendations Implementation Files

New package (won't conflict with upstream):
```
pkg/recommendation/
├── performer.go    # 525 lines - performer recommendation logic
└── scene.go        # 903 lines - scene recommendation logic
```

API resolvers:
```
internal/api/
├── resolver_query_performer_recommendations.go
├── resolver_query_scene_recommendations.go
├── resolver_performer_recommendations_result_type.go
└── resolver_scene_recommendations_result_type.go
```

---

## After Adding All Queries

Regenerate the GraphQL code:

```bash
go generate ./...
```

