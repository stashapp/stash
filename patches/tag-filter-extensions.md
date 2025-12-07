# Tag Filter Extensions

This fork adds `performers_filter` and `groups_filter` to the TagFilterType, allowing filtering tags by their related performers and groups.

## Files Modified

| File | Change |
|------|--------|
| `graphql/schema/types/filters.graphql` | Add filter fields to TagFilterType |
| `pkg/models/tag.go` | Add Go struct fields |
| `pkg/sqlite/tag.go` | Add join repositories |
| `pkg/sqlite/tag_filter.go` | Add filter handlers |

---

## 1. GraphQL Schema

**File:** `graphql/schema/types/filters.graphql`

Find `input TagFilterType` and add these fields after `galleries_filter`:

```graphql
input TagFilterType {
  # ... existing fields ...
  
  "Filter by related galleries that meet this criteria"
  galleries_filter: GalleryFilterType
  
  # ADD THESE TWO FIELDS:
  "Filter by related performers that meet this criteria"
  performers_filter: PerformerFilterType
  "Filter by related groups that meet this criteria"
  groups_filter: GroupFilterType
  
  "Filter by creation time"
  created_at: TimestampCriterionInput
  # ... rest of fields ...
}
```

---

## 2. Go Model

**File:** `pkg/models/tag.go`

Find `type TagFilterType struct` and add these fields after `GalleriesFilter`:

```go
type TagFilterType struct {
	// ... existing fields ...
	
	// Filter by related galleries that meet this criteria
	GalleriesFilter *GalleryFilterType `json:"galleries_filter"`
	
	// ADD THESE TWO FIELDS:
	// Filter by related performers that meet this criteria
	PerformersFilter *PerformerFilterType `json:"performers_filter"`
	// Filter by related groups that meet this criteria
	GroupsFilter *GroupFilterType `json:"groups_filter"`
	
	// Filter by created at
	CreatedAt *TimestampCriterionInput `json:"created_at"`
	// ... rest of fields ...
}
```

---

## 3. SQLite Repository

**File:** `pkg/sqlite/tag.go`

### 3a. Add join repository fields

Find `type tagRepositoryType struct` and add `performers` and `groups`:

```go
type tagRepositoryType struct {
	repository
	tableMgr *table

	aliases stringRepository

	scenes     joinRepository
	images     joinRepository
	galleries  joinRepository
	performers joinRepository  // ADD
	groups     joinRepository  // ADD
}
```

### 3b. Initialize the join repositories

Find the `tagRepository` initialization and add:

```go
var (
	tagRepository = tagRepositoryType{
		// ... existing fields ...
		
		galleries: joinRepository{
			// ... existing ...
		},
		
		// ADD THESE:
		performers: joinRepository{
			repository: repository{
				tableName: performersTagsTable,
				idColumn:  tagIDColumn,
			},
			fkColumn:     performerIDColumn,
			foreignTable: performerTable,
		},
		groups: joinRepository{
			repository: repository{
				tableName: groupsTagsTable,
				idColumn:  tagIDColumn,
			},
			fkColumn:     groupIDColumn,
			foreignTable: groupTable,
		},
	}
)
```

---

## 4. SQLite Filter Handler

**File:** `pkg/sqlite/tag_filter.go`

Find `func (qb *tagFilterHandler) criterionHandler()` and add these handlers after the `galleries_filter` handler:

```go
func (qb *tagFilterHandler) criterionHandler() criterionHandler {
	tagFilter := qb.tagFilter
	return criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
		// ... existing handlers ...
		
		// galleries_filter handler
		&relatedFilterHandler{
			// ... existing ...
		},
		
		// ADD THESE TWO HANDLERS:
		&relatedFilterHandler{
			relatedIDCol:   "performers_tags.performer_id",
			relatedRepo:    performerRepository.repository,
			relatedHandler: &performerFilterHandler{tagFilter.PerformersFilter},
			joinFn: func(f *filterBuilder) {
				tagRepository.performers.innerJoin(f, "", "tags.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "groups_tags.group_id",
			relatedRepo:    groupRepository.repository,
			relatedHandler: &groupFilterHandler{tagFilter.GroupsFilter},
			joinFn: func(f *filterBuilder) {
				tagRepository.groups.innerJoin(f, "", "tags.id")
			},
		},
	})
}
```

---

## 5. Regenerate

After making all changes:

```bash
go generate ./...
go build ./...
```

---

## Purpose

These filters enable queries like:
- Find tags used by performers from a specific country
- Find tags used in groups from a specific studio
- Filter tags based on complex performer/group criteria

