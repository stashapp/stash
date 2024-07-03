package sqlite

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

type tagFilterHandler struct {
	tagFilter *models.TagFilterType
}

func (qb *tagFilterHandler) validate() error {
	tagFilter := qb.tagFilter
	if tagFilter == nil {
		return nil
	}

	if err := validateFilterCombination(tagFilter.OperatorFilter); err != nil {
		return err
	}

	if subFilter := tagFilter.SubFilter(); subFilter != nil {
		sqb := &tagFilterHandler{tagFilter: subFilter}
		if err := sqb.validate(); err != nil {
			return err
		}
	}

	return nil
}

func (qb *tagFilterHandler) handle(ctx context.Context, f *filterBuilder) {
	tagFilter := qb.tagFilter
	if tagFilter == nil {
		return
	}

	if err := qb.validate(); err != nil {
		f.setError(err)
		return
	}

	sf := tagFilter.SubFilter()
	if sf != nil {
		sub := &tagFilterHandler{sf}
		handleSubFilter(ctx, sub, f, tagFilter.OperatorFilter)
	}

	f.handleCriterion(ctx, qb.criterionHandler())
}

func (qb *tagFilterHandler) criterionHandler() criterionHandler {
	tagFilter := qb.tagFilter
	return compoundHandler{
		stringCriterionHandler(tagFilter.Name, tagTable+".name"),
		qb.aliasCriterionHandler(tagFilter.Aliases),

		boolCriterionHandler(tagFilter.Favorite, tagTable+".favorite", nil),
		stringCriterionHandler(tagFilter.Description, tagTable+".description"),
		boolCriterionHandler(tagFilter.IgnoreAutoTag, tagTable+".ignore_auto_tag", nil),

		qb.isMissingCriterionHandler(tagFilter.IsMissing),
		qb.sceneCountCriterionHandler(tagFilter.SceneCount),
		qb.imageCountCriterionHandler(tagFilter.ImageCount),
		qb.galleryCountCriterionHandler(tagFilter.GalleryCount),
		qb.performerCountCriterionHandler(tagFilter.PerformerCount),
		qb.studioCountCriterionHandler(tagFilter.StudioCount),

		qb.groupCountCriterionHandler(tagFilter.GroupCount),
		qb.groupCountCriterionHandler(tagFilter.MovieCount),

		qb.markerCountCriterionHandler(tagFilter.MarkerCount),
		qb.parentsCriterionHandler(tagFilter.Parents),
		qb.childrenCriterionHandler(tagFilter.Children),
		qb.parentCountCriterionHandler(tagFilter.ParentCount),
		qb.childCountCriterionHandler(tagFilter.ChildCount),
		&timestampCriterionHandler{tagFilter.CreatedAt, "tags.created_at", nil},
		&timestampCriterionHandler{tagFilter.UpdatedAt, "tags.updated_at", nil},

		&relatedFilterHandler{
			relatedIDCol:   "scenes_tags.scene_id",
			relatedRepo:    sceneRepository.repository,
			relatedHandler: &sceneFilterHandler{tagFilter.ScenesFilter},
			joinFn: func(f *filterBuilder) {
				tagRepository.scenes.innerJoin(f, "", "tags.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "images_tags.image_id",
			relatedRepo:    imageRepository.repository,
			relatedHandler: &imageFilterHandler{tagFilter.ImagesFilter},
			joinFn: func(f *filterBuilder) {
				tagRepository.images.innerJoin(f, "", "tags.id")
			},
		},

		&relatedFilterHandler{
			relatedIDCol:   "galleries_tags.gallery_id",
			relatedRepo:    galleryRepository.repository,
			relatedHandler: &galleryFilterHandler{tagFilter.GalleriesFilter},
			joinFn: func(f *filterBuilder) {
				tagRepository.galleries.innerJoin(f, "", "tags.id")
			},
		},
	}
}

func (qb *tagFilterHandler) aliasCriterionHandler(alias *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		primaryTable: tagTable,
		primaryFK:    tagIDColumn,
		joinTable:    tagAliasesTable,
		stringColumn: tagAliasColumn,
		addJoinTable: func(f *filterBuilder) {
			tagRepository.aliases.join(f, "", "tags.id")
		},
	}

	return h.handler(alias)
}

func (qb *tagFilterHandler) isMissingCriterionHandler(isMissing *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if isMissing != nil && *isMissing != "" {
			switch *isMissing {
			case "image":
				f.addWhere("tags.image_blob IS NULL")
			default:
				f.addWhere("(tags." + *isMissing + " IS NULL OR TRIM(tags." + *isMissing + ") = '')")
			}
		}
	}
}

func (qb *tagFilterHandler) sceneCountCriterionHandler(sceneCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if sceneCount != nil {
			f.addLeftJoin("scenes_tags", "", "scenes_tags.tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct scenes_tags.scene_id)", *sceneCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *tagFilterHandler) imageCountCriterionHandler(imageCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if imageCount != nil {
			f.addLeftJoin("images_tags", "", "images_tags.tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct images_tags.image_id)", *imageCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *tagFilterHandler) galleryCountCriterionHandler(galleryCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if galleryCount != nil {
			f.addLeftJoin("galleries_tags", "", "galleries_tags.tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct galleries_tags.gallery_id)", *galleryCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *tagFilterHandler) performerCountCriterionHandler(performerCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performerCount != nil {
			f.addLeftJoin("performers_tags", "", "performers_tags.tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct performers_tags.performer_id)", *performerCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *tagFilterHandler) studioCountCriterionHandler(studioCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if studioCount != nil {
			f.addLeftJoin("studios_tags", "", "studios_tags.tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct studios_tags.studio_id)", *studioCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *tagFilterHandler) groupCountCriterionHandler(movieCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if movieCount != nil {
			f.addLeftJoin("movies_tags", "", "movies_tags.tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct movies_tags.movie_id)", *movieCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *tagFilterHandler) markerCountCriterionHandler(markerCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if markerCount != nil {
			f.addLeftJoin("scene_markers_tags", "", "scene_markers_tags.tag_id = tags.id")
			f.addLeftJoin("scene_markers", "", "scene_markers_tags.scene_marker_id = scene_markers.id OR scene_markers.primary_tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct scene_markers.id)", *markerCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *tagFilterHandler) parentsCriterionHandler(criterion *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if criterion != nil {
			tags := criterion.CombineExcludes()

			// validate the modifier
			switch tags.Modifier {
			case models.CriterionModifierIncludesAll, models.CriterionModifierIncludes, models.CriterionModifierExcludes, models.CriterionModifierIsNull, models.CriterionModifierNotNull:
				// valid
			default:
				f.setError(fmt.Errorf("invalid modifier %s for tag parent/children", criterion.Modifier))
			}

			if tags.Modifier == models.CriterionModifierIsNull || tags.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if tags.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				f.addLeftJoin("tags_relations", "parent_relations", "tags.id = parent_relations.child_id")

				f.addWhere(fmt.Sprintf("parent_relations.parent_id IS %s NULL", notClause))
				return
			}

			if len(tags.Value) == 0 && len(tags.Excludes) == 0 {
				return
			}

			if len(tags.Value) > 0 {
				var args []interface{}
				for _, val := range tags.Value {
					args = append(args, val)
				}

				depthVal := 0
				if tags.Depth != nil {
					depthVal = *tags.Depth
				}

				var depthCondition string
				if depthVal != -1 {
					depthCondition = fmt.Sprintf("WHERE depth < %d", depthVal)
				}

				query := `parents AS (
		SELECT parent_id AS root_id, child_id AS item_id, 0 AS depth FROM tags_relations WHERE parent_id IN` + getInBinding(len(tags.Value)) + `
		UNION
		SELECT root_id, child_id, depth + 1 FROM tags_relations INNER JOIN parents ON item_id = parent_id ` + depthCondition + `
	)`

				f.addRecursiveWith(query, args...)

				f.addLeftJoin("parents", "", "parents.item_id = tags.id")

				addHierarchicalConditionClauses(f, tags, "parents", "root_id")
			}

			if len(tags.Excludes) > 0 {
				var args []interface{}
				for _, val := range tags.Excludes {
					args = append(args, val)
				}

				depthVal := 0
				if tags.Depth != nil {
					depthVal = *tags.Depth
				}

				var depthCondition string
				if depthVal != -1 {
					depthCondition = fmt.Sprintf("WHERE depth < %d", depthVal)
				}

				query := `parents2 AS (
		SELECT parent_id AS root_id, child_id AS item_id, 0 AS depth FROM tags_relations WHERE parent_id IN` + getInBinding(len(tags.Excludes)) + `
		UNION
		SELECT root_id, child_id, depth + 1 FROM tags_relations INNER JOIN parents2 ON item_id = parent_id ` + depthCondition + `
	)`

				f.addRecursiveWith(query, args...)

				f.addLeftJoin("parents2", "", "parents2.item_id = tags.id")

				addHierarchicalConditionClauses(f, models.HierarchicalMultiCriterionInput{
					Value:    tags.Excludes,
					Depth:    tags.Depth,
					Modifier: models.CriterionModifierExcludes,
				}, "parents2", "root_id")
			}
		}
	}
}

func (qb *tagFilterHandler) childrenCriterionHandler(criterion *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if criterion != nil {
			tags := criterion.CombineExcludes()

			// validate the modifier
			switch tags.Modifier {
			case models.CriterionModifierIncludesAll, models.CriterionModifierIncludes, models.CriterionModifierExcludes, models.CriterionModifierIsNull, models.CriterionModifierNotNull:
				// valid
			default:
				f.setError(fmt.Errorf("invalid modifier %s for tag parent/children", criterion.Modifier))
			}

			if tags.Modifier == models.CriterionModifierIsNull || tags.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if tags.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				f.addLeftJoin("tags_relations", "child_relations", "tags.id = child_relations.parent_id")

				f.addWhere(fmt.Sprintf("child_relations.child_id IS %s NULL", notClause))
				return
			}

			if len(tags.Value) == 0 && len(tags.Excludes) == 0 {
				return
			}

			if len(tags.Value) > 0 {
				var args []interface{}
				for _, val := range tags.Value {
					args = append(args, val)
				}

				depthVal := 0
				if tags.Depth != nil {
					depthVal = *tags.Depth
				}

				var depthCondition string
				if depthVal != -1 {
					depthCondition = fmt.Sprintf("WHERE depth < %d", depthVal)
				}

				query := `children AS (
		SELECT child_id AS root_id, parent_id AS item_id, 0 AS depth FROM tags_relations WHERE child_id IN` + getInBinding(len(tags.Value)) + `
		UNION
		SELECT root_id, parent_id, depth + 1 FROM tags_relations INNER JOIN children ON item_id = child_id ` + depthCondition + `
	)`

				f.addRecursiveWith(query, args...)

				f.addLeftJoin("children", "", "children.item_id = tags.id")

				addHierarchicalConditionClauses(f, tags, "children", "root_id")
			}

			if len(tags.Excludes) > 0 {
				var args []interface{}
				for _, val := range tags.Excludes {
					args = append(args, val)
				}

				depthVal := 0
				if tags.Depth != nil {
					depthVal = *tags.Depth
				}

				var depthCondition string
				if depthVal != -1 {
					depthCondition = fmt.Sprintf("WHERE depth < %d", depthVal)
				}

				query := `children2 AS (
		SELECT child_id AS root_id, parent_id AS item_id, 0 AS depth FROM tags_relations WHERE child_id IN` + getInBinding(len(tags.Excludes)) + `
		UNION
		SELECT root_id, parent_id, depth + 1 FROM tags_relations INNER JOIN children2 ON item_id = child_id ` + depthCondition + `
	)`

				f.addRecursiveWith(query, args...)

				f.addLeftJoin("children2", "", "children2.item_id = tags.id")

				addHierarchicalConditionClauses(f, models.HierarchicalMultiCriterionInput{
					Value:    tags.Excludes,
					Depth:    tags.Depth,
					Modifier: models.CriterionModifierExcludes,
				}, "children2", "root_id")
			}
		}
	}
}

func (qb *tagFilterHandler) parentCountCriterionHandler(parentCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if parentCount != nil {
			f.addLeftJoin("tags_relations", "parents_count", "parents_count.child_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct parents_count.parent_id)", *parentCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *tagFilterHandler) childCountCriterionHandler(childCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if childCount != nil {
			f.addLeftJoin("tags_relations", "children_count", "children_count.parent_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct children_count.child_id)", *childCount)

			f.addHaving(clause, args...)
		}
	}
}
