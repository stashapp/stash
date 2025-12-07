# GraphQL Frontend Patches

This document describes the GraphQL schema and query modifications for the frontend that need to be reapplied after upgrading.

## Overview

| File | Change | Lines | Type |
|------|--------|-------|------|
| `data/facets.graphql` | Facet count fragments & queries | +177 | New file |
| `data/gallery.graphql` | Added `image_count` field | +1 | Modification |
| `data/group-slim.graphql` | Added counts + FilterGroupData | +11 | Modification |
| `data/performer-slim.graphql` | Added counts + FilterPerformerData | +12 | Modification |
| `data/studio.graphql` | Added counts + FilterStudioData | +12 | Modification |
| `data/tag.graphql` | Added counts + FilterTagData | +12 | Modification |
| `queries/movie.graphql` | Added FindGroupsForFilter | +13 | Modification |
| `queries/performer.graphql` | Added FindPerformersForFilter | +13 | Modification |
| `queries/scene.graphql` | Added SceneRecommendations queries | +24 | Modification |
| `queries/studio.graphql` | Added FindStudiosForFilter | +13 | Modification |
| `queries/tag.graphql` | Added FindTagsForFilter | +13 | Modification |

---

## New Files

### facets.graphql

**File:** `graphql/data/facets.graphql`

This is a **new file** containing all facet-related fragments and queries for sidebar filter counts.

**Fragments:**
- `FacetCountData` - Generic facet count
- `BooleanFacetCountData` - Boolean facets
- `ResolutionFacetCountData` - Resolution facets
- `OrientationFacetCountData` - Orientation facets
- `GenderFacetCountData` - Gender facets
- `RatingFacetCountData` - Rating facets
- `CircumcisedFacetCountData` - Circumcised facets
- `CaptionFacetCountData` - Caption/language facets

**Queries:**
- `SceneFacets` - Get facet counts for scenes
- `PerformerFacets` - Get facet counts for performers
- `GalleryFacets` - Get facet counts for galleries
- `ImageFacets` - Get facet counts for images
- `GroupFacets` - Get facet counts for groups
- `StudioFacets` - Get facet counts for studios
- `TagFacets` - Get facet counts for tags

---

## Data Fragment Modifications

### gallery.graphql

**File:** `graphql/data/gallery.graphql`

```diff
 fragment SelectGalleryData on Gallery {
   ...
+  image_count
 }
```

### group-slim.graphql

**File:** `graphql/data/group-slim.graphql`

```diff
 fragment SelectGroupData on Group {
   ...
   front_image_path
+  scene_count
+  scene_count_all: scene_count(depth: -1)
+  performer_count
+  sub_group_count
 }
+
+# Minimal fragment for sidebar filter - only essential fields for fast loading
+fragment FilterGroupData on Group {
+  id
+  name
+  aliases
+}
```

### performer-slim.graphql

**File:** `graphql/data/performer-slim.graphql`

```diff
 fragment SelectPerformerData on Performer {
   ...
   death_date
+  scene_count
+  image_count
+  gallery_count
+  group_count
 }
+
+# Minimal fragment for sidebar filter - only essential fields for fast loading
+fragment FilterPerformerData on Performer {
+  id
+  name
+  disambiguation
+  alias_list
+}
```

### studio.graphql

**File:** `graphql/data/studio.graphql`

```diff
 fragment SelectStudioData on Studio {
   ...
   image_path
+  scene_count
+  scene_count_all: scene_count(depth: -1)
+  image_count
+  gallery_count
+  performer_count
   ...
 }
+
+# Minimal fragment for sidebar filter - only essential fields for fast loading
+fragment FilterStudioData on Studio {
+  id
+  name
+  aliases
+}
```

### tag.graphql

**File:** `graphql/data/tag.graphql`

```diff
 fragment SelectTagData on Tag {
   ...
   image_path
+  scene_count
+  scene_count_all: scene_count(depth: -1)
+  image_count
+  gallery_count
+  performer_count
   ...
 }
+
+# Minimal fragment for sidebar filter - only essential fields for fast loading
+fragment FilterTagData on Tag {
+  id
+  name
+  aliases
+}
```

---

## Query Modifications

### movie.graphql (Groups)

**File:** `graphql/queries/movie.graphql`

```diff
+# Lightweight query for sidebar filter - minimal data for fast loading
+query FindGroupsForFilter(
+  $filter: FindFilterType
+  $group_filter: GroupFilterType
+) {
+  findGroups(filter: $filter, group_filter: $group_filter) {
+    count
+    groups {
+      ...FilterGroupData
+    }
+  }
+}
```

### performer.graphql

**File:** `graphql/queries/performer.graphql`

```diff
+# Lightweight query for sidebar filter - minimal data for fast loading
+query FindPerformersForFilter(
+  $filter: FindFilterType
+  $performer_filter: PerformerFilterType
+) {
+  findPerformers(filter: $filter, performer_filter: $performer_filter) {
+    count
+    performers {
+      ...FilterPerformerData
+    }
+  }
+}
```

### scene.graphql

**File:** `graphql/queries/scene.graphql`

```diff
+query SceneRecommendations($limit: Int) {
+  sceneRecommendations(limit: $limit) {
+    recommendations {
+      score
+      reasons
+      scene {
+        ...SlimSceneData
+      }
+    }
+  }
+}
+
+query SceneRecommendationsForScene($sceneId: ID!, $limit: Int) {
+  sceneRecommendationsForScene(scene_id: $sceneId, limit: $limit) {
+    recommendations {
+      score
+      reasons
+      scene {
+        ...SlimSceneData
+      }
+    }
+  }
+}
```

### studio.graphql

**File:** `graphql/queries/studio.graphql`

```diff
+# Lightweight query for sidebar filter - minimal data for fast loading
+query FindStudiosForFilter(
+  $filter: FindFilterType
+  $studio_filter: StudioFilterType
+) {
+  findStudios(filter: $filter, studio_filter: $studio_filter) {
+    count
+    studios {
+      ...FilterStudioData
+    }
+  }
+}
```

### tag.graphql

**File:** `graphql/queries/tag.graphql`

```diff
+# Lightweight query for sidebar filter - minimal data for fast loading
+query FindTagsForFilter(
+  $filter: FindFilterType
+  $tag_filter: TagFilterType
+) {
+  findTags(filter: $filter, tag_filter: $tag_filter) {
+    count
+    tags {
+      ...FilterTagData
+    }
+  }
+}
```

---

## Backend Requirements

These frontend GraphQL changes require corresponding **backend schema changes**:

1. **Facet types** must be defined in the backend schema
2. **SceneRecommendations** resolver must exist
3. **Count fields** (scene_count, image_count, etc.) must be available on entities

If upgrading from a version without these backend features, the frontend queries will fail until the backend is updated.

---

## Application Instructions

After upgrading upstream:

1. **Check backend compatibility** - Ensure backend schema supports facets and recommendations
2. **Create facets.graphql** - This is a new file
3. **Apply fragment additions** - Add count fields and Filter*Data fragments
4. **Apply query additions** - Add For Filter and Recommendations queries
5. **Run codegen** - Regenerate TypeScript types: `yarn generate`
6. **Test** - Verify facet counts and recommendations work

## Related Files

- `extensions/hooks/useFacetCounts.ts` - Uses facet queries
- `extensions/filters/` - Uses Filter*Data fragments
- `extensions/components/AISceneRecommendationRow.tsx` - Uses SceneRecommendations
- `core/StashService.ts` - May have additional query exports

